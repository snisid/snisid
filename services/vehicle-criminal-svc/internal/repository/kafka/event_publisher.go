package kafka

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

type EventPublisher struct {
	writer *kafka.Writer
}

func NewEventPublisher(writer *kafka.Writer) *EventPublisher {
	return &EventPublisher{writer: writer}
}

func (p *EventPublisher) Publish(ctx context.Context, topic string, event interface{}) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(topic),
		Value: data,
		Headers: []kafka.Header{
			{Key: "content-type", Value: []byte("application/json")},
		},
	})
}
