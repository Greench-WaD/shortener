package create

import (
	"context"
	"fmt"
	"github.com/Igorezka/shortener/internal/app/config"
	"io"
	"net/http"
	"net/url"
)

//go:generate go run github.com/vektra/mockery/v2@v2.49.0 --name=URLSaver
type URLSaver interface {
	SaveURL(ctx context.Context, url string) (string, error)
}

func New(cfg *config.Config, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(string(body)) <= 0 {
			http.Error(w, "URI required", http.StatusBadRequest)
			return
		}

		if _, err = url.ParseRequestURI(string(body)); err != nil {
			http.Error(w, "Only valid URI required", http.StatusBadRequest)
			return
		}

		id, err := urlSaver.SaveURL(r.Context(), string(body))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusCreated)

		_, err = w.Write([]byte(cfg.BaseURL + "/" + id))
		if err != nil {
			fmt.Println(err)
		}
	}
}
