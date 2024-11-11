package router

import (
	"github.com/Igorezka/shortener/internal/app/http-server/handlers/url/create"
	"github.com/Igorezka/shortener/internal/app/http-server/handlers/url/get"
	"github.com/Igorezka/shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New(store *storage.Store) chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(middleware.AllowContentType("text/plain"))
		r.Post("/", create.New(store))
	})

	r.Get("/{id}", get.New(store))

	return r
}
