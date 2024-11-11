package create

import (
	"fmt"
	"github.com/Igorezka/shortener/internal/app/storage"
	"io"
	"net/http"
	"net/url"
)

func New(store *storage.Store) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = url.ParseRequestURI(string(body))
		if err != nil {
			http.Error(res, "Only valid URI required", http.StatusBadRequest)
			return
		}
		id := store.DB.CreateURI(string(body))

		res.Header().Set("content-type", "text/plain; charset=utf-8")
		res.WriteHeader(http.StatusCreated)

		_, err = res.Write([]byte("http://localhost:8080/" + id))
		if err != nil {
			fmt.Println(err)
		}
	}
}
