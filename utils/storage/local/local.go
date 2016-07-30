/*
Copyright 2016 The ContainerOps Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package local

import (
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/liangchenye/update-service/utils"
)

const (
	localPrefix = "local"
)

var (
	// Parse "local://tmp/containerops" and get  "Path" : "/tmp/containerops"
	localRegexp = regexp.MustCompile(`^(.+):/(.+)$`)
)

type UpdateServiceStorageLocal struct {
	Path string

	kmURL string
}

func init() {
	utils.RegisterStorage(localPrefix, &UpdateServiceStorageLocal{})
}

func (ussl *UpdateServiceStorageLocal) Supported(url string) bool {
	return strings.HasPrefix(url, localPrefix+"://")
}

func (ussl *UpdateServiceStorageLocal) New(url string) (utils.UpdateServiceStorage, error) {
	parts := localRegexp.FindStringSubmatch(url)
	if len(parts) != 3 || parts[1] != localPrefix {
		return nil, errors.New("invalid url set in StorageLocal.New")
	}

	ussl.Path = parts[2]
	ussl.kmURL = ""

	return ussl, nil
}

func (ussl *UpdateServiceStorageLocal) String() string {
	return fmt.Sprintf("%s:/%s", localPrefix, ussl.Path)
}

func (ussl *UpdateServiceStorageLocal) SetKM(kmURL string) error {
	ussl.kmURL = kmURL

	return nil
}

// Key is "namespace/repository/appname"
func (ussl *UpdateServiceStorageLocal) Get(key string) ([]byte, error) {
	s := strings.Split(key, "/")
	if len(s) != 3 {
		return nil, errors.New("invalid key detected in StorageLocal.Get")
	}

	r, err := NewRepoWithKM(ussl.Path, strings.Join(s[:2], "/"), ussl.kmURL)
	if err != nil {
		return nil, err
	}

	return r.Get(s[2])
}

// Key is "namespace/repository"
func (ussl *UpdateServiceStorageLocal) GetMeta(key string) ([]byte, error) {
	s := strings.Split(key, "/")
	if len(s) != 2 {
		return nil, errors.New("invalid key detected in StorageLocal.GetMeta")
	}

	r, err := NewRepoWithKM(ussl.Path, key, ussl.kmURL)
	if err != nil {
		return nil, err
	}

	return r.GetMeta()
}

// Key is "namespace/repository"
func (ussl *UpdateServiceStorageLocal) GetMetaSign(key string) ([]byte, error) {
	s := strings.Split(key, "/")
	if len(s) != 2 {
		return nil, errors.New("invalid key detected in StorageLocal.GetMetaSign")
	}

	r, err := NewRepoWithKM(ussl.Path, key, ussl.kmURL)
	if err != nil {
		return nil, err
	}

	file := r.GetMetaSignFile()
	return ioutil.ReadFile(file)
}

// Key is "namespace/repository"
func (ussl *UpdateServiceStorageLocal) GetPublicKey(key string) ([]byte, error) {
	s := strings.Split(key, "/")
	if len(s) != 2 {
		return nil, errors.New("invalid key detected in StorageLocal.GetPublicKey")
	}

	r, err := NewRepoWithKM(ussl.Path, key, ussl.kmURL)
	if err != nil {
		return nil, err
	}

	file := r.GetPublicKeyFile()
	return ioutil.ReadFile(file)
}

// Key is "namespace/repository/appname"
func (ussl *UpdateServiceStorageLocal) Put(key string, content []byte) error {
	s := strings.Split(key, "/")
	if len(s) != 3 {
		return errors.New("invalid key detected in StorageLocal.Put")
	}

	r, err := NewRepoWithKM(ussl.Path, strings.Join(s[:2], "/"), ussl.kmURL)
	if err != nil {
		return err
	}
	return r.Add(s[2], content)
}

// Key is "namespace/repository"
func (ussl *UpdateServiceStorageLocal) Delete(key string) error {
	s := strings.Split(key, "/")
	if len(s) != 2 {
		return errors.New("invalid key detected in StorageLocal.Delete")
	}

	r, err := NewRepoWithKM(ussl.Path, strings.Join(s[:2], "/"), ussl.kmURL)
	if err != nil {
		return err
	}

	return r.Remove(s[2])
}

// Key is "namespace/repository"
func (ussl *UpdateServiceStorageLocal) List(key string) ([]string, error) {
	s := strings.Split(key, "/")
	if len(s) != 2 {
		return nil, errors.New("invalid key deteced in StorageLocal.List")
	}

	r, err := NewRepoWithKM(ussl.Path, key, ussl.kmURL)
	if err != nil {
		return nil, err
	}

	return r.List()
}
