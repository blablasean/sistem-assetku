package middlewares

import (
	"context"
	"net/http"
	"sistem-asetku-backend/utils"
	"strings"
)

const (
	UserContextKey = "user_claims"
	ClaimsKey      = "claims"
)

// AuthMiddleware validates JWT token from Authorization header
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.SendError(w, http.StatusUnauthorized, "Missing authorization header", "Authorization header required")
			return
		}

		// Extract bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.SendError(w, http.StatusUnauthorized, "Invalid authorization header format", "Expected 'Bearer <token>'")
			return
		}

		tokenString := parts[1]
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			utils.SendError(w, http.StatusUnauthorized, "Invalid or expired token", err.Error())
			return
		}

		// Add claims to context
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetClaimsFromContext retrieves claims from request context
func GetClaimsFromContext(r *http.Request) *utils.Claims {
	claims, ok := r.Context().Value(UserContextKey).(*utils.Claims)
	if !ok {
		return nil
	}
	return claims
}
