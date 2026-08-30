package service

import (
	"fmt"
	"sync"
	"time"
)

type LocalRecord struct {
	ID        string         `json:"id"`
	Data      map[string]any `json:"data"`
	Version   int64          `json:"version"`
	UpdatedAt time.Time      `json:"updated_at"`
	Synced    bool           `json:"synced"`
}

type LAPISyncService struct {
	mu    sync.Mutex
	cache *OfflineCache
}

func NewLAPISyncService(cache *OfflineCache) *LAPISyncService {
	return &LAPISyncService{cache: cache}
}

func (s *LAPISyncService) CreateRecord(data map[string]any) (*LocalRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	record := &LocalRecord{
		ID:        fmt.Sprintf("rec_%d", time.Now().UnixNano()),
		Data:      data,
		Version:   time.Now().UnixMilli(),
		UpdatedAt: time.Now(),
		Synced:    false,
	}

	if err := s.cache.Store(record); err != nil {
		return nil, err
	}
	return record, nil
}

func (s *LAPISyncService) UpdateRecord(id string, data map[string]any) (*LocalRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, ok := s.cache.Get(id)
	if !ok {
		return nil, fmt.Errorf("record %s not found", id)
	}

	existing.Data = data
	existing.Version = time.Now().UnixMilli()
	existing.UpdatedAt = time.Now()
	existing.Synced = false

	if err := s.cache.Store(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *LAPISyncService) Sync() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, r := range s.cache.GetPending() {
		backendVersion, err := s.sendToBackend(r)
		if err != nil {
			return fmt.Errorf("sync failed for %s: %w", r.ID, err)
		}
		if backendVersion > r.Version {
			r.Version = backendVersion
		}
		r.Synced = true
		s.cache.Store(r)
	}
	return nil
}

func (s *LAPISyncService) sendToBackend(r *LocalRecord) (int64, error) {
	if r.Version == 0 {
		return 0, fmt.Errorf("invalid record version")
	}
	return time.Now().UnixMilli(), nil
}
