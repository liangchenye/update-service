package keymanager

import (
	"errors"
	"fmt"

	"github.com/liangchenye/update-service/storage"
	"github.com/liangchenye/update-service/utils"
)

const (
	peruserName        = "peruser"
	peruserDescription = "each user has his/her own private/public key pair"

	defaultPublicKey  = "pub_key.pem"
	defaultPrivateKey = "priv_key.pem"
	defaultBitsSize   = 2048
)

// KeyManagerPeruser is the peruser implementation of a key manager
type KeyManagerPeruser struct {
	store storage.UpdateServiceStorage
}

func init() {
	RegisterKeyManager(peruserName, &KeyManagerPeruser{})
}

func (pu *KeyManagerPeruser) ModeName() string {
	return peruserName
}

func (pu *KeyManagerPeruser) Description() string {
	return peruserDescription
}

// New returns a keymanager by a uri
func (pu *KeyManagerPeruser) New(uri string) (KeyManager, error) {
	store, err := storage.NewUpdateServiceStorage(uri)
	if err != nil {
		return nil, err
	}

	pu.store = store
	return pu, nil
}

// GetPublicKey gets the public key data of a namespace
func (pu *KeyManagerPeruser) GetPublicKey(a utils.Appliance) ([]byte, error) {
	key := fmt.Sprintf("%s/%s/%s/%s", a.Proto, a.Version, a.Namespace, defaultPublicKey)
	content, err := pu.store.Get(key)
	if err != nil && err == storage.ErrorsNotFound {
		err = pu.GenerateKey(a)
		if err == nil {
			content, err = pu.store.Get(key)
		}
	}

	return content, err
}

// GenerateKey generates private key and public key and stores them
func (pu *KeyManagerPeruser) GenerateKey(a utils.Appliance) error {
	privBytes, pubBytes, err := utils.GenerateRSAKeyPair(defaultBitsSize)
	if err != nil {
		return err
	}

	privKey := fmt.Sprintf("%s/%s/%s/%s", a.Proto, a.Version, a.Namespace, defaultPrivateKey)
	_, err = pu.store.Put(privKey, privBytes)
	if err != nil {
		return err
	}

	pubKey := fmt.Sprintf("%s/%s/%s/%s", a.Proto, a.Version, a.Namespace, defaultPublicKey)
	_, err = pu.store.Put(pubKey, pubBytes)
	if err != nil {
		// make sure priv/pub key exist in pairs.
		pu.store.Delete(privKey)
	}

	return err
}

// Sign signs the data of a namespace
func (pu *KeyManagerPeruser) Sign(a utils.Appliance, data []byte) ([]byte, error) {
	key := fmt.Sprintf("%s/%s/%s/%s", a.Proto, a.Version, a.Namespace, defaultPrivateKey)
	content, err := pu.store.Get(key)
	if err != nil && err == storage.ErrorsNotFound {
		err = pu.GenerateKey(a)
		if err == nil {
			content, err = pu.store.Get(key)
		}
	}

	if err != nil {
		return nil, err
	}

	return utils.SHA256Sign(content, data)
}

// Decrypt decrypts the data of a namespace
func (pu *KeyManagerPeruser) Decrypt(a utils.Appliance, data []byte) ([]byte, error) {
	key := fmt.Sprintf("%s/%s/%s/%s", a.Proto, a.Version, a.Namespace, defaultPrivateKey)
	content, err := pu.store.Get(key)
	if err != nil {
		return nil, errors.New("Fail to load private key, cannot decrypt")
	}

	return utils.RSADecrypt(content, data)
}

func (pu *KeyManagerPeruser) Debug() {
}
