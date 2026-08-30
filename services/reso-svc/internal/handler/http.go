package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/reso-svc/internal/api/rest"
	"github.com/snisid/platform/services/reso-svc/internal/service"
)

func NewRouter(svc *service.ResoService, log *zap.Logger) *gin.Engine {
	return rest.SetupRouter(svc, log)
}
