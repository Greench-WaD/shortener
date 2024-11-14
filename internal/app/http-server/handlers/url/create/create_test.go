package create

import (
	"github.com/Igorezka/shortener/internal/app/config"
	"github.com/Igorezka/shortener/internal/app/storage"
	"github.com/Igorezka/shortener/internal/app/storage/memory"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}

	store := storage.New(memory.New())
	cfg := config.New()
	handler := New(cfg, store)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	tests := []struct {
		name        string
		want        want
		method      string
		request     string
		body        string
		contentType string
	}{
		{
			name: "Positive",
			want: want{
				code:        http.StatusCreated,
				contentType: "text/plain; charset=utf-8",
			},
			method:      http.MethodPost,
			request:     "/",
			body:        "https://ya.ru",
			contentType: "text/plain; charset=utf-8",
		},
		{
			name: "Negative URI required",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
			method:      http.MethodGet,
			request:     "/",
			body:        "",
			contentType: "text/plain; charset=utf-8",
		},
		{
			name: "Negative Invalid url",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
			method:      http.MethodPost,
			request:     "/",
			body:        "httpsyaru",
			contentType: "text/plain; charset=utf-8",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := resty.New().R().SetHeader("Content-Type", tt.contentType).SetBody(tt.body)
			req.Method = tt.method
			req.URL = srv.URL

			result, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, tt.want.code, result.StatusCode())
			assert.Equal(t, tt.want.contentType, result.Header().Get("Content-Type"))

			if result.StatusCode() == http.StatusCreated {
				parseURL, err := url.Parse(string(result.Body()))
				require.NoError(t, err)

				link, err := store.DB.GetLink(strings.ReplaceAll(parseURL.Path, "/", ""))
				assert.NoError(t, err)
				assert.Equal(t, link, tt.body)
			}
		})
	}
}
