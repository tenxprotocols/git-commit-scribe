# Reference

## CLI Flags

### Global Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--verbose` | `-v` | bool | false | Verbose output |
| `--no-cache` | | bool | false | Disable cache |
| `--help` | `-h` | | | Show help |

### Commit Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--yes` | `-y` | bool | false | Skip confirmation |
| `--type` | `-t` | string | | Commit type |
| `--scope` | `-s` | string | | Commit scope |
| `--breaking` | `-b` | bool | false | Breaking change |
| `--emoji` | `-e` | bool | false | Add emoji |
| `--one-line` | `-o` | bool | false | One-line message |
| `--dry-run` | `-d` | bool | false | Preview only |
| `--push` | `-p` | bool | false | Push after commit |
| `--max-files` | | int | 50 | Max files to analyze |

## Config Options

### API Settings

```yaml
provider: openrouter
model: anthropic/claude-3.5-sonnet
api_key: sk-...
```

### AI Settings

```yaml
ai:
  temperature: 0.7          # 0.0-1.0, creativity level
  timeout: 30s              # Request timeout
  max_retries: 3            # Retry attempts
  rate_limit_wait: 5s       # Wait between retries
```

### Commit Settings

```yaml
commit:
  emoji: false              # Include emoji
  one_line: false           # One-line messages
  description_length: 72    # Max description length
  max_files: 50             # Max files to analyze
  ignore_generated: true    # Skip auto-generated files
  ignore_whitespace: true   # Skip whitespace-only changes
  auto_push: false          # Auto push after commit
  confirm: true             # Ask before committing
```

### Cache Settings

```yaml
cache:
  enabled: true             # Enable caching
  ttl: 86400                # Time-to-live (seconds)
  max_memory_mb: 100        # Memory cache limit
  max_disk_mb: 1000         # Disk cache limit
```

### Commit Types

```yaml
types:
  feat: "New feature"
  fix: "Bug fix"
  docs: "Documentation"
  style: "Code style"
  refactor: "Refactoring"
  perf: "Performance"
  test: "Tests"
  build: "Build system"
  ci: "CI/CD"
  chore: "Maintenance"
  revert: "Revert commit"
```

## Environment Variables

All settings can be overridden with environment variables using the `GSCRIBE_` prefix:

```bash
GSCRIBE_API_KEY="sk-..."
GSCRIBE_MODEL="anthropic/claude-3.5-sonnet"
GSCRIBE_PROVIDER="openrouter"
GSCRIBE_VERBOSE="true"
GSCRIBE_NO_CACHE="true"
```

## Config Locations

1. **Global**: `~/.config/gscribe/config.yaml`
2. **Repository**: `.git/gscribe-config.yaml` (overrides global)

## Commands

### gscribe (commit)

Generate and create a commit:

```bash
gscribe [flags]
```

### gscribe config init

Initialize configuration:

```bash
gscribe config init [--force]
```

### gscribe config show

Display configuration:

```bash
gscribe config show [--global]
```

### gscribe config set

Set config value:

```bash
gscribe config set <key> <value> [--repo]
```

### gscribe config unset

Remove config value:

```bash
gscribe config unset <key> [--repo]
```

### gscribe cache stats

Show cache statistics:

```bash
gscribe cache stats
```

### gscribe cache clear

Clear cache:

```bash
gscribe cache clear [--all]
```

## Cache System

- **Memory cache**: Fast, LRU eviction
- **Disk cache**: Persistent, gzip compressed
- **Location**: `~/.cache/gscribe/`
- **Key**: SHA256 hash of diff + options

Cache entries include:
- Git diff content
- Commit type/scope
- Breaking change flag
- One-line format flag
- AI model name

## Exit Codes

| Code | Description |
|------|-------------|
| 0 | Success |
| 1 | General error |
| 2 | Configuration error |
| 3 | Git error |
| 4 | API error |

## Conventional Commits

Messages follow [Conventional Commits v1.0.0](https://www.conventionalcommits.org/):

```
<type>[optional scope][!]: <description>

[optional body]

[optional footer]
```

Examples:

```
feat(api): add user authentication
fix: resolve memory leak
docs: update README
feat(api)!: redesign authentication

BREAKING CHANGE: Old auth tokens no longer valid
