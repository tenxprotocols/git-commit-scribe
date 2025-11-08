package config

import (
	"time"
)

// Config represents the complete configuration
type Config struct {
	Provider string        `yaml:"provider"`
	Model    string        `yaml:"model"`
	APIKey   string        `yaml:"api_key"`
	Cache    CacheConfig   `yaml:"cache"`
	AI       AIConfig      `yaml:"ai"`
	Commit   CommitConfig  `yaml:"commit"`
	Types    map[string]string `yaml:"types,omitempty"`
}

// CacheConfig holds caching configuration
type CacheConfig struct {
	Enabled     bool   `yaml:"enabled"`
	TTL         int    `yaml:"ttl"`          // in seconds
	MaxMemoryMB int    `yaml:"max_memory_mb"`
	MaxDiskMB   int    `yaml:"max_disk_mb"`
}

// AIConfig holds AI provider configuration
type AIConfig struct {
	Temperature    float64       `yaml:"temperature"`
	Timeout        time.Duration `yaml:"timeout"`
	MaxRetries     int           `yaml:"max_retries"`
	RateLimitWait  time.Duration `yaml:"rate_limit_wait"`
}

// CommitConfig holds commit-related configuration
type CommitConfig struct {
	Emoji             bool   `yaml:"emoji"`
	OneLine           bool   `yaml:"one_line"`
	DescriptionLength int    `yaml:"description_length"`
	MaxFiles          int    `yaml:"max_files"`
	IgnoreGenerated   bool   `yaml:"ignore_generated"`
	IgnoreWhitespace  bool   `yaml:"ignore_whitespace"`
	AutoPush          bool   `yaml:"auto_push"`
	Confirm           bool   `yaml:"confirm"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Provider: "openrouter",
		Model:    "anthropic/claude-3.5-sonnet",
		Cache: CacheConfig{
			Enabled:     true,
			TTL:         86400, // 24 hours
			MaxMemoryMB: 100,
			MaxDiskMB:   1000,
		},
		AI: AIConfig{
			Temperature:   0.7,
			Timeout:       30 * time.Second,
			MaxRetries:    3,
			RateLimitWait: 5 * time.Second,
		},
		Commit: CommitConfig{
			Emoji:             false,
			OneLine:           false,
			DescriptionLength: 72,
			MaxFiles:          50,
			IgnoreGenerated:   true,
			IgnoreWhitespace:  true,
			AutoPush:          false,
			Confirm:           true,
		},
		Types: DefaultConventionalTypes(),
	}
}

// DefaultConventionalTypes returns the default commit types
// Based on conventional-changelog-metahub
func DefaultConventionalTypes() map[string]string {
	return map[string]string{
		"feat":     "A new feature",
		"fix":      "A bug fix",
		"docs":     "Documentation only changes",
		"style":    "Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)",
		"refactor": "A code change that neither fixes a bug nor adds a feature",
		"perf":     "A code change that improves performance",
		"test":     "Adding missing tests or correcting existing tests",
		"build":    "Changes that affect the build system or external dependencies",
		"ci":       "Changes to our CI configuration files and scripts",
		"chore":    "Other changes that don't modify src or test files",
		"revert":   "Reverts a previous commit",
	}
}
