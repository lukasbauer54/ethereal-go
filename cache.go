package main

import (
	"errors"
	"sync"
	"time"
)

// CacheItem represents a single cached item with expiration
type CacheItem struct {
	Value      interface{}
	Expiration time.Time
}

// Cache provides caching functionality
type Cache struct {
	items map[string]CacheItem
	mu    sync.RWMutex
	ttl   time.Duration
}

// NewCache creates a new Cache instance with a specified TTL (Time To Live)
func NewCache(ttl time.Duration) *Cache {
	return &Cache{
		items: make(map[string]CacheItem),
		ttl:   ttl,
	}
}

// Get retrieves a value from cache
func (c *Cache) Get(key string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, errors.New("key not found in cache")
	}

	if time.Now().After(item.Expiration) {
		delete(c.items, key)
		return nil, errors.New("cache item expired")
	}

	return item.Value, nil
}

// Set stores a value in cache
func (c *Cache) Set(key string, value interface{}) error {
	if value == nil {
		return errors.New("cannot cache nil value")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = CacheItem{
		Value:      value,
		Expiration: time.Now().Add(c.ttl),
	}
	return nil
}

// Delete removes a specific key from the cache
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Clear removes all items from the cache
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]CacheItem)
}

// Cleanup removes expired items from the cache
func (c *Cache) Cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, item := range c.items {
		if now.After(item.Expiration) {
			delete(c.items, key)
		}
	}
}
