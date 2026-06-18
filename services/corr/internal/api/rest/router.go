package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/snisid/platform/services/corr/internal/service"
)

func NewRouter(
	caseSvc *service.CaseService,
	investigationSvc *service.InvestigationService,
	evidenceSvc *service.EvidenceService,
	wbSvc *service.WhistleblowerService,
	alertSvc *service.AlertService,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(AuditMiddleware())

	caseHandler := NewCaseHandler(caseSvc)
	investigationHandler := NewInvestigationHandler(investigationSvc)
	evidenceHandler := NewEvidenceHandler(evidenceSvc)
	wbHandler := NewWhistleblowerHandler(wbSvc)
	alertHandler := NewAlertHandler(alertSvc)
	declHandler := NewDeclarationHandler(alertSvc)

	v1 := r.Group("/api/v1/corr")
	{
		cases := v1.Group("/cases")
		{
			cases.POST("", caseHandler.Create)
			cases.GET("/active", caseHandler.ListActive)
			cases.GET("/:id", caseHandler.GetByID)
			cases.POST("/:id/investigate", investigationHandler.Start)
			cases.POST("/:id/close", investigationHandler.Close)
		}

		v1.POST("/evidence", evidenceHandler.Create)
		v1.GET("/cases/:id/evidence", evidenceHandler.ListByCase)

		v1.POST("/whistleblower", wbHandler.Submit)
		v1.GET("/whistleblower/:token", wbHandler.GetByToken)

		v1.GET("/alerts/behavioral", alertHandler.ListBehavioral)
		v1.GET("/risk-scores", alertHandler.ListRiskScores)

		v1.POST("/declarations", declHandler.Submit)
		v1.GET("/declarations/flagged", declHandler.ListFlagged)
	}

	return r
}
