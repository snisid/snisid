package integration

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
)

func TestKafkaEventFlow(t *testing.T) {
	broker := "localhost:9092"
	topic := "test.topic"

	// 1. Produce
	writer := &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	
	msg := map[string]string{"event": "test_occured", "id": "123"}
	data, _ := json.Marshal(msg)
	
	err := writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte("123"),
		Value: data,
	})
	assert.NoError(t, err)

	// 2. Consume
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		GroupID: "test-group",
	})
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	received, err := reader.ReadMessage(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "123", string(received.Key))
}
