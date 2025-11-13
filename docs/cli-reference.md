# CLI Reference

Complete reference for all gscribe commands and options.

## Global Flags

These flags are available for all commands:

| Flag | Short | Environment Variable | Default | Description |
|------|-------|---------------------|---------|-------------|
| `--config-file` | | `GSCRIBE_CONFIG` | | Path to config file |
| `--config-dir` | | `GSCRIBE_CONFIG_DIR` | `~/.config/gscribe` | Path to config directory |
| `--provider` | | `GSCRIBE_PROVIDER` | `openrouter` | AI provider |
| `--model` | | `GSCRIBE_MODEL` | | Model to use |
| `--verbose` | `-v` | `GSCRIBE_VERBOSE` | `false` | Enable verbose logging |
| `--no-cache` | | `GSCRIBE_NO_CACHE` | `false` | Disable caching |
| `--help` | `-h` | | | Show help |

## Commands

### gscribe commit

Generate and create a git commit message.

**Usage:**
```bash
gscribe [commit] [flags]
```

Note: `commit` is the default command and can be omitted.

**Flags:**

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--yes` | `-y` | `false` | Skip confirmations |
| `--type` | `-t` | | Commit type (feat, fix, etc.) |
| `--scope` | `-s` | | Commit scope |
| `--no-scope` | | `false` | Explicitly exclude scope |
| `--breaking` | `-b` | `false` | Mark as breaking change |
| `--emoji` | `-e` | `false` | Include emoji in commit message |
| `--one-line` | `-o` | `false` | Generate one-line commit message |
| `--description` | | `72` | Max description length |
| `--max-files` | | `50` | Maximum number of files to analyze |
| `--ignore-generated` | | `true` | Ignore auto-generated files |
| `--ignore-whitespace` | | `true` | Ignore whitespace-only changes |
| `--dry-run` | `-d` | `false` | Generate message without creating commit |
| `--push` | `-p` | `false` | Push changes to remote after commit |

**Examples:**

```bash
# Basic usage
gscribe

# Feature with scope
gscribe -t feat -s api

# Breaking change with emoji
gscribe -b -e

# Quick commit and push
gscribe -yp

# Dry run to preview
gscribe -d
```

### gscribe config init

Initialize configuration with interactive prompts.

**Usage:**
```bash
gscribe config init [flags]
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--force` | `-f` | Overwrite existing configuration |

**Examples:**

```bash
# Initialize config
gscribe config init

# Force reinitialize
gscribe config init -f
```

### gscribe config show

Show current configuration values.

**Usage:**
```bash
gscribe config show [flags]
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--global` | `-g` | Show global configuration only |

**Examples:**

```bash
# Show merged config (global + repo)
gscribe config show

# Show only global config
gscribe config show -g
```

### gscribe config set

Set a configuration value.

**Usage:**
```bash
gscribe config set <key> <value> [flags]
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--repo` | `-r` | Set for current repository only |

**Examples:**

```bash
# Set global value
gscribe config set commit.emoji true

# Set for repository only
gscribe config set model anthropic/claude-3.5-sonnet --repo

# Set nested values
gscribe config set ai.temperature 0.8
```

**Available Keys:**

- `provider` - AI provider name
- `model` - Model identifier
- `api_key` - API key (use environment variable instead)
- `commit.emoji` - Enable emoji (true/false)
- `commit.one_line` - One-line commits (true/false)
- `commit.description_length` - Max length (number)
- `commit.max_files` - Max files (number)
- `commit.ignore_generated` - Ignore generated files (true/false)
- `commit.ignore_whitespace` - Ignore whitespace (true/false)
- `commit.auto_push` - Auto push (true/false)
- `commit.confirm` - Require confirmation (true/false)
- `ai.temperature` - Temperature (0.0-1.0)
- `ai.timeout` - Timeout duration
- `ai.max_retries` - Max retries (number)
- `cache.enabled` - Enable cache (true/false)
- `cache.ttl` - Cache TTL in seconds (number)

### gscribe config unset

Remove a configuration value.

**Usage:**
```bash
gscribe config unset <key> [flags]
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--repo` | `-r` | Unset for current repository only |

**Examples:**

```bash
# Remove global value
gscribe config unset commit.emoji

