package ethereal

import (
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	ttl := 5 * time.Minute
	cache := NewCache(ttl)

	if cache == nil {
		t.Error("NewCache returned nil")
	}
	if cache.ttl != ttl {
		t.Errorf("Expected TTL %v, got %v", ttl, cache.ttl)
	}
	if cache.items == nil {
		t.Error("Cache items map not initialized")
	}
}

func TestCacheSetAndGet(t *testing.T) {
	cache := NewCache(100 * time.Millisecond)

	// Test setting and getting a value
	err := cache.Set("key1", "value1")
	if err != nil {
		t.Errorf("Failed to set cache value: %v", err)
	}

	value, err := cache.Get("key1")
	if err != nil {
		t.Errorf("Failed to get cache value: %v", err)
	}
	if value != "value1" {
		t.Errorf("Expected 'value1', got %v", value)
	}

	// Test getting non-existent key
	_, err = cache.Get("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent key, got nil")
	}

	// Test setting nil value
	err = cache.Set("key2", nil)
	if err == nil {
		t.Error("Expected error when setting nil value, got nil")
	}
}

func TestCacheExpiration(t *testing.T) {
	cache := NewCache(100 * time.Millisecond)

	// Set a value
	err := cache.Set("key1", "value1")
	if err != nil {
		t.Errorf("Failed to set cache value: %v", err)
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Try to get expired value
	_, err = cache.Get("key1")
	if err == nil {
		t.Error("Expected error for expired key, got nil")
	}
}

func TestCacheDelete(t *testing.T) {
	cache := NewCache(time.Minute)

	// Set and delete a value
	cache.Set("key1", "value1")
	cache.Delete("key1")

	// Try to get deleted value
	_, err := cache.Get("key1")
	if err == nil {
		t.Error("Expected error after deletion, got nil")
	}
}

func TestCacheClear(t *testing.T) {
	cache := NewCache(time.Minute)

	// Set multiple values
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	// Clear the cache
	cache.Clear()

	// Try to get cleared values
	_, err1 := cache.Get("key1")
	_, err2 := cache.Get("key2")
	if err1 == nil || err2 == nil {
		t.Error("Expected errors after clearing cache, got nil")
	}
}

func TestCacheCleanup(t *testing.T) {
	cache := NewCache(100 * time.Millisecond)

	// Set multiple values
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Run cleanup
	cache.Cleanup()

	// Try to get cleaned up values
	_, err1 := cache.Get("key1")
	_, err2 := cache.Get("key2")
	if err1 == nil || err2 == nil {
		t.Error("Expected errors after cleanup, got nil")
	}
}

func TestCacheConcurrency(t *testing.T) {
	cache := NewCache(time.Minute)
	done := make(chan bool)

	// Concurrent reads and writes
	go func() {
		for i := 0; i < 100; i++ {
			cache.Set("key", i)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 100; i++ {
			cache.Get("key")
		}
		done <- true
	}()

	// Wait for both goroutines to complete
	<-done
	<-done
}
