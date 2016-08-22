package service

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/liangchenye/update-service/keymanager"
	"github.com/liangchenye/update-service/storage"
	"github.com/stretchr/testify/assert"
)

func TestNewUpdateService(t *testing.T) {
	tmpPath, err := ioutil.TempDir("", "us-test-")
	assert.Nil(t, err, "Fail to create a temp dir")
	defer os.RemoveAll(tmpPath)

	store, _ := storage.NewUpdateServiceStorage(tmpPath)
	km, _ := keymanager.NewKeyManager("peruser", tmpPath)

	cases := []struct {
		us       UpdateService
		store    storage.UpdateServiceStorage
		km       keymanager.KeyManager
		expected bool
	}{
		{us: UpdateService{Proto: "p", Version: "v", Namespace: "n", Repository: "r"}, store: store, km: km, expected: true},
		{us: UpdateService{Proto: "p", Version: "v", Namespace: "n", Repository: "r"}, store: store, km: nil, expected: true},
		{us: UpdateService{Proto: "p", Version: "v", Namespace: "n", Repository: "r"}, store: nil, km: km, expected: false},
		{us: UpdateService{Proto: "", Version: "v", Namespace: "n", Repository: "r"}, store: store, km: km, expected: false},
		{us: UpdateService{Proto: "p", Version: "", Namespace: "n", Repository: "r"}, store: store, km: km, expected: false},
		{us: UpdateService{Proto: "p", Version: "v", Namespace: "", Repository: "r"}, store: store, km: km, expected: false},
		{us: UpdateService{Proto: "p", Version: "v", Namespace: "n", Repository: ""}, store: store, km: km, expected: false},
	}

	for _, c := range cases {
		_, err := NewUpdateService(c.store, c.km, c.us.Proto, c.us.Version, c.us.Namespace, c.us.Repository)
		assert.Equal(t, c.expected, err == nil, "Error in creating update service")
	}
}

func TestUpdateServiceOper(t *testing.T) {
	tmpPath, err := ioutil.TempDir("", "us-test-")
	assert.Nil(t, err, "Fail to create a temp dir")
	defer os.RemoveAll(tmpPath)

	store, _ := storage.NewUpdateServiceStorage(tmpPath)
	km, _ := keymanager.NewKeyManager("peruser", tmpPath)

	// add an 'fn/sha0' item
	testService, _ := NewUpdateService(store, km, "p", "v", "n", "r")
	testItem, _ := NewUpdateServiceItem("fn", []string{"sha0"})
	err = testService.Put(testItem)
	assert.Nil(t, err, "Fail to add a test item")

	// query an 'fn' item and compare it
	newService, _ := NewUpdateService(store, km, "p", "v", "n", "r")
	_, err = newService.Get("invalidfn")
	assert.NotNil(t, err, "Should not load item with invalid fullname")
	retItem, err := newService.Get("fn")
	assert.Nil(t, err, "Fail to load exist item")
	assert.Equal(t, testItem.FullName, retItem.FullName, "Fail to load the correct fullname")
	assert.Equal(t, len(testItem.SHAS), len(retItem.SHAS), "Fail to load the correct SHAS count")
	assert.Equal(t, testItem.SHAS[0], retItem.SHAS[0], "Fail to load the correct SHAS value")

	// update 'fn' item with 'sha1'
	updatedItem, _ := NewUpdateServiceItem("fn", []string{"sha0-updated"})
	err = newService.Put(updatedItem)
	assert.Nil(t, err, "Fail to update a test item")
	newUpdatedService, _ := NewUpdateService(store, km, "p", "v", "n", "r")
	retUpdatedItem, err := newUpdatedService.Get("fn")
	assert.Nil(t, err, "Fail to load exist item")
	assert.Equal(t, len(updatedItem.SHAS), len(retUpdatedItem.SHAS), "Fail to load the correct SHAS count")
	assert.Equal(t, updatedItem.SHAS[0], retUpdatedItem.SHAS[0], "Fail to load the correct SHAS value")

	// get meta file
	_, err = newService.GetMeta()
	assert.Nil(t, err, "Fail to read meta file")

	// get meta sign file
	_, err = newService.GetMetaSign()
	assert.Nil(t, err, "Fail to read meta sign")

	// delete file
	err = newService.Delete("fn")
	assert.Nil(t, err, "Fail to delete meta item")
	err = newService.Delete("fn")
	assert.NotNil(t, err, "Should return error in deleting non exist item")
	_, err = newService.Get("fn")
	assert.NotNil(t, err, "Should return error in query deleted item")
}
