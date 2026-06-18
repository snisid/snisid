package service

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

const cacheDir = ".lapi_cache"

type OfflineCache struct {
	mu      sync.RWMutex
	records map[string]*LocalRecord
}

func NewOfflineCache() *OfflineCache {
	c := &OfflineCache{
		records: make(map[string]*LocalRecord),
	}
	c.loadFromDisk()
	return c
}

func (c *OfflineCache) Store(r *LocalRecord) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.records[r.ID] = r
	return c.saveToDisk()
}

func (c *OfflineCache) Get(id string) (*LocalRecord, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	r, ok := c.records[id]
	if !ok {
		return nil, false
	}
	return r, true
}

func (c *OfflineCache) GetPending() []*LocalRecord {
	c.mu.RLock()
	defer c.mu.RUnlock()
	var pending []*LocalRecord
	for _, r := range c.records {
		if !r.Synced {
			pending = append(pending, r)
		}
	}
	return pending
}

func (c *OfflineCache) saveToDisk() error {
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return err
	}
	data, err := json.Marshal(c.records)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(cacheDir, "records.json"), data, 0644)
}

func (c *OfflineCache) loadFromDisk() {
	data, err := os.ReadFile(filepath.Join(cacheDir, "records.json"))
	if err != nil {
		return
	}
	json.Unmarshal(data, &c.records)
}
