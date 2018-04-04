package storage

import (
	"errors"
)

type Platform struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

type Storage interface {
	AddPlatform(platform *Platform) error
	GetPlatforms() (platforms []Platform, err error)
	GetPlatform(id string) (platform *Platform, err error)
	DeletePlatform(id string) (deleted bool, err error)
	UpdatePlatform(platform *Platform) error
}

var ErrNotFound = errors.New("Not found")
