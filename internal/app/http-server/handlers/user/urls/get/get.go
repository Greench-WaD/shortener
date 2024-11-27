package get

import (
	"context"
	"errors"
	"github.com/Igorezka/shortener/internal/app/config"
	resp "github.com/Igorezka/shortener/internal/app/lib/api/response"
	ci "github.com/Igorezka/shortener/internal/app/lib/cipher"
	"github.com/Igorezka/shortener/internal/app/storage"
	"github.com/Igorezka/shortener/internal/app/storage/models"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"net/http"
)

type UrlGetter interface {
	GetUserURLS(ctx context.Context, baseURL string, userID string) ([]models.UserBatchLink, error)
}

func New(log *zap.Logger, cipher *ci.Cipher, cfg *config.Config, urlGetter UrlGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.urls.get.New"

		log := log.With(
			zap.String("op", op),
			zap.String("request_id", middleware.GetReqID(r.Context())),
		)

		token, err := r.Cookie("token")
		if err != nil {
			log.Error("token required", zap.String("error", err.Error()))
			resp.Status(r, http.StatusUnauthorized)
			resp.JSON(w, r, resp.Error("Unauthorized"))
			return
		}

		userID, err := cipher.Open(token.Value)
		if err != nil {
			log.Error("failed to decode token", zap.String("error", err.Error()))
			resp.Status(r, http.StatusInternalServerError)
			resp.JSON(w, r, resp.Error("Unauthorized"))
			return
		}

		res, err := urlGetter.GetUserURLS(r.Context(), cfg.BaseURL, userID)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				http.Error(w, "Links not found", http.StatusNoContent)
				return
			}

			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		resp.Status(r, http.StatusOK)
		resp.JSON(w, r, res)
	}
}
