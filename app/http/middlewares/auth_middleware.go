package middleware

import (
	"context"
	"net/http"
	"strings"

	helper "app/helpers"
)

type contextKey string

const userIDContextKey contextKey = "userID"

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := bearerToken(r)
		if !ok {
			helper.WriteResponse(w, http.StatusUnauthorized, "missing or invalid authorization header", nil)
			return
		}

		claims, err := helper.ValidateToken(token)
		if err != nil {
			helper.WriteResponse(w, http.StatusUnauthorized, "invalid or expired token", nil)
			return
		}

		ctx := context.WithValue(r.Context(), userIDContextKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserID(r *http.Request) (int64, bool) {
	id, ok := r.Context().Value(userIDContextKey).(int64)
	return id, ok
}

func bearerToken(r *http.Request) (string, bool) {
	const prefix = "Bearer "

	header := r.Header.Get("Authorization")
	if !strings.HasPrefix(header, prefix) {
		return "", false
	}

	token := strings.TrimSpace(strings.TrimPrefix(header, prefix))
	if token == "" {
		return "", false
	}

	return token, true
}
