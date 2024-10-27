package get

import (
	"errors"
	"github.com/Igorezka/shortener/internal/app/storage"
	"net/http"
)

func New(store *storage.Store) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			http.Error(res, "Only Get method required", http.StatusBadRequest)
			return
		}
		id := req.PathValue("id")
		link, err := store.GetLink(id)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				http.Error(res, "Link not found", http.StatusBadRequest)
				return
			}

			http.Error(res, "Server error", http.StatusInternalServerError)
			return
		}
		http.Redirect(res, req, link, http.StatusTemporaryRedirect)
	}
}
