package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/vehicle-criminal-svc/internal/domain"
)

type contextKey string

const (
	userIDKey   contextKey = "user_id"
	userUnitKey contextKey = "user_unit"
	userRoleKey contextKey = "user_role"
)

func AuditMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(wrapped, r)

		_ = domain.AuditEvent{
			AgentID:    r.Context().Value(userIDKey),
			Unit:       r.Context().Value(userUnitKey),
			Action:     r.Method + " " + r.URL.Path,
			ResourceID: r.URL.Query().Get("id"),
			ClientIP:   r.RemoteAddr,
			Timestamp:  start,
		}
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		userUnit := r.Header.Get("X-User-Unit")
		userRole := r.Header.Get("X-User-Role")

		if userID == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, parseUUID(userID))
		ctx = context.WithValue(ctx, userUnitKey, userUnit)
		ctx = context.WithValue(ctx, userRoleKey, userRole)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func parseUUID(s string) uuid.UUID {
	id, _ := uuid.Parse(s)
	return id
}
