package vault

import (
	"sync"
	"time"
)

// SecretCache is the thread-safe local cache for vault secrets.
//
// Author: Suresh Kumar
// Last Modified: 2026-02-16

type CacheStats struct {
	Hits   int64
	Misses int64
	Size   int
}

type SecretCache struct {
	entries map[string]*SecretEntry
	mu      sync.RWMutex
	hits    int64
	misses  int64
}

func NewSecretCache() *SecretCache {
	return &SecretCache{
		entries: make(map[string]*SecretEntry),
	}
}

func (sc *SecretCache) Get(key string) (*SecretEntry, bool) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	entry, exists := sc.entries[key]
	if exists {
		sc.hits++
		return entry, true
	}
	sc.misses++
	return nil, false
}

func (sc *SecretCache) Set(key string, entry *SecretEntry) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.entries[key] = entry
}

func (sc *SecretCache) Delete(key string) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	delete(sc.entries, key)
}

func (sc *SecretCache) Clear() {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.entries = make(map[string]*SecretEntry)
}

func (sc *SecretCache) Stats() CacheStats {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	return CacheStats{
		Hits:   sc.hits,
		Misses: sc.misses,
		Size:   len(sc.entries),
	}
}

// Placeholder to keep import used
var _ = time.Now
