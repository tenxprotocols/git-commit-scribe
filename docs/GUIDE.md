# Quick Guide

## Installation

```bash
go install github.com/tenxprotocols/git-commit-scribe/cmd/gscribe@latest
```

## Setup

```bash
gscribe config init
```

You'll need an API key from [OpenRouter](https://openrouter.ai/).

## Basic Usage

```bash
# Make changes
git add .

# Generate commit
gscribe

# Skip confirmation
gscribe -y

# Dry run (preview only)
gscribe -d
```

## Common Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--yes` | `-y` | Skip confirmation |
| `--type` | `-t` | Commit type (feat, fix, docs, etc.) |
| `--scope` | `-s` | Commit scope |
| `--breaking` | `-b` | Mark as breaking change |
| `--emoji` | `-e` | Add emoji |
| `--dry-run` | `-d` | Preview only |
| `--push` | `-p` | Push after commit |
| `--verbose` | `-v` | Show details |

## Examples

```bash
# Feature with scope
gscribe -t feat -s api

# Bug fix with emoji
gscribe -t fix -e

# Quick commit and push
gscribe -yp

# Breaking change
gscribe -t feat --breaking
```

## Configuration

Config file: `~/.config/gscribe/config.yaml`

### Essential Options

```yaml
# API
provider: openrouter
model: anthropic/claude-3.5-sonnet
api_key: your-key

# Commit preferences
commit:
  emoji: false
  confirm: true
  max_files: 50

# Cache
cache:
  enabled: true
  ttl: 86400  # seconds
```

### Config Commands

```bash
# View config
gscribe config show

# Set value
gscribe config set commit.emoji true

# Remove value
gscribe config unset commit.emoji
```

### Environment Variables

```bash
export GSCRIBE_API_KEY="your-key"
export GSCRIBE_MODEL="anthropic/claude-3.5-sonnet"
export GSCRIBE_VERBOSE="true"
```

## Workflows

### Solo Developer

```bash
# Daily workflow
git add .
gscribe -y
git push
```

### Feature Development

```bash
# Create branch
git checkout -b feat/new-feature

# Make commits
git add src/component.js
gscribe -t feat -s ui

git add tests/
gscribe -t test -s ui

# Push
git push -u origin feat/new-feature
```

### Team Project

Add `.gscribe-config.yaml` to your repo:

```yaml
commit:
  emoji: true
  confirm: true
  description_length: 72
```

### Quick Aliases

```bash
# Add to ~/.bashrc or ~/.zshrc
alias gc='git add . && gscribe -y'
alias gcf='git add . && gscribe -t feat -y'
alias gcx='git add . && gscribe -t fix -y'
```

## Cache

Cache stores generated messages to avoid repeated API calls.

```bash
# View stats
gscribe cache stats

# Clear cache
gscribe cache clear
```

Cache location: `~/.cache/gscribe/`

## Commit Types

- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation
- `style` - Code style/formatting
- `refactor` - Code refactoring
- `perf` - Performance improvement
- `test` - Tests
- `build` - Build system
- `ci` - CI/CD
- `chore` - Maintenance

## Tips

1. **Review before committing** - Keep `confirm: true`
2. **Use scopes** - Helps organize commits
3. **Enable cache** - Speeds up repeated operations
4. **Commit small changes** - Better than large dumps
5. **Use verbose mode** - When debugging: `gscribe -v`
