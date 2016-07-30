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
	_ "github.com/liangchenye/update-service/utils/storage/local"
)

const (
	appV1Prefix = "appV1"
)

type UpdateServiceProtocalAppV1 struct {
	storage utils.UpdateServiceStorage
}

func init() {
	utils.RegisterProtocal(appV1Prefix, &UpdateServiceProtocalAppV1{})
}

func (app *UpdateServiceProtocalAppV1) Supported(protocal string) bool {
	return protocal == appV1Prefix
}

func (app *UpdateServiceProtocalAppV1) New(protocal string) (utils.UpdateServiceProtocal, error) {
	if protocal != appV1Prefix {
		return nil, utils.ErrorsUSPNotSupported
	}

	//FIXME: read from config
	var err error
	app.storage, err = utils.NewUSStorage("")
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (app *UpdateServiceProtocalAppV1) Put(nr, name string, data []byte) error {
	key := nr + "/" + name
	return app.storage.Put(key, data)
}

func (app *UpdateServiceProtocalAppV1) Get(nr, name string) ([]byte, error) {
	key := nr + "/" + name
	return app.storage.Get(key)
}

func (app *UpdateServiceProtocalAppV1) List(nr string) ([]string, error) {
	return app.storage.List(nr)
}

func (app *UpdateServiceProtocalAppV1) GetPublicKey(nr string) ([]byte, error) {
	return app.storage.GetPublicKey(nr)
}

func (app *UpdateServiceProtocalAppV1) GetMeta(nr string) ([]byte, error) {
	return app.storage.GetMeta(nr)
}

func (app *UpdateServiceProtocalAppV1) GetMetaSign(nr string) ([]byte, error) {
	return app.storage.GetMetaSign(nr)
}
