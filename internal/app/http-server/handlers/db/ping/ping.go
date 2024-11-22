package ping

import (
	"context"
	resp "github.com/Igorezka/shortener/internal/app/lib/api/response"
	"go.uber.org/zap"
	"net/http"
)

type ConnectChecker interface {
	CheckConnect(ctx context.Context) error
}

func New(ctx context.Context, log *zap.Logger, checker ConnectChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := checker.CheckConnect(ctx)
		if err != nil {
			log.Error("failed to check db connect", zap.String("error", err.Error()))
			resp.Status(r, http.StatusInternalServerError)
			resp.JSON(w, r, resp.Error("internal server error"))
			return
		}

		resp.Status(r, http.StatusOK)
		resp.JSON(w, r, resp.OK())
	}
}
