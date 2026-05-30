package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/snisid/platform/backend/internal/platform/events"
	"github.com/snisid/platform/backend/internal/platform/logger"
	"github.com/snisid/platform/backend/internal/service/fraud"
	"github.com/snisid/platform/backend/internal/service/router"
	"go.uber.org/zap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	brokers := []string{os.Getenv("KAFKA_BROKERS")}
	if len(brokers[0]) == 0 {
		brokers = []string{"localhost:9092"}
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	// 1. Initialize Scoring Engine
	engine, err := fraud.NewScoringEngine(redisAddr)
	if err != nil {
		logger.Fatal(context.Background(), "Failed to init scoring engine", err)
	}

	// Load some heuristic rules
	_ = engine.ReloadRules([]router.Rule{
		{
			ID:         "suspicious-location",
			Expression: "event.metadata.location == 'untrusted'",
			Targets:    []string{"internal"},
		},
	})

	// 2. Initialize Producers/Consumers
	consumer := events.NewConsumer(brokers, "fraud-engine-group", "snisid.prod.fraud.v1.events")
	producer := events.NewProducer(brokers, "snisid.prod.soc.v1.alerts")

	logger.Info(context.Background(), "SNISID Fraud Engine starting...", zap.String("redis", redisAddr))

	// 3. Start Processing Loop
	err = consumer.Start(ctx, func(ctx context.Context, payload []byte) error {
		var event map[string]interface{}
		if err := json.Unmarshal(payload, &event); err != nil {
			return err
		}

		score, reason := engine.CalculateScore(ctx, event)
		
		if score > 80 {
			logger.Warn(ctx, "FRAUD ALERT GENERATED", zap.Int("score", score), zap.String("reason", reason))
			
			alert := map[string]interface{}{
				"identityId": event["identityId"],
				"fraudScore": score,
				"reason":     reason,
				"status":     "CRITICAL",
			}
			return producer.Publish(ctx, "alert", alert)
		}

		return nil
	})

	if err != nil {
		logger.Fatal(context.Background(), "Fraud Engine crashed", err)
	}

	logger.Info(context.Background(), "SNISID Fraud Engine shutting down")
}
