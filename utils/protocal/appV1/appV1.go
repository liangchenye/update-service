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

package appV1

import (
	"github.com/liangchenye/update-service/utils"
)

const (
	appV1Prefix = "appV1"
)

// UpdateServiceProtocalAppV1 is the appV1 implementation of the update service protocal
type UpdateServiceProtocalAppV1 struct {
	storage utils.UpdateServiceStorage
}

func init() {
	utils.RegisterProtocal(appV1Prefix, &UpdateServiceProtocalAppV1{})
}

// Supported checks if a protocal is 'appV1'
func (app *UpdateServiceProtocalAppV1) Supported(protocal string) bool {
	return protocal == appV1Prefix
}

// New creates a usp interface by an appV1 protocal
func (app *UpdateServiceProtocalAppV1) New(protocal string) (utils.UpdateServiceProtocal, error) {
	if protocal != appV1Prefix {
		return nil, utils.ErrorsUSPNotSupported
	}

	//FIXME: read from config
	var err error
	app.storage, err = utils.NewUSStorage("", "")
	if err != nil {
		return nil, err
	}

	return app, nil
}

// Put adds a appV1 file to a repository
func (app *UpdateServiceProtocalAppV1) Put(nr, name string, data []byte) error {
	key := nr + "/" + name
	return app.storage.Put(key, data)
}

// Get gets the appV1 file data of a repository
func (app *UpdateServiceProtocalAppV1) Get(nr, name string) ([]byte, error) {
	key := nr + "/" + name
	return app.storage.Get(key)
}

// List lists the applications of a repository
func (app *UpdateServiceProtocalAppV1) List(nr string) ([]string, error) {
	return app.storage.List(nr)
}

// GetPublicKey returns the public key data of a repository
func (app *UpdateServiceProtocalAppV1) GetPublicKey(nr string) ([]byte, error) {
	return app.storage.GetPublicKey(nr)
}

// GetMeta returns the meta data of a repository
func (app *UpdateServiceProtocalAppV1) GetMeta(nr string) ([]byte, error) {
	return app.storage.GetMeta(nr)
}

// GetMetaSign returns the meta signature data of a repository
func (app *UpdateServiceProtocalAppV1) GetMetaSign(nr string) ([]byte, error) {
	return app.storage.GetMetaSign(nr)
}
