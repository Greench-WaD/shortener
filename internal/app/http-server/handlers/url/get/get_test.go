package get

import (
	"github.com/Igorezka/shortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		store *storage.Store
	}
	type want struct {
		code        int
		contentType string
		location    string
	}
	tests := []struct {
		name       string
		args       args
		want       want
		method     string
		request    string
		url        string
		wantCreate bool
	}{
		{
			name: "Positive",
			args: args{
				store: storage.New(),
			},
			want: want{
				code:        http.StatusTemporaryRedirect,
				contentType: "text/html; charset=utf-8",
				location:    "https://ya.ru",
			},
			method:     http.MethodGet,
			request:    "/",
			url:        "https://ya.ru",
			wantCreate: true,
		},
		{
			name: "Negative Method not allowed",
			args: args{
				store: storage.New(),
			},
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/html; charset=utf-8",
				location:    "https://ya.ru",
			},
			method:     http.MethodPost,
			request:    "/",
			url:        "https://ya.ru",
			wantCreate: true,
		},
		{
			name: "Negative Link not found",
			args: args{
				store: storage.New(),
			},
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/html; charset=utf-8",
				location:    "https://ya.ru",
			},
			method:     http.MethodGet,
			request:    "/",
			url:        "https://ya.ru",
			wantCreate: false,
		},
	}
	for _, tt := range tests {
		id := ""
		if tt.wantCreate {
			id = tt.args.store.CreateURI(tt.url)
		}
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.method, tt.request+id, nil)
			r.SetPathValue("id", id)
			w := httptest.NewRecorder()
			h := New(tt.args.store)

			h(w, r)

			result := w.Result()

			assert.Equal(t, tt.want.code, result.StatusCode)
			if result.StatusCode == http.StatusTemporaryRedirect {
				assert.Equal(t, tt.want.location, result.Header.Get("Location"))
			}
		})
	}
}
