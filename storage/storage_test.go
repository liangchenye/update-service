package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	_, err := NewUSStorage("unknown://")
	assert.Equal(t, err, ErrorsUSSNotSupported)
}
