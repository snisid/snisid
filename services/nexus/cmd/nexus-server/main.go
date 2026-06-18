package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	nexusv1 "github.com/snisid/platform/services/nexus/api/proto/nexus/v1"
	"github.com/snisid/platform/services/nexus/core"
	"github.com/snisid/platform/services/nexus/internal/agent"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	state := core.NewInMemoryState()
	orchestrator := core.NewOrchestrator(10, logger, state)

	// Register a default Kai agent
	kai := agent.NewKaiAgent("kai-01", logger)
	orchestrator.RegisterAgent(kai)

	// Start orchestrator workers
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	orchestrator.Start(ctx)

	s := grpc.NewServer()
	nexusv1.RegisterNexusServiceServer(s, orchestrator)
	reflection.Register(s)

	logger.Info("nexus server starting", zap.String("addr", lis.Addr().String()))

	go func() {
		if err := s.Serve(lis); err != nil {
			logger.Fatal("failed to serve", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down nexus server")
	s.GracefulStop()
	cancel()
}
