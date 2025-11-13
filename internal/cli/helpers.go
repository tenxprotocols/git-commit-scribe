package cli

import (
	"fmt"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/tenxprotocols/git-commit-scribe/internal/config"
	"golang.org/x/term"
)

// readSecureInput reads input from stdin without echoing (for passwords/API keys)
func readSecureInput() (string, error) {
	bytePassword, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		return "", err
	}
	fmt.Println() // Print newline after hidden input
	return strings.TrimSpace(string(bytePassword)), nil
}

// displayConfig displays configuration in a readable format
func displayConfig(cfg *config.Config) {
	fmt.Printf("Provider: %s\n", cfg.Provider)
	fmt.Printf("Model: %s\n", cfg.Model)

	// Mask API key for security
	if cfg.APIKey != "" {
		masked := cfg.APIKey
		if len(masked) > 8 {
			masked = masked[:4] + "..." + masked[len(masked)-4:]
		}
		fmt.Printf("API Key: %s\n", masked)
	} else {
		fmt.Printf("API Key: (not set)\n")
	}

	fmt.Println("\nCache:")
	fmt.Printf("  Enabled: %v\n", cfg.Cache.Enabled)
	fmt.Printf("  TTL: %d seconds\n", cfg.Cache.TTL)
	fmt.Printf("  Max Memory: %d MB\n", cfg.Cache.MaxMemoryMB)
	fmt.Printf("  Max Disk: %d MB\n", cfg.Cache.MaxDiskMB)

	fmt.Println("\nAI:")
	fmt.Printf("  Temperature: %.2f\n", cfg.AI.Temperature)
	fmt.Printf("  Timeout: %s\n", cfg.AI.Timeout)
	fmt.Printf("  Max Retries: %d\n", cfg.AI.MaxRetries)
	fmt.Printf("  Rate Limit Wait: %s\n", cfg.AI.RateLimitWait)

	fmt.Println("\nCommit:")
	fmt.Printf("  Emoji: %v\n", cfg.Commit.Emoji)
	fmt.Printf("  One Line: %v\n", cfg.Commit.OneLine)
	fmt.Printf("  Description Length: %d\n", cfg.Commit.DescriptionLength)
	fmt.Printf("  Max Files: %d\n", cfg.Commit.MaxFiles)
	fmt.Printf("  Ignore Generated: %v\n", cfg.Commit.IgnoreGenerated)
	fmt.Printf("  Ignore Whitespace: %v\n", cfg.Commit.IgnoreWhitespace)
	fmt.Printf("  Auto Push: %v\n", cfg.Commit.AutoPush)
	fmt.Printf("  Confirm: %v\n", cfg.Commit.Confirm)

	if len(cfg.Types) > 0 {
		fmt.Println("\nCommit Types:")
		for t, desc := range cfg.Types {
			fmt.Printf("  %s: %s\n", t, desc)
		}
	}
}

// setConfigValue sets a configuration value based on a dotted key path
func setConfigValue(cfg *config.Config, key, value string) error {
	parts := strings.Split(key, ".")

	switch parts[0] {
	case "provider":
		cfg.Provider = value
	case "model":
		cfg.Model = value
	case "api_key":
		cfg.APIKey = value
	case "cache":
		if len(parts) < 2 {
			return fmt.Errorf("cache key requires a subkey (e.g., cache.enabled)")
		}
		return setCacheValue(&cfg.Cache, parts[1], value)
	case "ai":
		if len(parts) < 2 {
			return fmt.Errorf("ai key requires a subkey (e.g., ai.temperature)")
		}
		return setAIValue(&cfg.AI, parts[1], value)
	case "commit":
		if len(parts) < 2 {
			return fmt.Errorf("commit key requires a subkey (e.g., commit.emoji)")
		}
		return setCommitValue(&cfg.Commit, parts[1], value)
	case "types":
		if len(parts) < 2 {
			return fmt.Errorf("types key requires a type name (e.g., types.feat)")
		}
		if cfg.Types == nil {
			cfg.Types = make(map[string]string)
		}
		cfg.Types[parts[1]] = value
	default:
		return fmt.Errorf("unknown config key: %s", parts[0])
	}

	return nil
}

// setCacheValue sets a cache configuration value
func setCacheValue(cache *config.CacheConfig, key, value string) error {
	switch key {
	case "enabled":
		enabled, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid boolean value for cache.enabled: %s", value)
		}
		cache.Enabled = enabled
	case "ttl":
		ttl, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid integer value for cache.ttl: %s", value)
		}
		cache.TTL = ttl
	case "max_memory_mb":
		maxMemory, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid integer value for cache.max_memory_mb: %s", value)
		}
		cache.MaxMemoryMB = maxMemory
	case "max_disk_mb":
		maxDisk, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid integer value for cache.max_disk_mb: %s", value)
		}
		cache.MaxDiskMB = maxDisk
	default:
		return fmt.Errorf("unknown cache config key: %s", key)
	}
	return nil
}

