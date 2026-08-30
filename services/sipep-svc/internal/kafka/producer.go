package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			BatchTimeout: 10 * time.Millisecond,
			BatchSize:    1,
			RequiredAcks: kafka.RequireOne,
		},
	}
}

func (p *Producer) Publish(ctx context.Context, topic string, key string, msg interface{}) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}
	return p.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: data,
	})
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
