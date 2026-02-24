# Troubleshooting

## Installation Issues

### "command not found: gscribe"

Ensure Go bin is in PATH:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
source ~/.bashrc  # or ~/.zshrc
```

Verify installation:

```bash
which gscribe
go install github.com/tenxprotocols/git-commit-scribe/cmd/gscribe@latest
```

## Configuration Issues

### "API key not configured"

```bash
gscribe config init
# or
gscribe config set api_key "your-key"
# or
export GSCRIBE_API_KEY="your-key"
```

### "Failed to load configuration"

Check YAML syntax and permissions:

```bash
cat ~/.config/gscribe/config.yaml
chmod 600 ~/.config/gscribe/config.yaml
```

Reinitialize if needed:

```bash
gscribe config init --force
```

### Configuration not applied

Check precedence (highest to lowest):
1. CLI flags
2. Environment variables
3. Repository config (`.git/gscribe-config.yaml`)
4. Global config (`~/.config/gscribe/config.yaml`)

```bash
gscribe config show
```

## Git Issues

### "not a git repository"

```bash
git status  # verify you're in a git repo
git init    # or initialize one
```

### "no staged changes to commit"

```bash
git add <files>
git status  # verify changes are staged
```

### "failed to create commit"

Check git hooks and user config:

```bash
git config user.name "Your Name"
git config user.email "your.email@example.com"

# Remove lock if exists
rm .git/index.lock
```

## API Issues

### "401 Unauthorized"

Verify API key:

```bash
gscribe config show | grep api_key
```

Get new key from [OpenRouter](https://openrouter.ai/) if needed.

### "429 Too Many Requests"

Rate limited. Wait and retry:

```bash
sleep 10
gscribe
```

Enable caching to reduce API calls:

```bash
gscribe config set cache.enabled true
```

### "timeout"

Increase timeout or reduce files:

```bash
gscribe config set ai.timeout 60s
gscribe --max-files 20
```

### "Model not found"

Use valid model:

```bash
gscribe config set model "anthropic/claude-3.5-sonnet"
# or
gscribe config set model "openai/gpt-4"
```

See [OpenRouter Models](https://openrouter.ai/models) for options.

## Performance Issues

### Slow generation

```bash
# Enable cache
gscribe config set cache.enabled true

# Use faster model
gscribe config set model "openai/gpt-3.5-turbo"

# Reduce files
gscribe --max-files 30

# Commit in smaller batches
git add src/feature1/
gscribe -y
```

### High memory usage

```bash
# Reduce cache size
gscribe config set cache.max_memory_mb 50

# Clear cache
gscribe cache clear --all
```

## Message Quality Issues

### Poor messages

```bash
# Use better model
gscribe config set model "anthropic/claude-3.5-sonnet"

# Adjust temperature
gscribe config set ai.temperature 0.5

# Specify type and scope
gscribe -t feat -s auth

# Review staged changes
git diff --staged
```

### Wrong commit type

Specify explicitly:

```bash
gscribe -t fix  # or feat, docs, etc.
```

### Breaking changes not detected

Use flag explicitly:

```bash
gscribe --breaking
```

## Debug Mode

Enable verbose output:

```bash
gscribe -v
# or
export GSCRIBE_VERBOSE=true
```

Check cache:

```bash
gscribe cache stats
gscribe cache clear --all
```

Test without committing:

```bash
gscribe -d -v
```

## Quick Diagnostics

Run through this checklist:

```bash
# Is gscribe installed?
which gscribe

# Is config valid?
gscribe config show

# Is API key set?
gscribe config show | grep api_key

# Is it a git repo?
git status

# Are there staged changes?
git diff --staged

# Can we reach the API?
curl -I https://openrouter.ai/

# Does dry run work?
gscribe -d

# What does verbose show?
gscribe -v
```

## Still Having Issues?

Create a bug report with:

```bash
gscribe --version
go version
git --version
uname -a
gscribe config show  # redact API key
gscribe -v [your command]
```

Visit [GitHub Issues](https://github.com/tenxprotocols/git-commit-scribe/issues)
