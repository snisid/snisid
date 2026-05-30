package events

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisStreamManager struct {
	client *redis.Client
}

func NewRedisStreamManager(client *redis.Client) *RedisStreamManager {
	return &RedisStreamManager{client: client}
}

func (m *RedisStreamManager) Publish(ctx context.Context, stream string, data map[string]interface{}) error {
	return m.client.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		Values: data,
	}).Err()
}

func (m *RedisStreamManager) Subscribe(ctx context.Context, stream string, group string, handler func(data map[string]interface{}) error) error {
	// Simple Redis Stream consumer logic
	for {
		entries, err := m.client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    group,
			Consumer: "platform-consumer",
			Streams:  []string{stream, ">"},
			Count:    10,
			Block:    0,
		}).Result()
		if err != nil {
			return err
		}

		for _, entry := range entries {
			for _, msg := range entry.Messages {
				if err := handler(msg.Values); err != nil {
					return err
				}
				m.client.XAck(ctx, stream, group, msg.ID)
			}
		}
	}
}
