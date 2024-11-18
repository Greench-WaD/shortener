package router

import (
	"github.com/Igorezka/shortener/internal/app/config"
	"github.com/Igorezka/shortener/internal/app/http-server/handlers/url/create"
	createj "github.com/Igorezka/shortener/internal/app/http-server/handlers/url/create_json"
	"github.com/Igorezka/shortener/internal/app/http-server/handlers/url/get"
	mw "github.com/Igorezka/shortener/internal/app/http-server/middleware"
	"github.com/Igorezka/shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// New TODO: убрать жесткую привязку к chi
func New(log *zap.Logger, cfg *config.Config, store *storage.Store) chi.Router {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequestID)
		r.Use(mw.RequestLogger(log))
		r.Get("/{id}", get.New(store))
		r.Group(func(r chi.Router) {
			r.Use(middleware.AllowContentType("text/plain"))
			r.Post("/", create.New(cfg, store))
		})
		r.Group(func(r chi.Router) {
			r.Use(middleware.AllowContentType("application/json"))
			r.Post("/api/shorten", createj.New(log, cfg, store))
		})
	})

	return r
}
