package verification

import (
	"context"
	"fmt"
	"sync"

	"github.com/snisid/platform/internal/platform/logger"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type Orchestrator struct {
	connectors []Connector
	mu         sync.Mutex
}

func NewOrchestrator(connectors ...Connector) *Orchestrator {
	return &Orchestrator{
		connectors: connectors,
	}
}

func (o *Orchestrator) VerifyIdentity(ctx context.Context, data map[string]interface{}) (map[string]Result, error) {
	g, ctx := errgroup.WithContext(ctx)
	results := make(map[string]Result)
	var mu sync.Mutex

	logger.Info(ctx, "Starting multi-factor identity verification", zap.Int("connector_count", len(o.connectors)))

	for _, conn := range o.connectors {
		c := conn // Capture loop variable
		g.Go(func() error {
			res, err := c.Verify(ctx, data)
			if err != nil {
				logger.Error(ctx, fmt.Sprintf("Connector %s failed", c.Name()), err)
				return err
			}

			mu.Lock()
			results[c.Name()] = res
			mu.Unlock()
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("verification orchestration failed: %w", err)
	}

	return results, nil
}
