package cache

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// DiskCache implements a file-based cache with gzip compression
type DiskCache struct {
	dir     string
	maxSize int64
	ttl     time.Duration
}

type diskEntry struct {
	Value     string    `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

// NewDiskCache creates a new disk cache
func NewDiskCache(dir string, maxSizeMB int, ttl time.Duration) (*DiskCache, error) {
	// Create cache directory if it doesn't exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	return &DiskCache{
		dir:     dir,
		maxSize: int64(maxSizeMB) * 1024 * 1024,
		ttl:     ttl,
	}, nil
}

// Get retrieves a value from disk cache
func (dc *DiskCache) Get(ctx context.Context, key string) (string, bool, error) {
	filePath := dc.getFilePath(key)

	// Check if file exists
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return "", false, nil
	}
	if err != nil {
		return "", false, fmt.Errorf("failed to stat cache file: %w", err)
	}

	// Read and decompress file
	file, err := os.Open(filePath)
	if err != nil {
		return "", false, fmt.Errorf("failed to open cache file: %w", err)
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return "", false, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()

	var entry diskEntry
	if err := json.NewDecoder(gzReader).Decode(&entry); err != nil {
		return "", false, fmt.Errorf("failed to decode cache entry: %w", err)
	}

	// Check if entry has expired
	if dc.ttl > 0 && time.Since(entry.Timestamp) > dc.ttl {
		// Delete expired entry
		_ = os.Remove(filePath)
		return "", false, nil
	}

	return entry.Value, true, nil
}

// Set stores a value in disk cache
func (dc *DiskCache) Set(ctx context.Context, key string, value string) error {
	filePath := dc.getFilePath(key)

	// Create entry
	entry := diskEntry{
		Value:     value,
		Timestamp: time.Now(),
	}

	// Create temporary file
	tmpFile, err := os.CreateTemp(dc.dir, "cache-*.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath) // Clean up temp file if we fail

	// Write compressed data
	gzWriter := gzip.NewWriter(tmpFile)
	encoder := json.NewEncoder(gzWriter)
	if err := encoder.Encode(&entry); err != nil {
		tmpFile.Close()
		return fmt.Errorf("failed to encode entry: %w", err)
	}

	if err := gzWriter.Close(); err != nil {
		tmpFile.Close()
		return fmt.Errorf("failed to close gzip writer: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	// Atomic rename
	if err := os.Rename(tmpPath, filePath); err != nil {
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	// Evict old entries if necessary
	if dc.maxSize > 0 {
		go dc.evictIfNeeded()
	}

	return nil
}

// Clear removes all cache files
func (dc *DiskCache) Clear(ctx context.Context) error {
	entries, err := os.ReadDir(dc.dir)
	if err != nil {
		return fmt.Errorf("failed to read cache directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".gz" {
			path := filepath.Join(dc.dir, entry.Name())
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("failed to remove cache file: %w", err)
			}
		}
	}

	return nil
}

// Stats returns cache statistics
func (dc *DiskCache) Stats() *Stats {
	entries, err := os.ReadDir(dc.dir)
	if err != nil {
		return &Stats{}
	}

	var totalSize int64
	count := 0

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".gz" {
			info, err := entry.Info()
			if err != nil {
				continue
			}
			totalSize += info.Size()
			count++
		}
	}

	return &Stats{
		Entries:  count,
		DiskSize: totalSize,
	}
}

// getFilePath returns the file path for a cache key
func (dc *DiskCache) getFilePath(key string) string {
	return filepath.Join(dc.dir, key+".gz")
}

// evictIfNeeded removes old entries if cache size exceeds limit
func (dc *DiskCache) evictIfNeeded() {
	stats := dc.Stats()
	if stats.DiskSize <= dc.maxSize {
		return
	}

	// Get all cache files with timestamps
	type fileInfo struct {
		path    string
		modTime time.Time
		size    int64
	}

	var files []fileInfo
	entries, err := os.ReadDir(dc.dir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".gz" {
			info, err := entry.Info()
			if err != nil {
				continue
			}
			files = append(files, fileInfo{
				path:    filepath.Join(dc.dir, entry.Name()),
				modTime: info.ModTime(),
				size:    info.Size(),
			})
		}
	}

	// Sort by modification time (oldest first)
	for i := 0; i < len(files)-1; i++ {
		for j := i + 1; j < len(files); j++ {
			if files[i].modTime.After(files[j].modTime) {
				files[i], files[j] = files[j], files[i]
			}
		}
	}

	// Remove oldest files until under limit
	currentSize := stats.DiskSize
	for _, file := range files {
		if currentSize <= dc.maxSize {
			break
		}
		if err := os.Remove(file.path); err == nil {
			currentSize -= file.size
		}
	}
}
