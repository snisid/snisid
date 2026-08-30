package ml

import (
	"fmt"
	"log"
	"sync"
)

// ModelRegistry manages the lifecycle and versioning of AI models.
type ModelRegistry struct {
	mu     sync.RWMutex
	Models map[string]ModelMetadata
}

type ModelMetadata struct {
	Name      string
	Version   string
	Algorithm string
	Status    string // DEPLOYED, SHADOW, RETIRED
}

func NewModelRegistry() *ModelRegistry {
	return &ModelRegistry{
		Models: make(map[string]ModelMetadata),
	}
}

func (r *ModelRegistry) Register(name, version, algo string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Models[name] = ModelMetadata{
		Name:      name,
		Version:   version,
		Algorithm: algo,
		Status:    "DEPLOYED",
	}
	log.Printf("🤖 ML-REGISTRY: Registered model %s (v%s) [%s]", name, version, algo)
}

func (r *ModelRegistry) Get(name string) (ModelMetadata, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	m, ok := r.Models[name]
	if !ok {
		return ModelMetadata{}, fmt.Errorf("MODEL_NOT_FOUND: %s", name)
	}
	return m, nil
}

func (r *ModelRegistry) SwitchToShadow(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if m, ok := r.Models[name]; ok {
		m.Status = "SHADOW"
		r.Models[name] = m
		log.Printf("⚠️ ML-REGISTRY: Model %s moved to SHADOW mode (no production execution)", name)
	}
}
