package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Loader handles configuration loading with priority
type Loader struct {
	configDir  string
	configFile string
	repoConfig string
}

// NewLoader creates a new configuration loader
func NewLoader(configDir, configFile string) *Loader {
	// Expand home directory
	if configDir != "" {
		configDir = expandPath(configDir)
	} else {
		configDir = expandPath("~/.config/gscribe")
	}

	if configFile != "" {
		configFile = expandPath(configFile)
	} else {
		configFile = filepath.Join(configDir, "config.yaml")
	}

	// Check for repo-level config
	repoConfig := ""
	if gitDir := findGitDir(); gitDir != "" {
		repoConfig = filepath.Join(gitDir, "gscribe-config.yaml")
	}

	return &Loader{
		configDir:  configDir,
		configFile: configFile,
		repoConfig: repoConfig,
	}
}

// Load loads configuration with priority: user config -> repo config
func (l *Loader) Load() (*Config, error) {
	cfg := DefaultConfig()

	// Load user config if it exists
	if _, err := os.Stat(l.configFile); err == nil {
		if err := l.loadFromFile(l.configFile, cfg); err != nil {
			return nil, fmt.Errorf("failed to load user config: %w", err)
		}
	}

	// Load repo config if it exists (overrides user config)
	if l.repoConfig != "" {
		if _, err := os.Stat(l.repoConfig); err == nil {
			if err := l.loadFromFile(l.repoConfig, cfg); err != nil {
				return nil, fmt.Errorf("failed to load repo config: %w", err)
			}
		}
	}

	return cfg, nil
}

// loadFromFile loads configuration from a YAML file
func (l *Loader) loadFromFile(path string, cfg *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, cfg)
}

// Save saves configuration to the user config file
func (l *Loader) Save(cfg *Config) error {
	// Create config directory if it doesn't exist
	if err := os.MkdirAll(l.configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal to YAML
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write with secure permissions (0600 for API key security)
	if err := os.WriteFile(l.configFile, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// SaveRepo saves configuration to the repository config file
func (l *Loader) SaveRepo(cfg *Config) error {
	if l.repoConfig == "" {
		return fmt.Errorf("not in a git repository")
	}

	// Marshal to YAML
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write with secure permissions
	if err := os.WriteFile(l.repoConfig, data, 0600); err != nil {
		return fmt.Errorf("failed to write repo config file: %w", err)
	}

	return nil
}

// GetConfigDir returns the configuration directory
func (l *Loader) GetConfigDir() string {
	return l.configDir
}

// GetConfigFile returns the user configuration file path
func (l *Loader) GetConfigFile() string {
	return l.configFile
}

// GetRepoConfigFile returns the repository configuration file path
func (l *Loader) GetRepoConfigFile() string {
	return l.repoConfig
}

// expandPath expands ~ to home directory
func expandPath(path string) string {
	if len(path) == 0 || path[0] != '~' {
		return path
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}

	if len(path) == 1 {
		return home
	}

	return filepath.Join(home, path[1:])
}

// findGitDir finds the .git directory in current or parent directories
func findGitDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		gitDir := filepath.Join(dir, ".git")
		if info, err := os.Stat(gitDir); err == nil && info.IsDir() {
			return gitDir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root
			return ""
		}
		dir = parent
	}
}

// GetCacheDir returns the cache directory path
func GetCacheDir(configDir string) string {
	// Use XDG cache directory convention
	cacheDir := os.Getenv("XDG_CACHE_HOME")
	if cacheDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return filepath.Join(configDir, "cache")
		}
		cacheDir = filepath.Join(home, ".cache")
	}
	return filepath.Join(cacheDir, "gscribe")
}

// UnmarshalConfig unmarshals YAML data into a Config struct
func UnmarshalConfig(data []byte, cfg *Config) error {
	return yaml.Unmarshal(data, cfg)
}
