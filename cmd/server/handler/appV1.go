package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/macaron.v1"

	"github.com/liangchenye/update-service/keymanager"
	"github.com/liangchenye/update-service/service"
	"github.com/liangchenye/update-service/storage"
	"github.com/liangchenye/update-service/utils"
)

type httpListRet struct {
	Message string
	Content interface{}
}

//TODO: better http return result
func httpRet(head string, content interface{}, err error) (int, []byte) {
	var ret httpListRet
	var code int

	if err != nil {
		ret.Message = head + " fail"
		ret.Content = err.Error()
		code = http.StatusBadRequest
	} else {
		ret.Message = head
		ret.Content = content
		code = http.StatusOK
	}

	result, _ := json.Marshal(ret)
	return code, result
}

// AppListFileV1Handler lists  all the files in the namespace/repository
func AppListFileV1Handler(ctx *macaron.Context) (int, []byte) {
	namespace := ctx.Params(":namespace")
	repository := ctx.Params(":repository")

	us, _ := service.DefaultUpdateService("app", "v1", namespace, repository)
	apps, err := us.List()

	return httpRet("AppV1 List files", apps, err)
}

// AppGetPublicKeyV1Handler gets the public key of a appliance
func AppGetPublicKeyV1Handler(ctx *macaron.Context) (int, []byte) {
	namespace := ctx.Params(":namespace")
	a := utils.Appliance{Proto: "app", Version: "v1", Namespace: namespace}
	km, _ := keymanager.DefaultKeyManager()
	data, err := km.GetPublicKey(a)
	if err == nil {
		return http.StatusOK, data
	}

	return httpRet("AppV1 Get Public Key", nil, err)
}

// AppGetMetaV1Handler gets the meta data of all the namespace/repository
func AppGetMetaV1Handler(ctx *macaron.Context) (int, []byte) {
	namespace := ctx.Params(":namespace")
	repository := ctx.Params(":repository")

	us, _ := service.DefaultUpdateService("app", "v1", namespace, repository)
	data, err := us.GetMeta()
	if err == nil {
		return http.StatusOK, data
	}

	return httpRet("AppV1 Get Meta", nil, err)
}

// AppGetMetaSignV1Handler gets the meta signature data of all the namespace/repository
func AppGetMetaSignV1Handler(ctx *macaron.Context) (int, []byte) {
	namespace := ctx.Params(":namespace")
	repository := ctx.Params(":repository")

	us, _ := service.DefaultUpdateService("app", "v1", namespace, repository)
	data, err := us.GetMetaSign()
	if err != nil {
		return httpRet("AppV1 Get Meta Sign", data, err)
	}

	return http.StatusOK, data
}

// AppGetFileV1Handler gets the content of a certain app
func AppGetFileV1Handler(ctx *macaron.Context) (int, []byte) {
	namespace := ctx.Params(":namespace")
	repository := ctx.Params(":repository")
	name := ctx.Params(":name")

	key := fmt.Sprintf("%s/%s/%s", namespace, repository, name)
	store, _ := storage.DefaultUpdateServiceStorage()
	data, err := store.Get(key)
	if err != nil {
		return httpRet("AppV1 Get File", data, err)
	}

	return http.StatusOK, data
}

// AppPutFileV1Handler posts the content of a certain app
func AppPutFileV1Handler(ctx *macaron.Context) (int, []byte) {
	namespace := ctx.Params(":namespace")
	repository := ctx.Params(":repository")
	name := ctx.Params(":name")

	data, _ := ctx.Req.Body().Bytes()
	key := fmt.Sprintf("%s/%s/%s", namespace, repository, name)
	store, _ := storage.DefaultUpdateServiceStorage()
	err := store.Put(key, data)
	if err != nil {
		return httpRet("AppV1 Put data", nil, err)
	}

	sha, err := utils.SHA512(data)
	us, _ := service.DefaultUpdateService("app", "v1", namespace, repository)
	item, _ := service.NewUpdateServiceItem(name, []string{sha})
	err = us.Put(item)
	if err != nil {
		// remove the blob data either
		store.Delete(key)
		return httpRet("AppV1 Put data", nil, err)
	}

	return httpRet("AppV1 Put File", nil, nil)
}
