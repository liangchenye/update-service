package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/liangchenye/update-service/utils"
)

// expunge all the registed implementaions
func preTest() {
	for n, _ := range usStorages {
		delete(usStorages, n)
	}
}

func TestRegisterUpdateServiceStorage(t *testing.T) {
	preTest()

	cases := []struct {
		name     string
		f        UpdateServiceStorage
		expected bool
	}{
		{"", &UpdateServiceStorageLocal{}, false},
		{"testsname", nil, false},
		{"testsname", &UpdateServiceStorageLocal{}, true},
		{"testsname", &UpdateServiceStorageLocal{}, false},
	}

	for _, c := range cases {
		err := RegisterStorage(c.name, c.f)
		assert.Equal(t, c.expected, err == nil, "Fail to register key manager")
	}
}

func TestNewUpdateServiceStorage(t *testing.T) {
	preTest()

	RegisterStorage("local", &UpdateServiceStorageLocal{})
	_, err := NewUpdateServiceStorage("unknown://")
	assert.Equal(t, err, ErrorsNotSupported)
}

func TestDefaultUpdateServiceStorage(t *testing.T) {

	cases := []struct {
		uri      string
		expected bool
	}{
		{"/tmp", true},
		{"", false},
		{"unknown://", false},
	}

	for _, c := range cases {
		utils.SetSetting("storage-uri", c.uri)
		_, err := DefaultUpdateServiceStorage()
		assert.Equal(t, c.expected, err == nil, "Error in creating default update service")
	}
}
