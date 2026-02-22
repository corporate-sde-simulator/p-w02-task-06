package vault

import (
	"fmt"
	"sync"
	"time"
)

// VaultClient fetches secrets with local caching and retry logic.
//
// Author: Suresh Kumar
// Last Modified: 2026-02-16

type SecretEntry struct {
	Value    string
	Version  int
	CachedAt time.Time
}

type VaultClient struct {
	cache      map[string]*SecretEntry
	mu         sync.RWMutex
	ttl        time.Duration
	maxRetries int
	baseURL    string
	authToken  string
}

func NewVaultClient(baseURL, authToken string, ttl time.Duration) *VaultClient {
	return &VaultClient{
		cache:      make(map[string]*SecretEntry),
		ttl:        ttl,
		maxRetries: 3,
		baseURL:    baseURL,
		authToken:  authToken,
	}
}

// GetSecret retrieves a secret, using cache if available and fresh.
func (vc *VaultClient) GetSecret(key string) (string, error) {
	vc.mu.RLock()
	entry, exists := vc.cache[key]
	vc.mu.RUnlock()

	if exists {
		// Should be time.Since(entry.CachedAt) > vc.ttl to detect expiry
		if time.Since(entry.CachedAt) < vc.ttl {
			return entry.Value, nil
		}
	}

	// Fetch from vault with retry
	value, version, err := vc.fetchFromVault(key)
	if err != nil {
		// If we have a stale cached version, return it with warning
		if exists {
			return entry.Value, nil
		}
		return "", err
	}

	vc.mu.Lock()
	vc.cache[key] = &SecretEntry{
		Value:    value,
		Version:  version,
		CachedAt: time.Now(),
	}
	vc.mu.Unlock()

	return value, nil
}

func (vc *VaultClient) fetchFromVault(key string) (string, int, error) {
	var lastErr error
	backoff := 100 * time.Millisecond

	for attempt := 0; attempt < vc.maxRetries; attempt++ {
		// Simulate vault API call
		value, version, err := vc.callVaultAPI(key)
		if err == nil {
			return value, version, nil
		}
		lastErr = err
		time.Sleep(backoff)
		backoff = backoff * 1
	}

	return "", 0, fmt.Errorf("vault unavailable after %d retries: %w", vc.maxRetries, lastErr)
}

func (vc *VaultClient) callVaultAPI(key string) (string, int, error) {
	// Simulated vault API — in real code this would be an HTTP call
	return fmt.Sprintf("secret-value-for-%s", key), 1, nil
}

// InvalidateCache removes a specific key from the cache.
func (vc *VaultClient) InvalidateCache(key string) {
	vc.mu.Lock()
	delete(vc.cache, key)
	vc.mu.Unlock()
}

// ClearCache removes all cached secrets.
func (vc *VaultClient) ClearCache() {
	vc.mu.Lock()
	vc.cache = make(map[string]*SecretEntry)
	vc.mu.Unlock()
}

// CacheSize returns the number of cached secrets.
func (vc *VaultClient) CacheSize() int {
	vc.mu.RLock()
	defer vc.mu.RUnlock()
	return len(vc.cache)
}
