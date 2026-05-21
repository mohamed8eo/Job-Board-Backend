package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/mohamed8eo/jobBoard/internal/auth"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
	RoleKey   contextKey = "role"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			tokenSecret := os.Getenv("SECRET_JWT")

			claims, err := auth.ValidateJwt(tokenString, tokenSecret)
			if err == nil {
				ctx := context.WithValue(r.Context(), UserIDKey, claims.Subject)
				ctx = context.WithValue(ctx, RoleKey, claims.Role)
				r = r.WithContext(ctx)
			}
		}

		next.ServeHTTP(w, r)
	})
}

// RequireAuth checks if user is authenticated (any role)
func RequireAuth(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok || userID == "" {
		return "", errors.New("authentication required")
	}
	return userID, nil
}

// RequireCompany checks if user is authenticated and has company role
func RequireCompany(ctx context.Context) (string, error) {
	userID, err := RequireAuth(ctx)
	if err != nil {
		return "", err
	}

	role, ok := ctx.Value(RoleKey).(string)
	if !ok || role == "" {
		return "", errors.New("role information not found")
	}

	if role != "company" {
		return "", fmt.Errorf("access denied: this operation requires company role")
	}

	return userID, nil
}

// RequireUser checks if user is authenticated and has user role
func RequireUser(ctx context.Context) (string, error) {
	userID, err := RequireAuth(ctx)
	if err != nil {
		return "", err
	}

	role, ok := ctx.Value(RoleKey).(string)
	if !ok || role == "" {
		return "", errors.New("role information not found")
	}

	if role != "user" {
		return "", fmt.Errorf("access denied: this operation requires user role")
	}

	return userID, nil
}

// GetUserIDIfExists returns userID if exists, false otherwise
func GetUserIDIfExists(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok && userID != ""
}

// GetRoleIfExists returns role if exists, false otherwise
func GetRoleIfExists(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(RoleKey).(string)
	return role, ok && role != ""
}
