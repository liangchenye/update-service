package utils

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

	store storage.UpdateServiceStorage
	km    keymanager.KeyManager
}

// NewUpdateService creates/loads a UpdateService by a storage service, a key manager servic and 'proto', 'namespace' and 'repository'.
// key manager could be nil.
func NewUpdateService(store storage.UpdateServiceStorage, km keymanager.KeyManager, kmuri, p, v, n, r string) (us UpdateService, err error) {
	if p == "" || v == "" || n == "" || r == "" {
		return UpdateService{}, errors.New("Fail to create a update service with empty Proto/Version/Namespace/Repository")
	}

	if store == nil {
		return UpdateService{}, errors.New("Fail to create a update service with nil Storage interface")
	}

	if km == nil {
		fmt.Println("KeyManager is not set or is not enabled")
	}

	key := fmt.Sprintf("%s/%s/%s/%s/%s", p, v, n, r, defaultMetaFileName)
	data, err := store.Get(key)
	if err == nil {
		if err := json.Unmarshal(data, &us); err != nil {
			return UpdateService{}, err
		}
		us.store = store
		us.km = km
	} else if err == storage.ErrorsNotFound {
		us.Proto = p
		us.Version = v
		us.Namespace = n
		us.Repository = r
		us.store = store
		us.km = km
		us.save()
	} else {
		return UpdateService{}, err
	}

	us.save()
	return us, nil
}

// GetMeta provides meta bytes
func (us *UpdateService) GetMeta() ([]byte, error) {
	key := fmt.Sprintf("%s/%s/%s/%s/%s", us.Proto, us.Version, us.Namespace, us.Repository, defaultMetaFileName)
	return us.store.Get(key)
}

// GetMetaSign provides meta sign bytes
func (us *UpdateService) GetMetaSign() ([]byte, error) {
	key := fmt.Sprintf("%s/%s/%s/%s/%s", us.Proto, us.Version, us.Namespace, us.Repository, defaultMetaSignFileName)
	return us.store.Get(key)
}

// Get gets an UpdateServiceItem by 'fullname'
func (us *UpdateService) Get(fullname string) (UpdateServiceItem, error) {
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
	err := us.store.Put(key, content)
	if err != nil {
		return err
	}

	if us.km != nil {
		// write sign file, don't popup error even fail to saveSign
		us.saveSign(content)
	}

	return nil
}

// saveSign signs the meta data and save the signed data to local file
func (us *UpdateService) saveSign(content []byte) error {
	a := utils.Appliance{Proto: us.Proto, Version: us.Version, Namespace: us.Namespace, Repository: us.Repository}
	content, err := us.km.Sign(a, content)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s/%s/%s/%s/%s", us.Proto, us.Version, us.Namespace, us.Repository, defaultMetaSignFileName)
	return us.store.Put(key, content)
}
