package orchestrator

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockProducer struct {
	published []struct {
		Key   string
		Value interface{}
	}
}

func (m *mockProducer) Publish(ctx context.Context, key string, event interface{}) error {
	m.published = append(m.published, struct {
		Key   string
		Value interface{}
	}{Key: key, Value: event})
	return nil
}

func (m *mockProducer) Close() error { return nil }

func TestNewWorkflowManager(t *testing.T) {
	p := &mockProducer{}
	w := NewWorkflowManager(p)
	assert.NotNil(t, w)
	_, ok := w.stateStore.Load("nonexistent")
	assert.False(t, ok)
}

func TestHandleIdentityCreated_Success(t *testing.T) {
	p := &mockProducer{}
	w := NewWorkflowManager(p)

	evt := map[string]string{"identityId": "id-123"}
	data, _ := json.Marshal(evt)
	msg := kafka.Message{Value: data}

	err := w.HandleIdentityCreated(msg)
	require.NoError(t, err)

	val, ok := w.stateStore.Load("id-123")
	require.True(t, ok)
	state := val.(*IdentityState)
	assert.True(t, state.Created)
	assert.False(t, state.FraudScored)
	assert.False(t, state.Completed)
}

func TestHandleIdentityCreated_InvalidJSON(t *testing.T) {
	p := &mockProducer{}
	w := NewWorkflowManager(p)
	msg := kafka.Message{Value: []byte(`{bad json}`)}
	err := w.HandleIdentityCreated(msg)
	assert.Error(t, err)
}

func TestHandleFraudScored_Success(t *testing.T) {
	p := &mockProducer{}
	w := NewWorkflowManager(p)

	evt := map[string]interface{}{"identityId": "id-456", "riskScore": float64(30), "isFraud": false}
	data, _ := json.Marshal(evt)
	msg := kafka.Message{Value: data}

	err := w.HandleFraudScored(msg)
	require.NoError(t, err)

	val, ok := w.stateStore.Load("id-456")
	require.True(t, ok)
	state := val.(*IdentityState)
	assert.True(t, state.FraudScored)
	assert.Equal(t, 30, state.FraudRisk)
	assert.False(t, state.IsFraud)
}

func TestHandleFraudScored_InvalidJSON(t *testing.T) {
	p := &mockProducer{}
	w := NewWorkflowManager(p)
	msg := kafka.Message{Value: []byte(`{bad}`)}
	err := w.HandleFraudScored(msg)
	assert.Error(t, err)
}

func TestEvaluateState_NotReady(t *testing.T) {
	p := &mockProducer{}
	w := NewWorkflowManager(p)
	state := &IdentityState{IdentityID: "id-1", Created: true, FraudScored: false}
	err := w.evaluateState(context.Background(), state)
	require.NoError(t, err)
	assert.False(t, state.Completed)
	assert.Empty(t, p.published)
}

func TestEvaluateState_Approved(t *testing.T) {
	p := &mockProducer{}
	w := NewWorkflowManager(p)
	state := &IdentityState{IdentityID: "id-1", Created: true, FraudScored: true, FraudRisk: 30, IsFraud: false}
	err := w.evaluateState(context.Background(), state)
	require.NoError(t, err)
	assert.True(t, state.Completed)
	require.Len(t, p.published, 1)
	result := p.published[0].Value.(map[string]interface{})
	assert.Equal(t, "approved", result["status"])
}

func TestEvaluateState_RejectedByFraud(t *testing.T) {
	p := &mockProducer{}
	w := NewWorkflowManager(p)
	state := &IdentityState{IdentityID: "id-2", Created: true, FraudScored: true, FraudRisk: 90, IsFraud: true}
	err := w.evaluateState(context.Background(), state)
	require.NoError(t, err)
	assert.True(t, state.Completed)
	require.Len(t, p.published, 1)
	result := p.published[0].Value.(map[string]interface{})
	assert.Equal(t, "rejected", result["status"])
}

func TestEvaluateState_RejectedByHighRisk(t *testing.T) {
	p := &mockProducer{}
	w := NewWorkflowManager(p)
	state := &IdentityState{IdentityID: "id-3", Created: true, FraudScored: true, FraudRisk: 95, IsFraud: false}
	err := w.evaluateState(context.Background(), state)
	require.NoError(t, err)
	assert.True(t, state.Completed)
	require.Len(t, p.published, 1)
	result := p.published[0].Value.(map[string]interface{})
	assert.Equal(t, "rejected", result["status"])
}

func TestEvaluateState_AlreadyCompleted(t *testing.T) {
	p := &mockProducer{}
	w := NewWorkflowManager(p)
	state := &IdentityState{IdentityID: "id-4", Created: true, FraudScored: true, Completed: true}
	err := w.evaluateState(context.Background(), state)
	require.NoError(t, err)
	assert.Empty(t, p.published)
}

func TestGetState_Existing(t *testing.T) {
	w := NewWorkflowManager(&mockProducer{})
	w.stateStore.Store("id-1", &IdentityState{IdentityID: "id-1", Created: true})
	state := w.getState("id-1")
	assert.True(t, state.Created)
}

func TestGetState_New(t *testing.T) {
	w := NewWorkflowManager(&mockProducer{})
	state := w.getState("new-id")
	assert.Equal(t, "new-id", state.IdentityID)
	assert.False(t, state.Created)
	assert.False(t, state.FraudScored)
}

func TestFullWorkflowLifecycle(t *testing.T) {
	p := &mockProducer{}
	w := NewWorkflowManager(p)

	evt1 := map[string]string{"identityId": "wf-1"}
	d1, _ := json.Marshal(evt1)
	w.HandleIdentityCreated(kafka.Message{Value: d1})

	evt2 := map[string]interface{}{"identityId": "wf-1", "riskScore": float64(20), "isFraud": false}
	d2, _ := json.Marshal(evt2)
	w.HandleFraudScored(kafka.Message{Value: d2})

	state := w.getState("wf-1")
	assert.True(t, state.Completed)
	require.Len(t, p.published, 1)
}

func TestWorkflowManager_ConcurrentAccess(t *testing.T) {
	p := &mockProducer{}
	w := NewWorkflowManager(p)

	done := make(chan bool, 2)
	go func() {
		evt := map[string]string{"identityId": "con-1"}
		d, _ := json.Marshal(evt)
		w.HandleIdentityCreated(kafka.Message{Value: d})
		done <- true
	}()
	go func() {
		evt := map[string]interface{}{"identityId": "con-1", "riskScore": float64(40), "isFraud": false}
		d, _ := json.Marshal(evt)
		w.HandleFraudScored(kafka.Message{Value: d})
		done <- true
	}()
	<-done
	<-done

	state := w.getState("con-1")
	assert.True(t, state.Completed)
}
