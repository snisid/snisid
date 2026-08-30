package session

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"

	"github.com/snisid/platform/services/identity-provider/internal/oidc"
)

// Data represents an authenticated SSO session.
type Data struct {
	oidc.SessionData
	ID        string
	CreatedAt time.Time
	ExpiresAt time.Time
}

// Store manages SSO sessions in memory.
type Store struct {
	mu       sync.RWMutex
	sessions map[string]*Data
}

// NewStore creates a new SSO session store.
func NewStore() *Store {
	return &Store{
		sessions: make(map[string]*Data),
	}
}

func generateSessionID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return "sess_" + hex.EncodeToString(b)
}

// CreateSession creates a new SSO session and returns its ID.
func (s *Store) CreateSession(userID, clientID string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := generateSessionID()
	s.sessions[id] = &Data{
		SessionData: oidc.SessionData{
			UserID:   userID,
			ClientID: clientID,
			Subject:  userID,
			Scopes:   []string{"openid", "profile"},
		},
		ID:        id,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	return id
}

// GetSession retrieves an SSO session by ID.
func (s *Store) GetSession(sessionID string) *oidc.SessionData {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sess, ok := s.sessions[sessionID]
	if !ok {
		return nil
	}
	if time.Now().After(sess.ExpiresAt) {
		return nil
	}
	return &sess.SessionData
}

// DeleteSession removes an SSO session.
func (s *Store) DeleteSession(sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, sessionID)
}

// ListSessions returns all active sessions for a user.
func (s *Store) ListSessions(userID string) []*Data {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*Data
	for _, sess := range s.sessions {
		if sess.UserID == userID && time.Now().Before(sess.ExpiresAt) {
			result = append(result, sess)
		}
	}
	return result
}
