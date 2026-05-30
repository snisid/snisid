package fraud

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type StateStore struct {
	client *redis.Client
}

func NewStateStore(addr string) *StateStore {
	return &StateStore{
		client: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
	}
}

// IncrementVelocity tracks the frequency of an event within a window
func (s *StateStore) IncrementVelocity(ctx context.Context, key string, window time.Duration) (int64, error) {
	pipe := s.client.Pipeline()
	count := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, window)
	
	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}

	return count.Val(), nil
}

// GetState retrieves a generic state value
func (s *StateStore) GetState(ctx context.Context, key string) (string, error) {
	return s.client.Get(ctx, key).Result()
}

func (s *StateStore) Close() error {
	return s.client.Close()
}
