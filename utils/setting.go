package utils

import (
	"errors"
	"sync"
)

var (
	dsSettingsLock sync.Mutex
	dsSettings     = make(map[string]string)
)

// SetSetting sets the value by a key
func SetSetting(key string, value string) error {
	if key == "" {
		return errors.New("setting key should not be empty")
	}

	dsSettingsLock.Lock()
	defer dsSettingsLock.Unlock()

	dsSettings[key] = value
	return nil
}

// GetSetting gets the value from a key
func GetSetting(key string) (string, error) {
	if key == "" {
		return "", errors.New("setting key should not be empty")
	}

	if v, ok := dsSettings[key]; ok {
		return v, nil
	}

	return "", errors.New("setting key is not exist")
}
