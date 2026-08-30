module github.com/snisid/platform/services/planes

go 1.26.0

require (
	github.com/prometheus/client_golang v1.20.0
	github.com/segmentio/kafka-go v0.4.47
	github.com/snisid/platform/backend v0.0.0
	github.com/stretchr/testify v1.11.1
	go.opentelemetry.io/otel v1.43.0
	go.opentelemetry.io/otel/trace v1.43.0
)

replace github.com/snisid/platform/backend => ../../
