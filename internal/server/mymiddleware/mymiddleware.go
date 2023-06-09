package mymiddleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/EgorKo25/GophKeeper/internal/storage"

	"github.com/EgorKo25/GophKeeper/internal/database"

	"github.com/EgorKo25/GophKeeper/pkg/auth"
)

// MyMiddleware middleware struct
type MyMiddleware struct {
	au *auth.Auth
	db *database.ManagerDB
}

// NewMyMiddleware middleware struct constructor
func NewMyMiddleware(au *auth.Auth, db *database.ManagerDB) *MyMiddleware {
	return &MyMiddleware{
		au: au,
		db: db,
	}
}

func (m *MyMiddleware) CheckUserStatus(next http.Handler) http.Handler {

	var user storage.User
	ctx := context.Background()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		login, _ := r.Cookie("User")

		user.Login = login.Value
		user.Status = true

		_, err := m.db.Read(ctx, &user, login.Value)
		if err != nil {
			if err == database.ErrRace {
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = m.db.Update(ctx, &user, login.Value)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)

		user.Status = false

		err = m.db.Update(ctx, user, login.Value)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

// CheckCookie middleware for a check cookie
func (m *MyMiddleware) CheckCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		u, err := r.Cookie("User")
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		access, err := r.Cookie("Accesses-token")
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		refresh, err := r.Cookie("Accesses-token")
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		_, err = m.au.ParseWithClaims(refresh.Value)
		if err != nil {
			log.Printf("parse token error: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = m.au.ParseWithClaims(access.Value)
		if err != nil {
			log.Printf("parse token error: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if time.Until(access.Expires) < 5*time.Minute {

			cookies, err := m.au.RefreshTokens(access.Value, refresh.Value, u.Value)
			if err != nil {
				log.Printf("create cookie error: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			http.SetCookie(w, cookies[0])
			http.SetCookie(w, cookies[1])
			http.SetCookie(w, cookies[2])
			next.ServeHTTP(w, r)
			return
		}

	})
}
