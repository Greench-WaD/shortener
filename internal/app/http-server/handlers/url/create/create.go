package create

import (
	"context"
	"errors"
	"fmt"
	"github.com/Igorezka/shortener/internal/app/config"
	ci "github.com/Igorezka/shortener/internal/app/lib/cipher"
	"github.com/Igorezka/shortener/internal/app/storage"
	"io"
	"net/http"
	"net/url"
)

//go:generate go run github.com/vektra/mockery/v2@v2.49.0 --name=URLSaver
type URLSaver interface {
	SaveURL(ctx context.Context, url string, userID string) (string, error)
}

// TODO: логгирование и рефакторинг
func New(cfg *config.Config, cipher *ci.Cipher, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID, err := cipher.Open(token.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

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

		id, err := urlSaver.SaveURL(r.Context(), string(body), userID)
		if err != nil {
			if errors.Is(err, storage.ErrURLExist) {
				w.Header().Set("content-type", "text/plain")
				w.WriteHeader(http.StatusConflict)
				_, err = w.Write([]byte(cfg.BaseURL + "/" + id))
				if err != nil {
					fmt.Println(err)
				}
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusCreated)

		_, err = w.Write([]byte(cfg.BaseURL + "/" + id))
		if err != nil {
			fmt.Println(err)
		}
	}
}
