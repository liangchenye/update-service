package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/liangchenye/update-service/cmd/server/api"
	"github.com/liangchenye/update-service/service"
	"github.com/liangchenye/update-service/storage"
	"github.com/liangchenye/update-service/utils"
)

var (
	// ErrorsUCRepoAlreadyExist occurs when a repository is exist
	ErrorsUCRepoAlreadyExist = errors.New("repository is already exist")
	// ErrorsUCRepoNotExist occurs when a repository is not exist
	ErrorsUCRepoNotExist = errors.New("repository is not exist")
)

const (
	dirName    = ".update-service"
	configName = "config.json"
	cacheDir   = "cache"
)

// UpdateClientRepo is the saved repo
type UpdateClientRepo struct {
	Proto      string
	URL        string
	uri        string
	host       string
	namespace  string
	repository string

	//TODO: should use interface here
	protoRepo api.AppV1Repo

	store    storage.UpdateServiceStorage
	cacheDir string
}

func NewUpdateClientRepo(proto, uri string) (ucr UpdateClientRepo, err error) {
	if proto != "appv1" {
		return UpdateClientRepo{}, errors.New("Only appv1 supported now")
	}

	ucr.Proto = proto
	ucr.URL = uri
	u, err := url.Parse(uri)
	if err != nil {
		return UpdateClientRepo{}, err
	}

	if u.Scheme != "" {
		ucr.uri = fmt.Sprintf("%s://%s", u.Scheme, u.Host)
	} else {
		ucr.uri = u.Host
	}

	// u.Path will be '/namespace/repository'
	strs := strings.Split(u.Path, "/")
	if len(strs) != 3 {
		return UpdateClientRepo{}, errors.New("Invalid url type, should be 'https://host/namespace/repository'")

	}
	ucr.namespace = strs[1]
	ucr.repository = strs[2]
	ucr.protoRepo, err = api.NewAppV1Repo(ucr.uri, ucr.namespace, ucr.repository)
	if err != nil {
		return UpdateClientRepo{}, err
	}

	return
}

func (ucr *UpdateClientRepo) SetCacheDir(dir string) {
	// We can use this interface since we only have one instance
	ucr.store, _ = storage.NewUpdateServiceStorage(dir)
}

func (ucr *UpdateClientRepo) Put(name string, content []byte) error {
	_, err := ucr.protoRepo.PutFile(name, "", "", content)
	return err
}

func (ucr *UpdateClientRepo) List() ([]string, error) {
	key := fmt.Sprintf("%s/%s/%s/%s/%s", ucr.host, "app/v1", ucr.namespace, ucr.repository, "meta.json")
	metaBytes, err := ucr.store.Get(key)
	if err != nil {
		metaBytes, _, err = ucr.protoRepo.GetMeta("")
		if err != nil {
			return nil, err
		}
	}

	var meta service.UpdateService
	err = json.Unmarshal(metaBytes, &meta)
	if err != nil {
		return nil, err
	}

	var ret []string
	for _, item := range meta.Items {
		ret = append(ret, item.FullName)
	}

	return ret, nil
}

func (ucr *UpdateClientRepo) Sync() error {
	metaBytes, _, err := ucr.protoRepo.GetMeta("")
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%s/%s/%s/%s/%s", ucr.host, "app/v1", ucr.namespace, ucr.repository, "meta.json")
	_, err = ucr.store.Put(key, metaBytes)
	if err != nil {
		return err
	}

	metaSignBytes, _, err := ucr.protoRepo.GetMetaSign("")
	if err != nil {
		return err
	}
	key = fmt.Sprintf("%s/%s/%s/%s/%s", ucr.host, "app/v1", ucr.namespace, ucr.repository, "metasign")
	_, err = ucr.store.Put(key, metaSignBytes)
	if err != nil {
		//TODO: Need to rollback
		return err
	}

	pubBytes, _, err := ucr.protoRepo.GetPublicKey("")
	if err != nil {
		return err
	}
	key = fmt.Sprintf("%s/%s/%s/%s", ucr.host, "app/v1", ucr.namespace, "pubkey")
	_, err = ucr.store.Put(key, pubBytes)
	if err != nil {
		//TODO: Need to rollback
		return err
	}

	return utils.SHA256Verify(pubBytes, metaBytes, metaSignBytes)
}

