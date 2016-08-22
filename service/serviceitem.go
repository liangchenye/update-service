package service

import (
	"errors"
	"time"
)

const (
	//The default life circle for a software is half a year
	defaultLifecircle = time.Hour * 24 * 180
)

// UpdateServiceItem keeps the meta data of a vm/app/image
type UpdateServiceItem struct {
	// Full represents a uniq name of a file within a repo, for app, fullname means os/arch/appname/tag
	FullName string
	// SHAS represents a sha list of a file.
	// If a file is composed of several layers or parts, len of SHAS will be bigger than one
	SHAS []string
	// Created is the created data of a vm/app/image
	Created time.Time
	// Updated is the latest updated data of a vm/app/image
	Updated time.Time
	// Expired is used to check if a vm/app/image need to be upgraded
	Expired time.Time
}

// NewUpdateServiceItem creates a service item by a 'FullName' and a 'SHA' list
func NewUpdateServiceItem(fn string, shas []string) (usi UpdateServiceItem, err error) {
	usi.FullName, usi.SHAS = fn, shas
	usi.Created = time.Now()
	usi.Updated = usi.Created
	usi.Expired = usi.Created.Add(defaultLifecircle)

	if ok, err := usi.isValid(); !ok {
		return usi, err
	}

	return usi, nil
}

// isValid checks the fullname and SHAs
func (usi *UpdateServiceItem) isValid() (bool, error) {
	if usi.FullName == "" || len(usi.SHAS) == 0 {
		return false, errors.New("Fullname/SHA256 fields should not be empty")
	}

	return true, nil
}

// Equal compares fullname to see if these two items are the same
func (usi *UpdateServiceItem) Equal(item UpdateServiceItem) bool {
	return usi.FullName == item.FullName
}

// GetHash get the hash strings of a file
func (usi *UpdateServiceItem) GetSHAS() []string {
	return usi.SHAS
}

// GetCreated returns the created time of an application
func (usi *UpdateServiceItem) GetCreated() time.Time {
	return usi.Created
}

// SetCreated set the created time of an application
func (usi *UpdateServiceItem) SetCreated(t time.Time) {
	usi.Created = t
}

// GetExpired get the expired time of an application
func (usi *UpdateServiceItem) GetExpired() time.Time {
	return usi.Expired
}

// SetExpired set the expired time of an application
func (usi *UpdateServiceItem) SetExpired(t time.Time) {
	usi.Expired = t
}

// IsExpired tells if an application is expired
func (usi *UpdateServiceItem) IsExpired() bool {
	return usi.Expired.Before(time.Now())
}
