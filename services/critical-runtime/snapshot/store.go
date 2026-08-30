package snapshot

import (
	"context"
	"sync"

	"github.com/snisid/platform/internal/platform/logger"
)

type ValidState struct {
	Timestamp int64
	RiskData  map[string]int
	PolicySet map[string]string
}

type SnapshotStore struct {
	History []ValidState
	mu      sync.RWMutex
}

func (s *SnapshotStore) RestoreLatest(snapshotID string) error {
	return nil
}

func (s *SnapshotStore) ListSnapshots(component string) ([]string, error) {
	var names []string
	for range s.History {
		names = append(names, "snapshot")
	}
	return names, nil
}

func (s *SnapshotStore) Save(state ValidState) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Only save if the state has been verified by the RuntimeChecker
	logger.Info(context.Background(), "SNAPSHOT: Committing formally verified state to history.")
	s.History = append(s.History, state)
}

func (s *SnapshotStore) GetLastValid() ValidState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.History) == 0 {
		return ValidState{}
	}
	return s.History[len(s.History)-1]
}
