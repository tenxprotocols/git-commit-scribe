package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/tenxprotocols/git-commit-scribe/internal/ai"
	"github.com/tenxprotocols/git-commit-scribe/internal/cache"
	"github.com/tenxprotocols/git-commit-scribe/internal/config"
	"github.com/tenxprotocols/git-commit-scribe/internal/conventional"
	"github.com/tenxprotocols/git-commit-scribe/internal/git"
)

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

	// Prompt customization
	Prompt     string `help:"Custom prompt template (file path or inline text)" type:"string"`
	PromptFile string `help:"Path to custom prompt template file" type:"path"`

	// File handling
	MaxFiles         int  `help:"Maximum number of files to analyze" default:"50"`
	IgnoreGenerated  bool `help:"Ignore auto-generated files" default:"true"`
	IgnoreWhitespace bool `help:"Ignore whitespace-only changes" default:"true"`

	// Actions
	DryRun bool `short:"d" help:"Generate message without creating commit"`
	Push   bool `short:"p" help:"Push changes to remote after commit"`
}

// generateCommitMessage generates a commit message using AI
func (c *CommitCmd) generateCommitMessage(cfg *config.Config, diff string, additionalContext string) (*ai.CommitResult, error) {
	// Resolve provider, API key, and model
	resolvedProvider, apiKey, model, err := cfg.ResolveProvider()
	if err != nil {
		return nil, err
	}

	var provider ai.Provider
	switch resolvedProvider {
	case "anthropic":
		provider, err = ai.NewAnthropicProvider(apiKey, model)
	case "openrouter":
		provider, err = ai.NewOpenRouterProvider(apiKey, model)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", resolvedProvider)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create AI provider: %w", err)
	}

	// Determine which custom prompt to use
	customPrompt := c.PromptFile
	if c.Prompt != "" {
		customPrompt = c.Prompt
	}

	// Show progress spinner for AI generation
	spinner := NewSpinner(fmt.Sprintf("Generating commit message using %s/%s...", resolvedProvider, model))
	spinner.Start()

	genOpts := ai.GenerateOptions{
		Diff:              diff,
		Type:              c.Type,
		Scope:             c.Scope,
		Breaking:          c.Breaking,
		OneLine:           cfg.Commit.OneLine,
		MaxLength:         cfg.Commit.DescriptionLength,
		AvailableTypes:    cfg.Types,
		Temperature:       cfg.AI.Temperature,
		CustomPrompt:      customPrompt,
		AdditionalContext: additionalContext,
	}

	// Slightly higher temperature for variation when regenerating with context
	if additionalContext != "" {
		genOpts.Temperature += 0.1
	}

	aiCtx, cancel := context.WithTimeout(context.Background(), cfg.AI.Timeout)
	defer cancel()

	result, err := provider.GenerateCommitMessage(aiCtx, genOpts)
	if err != nil {
		spinner.Error("Failed to generate commit message")
		return nil, fmt.Errorf("failed to generate commit message: %w", err)
	}

	spinner.Success("Commit message generated")
	return result, nil
}

