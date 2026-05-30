package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"github.com/snisid/platform/backend/internal/platform/logger"
)

type Producer struct {
	Writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{
		Writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *Producer) Publish(ctx context.Context, key string, data interface{}) error {
	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}

	logger.Info("KAFKA: Publishing message to topic " + p.Writer.Topic)
	
	return p.Writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: msg,
	})
}
