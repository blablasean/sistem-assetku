package middlewares

import (
	"context"
	"net/http"
	"sistem-asetku-backend/models"
	"sistem-asetku-backend/utils"
	"strings"
	"time"

	"gorm.io/gorm"
)

const (
	UserContextKey = "user_claims"
	ClaimsKey      = "claims"
)

// AuthMiddlewareWithDB validates JWT token and checks active session token in DB
func AuthMiddlewareWithDB(db *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.SendError(w, http.StatusUnauthorized, "Missing authorization header", "Authorization header required")
				return
			}

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

			// Single Active Session Validation: Verify token matches database active token
			if db != nil && claims.UserID > 0 {
				var user models.User
				if err := db.Select("active_token, is_logged_in").First(&user, claims.UserID).Error; err == nil {
					if user.ActiveToken != "" && user.ActiveToken != tokenString {
						utils.SendError(w, http.StatusUnauthorized, "Sesi Anda telah berakhir karena akun ini sedang aktif di perangkat lain.", "Session invalidated")
						return
					}
					// Update last_seen_at timestamp and ensure is_logged_in = 1 on valid requests
					now := time.Now()
					db.Exec("UPDATE users SET is_logged_in = 1, last_seen_at = ? WHERE id = ?", now, claims.UserID)
				}
			}

			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AuthMiddleware(next http.Handler) http.Handler {
	return AuthMiddlewareWithDB(nil)(next)
}

// GetClaimsFromContext retrieves claims from request context
func GetClaimsFromContext(r *http.Request) *utils.Claims {
	claims, ok := r.Context().Value(UserContextKey).(*utils.Claims)
	if !ok {
		return nil
	}
	return claims
}
