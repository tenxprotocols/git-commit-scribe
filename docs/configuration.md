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

## Advanced Configuration

### Model-Specific Settings

Different models have different characteristics. Optimize your configuration based on the model you use:

#### For Claude Models (Recommended)

```yaml
provider: openrouter
model: anthropic/claude-3.5-sonnet
ai:
  temperature: 0.7
  timeout: 30s
  max_retries: 3
commit:
  description_length: 72
```

#### For GPT Models

```yaml
provider: openrouter
model: openai/gpt-4
ai:
  temperature: 0.6
  timeout: 45s
  max_retries: 3
commit:
  description_length: 80
```

#### For Fast/Cheap Models

```yaml
provider: openrouter
model: openai/gpt-3.5-turbo
ai:
  temperature: 0.5
  timeout: 20s
  max_retries: 2
commit:
  one_line: true  # Faster generation
  description_length: 60
cache:
  enabled: true
  ttl: 604800  # Cache for 7 days
```

### Custom Commit Types

Define your own commit types for your team or project:

```yaml
types:
  feat: "New feature"
  fix: "Bug fix"
  docs: "Documentation"
  style: "Code style"
  refactor: "Refactoring"
  perf: "Performance"
  test: "Testing"
  build: "Build system"
  ci: "CI/CD"
  chore: "Maintenance"
  # Custom types
  hotfix: "Critical production fix"
  wip: "Work in progress"
  experiment: "Experimental changes"
  security: "Security improvements"
```

Usage:
```bash
gscribe -t hotfix -s api
# hotfix(api): patch critical authentication bug

gscribe -t security -s auth
# security(auth): implement rate limiting
```

### Multi-Environment Configuration

#### Development Environment

```yaml
# ~/.config/gscribe/config.yaml
provider: openrouter
model: openai/gpt-3.5-turbo  # Faster for dev
api_key: ${GSCRIBE_API_KEY}

commit:
  emoji: true
  confirm: false  # Auto-approve in dev
  one_line: false

ai:
  temperature: 0.6
  timeout: 20s
  verbose: true  # More logging in dev

cache:
  enabled: true
  ttl: 86400
```

#### Production/CI Environment

```yaml
# .git/gscribe-config.yaml (repo-specific)
provider: openrouter
model: anthropic/claude-3.5-sonnet  # Better quality

commit:
  emoji: false
  confirm: true  # Require review
  max_files: 30

ai:
  temperature: 0.7
  timeout: 30s
  max_retries: 5  # More retries in CI

cache:
  enabled: false  # Don't cache in CI
```

### Project-Specific Patterns

#### Frontend Project

```yaml
# .git/gscribe-config.yaml
commit:
  emoji: true
  description_length: 80
  
types:
  feat: "New UI feature"
  fix: "UI bug fix"
  style: "Styling changes"
  refactor: "Component refactoring"
  test: "Component tests"
  a11y: "Accessibility improvements"
  perf: "Performance optimization"

ai:
  temperature: 0.7
```

#### Backend/API Project

```yaml
# .git/gscribe-config.yaml
commit:
  emoji: false
  description_length: 72
  
types:
  feat: "New API endpoint"
  fix: "API bug fix"
  refactor: "Code refactoring"
  perf: "Performance improvement"
  test: "API tests"
  security: "Security enhancement"
  db: "Database changes"

ai:
  temperature: 0.6
```

#### Documentation Project

```yaml
# .git/gscribe-config.yaml
commit:
  emoji: true
  one_line: true
  confirm: false
  
types:
  docs: "Documentation update"
  content: "Content changes"
  structure: "Structure changes"
  fix: "Fix typos/errors"

ai:
  temperature: 0.5
  timeout: 15s
```

### Monorepo Configuration

For monorepos with multiple packages:

```yaml
# Root .gscribe-config.yaml
commit:
  emoji: true
  max_files: 100  # Higher limit for monorepos
  description_length: 80

ai:
  temperature: 0.7
  timeout: 45s  # More time for larger diffs

cache:
  enabled: true
  max_disk_mb: 2000  # Larger cache for monorepo

# Define package-specific scopes
types:
  feat: "New feature"
  fix: "Bug fix"
  docs: "Documentation"
  # Workspace-specific
  frontend: "Frontend changes"
  backend: "Backend changes"
  shared: "Shared code changes"
  infra: "Infrastructure changes"
```

### Performance Tuning

#### Fast Configuration (Low Latency)

```yaml
provider: openrouter
model: openai/gpt-3.5-turbo

commit:
  one_line: true
  confirm: false
  max_files: 30

ai:
  temperature: 0.3
  timeout: 15s
  max_retries: 2

cache:
  enabled: true
  ttl: 604800  # 7 days
  max_memory_mb: 200
```

#### Quality Configuration (Better Messages)

```yaml
provider: openrouter
model: anthropic/claude-3.5-sonnet

commit:
  one_line: false
  confirm: true
  max_files: 50

ai:
  temperature: 0.7
  timeout: 45s
  max_retries: 5

cache:
  enabled: true
  ttl: 86400  # 1 day
```