// loadAndOverrideConfig loads configuration and applies CLI flag overrides
func (c *CommitCmd) loadAndOverrideConfig() (*config.Config, error) {
	loader := config.NewLoader(CLI.ConfigDir, CLI.ConfigFile)
	cfg, err := loader.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Override config with CLI flags (only override provider if explicitly set,
	// not when it's the default "auto" which Kong always populates)
	if CLI.Provider != "" && CLI.Provider != "auto" {
		cfg.Provider = CLI.Provider
	}
	if CLI.Model != "" {
		cfg.Model = CLI.Model
	}
	if CLI.APIKey != "" {
		cfg.APIKey = CLI.APIKey
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

	return cfg, nil
}

// getStagedDiffWithProgress gets the staged diff and shows progress
func (c *CommitCmd) getStagedDiffWithProgress() (string, error) {
	diffOpts := git.DiffOptions{
		MaxFiles:         c.MaxFiles,
		IgnoreGenerated:  c.IgnoreGenerated,
		IgnoreWhitespace: c.IgnoreWhitespace,
	}

	var spinner *Spinner
	if !CLI.Verbose {
		spinner = NewSpinner("Getting staged changes...")
		spinner.Start()
	} else {
		PrintInfo("Getting staged changes...")
	}

	diff, err := git.GetStagedDiff(diffOpts)
	if err != nil {
		if spinner != nil {
			spinner.Error("Failed to get staged diff")
		}
		return "", fmt.Errorf("failed to get staged diff: %w", err)
	}

	if spinner != nil {
		files, _ := git.GetStagedFiles()
		spinner.Success(fmt.Sprintf("Analyzing %d staged file(s)", len(files)))
	} else {
		files, _ := git.GetStagedFiles()
		PrintSuccess(fmt.Sprintf("Analyzing %d staged file(s)", len(files)))
	}

	return diff, nil
}

// initializeCache creates and initializes the cache if enabled
func (c *CommitCmd) initializeCache(cfg *config.Config, loader *config.Loader) cache.Cache {
	if !cfg.Cache.Enabled {
		return nil
	}

	cacheDir := config.GetCacheDir(loader.GetConfigDir())

	// Create memory cache
	memCache := cache.NewMemoryCache(cfg.Cache.MaxMemoryMB, time.Duration(cfg.Cache.TTL)*time.Second)

	// Create disk cache
	diskCache, err := cache.NewDiskCache(cacheDir, cfg.Cache.MaxDiskMB, time.Duration(cfg.Cache.TTL)*time.Second)
	if err != nil {
		if CLI.Verbose {
			PrintWarning(fmt.Sprintf("Failed to initialize disk cache: %v", err))
		}
		// Fall back to memory-only cache
		return memCache
	}

	// Use multi-level cache
	if CLI.Verbose {
		PrintInfo("Cache initialized")
	}
	return cache.NewMultiCache(memCache, diskCache)
}

// getOrGenerateCommitMessage retrieves from cache or generates a new commit message
func (c *CommitCmd) getOrGenerateCommitMessage(cfg *config.Config, diff string, cacheInstance cache.Cache, cacheKey string) (*ai.CommitResult, error) {
	// Try to get from cache
	if cfg.Cache.Enabled && cacheInstance != nil {
		cacheCtx := context.Background()
		if cached, found, err := cacheInstance.Get(cacheCtx, cacheKey); err == nil && found {
			if CLI.Verbose {
				PrintSuccess("Using cached commit message")
			}

			// Deserialize cached result
			var cachedResult ai.CommitResult
			if err := json.Unmarshal([]byte(cached), &cachedResult); err == nil {
				return &cachedResult, nil
			} else if CLI.Verbose {
				PrintWarning(fmt.Sprintf("Failed to deserialize cached result: %v", err))
			}
		}
	}

	// If not in cache, generate with AI
	if CLI.Verbose {
		if c.PromptFile != "" {
			PrintInfo(fmt.Sprintf("Using custom prompt from file: %s", c.PromptFile))
		} else if c.Prompt != "" {
			PrintInfo("Using custom inline prompt")
		}
	}

	result, err := c.generateCommitMessage(cfg, diff, "")
	if err != nil {
		return nil, err
	}

	// Store in cache
	if cfg.Cache.Enabled && cacheInstance != nil {
		cacheCtx := context.Background()
		if serialized, err := json.Marshal(result); err == nil {
			if err := cacheInstance.Set(cacheCtx, cacheKey, string(serialized)); err != nil && CLI.Verbose {
				PrintWarning(fmt.Sprintf("Failed to cache result: %v", err))
			}
		}
	}

	return result, nil
}

// formatCommitMessage formats the AI result into a commit message
func (c *CommitCmd) formatCommitMessage(result *ai.CommitResult, cfg *config.Config) string {
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
		return ""
	}

	// Format the final message
	if cfg.Commit.OneLine {
		return msg.FormatOneLine()
	}
	return msg.Format()
}

// handleUserInteraction handles the interactive prompt loop for commit confirmation
func (c *CommitCmd) handleUserInteraction(commitMsg string, cfg *config.Config, diff string) (string, bool, error) {
	for {
		fmt.Println()
		choice, err := PromptChoice("What would you like to do?", []string{
			"Create commit with this message",
			"Edit message before committing",
			"Regenerate message",
			"Cancel",
		}, 1) // Default to option 1
		if err != nil {
			return "", false, fmt.Errorf("failed to read input: %w", err)
		}

		switch choice {
		case 1:
			// Continue with commit
			return commitMsg, true, nil
		case 2:
			// Edit message
			edited, err := EditText(commitMsg)
			if err != nil {
				return "", false, fmt.Errorf("failed to edit message: %w", err)
			}
			if edited == "" {
				PrintWarning("Empty commit message, cancelling")
				return "", false, nil
			}
			commitMsg = edited

			// Show edited message
			FormatSubHeader("Edited Commit Message")
			PrintSeparator()
			PrintCommitMessage(commitMsg)
			PrintSeparator()
			// Loop back to show choices again
		case 3:
			// Regenerate with optional additional context
			fmt.Println()
			additionalContext, err := PromptMultilineText("Provide additional context for regeneration (optional):")
			if err != nil {
				return "", false, fmt.Errorf("failed to read additional context: %w", err)
			}

			result, err := c.generateCommitMessage(cfg, diff, additionalContext)
			if err != nil {
				return "", false, err
			}

			commitMsg = c.formatCommitMessage(result, cfg)
			if commitMsg == "" {
				return "", false, fmt.Errorf("failed to format commit message")
			}

			// Display the regenerated message
			FormatSubHeader("Regenerated Commit Message")
			PrintSeparator()
			PrintCommitMessage(commitMsg)
			PrintSeparator()
			// Loop back to show choices again
		case 0, 4:
			// Cancel
			PrintInfo("Commit cancelled")
			return "", false, nil
		}
	}
}

// Run executes the commit command
func (c *CommitCmd) Run(_ *Context) error {
	// Validate git repository and staged changes first (fail early)
	if err := git.ValidateRepository(); err != nil {
		return err
	}

	// Load configuration
	cfg, err := c.loadAndOverrideConfig()
	if err != nil {
		return err
	}

	// Get staged diff
	diff, err := c.getStagedDiffWithProgress()
	if err != nil {
		return err
	}

	// Initialize cache if enabled
	loader := config.NewLoader(CLI.ConfigDir, CLI.ConfigFile)
	cacheInstance := c.initializeCache(cfg, loader)

	// Generate cache key from diff and options
	cacheKey := cache.GenerateKey(
		diff,
		c.Type,
		c.Scope,
		fmt.Sprintf("%v", c.Breaking),
		fmt.Sprintf("%v", cfg.Commit.OneLine),
		fmt.Sprintf("%d", cfg.Commit.DescriptionLength),
		cfg.Model,
	)

	// Get or generate commit message
	result, err := c.getOrGenerateCommitMessage(cfg, diff, cacheInstance, cacheKey)
	if err != nil {
		return err
	}

	// Format the commit message
	commitMsg := c.formatCommitMessage(result, cfg)
	if commitMsg == "" {
		return fmt.Errorf("failed to format commit message")
	}

	// Display the generated message with syntax highlighting
	FormatSubHeader("Generated Commit Message")
	PrintSeparator()
	PrintCommitMessage(commitMsg)
	PrintSeparator()

	// Dry run mode - just show the message
	if c.DryRun {
		PrintInfo("\nDry run mode - no commit created")
		return nil
	}

	// Confirm with user unless --yes flag is set
	if cfg.Commit.Confirm {
		var proceed bool
		commitMsg, proceed, err = c.handleUserInteraction(commitMsg, cfg, diff)
		if err != nil {
			return err
		}
		if !proceed {
			return nil
		}
	}

	// Create the commit
	spinner := NewSpinner("Creating commit...")
	spinner.Start()

	if err := git.CreateCommit(commitMsg); err != nil {
		spinner.Error("Failed to create commit")
		return fmt.Errorf("failed to create commit: %w", err)
	}

	spinner.Success("Commit created successfully")

	// Push if requested
	if cfg.Commit.AutoPush {
		if !git.HasRemote() {
			PrintWarning("No remote configured, skipping push")
			return nil
		}

		spinner = NewSpinner("Pushing changes...")
		spinner.Start()

		if err := git.PushChanges(); err != nil {
			spinner.Error("Failed to push changes")
			return fmt.Errorf("failed to push changes: %w", err)
		}

		spinner.Success("Changes pushed successfully")
	}

	return nil
}
