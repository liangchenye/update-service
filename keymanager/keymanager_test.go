package keymanager

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/liangchenye/update-service/utils"
)

// expunge all the registed implementaions
func preTest() {
	for _, f := range kms {
		delete(kms, f.ModeName())
	}
}

func TestRegisterKeyManager(t *testing.T) {
	preTest()

	cases := []struct {
		name     string
		f        KeyManager
		expected bool
	}{
		{"", &KeyManagerPeruser{}, false},
		{"testkmname", nil, false},
		{"testkmname", &KeyManagerPeruser{}, true},
		{"testkmname", &KeyManagerPeruser{}, false},
	}

	for _, c := range cases {
		err := RegisterKeyManager(c.name, c.f)
		assert.Equal(t, c.expected, err == nil, "Fail to register key manager")
	}
}

func TestNewKeyManager(t *testing.T) {
	preTest()

	_, err := NewKeyManager("unknown", "/tmp")
	assert.Equal(t, err, ErrorsKMNotSupported)
}

func TestDefaultKeyManager(t *testing.T) {
	preTest()
	RegisterKeyManager("peruser", &KeyManagerPeruser{})

	cases := []struct {
		mode     string
		uri      string
		expected bool
	}{
		{"peruser", "/tmp", true},
		{"", "", true},
		{"unknown", "/tmp", false},
		{"", "/tmp", false},
		{"peruser", "", false},
	}

	for _, c := range cases {
		utils.SetSetting("keymanager-mode", c.mode)
		utils.SetSetting("keymanager-uri", c.uri)
		_, err := DefaultKeyManager()
		t.Log(err)
		assert.Equal(t, c.expected, err == nil, "Error in creating default key manager")
	}
}
