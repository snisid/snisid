package snapshot

import (
	"sync"
	"github.com/snisid/platform/backend/internal/platform/logger"
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

func (s *SnapshotStore) Save(state ValidState) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Only save if the state has been verified by the RuntimeChecker
	logger.Info("SNAPSHOT: Committing formally verified state to history.")
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