func (ucr *UpdateClientRepo) GetSHAS(name string) (string, error) {
	key := fmt.Sprintf("%s/%s/%s/%s/%s", ucr.host, "app/v1", ucr.namespace, ucr.repository, "meta.json")
	metaBytes, err := ucr.store.Get(key)
	if err != nil {
		metaBytes, _, err = ucr.protoRepo.GetMeta("")
		if err != nil {
			return "", err
		}
	}

	var meta service.UpdateService
	err = json.Unmarshal(metaBytes, &meta)
	if err != nil {
		return "", err
	}
	for _, item := range meta.Items {
		if item.FullName == name {
			return item.SHAS[0], nil
		}
	}

	return "", errors.New("Cannot find the appliance")
}

func (ucr *UpdateClientRepo) Get(name string) (string, error) {
	content, _, err := ucr.protoRepo.Pull(name, "")
	if err != nil {
		return "", err
	}

	key := fmt.Sprintf("%s/%s/%s/%s/blob/%s", ucr.host, "app/v1", ucr.namespace, ucr.repository, name)
	return ucr.store.Put(key, content)
}

// UpdateClientConfig is the local configuation of a update client
type UpdateClientConfig struct {
	CacheDir    string
	DefaultRepo UpdateClientRepo
	Repos       []UpdateClientRepo

	topDir string
}

func NewUpdateClientConfig(dir string) (ucc UpdateClientConfig, err error) {
	configFile := filepath.Join(dir, configName)
	if utils.IsFileExist(configFile) {
		content, err := ioutil.ReadFile(configFile)
		if err != nil {
			return ucc, err
		}

		if err := json.Unmarshal(content, &ucc); err != nil {
			return ucc, err
		}

	}

	ucc.topDir = dir

	return ucc, err
}

func DefaultUpdateClientConfig() (UpdateClientConfig, error) {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		return UpdateClientConfig{}, errors.New("Cannot get home directory")
	}

	dir := filepath.Join(homeDir, dirName)
	return NewUpdateClientConfig(dir)
}

func (ucc *UpdateClientConfig) GetCacheDir() string {
	if ucc.CacheDir == "" {
		ucc.CacheDir = filepath.Join(ucc.topDir, cacheDir)
	}

	if !utils.IsDirExist(ucc.CacheDir) {
		os.MkdirAll(ucc.CacheDir, 0755)
	}

	return ucc.CacheDir
}

func (ucc *UpdateClientConfig) save() error {
	data, err := json.MarshalIndent(ucc, "", "\t")
	if err != nil {
		return err
	}

	if !utils.IsDirExist(ucc.topDir) {
		os.MkdirAll(ucc.topDir, 0755)
	}

	configFile := filepath.Join(ucc.topDir, configName)
	if err := ioutil.WriteFile(configFile, data, 0644); err != nil {
		return err
	}

	return nil
}

// Add adds a repo url to the config file
func (ucc *UpdateClientConfig) Add(proto, url string) error {
	if proto == "" || url == "" {
		return errors.New("Proto and URL cannot be empty")
	}

	for _, repo := range ucc.Repos {
		if repo.Proto == proto && repo.URL == url {
			return ErrorsUCRepoAlreadyExist
		}
	}
	ucc.Repos = append(ucc.Repos, UpdateClientRepo{Proto: proto, URL: url})

	return ucc.save()
}

// Remove removes a repo url from the config file
func (ucc *UpdateClientConfig) Remove(proto, url string) error {
	if url == "" {
		return errors.New("Proto and URL cannot be empty")
	}

	found := false
	for i := range ucc.Repos {
		if ucc.Repos[i].Proto == proto && ucc.Repos[i].URL == url {
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
