package indexes

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/snisid/platform/services/bio-adn/internal/db"
	"github.com/snisid/platform/services/bio-adn/pkg/models"
)

const LAPIMaxResponseMs = 200

type LAPICache struct {
	redis *db.RedisCache
}

func NewLAPICache(redisAddr string) *LAPICache {
	return &LAPICache{redis: db.NewRedisCache(redisAddr)}
}

func (c *LAPICache) GetPlate(ctx context.Context, plate string) (*models.PlateHitResult, error) {
	cacheKey := fmt.Sprintf("lapi:plate:%s", plate)
	data, err := c.redis.Get(ctx, cacheKey)
	if err != nil {
		return nil, err
	}
	var hit models.PlateHitResult
	if err := json.Unmarshal([]byte(data), &hit); err != nil {
		return nil, err
	}
	return &hit, nil
}

func (c *LAPICache) SetPlateHit(ctx context.Context, plate string, hit *models.PlateHitResult) error {
	data, err := json.Marshal(hit)
	if err != nil {
		return err
	}
	return c.redis.Set(ctx, fmt.Sprintf("lapi:plate:%s", plate), string(data), 30*time.Second)
}

func (c *LAPICache) SetPlateNoHit(ctx context.Context, plate string) error {
	return c.redis.Set(ctx, fmt.Sprintf("lapi:plate:%s", plate),
		`{"hit_found":false}`, 5*time.Minute)
}

func (c *LAPICache) publishSLABreach(ctx context.Context, plate string, responseMs int) {
	_ = c.redis.Set(ctx, fmt.Sprintf("lapi:sla:breach:%s", plate),
		fmt.Sprintf(`{"plate":"%s","response_ms":%d,"timestamp":%d}`,
			plate, responseMs, time.Now().UnixMilli()), 1*time.Hour)
}

func (c *LAPICache) RecordPlateSighting(ctx context.Context, plate, cameraID string) error {
	key := fmt.Sprintf("lapi:sighting:%s", plate)
	sighting := fmt.Sprintf(`{"plate":"%s","camera":"%s","timestamp":%d}`,
		plate, cameraID, time.Now().UnixMilli())
	return c.redis.Set(ctx, key, sighting, 5*time.Minute)
}

func (c *LAPICache) GetRecentSighting(ctx context.Context, plate string) (string, error) {
	return c.redis.Get(ctx, fmt.Sprintf("lapi:sighting:%s", plate))
}

type LAPIQueryService struct {
	cache *LAPICache
	db    models.Database
}

func NewLAPIQueryService(redisAddr string, db models.Database) *LAPIQueryService {
	return &LAPIQueryService{
		cache: NewLAPICache(redisAddr),
		db:    db,
	}
}

func (s *LAPIQueryService) QueryPlate(ctx context.Context, plate string) (*models.PlateHitResult, error) {
	start := time.Now()
	responseMs := func() int { return int(time.Since(start).Milliseconds()) }

	cached, err := s.cache.GetPlate(ctx, plate)
	if err == nil && cached != nil {
		cached.ResponseMs = responseMs()
		return cached, nil
	}

	queryCtx, cancel := context.WithTimeout(ctx, 150*time.Millisecond)
	defer cancel()

	hit, err := s.db.QueryPlateIndex(queryCtx, plate)
	if err != nil {
		return &models.PlateHitResult{HitFound: false, ResponseMs: responseMs()}, nil
	}
	if hit == nil {
		s.cache.SetPlateNoHit(ctx, plate)
		return &models.PlateHitResult{HitFound: false, ResponseMs: responseMs()}, nil
	}

	s.cache.SetPlateHit(ctx, plate, hit)
	hit.ResponseMs = responseMs()

	if hit.ResponseMs > LAPIMaxResponseMs {
		s.cache.publishSLABreach(ctx, plate, hit.ResponseMs)
	}
	return hit, nil
}

func (s *LAPIQueryService) QueryVIN(ctx context.Context, vin string) (*models.PlateHitResult, error) {
	start := time.Now()
	responseMs := func() int { return int(time.Since(start).Milliseconds()) }

	queryCtx, cancel := context.WithTimeout(ctx, 150*time.Millisecond)
	defer cancel()

	hit, err := s.db.QueryVINIndex(queryCtx, vin)
	if err != nil {
		return &models.PlateHitResult{HitFound: false, ResponseMs: responseMs()}, nil
	}
	if hit == nil {
		return &models.PlateHitResult{HitFound: false, ResponseMs: responseMs()}, nil
	}
	hit.ResponseMs = responseMs()
	return hit, nil
}
