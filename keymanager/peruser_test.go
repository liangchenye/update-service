package keymanager

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/liangchenye/update-service/utils"
)

func TestPeruserNew(t *testing.T) {
	var keymanager KeyManagerPeruser
	validURI := "/tmp/containerops_km_cache"
	invalidURI := "keymanagerInvalid://tmp/containerops_km_cache"

	_, err := keymanager.New(validURI)
	assert.Nil(t, err, "Fail to setup a key manager")

	_, err = keymanager.New(invalidURI)
	assert.NotNil(t, err, "Should not setup an invalid key manager")
}

func TestPeruserGetPublicKey(t *testing.T) {
	tmpPath, err := ioutil.TempDir("", "dus-test-")
	defer os.RemoveAll(tmpPath)
	assert.Nil(t, err, "Fail to create temp dir")

	l, err := NewKeyManager("peruser", tmpPath)
	assert.Nil(t, err, "Fail to setup a keymanager test key manager")

	a := utils.Appliance{
		Proto:      "app",
		Version:    "v1",
		Namespace:  "containerops",
		Repository: "official",
	}

	_, err = l.GetPublicKey(a)
	assert.Nil(t, err, "Fail to get public key")
}

func TestPeruserSign(t *testing.T) {
	_, path, _, _ := runtime.Caller(0)
	realPath := filepath.Join(filepath.Dir(path), "testdata")
	l, _ := NewKeyManager("peruser", realPath)
	a := utils.Appliance{
		Proto:      "app",
		Version:    "v1",
		Namespace:  "containerops",
		Repository: "official",
	}

	expectedFile := filepath.Join(realPath, "hello.sig")
	expectedBytes, _ := ioutil.ReadFile(expectedFile)
	testFile := filepath.Join(realPath, "hello.txt")
	testBytes, _ := ioutil.ReadFile(testFile)

	data, err := l.Sign(a, testBytes)
	assert.Nil(t, err, "Fail to sign")
	assert.Equal(t, expectedBytes, data, "Fail to sign correctly")
}

func TestPeruserDecrypt(t *testing.T) {
	_, path, _, _ := runtime.Caller(0)
	realPath := filepath.Join(filepath.Dir(path), "testdata")
	l, _ := NewKeyManager("peruser", realPath)
	a := utils.Appliance{
		Proto:      "app",
		Version:    "v1",
		Namespace:  "containerops",
		Repository: "official",
	}

	expectedFile := filepath.Join(realPath, "hello.txt")
	expectedByte, _ := ioutil.ReadFile(expectedFile)
	decryptFile := filepath.Join(realPath, "hello.encrypt")
	decryptByte, _ := ioutil.ReadFile(decryptFile)

	data, err := l.Decrypt(a, decryptByte)
	assert.Nil(t, err, "Fail to decrypt")
	assert.Equal(t, expectedByte, data, "Fail to decrypt correctly")
}
