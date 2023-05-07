package myrouter

import (
	"github.com/EgorKo25/GophKeeper/internal/server/handlers"
	"github.com/EgorKo25/GophKeeper/internal/server/mymiddleware"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// NewRouter router for a server
func NewRouter(handler *handlers.Handler, middle *mymiddleware.MyMiddleware) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Group(func(r chi.Router) {
		r.Post("/api/user/register", handler.Register)
		r.Post("/api/user/login", handler.Login)

	})
	r.Group(func(r chi.Router) {
		r.Use(middle.CheckCookie)
	})

	return r
}
