package memory

import (
	"fmt"
	"github.com/Igorezka/shortener/internal/app/storage"
	"github.com/lithammer/shortuuid"
)

type Storage struct {
	links map[string]string
}

func New() *Storage {
	return &Storage{links: map[string]string{}}
}

func (s *Storage) CreateURI(link string) string {
	id := shortuuid.New()
	s.links[id] = link
	return id
}

func (s *Storage) GetLink(id string) (string, error) {
	const op = "storage.GetLink"
	link, ok := s.links[id]
	if !ok {
		return "", fmt.Errorf("%s: %w", op, storage.ErrNotFound)
	}
	return link, nil
}
