package client

import (
	"os"
	"sync"
)

// Client represents an OAuth2 client (agency or third-party application).
type Client struct {
	ID           string
	Secret       string
	Name         string
	RedirectURIs []string
	Scopes       []string
	Active       bool
}

// Manager handles OAuth2 client registration and validation.
type Manager struct {
	mu      sync.RWMutex
	clients map[string]*Client
}

// NewManager creates a new client manager.
func NewManager() *Manager {
	m := &Manager{
		clients: make(map[string]*Client),
	}
	m.registerDefaultClients()
	return m
}

func (m *Manager) registerDefaultClients() {
	m.clients["snisid-web"] = &Client{
		ID:           "snisid-web",
		Secret:       getSecret("SNISID_WEB_SECRET", "snisid-web-secret"),
		Name:         "SNISID Web Portal",
		RedirectURIs: []string{getEnv("SNISID_WEB_REDIRECT_URI", "http://localhost:3000/callback")},
		Scopes:       []string{"openid", "profile", "snisid:identity"},
		Active:       true,
	}
	m.clients["snisid-mobile"] = &Client{
		ID:           "snisid-mobile",
		Secret:       getSecret("SNISID_MOBILE_SECRET", "snisid-mobile-secret"),
		Name:         "SNISID Mobile App",
		RedirectURIs: []string{getEnv("SNISID_MOBILE_REDIRECT_URI", "snisid://callback")},
		Scopes:       []string{"openid", "profile", "snisid:identity", "offline_access"},
		Active:       true,
	}
	m.clients["snisid-agency-api"] = &Client{
		ID:           "snisid-agency-api",
		Secret:       getSecret("SNISID_AGENCY_API_SECRET", "snisid-agency-api-secret"),
		Name:         "SNISID Agency API",
		RedirectURIs: []string{getEnv("SNISID_AGENCY_REDIRECT_URI", "http://localhost:8080/callback")},
		Scopes:       []string{"openid", "snisid:identity", "snisid:biometric", "snisid:document"},
		Active:       true,
	}
}

func getSecret(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// ValidateClient checks client credentials.
func (m *Manager) ValidateClient(clientID, clientSecret string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	client, ok := m.clients[clientID]
	if !ok || !client.Active {
		return false
	}
	if clientSecret != "" && client.Secret != clientSecret {
		return false
	}
	return true
}

// GetRedirectURI returns the first redirect URI for a client.
func (m *Manager) GetRedirectURI(clientID string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	client, ok := m.clients[clientID]
	if !ok || len(client.RedirectURIs) == 0 {
		return ""
	}
	return client.RedirectURIs[0]
}

// GetClientScopes returns the scopes for a client.
func (m *Manager) GetClientScopes(clientID string) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	client, ok := m.clients[clientID]
	if !ok {
		return nil
	}
	return client.Scopes
}

// RegisterClient registers a new OAuth2 client.
func (m *Manager) RegisterClient(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.clients[client.ID] = client
}

// GetClient returns a client by ID.
func (m *Manager) GetClient(clientID string) *Client {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.clients[clientID]
}
