package handlers

import (
	"context"
	"net/http"
	"prodyo-backend/cmd/internal/models"
	"prodyo-backend/cmd/internal/usecases"
)

type contextKey string

const UserContextKey contextKey = "user"

func AuthMiddleware(authUseCase *usecases.AuthUseCase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")

			if token == "" {
				http.Error(w, "Authorization token required", http.StatusUnauthorized)
				return
			}

			// Remove "Bearer " prefix if present
			if len(token) > 7 && token[:7] == "Bearer " {
				token = token[7:]
			}

			ctx := r.Context()
			user, err := authUseCase.ValidateSession(ctx, token)
			if err != nil {
				http.Error(w, "Invalid or expired session", http.StatusUnauthorized)
				return
			}
			ctx = context.WithValue(ctx, UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserFromContext(r *http.Request) (models.User, bool) {
	user, ok := r.Context().Value(UserContextKey).(models.User)
	return user, ok
}
