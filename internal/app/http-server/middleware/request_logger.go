package middleware

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func RequestLogger(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			responseData := &responseData{
				status: 0,
				size:   0,
			}
			lw := loggingResponseWriter{
				ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
				responseData:   responseData,
			}
			next.ServeHTTP(&lw, r)
			logger.Info("got incoming HTTP request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("duration", time.Since(start).String()),
				zap.Int("status", responseData.status),
				zap.Int("size", responseData.size),
			)
		}
		return http.HandlerFunc(fn)
	}
}
