package router

import (
	"errors"
	"github.com/Igorezka/shortener/internal/app/config"
	"github.com/Igorezka/shortener/internal/app/http-server/router/mocks"
	"github.com/Igorezka/shortener/internal/app/storage"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-resty/resty/v2"
	"github.com/lithammer/shortuuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {
	type want struct {
		code        int
		contentType string
		response    string
	}
	type req struct {
		srvURL      string
		method      string
		contentType string
		body        string
		id          string
	}
	cfg := &config.Config{
		RunAddr:  "localhost:8080",
		BaseURL:  "http://localhost:8080",
		LogLevel: "Info",
	}

	tests := []struct {
		name string
		want want
		req  req
	}{
		{
			name: "Positive create",
			want: want{
				code:        http.StatusCreated,
				contentType: "text/plain",
			},
			req: req{
				srvURL:      "",
				method:      http.MethodPost,
				contentType: "text/plain",
			},
		},
		{
			name: "Negative create invalid method",
			want: want{
				code:        http.StatusMethodNotAllowed,
				contentType: "",
			},
			req: req{
				srvURL:      "",
				method:      http.MethodGet,
				contentType: "text/plain",
			},
		},
		{
			name: "Negative create invalid content/type",
			want: want{
				code:        http.StatusUnsupportedMediaType,
				contentType: "",
			},
			req: req{
				srvURL:      "",
				method:      http.MethodPost,
				contentType: "application/json",
			},
		},
		{
			name: "Positive get",
			want: want{
				code:        http.StatusTemporaryRedirect,
				contentType: "text/html",
			},
			req: req{
				srvURL:      "",
				method:      http.MethodGet,
				contentType: "application/json",
				id:          shortuuid.New(),
			},
		},
		{
			name: "Negative get link not found",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain",
			},
			req: req{
				srvURL:      "",
				method:      http.MethodGet,
				contentType: "application/json",
				id:          shortuuid.New(),
			},
		},
		{
			name: "Positive create json",
			want: want{
				code:        http.StatusCreated,
				contentType: "application/json",
			},
			req: req{
				srvURL:      "/api/shorten",
				method:      http.MethodPost,
				contentType: "application/json",
			},
		},
		{
			name: "Negative create json invalid method",
			want: want{
				code:        http.StatusMethodNotAllowed,
				contentType: "",
			},
			req: req{
				srvURL:      "/api/shorten",
				method:      http.MethodGet,
				contentType: "application/json",
			},
		},
		{
			name: "Negative create json invalid content/type",
			want: want{
				code:        http.StatusUnsupportedMediaType,
				contentType: "",
			},
			req: req{
				srvURL:      "/api/shorten",
				method:      http.MethodPost,
				contentType: "text/plain",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := mocks.NewStorage(t)
			log := zap.NewNop()
			router := New(log, cfg, store)
			srv := httptest.NewServer(router)
			defer srv.Close()

			if tt.req.method == http.MethodPost && tt.want.code == http.StatusCreated {
				url := gofakeit.URL()
				if tt.req.contentType == "text/plain" {
					tt.req.body = url
				} else if tt.req.contentType == "application/json" {
					tt.req.body = `{"url":"` + url + `"}`
				}
				store.On("SaveURL", mock.Anything, url).Return(shortuuid.New(), nil)
			}

			if tt.req.method == http.MethodGet && tt.want.code == http.StatusTemporaryRedirect {
				tt.req.srvURL += "/" + tt.req.id
				store.On("GetURL", mock.Anything, tt.req.id).Return(gofakeit.URL(), nil)
			} else if tt.req.method == http.MethodGet && tt.want.code == http.StatusBadRequest {
				tt.req.srvURL += "/" + tt.req.id
				store.On("GetURL", mock.Anything, tt.req.id).Return("", storage.ErrNotFound)
			}

			req := resty.New().SetRedirectPolicy(resty.NoRedirectPolicy()).R().SetHeader("Content-Type", tt.req.contentType).SetHeader("Accept-Encoding", "gzip")
			if tt.req.method == http.MethodPost {
				req.SetBody(tt.req.body)
			}
			req.Method = tt.req.method
			req.URL = srv.URL + tt.req.srvURL

			result, err := req.Send()
			if !errors.Is(err, resty.ErrAutoRedirectDisabled) {
				assert.NoError(t, err, "error making HTTP request")
			}

			assert.Equal(t, tt.want.code, result.StatusCode())
			assert.Contains(t, result.Header().Get("Content-Type"), tt.want.contentType)
		})
	}
}
