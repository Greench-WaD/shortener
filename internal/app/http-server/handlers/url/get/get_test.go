package get

import (
	"github.com/Igorezka/shortener/internal/app/storage"
	"github.com/Igorezka/shortener/internal/app/storage/memory"
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

	store := storage.New(memory.New())
	mux := http.NewServeMux()
	mux.HandleFunc(`/{id}`, New(store))
	srv := httptest.NewServer(mux)
	defer srv.Close()

	tests := []struct {
		name       string
		want       want
		method     string
		request    string
		wantCreate bool
	}{
		{
			name: "Positive",
			want: want{
				code:        http.StatusTemporaryRedirect,
				contentType: "text/html; charset=utf-8",
				location:    "https://ya.ru",
			},
			method:     http.MethodGet,
			request:    "/",
			wantCreate: true,
		},
		{
			name: "Negative Method not allowed",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/html; charset=utf-8",
				location:    "https://tarkov.help",
			},
			method:     http.MethodPost,
			request:    "/",
			wantCreate: false,
		},
		{
			name: "Negative Link not found",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/html; charset=utf-8",
				location:    "https://practicum.yandex.ru",
			},
			method:     http.MethodGet,
			request:    "/",
			wantCreate: false,
		},
	}
	for _, tt := range tests {
		id := shortuuid.New()
		if tt.wantCreate {
			id = store.DB.CreateURI(tt.want.location)
		}
		t.Run(tt.name, func(t *testing.T) {
			req := resty.New().SetRedirectPolicy(resty.NoRedirectPolicy()).R()
			req.Method = tt.method
			req.URL = srv.URL + "/" + id

			result, _ := req.Send()
			assert.Equal(t, tt.want.code, result.StatusCode())

			if result.StatusCode() == http.StatusTemporaryRedirect {
				assert.Equal(t, tt.want.location, result.Header().Get("Location"))
			}
		})
	}
}
