package handlers

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/mike-testut/task-api/internal/httpjson"
	"github.com/mike-testut/task-api/internal/service"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		log.Printf(
			"%s %s %v",
			r.Method,
			r.URL.Path,
			time.Since(start),
		)
	})
}

type AuthMiddleware struct {
	authService *service.AuthService
}

func NewAuthMiddleware(as *service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authService: as}
}

type contextKey string

const userContextKey = contextKey("user_id")

func (am *AuthMiddleware) Protect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			httpjson.ErrorJSON(w, http.StatusUnauthorized, "authorization header required")
			return
		}
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			httpjson.ErrorJSON(w, http.StatusUnauthorized, "Invalid Authorization header format")
			return
		}

		tokenString := headerParts[1]

		claims, err := am.authService.ValidateToken(tokenString)
		if err != nil {
			httpjson.ErrorJSON(w, http.StatusUnauthorized, err.Error())
			return
		}
		ctx := context.WithValue(r.Context(), userContextKey, claims.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
