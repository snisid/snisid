package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/snisid/platform/internal/platform/logger"
	"github.com/snisid/platform/services/governance-intervention/approver"
	"github.com/snisid/platform/services/governance-intervention/proposer"
	"go.uber.org/zap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	port := getEnv("PORT", "8093")
	proposalEngine := &proposer.ProposalEngine{}
	approvalGate := &approver.ApprovalGate{}

	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })

	r.POST("/v1/propose", func(c *gin.Context) {
		var req struct {
			Target    string   `json:"target"`
			RiskScore float64  `json:"risk_score"`
			Reasons   []string `json:"reasons"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		intervention := proposalEngine.Propose(req.Target, req.RiskScore, req.Reasons)
		c.JSON(http.StatusOK, intervention)
	})

	r.POST("/v1/approve", func(c *gin.Context) {
		var req struct {
			Intervention  proposer.Intervention `json:"intervention"`
			UserRole      string                `json:"user_role"`
			Justification string                `json:"justification"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := approvalGate.Authorize(&req.Intervention, req.UserRole, req.Justification); err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, req.Intervention)
	})

	r.POST("/v1/reject", func(c *gin.Context) {
		var req struct {
			Intervention proposer.Intervention `json:"intervention"`
			Reason       string                `json:"reason"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		approvalGate.Reject(&req.Intervention, req.Reason)
		c.JSON(http.StatusOK, req.Intervention)
	})

	srv := &http.Server{Addr: ":" + port, Handler: r}

	go func() {
		logger.Info(ctx, "governance-intervention starting", zap.String("port", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(ctx, "governance-intervention http server failed", err)
		}
	}()

	<-ctx.Done()
	logger.Info(ctx, "governance-intervention shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(shutdownCtx)
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
