package createj

import (
	"encoding/json"
	"github.com/Igorezka/shortener/internal/app/config"
	"github.com/Igorezka/shortener/internal/app/http-server/handlers/url/create_json/mocks"
	"github.com/go-resty/resty/v2"
	"github.com/lithammer/shortuuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"net/url"
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
	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
	}
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
				url:         "/api/shorten",
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
				url:         "/api/shorten",
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
				url:         "/api/shorten",
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
				url:         "/api/shorten",
				body:        `{"url":"httpsyaru"}`,
				contentType: "application/json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := mocks.NewURLSaver(t)
			if tt.want.code != http.StatusBadRequest {
				store.On("SaveURL", mock.Anything, tt.want.link).Return(shortuuid.New(), nil)
			}
			log := zap.NewNop()
			handler := New(log, cfg, store)
			srv := httptest.NewServer(handler)
			defer srv.Close()

			req := resty.New().R().SetHeader("Content-Type", tt.request.contentType).SetBody(tt.request.body)
			req.Method = tt.request.method
			req.URL = srv.URL + tt.request.url

			result, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, tt.want.code, result.StatusCode())
			assert.Equal(t, tt.want.contentType, result.Header().Get("Content-Type"))

			if result.StatusCode() == http.StatusCreated {
				var r Response
				err = json.Unmarshal(result.Body(), &r)
				assert.NoError(t, err, "error unmarshal json")

				_, err := url.Parse(r.Result)
				require.NoError(t, err)
			}
		})
	}
}
