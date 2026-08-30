package service

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Template struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Data      []byte    `json:"data"`
	Quality   float64   `json:"quality"`
	CreatedAt time.Time `json:"created_at"`
}

type TemplateManager struct {
	mu        sync.RWMutex
	templates map[string]*Template
}

func NewTemplateManager() *TemplateManager {
	return &TemplateManager{
		templates: make(map[string]*Template),
	}
}

func (m *TemplateManager) Create(userID, fingerprint string, quality float64) (*Template, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := fmt.Sprintf("tmpl_%d", time.Now().UnixNano())
	tmpl := &Template{
		ID:        id,
		UserID:    userID,
		Data:      []byte(fingerprint),
		Quality:   quality,
		CreatedAt: time.Now(),
	}
	m.templates[userID] = tmpl
	m.templates[id] = tmpl
	return tmpl, nil
}

func (m *TemplateManager) Get(key string) (*Template, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	t, ok := m.templates[key]
	return t, ok
}

func (m *TemplateManager) List() []*Template {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var list []*Template
	seen := make(map[string]bool)
	for _, t := range m.templates {
		if !seen[t.ID] {
			list = append(list, t)
			seen[t.ID] = true
		}
	}
	return list
}

func (m *TemplateManager) AssessQuality(imageData string) (float64, error) {
	if len(imageData) == 0 {
		return 0, fmt.Errorf("empty image data")
	}
	quality := 50.0 + rand.Float64()*50.0
	return quality, nil
}
