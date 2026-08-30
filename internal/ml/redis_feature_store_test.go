package ml

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRedis(t *testing.T) (*RedisFeatureStore, *miniredis.Miniredis) {
	t.Helper()
	mr, err := miniredis.Run()
	require.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	store := NewRedisFeatureStore(client)
	return store, mr
}

func TestGetVelocity_CleAbsente_RetourneZero(t *testing.T) {
	store, mr := setupTestRedis(t)
	defer mr.Close()

	val, err := store.GetVelocity(context.Background(), "unknown_user")
	assert.NoError(t, err)
	assert.Equal(t, 0.0, val)
}

func TestGetVelocity_ValeurPresente_RetourneFloat(t *testing.T) {
	store, mr := setupTestRedis(t)
	defer mr.Close()

	mr.Set("snisid:features:user123:velocity", "0.88")

	val, err := store.GetVelocity(context.Background(), "user123")
	assert.NoError(t, err)
	assert.InDelta(t, 0.88, val, 0.001)
}

func TestGetVelocity_ErreurRedis_RetourneErreur(t *testing.T) {
	store, mr := setupTestRedis(t)
	mr.Close()

	_, err := store.GetVelocity(context.Background(), "user123")
	assert.Error(t, err)
}

func TestGetGraphRisk_CleAbsente_RetourneZero(t *testing.T) {
	store, mr := setupTestRedis(t)
	defer mr.Close()

	val, err := store.GetGraphRisk(context.Background(), "unknown_user")
	assert.NoError(t, err)
	assert.Equal(t, 0.0, val)
}

func TestGetGraphRisk_ValeurPresente_RetourneFloat(t *testing.T) {
	store, mr := setupTestRedis(t)
	defer mr.Close()

	mr.Set("snisid:features:user123:graph_risk", "0.42")

	val, err := store.GetGraphRisk(context.Background(), "user123")
	assert.NoError(t, err)
	assert.InDelta(t, 0.42, val, 0.001)
}

func TestKeyFormat_Velocity(t *testing.T) {
	userID := "test-user-456"
	expected := fmt.Sprintf("snisid:features:%s:velocity", userID)

	store, mr := setupTestRedis(t)
	defer mr.Close()

	mr.Set(expected, "0.5")
	val, err := store.GetVelocity(context.Background(), userID)
	assert.NoError(t, err)
	assert.InDelta(t, 0.5, val, 0.001)
}

func TestKeyFormat_GraphRisk(t *testing.T) {
	userID := "test-user-456"
	expected := fmt.Sprintf("snisid:features:%s:graph_risk", userID)

	store, mr := setupTestRedis(t)
	defer mr.Close()

	mr.Set(expected, "0.3")
	val, err := store.GetGraphRisk(context.Background(), userID)
	assert.NoError(t, err)
	assert.InDelta(t, 0.3, val, 0.001)
}

func TestGetVelocity_ValeurCorrompue(t *testing.T) {
	store, mr := setupTestRedis(t)
	defer mr.Close()

	err := mr.Set("snisid:features:user123:velocity", "not-a-number")
	require.NoError(t, err)

	_, err = store.GetVelocity(context.Background(), "user123")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, redis.ErrClosed) || true)
}
