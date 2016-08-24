package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/liangchenye/update-service/keymanager"
	"github.com/liangchenye/update-service/storage"
	"github.com/liangchenye/update-service/utils"
)

const (
	defaultMetaFileName     = "meta.json"
	defaultMetaSignFileName = "meta.sign"
)

// UpdateService represents the meta info of a repository
type UpdateService struct {
	Proto      string
	Version    string
	Namespace  string
	Repository string
	Items      []UpdateServiceItem
	Updated    time.Time

	storageURI string
	kmURI      string
	kmMode     string
}

// DefaultUpdateService creates/loads a UpdateService from setting
func DefaultUpdateService(p, v, n, r string) (UpdateService, error) {
	storageURI, err := utils.GetSetting("storage-uri")
	if err != nil {
		return UpdateService{}, err
	}

	kmURI, _ := utils.GetSetting("keymanager-uri")
	kmMode, _ := utils.GetSetting("keymanager-mode")
	return NewUpdateService(storageURI, kmURI, kmMode, p, v, n, r)
}

// NewUpdateService creates/loads a UpdateService by a storage service, a key manager servic and 'proto', 'namespace' and 'repository'.
// key manager could be nil.
func NewUpdateService(storageURI, kmURI, kmMode, p, v, n, r string) (us UpdateService, err error) {
	if p == "" || v == "" || n == "" || r == "" {
		return UpdateService{}, errors.New("Fail to create a update service with empty Proto/Version/Namespace/Repository")
	}

	if storageURI == "" {
		return UpdateService{}, errors.New("Fail to create a update service with nil Storage interface")
	}

	key := fmt.Sprintf("%s/%s/%s/%s/%s", p, v, n, r, defaultMetaFileName)
	store, err := storage.NewUpdateServiceStorage(storageURI)
	if err != nil {
		return UpdateService{}, err
	}

	data, err := store.Get(key)
	if err == nil {
		err := json.Unmarshal(data, &us)
		if err != nil {
			return UpdateService{}, err
		}
		us.storageURI = storageURI
		us.kmURI = kmURI
		us.kmMode = kmMode
	} else if err == storage.ErrorsNotFound {
		us.storageURI = storageURI
		us.kmURI = kmURI
		us.kmMode = kmMode
		us.Proto = p
		us.Version = v
		us.Namespace = n
		us.Repository = r
		us.save()
	} else {
		return UpdateService{}, err
	}

	return us, nil
}

func (us *UpdateService) GetKM() keymanager.KeyManager {
	km, _ := keymanager.NewKeyManager(us.kmMode, us.kmURI)
	return km
}

func (us *UpdateService) GetStorage() storage.UpdateServiceStorage {
	store, _ := storage.NewUpdateServiceStorage(us.storageURI)
	return store
}

// GetMeta provides meta bytes
func (us *UpdateService) GetMeta() ([]byte, error) {
	key := fmt.Sprintf("%s/%s/%s/%s/%s", us.Proto, us.Version, us.Namespace, us.Repository, defaultMetaFileName)
	return us.GetStorage().Get(key)
}

// GetMetaSign provides meta sign bytes
func (us *UpdateService) GetMetaSign() ([]byte, error) {
	key := fmt.Sprintf("%s/%s/%s/%s/%s", us.Proto, us.Version, us.Namespace, us.Repository, defaultMetaSignFileName)
	return us.GetStorage().Get(key)
}

// TODO: this should not be in the update service, update service now is just handling meta/sign issues
// Get provides appliance data bytes
func (us *UpdateService) Get(fullname string) ([]byte, error) {
	key := fmt.Sprintf("%s/%s/%s/%s/%s", us.Proto, us.Version, us.Namespace, us.Repository, fullname)
	return us.GetStorage().Get(key)
}

// GetItem gets an UpdateServiceItem by 'fullname'
func (us *UpdateService) GetItem(fullname string) (UpdateServiceItem, error) {
	if us.Proto == "" || us.Namespace == "" || us.Repository == "" {
		return UpdateServiceItem{}, errors.New("Fail to get a meta with empty Proto/Namespace/Repository")
	}

	if fullname == "" {
		return UpdateServiceItem{}, errors.New("'FullName' should not be empty")
	}

	for _, item := range us.Items {
		if item.FullName == fullname {
			return item, nil
		}
	}

	return UpdateServiceItem{}, fmt.Errorf("Cannot find the meta item: %s", fullname)
}

// List gets files under a repo
func (us *UpdateService) List() ([]string, error) {
	//TODO

	return nil, nil
}

// Put adds an UpdateServiceItem to meta data, save both meta file and sign file
func (us *UpdateService) Put(usi UpdateServiceItem) error {
	exist := false
	for i := range us.Items {
		if us.Items[i].Equal(usi) {
			us.Items[i] = usi
			exist = true
		}
	}

	if !exist {
		us.Items = append(us.Items, usi)
	}

	if err := us.save(); err != nil {
		return err
	}

	return nil
}

// Delete removes an UpdateServiceItem from meta data, save both meta file and sign file after that
func (us *UpdateService) Delete(fullname string) error {
	exist := false
	for i := range us.Items {
		if us.Items[i].FullName == fullname {
			us.Items = append(us.Items[:i], us.Items[i+1:]...)
			exist = true
			break
		}
	}

	if !exist {
		return errors.New("Cannot find the meta item")
	}

	if err := us.save(); err != nil {
		return err
	}

	return nil
}

// save saves meta data to local file
func (us *UpdateService) save() error {
	us.Updated = time.Now()
	content, _ := json.Marshal(us)
	key := fmt.Sprintf("%s/%s/%s/%s/%s", us.Proto, us.Version, us.Namespace, us.Repository, defaultMetaFileName)
	_, err := us.GetStorage().Put(key, content)
	if err != nil {
		return err
	}

	if us.kmURI != "" {
		// write sign file, don't popup error even fail to saveSign
		us.saveSign(content)
	}

	return nil
}

// saveSign signs the meta data and save the signed data to local file
func (us *UpdateService) saveSign(content []byte) error {
	a := utils.Appliance{Proto: us.Proto, Version: us.Version, Namespace: us.Namespace, Repository: us.Repository}
	content, err := us.GetKM().Sign(a, content)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s/%s/%s/%s/%s", us.Proto, us.Version, us.Namespace, us.Repository, defaultMetaSignFileName)
	_, err = us.GetStorage().Put(key, content)
	return err
}

func (us *UpdateService) Debug() {
}
