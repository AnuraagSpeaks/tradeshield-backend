package middleware

import (
	"context"
	"net/http"
	"strings"

	"tradeshield-backend/internal/pkg/response"
	"tradeshield-backend/internal/pkg/token"
)

type contextKey string

const UserContextKey = contextKey("user_claims")

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.Error(w, http.StatusUnauthorized, "Missing Authorization header")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(w, http.StatusUnauthorized, "Invalid Authorization header format. Expected Bearer <token>")
			return
		}

		claims, err := token.ValidateToken(parts[1])
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserClaims(ctx context.Context) *token.Claims {
	if claims, ok := ctx.Value(UserContextKey).(*token.Claims); ok {
		return claims
	}
	return nil
}
