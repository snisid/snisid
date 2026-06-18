package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"fmt"
	"math"

	"github.com/gin-gonic/gin"
	"github.com/snisid/platform/internal/ml"
	"github.com/snisid/platform/internal/platform/events"
	"github.com/snisid/platform/internal/platform/logger"
	"github.com/snisid/platform/internal/service/fraud"
	"github.com/snisid/platform/internal/service/router"
	"go.uber.org/zap"
)

// adapter conforming *fraud.StateStore to ml.FeatureStore interface
type mlFeatureStoreAdapter struct {
	state *fraud.StateStore
}

func (a *mlFeatureStoreAdapter) GetVelocity(ctx context.Context, userID string) (float64, error) {
	val, err := a.state.GetState(ctx, fmt.Sprintf("snisid:features:%s:velocity", userID))
	if err != nil {
		return 0, err
	}
	var v float64
	fmt.Sscanf(val, "%f", &v)
	return math.Min(v/10.0, 1.0), nil
}

func (a *mlFeatureStoreAdapter) GetGraphRisk(ctx context.Context, userID string) (float64, error) {
	val, err := a.state.GetState(ctx, fmt.Sprintf("snisid:features:%s:graph_risk", userID))
	if err != nil {
		return 0, err
	}
	var v float64
	fmt.Sscanf(val, "%f", &v)
	return v, nil
}

// variable-level guard to enforce adapter implements ml.FeatureStore at compile time
var _ ml.FeatureStore = (*mlFeatureStoreAdapter)(nil)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	brokers := strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ",")
	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
	aiEndpoint := getEnv("AI_ENDPOINT", "http://ai-worker:8000/predict")
	port := getEnv("PORT", "8082")

	// ── Real Redis-backed components ─────────────────────────────────
	stateStore := fraud.NewStateStore(redisAddr)
	featureExtractor := ml.NewFeatureExtractor(&mlFeatureStoreAdapter{state: stateStore}, logger.Log)
	aiClient := fraud.NewDefaultAIClient(aiEndpoint)

	engine, err := fraud.NewScoringEngine(redisAddr, aiClient)
	if err != nil {
		logger.Fatal(context.Background(), "Failed to init scoring engine", err)
	}

	_ = engine.ReloadRules([]router.Rule{
		{ID: "suspicious-location", Expression: "event.metadata.location == 'untrusted'", Targets: []string{"internal"}},
		{ID: "rapid-enrollment", Expression: "event.action == 'enroll' && event.velocity > 5", Targets: []string{"internal"}},
		{ID: "identity-tampering", Expression: "event.action == 'update' && event.metadata.force == 'true'", Targets: []string{"security"}},
		{ID: "duplicate-biometric", Expression: "event.biometricDuplicate == true", Targets: []string{"internal", "security"}},
		{ID: "expired-document", Expression: "event.documentExpired == true", Targets: []string{"agency"}},
	})

	consumer := events.NewConsumer(brokers, "fraud-engine-group", "snisid.prod.fraud.v1.events")
	producer := events.NewProducer(brokers, "snisid.prod.soc.v1.alerts")
	riskProducer := events.NewProducer(brokers, "snisid.prod.risk.v1.updates")

	go func() {
		r := gin.Default()
		r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
		r.GET("/metrics", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"rules_loaded":    len(engine.Rules()),
				"ai_connected":    aiClient != nil,
				"redis_connected": stateStore != nil,
				"status":          "running",
			})
		})
		r.POST("/v1/score", func(c *gin.Context) {
			var event map[string]interface{}
			if err := c.BindJSON(&event); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			score, reason, riskLevel := engine.CalculateScore(c.Request.Context(), event)
			fv, err := featureExtractor.ExtractFeatures(c.Request.Context(), event)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"score":      score,
				"reason":     reason,
				"risk_level": riskLevel,
				"features":   fv,
			})
		})
		srv := &http.Server{Addr: ":" + port, Handler: r}
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(ctx, "fraud-engine http server failed", err)
		}
	}()

	logger.Info(ctx, "SNISID Fraud Engine starting",
		zap.String("redis", redisAddr),
		zap.String("ai", aiEndpoint),
	)

	err = consumer.Start(ctx, func(ctx context.Context, payload []byte) error {
		var event map[string]interface{}
		if err := json.Unmarshal(payload, &event); err != nil {
			return err
		}

		// 1. Extract real features from Redis-backed state
		fv, err := featureExtractor.ExtractFeatures(ctx, event)
		if err != nil {
			logger.Warn(ctx, "feature extraction failed, using degraded mode", zap.Error(err))
			fv = &ml.FeatureVector{}
		}

		// 2. Score the event (uses real Redis velocity via StateStore)
		score, reason, riskLevel := engine.CalculateScore(ctx, event)

		alertStatus := "ELEVATED_RISK"
		if riskLevel == "CRITICAL" || score > 80 {
			alertStatus = "CRITICAL"
		} else if riskLevel == "HIGH" || score > 40 {
			alertStatus = "HIGH"
		}

		if alertStatus == "CRITICAL" {
			logger.Warn(ctx, "CRITICAL FRAUD ALERT",
				zap.Int("score", score),
				zap.String("reason", reason),
				zap.String("risk_level", riskLevel),
			)
			alert := map[string]interface{}{
				"identityId": event["identityId"],
				"fraudScore": score,
				"reason":     reason,
				"riskLevel":  riskLevel,
				"features":   fv,
				"status":     "CRITICAL",
				"timestamp":  time.Now().UTC().Format(time.RFC3339),
			}
			if err := producer.Publish(ctx, "alert", alert); err != nil {
				logger.Error(ctx, "failed to publish fraud alert", err)
			}
		} else if alertStatus == "HIGH" {
			riskUpdate := map[string]interface{}{
				"identityId": event["identityId"],
				"score":      score,
				"reason":     reason,
				"riskLevel":  riskLevel,
				"features":   fv,
				"status":     "HIGH",
			}
			_ = riskProducer.Publish(ctx, "risk-update", riskUpdate)
		}

		return nil
	})

	if err != nil {
		logger.Fatal(ctx, "Fraud Engine crashed", err)
	}

	logger.Info(ctx, "SNISID Fraud Engine shutting down")
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
