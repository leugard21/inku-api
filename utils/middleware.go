package utils

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

type contextKey string

const userIDKey contextKey = "userID"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			WriteError(w, http.StatusUnauthorized, errors.New("missing authorization header"))
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			WriteError(w, http.StatusUnauthorized, errors.New("invalid authorization format"))
			return
		}

		claims, err := ParseToken(parts[1])
		if err != nil {
			WriteError(w, http.StatusUnauthorized, errors.New("invalid or expired token"))
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(r *http.Request) (int64, bool) {
	id, ok := r.Context().Value(userIDKey).(int64)
	return id, ok
}
