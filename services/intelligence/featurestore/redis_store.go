package featurestore

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type Store struct {
	Client *redis.Client
}

func NewStore(addr string) *Store {
	return &Store{
		Client: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
	}
}

func (s *Store) SaveFeature(ctx context.Context, key string, value string) error {
	return s.Client.Set(ctx, key, value, 0).Err()
}

func (s *Store) GetFeature(ctx context.Context, key string) (string, error) {
	return s.Client.Get(ctx, key).Result()
}
