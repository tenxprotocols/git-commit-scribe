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
