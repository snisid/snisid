module github.com/snisid/platform/services/enrollment-service

go 1.26.0

require (
	github.com/gin-gonic/gin v1.10.0
	github.com/snisid/platform/backend v0.0.0
	github.com/stretchr/testify v1.11.1
	go.uber.org/zap v1.27.1
	gorm.io/driver/postgres v1.6.0
	gorm.io/driver/sqlite v1.5.6
	gorm.io/gorm v1.31.1
)

replace github.com/snisid/platform/backend => ../../
