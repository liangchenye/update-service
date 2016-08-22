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

package storage

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"

	"github.com/containerops/dockyard/utils"
)

const (
	localName = "local"
)

// UpdateServiceStorageLocal is the local file implementation of storage service
type UpdateServiceStorageLocal struct {
	Path string
}

func init() {
	RegisterStorage(localName, &UpdateServiceStorageLocal{})
}

// Supported checks if a uri is a local path
func (ussl *UpdateServiceStorageLocal) Supported(uri string) bool {
	if uri == "" {
		return false
	}

	if u, err := url.Parse(uri); err != nil {
		return false
	} else if u.Scheme == "" {
		return true
	}

	return false
}

// New creates an UpdateServceStorage interface with a local implmentation
func (ussl *UpdateServiceStorageLocal) New(uri string) (UpdateServiceStorage, error) {
	if !ussl.Supported(uri) {
		return nil, fmt.Errorf("invalid uri set in StorageLocal.New: %s", uri)
	}

	ussl.Path = uri

	return ussl, nil
}

// Get the data of an input key. Key could be "app/v1/namespace/repository/fullname"
func (ussl *UpdateServiceStorageLocal) Get(key string) ([]byte, error) {
	file := filepath.Join(ussl.Path, key)
	return ioutil.ReadFile(file)
}

// Put adds a file with a key. Key could be "app/v1/namespace/repository/fullname"
func (ussl *UpdateServiceStorageLocal) Put(key string, content []byte) error {
	file := filepath.Join(ussl.Path, key)
	if !utils.IsDirExist(filepath.Dir(file)) {
		err := os.MkdirAll(filepath.Dir(file), 0755)
		if err != nil {
			return err
		}
	}

	return ioutil.WriteFile(file, content, 0644)
}

// Delete removes a file by a key. Key could be "app/v1/namespace/repository/fullname"
func (ussl *UpdateServiceStorageLocal) Delete(key string) error {
	file := filepath.Join(ussl.Path, key)
	if !utils.IsFileExist(file) {
		return fmt.Errorf("Cannot find %s.", key)
	}

	return os.Remove(file)
}
