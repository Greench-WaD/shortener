package router

import (
	"context"
	"github.com/Igorezka/shortener/internal/app/config"
	"github.com/Igorezka/shortener/internal/app/http-server/handlers/db/ping"
	"github.com/Igorezka/shortener/internal/app/http-server/handlers/url/create"
	cb "github.com/Igorezka/shortener/internal/app/http-server/handlers/url/create_batch"
	cj "github.com/Igorezka/shortener/internal/app/http-server/handlers/url/create_json"
	"github.com/Igorezka/shortener/internal/app/http-server/handlers/url/get"
	getu "github.com/Igorezka/shortener/internal/app/http-server/handlers/user/urls/get"
	mw "github.com/Igorezka/shortener/internal/app/http-server/middleware"
	ci "github.com/Igorezka/shortener/internal/app/lib/cipher"
	"github.com/Igorezka/shortener/internal/app/storage/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

//go:generate go run github.com/vektra/mockery/v2@v2.49.0 --name=Storage
type Storage interface {
	SaveURL(ctx context.Context, link string, userID string) (string, error)
	GetURL(ctx context.Context, id string) (string, error)
	GetUserURLS(ctx context.Context, baseURL string, userID string) ([]models.UserBatchLink, error)
	CheckConnect(ctx context.Context) error
	SaveBatchURL(ctx context.Context, baseURL string, batch []models.BatchLinkRequest, userID string) ([]models.BatchLinkResponse, error)
}

func New(log *zap.Logger, cfg *config.Config, store Storage, cipher *ci.Cipher) chi.Router {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequestSize(500000))
		r.Use(mw.Authentication(cipher))
		r.Use(middleware.RequestID)
		r.Use(mw.RequestLogger(log))
		r.Get("/{id}", get.New(store))
		r.Get("/ping", ping.New(log, store))
		r.Get("/api/user/urls", getu.New(log, cipher, cfg, store))
		r.Group(func(r chi.Router) {
			r.Use(middleware.AllowContentType("text/plain", "application/x-gzip"))
			r.Use(mw.GzipMiddleware)
			r.Post("/", create.New(cfg, cipher, store))
		})
		r.Group(func(r chi.Router) {
			r.Use(middleware.AllowContentType("application/json"))
			r.Use(mw.GzipMiddleware)
			r.Post("/api/shorten", cj.New(log, cfg, cipher, store))
			r.Post("/api/shorten/batch", cb.New(log, cfg, cipher, store))
		})
	})

	return r
}
