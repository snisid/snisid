package commandcenter

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBroadcast_DeliversToSubscriber(t *testing.T) {
	subMu.Lock()
	subscribers = make(map[chan LiveEvent]bool)
	subMu.Unlock()

	ch := make(chan LiveEvent, 10)
	subMu.Lock()
	subscribers[ch] = true
	subMu.Unlock()

	event := LiveEvent{
		Type:     "TEST",
		Severity: "INFO",
		Source:   "test",
		Payload:  "hello",
	}
	Broadcast(event)

	select {
	case received := <-ch:
		assert.Equal(t, event.Type, received.Type)
		assert.Equal(t, event.Severity, received.Severity)
		assert.Equal(t, event.Source, received.Source)
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for event")
	}
}

func TestBroadcast_MultipleSubscribers(t *testing.T) {
	subMu.Lock()
	subscribers = make(map[chan LiveEvent]bool)
	subMu.Unlock()

	ch1 := make(chan LiveEvent, 10)
	ch2 := make(chan LiveEvent, 10)
	subMu.Lock()
	subscribers[ch1] = true
	subscribers[ch2] = true
	subMu.Unlock()

	Broadcast(LiveEvent{Type: "ALERT", Severity: "HIGH"})

	<-ch1
	<-ch2
}

func TestBroadcast_SkipsFullBuffer(t *testing.T) {
	subMu.Lock()
	subscribers = make(map[chan LiveEvent]bool)
	subMu.Unlock()

	ch := make(chan LiveEvent, 1)
	ch <- LiveEvent{Type: "full"}
	subMu.Lock()
	subscribers[ch] = true
	subMu.Unlock()

	Broadcast(LiveEvent{Type: "dropped"})
	assert.Len(t, ch, 1)
}

func TestBroadcast_NoSubscribers(t *testing.T) {
	subMu.Lock()
	subscribers = make(map[chan LiveEvent]bool)
	subMu.Unlock()

	Broadcast(LiveEvent{Type: "test"})
}

func TestTriggerAlert(t *testing.T) {
	subMu.Lock()
	subscribers = make(map[chan LiveEvent]bool)
	subMu.Unlock()

	ch := make(chan LiveEvent, 10)
	subMu.Lock()
	subscribers[ch] = true
	subMu.Unlock()

	TriggerAlert("FRAUD", "CIT-001")

	select {
	case evt := <-ch:
		assert.Equal(t, "CRITICAL_ALERT", evt.Type)
		assert.Equal(t, "HIGH", evt.Severity)
		assert.Equal(t, "FRAUD_ENGINE", evt.Source)
		payload, ok := evt.Payload.(map[string]string)
		require.True(t, ok)
		assert.Equal(t, "FRAUD", payload["type"])
		assert.Equal(t, "CIT-001", payload["citizen_id"])
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for alert")
	}
}

func TestLiveEvent_JSONSerialization(t *testing.T) {
	event := LiveEvent{
		Type:     "UPDATE",
		Severity: "LOW",
		Source:   "system",
		Payload:  map[string]int{"count": 42},
	}
	data, err := json.Marshal(event)
	require.NoError(t, err)

	var decoded LiveEvent
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	assert.Equal(t, event.Type, decoded.Type)
	assert.Equal(t, event.Severity, decoded.Severity)
	assert.Equal(t, event.Source, decoded.Source)
}

func TestConcurrentBroadcast(t *testing.T) {
	subMu.Lock()
	subscribers = make(map[chan LiveEvent]bool)
	subMu.Unlock()

	ch := make(chan LiveEvent, 100)
	subMu.Lock()
	subscribers[ch] = true
	subMu.Unlock()

	done := make(chan bool)
	go func() {
		for i := 0; i < 50; i++ {
			Broadcast(LiveEvent{Type: "concurrent", Severity: "INFO"})
		}
		done <- true
	}()
	go func() {
		for i := 0; i < 50; i++ {
			Broadcast(LiveEvent{Type: "concurrent", Severity: "INFO"})
		}
		done <- true
	}()

	<-done
	<-done

	count := 0
	for len(ch) > 0 {
		<-ch
		count++
	}
	assert.Equal(t, 100, count)
}
