package api

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/liangchenye/update-service/utils"
)

type AppV1Repo struct {
	URI        string
	Namespace  string
	Repository string

	host string
}

func NewAppV1Repo(uri, n, r string) (AppV1Repo, error) {
	if uri == "" {
		return AppV1Repo{}, errors.New("URI should not be empty")
	}

	u, err := url.Parse(uri)
	if err != nil {
		return AppV1Repo{}, err
	}

	var o AppV1Repo
	o.URI = uri
	o.Namespace = n
	o.Repository = r
	o.host = u.Host
	return o, nil
}

func (o *AppV1Repo) pullData(rawurl, token string) ([]byte, int, error) {
	header := map[string]string{
		"Host":          o.host,
		"Authorization": token,
	}

	resp, err := sendHttpRequest("GET", rawurl, nil, header)
	if err != nil {
		return nil, 0, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return body, resp.StatusCode, nil
}

func (o *AppV1Repo) GetMeta(token string) ([]byte, int, error) {
	rawurl := fmt.Sprintf("%s/app/v1/%s/%s/meta", o.URI, o.Namespace, o.Repository)

	return o.pullData(rawurl, token)
}

func (o *AppV1Repo) GetMetaSign(token string) ([]byte, int, error) {
	rawurl := fmt.Sprintf("%s/app/v1/%s/%s/metasign", o.URI, o.Namespace, o.Repository)

	return o.pullData(rawurl, token)
}

func (o *AppV1Repo) GetPublicKey(token string) ([]byte, int, error) {
	rawurl := fmt.Sprintf("%s/app/v1/%s/pubkey", o.URI, o.Namespace)

	return o.pullData(rawurl, token)
}

func (o *AppV1Repo) Pull(name string, token string) ([]byte, int, error) {
	rawurl := fmt.Sprintf("%s/app/v1/%s/%s/blob/%s", o.URI, o.Namespace, o.Repository, name)

	return o.pullData(rawurl, token)
}

func (o *AppV1Repo) PutFile(name string, token, uuid string, fileBytes []byte) (int, error) {
	rawurl := fmt.Sprintf("%s/app/v1/%s/%s/%s", o.URI, o.Namespace, o.Repository, name)

	sha512Sum, err := utils.SHA512(fileBytes)
	if err != nil {
		return 0, err
	}

	digest := fmt.Sprintf("%s:%s", "sha512", sha512Sum)
	header := map[string]string{
		"Host":            o.host,
		"Authorization":   token,
		"App-Upload-UUID": uuid,
		"Digest":          digest,
	}
	resp, err := sendHttpRequest("PUT", rawurl, bytes.NewReader(fileBytes), header)
	if err != nil {
		return 0, err
	}
	return resp.StatusCode, nil
}

func (o *AppV1Repo) Delete(name string, token string) (int, error) {
	rawurl := fmt.Sprintf("%s/app/v1/%s/%s/%s", o.URI, o.Namespace, o.Repository, name)
	header := map[string]string{
		"Host":          o.host,
		"Authorization": token,
	}
	resp, err := sendHttpRequest("DELETE", rawurl, nil, header)
	if err != nil {
		return 0, err
	}

	return resp.StatusCode, nil
}
