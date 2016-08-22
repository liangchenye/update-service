package storage

import (
	"errors"
	"fmt"
	"sync"

	"github.com/liangchenye/update-service/utils"
)

// UpdateServiceStorage represents the storage interface
type UpdateServiceStorage interface {
	// uri is the database address or local directory (/data)
	New(uri string) (UpdateServiceStorage, error)
	Supported(uri string) bool
	Get(key string) ([]byte, error)
	Put(key string, data []byte) error
	Delete(key string) error
}

var (
	usStoragesLock sync.Mutex
	usStorages     = make(map[string]UpdateServiceStorage)

	// ErrorsNotSupported occurs if a type is not supported
	ErrorsNotSupported = errors.New("storage type is not supported")
	// ErrorsNotFound occurs if cannot find a key value
	ErrorsNotFound = errors.New("cannot find the value of the key")
)

// RegisterStorage provides a way to dynamically register an implementation of a
// storage type.
func RegisterStorage(name string, f UpdateServiceStorage) error {
	if name == "" {
		return errors.New("Could not register a Storage with an empty name")
	}
	if f == nil {
		return errors.New("Could not register a nil Storage")
	}

	usStoragesLock.Lock()
	defer usStoragesLock.Unlock()

	if _, alreadyExists := usStorages[name]; alreadyExists {
		return fmt.Errorf("Storage type '%s' is already registered", name)
	}

	usStorages[name] = f

	return nil
}

// NewUpdateServiceStorage creates a storage interface by a uri
func NewUpdateServiceStorage(uri string) (UpdateServiceStorage, error) {
	for _, f := range usStorages {
		if f.Supported(uri) {
			return f.New(uri)
		}
	}

	return nil, ErrorsNotSupported
}

// DefaultUpdateServiceStorage load default update service storage with uri from setting
func DefaultUpdateServiceStorage() (UpdateServiceStorage, error) {
	uri, err := utils.GetSetting("storage-uri")
	if err != nil {
		return nil, err
	}
	return NewUpdateServiceStorage(uri)
}
