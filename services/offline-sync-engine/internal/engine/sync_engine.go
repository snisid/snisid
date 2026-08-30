package engine

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/platform/services/offline-sync-engine/internal/models"
	"gorm.io/gorm"
)

type SyncCallback func(event *models.OfflineEvent) error

type SyncResult struct {
	Total     int               `json:"total"`
	Synced    int               `json:"synced"`
	Conflicts int               `json:"conflicts"`
	Failed    int               `json:"failed"`
	Skipped   int               `json:"skipped"`
	Details   []SyncEventResult `json:"details"`
}

type SyncEventResult struct {
	EventID  string `json:"event_id"`
	Status   string `json:"status"`
	Conflict string `json:"conflict,omitempty"`
	Error    string `json:"error,omitempty"`
}

type QueueStatus struct {
	Pending  int64 `json:"pending"`
	Synced   int64 `json:"synced"`
	Failed   int64 `json:"failed"`
	Conflict int64 `json:"conflict"`
	Syncing  int64 `json:"syncing"`
	Total    int64 `json:"total"`
}

type SyncEngine struct {
	db       *gorm.DB
	callback SyncCallback
}

func NewSyncEngine(db *gorm.DB) *SyncEngine {
	return &SyncEngine{db: db}
}

func (e *SyncEngine) SetSyncCallback(cb SyncCallback) {
	e.callback = cb
}

func (e *SyncEngine) QueueEvent(event *models.OfflineEvent) error {
	if event.EventType == "" {
		return fmt.Errorf("event_type is required")
	}
	if event.Payload == "" {
		return fmt.Errorf("payload is required")
	}
	if event.TerminalID == "" {
		return fmt.Errorf("terminal_id is required")
	}

	if event.ID == "" {
		event.ID = uuid.New().String()
	}
	event.Status = "pending"
	now := time.Now().UTC()
	event.CreatedAt = now
	event.UpdatedAt = now
	event.RetryCount = 0
	if event.MaxRetries <= 0 {
		event.MaxRetries = 3
	}

	vc := make(VectorClock)
	vc.Increment(event.TerminalID)

	if event.VectorClock != "" && event.VectorClock != "{}" {
		existingVC := DeserializeVectorClock(event.VectorClock)
		existingVC.Merge(vc)
		event.VectorClock = existingVC.Serialize()
	} else {
		event.VectorClock = vc.Serialize()
	}

	return e.db.Create(event).Error
}

