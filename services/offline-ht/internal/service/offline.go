package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/offline-ht/internal/domain"
	"github.com/snisid/offline-ht/internal/kafka"
	"github.com/snisid/offline-ht/internal/repository"
)

type OfflineService struct {
	repo     repository.Repository
	producer *kafka.Producer
}

func NewOfflineService(repo repository.Repository, producer *kafka.Producer) *OfflineService {
	return &OfflineService{repo: repo, producer: producer}
}

func (s *OfflineService) PushQueue(ctx context.Context, req domain.PushQueueRequest) (*domain.SyncQueueItem, error) {
	terminalID, err := uuid.Parse(req.TerminalID)
	if err != nil {
		return nil, fmt.Errorf("invalid terminal_id: %w", err)
	}

	now := time.Now().UTC()
	item := &domain.SyncQueueItem{
		ID:         uuid.New(),
		TerminalID: terminalID,
		EntityType: req.EntityType,
		EntityID:   req.EntityID,
		Action:     req.Action,
		Payload:    req.Payload,
		Status:     domain.SyncPending,
		RetryCount: 0,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.repo.PushQueue(ctx, item); err != nil {
		return nil, fmt.Errorf("push queue: %w", err)
	}

	s.publishEvent(ctx, "offline.queue.pushed", item.ID.String(), "sync_queue", item)
	return item, nil
}

func (s *OfflineService) SyncTerminal(ctx context.Context, terminalIDStr string) ([]domain.SyncQueueItem, error) {
	terminalID, err := uuid.Parse(terminalIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid terminal_id: %w", err)
	}

	items, err := s.repo.SyncTerminal(ctx, terminalID)
	if err != nil {
		return nil, fmt.Errorf("sync terminal: %w", err)
	}

	for i := range items {
		items[i].Status = domain.SyncSyncing
		if err := s.repo.UpdateQueueItemStatus(ctx, items[i].ID, domain.SyncSyncing, nil); err != nil {
			log.Printf("failed to update item %s status: %v", items[i].ID, err)
		}
	}

	s.markTerminalSynced(ctx, terminalID)
	s.publishEvent(ctx, "offline.sync.started", terminalID.String(), "terminal", map[string]any{
		"terminal_id": terminalID.String(),
		"item_count":  len(items),
	})
	return items, nil
}

func (s *OfflineService) GetConflicts(ctx context.Context) ([]domain.SyncQueueItem, error) {
	return s.repo.GetConflictItems(ctx)
}

func (s *OfflineService) GetTerminalsStatus(ctx context.Context) ([]domain.OfflineTerminal, error) {
	return s.repo.GetTerminalsStatus(ctx)
}

func (s *OfflineService) markTerminalSynced(ctx context.Context, terminalID uuid.UUID) {
	now := time.Now().UTC()
	t := &domain.OfflineTerminal{
		ID:         terminalID,
		LastSyncAt: &now,
		IsOnline:   true,
		UpdatedAt:  now,
	}
	if err := s.repo.UpsertTerminal(ctx, t); err != nil {
		log.Printf("failed to upsert terminal %s: %v", terminalID, err)
	}
}

func (s *OfflineService) publishEvent(ctx context.Context, eventType, entityID, entityType string, data any) {
	if s.producer == nil {
		return
	}
	evt := kafka.Event{
		EventType:  eventType,
		EntityID:   entityID,
		EntityType: entityType,
		Timestamp:  time.Now().UTC(),
		Data:       data,
	}
	if err := s.producer.Publish(ctx, evt); err != nil {
		log.Printf("failed to publish event: %v", err)
	}
}
