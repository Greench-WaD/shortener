package get

import (
	"errors"
	"github.com/Igorezka/shortener/internal/app/storage"
	"net/http"
)

func New(store *storage.Store) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		id := req.PathValue("id")
		if id == "" {
			http.Error(res, "Id parameter required", http.StatusBadRequest)
			return
		}
		link, err := store.DB.GetLink(id)
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
