package main

import (
	"context"
	"github.com/Igorezka/shortener/internal/app/config"
	"github.com/Igorezka/shortener/internal/app/http-server/router"
	"github.com/Igorezka/shortener/internal/app/logger"
	"github.com/Igorezka/shortener/internal/app/storage/memory"
	"github.com/Igorezka/shortener/internal/app/storage/postgres"
	"go.uber.org/zap"
	"net/http"
)

type Storage interface {
	SaveURL(ctx context.Context, link string) (string, error)
	GetURL(ctx context.Context, id string) (string, error)
	CheckConnect(ctx context.Context) error
	Close() error
}

func main() {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.New()
	log, err := logger.New(cfg.LogLevel)
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	store, err := InitDB(log, cfg.DatabaseDSN, cfg.FileStoragePath)
	if err != nil {
		panic(err)
	}
	defer store.Close()

	log.Info(
		"starting server",
		zap.String("Address", cfg.RunAddr),
		zap.String("Base URL", cfg.BaseURL),
	)

	err = http.ListenAndServe(cfg.RunAddr, router.New(log, cfg, store))
	if err != nil {
		panic(err)
	}
}

func InitDB(log *zap.Logger, dsn string, storagePath string) (Storage, error) {
	switch len(dsn) > 0 {
	case true:
		store, err := postgres.New(dsn)
		if err != nil {
			return nil, err
		}
		log.Info(
			"init db",
			zap.String("type", "Postgres"),
		)
		return store, nil
	default:
		store, err := memory.New(storagePath)
		if err != nil {
			return nil, err
		}
		log.Info(
			"init db",
			zap.String("type", "memory"),
		)
		return store, nil
	}
}
