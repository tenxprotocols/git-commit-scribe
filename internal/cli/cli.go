package cli

// Context holds runtime context for commands
type Context struct {
	Debug bool
}

// CLI is the main command structure
var CLI struct {
	// Global flags
	ConfigFile string `help:"Path to config file" type:"path" env:"GSCRIBE_CONFIG"`
	ConfigDir  string `help:"Path to config directory" type:"path" env:"GSCRIBE_CONFIG_DIR" default:"~/.config/gscribe"`
	Provider   string `help:"AI provider (currently only openrouter)" default:"openrouter" enum:"openrouter" env:"GSCRIBE_PROVIDER"`
	Model      string `help:"Model to use" env:"GSCRIBE_MODEL"`
	APIKey     string `help:"API key for AI provider" env:"GSCRIBE_API_KEY"`
	Verbose    bool   `short:"v" help:"Enable verbose logging" env:"GSCRIBE_VERBOSE"`
	NoCache    bool   `help:"Disable caching" env:"GSCRIBE_NO_CACHE"`

	// Commands
	Commit CommitCmd `cmd:"" help:"Generate and create a git commit" default:"1"`
	Config ConfigCmd `cmd:"" help:"Manage configuration"`
	Cache  CacheCmd  `cmd:"" help:"Manage cache"`
}
