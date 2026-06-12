package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// ─── Models ───────────────────────────────────────────────────────────────────

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Provider  string    `json:"provider"`
	CreatedAt time.Time `json:"created_at"`
}

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// ─── JWT ──────────────────────────────────────────────────────────────────────

func GenerateToken(userID, username string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateToken(tokenStr string) (*Claims, error) {
	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

// ─── Middleware ───────────────────────────────────────────────────────────────

type contextKey string

const UserContextKey contextKey = "user"

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		claims, err := ValidateToken(strings.TrimPrefix(header, "Bearer "))
		if err != nil {
			http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ClaimsFromContext(ctx context.Context) (*Claims, bool) {
	c, ok := ctx.Value(UserContextKey).(*Claims)
	return c, ok
}

// ─── Handler ──────────────────────────────────────────────────────────────────

type Handler struct {
	db *pgxpool.Pool
}

func NewHandler(db *pgxpool.Pool) *Handler {
	return &Handler{db: db}
}

type registerRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Username == "" || req.Password == "" {
		writeError(w, "email, username and password are required", http.StatusBadRequest)
		return
	}

	if len(req.Password) < 8 {
		writeError(w, "password must be at least 8 characters", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeError(w, "internal error", http.StatusInternalServerError)
		return
	}

	var user User
	err = h.db.QueryRow(r.Context(),
		`INSERT INTO users (email, username, password_hash)
		 VALUES ($1, $2, $3)
		 RETURNING id, email, username, provider, created_at`,
		req.Email, req.Username, string(hash),
	).Scan(&user.ID, &user.Email, &user.Username, &user.Provider, &user.CreatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "unique") {
			writeError(w, "email or username already in use", http.StatusConflict)
			return
		}
		writeError(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	token, err := GenerateToken(user.ID, user.Username)
	if err != nil {
		writeError(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	writeJSON(w, authResponse{Token: token, User: user}, http.StatusCreated)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	var user User
	var passwordHash string
	err := h.db.QueryRow(r.Context(),
		`SELECT id, email, username, provider, created_at, password_hash
		 FROM users WHERE email = $1`,
		req.Email,
	).Scan(&user.ID, &user.Email, &user.Username, &user.Provider, &user.CreatedAt, &passwordHash)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		writeError(w, "internal error", http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		writeError(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := GenerateToken(user.ID, user.Username)
	if err != nil {
		writeError(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	writeJSON(w, authResponse{Token: token, User: user}, http.StatusOK)
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var user User
	err := h.db.QueryRow(r.Context(),
		`SELECT id, email, username, provider, created_at FROM users WHERE id = $1`,
		claims.UserID,
	).Scan(&user.ID, &user.Email, &user.Username, &user.Provider, &user.CreatedAt)

	if err != nil {
		writeError(w, "user not found", http.StatusNotFound)
		return
	}

	writeJSON(w, user, http.StatusOK)
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func writeJSON(w http.ResponseWriter, v any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, msg string, status int) {
	writeJSON(w, map[string]string{"error": msg}, status)
}
