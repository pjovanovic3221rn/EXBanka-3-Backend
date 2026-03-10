package middleware

import (
	"context"
	"net/http"
	"strings"

	"auth-service/services"
)

type contextKey string

const (
	ContextCredentialIDKey contextKey = "credential_id"
	ContextEmployeeIDKey   contextKey = "employee_id"
	ContextEmailKey        contextKey = "email"
)

func AuthMiddleware(jwtService *services.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			claims, err := jwtService.ValidateAccessToken(tokenString)
			if err != nil {
				http.Error(w, "invalid or expired access token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ContextCredentialIDKey, claims.CredentialID)
			ctx = context.WithValue(ctx, ContextEmployeeIDKey, claims.EmployeeID)
			ctx = context.WithValue(ctx, ContextEmailKey, claims.Email)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}