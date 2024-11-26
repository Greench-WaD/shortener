package createb

import (
	"context"
	"github.com/Igorezka/shortener/internal/app/config"
	"github.com/Igorezka/shortener/internal/app/lib/api/request"
	resp "github.com/Igorezka/shortener/internal/app/lib/api/response"
	"github.com/Igorezka/shortener/internal/app/storage/models"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"net/http"
	"net/url"
)

type URLSaver interface {
	SaveBatchURL(ctx context.Context, baseURL string, batch []models.BatchLinkRequest) ([]models.BatchLinkResponse, error)
}

func New(log *zap.Logger, cfg *config.Config, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.create_batch.New"

		log = log.With(
			zap.String("op", op),
			zap.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req []models.BatchLinkRequest
		err := request.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request", zap.String("error", err.Error()))
			resp.Status(r, http.StatusBadRequest)
			resp.JSON(w, r, resp.Error("internal server error"))
			return
		}

		if len(req) <= 0 {
			log.Info("urls required")
			resp.Status(r, http.StatusBadRequest)
			resp.JSON(w, r, resp.Error("urls required"))
			return
		}

		for _, b := range req {
			if _, err := url.ParseRequestURI(b.OriginalURL); err != nil {
				log.Info("only valid urls required")
				resp.Status(r, http.StatusBadRequest)
				resp.JSON(w, r, resp.Error("only valid urls required: "+b.OriginalURL))
				return
			}
		}

		res, err := urlSaver.SaveBatchURL(r.Context(), cfg.BaseURL, req)
		if err != nil {
			log.Error("failed to store link", zap.String("error", err.Error()))
			resp.Status(r, http.StatusInternalServerError)
			resp.JSON(w, r, resp.Error("internal server error"))
			return
		}

		resp.Status(r, http.StatusCreated)
		resp.JSON(w, r, res)
	}
}
