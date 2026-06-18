package verification

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockBiometricConnector_Name(t *testing.T) {
	c := &MockBiometricConnector{}
	assert.Equal(t, "biometric", c.Name())
}

func TestMockBiometricConnector_Verify_Success(t *testing.T) {
	c := &MockBiometricConnector{}
	data := map[string]interface{}{
		"biometricData": "face-encoding-123",
	}
	result, err := c.Verify(context.Background(), data)
	assert.NoError(t, err)
	assert.Equal(t, StatusSuccess, result.Status)
	assert.Equal(t, 98, result.Score)
	assert.Contains(t, result.Reason, "Face match confirmed")
}

func TestMockBiometricConnector_Verify_NoBiometricData(t *testing.T) {
	c := &MockBiometricConnector{}
	result, err := c.Verify(context.Background(), map[string]interface{}{})
	assert.NoError(t, err)
	assert.Equal(t, StatusFailed, result.Status)
	assert.Equal(t, 0, result.Score)
	assert.Contains(t, result.Reason, "No biometric data")
}

func TestMockAgencyConnector_Name(t *testing.T) {
	c := &MockAgencyConnector{AgencyName: "dgi"}
	assert.Equal(t, "dgi", c.Name())
}

func TestMockAgencyConnector_Verify_Success(t *testing.T) {
	c := &MockAgencyConnector{AgencyName: "oni"}
	result, err := c.Verify(context.Background(), map[string]interface{}{
		"identityId": "valid-id",
	})
	assert.NoError(t, err)
	assert.Equal(t, StatusSuccess, result.Status)
	assert.Equal(t, 100, result.Score)
}

func TestMockAgencyConnector_Verify_RevokedIdentity(t *testing.T) {
	c := &MockAgencyConnector{AgencyName: "oni"}
	result, err := c.Verify(context.Background(), map[string]interface{}{
		"identityId": "revoked-id",
	})
	assert.NoError(t, err)
	assert.Equal(t, StatusFailed, result.Status)
	assert.Equal(t, 0, result.Score)
	assert.Contains(t, result.Reason, "revoked")
}

func TestMockAgencyConnector_Verify_EmptyData(t *testing.T) {
	c := &MockAgencyConnector{AgencyName: "police"}
	result, err := c.Verify(context.Background(), map[string]interface{}{})
	assert.NoError(t, err)
	assert.Equal(t, StatusSuccess, result.Status)
}
