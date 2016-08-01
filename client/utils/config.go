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
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	dyutils "github.com/liangchenye/update-service/utils"
)

var (
	ErrorsUCConfigExist  = errors.New("update-service update client configuration is already exist")
	ErrorsUCEmptyURL     = errors.New("invalid repository url")
	ErrorsUCRepoExist    = errors.New("repository is already exist")
	ErrorsUCRepoNotExist = errors.New("repository is not exist")
)

const (
	topDir     = ".update-service"
	configName = "config.json"
	cacheDir   = "cache"
)

type UpdateClientConfig struct {
	DefaultServer string
	CacheDir      string
	Repos         []string
}

func (ucc *UpdateClientConfig) exist() bool {
	configFile := filepath.Join(os.Getenv("HOME"), topDir, configName)
	return dyutils.IsFileExist(configFile)
}

func (ucc *UpdateClientConfig) Init() error {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		return errors.New("Cannot get home directory")
	}

	topURL := filepath.Join(homeDir, topDir)
	cacheURL := filepath.Join(topURL, cacheDir)
	if !dyutils.IsDirExist(cacheURL) {
		if err := os.MkdirAll(cacheURL, os.ModePerm); err != nil {
			return err
		}
	}

	ucc.CacheDir = cacheURL

	if !ucc.exist() {
		return ucc.save()
	}
	return nil
}

func (ucc *UpdateClientConfig) save() error {
	data, err := json.MarshalIndent(ucc, "", "\t")
	if err != nil {
		return err
	}

	configFile := filepath.Join(os.Getenv("HOME"), topDir, configName)
	if err := ioutil.WriteFile(configFile, data, 0666); err != nil {
		return err
	}

	return nil
}

func (ucc *UpdateClientConfig) Load() error {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		return errors.New("Cannot get home directory")
	}

	content, err := ioutil.ReadFile(filepath.Join(homeDir, topDir, configName))
	if err != nil {
		return err
	}

	if err := json.Unmarshal(content, &ucc); err != nil {
		return err
	}

	if ucc.CacheDir == "" {
		ucc.CacheDir = filepath.Join(homeDir, topDir, cacheDir)
	}

	return nil
}

func (ucc *UpdateClientConfig) Add(url string) error {
	if url == "" {
		return ErrorsUCEmptyURL
	}

	var err error
	if !ucc.exist() {
		err = ucc.Init()
	} else {
		err = ucc.Load()
	}
	if err != nil {
		return err
	}

	for _, repo := range ucc.Repos {
		if repo == url {
			return ErrorsUCRepoExist
		}
	}
	ucc.Repos = append(ucc.Repos, url)

	return ucc.save()
}

func (ucc *UpdateClientConfig) Remove(url string) error {
	if url == "" {
		return ErrorsUCEmptyURL
	}

	if !ucc.exist() {
		return ErrorsUCRepoNotExist
	}

	if err := ucc.Load(); err != nil {
		return err
	}
	found := false
	for i := range ucc.Repos {
		if ucc.Repos[i] == url {
			found = true
			ucc.Repos = append(ucc.Repos[:i], ucc.Repos[i+1:]...)
			break
		}
	}
	if !found {
		return ErrorsUCRepoNotExist
	}

	return ucc.save()
}
