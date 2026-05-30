package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Producer struct {
	writer *kafka.Writer
	tracer trace.Tracer
}

func NewProducer(config Config) *Producer {
	return &Producer{
		writer: kafka.NewWriter(config.WriterConfig()),
		tracer: otel.Tracer("kafka-producer"),
	}
}

func (p *Producer) Publish(ctx context.Context, key string, value []byte) error {
	ctx, span := p.tracer.Start(ctx, "kafka.publish", trace.WithAttributes(
		attribute.String("messaging.system", "kafka"),
		attribute.String("messaging.destination", p.writer.Topic),
	))
	defer span.End()

	msg := kafka.Message{
		Key:   []byte(key),
		Value: value,
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
