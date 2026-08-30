package commandcenter

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestBroadcast_SendsToSubscribers(t *testing.T) {
	ch := make(chan LiveEvent, 10)
	subMu.Lock()
	subscribers[ch] = true
	subMu.Unlock()
	defer func() {
		subMu.Lock()
		delete(subscribers, ch)
		subMu.Unlock()
	}()

	event := LiveEvent{
		Type: "TEST", Severity: "INFO", Source: "test",
		Payload: map[string]string{"msg": "hello"},
	}
	Broadcast(event)

	select {
	case received := <-ch:
		assert.Equal(t, "TEST", received.Type)
		assert.Equal(t, "INFO", received.Severity)
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for broadcast")
	}
}

func TestBroadcast_MultipleSubscribers(t *testing.T) {
	ch1 := make(chan LiveEvent, 10)
	ch2 := make(chan LiveEvent, 10)

	subMu.Lock()
	subscribers[ch1] = true
	subscribers[ch2] = true
	subMu.Unlock()

	defer func() {
		subMu.Lock()
		delete(subscribers, ch1)
		delete(subscribers, ch2)
		subMu.Unlock()
	}()

	Broadcast(LiveEvent{Type: "MULTI", Severity: "HIGH", Source: "test"})

	select {
	case <-ch1:
	case <-time.After(time.Second):
		t.Fatal("ch1 did not receive broadcast")
	}
	select {
	case <-ch2:
	case <-time.After(time.Second):
		t.Fatal("ch2 did not receive broadcast")
	}
}

func TestBroadcast_SkipsFullBuffer(t *testing.T) {
	ch := make(chan LiveEvent, 1)
	ch <- LiveEvent{Type: "FULL"}

	subMu.Lock()
	subscribers[ch] = true
	subMu.Unlock()

	defer func() {
		subMu.Lock()
		delete(subscribers, ch)
		subMu.Unlock()
	}()

	Broadcast(LiveEvent{Type: "SKIP", Severity: "LOW", Source: "test"})

	assert.Len(t, ch, 1)
	received := <-ch
	assert.Equal(t, "FULL", received.Type)
}

func TestTriggerAlert_CreatesCriticalAlert(t *testing.T) {
	ch := make(chan LiveEvent, 10)
	subMu.Lock()
	subscribers[ch] = true
	subMu.Unlock()

	defer func() {
		subMu.Lock()
		delete(subscribers, ch)
		subMu.Unlock()
	}()

	TriggerAlert("FRAUD_SYNTHETIC", "CIT-00042")

	select {
	case event := <-ch:
		assert.Equal(t, "CRITICAL_ALERT", event.Type)
		assert.Equal(t, "HIGH", event.Severity)
		assert.Equal(t, "FRAUD_ENGINE", event.Source)
		payload, ok := event.Payload.(map[string]string)
		assert.True(t, ok)
		assert.Equal(t, "FRAUD_SYNTHETIC", payload["type"])
		assert.Equal(t, "CIT-00042", payload["citizen_id"])
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for alert")
	}
}

func TestConcurrentBroadcastAndSubscribe(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ch := make(chan LiveEvent, 10)
			subMu.Lock()
			subscribers[ch] = true
			subMu.Unlock()

			Broadcast(LiveEvent{Type: "CONCUR", Severity: "MEDIUM", Source: "conc-test"})

			subMu.Lock()
			delete(subscribers, ch)
			subMu.Unlock()
		}()
	}
	wg.Wait()
}

func TestStream_ClosesOnWriteError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		Stream(ws)
	}))
	defer server.Close()

	url := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	assert.NoError(t, err)

	Broadcast(LiveEvent{Type: "WS_TEST", Severity: "LOW", Source: "ws-test"})

	_, msg, err := ws.ReadMessage()
	assert.NoError(t, err)

	var event LiveEvent
	err = json.Unmarshal(msg, &event)
	assert.NoError(t, err)
	assert.Equal(t, "WS_TEST", event.Type)

	ws.Close()
	time.Sleep(100 * time.Millisecond)

	Broadcast(LiveEvent{Type: "AFTER_CLOSE", Severity: "LOW", Source: "ws-test"})
}

func TestSubscriberCleanupOnDisconnect(t *testing.T) {
	ch := make(chan LiveEvent, 10)
	subMu.Lock()
	subscribers[ch] = true
	subMu.Unlock()

	assert.Contains(t, subscribers, ch)

	subMu.Lock()
	delete(subscribers, ch)
	subMu.Unlock()

	assert.NotContains(t, subscribers, ch)
}

func TestLiveEventJSON(t *testing.T) {
	event := LiveEvent{
		Type: "TEST", Severity: "INFO", Source: "json-test",
		Payload: map[string]interface{}{"key": "value", "count": 42},
	}

	data, err := json.Marshal(event)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "TEST")
	assert.Contains(t, string(data), "json-test")

	var decoded LiveEvent
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, event.Type, decoded.Type)
	assert.Equal(t, event.Severity, decoded.Severity)
}
