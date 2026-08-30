package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer interface {
	Publish(ctx context.Context, key string, msg interface{}) error
}

type producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) Producer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.Murmur2Balancer{},
		BatchTimeout: 10 * time.Millisecond,
	}
	return &producer{writer: w}
}

func (p *producer) Publish(ctx context.Context, key string, msg interface{}) error {
	val, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: val,
	})
}

type Consumer interface {
	Start(ctx context.Context, handler func(ctx context.Context, msg []byte) error) error
}

type consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, groupID string, topic string) Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10,
		MaxBytes: 10e6,
	})
	return &consumer{reader: r}
}

func (c *consumer) Start(ctx context.Context, handler func(ctx context.Context, msg []byte) error) error {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			return err
		}
		if err := handler(ctx, m.Value); err != nil {
			return err
		}
	}
}
