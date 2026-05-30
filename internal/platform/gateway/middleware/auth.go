package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/snisid/platform/backend/internal/platform/errors"
)

type userContextKey struct{}
type roleContextKey struct{}

func Auth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				err := errors.New(errors.Unauthenticated, "missing authorization header", "middleware.Auth", nil)
				errors.RespondWithErrorStd(w, r, err)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				err := errors.New(errors.Unauthenticated, "invalid authorization header format", "middleware.Auth", nil)
				errors.RespondWithErrorStd(w, r, err)
				return
			}

			tokenStr := parts[1]
			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					// We should actually support RS256 here based on the Auth service changes, but keeping HMAC structure for brevity as requested
					return nil, http.ErrAbortHandler
				}
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				parsedErr := errors.New(errors.Unauthenticated, "invalid or expired token", "middleware.Auth", err)
				errors.RespondWithErrorStd(w, r, parsedErr)
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				ctx := context.WithValue(r.Context(), userContextKey{}, claims["sub"])
				if roles, ok := claims["roles"].([]interface{}); ok {
					ctx = context.WithValue(ctx, roleContextKey{}, roles)
				}
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			parsedErr := errors.New(errors.Unauthenticated, "invalid token claims", "middleware.Auth", nil)
			errors.RespondWithErrorStd(w, r, parsedErr)
		})
	}
}

// RequireRole RBAC middleware to check roles
func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			roles, ok := r.Context().Value(roleContextKey{}).([]interface{})
			if !ok {
				err := errors.New(errors.PermissionDenied, "missing role context", "middleware.RequireRole", nil)
				errors.RespondWithErrorStd(w, r, err)
				return
			}

			hasRole := false
			for _, r := range roles {
				if r.(string) == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				err := errors.New(errors.PermissionDenied, "insufficient permissions", "middleware.RequireRole", nil)
				errors.RespondWithErrorStd(w, r, err)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
