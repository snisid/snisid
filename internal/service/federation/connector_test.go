package federation

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHTTPConnector(t *testing.T) {
	c := NewHTTPConnector("oni", "http://localhost:8080", "test-key", 5*time.Second)
	assert.NotNil(t, c)
	assert.Equal(t, "oni", c.AgencyName)
	assert.Equal(t, "http://localhost:8080", c.BaseURL)
	assert.Equal(t, "test-key", c.APIKey)
	assert.NotNil(t, c.Client)
}

func TestHTTPConnector_Name(t *testing.T) {
	c := NewHTTPConnector("dgi", "http://dgi:8080", "key", 5*time.Second)
	assert.Equal(t, "dgi", c.Name())
}

func TestHTTPConnector_Fetch_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer test-key", r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"name": "John", "status": "active"}`))
	}))
	defer server.Close()

	c := NewHTTPConnector("oni", server.URL, "test-key", 5*time.Second)
	result, err := c.Fetch(context.Background(), "NNU-12345")
	require.NoError(t, err)
	assert.Equal(t, "oni", result.Source)
	assert.Equal(t, "John", result.Data["name"])
	assert.Equal(t, "active", result.Data["status"])
}

func TestHTTPConnector_Fetch_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "not found"}`))
	}))
	defer server.Close()

	c := NewHTTPConnector("dgi", server.URL, "key", 5*time.Second)
	_, err := c.Fetch(context.Background(), "unknown")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "returned status 404")
}

func TestHTTPConnector_Fetch_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close()

	c := NewHTTPConnector("oni", server.URL, "key", 1*time.Millisecond)
	_, err := c.Fetch(context.Background(), "test")
	assert.Error(t, err)
}

func TestHTTPConnector_Fetch_ContextCancelled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close()

	c := NewHTTPConnector("oni", server.URL, "key", 5*time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := c.Fetch(ctx, "test")
	assert.Error(t, err)
}
