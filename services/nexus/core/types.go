package core

import (
	"context"
	"time"

	nexusv1 "github.com/snisid/platform/api/proto/nexus/v1"
)

// Task wraps the protobuf TaskDefinition with internal metadata.
type Task struct {
	Definition *nexusv1.TaskDefinition
	Status     nexusv1.TaskStatus
	UpdatedAt  time.Time
	Retries    int
}

// Agent represents an execution unit (Kai).
type Agent interface {
	ID() string
	Type() string
	CanHandle(task *Task) bool
	Execute(ctx context.Context, task *Task) (*nexusv1.AgentSignal, error)
}

// StateManager handles persistence and distributed state.
type StateManager interface {
	SaveTask(ctx context.Context, task *Task) error
	GetTask(ctx context.Context, id string) (*Task, error)
	UpdateTaskStatus(ctx context.Context, id string, status nexusv1.TaskStatus) error
}
