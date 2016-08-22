package storage

import (
	"errors"
	"fmt"
	"sync"
)

// UpdateServiceStorage represents the storage interface
type UpdateServiceStorage interface {
	// url is the database address or local directory (/data)
	New(url string) (UpdateServiceStorage, error)
	Supported(url string) bool
	Get(key string) ([]byte, error)
	Put(key string, data []byte) error
	Delete(key string) error
}

var (
	usStoragesLock sync.Mutex
	usStorages     = make(map[string]UpdateServiceStorage)

	// ErrorsUSSNotSupported occurs if a type is not supported
	ErrorsUSSNotSupported = errors.New("storage type is not supported")
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

// NewUSStorage creates a storage interface by a url
func NewUSStorage(url string) (UpdateServiceStorage, error) {
	for _, f := range usStorages {
		if f.Supported(url) {
			return f.New(url)
		}
	}

	return nil, ErrorsUSSNotSupported
}
