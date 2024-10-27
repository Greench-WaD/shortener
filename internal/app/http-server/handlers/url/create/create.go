package create

import (
	"fmt"
	"github.com/Igorezka/shortener/internal/app/storage"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func New(store *storage.Store) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(res, "Only POST method required", http.StatusBadRequest)
			return
		}
		if !strings.Contains(req.Header.Get("Content-type"), "text/plain") {
			http.Error(res, "Only text/plain body required", http.StatusBadRequest)
			return
		}
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
		id := store.CreateURI(string(body))

		res.Header().Set("content-type", "text/plain")
		res.WriteHeader(http.StatusCreated)

		_, err = res.Write([]byte("http://localhost:8080/" + id))
		if err != nil {
			fmt.Println(err)
		}
	}
}
