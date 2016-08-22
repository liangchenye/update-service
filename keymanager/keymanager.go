package keymanager

import (
	"errors"
	"fmt"
	"sync"

	"github.com/liangchenye/update-service/utils"
)

// KeyManager provides interfaces to handle key issues.
type KeyManager interface {
	Name() string
	Description() string
	New(url string) (KeyManager, error)
	GenerateKey(a utils.Appliance) error
	GetPublicKey(a utils.Appliance) ([]byte, error)
	Sign(a utils.Appliance, data []byte) ([]byte, error)
	Decrypt(a utils.Appliance, data []byte) ([]byte, error)
}

var (
	kmsLock sync.Mutex
	kms     = make(map[string]KeyManager)

	// ErrorsKMNotSupported occurs when the km type is not supported
	ErrorsKMNotSupported = errors.New("key manager type is not supported")
)

// RegisterKeyManager provides a way to dynamically register an implementation of a
// key manager.
//
// If RegisterKeyManager is called twice with the same name,
// if the 'name' is blank,
// or if the 'keymanager inteface' is nil,
// it returns an error.
func RegisterKeyManager(name string, f KeyManager) error {
	if name == "" {
		return errors.New("Could not register a KeyManager with an empty name")
	}
	if f == nil {
		return errors.New("Could not register a nil KeyManager")
	}

	kmsLock.Lock()
	defer kmsLock.Unlock()

	if _, alreadyExists := kms[name]; alreadyExists {
		return fmt.Errorf("KeyManager type '%s' is already registered", name)
	}
	kms[name] = f

	return nil
}

// NewKeyManager create a key manager by its name and a storage url
func NewKeyManager(name, url string) (KeyManager, error) {
	for _, f := range kms {
		if f.Name() == name {
			return f.New(url)
		}
	}

	return nil, ErrorsKMNotSupported
}
