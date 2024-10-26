package storage

import (
	"errors"
	"fmt"
	"github.com/lithammer/shortuuid"
)

var ErrNotFound = errors.New("not found")

type Store struct {
	links map[string]string
}

func New() *Store {
	return &Store{links: map[string]string{}}
}

func (s *Store) CreateUri(link string) string {
	id := shortuuid.New()
	s.links[id] = link
	return id
}

func (s *Store) GetLink(id string) (string, error) {
	const op = "storage.GetLink"
	link, ok := s.links[id]
	if !ok {
		return "", fmt.Errorf("%s: %w", op, ErrNotFound)
	}
	return link, nil
}