# Remove repository value
gscribe config unset model --repo
```

### gscribe config repos list

List repositories with custom configuration.

**Usage:**
```bash
gscribe config repos list
```

**Examples:**

```bash
gscribe config repos list
```

### gscribe cache clear

Clear the cache.

**Usage:**
```bash
gscribe cache clear [flags]
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--all` | `-a` | Clear both memory and disk cache |

**Examples:**

```bash
# Clear cache
gscribe cache clear

# Clear all caches
gscribe cache clear -a
```

### gscribe cache stats

Show cache statistics.

**Usage:**
```bash
gscribe cache stats
```

**Examples:**

```bash
gscribe cache stats
```

## Exit Codes

| Code | Description |
|------|-------------|
| 0 | Success |
| 1 | General error |
| 2 | Configuration error |
| 3 | Git error |
| 4 | API error |

## Environment Variables

All global flags can be set via environment variables with the `GSCRIBE_` prefix:

```bash
export GSCRIBE_API_KEY="sk-..."
export GSCRIBE_MODEL="anthropic/claude-3.5-sonnet"
export GSCRIBE_PROVIDER="openrouter"
export GSCRIBE_VERBOSE="true"
export GSCRIBE_NO_CACHE="true"
export GSCRIBE_CONFIG_DIR="$HOME/.config/gscribe"
```

## Configuration File Locations

1. **Global**: `~/.config/gscribe/config.yaml`
2. **Repository**: `.git/gscribe-config.yaml`

Repository configuration overrides global configuration.

## Commit Types

Default conventional commit types:

| Type | Description |
|------|-------------|
| `feat` | A new feature |
| `fix` | A bug fix |
| `docs` | Documentation only changes |
| `style` | Code style changes (formatting, etc) |
| `refactor` | Code refactoring |
| `perf` | Performance improvements |
| `test` | Adding or updating tests |
| `build` | Build system changes |
| `ci` | CI configuration changes |
| `chore` | Other changes |
| `revert` | Revert a previous commit |

## Examples

### Common Workflows

```bash
# Quick commit workflow
git add .
gscribe -y

# Feature development
git add src/auth/
gscribe -t feat -s auth -e

# Bug fix with push
git add src/utils/validator.js
gscribe -t fix -p

# Documentation update
git add README.md docs/
gscribe -t docs -o -y

# Breaking change
git add .
gscribe -b -t feat -s api
```

### Configuration Examples

```bash
# Set up for emoji commits
gscribe config set commit.emoji true
gscribe config set commit.confirm false

# Increase file limit for monorepos
gscribe config set commit.max_files 100

# Fast commits in a specific repo
cd my-project
gscribe config set commit.confirm false --repo
gscribe config set commit.one_line true --repo
```

### Environment-Based Configuration

```bash
# Development environment
export GSCRIBE_VERBOSE=true
export GSCRIBE_NO_CACHE=true

# CI environment
export GSCRIBE_API_KEY="${SECRET_API_KEY}"
export GSCRIBE_MODEL="anthropic/claude-3.5-sonnet"
gscribe -y -t ci
```

## Advanced Usage Patterns

### Chain Multiple Operations

```bash
# Commit and immediately create PR
gscribe -yp && gh pr create --fill

# Commit, tag, and push
gscribe -y && git tag v1.0.0 && git push --tags

# Auto-format, commit, and push
go fmt ./... && gscribe -yp
```

### Conditional Commits

```bash
# Only commit if there are changes
if [[ -n $(git status -s) ]]; then
  gscribe -y
fi

# Commit only specific file types
if git diff --staged --name-only | grep -q '\.go$'; then
  gscribe -t feat -s backend -y
fi
```

### Scripted Workflows

```bash
#!/bin/bash
# auto-commit.sh - Automated commit script

# Check for staged changes
if [[ -z $(git diff --staged) ]]; then
  echo "No staged changes"
  exit 0
fi

# Determine commit type based on files
if git diff --staged --name-only | grep -q '^docs/'; then
  TYPE="docs"
elif git diff --staged --name-only | grep -q '\.test\.'; then
  TYPE="test"
else
  TYPE="feat"
fi

# Generate and commit
gscribe -t "$TYPE" -y