#### Balanced Configuration

```yaml
provider: openrouter
model: anthropic/claude-3.5-sonnet

commit:
  one_line: false
  confirm: true
  max_files: 40

ai:
  temperature: 0.6
  timeout: 30s
  max_retries: 3

cache:
  enabled: true
  ttl: 259200  # 3 days
  max_memory_mb: 100
  max_disk_mb: 1000
```

## Configuration Recipes

### Strict Mode (Team Lead)

Review everything, enforce quality:

```yaml
commit:
  confirm: true
  emoji: false
  max_files: 30
  ignore_generated: true
  ignore_whitespace: true

ai:
  temperature: 0.7
  timeout: 45s
  max_retries: 5
```

### Speed Mode (Rapid Development)

Fast commits, minimal friction:

```yaml
commit:
  confirm: false
  one_line: true
  emoji: true
  auto_push: false

ai:
  temperature: 0.5
  timeout: 20s

cache:
  enabled: true
  ttl: 604800
```

### CI/CD Mode

Automated, reliable commits:

```yaml
commit:
  confirm: false
  emoji: false
  max_files: 100

ai:
  temperature: 0.6
  timeout: 60s
  max_retries: 10
  rate_limit_wait: 15s

cache:
  enabled: false
```

### Documentation Mode

Quick, simple commits for docs:

```yaml
commit:
  confirm: false
  one_line: true
  emoji: true
  description_length: 100

ai:
  temperature: 0.5
  timeout: 15s
```

## Configuration Management

### Backup Configuration

```bash
# Backup global config
cp ~/.config/gscribe/config.yaml ~/.config/gscribe/config.yaml.bak

# Backup all repo configs
find ~/projects -name "gscribe-config.yaml" -exec cp {} {}.bak \;
```

### Sync Configuration Across Machines

```bash
# Export configuration
cat ~/.config/gscribe/config.yaml

# On new machine, import (without API key)
cat > ~/.config/gscribe/config.yaml <<'EOF'
provider: openrouter
model: anthropic/claude-3.5-sonnet
commit:
  emoji: true
  confirm: true
ai:
  temperature: 0.7
EOF

# Set API key separately
gscribe config set api_key "your-key-here"
```

### Version Control for Team Config

```bash
# Create shareable config (no secrets)
cat > .gscribe-config.yaml <<'EOF'
commit:
  emoji: true
  confirm: true
  description_length: 72
  max_files: 50

ai:
  temperature: 0.7
  timeout: 30s

types:
  feat: "New feature"
  fix: "Bug fix"
  docs: "Documentation"
  # ... more types
EOF

# Add to git
git add .gscribe-config.yaml
git commit -m "chore: add gscribe team configuration"
```

### Migrate Configuration

```bash
# Export current config
gscribe config show > current-config.yaml

# Edit as needed
vim current-config.yaml

# Apply to new location
cp current-config.yaml .git/gscribe-config.yaml
```

## Troubleshooting Configuration

### Verify Active Configuration

```bash
# Show merged configuration
gscribe config show

# Show with sources
gscribe config show -v

# Show only global
gscribe config show --global
```

### Test Configuration

```bash
# Test with dry run
echo "test" > test.txt
git add test.txt
gscribe -d -v
git reset HEAD test.txt
rm test.txt
```

### Reset Configuration

```bash
# Remove global config
rm ~/.config/gscribe/config.yaml

# Reinitialize
gscribe config init

# Remove repo config
rm .git/gscribe-config.yaml
```

### Debug Configuration Issues

```bash
# Check file permissions
ls -la ~/.config/gscribe/config.yaml

# Should be: -rw------- (600)
chmod 600 ~/.config/gscribe/config.yaml

# Validate YAML syntax
cat ~/.config/gscribe/config.yaml | python -c "import yaml, sys; yaml.safe_load(sys.stdin)"

# Check environment variables
env | grep GSCRIBE
```

## Best Practices

1. **Keep API keys out of version control**
   - Use environment variables
   - Add config files with secrets to `.gitignore`

2. **Use repository configs for project-specific settings**
   - Keep global config minimal
   - Override with repo-specific needs

3. **Enable caching for better performance**
   - Especially useful for repeated operations
   - Clear cache periodically

4. **Adjust temperature based on use case**
   - Lower (0.3-0.5) for consistency
   - Higher (0.7-0.9) for creativity

5. **Set appropriate timeouts**
   - Longer for large repositories
   - Shorter for quick commits

6. **Configure max_files appropriately**
   - Lower for atomic commits
   - Higher for batch operations

## See Also

- [Usage Guide](usage.md) - How to use gscribe effectively
- [CLI Reference](cli-reference.md) - Complete command reference
- [Troubleshooting](troubleshooting.md) - Common configuration issues
- [Workflows](workflows.md) - Team configuration patterns
