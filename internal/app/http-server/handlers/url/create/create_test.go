package create

import (
	"github.com/Igorezka/shortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		store *storage.Store
	}
	type want struct {
		code        int
		contentType string
		response    string
	}
	tests := []struct {
		name        string
		args        args
		want        want
		method      string
		request     string
		body        string
		contentType string
	}{
		{
			name: "Positive",
			args: args{
				store: storage.New(),
			},
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
			name: "Negative Method not supported",
			args: args{
				store: storage.New(),
			},
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
			method:      http.MethodGet,
			request:     "/",
			body:        "https://ya.ru",
			contentType: "text/plain; charset=utf-8",
		},
		{
			name: "Negative Content type not supported",
			args: args{
				store: storage.New(),
			},
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
			method:      http.MethodPost,
			request:     "/",
			body:        "https://ya.ru",
			contentType: "application/json",
		},
		{
			name: "Negative Invalid url",
			args: args{
				store: storage.New(),
			},
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
			r := httptest.NewRequest(tt.method, tt.request, strings.NewReader(tt.body))
			r.Header.Add("Content-Type", tt.contentType)
			w := httptest.NewRecorder()
			h := New(tt.args.store)

			h(w, r)

			result := w.Result()

			assert.Equal(t, tt.want.code, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			urlResult, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			if result.StatusCode == http.StatusCreated {
				parseURL, err := url.Parse(string(urlResult))
				require.NoError(t, err)

				link, err := tt.args.store.GetLink(strings.ReplaceAll(parseUrl.Path, "/", ""))
				assert.NoError(t, err)
				assert.Equal(t, link, tt.body)
			}
		})
	}
}