echo "✓ Committed as $TYPE"
```

### Interactive Selection

```bash
# Select files interactively
git add -i

# Generate commit for selection
gscribe -d

# Review and commit
gscribe
```

### Bulk Operations

```bash
# Process multiple directories
for dir in packages/*/; do
  cd "$dir"
  if [[ -n $(git status -s) ]]; then
    git add .
    gscribe -t feat -s "$(basename $dir)" -y
  fi
  cd -
done

# Commit files by pattern
find . -name "*.md" -type f | while read f; do
  git add "$f"
  gscribe -t docs -o -y
done
```

## Flag Combinations Reference

### Common Combinations

| Flags | Use Case | Result |
|-------|----------|--------|
| `-yp` | Quick commit and push | Auto-approve, then push |
| `-ye` | Quick emoji commit | Auto-approve with emoji |
| `-d -v` | Debug/preview | Dry run with verbose output |
| `-t feat -s api` | Scoped feature | Feature commit for API |
| `-b -t feat` | Breaking feature | Breaking change marker |
| `-o -y` | Fast one-liner | One-line auto-approved |
| `-e -t fix -s ui` | UI fix with emoji | Emoji fix for UI scope |

### Advanced Combinations

```bash
# Maximum speed (auto-everything)
gscribe -yop

# Maximum quality (review everything)
gscribe -v

# Feature development
gscribe -t feat -s auth -e -v

# Hotfix
gscribe -t fix --breaking -y

# Documentation sprint
gscribe -t docs -o -ye

# CI/CD commit
gscribe -t ci --no-scope -y --max-files 100
```

## Command Chaining Examples

### Development Workflow

```bash
# Format, test, commit, push
npm run format && npm test && gscribe -yp

# Build, test, commit
go build && go test ./... && gscribe -t build -y

# Lint, fix, commit
npm run lint:fix && gscribe -t style -y
```

### Release Workflow

```bash
# Update version, commit, tag, push
npm version patch && \
gscribe -t chore -s release -y && \
git push --follow-tags

# Changelog, commit, tag
npx conventional-changelog -p angular -i CHANGELOG.md -s && \
git add CHANGELOG.md && \
gscribe -t docs -s changelog -y && \
git tag -a v$(node -p "require('./package.json').version") -m "Release"
```

### Pre-commit Pipeline

```bash
# Full quality pipeline
prettier --write . && \
eslint --fix . && \
npm test && \
gscribe -ye
```

## Shell Integration

### Bash Functions

```bash
# Add to ~/.bashrc

# Smart commit (auto-detect type)
gc() {
  local type=""
  
  # Detect type from staged files
  if git diff --staged --name-only | grep -q '^docs/'; then
    type="docs"
  elif git diff --staged --name-only | grep -q '\.test\.'; then
    type="test"
  elif git diff --staged --name-only | grep -q '\.css$\|\.scss$'; then
    type="style"
  else
    type="feat"
  fi
  
  gscribe -t "$type" -y
}

# Commit with scope
gcs() {
  gscribe -t "${1:-feat}" -s "${2:-}" -y
}

# Quick fix
gfix() {
  git add .
  gscribe -t fix -s "${1:-}" -y
}

# Documentation commit
gdoc() {
  git add docs/ *.md
  gscribe -t docs -o -y
}
```

### Zsh Functions

```zsh
# Add to ~/.zshrc

# Commit with emoji
gce() {
  gscribe -e -y
}

# Scoped commit with completion
gcs() {
  local scope=$1
  shift
  gscribe -t feat -s "$scope" "$@"
}

# Auto-complete scopes
_gscribe_scopes() {
  local -a scopes
  scopes=(api ui auth db cache docs tests utils)
  _describe 'scope' scopes
}
compdef _gscribe_scopes gcs
```

### Fish Shell

```fish
# Add to ~/.config/fish/functions/

# Quick commit
function gc
  git add .
  gscribe -y
end

# Feature commit
function gcf
  gscribe -t feat -s $argv[1] -y
end

# Fix commit
function gcx
  gscribe -t fix -s $argv[1] -y
end
```

## Configuration Override Examples

### Temporary Overrides

```bash
# One-time model change
GSCRIBE_MODEL="openai/gpt-4" gscribe

# One-time verbose
GSCRIBE_VERBOSE=true gscribe

# Disable cache for this commit
GSCRIBE_NO_CACHE=true gscribe

# Multiple overrides
GSCRIBE_MODEL="openai/gpt-3.5-turbo" \
GSCRIBE_VERBOSE=true \
gscribe -d
```

### Environment-Specific

```bash
# Development profile
dev_commit() {
  GSCRIBE_MODEL="openai/gpt-3.5-turbo" \
  GSCRIBE_VERBOSE=true \
  gscribe "$@"
}

# Production profile
prod_commit() {
  GSCRIBE_MODEL="anthropic/claude-3.5-sonnet" \
  GSCRIBE_VERBOSE=false \
  gscribe -y "$@"
}
```

## Error Handling

### Safe Commit Script

```bash
#!/bin/bash
# safe-commit.sh

set -e  # Exit on error

# Function to cleanup on error
cleanup() {
  echo "Error occurred, cleaning up..."
  git reset HEAD
}
trap cleanup ERR

# Verify staged changes
if [[ -z $(git diff --staged) ]]; then
  echo "Error: No staged changes"
  exit 1
fi

# Generate commit
if ! gscribe -d; then
  echo "Error: Failed to generate commit message"
  exit 1
fi

# Commit
gscribe -y

echo "✓ Commit successful"
```

### Retry Logic

```bash
# Retry on failure
commit_with_retry() {
  local max_attempts=3
  local attempt=1
  
  while [ $attempt -le $max_attempts ]; do
    echo "Attempt $attempt of $max_attempts..."
    
    if gscribe "$@"; then
      echo "✓ Success"
      return 0
    fi
    
    echo "✗ Failed, waiting before retry..."
    sleep 5
    ((attempt++))
  done
  
  echo "✗ Failed after $max_attempts attempts"
  return 1
}
```

## Platform-Specific Notes

### macOS

```bash
# Use with pbcopy to copy commit message
gscribe -d | pbcopy
# Now paste into IDE or review

# Integration with Alfred/Raycast
gscribe -d -o | pbcopy
```

### Windows (PowerShell)

```powershell
# Copy to clipboard
gscribe -d | Set-Clipboard

# Git alias
git config --global alias.cm "!gscribe -y"

# Function in profile
function Quick-Commit {
  git add .
  gscribe -y
}
Set-Alias gc Quick-Commit
```

### Linux

```bash
# Use with xclip
gscribe -d | xclip -selection clipboard

# Desktop notification on completion
gscribe -y && notify-send "Commit" "Successfully committed"
```

## Performance Tips

### Speed Optimization

```bash
# Use cache aggressively
gscribe config set cache.enabled true
gscribe config set cache.ttl 604800

# Use faster model
export GSCRIBE_MODEL="openai/gpt-3.5-turbo"

# Limit files
gscribe --max-files 30

# One-line mode
gscribe -o
```

### Quality Optimization

```bash
# Use better model
export GSCRIBE_MODEL="anthropic/claude-3.5-sonnet"

# Allow more time
gscribe config set ai.timeout 60s

# Verbose for debugging
gscribe -v

# Always review
gscribe config set commit.confirm true
```

## Debugging Commands

### Diagnostic Commands

```bash
# Check version
gscribe --version

# Verify configuration
gscribe config show

# Test with verbose
gscribe -d -v

# Check cache stats
gscribe cache stats

# Clear cache
gscribe cache clear --all

# Test API connectivity
echo "test" > /tmp/test.txt
cd /tmp && git init && git add test.txt
gscribe -d -v
```

### Environment Check

```bash
# Complete diagnostic
echo "=== gscribe Diagnostics ==="
echo "Version: $(gscribe --version)"
echo "Go Version: $(go version)"
echo "Git Version: $(git --version)"
echo ""
echo "=== Environment Variables ==="
env | grep GSCRIBE
echo ""
echo "=== Configuration ==="
gscribe config show
echo ""
echo "=== Cache Stats ==="
gscribe cache stats
```

## See Also

- [Usage Guide](usage.md) - Detailed usage examples and workflows
- [Configuration](configuration.md) - Complete configuration options
- [Troubleshooting](troubleshooting.md) - Common issues and solutions
- [Workflows](workflows.md) - Team workflows and best practices
