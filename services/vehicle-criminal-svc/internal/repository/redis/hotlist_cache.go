package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/snisid/vehicle-criminal-svc/internal/domain"
)

const (
	keyPrefixPlateAlert  = "sivc:plate:"
	keyPrefixStolenPlate = "sivc:stolen:"
	keyHotlistVersion    = "sivc:hotlist:version"
)

type HotlistCache struct {
	client *redis.Client
}

func NewHotlistCache(client *redis.Client) *HotlistCache {
	return &HotlistCache{client: client}
}

func (h *HotlistCache) SetPlateAlert(
	ctx context.Context,
	plate string,
	alert *domain.CriminalAlert,
	ttl time.Duration,
) error {
	key := fmt.Sprintf("%s%s", keyPrefixPlateAlert, plate)
	data, err := json.Marshal(alert)
	if err != nil {
		return err
	}
	return h.client.Set(ctx, key, data, ttl).Err()
}

func (h *HotlistCache) GetPlateAlert(
	ctx context.Context,
	plate string,
) (*domain.CriminalAlert, error) {
	key := fmt.Sprintf("%s%s", keyPrefixPlateAlert, plate)
	data, err := h.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	var alert domain.CriminalAlert
	if err := json.Unmarshal(data, &alert); err != nil {
		return nil, err
	}
	return &alert, nil
}

func (h *HotlistCache) DeletePlateAlert(ctx context.Context, plate string) error {
	key := fmt.Sprintf("%s%s", keyPrefixPlateAlert, plate)
	return h.client.Del(ctx, key).Err()
}

func (h *HotlistCache) BulkLoadHotlist(
	ctx context.Context,
	alerts []*domain.CriminalAlert,
) error {
	pipe := h.client.Pipeline()
	for _, alert := range alerts {
		ttl := 90 * 24 * time.Hour
		if alert.ExpiryDate != nil {
			remaining := time.Until(*alert.ExpiryDate)
			if remaining > 0 {
				ttl = remaining
			}
		}
		key := fmt.Sprintf("%s%s", keyPrefixPlateAlert, alert.PlateNumber)
		data, _ := json.Marshal(alert)
		pipe.Set(ctx, key, data, ttl)
	}
	_, err := pipe.Exec(ctx)
	return err
}
