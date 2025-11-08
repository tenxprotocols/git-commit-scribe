package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/tenxprotocols/git-commit-scribe/internal/ai"
	"github.com/tenxprotocols/git-commit-scribe/internal/config"
	"github.com/tenxprotocols/git-commit-scribe/internal/conventional"
	"github.com/tenxprotocols/git-commit-scribe/internal/git"
)

// Context holds runtime context for commands
type Context struct {
	Debug bool
}

// CLI is the main command structure
var CLI struct {
	// Global flags
	ConfigFile string `help:"Path to config file" type:"path" env:"GSCRIBE_CONFIG"`
	ConfigDir  string `help:"Path to config directory" type:"path" env:"GSCRIBE_CONFIG_DIR" default:"~/.config/gscribe"`
	Provider  string `help:"AI provider (currently only openrouter)" default:"openrouter" enum:"openrouter" env:"GSCRIBE_PROVIDER"`
	Model     string `help:"Model to use" env:"GSCRIBE_MODEL"`
	Verbose   bool   `short:"v" help:"Enable verbose logging" env:"GSCRIBE_VERBOSE"`
	NoCache   bool   `help:"Disable caching" env:"GSCRIBE_NO_CACHE"`

	// Commands
	Commit CommitCmd `cmd:"" help:"Generate and create a git commit" default:"1"`
	Config ConfigCmd `cmd:"" help:"Manage configuration"`
	Cache  CacheCmd  `cmd:"" help:"Manage cache"`
}

// CommitCmd handles commit generation
type CommitCmd struct {
	// Confirmation
	Yes bool `short:"y" help:"Skip confirmations"`

	// Commit structure
	Type        string `short:"t" help:"Commit type (feat, fix, etc.)"`
	Scope       string `short:"s" help:"Commit scope"`
	NoScope     bool   `help:"Explicitly exclude scope"`
	Breaking    bool   `short:"b" help:"Mark as breaking change"`
	Emoji       bool   `short:"e" help:"Include emoji in commit message"`
	OneLine     bool   `short:"o" help:"Generate a one-line commit message"`
	Description int    `help:"Max description length" default:"72"`

	// File handling
	MaxFiles         int  `help:"Maximum number of files to analyze" default:"50"`
	IgnoreGenerated  bool `help:"Ignore auto-generated files" default:"true"`
	IgnoreWhitespace bool `help:"Ignore whitespace-only changes" default:"true"`

	// Actions
	DryRun bool `short:"d" help:"Generate message without creating commit"`
	Push   bool `short:"p" help:"Push changes to remote after commit"`
}

// Run executes the commit command
func (c *CommitCmd) Run(ctx *Context) error {
	// Load configuration
	loader := config.NewLoader(CLI.ConfigDir, CLI.ConfigFile)
	cfg, err := loader.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Override config with CLI flags
	if CLI.Provider != "" {
		cfg.Provider = CLI.Provider
	}
	if CLI.Model != "" {
		cfg.Model = CLI.Model
	}
	if CLI.NoCache {
		cfg.Cache.Enabled = false
	}

	// Override commit settings with command flags
	if c.Emoji {
		cfg.Commit.Emoji = true
	}
	if c.OneLine {
		cfg.Commit.OneLine = true
	}
	if c.Description > 0 {
		cfg.Commit.DescriptionLength = c.Description
	}
	if c.Yes {
		cfg.Commit.Confirm = false
	}
	if c.Push {
		cfg.Commit.AutoPush = true
	}

	// Get staged diff
	diffOpts := git.DiffOptions{
		MaxFiles:         c.MaxFiles,
		IgnoreGenerated:  c.IgnoreGenerated,
		IgnoreWhitespace: c.IgnoreWhitespace,
	}

	if CLI.Verbose {
		fmt.Println("Getting staged changes...")
	}

	diff, err := git.GetStagedDiff(diffOpts)
	if err != nil {
		return fmt.Errorf("failed to get staged diff: %w", err)
	}

	if CLI.Verbose {
		files, _ := git.GetStagedFiles()
		fmt.Printf("Analyzing %d staged file(s)...\n", len(files))
	}

	// Initialize AI provider
	var provider ai.Provider
	switch cfg.Provider {
	case "openrouter":
		if cfg.APIKey == "" {
			return fmt.Errorf("API key not configured. Run 'gscribe config init' or set GSCRIBE_API_KEY environment variable")
		}
		provider, err = ai.NewOpenRouterProvider(cfg.APIKey, cfg.Model)
		if err != nil {
			return fmt.Errorf("failed to create AI provider: %w", err)
		}
	default:
		return fmt.Errorf("unsupported provider: %s", cfg.Provider)
	}

	// Generate commit message
	if CLI.Verbose {
		fmt.Printf("Generating commit message using %s...\n", cfg.Model)
	}

	genOpts := ai.GenerateOptions{
		Diff:           diff,
		Type:           c.Type,
		Scope:          c.Scope,
		Breaking:       c.Breaking,
		OneLine:        cfg.Commit.OneLine,
		MaxLength:      cfg.Commit.DescriptionLength,
		AvailableTypes: cfg.Types,
		Temperature:    cfg.AI.Temperature,
	}

	aiCtx, cancel := context.WithTimeout(context.Background(), cfg.AI.Timeout)
	defer cancel()

	result, err := provider.GenerateCommitMessage(aiCtx, genOpts)
	if err != nil {
		return fmt.Errorf("failed to generate commit message: %w", err)
	}

	// Format the commit message
	msg := &conventional.CommitMessage{
		Type:        result.Type,
		Scope:       result.Scope,
		Breaking:    result.Breaking || c.Breaking,
		Description: result.Description,
		Body:        result.Body,
		Emoji:       cfg.Commit.Emoji,
	}

	// Ensure description length limit
	msg.TruncateDescription(cfg.Commit.DescriptionLength)

	// Validate the message
	if err := msg.Validate(cfg.Types); err != nil {
		return fmt.Errorf("invalid commit message: %w", err)
	}

	// Format the final message
	var commitMsg string
	if cfg.Commit.OneLine {
		commitMsg = msg.FormatOneLine()
	} else {
		commitMsg = msg.Format()
	}

	// Display the generated message
	fmt.Println("\nGenerated commit message:")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println(commitMsg)
	fmt.Println("─────────────────────────────────────────")

	// Dry run mode - just show the message
	if c.DryRun {
		fmt.Println("\nDry run mode - no commit created")
		return nil
	}

	// Confirm with user unless --yes flag is set
	if cfg.Commit.Confirm {
		fmt.Print("\nCreate commit with this message? [Y/n]: ")
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		response = strings.ToLower(strings.TrimSpace(response))
		if response != "" && response != "y" && response != "yes" {
			fmt.Println("Commit cancelled")
			return nil
		}
	}

	// Create the commit
	if CLI.Verbose {
		fmt.Println("Creating commit...")
	}

	if err := git.CreateCommit(commitMsg); err != nil {
		return fmt.Errorf("failed to create commit: %w", err)
	}

	fmt.Println("✓ Commit created successfully")

	// Push if requested
	if cfg.Commit.AutoPush {
		if !git.HasRemote() {
			fmt.Println("⚠ No remote configured, skipping push")
			return nil
		}

		if CLI.Verbose {
			fmt.Println("Pushing changes...")
		}

		if err := git.PushChanges(); err != nil {
			return fmt.Errorf("failed to push changes: %w", err)
		}

		fmt.Println("✓ Changes pushed successfully")
	}

	return nil
}

