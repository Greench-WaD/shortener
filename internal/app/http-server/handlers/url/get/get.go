package get

import (
	"errors"
	"github.com/Igorezka/shortener/internal/app/storage"
	"net/http"
)

//go:generate go run github.com/vektra/mockery/v2@v2.49.0 --name=URLGetter
type URLGetter interface {
	GetURL(id string) (string, error)
}

func New(urlGetter URLGetter) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		id := req.PathValue("id")
		if id == "" {
			http.Error(res, "Id parameter required", http.StatusBadRequest)
			return
		}
		link, err := urlGetter.GetURL(id)
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
