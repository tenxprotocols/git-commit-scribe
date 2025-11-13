package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tenxprotocols/git-commit-scribe/internal/config"
)

// ConfigCmd manages configuration
type ConfigCmd struct {
	Init  ConfigInitCmd  `cmd:"" help:"Initialize configuration"`
	Show  ConfigShowCmd  `cmd:"" help:"Show current configuration"`
	Set   ConfigSetCmd   `cmd:"" help:"Set a configuration value"`
	Unset ConfigUnsetCmd `cmd:"" help:"Unset a configuration value"`
	Repos ConfigReposCmd `cmd:"" help:"Manage per-repository configuration"`
}

// ConfigInitCmd initializes configuration
type ConfigInitCmd struct {
	Force bool `short:"f" help:"Overwrite existing configuration"`
}

// Run executes the config init command
func (c *ConfigInitCmd) Run(_ *Context) error {
	loader := config.NewLoader(CLI.ConfigDir, CLI.ConfigFile)
	configFile := loader.GetConfigFile()

	// Check if config already exists
	if !c.Force {
		if _, err := os.Stat(configFile); err == nil {
			return fmt.Errorf("configuration already exists at %s. Use --force to overwrite", configFile)
		}
	}

	FormatHeader("Initialize Configuration")
	PrintDim(fmt.Sprintf("Config file: %s", configFile))
	fmt.Println()

	// Create default config
	cfg := config.DefaultConfig()

	// Prompt for API key
	PrintInfo("OpenRouter API Key")
	PrintDim("  Get your API key from: https://openrouter.ai/keys")
	fmt.Print(promptColor.Sprint("  Enter API key (input hidden): "))

	apiKey, err := readSecureInput()
	if err != nil {
		return fmt.Errorf("failed to read API key: %w", err)
	}

	if apiKey == "" {
		return fmt.Errorf("API key is required")
	}

	cfg.APIKey = apiKey

	// Prompt for model
	fmt.Print(promptColor.Sprint("\nModel [anthropic/claude-3.5-sonnet]: "))
	reader := bufio.NewReader(os.Stdin)
	model, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read model: %w", err)
	}

	model = strings.TrimSpace(model)
	if model != "" {
		cfg.Model = model
	}

	// Save configuration
	spinner := NewSpinner("Saving configuration...")
	spinner.Start()

	if err := loader.Save(cfg); err != nil {
		spinner.Error("Failed to save configuration")
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	spinner.Success(fmt.Sprintf("Configuration saved to %s", configFile))
	fmt.Println()
	PrintSuccess("You can now use gscribe to generate commit messages!")
	PrintDim("Try: gscribe commit --help")

	return nil
}

// ConfigShowCmd shows current configuration
type ConfigShowCmd struct {
	Global bool `short:"g" help:"Show global configuration only"`
}

// Run executes the config show command
func (c *ConfigShowCmd) Run(_ *Context) error {
	loader := config.NewLoader(CLI.ConfigDir, CLI.ConfigFile)

	if c.Global {
		// Show only global config
		configFile := loader.GetConfigFile()
		if _, err := os.Stat(configFile); err != nil {
			return fmt.Errorf("no global configuration found at %s", configFile)
		}

		data, err := os.ReadFile(configFile)
		if err != nil {
			return fmt.Errorf("failed to read config file: %w", err)
		}

		FormatHeader("Global Configuration")
		PrintDim(configFile)
		fmt.Println()
		fmt.Println(string(data))
		return nil
	}

	// Load merged configuration
	cfg, err := loader.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	FormatHeader("Current Configuration")
	PrintDim("(merged from all sources)")
	fmt.Println()
	displayConfig(cfg)

	FormatSubHeader("Configuration Sources")
	PrintDim("Priority order (highest to lowest):")
	fmt.Printf("  1. Command-line flags\n")
	fmt.Printf("  2. Environment variables (GSCRIBE_*)\n")
	fmt.Printf("  3. Repository config: %s\n", loader.GetRepoConfigFile())
	fmt.Printf("  4. User config: %s\n", loader.GetConfigFile())
	fmt.Printf("  5. Default values\n")

	return nil
}

// ConfigSetCmd sets a configuration value
type ConfigSetCmd struct {
	Key   string `arg:"" help:"Configuration key (e.g., model, provider, commit.emoji)"`
	Value string `arg:"" help:"Configuration value"`
	Repo  bool   `short:"r" help:"Set for current repository only"`
}

// Run executes the config set command
func (c *ConfigSetCmd) Run(_ *Context) error {
	loader := config.NewLoader(CLI.ConfigDir, CLI.ConfigFile)

	var cfg *config.Config

	if c.Repo {
		// Set in repository config
		repoConfig := loader.GetRepoConfigFile()
		if repoConfig == "" {
			return fmt.Errorf("not in a git repository")
		}

		// Load existing repo config or create new one
		cfg = config.DefaultConfig()
		if _, err := os.Stat(repoConfig); err == nil {
			data, err := os.ReadFile(repoConfig)
			if err != nil {
				return fmt.Errorf("failed to read repo config: %w", err)
			}
			if err := config.UnmarshalConfig(data, cfg); err != nil {
				return fmt.Errorf("failed to parse repo config: %w", err)
			}
		}
	} else {
		// Set in user config
		targetFile := loader.GetConfigFile()

		// Load existing user config or create new one
		cfg = config.DefaultConfig()
		if _, err := os.Stat(targetFile); err == nil {
			data, err := os.ReadFile(targetFile)
			if err != nil {
				return fmt.Errorf("failed to read user config: %w", err)
			}
			if err := config.UnmarshalConfig(data, cfg); err != nil {
				return fmt.Errorf("failed to parse user config: %w", err)
			}
		}
	}

	// Update the configuration value
	if err := setConfigValue(cfg, c.Key, c.Value); err != nil {
		return err
	}

	// Save the configuration
	var err error
	if c.Repo {
		err = loader.SaveRepo(cfg)
	} else {
		err = loader.Save(cfg)
	}

	if err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	scope := "user"
	if c.Repo {
		scope = "repository"
	}
	PrintSuccess(fmt.Sprintf("Set %s config: %s = %s", scope, c.Key, c.Value))

	return nil
}

