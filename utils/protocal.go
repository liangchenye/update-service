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

package utils

import (
	"errors"
	"fmt"
	"sync"
)

// UpdateServiceProtocal represents the update service interface
type UpdateServiceProtocal interface {
	Supported(protocal string) bool
	New(protocal string) (UpdateServiceProtocal, error)
	List(nr string) ([]string, error)
	GetPublicKey(nr string) ([]byte, error)
	GetMeta(nr string) ([]byte, error)
	GetMetaSign(nr string) ([]byte, error)
	Get(nr string, name string) ([]byte, error)
	Put(nr string, name string, data []byte) error
}

var (
	usProtocalsLock sync.Mutex
	usProtocals     = make(map[string]UpdateServiceProtocal)

	// ErrorsUSPNotSupported occurs when a protocal is not supported
	ErrorsUSPNotSupported = errors.New("protocal is not supported")
)

// RegisterProtocal provides a way to dynamically register an implementation of a
// protocal.
//
// If RegisterProtocal is called twice with the same name if 'protocal' is nil,
// or if the name is blank, it panics.
func RegisterProtocal(name string, f UpdateServiceProtocal) {
	if name == "" {
		panic("Could not register a Protocal with an empty name")
	}
	if f == nil {
		panic("Could not register a nil Protocal")
	}

	usProtocalsLock.Lock()
	defer usProtocalsLock.Unlock()

	if _, alreadyExists := usProtocals[name]; alreadyExists {
		panic(fmt.Sprintf("Protocal type '%s' is already registered", name))
	}
	usProtocals[name] = f
}

// NewUSProtocal create a update service protocal interface by a protocal type
func NewUSProtocal(protocal string) (UpdateServiceProtocal, error) {
	for _, f := range usProtocals {
		if f.Supported(protocal) {
			return f.New(protocal)
		}
	}

	return nil, ErrorsUSPNotSupported
}
