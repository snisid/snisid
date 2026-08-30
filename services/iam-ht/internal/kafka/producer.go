package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

type Event struct {
	EventType string    `json:"event_type"`
	UserID    string    `json:"user_id,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Data      any       `json:"data"`
}

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	w := &kafka.Writer{Addr: kafka.TCP(brokers...), Topic: topic, Balancer: &kafka.LeastBytes{}, RequiredAcks: kafka.RequireOne, Compression: kafka.Snappy, Async: true}
	return &Producer{writer: w}
}

func (p *Producer) Publish(ctx context.Context, event Event) error {
	payload, _ := json.Marshal(event)
	return p.writer.WriteMessages(ctx, kafka.Message{Key: []byte(event.UserID), Value: payload})
}

func (p *Producer) Close() error { return p.writer.Close() }
