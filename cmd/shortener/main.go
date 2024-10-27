package main

import (
	"github.com/Igorezka/shortener/internal/app/http-server/handlers/url/create"
	"github.com/Igorezka/shortener/internal/app/http-server/handlers/url/get"
	"github.com/Igorezka/shortener/internal/app/storage"
	"net/http"
)

func main() {
	store := storage.New()

	mux := http.NewServeMux()
	mux.HandleFunc(`/`, create.New(store))
	mux.HandleFunc(`/{id}`, get.New(store))

	err := http.ListenAndServe(`localhost:8080`, mux)
	if err != nil {
		panic(err)
	}
}
