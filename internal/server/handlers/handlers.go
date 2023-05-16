package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/EgorKo25/GophKeeper/pkg/mycrypto"

	"github.com/EgorKo25/GophKeeper/pkg/auth"

	"github.com/EgorKo25/GophKeeper/internal/database"
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
	cr *mycrypto.Crypto
}

// NewHandler Handler constructor
func NewHandler(db *database.ManagerDB, au *auth.Auth, cr *mycrypto.Crypto) *Handler {
	return &Handler{
		db: db,
		au: au,
		cr: cr,
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

	if user.Login == "" || user.Email == "" || user.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = h.db.Read(ctx, &user, user.Login)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error: %s", err)
		return
	}

	exp := time.Now()

	user.CreatedAt, user.UpdatedAt = exp, exp

	err = h.db.Add(ctx, &user, user.Login)
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

	w.WriteHeader(http.StatusOK)
}

// Login authorize user
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
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

	isUserExist, err := h.db.CheckUser(ctx, &user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s", err)
		return
	}

	if !isUserExist {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	cookies, err = h.au.GenerateTokensAndCreateCookie(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("create cookie error: %s", err)
		return
	}

	log.Println(cookies)
	http.SetCookie(w, cookies[0])
	http.SetCookie(w, cookies[1])
	http.SetCookie(w, cookies[2])

	w.WriteHeader(http.StatusOK)
}

// Add user data to database
func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf(cantRead, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resType := r.Header.Get("Data-Type")
	data, err := myUnmarshal(resType, body)
	if err != nil {
		log.Printf("%s: %s", err, resType)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cook, err := r.Cookie("User")
	if err != nil {
		log.Printf("%s", err)
	}

	err = h.db.Add(ctx, data, cook.Value)
	if err != nil {
		log.Printf("%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func myUnmarshal(t string, body []byte) (any, error) {
	switch t {
	case "card":
		res := storage.Card{}
		err := json.Unmarshal(body, &res)
		if err != nil {
			log.Printf(cantUnmarshal, err)
			return nil, err
		}
		return res, nil
	case "password":
		res := storage.Password{}
		err := json.Unmarshal(body, &res)
		if err != nil {
			log.Printf(cantUnmarshal, err)
			return nil, err
		}
		return res, nil
	case "bin":
		res := storage.BinaryData{}
		err := json.Unmarshal(body, &res)
		if err != nil {
			log.Printf(cantUnmarshal, err)
			return nil, err
		}
		return res, nil
	}

	return nil, errors.New("unknown type")
}

// Read user data to database
func (h *Handler) Read(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf(cantRead, err)
		return
	}

	resType := r.Header.Get("Data-Type")
	data, err := myUnmarshal(resType, body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cook, _ := r.Cookie("User")

	res, err := h.db.Read(ctx, data, cook.Value)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// Update user data to database
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf(cantRead, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resType := r.Header.Get("Data-Type")
	data, err := myUnmarshal(resType, body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cook, _ := r.Cookie("User")

	err = h.db.Update(ctx, data, cook.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(&data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// Delete user data to database
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf(cantRead, err)
		return
	}

	resType := r.Header.Get("Data-Type")
	data, err := myUnmarshal(resType, body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cook, _ := r.Cookie("User")

	err = h.db.Delete(ctx, data, cook.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
