package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      brokers,
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		WriteTimeout: 10 * time.Second,
	})
	return &Producer{writer: w}
}

func (p *Producer) PublishWarrantCreated(warrant interface{}) {
	data, _ := json.Marshal(warrant)
	p.writer.WriteMessages(context.Background(), kafka.Message{Value: data})
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
