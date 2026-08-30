module github.com/snisid/platform/services/intelligence

go 1.26.0

require (
	github.com/gin-gonic/gin v1.10.0
	github.com/redis/go-redis/v9 v9.6.1
	github.com/snisid/platform/backend v0.0.0
	github.com/stretchr/testify v1.11.1
	go.uber.org/zap v1.27.1
)

replace github.com/snisid/platform/backend => ../../
