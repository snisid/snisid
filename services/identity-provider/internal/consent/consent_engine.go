package consent

import (
	"sync"
)

// Record represents a user's consent grant for a specific client and scope.
type Record struct {
	UserID   string
	ClientID string
	Scope    string
}

// Engine manages user consent records for OAuth2 client access.
type Engine struct {
	mu      sync.RWMutex
	records map[string]*Record
}

// NewEngine creates a new consent engine.
func NewEngine() *Engine {
	return &Engine{
		records: make(map[string]*Record),
	}
}

func consentKey(userID, clientID, scope string) string {
	return userID + ":" + clientID + ":" + scope
}

// HasConsent checks whether a user has granted consent for a client and scope.
func (e *Engine) HasConsent(userID, clientID, scope string) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	_, ok := e.records[consentKey(userID, clientID, scope)]
	return ok
}

// RecordConsent stores a user's consent grant.
func (e *Engine) RecordConsent(userID, clientID, scope string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.records[consentKey(userID, clientID, scope)] = &Record{
		UserID:   userID,
		ClientID: clientID,
		Scope:    scope,
	}
}

// RevokeConsent removes a user's consent grant.
func (e *Engine) RevokeConsent(userID, clientID, scope string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.records, consentKey(userID, clientID, scope))
}

// ListConsents returns all consent records for a user.
func (e *Engine) ListConsents(userID string) []*Record {
	e.mu.RLock()
	defer e.mu.RUnlock()
	var result []*Record
	for _, rec := range e.records {
		if rec.UserID == userID {
			result = append(result, rec)
		}
	}
	return result
}
