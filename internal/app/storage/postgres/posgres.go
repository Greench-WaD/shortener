package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Igorezka/shortener/internal/app/storage"
	"github.com/Igorezka/shortener/internal/app/storage/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
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
		original_url TEXT NOT NULL UNIQUE,
	    user_id TEXT NULL);
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

func (s *Storage) SaveURL(ctx context.Context, link string, userID string) (string, error) {
	const op = "storage.postgres.SaveURL"

	id := shortuuid.New()
	_, err := s.db.ExecContext(ctx, "INSERT INTO links (short_url, original_url, user_id) VALUES ($1,$2,$3)", id, link, userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			row := s.db.QueryRowContext(ctx, "SELECT short_url FROM links WHERE original_url = $1", link)
			var i string
			err := row.Scan(&i)
			if err != nil {
				return "", fmt.Errorf("%s: %w", op, err)
			}

			return i, fmt.Errorf("%s: %w", op, storage.ErrURLExist)
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// TODO: Обработка ошибок бд
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

func (s *Storage) GetUserURLS(ctx context.Context, baseURL string, userID string) ([]models.UserBatchLink, error) {
	const op = "storage.postgres.GetUserURLS"

	rows, err := s.db.QueryContext(ctx, "SELECT short_url, original_url FROM links WHERE user_id = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	links := make([]models.UserBatchLink, 0)

	for rows.Next() {
		var link models.UserBatchLink
		if err := rows.Scan(&link.ShortURL, &link.OriginalURL); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		link.ShortURL = baseURL + "/" + link.ShortURL
		links = append(links, link)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return links, nil
}

func (s *Storage) SaveBatchURL(ctx context.Context, baseURL string, batch []models.BatchLinkRequest, userID string) ([]models.BatchLinkResponse, error) {
	const op = "storage.postgres.SaveBatchURL"

	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO links (short_url, original_url, user_id) VALUES ($1,$2,$3)")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var res []models.BatchLinkResponse
	for _, b := range batch {
		id := shortuuid.New()
		_, err := stmt.ExecContext(ctx, id, b.OriginalURL, userID)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		res = append(res, models.BatchLinkResponse{
			CorrelationID: b.CorrelationID,
			ShortURL:      baseURL + "/" + id,
		})
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return res, nil
}
