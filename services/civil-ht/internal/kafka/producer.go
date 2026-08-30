package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type Event struct {
	EventType string    `json:"event_type"`
	ActID     string    `json:"act_id,omitempty"`
	ActNumber string    `json:"act_number,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Data      any       `json:"data"`
}

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
		Compression:  kafka.Snappy,
		Async:        true,
	}
	return &Producer{writer: w}
}

func (p *Producer) Publish(ctx context.Context, event Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}
	msg := kafka.Message{
		Key:   []byte(event.ActID),
		Value: payload,
	}
	return p.writer.WriteMessages(ctx, msg)
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
