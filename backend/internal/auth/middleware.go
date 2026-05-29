package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"github.com/ZyphenSVC/SecureOps/backend/internal/httpx"
)

type contextKey string

const UserIDContextKey contextKey = "userID"

func (s *Service) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			httpx.WriteError(w, http.StatusUnauthorized, "missing authorization header")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			httpx.WriteError(w, http.StatusUnauthorized, "invalid authorization header")
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			return s.jwtSecret, nil
		})
		if err != nil || !token.Valid {
			httpx.WriteError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "invalid token claims")
			return
		}

		userID, ok := claims["sub"].(string)
		if !ok || userID == "" {
			httpx.WriteError(w, http.StatusUnauthorized, "invalid token subject")
			return
		}

		ctx := context.WithValue(r.Context(), UserIDContextKey, userID)
		next(w, r.WithContext(ctx))
	}
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDContextKey).(string)
	return userID, ok
}

func (s *Service) RequirePermission(permissionKey string, next http.HandlerFunc) http.HandlerFunc {
	return s.Authenticate(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := UserIDFromContext(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "unauthenticated")
			return
		}

		allowed, err := s.rbac.UserHasPermission(r.Context(), userID, permissionKey)
		if err != nil {
			httpx.WriteError(w, http.StatusInternalServerError, "permission check failed")
			return
		}

		if !allowed {
			httpx.WriteError(w, http.StatusForbidden, "forbidden")
			return
		}

		next(w, r)
	})
}
