package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/terr-svc/internal/api/rest"
	"github.com/snisid/platform/services/terr-svc/internal/service"
)

func NewRouter(svc *service.TerritoryService, log *zap.Logger) *gin.Engine {
	return rest.NewRouter(svc, log)
}
