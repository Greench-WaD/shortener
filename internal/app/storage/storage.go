package storage

import (
	"errors"
)

type Storage interface {
	CreateURI(link string) string
	GetLink(id string) (string, error)
}

type Store struct {
	DB Storage
}

func New(storage Storage) *Store {
	return &Store{DB: storage}
}

var ErrNotFound = errors.New("not found")
