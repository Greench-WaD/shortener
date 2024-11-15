package create_json

import (
	"encoding/json"
	"github.com/Igorezka/shortener/internal/app/config"
	"github.com/Igorezka/shortener/internal/app/storage"
	"net/http"
	"net/url"
)

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}

func New(cfg *config.Config, store *storage.Store) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var r Request
		dec := json.NewDecoder(req.Body)
		if err := dec.Decode(&r); err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(r.URL) <= 0 {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		if _, err := url.ParseRequestURI(r.URL); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		id := store.DB.CreateURI(r.URL)

		resp := Response{
			Result: cfg.BaseURL + "/" + id,
		}

		res.Header().Set("Content-Type", "application/json; charset=utf-8")
		res.WriteHeader(http.StatusCreated)

		enc := json.NewEncoder(res)
		if err := enc.Encode(resp); err != nil {
			return
		}
	}
}
