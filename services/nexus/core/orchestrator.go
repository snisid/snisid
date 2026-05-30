package core

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	nexusv1 "github.com/snisid/platform/api/proto/nexus/v1"
	"github.com/snisid/platform/pkg/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Orchestrator struct {
	nexusv1.UnimplementedNexusServiceServer

	mu      sync.RWMutex
	logger  *zap.Logger
	tracer  trace.Tracer
	state   StateManager
	agents  map[string]Agent
	queue   chan *Task
	workers int
}

func NewOrchestrator(workers int, logger *zap.Logger, state StateManager) *Orchestrator {
	return &Orchestrator{
		logger:  logger,
		tracer:  telemetry.GetTracer(),
		state:   state,
		agents:  make(map[string]Agent),
		queue:   make(chan *Task, 10000),
		workers: workers,
	}
}

func (o *Orchestrator) RegisterAgent(agent Agent) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.agents[agent.ID()] = agent
	o.logger.Info("agent registered", zap.String("id", agent.ID()), zap.String("type", agent.Type()))
}

func (o *Orchestrator) SubmitTask(ctx context.Context, def *nexusv1.TaskDefinition) (*nexusv1.TaskResponse, error) {
	ctx, span := o.tracer.Start(ctx, "nexus.SubmitTask", trace.WithAttributes(
		attribute.String("task.type", def.Type),
	))
	defer span.End()

	if def.Id == "" {
		def.Id = uuid.NewString()
	}
	span.SetAttributes(attribute.String("task.id", def.Id))

	task := &Task{
		Definition: def,
		Status:     nexusv1.TaskStatus_TASK_STATUS_PENDING,
		UpdatedAt:  time.Now(),
	}

	if err := o.state.SaveTask(ctx, task); err != nil {
		return nil, fmt.Errorf("failed to save task: %w", err)
	}

	o.queue <- task

	return &nexusv1.TaskResponse{
		TaskId:   def.Id,
		Accepted: true,
		Message:  "Task submitted and queued",
	}, nil
}

func (o *Orchestrator) Start(ctx context.Context) {
	o.logger.Info("starting nexus orchestrator", zap.Int("workers", o.workers))

	for i := 0; i < o.workers; i++ {
		go o.worker(ctx, i)
	}
}

func (o *Orchestrator) worker(ctx context.Context, id int) {
	o.logger.Debug("worker started", zap.Int("id", id))

	for {
		select {
		case <-ctx.Done():
			o.logger.Debug("worker stopping", zap.Int("id", id))
			return

		case task := <-o.queue:
			o.processTask(ctx, task)
		}
	}
}

func (o *Orchestrator) processTask(ctx context.Context, task *Task) {
	o.logger.Info("processing task", zap.String("id", task.Definition.Id), zap.String("type", task.Definition.Type))

	o.mu.RLock()
	var selectedAgent Agent
	for _, agent := range o.agents {
		if agent.CanHandle(task) {
			selectedAgent = agent
			break
		}
	}
	o.mu.RUnlock()

	if selectedAgent == nil {
		o.logger.Warn("no agent available for task", zap.String("id", task.Definition.Id))
		// Optional: Re-queue or fail after timeout
		return
	}

	o.state.UpdateTaskStatus(ctx, task.Definition.Id, nexusv1.TaskStatus_TASK_STATUS_RUNNING)

	signal, err := selectedAgent.Execute(ctx, task)
	if err != nil {
		o.logger.Error("task execution failed", zap.String("id", task.Definition.Id), zap.Error(err))
		o.state.UpdateTaskStatus(ctx, task.Definition.Id, nexusv1.TaskStatus_TASK_STATUS_FAILED)
		return
	}

	o.logger.Info("task completed successfully", zap.String("id", task.Definition.Id), zap.String("status", signal.Status.String()))
	o.state.UpdateTaskStatus(ctx, task.Definition.Id, signal.Status)
}
