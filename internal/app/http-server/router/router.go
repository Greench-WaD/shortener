package router

import (
	"context"
	"github.com/Igorezka/shortener/internal/app/config"
	"github.com/Igorezka/shortener/internal/app/http-server/handlers/db/ping"
	"github.com/Igorezka/shortener/internal/app/http-server/handlers/url/create"
	cb "github.com/Igorezka/shortener/internal/app/http-server/handlers/url/create_batch"
	cj "github.com/Igorezka/shortener/internal/app/http-server/handlers/url/create_json"
	"github.com/Igorezka/shortener/internal/app/http-server/handlers/url/get"
	mw "github.com/Igorezka/shortener/internal/app/http-server/middleware"
	"github.com/Igorezka/shortener/internal/app/storage/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

//go:generate go run github.com/vektra/mockery/v2@v2.49.0 --name=Storage
type Storage interface {
	SaveURL(ctx context.Context, link string) (string, error)
	GetURL(ctx context.Context, id string) (string, error)
	CheckConnect(ctx context.Context) error
	SaveBatchURL(ctx context.Context, baseURL string, batch []models.BatchLinkRequest) ([]models.BatchLinkResponse, error)
}

func New(log *zap.Logger, cfg *config.Config, store Storage) chi.Router {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequestSize(500000))
		r.Use(middleware.RequestID)
		r.Use(mw.RequestLogger(log))
		r.Get("/{id}", get.New(store))
		r.Get("/ping", ping.New(log, store))
		r.Group(func(r chi.Router) {
			r.Use(middleware.AllowContentType("text/plain", "application/x-gzip"))
			r.Use(mw.GzipMiddleware)
			r.Post("/", create.New(cfg, store))
		})
		r.Group(func(r chi.Router) {
			r.Use(middleware.AllowContentType("application/json"))
			r.Use(mw.GzipMiddleware)
			r.Post("/api/shorten", cj.New(log, cfg, store))
			r.Post("/api/shorten/batch", cb.New(log, cfg, store))
		})
	})

	return r
}