// setAIValue sets an AI configuration value
func setAIValue(ai *config.AIConfig, key, value string) error {
	switch key {
	case "temperature":
		temp, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("invalid float value for ai.temperature: %s", value)
		}
		ai.Temperature = temp
	case "timeout":
		timeout, err := time.ParseDuration(value)
		if err != nil {
			return fmt.Errorf("invalid duration for ai.timeout: %s", value)
		}
		ai.Timeout = timeout
	case "max_retries":
		retries, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid integer value for ai.max_retries: %s", value)
		}
		ai.MaxRetries = retries
	case "rate_limit_wait":
		wait, err := time.ParseDuration(value)
		if err != nil {
			return fmt.Errorf("invalid duration for ai.rate_limit_wait: %s", value)
		}
		ai.RateLimitWait = wait
	default:
		return fmt.Errorf("unknown AI config key: %s", key)
	}
	return nil
}

// setCommitValue sets a commit configuration value
func setCommitValue(commit *config.CommitConfig, key, value string) error {
	switch key {
	case "emoji":
		emoji, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid boolean value for commit.emoji: %s", value)
		}
		commit.Emoji = emoji
	case "one_line":
		oneLine, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid boolean value for commit.one_line: %s", value)
		}
		commit.OneLine = oneLine
	case "description_length":
		length, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid integer value for commit.description_length: %s", value)
		}
		commit.DescriptionLength = length
	case "max_files":
		maxFiles, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid integer value for commit.max_files: %s", value)
		}
		commit.MaxFiles = maxFiles
	case "ignore_generated":
		ignoreGen, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid boolean value for commit.ignore_generated: %s", value)
		}
		commit.IgnoreGenerated = ignoreGen
	case "ignore_whitespace":
		ignoreWS, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid boolean value for commit.ignore_whitespace: %s", value)
		}
		commit.IgnoreWhitespace = ignoreWS
	case "auto_push":
		autoPush, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid boolean value for commit.auto_push: %s", value)
		}
		commit.AutoPush = autoPush
	case "confirm":
		confirm, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid boolean value for commit.confirm: %s", value)
		}
		commit.Confirm = confirm
	default:
		return fmt.Errorf("unknown commit config key: %s", key)
	}
	return nil
}

// unsetConfigValue resets a configuration value to its default
func unsetConfigValue(cfg, defaultCfg *config.Config, key string) error {
	parts := strings.Split(key, ".")

	switch parts[0] {
	case "provider":
		cfg.Provider = defaultCfg.Provider
	case "model":
		cfg.Model = defaultCfg.Model
	case "api_key":
		cfg.APIKey = ""
	case "cache":
		if len(parts) < 2 {
			cfg.Cache = defaultCfg.Cache
		} else {
			return unsetCacheValue(&cfg.Cache, &defaultCfg.Cache, parts[1])
		}
	case "ai":
		if len(parts) < 2 {
			cfg.AI = defaultCfg.AI
		} else {
			return unsetAIValue(&cfg.AI, &defaultCfg.AI, parts[1])
		}
	case "commit":
		if len(parts) < 2 {
			cfg.Commit = defaultCfg.Commit
		} else {
			return unsetCommitValue(&cfg.Commit, &defaultCfg.Commit, parts[1])
		}
	case "types":
		if len(parts) < 2 {
			cfg.Types = config.DefaultConventionalTypes()
		} else {
			delete(cfg.Types, parts[1])
		}
	default:
		return fmt.Errorf("unknown config key: %s", parts[0])
	}

	return nil
}

// unsetCacheValue resets a cache configuration value to its default
func unsetCacheValue(cache, defaultCache *config.CacheConfig, key string) error {
	switch key {
	case "enabled":
		cache.Enabled = defaultCache.Enabled
	case "ttl":
		cache.TTL = defaultCache.TTL
	case "max_memory_mb":
		cache.MaxMemoryMB = defaultCache.MaxMemoryMB
	case "max_disk_mb":
		cache.MaxDiskMB = defaultCache.MaxDiskMB
	default:
		return fmt.Errorf("unknown cache config key: %s", key)
	}
	return nil
}

// unsetAIValue resets an AI configuration value to its default
func unsetAIValue(ai, defaultAI *config.AIConfig, key string) error {
	switch key {
	case "temperature":
		ai.Temperature = defaultAI.Temperature
	case "timeout":
		ai.Timeout = defaultAI.Timeout
	case "max_retries":
		ai.MaxRetries = defaultAI.MaxRetries
	case "rate_limit_wait":
		ai.RateLimitWait = defaultAI.RateLimitWait
	default:
		return fmt.Errorf("unknown AI config key: %s", key)
	}
	return nil
}

// unsetCommitValue resets a commit configuration value to its default
func unsetCommitValue(commit, defaultCommit *config.CommitConfig, key string) error {
	switch key {
	case "emoji":
		commit.Emoji = defaultCommit.Emoji
	case "one_line":
		commit.OneLine = defaultCommit.OneLine
	case "description_length":
		commit.DescriptionLength = defaultCommit.DescriptionLength
	case "max_files":
		commit.MaxFiles = defaultCommit.MaxFiles
	case "ignore_generated":
		commit.IgnoreGenerated = defaultCommit.IgnoreGenerated
	case "ignore_whitespace":
		commit.IgnoreWhitespace = defaultCommit.IgnoreWhitespace
	case "auto_push":
		commit.AutoPush = defaultCommit.AutoPush
	case "confirm":
		commit.Confirm = defaultCommit.Confirm
	default:
		return fmt.Errorf("unknown commit config key: %s", key)
	}
	return nil
}
