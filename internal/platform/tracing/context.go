package tracing

import (
	"context"

	"github.com/google/uuid"
)

type correlationIDKey struct{}

// CorrelationHeader is the standard HTTP/gRPC header key
const CorrelationHeader = "X-Correlation-ID"

// WithCorrelationID injects a specific correlation ID into the context.
func WithCorrelationID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, correlationIDKey{}, id)
}

// ExtractCorrelationID retrieves the correlation ID from context.
// If none exists, it generates a new one.
func ExtractCorrelationID(ctx context.Context) string {
	if id, ok := ctx.Value(correlationIDKey{}).(string); ok && id != "" {
		return id
	}
	return uuid.NewString()
}

// GenerateCorrelationID unconditionally creates a new correlation ID
func GenerateCorrelationID() string {
	return uuid.NewString()
}
