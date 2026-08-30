package middleware

import (
	"net/http"
	"time"

	"github.com/snisid/platform/internal/platform/logger"
	"github.com/snisid/platform/internal/platform/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type responseRecorder struct {
	http.ResponseWriter
	status int
	size   int
}

func (rw *responseRecorder) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *responseRecorder) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

func Telemetry() func(next http.Handler) http.Handler {
	tracer := otel.Tracer("api-gateway")
	propagator := otel.GetTextMapPropagator()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Extract correlation ID or generate new
			correlationID := r.Header.Get(tracing.CorrelationHeader)
			if correlationID == "" {
				correlationID = tracing.GenerateCorrelationID()
			}

			// Add Correlation ID to response header
			w.Header().Set(tracing.CorrelationHeader, correlationID)

			// Setup context with Correlation ID and OTel propagation
			ctx := tracing.WithCorrelationID(r.Context(), correlationID)
			ctx = propagator.Extract(ctx, propagation.HeaderCarrier(r.Header))

			// Start Span
			ctx, span := tracer.Start(ctx, r.URL.Path, trace.WithAttributes(
				attribute.String("http.method", r.Method),
				attribute.String("http.url", r.URL.String()),
				attribute.String("correlation.id", correlationID),
			))
			defer span.End()

			// Update request with new context
			r = r.WithContext(ctx)

			// Record Response
			rw := &responseRecorder{ResponseWriter: w, status: http.StatusOK}

			// Next
			next.ServeHTTP(rw, r)

			// Log Request
			duration := time.Since(start)
			span.SetAttributes(attribute.Int("http.status_code", rw.status))

			// Avoid verbose logging for health checks
			if r.URL.Path != "/healthz" {
				logger.Info(ctx, "HTTP Access Log",
					zap.String("method", r.Method),
					zap.String("path", r.URL.Path),
					zap.Int("status_code", rw.status),
					zap.String("duration", duration.String()),
					zap.String("user_agent", r.UserAgent()),
					zap.String("client_ip", r.RemoteAddr),
				)
			}
		})
	}
}
