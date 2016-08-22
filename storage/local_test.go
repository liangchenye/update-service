package storage

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalSupportAndNew(t *testing.T) {
	cases := []struct {
		url      string
		expected bool
	}{
		{"", false},
		{"/tmp", true},
		{"invalid://tmp", false},
	}

	var local UpdateServiceStorageLocal
	for _, c := range cases {
		assert.Equal(t, c.expected, local.Supported(c.url), "Fail to get support status")

		_, err := local.New(c.url)
		assert.Equal(t, c.expected, err == nil, "Fail to create a new local storage interface")
	}
}

func TestLocalGet(t *testing.T) {
	var local UpdateServiceStorageLocal
	_, path, _, _ := runtime.Caller(0)
	realPath := filepath.Join(filepath.Dir(path), "testdata")
	l, _ := local.New(realPath)

	expected := "This is the content of appA.\n"
	key := "containerops/official/appA"
	content, err := l.Get(key)
	assert.Nil(t, err, "Fail to call `get`")
	assert.Equal(t, []byte(expected), content, "Fail to get the correct data")
}

func TestLocalPut(t *testing.T) {
	tmpPath, err := ioutil.TempDir("", "dus-test-")
	defer os.RemoveAll(tmpPath)
	assert.Nil(t, err, "Fail to create temp dir")

	var local UpdateServiceStorageLocal
	testData := "this is test DATA, you can put in anything here"
	key := "containerops/official/appA"
	l, _ := local.New(tmpPath)

	err = l.Put(key, []byte(testData))
	assert.Nil(t, err, "Fail to put key")

	content, _ := l.Get(key)
	assert.Equal(t, []byte(testData), content, "Fail to put correct file content")
}

func TestLocalDelete(t *testing.T) {
	tmpPath, err := ioutil.TempDir("", "dus-test-")
	defer os.RemoveAll(tmpPath)
	assert.Nil(t, err, "Fail to create temp dir")

	var local UpdateServiceStorageLocal
	testData := "this is test DATA, you can put in anything here"
	key := "containerops/official/appA"
	l, _ := local.New(tmpPath)
	l.Put(key, []byte(testData))

	err = l.Delete(key)
	assert.Nil(t, err, "Fail to delete")

	err = l.Delete(key)
	assert.NotNil(t, err, "Should not be able to delete")
}
