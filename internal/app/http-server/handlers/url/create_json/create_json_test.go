package createj

import (
	"encoding/json"
	"github.com/Igorezka/shortener/internal/app/config"
	"github.com/Igorezka/shortener/internal/app/logger"
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
		link        string
		contentType string
	}
	type request struct {
		method      string
		url         string
		body        string
		contentType string
	}

	store := storage.New(memory.New())
	cfg := config.New()
	log, _ := logger.New(cfg.LogLevel)
	handler := New(log, cfg, store)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	tests := []struct {
		name        string
		want        want
		method      string
		request     request
		body        string
		contentType string
	}{
		{
			name: "Positive",
			want: want{
				code:        http.StatusCreated,
				link:        "https://ya.ru",
				contentType: "application/json",
			},
			request: request{
				method:      http.MethodPost,
				url:         srv.URL + "/api/shorten",
				body:        `{"url":"https://ya.ru"}`,
				contentType: "application/json",
			},
		},
		{
			name: "Negative failed to decode json",
			want: want{
				code:        http.StatusBadRequest,
				link:        "https://ya.ru",
				contentType: "application/json",
			},
			request: request{
				method:      http.MethodPost,
				url:         srv.URL + "/api/shorten",
				body:        `{"url":"https://ya.ru}`,
				contentType: "application/json",
			},
		},
		{
			name: "Negative url field required",
			want: want{
				code:        http.StatusBadRequest,
				link:        "https://ya.ru",
				contentType: "application/json",
			},
			request: request{
				method:      http.MethodPost,
				url:         srv.URL + "/api/shorten",
				body:        `{"link":"https://ya.ru"}`,
				contentType: "application/json",
			},
		},
		{
			name: "Negative invalid url",
			want: want{
				code:        http.StatusBadRequest,
				link:        "https://ya.ru",
				contentType: "application/json",
			},
			request: request{
				method:      http.MethodPost,
				url:         srv.URL + "/api/shorten",
				body:        `{"url":"httpsyaru"}`,
				contentType: "application/json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := resty.New().R().SetHeader("Content-Type", tt.request.contentType).SetBody(tt.request.body)
			req.Method = tt.request.method
			req.URL = tt.request.url

			result, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, tt.want.code, result.StatusCode())
			assert.Equal(t, tt.want.contentType, result.Header().Get("Content-Type"))

			if result.StatusCode() == http.StatusCreated {
				var r Response
				err = json.Unmarshal(result.Body(), &r)
				assert.NoError(t, err, "error unmarshal json")

				parseURL, err := url.Parse(r.Result)
				require.NoError(t, err)
				link, err := store.DB.GetLink(strings.ReplaceAll(parseURL.Path, "/", ""))
				assert.NoError(t, err)
				assert.Equal(t, link, tt.want.link)
			}
		})
	}
}
