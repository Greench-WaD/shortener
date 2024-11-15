package create_json

import (
	"encoding/json"
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
				contentType: "application/json; charset=utf-8",
			},
			method:      http.MethodPost,
			request:     "/api/shorten",
			body:        `{"url":"https://ya.ru"}`,
			contentType: "application/json; charset=utf-8",
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
				var r Request
				err = json.Unmarshal(result.Body(), &r)
				assert.NoError(t, err, "error unmarshal json")

				parseURL, err := url.Parse(r.URL)
				require.NoError(t, err)

				link, err := store.DB.GetLink(strings.ReplaceAll(parseURL.Path, "/", ""))
				assert.NoError(t, err)
				assert.Equal(t, link, tt.body)
			}
		})
	}
}
