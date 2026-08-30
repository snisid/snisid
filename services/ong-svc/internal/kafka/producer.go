package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Producer struct {
	writer *kafka.Writer
	log    *zap.Logger
}

func NewProducer(brokers string, log *zap.Logger) *Producer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(brokers),
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
	}
	return &Producer{writer: w, log: log}
}

func (p *Producer) Publish(ctx context.Context, topic string, msg interface{}) error {
	data, err := json.Marshal(msg)
	if err != nil {
		p.log.Error("kafka marshal error", zap.Error(err))
		return err
	}
	return p.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Value: data,
	})
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
