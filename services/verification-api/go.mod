module github.com/snisid/platform/services/verification-api

go 1.26.0

require (
	github.com/google/uuid v1.6.0
	github.com/snisid/platform/backend v0.0.0
	github.com/stretchr/testify v1.11.1
	go.uber.org/zap v1.27.1
)

replace github.com/snisid/platform/backend => ../../
