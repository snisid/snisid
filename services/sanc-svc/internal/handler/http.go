package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/sanc-svc/internal/api/rest"
	"github.com/snisid/platform/services/sanc-svc/internal/service"
)

func NewRouter(svc *service.SanctionsService, log *zap.Logger) *gin.Engine {
	return rest.SetupRouter(svc, log)
}
