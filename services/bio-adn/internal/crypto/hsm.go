package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"
)

type STRLociData map[string]LociPair

type LociPair struct {
	Value1 string `json:"value1"`
	Value2 string `json:"value2"`
}

type HSMBackend interface {
	WrapKey(dataKey []byte) ([]byte, error)
	UnwrapKey(wrappedKey []byte) ([]byte, error)
	Close() error
}

type SoftwareBackend struct{}

func (s *SoftwareBackend) WrapKey(dataKey []byte) ([]byte, error) {
	out := make([]byte, len(dataKey))
	copy(out, dataKey)
	return out, nil
}

func (s *SoftwareBackend) UnwrapKey(wrappedKey []byte) ([]byte, error) {
	out := make([]byte, len(wrappedKey))
	copy(out, wrappedKey)
	return out, nil
}

func (s *SoftwareBackend) Close() error { return nil }

type HSMCrypto struct {
	mu      sync.RWMutex
	backend HSMBackend
	slotID  uint
}

func NewHSMCrypto(slotID uint) *HSMCrypto {
	return &HSMCrypto{
		backend: &SoftwareBackend{},
		slotID:  slotID,
	}
}

func (h *HSMCrypto) EncryptSTRProfile(profile STRLociData) ([]byte, error) {
	plaintext, err := json.Marshal(profile)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}

	dataKey := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, dataKey); err != nil {
		return nil, fmt.Errorf("key gen: %w", err)
	}

	block, err := aes.NewCipher(dataKey)
	if err != nil {
		return nil, fmt.Errorf("new cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("new gcm: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	h.mu.RLock()
	wrappedKey, err := h.backend.WrapKey(dataKey)
	h.mu.RUnlock()
	if err != nil {
		return nil, fmt.Errorf("hsm wrap: %w", err)
	}

	result := make([]byte, 4+len(wrappedKey)+len(ciphertext))
	binary.BigEndian.PutUint32(result[:4], uint32(len(wrappedKey)))
	copy(result[4:], wrappedKey)
	copy(result[4+len(wrappedKey):], ciphertext)

	return result, nil
}

func (h *HSMCrypto) DecryptSTRProfile(data []byte) (STRLociData, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("invalid data: too short")
	}

	keyLen := binary.BigEndian.Uint32(data[:4])
	if len(data) < 4+int(keyLen) {
		return nil, fmt.Errorf("invalid data: key length mismatch")
	}

	wrappedKey := data[4 : 4+keyLen]
	ciphertext := data[4+keyLen:]

	h.mu.RLock()
	dataKey, err := h.backend.UnwrapKey(wrappedKey)
	h.mu.RUnlock()
	if err != nil {
		return nil, fmt.Errorf("hsm unwrap: %w", err)
	}

	block, err := aes.NewCipher(dataKey)
	if err != nil {
		return nil, fmt.Errorf("new cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("new gcm: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ct := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return nil, fmt.Errorf("decrypt: %w", err)
	}

	var profile STRLociData
	if err := json.Unmarshal(plaintext, &profile); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return profile, nil
}

func (h *HSMCrypto) SetBackend(backend HSMBackend) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.backend = backend
}

func (h *HSMCrypto) Close() error {
	return h.backend.Close()
}

// ── HSM Session with key lifetime management ──────────────────────────────

type HSMSession struct {
	Backend   HSMBackend
	CreatedAt time.Time
}

var hsmSession *HSMSession
var hsmMu sync.Mutex

func GetHSMSession() (*HSMSession, error) {
	hsmMu.Lock()
	defer hsmMu.Unlock()

	if hsmSession != nil && time.Since(hsmSession.CreatedAt) < 30*time.Second {
		return hsmSession, nil
	}

	if hsmSession != nil {
		_ = hsmSession.Backend.Close()
	}

	hsmSession = &HSMSession{
		Backend:   &SoftwareBackend{},
		CreatedAt: time.Now(),
	}
	return hsmSession, nil
}
