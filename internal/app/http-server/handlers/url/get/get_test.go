package get

import (
	"github.com/Igorezka/shortener/internal/app/http-server/handlers/url/get/mocks"
	"github.com/Igorezka/shortener/internal/app/storage"
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
		location    string
	}
	type request struct {
		id     string
		method string
	}

	tests := []struct {
		name    string
		want    want
		request request
	}{
		{
			name: "Positive",
			want: want{
				code:        http.StatusTemporaryRedirect,
				contentType: "text/html; charset=utf-8",
				location:    "https://ya.ru",
			},
			request: request{
				id:     shortuuid.New(),
				method: http.MethodGet,
			},
		},
		{
			name: "Negative Link not found",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/html; charset=utf-8",
				location:    "https://practicum.yandex.ru",
			},
			request: request{
				id:     shortuuid.New(),
				method: http.MethodGet,
			},
		},
	}
	for _, tt := range tests {
		store := mocks.NewURLGetter(t)
		if tt.want.code != http.StatusBadRequest {
			store.On("GetURL", tt.request.id).Return(tt.want.location, nil)
		} else {
			store.On("GetURL", tt.request.id).Return("", storage.ErrNotFound)
		}
		mux := http.NewServeMux()
		mux.HandleFunc(`/{id}`, New(store))
		srv := httptest.NewServer(mux)
		defer srv.Close()

		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc(`/{id}`, New(store))
			srv := httptest.NewServer(mux)
			defer srv.Close()

			req := resty.New().SetRedirectPolicy(resty.NoRedirectPolicy()).R()
			req.Method = tt.request.method
			req.URL = srv.URL + "/" + tt.request.id

			result, _ := req.Send()
			assert.Equal(t, tt.want.code, result.StatusCode())

			if result.StatusCode() == http.StatusTemporaryRedirect {
				assert.Equal(t, tt.want.location, result.Header().Get("Location"))
			}
		})
	}
}
