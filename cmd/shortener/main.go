package main

import (
	"fmt"
	"github.com/Igorezka/shortener/internal/app/config"
	"github.com/Igorezka/shortener/internal/app/http-server/router"
	"github.com/Igorezka/shortener/internal/app/storage"
	"github.com/Igorezka/shortener/internal/app/storage/memory"
	"net/http"
)

func main() {
	cfg := config.New()
	store := storage.New(memory.New())

	fmt.Println("starting server on " + cfg.RunAddr)
	err := http.ListenAndServe(cfg.RunAddr, router.New(cfg, store))
	if err != nil {
		panic(err)
	}
}
