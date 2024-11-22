package create

import (
	"github.com/Igorezka/shortener/internal/app/config"
	"github.com/Igorezka/shortener/internal/app/http-server/handlers/url/create/mocks"
	"github.com/go-resty/resty/v2"
	"github.com/lithammer/shortuuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestNew(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	type request struct {
		method      string
		body        string
		contentType string
	}
	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
	}
	tests := []struct {
		name    string
		want    want
		request request
	}{
		{
			name: "Positive",
			want: want{
				code:        http.StatusCreated,
				contentType: "text/plain; charset=utf-8",
			},
			request: request{
				method:      http.MethodPost,
				body:        "https://ya.ru",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Negative URI required",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
			request: request{
				method:      http.MethodGet,
				body:        "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Negative Invalid url",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
			request: request{
				method:      http.MethodPost,
				body:        "httpsyaru",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := mocks.NewURLSaver(t)
			if tt.want.code != http.StatusBadRequest {
				store.On("SaveURL", tt.request.body).Return(shortuuid.New(), nil)
			}
			handler := New(cfg, store)
			srv := httptest.NewServer(handler)
			defer srv.Close()

			req := resty.New().R().SetHeader("Content-Type", tt.request.contentType).SetBody(tt.request.body)
			req.Method = tt.request.method
			req.URL = srv.URL

			result, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, tt.want.code, result.StatusCode())
			assert.Equal(t, tt.want.contentType, result.Header().Get("Content-Type"))

			if result.StatusCode() == http.StatusCreated {
				_, err := url.Parse(string(result.Body()))
				require.NoError(t, err)
			}
		})
	}
}
