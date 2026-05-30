package analytics

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/snisid/platform/backend/internal/platform/logger"
	"go.uber.org/zap"
)

type ClickHouseBridge struct {
	conn      clickhouse.Conn
	batchSize int
	buffer    []FusedEvent
	mu        sync.Mutex
	flushInt  time.Duration
}

func NewClickHouseBridge(addr string) (*ClickHouseBridge, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{addr},
		Auth: clickhouse.Auth{
			Database: "snisid",
			Username: "default",
			Password: "",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to clickhouse: %w", err)
	}

	return &ClickHouseBridge{
		conn:      conn,
		batchSize: 1000,
		buffer:    make([]FusedEvent, 0, 1000),
		flushInt:  5 * time.Second,
	}, nil
}

func (b *ClickHouseBridge) Start(ctx context.Context) {
	ticker := time.NewTicker(b.flushInt)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			b.Flush(context.Background())
			return
		case <-ticker.C:
			b.Flush(ctx)
		}
	}
}

func (b *ClickHouseBridge) BufferEvent(event FusedEvent) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.buffer = append(b.buffer, event)
	if len(b.buffer) >= b.batchSize {
		go b.Flush(context.Background())
	}
}

func (b *ClickHouseBridge) Flush(ctx context.Context) {
	b.mu.Lock()
	if len(b.buffer) == 0 {
		b.mu.Unlock()
		return
	}
	events := b.buffer
	b.buffer = make([]FusedEvent, 0, b.batchSize)
	b.mu.Unlock()

	logger.Info(ctx, "Flushing events to ClickHouse", zap.Int("count", len(events)))

	batch, err := b.conn.PrepareBatch(ctx, "INSERT INTO fused_events")
	if err != nil {
		logger.Error(ctx, "Failed to prepare clickhouse batch", err)
		return
	}

	for _, e := range events {
		err := batch.Append(
			e.EventID,
			e.CorrelationID,
			e.Type,
			e.Source,
			e.Timestamp,
			fmt.Sprintf("%v", e.Data),
		)
		if err != nil {
			logger.Error(ctx, "Failed to append to batch", err)
		}
	}

	if err := batch.Send(); err != nil {
		logger.Error(ctx, "Failed to send clickhouse batch", err)
	}
}
