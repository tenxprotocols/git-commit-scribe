package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// Cache is the interface for caching commit messages
type Cache interface {
	// Get retrieves a cached commit message
	Get(ctx context.Context, key string) (string, bool, error)

	// Set stores a commit message in the cache
	Set(ctx context.Context, key string, value string) error

	// Clear clears the cache
	Clear(ctx context.Context) error

	// Stats returns cache statistics
	Stats() *Stats
}

// Stats holds cache statistics
type Stats struct {
	Entries    int
	MemorySize int64
	DiskSize   int64
	Hits       int64
	Misses     int64
}

// GenerateKey generates a cache key from parameters
func GenerateKey(params ...string) string {
	h := sha256.New()
	for _, p := range params {
		h.Write([]byte(p))
	}
	return hex.EncodeToString(h.Sum(nil))
}

// MultiCache combines memory and disk caches
type MultiCache struct {
	memory Cache
	disk   Cache
}

// NewMultiCache creates a new multi-level cache
func NewMultiCache(memory, disk Cache) *MultiCache {
	return &MultiCache{
		memory: memory,
		disk:   disk,
	}
}

// Get checks memory first, then disk
func (mc *MultiCache) Get(ctx context.Context, key string) (string, bool, error) {
	// Check memory cache first
	if mc.memory != nil {
		if value, found, err := mc.memory.Get(ctx, key); found {
			return value, true, err
		}
	}

	// Check disk cache
	if mc.disk != nil {
		if value, found, err := mc.disk.Get(ctx, key); found {
			// Promote to memory cache
			if mc.memory != nil {
				_ = mc.memory.Set(ctx, key, value)
			}
			return value, true, err
		}
	}

	return "", false, nil
}

// Set stores in both caches
func (mc *MultiCache) Set(ctx context.Context, key string, value string) error {
	var errs []error

	if mc.memory != nil {
		if err := mc.memory.Set(ctx, key, value); err != nil {
			errs = append(errs, fmt.Errorf("memory cache: %w", err))
		}
	}

	if mc.disk != nil {
		if err := mc.disk.Set(ctx, key, value); err != nil {
			errs = append(errs, fmt.Errorf("disk cache: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("cache errors: %v", errs)
	}

	return nil
}

// Clear clears both caches
func (mc *MultiCache) Clear(ctx context.Context) error {
	var errs []error

	if mc.memory != nil {
		if err := mc.memory.Clear(ctx); err != nil {
			errs = append(errs, fmt.Errorf("memory cache: %w", err))
		}
	}

	if mc.disk != nil {
		if err := mc.disk.Clear(ctx); err != nil {
			errs = append(errs, fmt.Errorf("disk cache: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("cache errors: %v", errs)
	}

	return nil
}

// Stats returns combined statistics
func (mc *MultiCache) Stats() *Stats {
	stats := &Stats{}

	if mc.memory != nil {
		memStats := mc.memory.Stats()
		if memStats != nil {
			stats.MemorySize = memStats.MemorySize
			stats.Hits += memStats.Hits
			stats.Misses += memStats.Misses
		}
	}

	if mc.disk != nil {
		diskStats := mc.disk.Stats()
		if diskStats != nil {
			stats.DiskSize = diskStats.DiskSize
			stats.Entries = diskStats.Entries
		}
	}

	return stats
}
