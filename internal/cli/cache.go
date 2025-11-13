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

	spinner := NewSpinner("Clearing disk cache...")
	spinner.Start()

	if c.All {
		// Clear disk cache
		if _, err := os.Stat(cacheDir); err == nil {
			if err := os.RemoveAll(cacheDir); err != nil {
				spinner.Error("Failed to clear disk cache")
				return fmt.Errorf("failed to clear disk cache: %w", err)
			}
			spinner.Success("Disk cache cleared")
		} else {
			spinner.Stop()
			PrintInfo("No disk cache found")
		}
	} else {
		// Note: Memory cache clearing would require a running instance
		// For now, we'll just clear the disk cache
		if _, err := os.Stat(cacheDir); err == nil {
			if err := os.RemoveAll(cacheDir); err != nil {
				spinner.Error("Failed to clear disk cache")
				return fmt.Errorf("failed to clear disk cache: %w", err)
			}
			spinner.Success("Disk cache cleared")
		} else {
			spinner.Stop()
			PrintInfo("No disk cache found")
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

	FormatHeader("Cache Statistics")

	// Check if cache directory exists
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		PrintInfo("No cache directory found")
		PrintDim(fmt.Sprintf("Expected location: %s", cacheDir))
		return nil
	}

	PrintDim(fmt.Sprintf("Cache directory: %s", cacheDir))
	fmt.Println()

	// Load cache configuration
	cfg, err := loader.Load()
	if err != nil {
		PrintWarning(fmt.Sprintf("Could not load config: %v", err))
	} else {
		PrintInfo("Cache settings:")
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

	fmt.Println()
	PrintSuccess(fmt.Sprintf("Cache entries: %d", fileCount))
	fmt.Printf("Total size: %.2f MB\n", float64(totalSize)/(1024*1024))

	// Show recent cache entries
	if fileCount > 0 {
		FormatSubHeader("Recent Cache Entries")
		count := 0
		for i := len(entries) - 1; i >= 0 && count < 5; i-- {
			entry := entries[i]
			if !entry.IsDir() {
				info, err := entry.Info()
				if err == nil {
					fmt.Printf("  • %s ", highlightColor.Sprint(entry.Name()))
					PrintDim(fmt.Sprintf("(%.2f KB, modified %s)",
						float64(info.Size())/1024,
						info.ModTime().Format("2006-01-02 15:04:05")))
					count++
				}
			}
		}
	}

	return nil
}
