package createj

import (
	"github.com/Igorezka/shortener/internal/app/config"
	resp "github.com/Igorezka/shortener/internal/app/lib/api/response"
	"github.com/Igorezka/shortener/internal/app/storage"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"go.uber.org/zap"
	"net/http"
	"net/url"
)

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	resp.Response
	Result string `json:"result"`
}

func New(log *zap.Logger, cfg *config.Config, store *storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.create_json.New"
		log = log.With(
			zap.String("op", op),
			zap.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request", zap.String("error", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("internal server error"))
			return
		}

		if len(req.URL) <= 0 {
			log.Info("url field required")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("url required"))
			return
		}

		if _, err := url.ParseRequestURI(req.URL); err != nil {
			log.Info("invalid url", zap.String("url", req.URL))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("only valid url required"))
			return
		}

		id := store.DB.CreateURI(req.URL)

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, Response{
			Response: resp.OK(),
			Result:   cfg.BaseURL + "/" + id,
		})
	}
}
