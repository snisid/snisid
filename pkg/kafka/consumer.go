package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Consumer struct {
	reader *kafka.Reader
	writer *kafka.Writer // For DLQ
	logger *zap.Logger
	tracer trace.Tracer
}

func NewConsumer(config Config, logger *zap.Logger) *Consumer {
	var dlqWriter *kafka.Writer
	if config.Topic != "" {
		dlqWriter = kafka.NewWriter(kafka.WriterConfig{
			Brokers: config.Brokers,
			Topic:   config.Topic + ".dlq",
		})
	}

	return &Consumer{
		reader: kafka.NewReader(config.ReaderConfig()),
		writer: dlqWriter,
		logger: logger,
		tracer: otel.Tracer("kafka-consumer"),
	}
}

func (c *Consumer) Start(ctx context.Context, handler func(context.Context, kafka.Message) error) {
	c.logger.Info("starting kafka consumer", zap.String("topic", c.reader.Config().Topic))

	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			c.logger.Error("failed to fetch message", zap.Error(err))
			continue
		}

		c.processMessage(ctx, msg, handler)
	}
}

func (c *Consumer) processMessage(ctx context.Context, msg kafka.Message, handler func(context.Context, kafka.Message) error) {
	ctx, span := c.tracer.Start(ctx, "kafka.consume", trace.WithAttributes(
		attribute.String("messaging.system", "kafka"),
		attribute.String("messaging.destination", msg.Topic),
		attribute.String("messaging.kafka.partition", fmt.Sprintf("%d", msg.Partition)),
	))
	defer span.End()

	if err := handler(ctx, msg); err != nil {
		c.logger.Error("failed to handle message", zap.Error(err), zap.String("key", string(msg.Key)))
		span.RecordError(err)

		// Send to DLQ if configured
		if c.writer != nil {
			dlqErr := c.writer.WriteMessages(ctx, kafka.Message{
				Key:   msg.Key,
				Value: msg.Value,
				Headers: []kafka.Header{
					{Key: "error", Value: []byte(err.Error())},
				},
			})
			if dlqErr != nil {
				c.logger.Error("failed to send to DLQ", zap.Error(dlqErr))
			}
		}
	} else {
		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			c.logger.Error("failed to commit message", zap.Error(err))
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
