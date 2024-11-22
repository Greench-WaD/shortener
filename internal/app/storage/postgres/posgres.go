package postgres

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

type Storage struct {
	db *sql.DB
}

func New(ctx context.Context, dsn string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	//stmt, err := db.Prepare(`
	//	CREATE TABLE IF NOT EXISTS shortener(
	//		uuid INTEGER PRIMARY KEY,
	//		short_url TEXT NOT NULL UNIQUE,
	//		original_url TEXT NOT NULL);
	//	CREATE INDEX IF NOT EXISTS idx_short_url ON shortener(short_url);
	//`)
	//
	//_, err = stmt.Exec()
	//if err != nil {
	//	return nil, fmt.Errorf("%s: %w", op, err)
	//}

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

func (s *Storage) SaveURL(link string) (string, error) {
	return "", nil
}

func (s *Storage) GetURL(id string) (string, error) {
	return "", nil
}
