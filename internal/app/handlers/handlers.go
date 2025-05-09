package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"net/http"
	"regexp"
	"time"
	"zkKYC-backend/internal/app/config"
	"zkKYC-backend/internal/app/db"
	"zkKYC-backend/internal/app/helpers"
	"zkKYC-backend/internal/app/storage"
)

// CustomClaims
type CustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type ZkKYCHandler struct {
	Repo storage.Repository
	cfg  config.Config
	DB   *sql.DB
}

// Create new instance of ZkKYCHandler
func NewZkKYCHandler(c config.Config) *ZkKYCHandler {
	h := &ZkKYCHandler{
		cfg: c,
		DB:  db.NewDBConnection(c),
	}

	if c.DatabaseDSN != "" {
		h.Repo = storage.NewDBStorage(h.DB)
	}

	return h
}

// API Endpoint for creating user
func (h *ZkKYCHandler) APICreateUser(w http.ResponseWriter, r *http.Request) {

	fmt.Println("[+]POST /api/user")
	input := storage.User{}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{}"))
		return
	}

	did := helpers.CreateDid(input)

	input.DID = did

	err := h.Repo.Add(r.Context(), &input)
	var pge *pgconn.PgError
	if errors.As(err, &pge) {
		if pge.Code == pgerrcode.UniqueViolation {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			u, ok := h.Repo.Get(r.Context(), input.EthAddress)
			if ok {
				input = *u.(*storage.User)
			}
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
	}

	result := struct {
		Id  int    `json:"id"`
		Did string `json:"did"`
	}{Id: input.ID, Did: input.DID}

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)

	encoder.Encode(result)
}

// API Endpoint for getting exiting user
func (h *ZkKYCHandler) APIGetExitingUser(w http.ResponseWriter, r *http.Request) {

	ethAddress := chi.URLParam(r, "eth")

	fmt.Printf("[+]GET /api/user/%s\n", ethAddress)

	re := regexp.MustCompile(`^[a-zA-Z0-9]+$`)

	if !re.MatchString(ethAddress) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u, ok := h.Repo.Get(r.Context(), ethAddress)

	if ok {

		u := u.(*storage.User)

		result := struct {
			Id  int    `json:"id"`
			Did string `json:"did"`
		}{Id: u.ID, Did: u.DID}

		encoder := json.NewEncoder(w)
		encoder.SetEscapeHTML(false)

		encoder.Encode(result)

		return
	}

	w.WriteHeader(http.StatusNotFound)

	return

}

// Login handler for authentication
func (h *ZkKYCHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("[+]POST /api/regulator/login")

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username != "admin" || password != "admin" {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	claims := CustomClaims{
		UserID: "1",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(h.cfg.JWTSecret))
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(tokenString))
}

// API Endpoint for getting exiting user for regulator
func (h *ZkKYCHandler) APIGetExitingUserForRegulator(w http.ResponseWriter, r *http.Request) {

	val := r.Context().Value("user_id")
	fmt.Println(val)

	ethAddress := chi.URLParam(r, "eth")

	fmt.Printf("[+]GET /api/regulator/user/%s\n", ethAddress)

	re := regexp.MustCompile(`^[a-zA-Z0-9]+$`)

	if !re.MatchString(ethAddress) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u, ok := h.Repo.Get(r.Context(), ethAddress)

	if ok {

		u := u.(*storage.User)

		encoder := json.NewEncoder(w)
		encoder.SetEscapeHTML(false)

		encoder.Encode(u)

		return
	}

	w.WriteHeader(http.StatusNotFound)

	return

}
