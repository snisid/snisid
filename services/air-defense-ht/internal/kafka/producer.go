package kafka

import (
	"context"
	"log"
	"time"

	kafkago "github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafkago.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	w := &kafkago.Writer{
		Addr:         kafkago.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafkago.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
	}
	return &Producer{writer: w}
}

func (p *Producer) Publish(key string, value string) {
	msg := kafkago.Message{
		Key:   []byte(key),
		Value: []byte(value),
		Time:  time.Now(),
	}
	if err := p.writer.WriteMessages(context.Background(), msg); err != nil {
		log.Printf("kafka publish error: %v", err)
	}
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
