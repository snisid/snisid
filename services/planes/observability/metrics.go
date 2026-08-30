package observability

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var (
	RiskRequests = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "snisid_risk_requests_total",
			Help: "Total number of risk evaluations processed",
		},
	)

	RiskLatency = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "snisid_risk_latency_seconds",
			Help:    "Latency of risk evaluation requests",
			Buckets: prometheus.DefBuckets,
		},
	)
)

var Tracer = otel.Tracer("snisid-risk-engine")

func StartSpan(ctx trace.Context, name string) (trace.Context, trace.Span) {
	return Tracer.Start(ctx, name)
}
