package database

import (
	"context"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/redis/go-redis/v9"
	"github.com/snisid/platform/backend/internal/platform/logger"
)

func NewRedisClient(addr string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.Error("failed to connect to redis", err)
	}
	return rdb
}

func NewNeo4jDriver(uri, user, pass string) (neo4j.DriverWithContext, error) {
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(user, pass, ""))
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := driver.VerifyConnectivity(ctx); err != nil {
		return nil, err
	}
	return driver, nil
}

// Postgres connection logic (using GORM would be better for production)
// For now, I'll stick to basic pgx or sql if needed, but the prompt mentioned GORM in my thoughts.
// Actually, let's keep it simple with standard library for now or just placeholders for the services to implement their own.
// I'll provide a generic helper.
