package keymanager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// expunge all the registed implementaions
func preTest() {
	for _, f := range kms {
		delete(kms, f.Name())
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
