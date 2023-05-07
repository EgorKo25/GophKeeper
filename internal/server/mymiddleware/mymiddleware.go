package mymiddleware

import (
	"log"
	"net/http"
	"time"

	"github.com/EgorKo25/GophKeeper/internal/server/auth"
	"github.com/EgorKo25/GophKeeper/internal/storage"

	"github.com/dgrijalva/jwt-go"
)

// MyMiddleware middleware struct
type MyMiddleware struct {
	au *auth.Auth
}

// NewMyMiddleware middleware struct constructor
func NewMyMiddleware(au *auth.Auth) *MyMiddleware {
	return &MyMiddleware{au: au}
}

// CheckCookie middleware for a check cookie
func (m *MyMiddleware) CheckCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		u, err := r.Cookie("User")
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		user := &storage.User{Login: u.Value}

		access, err := r.Cookie("Accesses-token")
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		_, err = jwt.Parse(access.Value, nil)
		if err != nil {
			log.Printf("parse token error: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		refresh, err := r.Cookie("Accesses-token")
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		tokenRefresh, err := jwt.Parse(refresh.Value, nil)
		if err != nil {
			log.Printf("parse token error: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if access.Expires.Sub(time.Now()) < 5*time.Minute {
			if tokenRefresh.Valid && tokenRefresh != nil {

				cookies, err := m.au.GenerateTokensAndCreateCookie(user)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					log.Printf("create cookie error: %s", err)
					return
				}
				http.SetCookie(w, cookies[0])
				http.SetCookie(w, cookies[1])
				http.SetCookie(w, cookies[2])
				next.ServeHTTP(w, r)
			}
		}

	})
}
