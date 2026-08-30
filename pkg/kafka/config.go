package kafka

import (
	"time"

	"github.com/segmentio/kafka-go"
)

type Config struct {
	Brokers       []string
	Topic         string
	GroupID       string
	MaxRetries    int
	RetryInterval time.Duration
	EnableTLS     bool
}

func (c *Config) WriterConfig() kafka.WriterConfig {
	return kafka.WriterConfig{
		Brokers:  c.Brokers,
		Topic:    c.Topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func (c *Config) ReaderConfig() kafka.ReaderConfig {
	return kafka.ReaderConfig{
		Brokers:  c.Brokers,
		Topic:    c.Topic,
		GroupID:  c.GroupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	}
}
