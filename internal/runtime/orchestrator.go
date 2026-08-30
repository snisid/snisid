package runtime

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	Pending   TaskStatus = "pending"
	Running   TaskStatus = "running"
	Completed TaskStatus = "completed"
	Failed    TaskStatus = "failed"
)

type Task struct {
	ID        string
	Type      string
	Payload   []byte
	Priority  int
	Status    TaskStatus
	CreatedAt time.Time
}

type Agent interface {
	Name() string
	CanHandle(task Task) bool
	Execute(ctx context.Context, task Task) error
}

type Orchestrator struct {
	mu      sync.RWMutex
	agents  []Agent
	queue   chan Task
	workers int
}

func NewOrchestrator(workers int) *Orchestrator {
	return &Orchestrator{
		queue:   make(chan Task, 1000),
		workers: workers,
	}
}

func (o *Orchestrator) Register(agent Agent) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.agents = append(o.agents, agent)
}

func (o *Orchestrator) Submit(taskType string, payload []byte, priority int) {
	task := Task{
		ID:        uuid.NewString(),
		Type:      taskType,
		Payload:   payload,
		Priority:  priority,
		Status:    Pending,
		CreatedAt: time.Now(),
	}

	o.queue <- task
}

func (o *Orchestrator) Start(ctx context.Context) {
	for i := 0; i < o.workers; i++ {
		go o.worker(ctx)
	}
}

func (o *Orchestrator) worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case task := <-o.queue:
			o.dispatch(ctx, task)
		}
	}
}

func (o *Orchestrator) dispatch(ctx context.Context, task Task) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	for _, agent := range o.agents {
		if agent.CanHandle(task) {
			log.Printf("dispatching task=%s to=%s", task.ID, agent.Name())

			if err := agent.Execute(ctx, task); err != nil {
				log.Printf("task failed=%v", err)
			}

			return
		}
	}

	log.Printf("no agent found for task=%s", task.ID)
}
