package vault

import (
	"testing"
	"time"
)

func TestGetSecret_FreshCache(t *testing.T) {
	client := NewVaultClient("http://vault:8200", "test-token", 5*time.Minute)
	val, err := client.GetSecret("db_password")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val == "" {
		t.Error("expected non-empty secret value")
	}
}

func TestGetSecret_UsesCache(t *testing.T) {
	client := NewVaultClient("http://vault:8200", "test-token", 5*time.Minute)
	val1, _ := client.GetSecret("api_key")
	val2, _ := client.GetSecret("api_key")
	if val1 != val2 {
		t.Error("expected same value from cache")
	}
}

func TestGetSecret_CacheExpiry(t *testing.T) {
	client := NewVaultClient("http://vault:8200", "test-token", 1*time.Millisecond)
	client.GetSecret("api_key")
	time.Sleep(10 * time.Millisecond)
	// After TTL, should re-fetch (not use stale cache)
	val, err := client.GetSecret("api_key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val == "" {
		t.Error("expected non-empty value after re-fetch")
	}
}

func TestInvalidateCache(t *testing.T) {
	client := NewVaultClient("http://vault:8200", "test-token", 5*time.Minute)
	client.GetSecret("secret_1")
	if client.CacheSize() != 1 {
		t.Errorf("expected cache size 1, got %d", client.CacheSize())
	}
	client.InvalidateCache("secret_1")
	if client.CacheSize() != 0 {
		t.Errorf("expected cache size 0 after invalidation, got %d", client.CacheSize())
	}
}

func TestClearCache(t *testing.T) {
	client := NewVaultClient("http://vault:8200", "test-token", 5*time.Minute)
	client.GetSecret("s1")
	client.GetSecret("s2")
	client.ClearCache()
	if client.CacheSize() != 0 {
		t.Errorf("expected empty cache after clear, got %d", client.CacheSize())
	}
}

func TestConcurrentAccess(t *testing.T) {
	client := NewVaultClient("http://vault:8200", "test-token", 5*time.Minute)
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			client.GetSecret("concurrent_key")
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}
