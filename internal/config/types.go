package config

import (
	"fmt"
	"os"
	"time"
)

const (
	DefaultAnthropicModel  = "claude-sonnet-4-6"
	DefaultOpenRouterModel = "anthropic/claude-3.5-sonnet"
)

// Config represents the complete configuration
type Config struct {
	Provider         string            `yaml:"provider"`
	Model            string            `yaml:"model"`
	APIKey           string            `yaml:"api_key,omitempty"` // backwards compat fallback
	AnthropicAPIKey  string            `yaml:"anthropic_api_key,omitempty"`
	OpenRouterAPIKey string            `yaml:"openrouter_api_key,omitempty"`
	Cache            CacheConfig       `yaml:"cache"`
	AI               AIConfig          `yaml:"ai"`
	Commit           CommitConfig      `yaml:"commit"`
	Types            map[string]string `yaml:"types,omitempty"`
}

// CacheConfig holds caching configuration
type CacheConfig struct {
	Enabled     bool `yaml:"enabled"`
	TTL         int  `yaml:"ttl"` // in seconds
	MaxMemoryMB int  `yaml:"max_memory_mb"`
	MaxDiskMB   int  `yaml:"max_disk_mb"`
}

// AIConfig holds AI provider configuration
type AIConfig struct {
	Temperature   float64       `yaml:"temperature"`
	Timeout       time.Duration `yaml:"timeout"`
	MaxRetries    int           `yaml:"max_retries"`
	RateLimitWait time.Duration `yaml:"rate_limit_wait"`
}

// CommitConfig holds commit-related configuration
type CommitConfig struct {
	Emoji             bool `yaml:"emoji"`
	OneLine           bool `yaml:"one_line"`
	DescriptionLength int  `yaml:"description_length"`
	MaxFiles          int  `yaml:"max_files"`
	IgnoreGenerated   bool `yaml:"ignore_generated"`
	IgnoreWhitespace  bool `yaml:"ignore_whitespace"`
	AutoPush          bool `yaml:"auto_push"`
	Confirm           bool `yaml:"confirm"`
}

// ResolveProvider determines the active provider, API key, and model.
func (c *Config) ResolveProvider() (provider, apiKey, model string, err error) {
	provider = c.Provider

	if provider == "auto" {
		switch {
		case os.Getenv("ANTHROPIC_API_KEY") != "":
			provider = "anthropic"
			apiKey = os.Getenv("ANTHROPIC_API_KEY")
		case os.Getenv("OPENROUTER_API_KEY") != "":
			provider = "openrouter"
			apiKey = os.Getenv("OPENROUTER_API_KEY")
		case c.AnthropicAPIKey != "":
			provider = "anthropic"
			apiKey = c.AnthropicAPIKey
		case c.OpenRouterAPIKey != "":
			provider = "openrouter"
			apiKey = c.OpenRouterAPIKey
		case c.APIKey != "":
			provider = "openrouter"
			apiKey = c.APIKey
		default:
			return "", "", "", fmt.Errorf("no API key configured. Run 'gscribe config init' or set ANTHROPIC_API_KEY or OPENROUTER_API_KEY")
		}
	} else {
		switch provider {
		case "anthropic":
			apiKey = c.AnthropicAPIKey
			if apiKey == "" {
				apiKey = os.Getenv("ANTHROPIC_API_KEY")
			}
			if apiKey == "" {
				apiKey = c.APIKey
			}
		case "openrouter":
			apiKey = c.OpenRouterAPIKey
			if apiKey == "" {
				apiKey = os.Getenv("OPENROUTER_API_KEY")
			}
			if apiKey == "" {
				apiKey = c.APIKey
			}
		default:
			return "", "", "", fmt.Errorf("unsupported provider: %s", provider)
		}
		if apiKey == "" {
			return "", "", "", fmt.Errorf("API key not configured for provider %s", provider)
		}
	}

	model = c.Model
	if model == "" {
		switch provider {
		case "anthropic":
			model = DefaultAnthropicModel
		case "openrouter":
			model = DefaultOpenRouterModel
		}
	}

	return provider, apiKey, model, nil
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Provider: "auto",
		// Model intentionally empty — set per-provider at resolve time
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
