package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/snisid/platform/services/corr/internal/domain"
	"github.com/snisid/platform/services/corr/internal/service"
)

type DeclarationHandler struct {
	svc *service.AlertService
}

func NewDeclarationHandler(svc *service.AlertService) *DeclarationHandler {
	return &DeclarationHandler{svc: svc}
}

func (h *DeclarationHandler) Submit(c *gin.Context) {
	var req domain.CreateDeclarationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "requête invalide"})
		return
	}
	now := req.RealEstateUSD + req.VehiclesUSD + req.BankAccountsUSD + req.OtherAssetsUSD
	unexplained := now - req.KnownSalaryAnnualUSD
	if unexplained < 0 {
		unexplained = 0
	}
	d := &domain.AssetDeclaration{
		DeclarationID:        uuid.New(),
		OfficerSnisidID:      req.OfficerSnisidID,
		DeclarationYear:      req.DeclarationYear,
		RealEstateUSD:        req.RealEstateUSD,
		VehiclesUSD:          req.VehiclesUSD,
		BankAccountsUSD:      req.BankAccountsUSD,
		OtherAssetsUSD:       req.OtherAssetsUSD,
		TotalAssetsUSD:       now,
		KnownSalaryAnnualUSD: req.KnownSalaryAnnualUSD,
		UnexplainedWealthUSD: unexplained,
		IsFlagged:            unexplained > req.KnownSalaryAnnualUSD*3,
	}
	if err := h.svc.CreateDeclaration(c.Request.Context(), d); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, d)
}

func (h *DeclarationHandler) ListFlagged(c *gin.Context) {
	declarations, err := h.svc.ListFlaggedDeclarations(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur requête"})
		return
	}
	c.JSON(http.StatusOK, declarations)
}
