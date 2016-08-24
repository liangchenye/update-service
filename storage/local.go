package storage

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"

	"github.com/liangchenye/update-service/utils"
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
	if !utils.IsFileExist(file) {
		return nil, ErrorsNotFound
	}

	return ioutil.ReadFile(file)
}

// Put adds a file with a key. Key could be "app/v1/namespace/repository/fullname"
func (ussl *UpdateServiceStorageLocal) Put(key string, content []byte) (string, error) {
	file := filepath.Join(ussl.Path, key)
	if !utils.IsDirExist(filepath.Dir(file)) {
		err := os.MkdirAll(filepath.Dir(file), 0755)
		if err != nil {
			return "", err
		}
	}

	err := ioutil.WriteFile(file, content, 0644)
	if err != nil {
		return "", err
	}
	return file, nil
}

// Delete removes a file by a key. Key could be "app/v1/namespace/repository/fullname"
func (ussl *UpdateServiceStorageLocal) Delete(key string) error {
	file := filepath.Join(ussl.Path, key)
	if !utils.IsFileExist(file) {
		return fmt.Errorf("Cannot find %s.", key)
	}

	return os.Remove(file)
}

func (ussl *UpdateServiceStorageLocal) Debug() {
}
