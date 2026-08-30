package rest

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/crypt-svc/internal/service"
)

func SetupRouter(svc *service.CryptService, log *zap.Logger) *gin.Engine {
	handler := NewCryptHandler(svc, log)
	r := gin.Default()

	api := r.Group("/api/v1/crypt")
	{
		api.GET("/check/:address", handler.CheckAddress)
		api.POST("/wallets", handler.CreateWallet)
		api.POST("/wallets/:id/transactions", handler.AddTransaction)
		api.GET("/wallets/sanctioned", handler.GetSanctionedWallets)
		api.GET("/wallets/gang/:id", handler.GetWalletsByGang)
		api.GET("/stats/by-asset", handler.GetStatsByAsset)
	}

	return r
}
