package events

import (
	"context"
	"errors"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func newTestRedisClient(t *testing.T) (*redis.Client, *miniredis.Miniredis) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	t.Cleanup(func() {
		client.Close()
		s.Close()
	})
	return client, s
}

func TestNewRedisStreamManager(t *testing.T) {
	client, _ := newTestRedisClient(t)
	mgr := NewRedisStreamManager(client)
	if mgr == nil {
		t.Fatal("NewRedisStreamManager returned nil")
	}
}

func TestPublish_Success(t *testing.T) {
	client, _ := newTestRedisClient(t)
	mgr := NewRedisStreamManager(client)

	err := mgr.Publish(context.Background(), "test-stream", map[string]interface{}{
		"event": "identity.created",
		"id":    "ID-001",
	})
	if err != nil {
		t.Fatalf("Publish failed: %v", err)
	}
}

func TestPublish_EmptyStream(t *testing.T) {
	client, _ := newTestRedisClient(t)
	mgr := NewRedisStreamManager(client)

	// Publish without stream name
	err := mgr.Publish(context.Background(), "", map[string]interface{}{"key": "val"})
	if err != nil {
		t.Logf("Publish to empty stream returned: %v", err)
	}
}

func TestPublish_MultipleMessages(t *testing.T) {
	client, _ := newTestRedisClient(t)
	mgr := NewRedisStreamManager(client)

	for i := 0; i < 5; i++ {
		err := mgr.Publish(context.Background(), "multi-stream", map[string]interface{}{
			"index": i,
		})
		if err != nil {
			t.Fatalf("Publish message %d failed: %v", i, err)
		}
	}

	// Verify stream length
	len, err := client.XLen(context.Background(), "multi-stream").Result()
	if err != nil {
		t.Fatalf("XLen failed: %v", err)
	}
	if len != 5 {
		t.Errorf("Stream length = %d, want 5", len)
	}
}

func TestSubscribe_ContextCancelled(t *testing.T) {
	client, _ := newTestRedisClient(t)
	mgr := NewRedisStreamManager(client)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Create group first
	err := client.XGroupCreate(ctx, "test-stream", "test-group", "0").Err()
	if err != nil && !errors.Is(err, redis.Nil) {
		t.Logf("XGroupCreate result: %v", err)
	}

	err = mgr.Subscribe(ctx, "test-stream", "test-group", func(data map[string]interface{}) error {
		return nil
	})
	if err == nil {
		t.Log("Subscribe returned nil on cancelled context (expected)")
	}
}
