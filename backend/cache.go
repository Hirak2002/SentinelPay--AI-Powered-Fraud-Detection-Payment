package main

import (
	"sync"
	"time"
)

type CacheEntry struct {
	Value      interface{}
	ExpiresAt  time.Time
}

type CacheService struct {
	data map[string]*CacheEntry
	mu   sync.RWMutex
}

func NewCacheService() *CacheService {
	cs := &CacheService{
		data: make(map[string]*CacheEntry),
	}
	
	go cs.cleanupExpired()
	return cs
}

func (cs *CacheService) Set(key string, value interface{}, ttl time.Duration) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	cs.data[key] = &CacheEntry{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}
}

func (cs *CacheService) Get(key string) (interface{}, bool) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	entry, exists := cs.data[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(entry.ExpiresAt) {
		return nil, false
	}

	return entry.Value, true
}

func (cs *CacheService) Delete(key string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	delete(cs.data, key)
}

func (cs *CacheService) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		cs.mu.Lock()
		now := time.Now()
		for key, entry := range cs.data {
			if now.After(entry.ExpiresAt) {
				delete(cs.data, key)
			}
		}
		cs.mu.Unlock()
	}
}
