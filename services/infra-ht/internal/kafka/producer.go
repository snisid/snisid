package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

type Event struct {
	EventType string    `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
	Data      any       `json:"data"`
}

type Producer struct{ writer *kafka.Writer }

func NewProducer(b []string, t string) *Producer {
	return &Producer{writer: &kafka.Writer{Addr: kafka.TCP(b...), Topic: t, Balancer: &kafka.LeastBytes{}, RequiredAcks: kafka.RequireOne, Async: true}}
}
func (p *Producer) Publish(ctx context.Context, e Event) error {
	b, _ := json.Marshal(e)
	return p.writer.WriteMessages(ctx, kafka.Message{Key: []byte(e.EventType), Value: b})
}
func (p *Producer) Close() error { return p.writer.Close() }
