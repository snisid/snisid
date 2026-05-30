package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/snisid/platform/backend/internal/domain/orchestrator"
	"github.com/snisid/platform/backend/internal/platform/events"
	"github.com/snisid/platform/backend/internal/platform/logger"
)

func main() {
	broker := getEnv("KAFKA_BROKER", "localhost:9092")

	// Producer for workflow output events
	producer := events.NewProducer([]string{broker}, "identity.workflow.completed")
	defer producer.Close()

	manager := orchestrator.NewWorkflowManager(producer)

	// Consumers
	identityConsumer := events.NewConsumer([]string{broker}, "orchestrator-group", "identity.created")
	defer identityConsumer.Close()

	fraudConsumer := events.NewConsumer([]string{broker}, "orchestrator-group", "fraud.scored")
	defer fraudConsumer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		logger.Info("shutting down orchestrator...", nil)
		cancel()
	}()

	// Run consumers in parallel
	go func() {
		if err := identityConsumer.Read(ctx, manager.HandleIdentityCreated); err != nil && err != context.Canceled {
			logger.Error("identity consumer error", err)
		}
	}()

	go func() {
		if err := fraudConsumer.Read(ctx, manager.HandleFraudScored); err != nil && err != context.Canceled {
			logger.Error("fraud consumer error", err)
		}
	}()

	logger.Info("orchestrator started", nil)
	<-ctx.Done()
}

func getEnv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}
