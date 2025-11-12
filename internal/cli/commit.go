package cli

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
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

	// Initialize cache if enabled
	var cacheInstance cache.Cache
	if cfg.Cache.Enabled {
		cacheDir := config.GetCacheDir(loader.GetConfigDir())
		
		// Create memory cache
		memCache := cache.NewMemoryCache(cfg.Cache.MaxMemoryMB, time.Duration(cfg.Cache.TTL)*time.Second)
		
		// Create disk cache
		diskCache, err := cache.NewDiskCache(cacheDir, cfg.Cache.MaxDiskMB, time.Duration(cfg.Cache.TTL)*time.Second)
		if err != nil {
			if CLI.Verbose {
				fmt.Printf("Warning: Failed to initialize disk cache: %v\n", err)
			}
			// Fall back to memory-only cache
			cacheInstance = memCache
		} else {
			// Use multi-level cache
			cacheInstance = cache.NewMultiCache(memCache, diskCache)
		}
		
		if CLI.Verbose {
			fmt.Println("Cache initialized")
		}
	}

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

	// Try to get from cache
	var result *ai.CommitResult
	if cfg.Cache.Enabled && cacheInstance != nil {
		cacheCtx := context.Background()
		if cached, found, err := cacheInstance.Get(cacheCtx, cacheKey); err == nil && found {
			if CLI.Verbose {
				fmt.Println("✓ Using cached commit message")
			}
			
			// Deserialize cached result
			var cachedResult ai.CommitResult
			if err := json.Unmarshal([]byte(cached), &cachedResult); err == nil {
				result = &cachedResult
			} else if CLI.Verbose {
				fmt.Printf("Warning: Failed to deserialize cached result: %v\n", err)
			}
		}
	}

	// If not in cache, generate with AI
	if result == nil {
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

		result, err = provider.GenerateCommitMessage(aiCtx, genOpts)
		if err != nil {
			return fmt.Errorf("failed to generate commit message: %w", err)
		}

		// Store in cache
		if cfg.Cache.Enabled && cacheInstance != nil {
			cacheCtx := context.Background()
			if serialized, err := json.Marshal(result); err == nil {
				if err := cacheInstance.Set(cacheCtx, cacheKey, string(serialized)); err != nil && CLI.Verbose {
					fmt.Printf("Warning: Failed to cache result: %v\n", err)
				}
			}
		}
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
