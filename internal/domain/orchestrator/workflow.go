package orchestrator

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/snisid/platform/internal/platform/events"
	"github.com/snisid/platform/internal/platform/logger"
)

type IdentityState struct {
	IdentityID   string
	Created      bool
	FraudScored  bool
	FraudRisk    int
	IsFraud      bool
	Completed    bool
	LastModified time.Time
}

type WorkflowManager struct {
	stateStore sync.Map // Simple in-memory state store for the initial batch
	producer   *events.Producer
}

func NewWorkflowManager(producer *events.Producer) *WorkflowManager {
	return &WorkflowManager{
		producer: producer,
	}
}

func (w *WorkflowManager) HandleIdentityCreated(msg kafka.Message) error {
	var evt struct {
		IdentityID string `json:"identityId"`
	}
	if err := json.Unmarshal(msg.Value, &evt); err != nil {
		return err
	}

	state := w.getState(evt.IdentityID)
	state.Created = true
	state.LastModified = time.Now()
	w.stateStore.Store(evt.IdentityID, state)

	return w.evaluateState(context.Background(), state)
}

func (w *WorkflowManager) HandleFraudScored(msg kafka.Message) error {
	var evt struct {
		IdentityID string `json:"identityId"`
		RiskScore  int    `json:"riskScore"`
		IsFraud    bool   `json:"isFraud"`
	}
	if err := json.Unmarshal(msg.Value, &evt); err != nil {
		return err
	}

	state := w.getState(evt.IdentityID)
	state.FraudScored = true
	state.FraudRisk = evt.RiskScore
	state.IsFraud = evt.IsFraud
	state.LastModified = time.Now()
	w.stateStore.Store(evt.IdentityID, state)

	return w.evaluateState(context.Background(), state)
}

func (w *WorkflowManager) getState(id string) *IdentityState {
	val, ok := w.stateStore.Load(id)
	if ok {
		return val.(*IdentityState)
	}
	return &IdentityState{IdentityID: id}
}

func (w *WorkflowManager) evaluateState(ctx context.Context, state *IdentityState) error {
	if state.Completed {
		return nil
	}

	if state.Created && state.FraudScored {
		// Workflow is complete
		state.Completed = true
		w.stateStore.Store(state.IdentityID, state)

		status := "approved"
		if state.IsFraud || state.FraudRisk > 80 {
			status = "rejected"
		}

		resultEvt := map[string]any{
			"identityId": state.IdentityID,
			"status":     status,
			"reason":     "Workflow completed",
			"timestamp":  time.Now().UTC(),
		}

		logger.Info(ctx, "workflow completed for identity")
		return w.producer.Publish(ctx, state.IdentityID, resultEvt)
	}

	return nil
}
