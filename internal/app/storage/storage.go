package storage

import "errors"

var (
	ErrNotFound   = errors.New("not found")
	ErrCreateLink = errors.New("failed to create link")
)

type Storage interface {
	CreateURI(link string) (string, error)
	GetLink(id string) (string, error)
}

type Store struct {
	DB Storage
}

func New(storage Storage) *Store {
	return &Store{DB: storage}
}
