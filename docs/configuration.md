# Configuration

## Configuration Files

gscribe uses a hierarchical configuration system:

1. **Global config**: `~/.config/gscribe/config.yaml`
2. **Repository config**: `.git/gscribe-config.yaml` (optional, overrides global)

## Initial Setup

Initialize your configuration:

```bash
gscribe config init
```

This will prompt you for:
- OpenRouter API key
- Default model (recommended: `anthropic/claude-3.5-sonnet`)

## Configuration Options

### API Settings

```yaml
provider: openrouter
model: anthropic/claude-3.5-sonnet
api_key: your-api-key-here
```

### AI Behavior

```yaml
ai:
  temperature: 0.7          # Creativity level (0.0-1.0)
  timeout: 30s              # Request timeout
  max_retries: 3            # Retry attempts on failure
  rate_limit_wait: 5s       # Wait time between retries
```

### Commit Preferences

```yaml
commit:
  emoji: false              # Include emoji in messages
  one_line: false           # Generate only one-line messages
  description_length: 72    # Max description length
  max_files: 50             # Max files to analyze
  ignore_generated: true    # Skip auto-generated files
  ignore_whitespace: true   # Skip whitespace-only changes
  auto_push: false          # Automatically push after commit
  confirm: true             # Ask for confirmation before committing
```

### Caching

```yaml
cache:
  enabled: true
  ttl: 86400                # Cache duration in seconds (24 hours)
  max_memory_mb: 100        # Memory cache limit
  max_disk_mb: 1000         # Disk cache limit
```

### Commit Types

Customize available commit types:

```yaml
types:
  feat: "A new feature"
  fix: "A bug fix"
  docs: "Documentation only changes"
  style: "Code style changes (formatting, etc)"
  refactor: "Code refactoring"
  perf: "Performance improvements"
  test: "Adding or updating tests"
  build: "Build system changes"
  ci: "CI configuration changes"
  chore: "Other changes"
  revert: "Revert a previous commit"
```

## Environment Variables

Override configuration with environment variables:

```bash
export GSCRIBE_API_KEY="your-api-key"
export GSCRIBE_MODEL="anthropic/claude-3.5-sonnet"
export GSCRIBE_PROVIDER="openrouter"
export GSCRIBE_VERBOSE="true"
export GSCRIBE_NO_CACHE="true"
```

## Per-Repository Configuration

Create a repository-specific config:

```bash
# In your repository
mkdir -p .git
cat > .git/gscribe-config.yaml <<EOF
commit:
  emoji: true
  one_line: true
model: anthropic/claude-3.5-sonnet
EOF
```

This will override your global settings for this repository only.

## Managing Configuration

```bash
# View current configuration
gscribe config show

# Set a value
gscribe config set commit.emoji true

# Set for repository only
gscribe config set commit.emoji true --repo

# Remove a value
gscribe config unset commit.emoji
```

## API Key Security

- API keys are stored with `0600` permissions (user read/write only)
- Never commit config files containing API keys to version control
- Use environment variables for CI/CD environments
