module github.com/snisid/platform/services/nexus

go 1.26.0

require (
	github.com/google/uuid v1.6.0
	github.com/snisid/platform/backend v0.0.0
	github.com/stretchr/testify v1.11.1
	go.opentelemetry.io/otel v1.43.0
	go.opentelemetry.io/otel/trace v1.43.0
	go.uber.org/zap v1.27.1
	google.golang.org/grpc v1.80.0
)

replace github.com/snisid/platform/backend => ../../
