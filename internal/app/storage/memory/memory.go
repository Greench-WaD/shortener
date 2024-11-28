package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Igorezka/shortener/internal/app/storage"
	"github.com/Igorezka/shortener/internal/app/storage/models"
	"github.com/lithammer/shortuuid"
	"os"
	"path/filepath"
	"slices"
	"strconv"
)

type Link struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	UserID      string
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
	if s.file != nil {
		return s.file.Close()
	}
	return nil
}

func New(fileName string) (*Storage, error) {
	const op = "storage.memory.New"
	s := &Storage{
		file:    nil,
		encoder: nil,
		pointer: 1,
	}

	if len(fileName) > 0 {
		dir, err := filepath.Abs(filepath.Dir(fileName))
		fmt.Println(dir)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err := os.Mkdir(dir, 0666)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", op, err)
			}
		}
		file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		s.file = file
		s.CreateEncoder(file)
		decoder := json.NewDecoder(file)
		for decoder.More() {
			link := Link{}
			if err := decoder.Decode(&link); err != nil {
				return nil, fmt.Errorf("%s: %w", op, err)
			}
			s.links = append(s.links, link)
			s.pointer += 1
		}
	}

	return s, nil
}

func (s *Storage) SaveURL(ctx context.Context, link string, userID string) (string, error) {
	const op = "storage.memory.SaveURL"
	l := Link{
		UUID:        strconv.Itoa(s.pointer),
		ShortURL:    shortuuid.New(),
		OriginalURL: link,
		UserID:      userID,
	}
	s.links = append(s.links, l)
	if s.encoder != nil {
		err := s.encoder.Encode(l)
		if err != nil {
			return "", fmt.Errorf("%s: %w", op, storage.ErrCreateLink)
		}
	}
	s.pointer += 1
	return l.ShortURL, nil
}

func (s *Storage) GetURL(ctx context.Context, id string) (string, error) {
	const op = "storage.memory.GetURL"
	idx := slices.IndexFunc(s.links, func(l Link) bool { return l.ShortURL == id })
	if idx == -1 {
		return "", fmt.Errorf("%s: %w", op, storage.ErrNotFound)
	}

	return s.links[idx].OriginalURL, nil
}

func (s *Storage) GetUserURLS(ctx context.Context, baseURL string, userID string) ([]models.UserBatchLink, error) {
	return []models.UserBatchLink{}, nil
}

func (s *Storage) CheckConnect(ctx context.Context) error {
	return nil
}

func (s *Storage) SaveBatchURL(ctx context.Context, baseURL string, batch []models.BatchLinkRequest, userID string) ([]models.BatchLinkResponse, error) {
	const op = "storage.memory.SaveBatchURL"

	var links []Link
	var res []models.BatchLinkResponse
	for _, b := range batch {
		id := shortuuid.New()
		links = append(links, Link{
			UUID:        strconv.Itoa(s.pointer),
			ShortURL:    id,
			OriginalURL: b.OriginalURL,
			UserID:      userID,
		})
		res = append(res, models.BatchLinkResponse{
			CorrelationID: b.CorrelationID,
			ShortURL:      baseURL + "/" + id,
		})
	}
	s.links = append(s.links, links...)
	if s.encoder != nil {
		for _, l := range links {
			err := s.encoder.Encode(l)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", op, storage.ErrCreateLink)
			}
		}
	}
	fmt.Println(s.links)
	return res, nil
}
