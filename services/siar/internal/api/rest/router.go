package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/snisid/platform/services/siar/internal/service"
)

func NewRouter(
	firearmSvc *service.FirearmService,
	licenseSvc *service.LicenseService,
	transferSvc *service.TransferService,
	dealerSvc *service.DealerService,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(AuditMiddleware())

	firearmHandler := NewFirearmHandler(firearmSvc)
	licenseHandler := NewLicenseHandler(licenseSvc)
	seizureHandler := NewSeizureHandler(firearmSvc, transferSvc)
	dealerHandler := NewDealerHandler(dealerSvc)

	v1 := r.Group("/api/v1/siar")
	{
		firearms := v1.Group("/firearms")
		{
			firearms.POST("", firearmHandler.Create)
			firearms.GET("/:id", firearmHandler.GetByID)
		}

		v1.GET("/check/serial/:sn", firearmHandler.CheckSerial)
		v1.GET("/stats/by-type", firearmHandler.StatsByType)

		v1.POST("/seizures", seizureHandler.ReportSeizure)
		v1.POST("/stolen", seizureHandler.ReportStolen)

		licenses := v1.Group("/licenses")
		{
			licenses.POST("", licenseHandler.Create)
			licenses.GET("/:person", licenseHandler.GetByPerson)
		}

		dealers := v1.Group("/dealers")
		{
			dealers.POST("", dealerHandler.Create)
			dealers.GET("", dealerHandler.List)
			dealers.GET("/:id", dealerHandler.GetByID)
		}
	}

	return r
}
