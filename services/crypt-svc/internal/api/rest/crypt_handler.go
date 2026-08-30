package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/snisid/platform/services/crypt-svc/internal/domain"
	"github.com/snisid/platform/services/crypt-svc/internal/service"
)

type CryptHandler struct {
	svc *service.CryptService
	log *zap.Logger
}

func NewCryptHandler(svc *service.CryptService, log *zap.Logger) *CryptHandler {
	return &CryptHandler{svc: svc, log: log}
}

func (h *CryptHandler) CheckAddress(c *gin.Context) {
	address := c.Param("address")
	if address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "address is required"})
		return
	}

	report, err := h.svc.AnalyzeWalletRisk(address)
	if err != nil {
		h.log.Error("analyze wallet risk failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, report)
}

func (h *CryptHandler) CreateWallet(c *gin.Context) {
	var req domain.CreateWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet, err := h.svc.FlagWallet(&req)
	if err != nil {
		h.log.Error("flag wallet failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, wallet)
}

func (h *CryptHandler) AddTransaction(c *gin.Context) {
	idStr := c.Param("id")
	walletID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid wallet id"})
		return
	}

	var req domain.AddTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := h.svc.AddTransaction(walletID, &req)
	if err != nil {
		h.log.Error("add transaction failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, tx)
}

func (h *CryptHandler) GetSanctionedWallets(c *gin.Context) {
	wallets, err := h.svc.GetSanctionedWallets()
	if err != nil {
		h.log.Error("get sanctioned wallets failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, wallets)
}

func (h *CryptHandler) GetWalletsByGang(c *gin.Context) {
	idStr := c.Param("id")
	gangID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid gang id"})
		return
	}

	wallets, err := h.svc.GetWalletsByGang(gangID)
	if err != nil {
		h.log.Error("get wallets by gang failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, wallets)
}

func (h *CryptHandler) GetStatsByAsset(c *gin.Context) {
	stats, err := h.svc.GetStatsByAsset()
	if err != nil {
		h.log.Error("get stats by asset failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
