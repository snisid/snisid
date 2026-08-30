package snapshot

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveAndGetLastValid(t *testing.T) {
	s := &SnapshotStore{}
	state := ValidState{Timestamp: 100, RiskData: map[string]int{"node-1": 50}}
	s.Save(state)

	last := s.GetLastValid()
	assert.Equal(t, int64(100), last.Timestamp)
	assert.Equal(t, 50, last.RiskData["node-1"])
}

func TestGetLastValid_Empty(t *testing.T) {
	s := &SnapshotStore{}
	state := s.GetLastValid()
	assert.Equal(t, int64(0), state.Timestamp)
	assert.Nil(t, state.RiskData)
	assert.Nil(t, state.PolicySet)
}

func TestSave_MultipleStates(t *testing.T) {
	s := &SnapshotStore{}
	for i := 0; i < 5; i++ {
		s.Save(ValidState{Timestamp: int64(i * 100)})
	}

	require.Len(t, s.History, 5)
	last := s.GetLastValid()
	assert.Equal(t, int64(400), last.Timestamp)
}

func TestSave_OrderPreserved(t *testing.T) {
	s := &SnapshotStore{}
	s.Save(ValidState{Timestamp: 300})
	s.Save(ValidState{Timestamp: 100})
	s.Save(ValidState{Timestamp: 200})

	require.Len(t, s.History, 3)
	assert.Equal(t, int64(300), s.History[0].Timestamp)
	assert.Equal(t, int64(100), s.History[1].Timestamp)
	assert.Equal(t, int64(200), s.History[2].Timestamp)
	assert.Equal(t, int64(200), s.GetLastValid().Timestamp)
}

func TestSave_WithPolicySet(t *testing.T) {
	s := &SnapshotStore{}
	state := ValidState{
		Timestamp: 1,
		RiskData:  map[string]int{"n1": 10, "n2": 20},
		PolicySet: map[string]string{"n1": "ALLOW", "n2": "DENY"},
	}
	s.Save(state)

	last := s.GetLastValid()
	assert.Equal(t, 10, last.RiskData["n1"])
	assert.Equal(t, "ALLOW", last.PolicySet["n1"])
}

func TestConcurrentAccess(t *testing.T) {
	s := &SnapshotStore{}
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			s.Save(ValidState{Timestamp: int64(i)})
		}(i)

		wg.Add(1)
		go func() {
			defer wg.Done()
			s.GetLastValid()
		}()
	}
	wg.Wait()

	assert.Len(t, s.History, 50)
}