func (e *SyncEngine) Sync() (*SyncResult, error) {
	e.db.Model(&models.OfflineEvent{}).
		Where("status = ?", "syncing").
		Update("status", "pending")

	var events []models.OfflineEvent
	if err := e.db.Where("status IN ?", []string{"pending", "failed"}).
		Where("retry_count < max_retries OR retry_count = 0").
		Order("priority DESC, created_at ASC").
		Find(&events).Error; err != nil {
		return nil, fmt.Errorf("failed to query events: %w", err)
	}

	result := &SyncResult{Total: len(events)}
	aggregateClocks := make(map[string]VectorClock)

	for i := range events {
		evt := &events[i]

		if evt.RetryCount >= evt.MaxRetries {
			result.Skipped++
			continue
		}

		e.db.Model(evt).Update("status", "syncing")

		aggID := evt.AggregateID
		eventVC := DeserializeVectorClock(evt.VectorClock)

		conflict := ""
		shouldApply := true

		if aggID != "" {
			currentClock, exists := aggregateClocks[aggID]
			if exists {
				switch eventVC.Compare(currentClock) {
				case AFTER:
					aggregateClocks[aggID].Merge(eventVC)
				case EQUAL:
					aggregateClocks[aggID].Merge(eventVC)
				case BEFORE:
					shouldApply = false
					conflict = "vector_clock_outdated"
					evt.Status = "conflict"
					evt.ErrorMessage = "Event vector clock is behind current aggregate state"
				case CONCURRENT:
					if evt.CreatedAt.After(time.Now().UTC().Add(-5 * time.Minute)) {
						aggregateClocks[aggID].Merge(eventVC)
					} else {
						shouldApply = false
						conflict = "concurrent_modification"
						evt.Status = "conflict"
						evt.ErrorMessage = "Concurrent modification detected, manual review required"
					}
				}
			} else {
				aggregateClocks[aggID] = eventVC
			}
		}

		if evt.Status != "conflict" && shouldApply {
			if e.callback != nil {
				if err := e.callback(evt); err != nil {
					evt.Status = "failed"
					evt.ErrorMessage = err.Error()
					evt.RetryCount++
					e.db.Model(evt).Updates(map[string]interface{}{
						"status":        "failed",
						"error_message": err.Error(),
						"retry_count":   evt.RetryCount,
						"updated_at":    time.Now().UTC(),
					})
					result.Failed++
				} else {
					now := time.Now().UTC()
					evt.Status = "synced"
					evt.SyncedAt = &now
					evt.ErrorMessage = ""
					e.db.Model(evt).Updates(map[string]interface{}{
						"status":        "synced",
						"synced_at":     now,
						"error_message": "",
						"updated_at":    now,
					})
					result.Synced++
				}
			} else {
				now := time.Now().UTC()
				evt.Status = "synced"
				evt.SyncedAt = &now
				e.db.Model(evt).Updates(map[string]interface{}{
					"status":     "synced",
					"synced_at":  now,
					"updated_at": now,
				})
				result.Synced++
			}
		} else if evt.Status == "conflict" {
			e.db.Model(evt).Updates(map[string]interface{}{
				"status":        "conflict",
				"error_message": evt.ErrorMessage,
				"updated_at":    time.Now().UTC(),
			})
			result.Conflicts++
		}

		result.Details = append(result.Details, SyncEventResult{
			EventID:  evt.ID,
			Status:   evt.Status,
			Conflict: conflict,
			Error:    evt.ErrorMessage,
		})
	}

	return result, nil
}

func (e *SyncEngine) GetQueueStatus() (*QueueStatus, error) {
	status := &QueueStatus{}

	if err := e.db.Model(&models.OfflineEvent{}).
		Select("COUNT(*)").Where("status = 'pending'").Scan(&status.Pending).Error; err != nil {
		return nil, err
	}
	if err := e.db.Model(&models.OfflineEvent{}).
		Select("COUNT(*)").Where("status = 'synced'").Scan(&status.Synced).Error; err != nil {
		return nil, err
	}
	if err := e.db.Model(&models.OfflineEvent{}).
		Select("COUNT(*)").Where("status = 'failed'").Scan(&status.Failed).Error; err != nil {
		return nil, err
	}
	if err := e.db.Model(&models.OfflineEvent{}).
		Select("COUNT(*)").Where("status = 'conflict'").Scan(&status.Conflict).Error; err != nil {
		return nil, err
	}
	if err := e.db.Model(&models.OfflineEvent{}).
		Select("COUNT(*)").Where("status = 'syncing'").Scan(&status.Syncing).Error; err != nil {
		return nil, err
	}

	status.Total = status.Pending + status.Synced + status.Failed + status.Conflict + status.Syncing
	return status, nil
}

func (e *SyncEngine) RemoveEvent(id string) error {
	result := e.db.Delete(&models.OfflineEvent{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("event not found: %s", id)
	}
	return nil
}

func (e *SyncEngine) GetPendingCount() int {
	var count int64
	e.db.Model(&models.OfflineEvent{}).Where("status = ?", "pending").Count(&count)
	return int(count)
}

func (e *SyncEngine) ListEvents(status string, page, pageSize int) ([]models.OfflineEvent, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	query := e.db.Model(&models.OfflineEvent{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var events []models.OfflineEvent
	if err := query.Order("priority DESC, created_at ASC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

func (e *SyncEngine) ResetStuckEvents() error {
	return e.db.Model(&models.OfflineEvent{}).
		Where("status = ?", "syncing").
		Update("status", "pending").Error
}

func (e *SyncEngine) LastSynced(event *models.OfflineEvent) error {
	return e.db.Where("status = ?", "synced").
		Order("synced_at DESC").
		First(event).Error
}
