package cli

import (
	"fmt"
	"os"

	"github.com/tenxprotocols/git-commit-scribe/internal/config"
)

// CacheCmd manages cache
type CacheCmd struct {
	Clear CacheClearCmd `cmd:"" help:"Clear cache"`
	Stats CacheStatsCmd `cmd:"" help:"Show cache statistics"`
}

// CacheClearCmd clears the cache
type CacheClearCmd struct {
	All bool `short:"a" help:"Clear both memory and disk cache"`
}

// Run executes the cache clear command
func (c *CacheClearCmd) Run(ctx *Context) error {
	loader := config.NewLoader(CLI.ConfigDir, CLI.ConfigFile)
	cacheDir := config.GetCacheDir(loader.GetConfigDir())

	if c.All {
		// Clear disk cache
		if _, err := os.Stat(cacheDir); err == nil {
			if err := os.RemoveAll(cacheDir); err != nil {
				return fmt.Errorf("failed to clear disk cache: %w", err)
			}
			fmt.Println("✓ Disk cache cleared")
		} else {
			fmt.Println("No disk cache found")
		}
	}

	// Note: Memory cache clearing would require a running instance
	// For now, we'll just clear the disk cache
	if !c.All {
		fmt.Println("Clearing disk cache...")
		if _, err := os.Stat(cacheDir); err == nil {
			if err := os.RemoveAll(cacheDir); err != nil {
				return fmt.Errorf("failed to clear disk cache: %w", err)
			}
			fmt.Println("✓ Disk cache cleared")
		} else {
			fmt.Println("No disk cache found")
		}
	}

	return nil
}

// CacheStatsCmd shows cache statistics
type CacheStatsCmd struct{}

// Run executes the cache stats command
func (c *CacheStatsCmd) Run(ctx *Context) error {
	loader := config.NewLoader(CLI.ConfigDir, CLI.ConfigFile)
	cacheDir := config.GetCacheDir(loader.GetConfigDir())

	fmt.Println("Cache Statistics")
	fmt.Println("═══════════════════════════════════════")

	// Check if cache directory exists
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		fmt.Println("\nNo cache directory found")
		fmt.Printf("Expected location: %s\n", cacheDir)
		return nil
	}

	fmt.Printf("\nCache directory: %s\n", cacheDir)

	// Load cache configuration
	cfg, err := loader.Load()
	if err != nil {
		fmt.Printf("\nWarning: Could not load config: %v\n", err)
	} else {
		fmt.Printf("\nCache settings:\n")
		fmt.Printf("  Enabled: %v\n", cfg.Cache.Enabled)
		fmt.Printf("  TTL: %d seconds\n", cfg.Cache.TTL)
		fmt.Printf("  Max memory: %d MB\n", cfg.Cache.MaxMemoryMB)
		fmt.Printf("  Max disk: %d MB\n", cfg.Cache.MaxDiskMB)
	}

	// Count cache entries
	entries, err := os.ReadDir(cacheDir)
	if err != nil {
		return fmt.Errorf("failed to read cache directory: %w", err)
	}

	var totalSize int64
	fileCount := 0

	for _, entry := range entries {
		if !entry.IsDir() {
			fileCount++
			info, err := entry.Info()
			if err == nil {
				totalSize += info.Size()
			}
		}
	}

	fmt.Printf("\nCache entries: %d\n", fileCount)
	fmt.Printf("Total size: %.2f MB\n", float64(totalSize)/(1024*1024))

	// Show recent cache entries
	if fileCount > 0 {
		fmt.Println("\nRecent cache entries:")
		count := 0
		for i := len(entries) - 1; i >= 0 && count < 5; i-- {
			entry := entries[i]
			if !entry.IsDir() {
				info, err := entry.Info()
				if err == nil {
					fmt.Printf("  • %s (%.2f KB, modified %s)\n",
						entry.Name(),
						float64(info.Size())/1024,
						info.ModTime().Format("2006-01-02 15:04:05"))
					count++
				}
			}
		}
	}

	return nil
}
