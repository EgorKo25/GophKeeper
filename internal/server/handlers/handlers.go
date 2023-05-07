package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/EgorKo25/GophKeeper/internal/database"
	"github.com/EgorKo25/GophKeeper/internal/server/auth"
	"github.com/EgorKo25/GophKeeper/internal/storage"
)

var (
	cantRead      = "can't read request body: %s"
	cantUnmarshal = "can't unmarshal json obj: %s"
)

// Handler handler struct
type Handler struct {
	db *database.ManagerDB
	au *auth.Auth
}

// NewHandler Handler constructor
func NewHandler(db *database.ManagerDB, au *auth.Auth) *Handler {
	return &Handler{
		db: db,
		au: au,
	}
}

// Register register new user
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {

	var user storage.User
	var cookies []*http.Cookie

	ctx := context.Background()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf(cantRead, err)
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf(cantUnmarshal, err)
		return
	}

	exp := time.Now()

	user.CreatedAt, user.UpdatedAt = exp, exp

	err = h.db.AddUser(ctx, &user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("add user error: %s", err)
		return
	}

	cookies, err = h.au.GenerateTokensAndCreateCookie(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("create cookie error: %s", err)
		return
	}

	http.SetCookie(w, cookies[0])
	http.SetCookie(w, cookies[1])
	http.SetCookie(w, cookies[2])

	w.Header().Set("Set-Cookie", "Cookie Set")
	w.WriteHeader(http.StatusOK)
}

// Login authorize user
func (h *Handler) Login(_ http.ResponseWriter, _ *http.Request) {

}
