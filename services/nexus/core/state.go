package core

import (
	"context"
	"fmt"
	"sync"
	"time"

	nexusv1 "github.com/snisid/platform/api/proto/nexus/v1"
)

type InMemoryState struct {
	mu    sync.RWMutex
	tasks map[string]*Task
}

func NewInMemoryState() *InMemoryState {
	return &InMemoryState{
		tasks: make(map[string]*Task),
	}
}

func (s *InMemoryState) SaveTask(ctx context.Context, task *Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tasks[task.Definition.Id] = task
	return nil
}

func (s *InMemoryState) GetTask(ctx context.Context, id string) (*Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, ok := s.tasks[id]
	if !ok {
		return nil, fmt.Errorf("task not found: %s", id)
	}
	return task, nil
}

func (s *InMemoryState) UpdateTaskStatus(ctx context.Context, id string, status nexusv1.TaskStatus) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.tasks[id]
	if !ok {
		return fmt.Errorf("task not found: %s", id)
	}

	task.Status = status
	task.UpdatedAt = time.Now()
	return nil
}