// ConfigCmd manages configuration
type ConfigCmd struct {
	Init   ConfigInitCmd   `cmd:"" help:"Initialize configuration"`
	Show   ConfigShowCmd   `cmd:"" help:"Show current configuration"`
	Set    ConfigSetCmd    `cmd:"" help:"Set a configuration value"`
	Unset  ConfigUnsetCmd  `cmd:"" help:"Unset a configuration value"`
	Repos  ConfigReposCmd  `cmd:"" help:"Manage per-repository configuration"`
}

// ConfigInitCmd initializes configuration
type ConfigInitCmd struct {
	Force bool `short:"f" help:"Overwrite existing configuration"`
}

// Run executes the config init command
func (c *ConfigInitCmd) Run(ctx *Context) error {
	if ctx.Debug {
		fmt.Println("Running config init command...")
	}
	// TODO: Implement config init logic
	return fmt.Errorf("config init command not yet implemented")
}

// ConfigShowCmd shows current configuration
type ConfigShowCmd struct {
	Global bool `short:"g" help:"Show global configuration only"`
}

// Run executes the config show command
func (c *ConfigShowCmd) Run(ctx *Context) error {
	if ctx.Debug {
		fmt.Println("Running config show command...")
	}
	// TODO: Implement config show logic
	return fmt.Errorf("config show command not yet implemented")
}

// ConfigSetCmd sets a configuration value
type ConfigSetCmd struct {
	Key   string `arg:"" help:"Configuration key"`
	Value string `arg:"" help:"Configuration value"`
	Repo  bool   `short:"r" help:"Set for current repository only"`
}

// Run executes the config set command
func (c *ConfigSetCmd) Run(ctx *Context) error {
	if ctx.Debug {
		fmt.Printf("Setting %s = %s\n", c.Key, c.Value)
	}
	// TODO: Implement config set logic
	return fmt.Errorf("config set command not yet implemented")
}

// ConfigUnsetCmd unsets a configuration value
type ConfigUnsetCmd struct {
	Key  string `arg:"" help:"Configuration key"`
	Repo bool   `short:"r" help:"Unset for current repository only"`
}

// Run executes the config unset command
func (c *ConfigUnsetCmd) Run(ctx *Context) error {
	if ctx.Debug {
		fmt.Printf("Unsetting %s\n", c.Key)
	}
	// TODO: Implement config unset logic
	return fmt.Errorf("config unset command not yet implemented")
}

// ConfigReposCmd manages per-repository configuration
type ConfigReposCmd struct {
	List ConfigReposListCmd `cmd:"" help:"List repositories with custom configuration"`
}

// ConfigReposListCmd lists repositories with custom configuration
type ConfigReposListCmd struct{}

// Run executes the config repos list command
func (c *ConfigReposListCmd) Run(ctx *Context) error {
	if ctx.Debug {
		fmt.Println("Listing repositories with custom configuration...")
	}
	// TODO: Implement config repos list logic
	return fmt.Errorf("config repos list command not yet implemented")
}

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
	if ctx.Debug {
		fmt.Println("Clearing cache...")
	}
	// TODO: Implement cache clear logic
	return fmt.Errorf("cache clear command not yet implemented")
}

// CacheStatsCmd shows cache statistics
type CacheStatsCmd struct{}

// Run executes the cache stats command
func (c *CacheStatsCmd) Run(ctx *Context) error {
	if ctx.Debug {
		fmt.Println("Showing cache statistics...")
	}
	// TODO: Implement cache stats logic
	return fmt.Errorf("cache stats command not yet implemented")
}
