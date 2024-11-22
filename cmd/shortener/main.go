package main

import (
	"context"
	"github.com/Igorezka/shortener/internal/app/config"
	"github.com/Igorezka/shortener/internal/app/http-server/router"
	"github.com/Igorezka/shortener/internal/app/logger"
	"github.com/Igorezka/shortener/internal/app/storage/postgres"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.New()
	log, err := logger.New(cfg.LogLevel)
	if err != nil {
		panic(err)
	}
	defer log.Sync()
	store, err := postgres.New(ctx, cfg.DatabaseDSN)
	if err != nil {
		panic(err)
	}
	defer store.Close()

	log.Info(
		"starting server",
		zap.String("Address", cfg.RunAddr),
		zap.String("Base URL", cfg.BaseURL),
	)
	err = http.ListenAndServe(cfg.RunAddr, router.New(ctx, log, cfg, store))
	if err != nil {
		panic(err)
	}
}
