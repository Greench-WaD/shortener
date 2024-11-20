package memory

import (
	"encoding/json"
	"fmt"
	"github.com/Igorezka/shortener/internal/app/storage"
	"github.com/lithammer/shortuuid"
	"os"
	"path/filepath"
	"slices"
	"strconv"
)

type Link struct {
	Uuid        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Storage struct {
	file    *os.File
	encoder *json.Encoder
	pointer int
	links   []Link
}

func (s *Storage) CreateEncoder(file *os.File) {
	s.encoder = json.NewEncoder(file)
}

func (s *Storage) Close() error {
	return s.file.Close()
}

func New(fileName string) (*Storage, error) {
	s := &Storage{
		file:    nil,
		encoder: nil,
		pointer: 1,
	}

	if len(fileName) > 0 {
		dir, err := filepath.Abs(filepath.Dir(fileName))
		if err != nil {
			return nil, err
		}
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err := os.Mkdir(dir, 0666)
			if err != nil {
				return nil, err
			}
		}
		file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
		s.file = file
		s.CreateEncoder(file)
		decoder := json.NewDecoder(file)
		for decoder.More() {
			link := Link{}
			if err := decoder.Decode(&link); err != nil {
				return nil, err
			}
			s.links = append(s.links, link)
			s.pointer += 1
		}
	}

	return s, nil
}

func (s *Storage) CreateURI(link string) (string, error) {
	const op = "storage.memory.CreateURI"
	id := shortuuid.New()
	l := Link{
		Uuid:        strconv.Itoa(s.pointer),
		ShortURL:    id,
		OriginalURL: link,
	}
	s.links = append(s.links, l)
	if s.encoder != nil {
		err := s.encoder.Encode(l)
		if err != nil {
			return "", fmt.Errorf("%s: %w", op, storage.ErrCreateLink)
		}
	}
	s.pointer += 1
	return id, nil
}

func (s *Storage) GetLink(id string) (string, error) {
	const op = "storage.memory.GetLink"
	idx := slices.IndexFunc(s.links, func(l Link) bool { return l.Uuid == id })
	if idx == -1 {
		return "", fmt.Errorf("%s: %w", op, storage.ErrNotFound)
	}

	return s.links[idx].OriginalURL, nil
}
