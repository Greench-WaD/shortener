package get

import (
	"context"
	"errors"
	"github.com/Igorezka/shortener/internal/app/storage"
	"net/http"
)

//go:generate go run github.com/vektra/mockery/v2@v2.49.0 --name=URLGetter
type URLGetter interface {
	GetURL(ctx context.Context, id string) (string, error)
}

func New(urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "Id parameter required", http.StatusBadRequest)
			return
		}
		link, err := urlGetter.GetURL(r.Context(), id)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				http.Error(w, "Link not found", http.StatusBadRequest)
				return
			}

			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, link, http.StatusTemporaryRedirect)
	}
}
