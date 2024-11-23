package postgres

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lithammer/shortuuid"
	"time"
)

type Storage struct {
	db *sql.DB
}

func New(dsn string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS links(
			uuid SERIAL PRIMARY KEY,
			short_url TEXT NOT NULL UNIQUE,
			original_url TEXT NOT NULL);
		CREATE INDEX IF NOT EXISTS idx_short_url ON links(short_url);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	err := s.db.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) CheckConnect(ctx context.Context) error {
	const op = "storage.postgres.CheckConnect"

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) SaveURL(ctx context.Context, link string) (string, error) {
	const op = "storage.postgres.SaveURL"
	id := shortuuid.New()
	_, err := s.db.ExecContext(ctx, "INSERT INTO links (short_url, original_url) VALUES ($1,$2)", id, link)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) GetURL(ctx context.Context, id string) (string, error) {
	const op = "storage.postgres.GetURL"
	row := s.db.QueryRowContext(ctx, "SELECT original_url FROM links WHERE short_url = $1", id)
	var url string
	err := row.Scan(&url)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return url, nil
}
