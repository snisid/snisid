package verification_test

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/snisid/platform/internal/service/verification"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockConnector struct {
	name   string
	result verification.Result
	err    error
	delay  chan struct{}
}

func (c *mockConnector) Name() string { return c.name }

func (c *mockConnector) Verify(ctx context.Context, data map[string]interface{}) (verification.Result, error) {
	if c.delay != nil {
		select {
		case <-c.delay:
		case <-ctx.Done():
			return verification.Result{}, ctx.Err()
		}
	}
	if c.err != nil {
		return verification.Result{}, c.err
	}
	return c.result, nil
}

func TestOrchestrator_VerifyIdentity_AllSuccess(t *testing.T) {
	conn1 := &mockConnector{name: "biometric", result: verification.Result{Status: verification.StatusSuccess, Score: 98}}
	conn2 := &mockConnector{name: "dgi", result: verification.Result{Status: verification.StatusSuccess, Score: 100}}

	o := verification.NewOrchestrator(conn1, conn2)
	results, err := o.VerifyIdentity(context.Background(), map[string]interface{}{"id": "123"})
	require.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, 98, results["biometric"].Score)
	assert.Equal(t, 100, results["dgi"].Score)
}

func TestOrchestrator_VerifyIdentity_OneFails(t *testing.T) {
	conn1 := &mockConnector{name: "good", result: verification.Result{Status: verification.StatusSuccess, Score: 95}}
	conn2 := &mockConnector{name: "bad", err: errors.New("service unavailable")}

	o := verification.NewOrchestrator(conn1, conn2)
	_, err := o.VerifyIdentity(context.Background(), map[string]interface{}{"id": "123"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "verification orchestration failed")
}

func TestOrchestrator_VerifyIdentity_AllFail(t *testing.T) {
	conn1 := &mockConnector{name: "a", err: errors.New("err1")}
	conn2 := &mockConnector{name: "b", err: errors.New("err2")}

	o := verification.NewOrchestrator(conn1, conn2)
	_, err := o.VerifyIdentity(context.Background(), map[string]interface{}{})
	assert.Error(t, err)
}

func TestOrchestrator_VerifyIdentity_EmptyConnectors(t *testing.T) {
	o := verification.NewOrchestrator()
	results, err := o.VerifyIdentity(context.Background(), map[string]interface{}{"id": "123"})
	require.NoError(t, err)
	assert.Empty(t, results)
}

func TestOrchestrator_VerifyIdentity_ContextCancelled(t *testing.T) {
	conn := &mockConnector{
		name:  "slow",
		delay: make(chan struct{}),
	}

	o := verification.NewOrchestrator(conn)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := o.VerifyIdentity(ctx, map[string]interface{}{"id": "123"})
	assert.Error(t, err)
}

func TestOrchestrator_VerifyIdentity_PartialResultsOnCancel(t *testing.T) {
	fast := &mockConnector{
		name:   "fast",
		result: verification.Result{Status: verification.StatusSuccess, Score: 100},
	}
	slow := &mockConnector{
		name:  "slow",
		delay: make(chan struct{}),
	}

	o := verification.NewOrchestrator(fast, slow)
	ctx, cancel := context.WithCancel(context.Background())

	errCh := make(chan error, 1)
	go func() {
		_, err := o.VerifyIdentity(ctx, map[string]interface{}{"id": "123"})
		errCh <- err
	}()

	cancel()
	err := <-errCh
	assert.Error(t, err)
}

func TestOrchestrator_ConcurrentSafety(t *testing.T) {
	conn := &mockConnector{
		name:   "safe",
		result: verification.Result{Status: verification.StatusSuccess, Score: 50},
	}

	o := verification.NewOrchestrator(conn)
	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := o.VerifyIdentity(context.Background(), map[string]interface{}{"id": "concurrent"})
			assert.NoError(t, err)
		}()
	}
	wg.Wait()
}

func TestMockBiometricConnector_MissingData(t *testing.T) {
	c := &verification.MockBiometricConnector{}
	result, err := c.Verify(context.Background(), map[string]interface{}{"no_biometric": "data"})
	require.NoError(t, err)
	assert.Equal(t, verification.StatusFailed, result.Status)
	assert.Equal(t, 0, result.Score)
	assert.Contains(t, result.Reason, "No biometric data")
}

func TestMockAgencyConnector_RevokedID(t *testing.T) {
	c := &verification.MockAgencyConnector{AgencyName: "police"}
	result, err := c.Verify(context.Background(), map[string]interface{}{
		"identityId": "revoked-id",
	})
	require.NoError(t, err)
	assert.Equal(t, verification.StatusFailed, result.Status)
	assert.Equal(t, 0, result.Score)
}

func TestMockAgencyConnector_ValidID(t *testing.T) {
	c := &verification.MockAgencyConnector{AgencyName: "oni"}
	result, err := c.Verify(context.Background(), map[string]interface{}{
		"identityId": "valid-id-123",
	})
	require.NoError(t, err)
	assert.Equal(t, verification.StatusSuccess, result.Status)
	assert.Equal(t, 100, result.Score)
}

func TestMockBiometricConnector_WithData(t *testing.T) {
	c := &verification.MockBiometricConnector{}
	result, err := c.Verify(context.Background(), map[string]interface{}{
		"biometricData": "face-vector-xyz",
	})
	require.NoError(t, err)
	assert.Equal(t, verification.StatusSuccess, result.Status)
	assert.Equal(t, 98, result.Score)
}

func TestConnectorInterface(t *testing.T) {
	var conn verification.Connector = &mockConnector{
		name:   "test-connector",
		result: verification.Result{Status: verification.StatusSuccess, Score: 100},
	}
	assert.Equal(t, "test-connector", conn.Name())
}

func TestStatusConstants(t *testing.T) {
	assert.Equal(t, verification.CheckStatus("SUCCESS"), verification.StatusSuccess)
	assert.Equal(t, verification.CheckStatus("FAILED"), verification.StatusFailed)
	assert.Equal(t, verification.CheckStatus("ERROR"), verification.StatusError)
}