// ConfigUnsetCmd unsets a configuration value
type ConfigUnsetCmd struct {
	Key  string `arg:"" help:"Configuration key (e.g., model, provider, commit.emoji)"`
	Repo bool   `short:"r" help:"Unset for current repository only"`
}

// Run executes the config unset command
func (c *ConfigUnsetCmd) Run(_ *Context) error {
	loader := config.NewLoader(CLI.ConfigDir, CLI.ConfigFile)

	var cfg *config.Config

	if c.Repo {
		// Unset in repository config
		repoConfig := loader.GetRepoConfigFile()
		if repoConfig == "" {
			return fmt.Errorf("not in a git repository")
		}

		if _, err := os.Stat(repoConfig); err != nil {
			return fmt.Errorf("no repository configuration found")
		}

		cfg = config.DefaultConfig()
		data, err := os.ReadFile(repoConfig)
		if err != nil {
			return fmt.Errorf("failed to read repo config: %w", err)
		}
		if err := config.UnmarshalConfig(data, cfg); err != nil {
			return fmt.Errorf("failed to parse repo config: %w", err)
		}
	} else {
		// Unset in user config
		targetFile := loader.GetConfigFile()

		if _, err := os.Stat(targetFile); err != nil {
			return fmt.Errorf("no user configuration found")
		}

		cfg = config.DefaultConfig()
		data, err := os.ReadFile(targetFile)
		if err != nil {
			return fmt.Errorf("failed to read user config: %w", err)
		}
		if err := config.UnmarshalConfig(data, cfg); err != nil {
			return fmt.Errorf("failed to parse user config: %w", err)
		}
	}

	// Reset the value to default
	defaultCfg := config.DefaultConfig()
	if err := unsetConfigValue(cfg, defaultCfg, c.Key); err != nil {
		return err
	}

	// Save the configuration
	var err error
	if c.Repo {
		err = loader.SaveRepo(cfg)
	} else {
		err = loader.Save(cfg)
	}

	if err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	scope := "user"
	if c.Repo {
		scope = "repository"
	}
	PrintSuccess(fmt.Sprintf("Unset %s config: %s", scope, c.Key))

	return nil
}

// ConfigReposCmd manages per-repository configuration
type ConfigReposCmd struct {
	List ConfigReposListCmd `cmd:"" help:"List repositories with custom configuration"`
}

// ConfigReposListCmd lists repositories with custom configuration
type ConfigReposListCmd struct{}

// Run executes the config repos list command
func (c *ConfigReposListCmd) Run(_ *Context) error {
	loader := config.NewLoader(CLI.ConfigDir, CLI.ConfigFile)
	configDir := loader.GetConfigDir()

	var repoConfigs []string

	// Check current repository
	currentRepoConfig := loader.GetRepoConfigFile()
	if currentRepoConfig != "" {
		if _, err := os.Stat(currentRepoConfig); err == nil {
			repoConfigs = append(repoConfigs, currentRepoConfig)
		}
	}

	// Check for saved repository configs
	reposDir := filepath.Join(configDir, "repos")
	if _, err := os.Stat(reposDir); err == nil {
		entries, err := os.ReadDir(reposDir)
		if err == nil {
			for _, entry := range entries {
				if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".yaml") {
					repoConfigs = append(repoConfigs, filepath.Join(reposDir, entry.Name()))
				}
			}
		}
	}

	if len(repoConfigs) == 0 {
		PrintInfo("No repositories with custom configuration found")
		return nil
	}

	FormatHeader("Repositories with Custom Configuration")

	for _, configPath := range repoConfigs {
		// Try to determine repository path from config location
		repoPath := "unknown"
		if strings.Contains(configPath, ".git") {
			repoPath = filepath.Dir(filepath.Dir(configPath))
		}

		fmt.Printf("  • %s\n", highlightColor.Sprint(repoPath))
		PrintDim(fmt.Sprintf("    Config: %s", configPath))

		// Load and show a summary of the config
		cfg := config.DefaultConfig()
		data, err := os.ReadFile(configPath)
		if err == nil {
			if err := config.UnmarshalConfig(data, cfg); err == nil {
				if cfg.Model != "" && cfg.Model != config.DefaultConfig().Model {
					fmt.Printf("    Model: %s\n", cfg.Model)
				}
				if cfg.Commit.Emoji {
					fmt.Printf("    Emoji: enabled\n")
				}
				if cfg.Commit.OneLine {
					fmt.Printf("    Format: one-line\n")
				}
			}
		}
		fmt.Println()
	}

	return nil
}
