package router

import (
	"errors"
	"github.com/Igorezka/shortener/internal/app/config"
	"github.com/Igorezka/shortener/internal/app/storage"
	"github.com/Igorezka/shortener/internal/app/storage/memory"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-resty/resty/v2"
	"github.com/lithammer/shortuuid"
	"github.com/stretchr/testify/assert"
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
		method      string
		contentType string
		url         string
		create      bool
	}

	store := storage.New(memory.New())
	cfg := config.New()
	router := New(cfg, store)
	srv := httptest.NewServer(router)
	defer srv.Close()

	tests := []struct {
		name string
		want want
		req  req
	}{
		{
			name: "Positive create",
			want: want{
				code:        http.StatusCreated,
				contentType: "text/plain; charset=utf-8",
			},
			req: req{
				method:      http.MethodPost,
				contentType: "text/plain; charset=utf-8",
				url:         gofakeit.URL(),
				create:      false,
			},
		},
		{
			name: "Negative create invalid method",
			want: want{
				code:        http.StatusMethodNotAllowed,
				contentType: "",
			},
			req: req{
				method:      http.MethodGet,
				contentType: "text/plain; charset=utf-8",
				url:         gofakeit.URL(),
				create:      false,
			},
		},
		{
			name: "Negative create invalid content/type",
			want: want{
				code:        http.StatusUnsupportedMediaType,
				contentType: "",
			},
			req: req{
				method:      http.MethodPost,
				contentType: "application/json; charset=utf-8",
				url:         "{'url':'https://ya.ru'}",
				create:      false,
			},
		},
		{
			name: "Positive get",
			want: want{
				code:        http.StatusTemporaryRedirect,
				contentType: "text/html; charset=utf-8",
			},
			req: req{
				method:      http.MethodGet,
				contentType: "application/json; charset=utf-8",
				url:         gofakeit.URL(),
				create:      true,
			},
		},
		{
			name: "Negative get link not found",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
			req: req{
				method:      http.MethodGet,
				contentType: "application/json; charset=utf-8",
				url:         gofakeit.URL(),
				create:      false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := srv.URL
			if tt.req.method == http.MethodGet && tt.want.code != http.StatusMethodNotAllowed {
				id := shortuuid.New()
				if tt.req.create {
					id = store.DB.CreateURI(tt.req.url)
				}
				url += "/" + id
			}

			req := resty.New().SetRedirectPolicy(resty.NoRedirectPolicy()).R().SetHeader("Content-Type", tt.req.contentType).SetBody(tt.req.url)
			req.Method = tt.req.method
			req.URL = url

			result, err := req.Send()
			if !errors.Is(err, resty.ErrAutoRedirectDisabled) {
				assert.NoError(t, err, "error making HTTP request")
			}

			assert.Equal(t, tt.want.code, result.StatusCode())
			assert.Equal(t, tt.want.contentType, result.Header().Get("Content-Type"))

			if result.StatusCode() == http.StatusTemporaryRedirect {
				assert.Equal(t, tt.req.url, result.Header().Get("Location"))
			}
		})
	}
}
