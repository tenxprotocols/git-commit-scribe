package cache

import (
	"context"
	"sync"
	"time"
)

// MemoryCache implements an in-memory LRU cache
type MemoryCache struct {
	mu          sync.RWMutex
	entries     map[string]*cacheEntry
	maxSize     int64
	currentSize int64
	ttl         time.Duration
	hits        int64
	misses      int64
}

type cacheEntry struct {
	value     string
	size      int64
	timestamp time.Time
}

// NewMemoryCache creates a new memory cache
func NewMemoryCache(maxSizeMB int, ttl time.Duration) *MemoryCache {
	return &MemoryCache{
		entries: make(map[string]*cacheEntry),
		maxSize: int64(maxSizeMB) * 1024 * 1024, // Convert MB to bytes
		ttl:     ttl,
	}
}

// Get retrieves a value from the cache
func (mc *MemoryCache) Get(_ context.Context, key string) (string, bool, error) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	entry, found := mc.entries[key]
	if !found {
		mc.misses++
		return "", false, nil
	}

	// Check if entry has expired
	if mc.ttl > 0 && time.Since(entry.timestamp) > mc.ttl {
		mc.misses++
		return "", false, nil
	}

	mc.hits++
	return entry.value, true, nil
}

// Set stores a value in the cache
func (mc *MemoryCache) Set(_ context.Context, key string, value string) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	size := int64(len(value))

	// Remove old entry if exists
	if oldEntry, exists := mc.entries[key]; exists {
		mc.currentSize -= oldEntry.size
	}

	// Evict entries if necessary
	for mc.currentSize+size > mc.maxSize && len(mc.entries) > 0 {
		mc.evictOldest()
	}

	// Add new entry
	mc.entries[key] = &cacheEntry{
		value:     value,
		size:      size,
		timestamp: time.Now(),
	}
	mc.currentSize += size

	return nil
}

// Clear clears all entries
func (mc *MemoryCache) Clear(_ context.Context) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.entries = make(map[string]*cacheEntry)
	mc.currentSize = 0
	return nil
}

// Stats returns cache statistics
func (mc *MemoryCache) Stats() *Stats {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	return &Stats{
		Entries:    len(mc.entries),
		MemorySize: mc.currentSize,
		Hits:       mc.hits,
		Misses:     mc.misses,
	}
}

// evictOldest removes the oldest entry from the cache
func (mc *MemoryCache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time

	for key, entry := range mc.entries {
		if oldestKey == "" || entry.timestamp.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.timestamp
		}
	}

	if oldestKey != "" {
		mc.currentSize -= mc.entries[oldestKey].size
		delete(mc.entries, oldestKey)
	}
}
