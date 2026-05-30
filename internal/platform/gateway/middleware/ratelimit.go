package middleware

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	client *redis.Client
	limit  int
	window time.Duration
}

func NewRateLimiter(client *redis.Client, limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		client: client,
		limit:  limit,
		window: window,
	}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr // In production, parse X-Forwarded-For

		ctx := r.Context()
		key := "ratelimit:" + ip

		// Using Redis INCR and EXPIRE for a simple fixed-window rate limiter
		count, err := rl.client.Incr(ctx, key).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				count = 1
			} else {
				// Fail open on redis error
				next.ServeHTTP(w, r)
				return
			}
		}

		if count == 1 {
			rl.client.Expire(ctx, key, rl.window)
		}

		if count > int64(rl.limit) {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
