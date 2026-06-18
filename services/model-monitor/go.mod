module github.com/snisid/platform/services/model-monitor

go 1.26.0

require (
	github.com/prometheus/client_golang v1.20.0
	github.com/snisid/platform/backend v0.0.0
	github.com/stretchr/testify v1.11.1
)

replace github.com/snisid/platform/backend => ../../
