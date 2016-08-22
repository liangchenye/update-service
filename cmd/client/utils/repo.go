package utils

import (
	"errors"
	"fmt"
	"sync"
)

// UpdateClientRepo reprensents the local repo interface
type UpdateClientRepo interface {
	Supported(url string) bool
	New(url string) (UpdateClientRepo, error)
	List() ([]string, error)
	GetFile(name string) ([]byte, error)
	GetPublicKey() ([]byte, error)
	GetMeta() ([]byte, error)
	GetMetaSign() ([]byte, error)
	Put(name string, content []byte) error
	NRString() string
	String() string
}

var (
	ucReposLock sync.Mutex
	ucRepos     = make(map[string]UpdateClientRepo)

	// ErrorsUCRepoInvalid occurs when a repository is invalid
	ErrorsUCRepoInvalid = errors.New("repository is invalid")
	// ErrorsUCRepoNotSupported occurs when a url is not supported by existed implementations
	ErrorsUCRepoNotSupported = errors.New("repository protocal is not supported")
)

// RegisterRepo provides a way to dynamically register an implementation of a
// Repo.
//
// If RegisterRepo is called twice with the same name if Repo is nil,
// or if the name is blank, it panics.
func RegisterRepo(name string, f UpdateClientRepo) {
	if name == "" {
		panic("Could not register a Repo with an empty name")
	}
	if f == nil {
		panic("Could not register a nil Repo")
	}

	ucReposLock.Lock()
	defer ucReposLock.Unlock()

	if _, alreadyExists := ucRepos[name]; alreadyExists {
		panic(fmt.Sprintf("Repo type '%s' is already registered", name))
	}
	ucRepos[name] = f
}

// NewUCRepo creates a update client repo by a url
func NewUCRepo(url string) (UpdateClientRepo, error) {
	for _, f := range ucRepos {
		if f.Supported(url) {
			return f.New(url)
		}
	}

	return nil, ErrorsUCRepoNotSupported
}
