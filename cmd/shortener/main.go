package main

import (
	"github.com/Igorezka/shortener/internal/app/http-server/router"
	"github.com/Igorezka/shortener/internal/app/storage"
	"github.com/Igorezka/shortener/internal/app/storage/memory"
	"net/http"
)

func main() {
	store := storage.New(memory.New())

	err := http.ListenAndServe(`localhost:8080`, router.New(store))
	if err != nil {
		panic(err)
	}
}
