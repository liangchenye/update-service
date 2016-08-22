package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTmpHome(t *testing.T) (string, string) {
	tmpHome, err := ioutil.TempDir("", "duc-test-")
	assert.Nil(t, err, "Fail to create temp directory")

	savedHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpHome)

	return tmpHome, savedHome
}

// TestInitConfig tests the Init function
func TestInitConfig(t *testing.T) {
	tmpHome, savedHome := createTmpHome(t)
	defer os.RemoveAll(tmpHome)
	defer os.Setenv("HOME", savedHome)

	var conf UpdateClientConfig
	err := conf.Init()
	assert.Nil(t, err, "Fail to init config")
}

// TestLoadConfig tests the testdata/home/.update-service/config.json file
func TestLoadConfig(t *testing.T) {
	_, path, _, _ := runtime.Caller(0)
	testHome := "/testdata/home"
	savedHome := os.Getenv("HOME")
	defer os.Setenv("HOME", savedHome)
	os.Setenv("HOME", filepath.Join(filepath.Dir(path), testHome))

	var conf UpdateClientConfig
	err := conf.Load()
	assert.Nil(t, err, "Fail to load config")
	assert.Equal(t, conf.DefaultServer, "containerops.me", "Fail to load 'DefaultServer'")
	assert.Equal(t, conf.CacheDir, "/tmp/containeropsCache", "Fail to load 'CacheDir'")
}

func TestAddRemoveConfig(t *testing.T) {
	tmpHome, savedHome := createTmpHome(t)
	defer os.RemoveAll(tmpHome)
	defer os.Setenv("HOME", savedHome)

	var conf UpdateClientConfig
	invalidURL := ""
	validURL := "app://containerops/official/duc.rpm"

	// 'add'
	err := conf.Add(invalidURL)
	assert.Equal(t, err, ErrorsUCEmptyURL)
	err = conf.Add(validURL)
	assert.Nil(t, err, "Failed to add repository")
	err = conf.Add(validURL)
	assert.Equal(t, err, ErrorsUCRepoExist)

	// 'remove'
	err = conf.Remove(validURL)
	assert.Nil(t, err, "Failed to remove repository")
	err = conf.Remove(validURL)
	assert.Equal(t, err, ErrorsUCRepoNotExist)
}
