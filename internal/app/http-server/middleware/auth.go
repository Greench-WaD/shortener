package middleware

import (
	ci "github.com/Igorezka/shortener/internal/app/lib/cipher"
	"github.com/lithammer/shortuuid"
	"net/http"
	"time"
)

func setCookie(w http.ResponseWriter, r *http.Request, cipher *ci.Cipher) {
	c := http.Cookie{
		Name:    "token",
		Value:   cipher.Sile(shortuuid.New()),
		Expires: time.Now().Add(time.Minute * 15),
	}
	r.Header.Add("Cookie", c.String())
	http.SetCookie(w, &c)
}

func Authentication(cipher *ci.Cipher) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			token, err := r.Cookie("token")
			if err != nil {
				setCookie(w, r, cipher)

			} else {
				_, err := cipher.Open(token.Value)
				if err != nil {
					setCookie(w, r, cipher)
				}
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
