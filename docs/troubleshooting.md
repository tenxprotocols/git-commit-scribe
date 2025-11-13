# Troubleshooting Guide

This guide covers common issues and their solutions when using git-commit-scribe.

## Table of Contents

- [Installation Issues](#installation-issues)
- [Configuration Issues](#configuration-issues)
- [Git-Related Issues](#git-related-issues)
- [API Issues](#api-issues)
- [Performance Issues](#performance-issues)
- [Generate Message Issues](#generate-message-issues)
- [Platform-Specific Issues](#platform-specific-issues)
- [Debug Mode](#debug-mode)

---

## Installation Issues

### "command not found: gscribe"

**Problem:** After installation, `gscribe` command is not recognized.

**Solutions:**

1. **Ensure Go bin is in PATH:**
   ```bash
   # Add to ~/.bashrc, ~/.zshrc, or equivalent
   export PATH="$PATH:$(go env GOPATH)/bin"
   
   # Reload shell
   source ~/.bashrc  # or ~/.zshrc
   ```

2. **Verify installation:**
   ```bash
   which gscribe
   # Should show path like: /Users/username/go/bin/gscribe
   ```

3. **Reinstall:**
   ```bash
   go install github.com/tenxprotocols/git-commit-scribe/cmd/gscribe@latest
   ```

### Build Errors

**Problem:** Compilation fails when building from source.

**Solutions:**

1. **Check Go version:**
   ```bash
   go version
   # Requires Go 1.21 or higher
   ```

2. **Update dependencies:**
   ```bash
   go mod download
   go mod tidy
   ```

3. **Clean and rebuild:**
   ```bash
   go clean -cache
   ./scripts/build.sh
   ```

---

## Configuration Issues

### "API key not configured"

**Problem:**
```
Error: API key not configured
```

**Solutions:**

1. **Initialize configuration:**
   ```bash
   gscribe config init
   ```

2. **Manually set API key:**
   ```bash
   gscribe config set api_key "your-api-key-here"
   ```

3. **Use environment variable:**
   ```bash
   export GSCRIBE_API_KEY="your-api-key-here"
   gscribe
   ```

### "Failed to load configuration"

**Problem:** Configuration file is corrupted or invalid.

**Solutions:**

1. **Check configuration syntax:**
   ```bash
   cat ~/.config/gscribe/config.yaml
   # Ensure valid YAML format
   ```

2. **Reinitialize configuration:**
   ```bash
   gscribe config init --force
   ```

3. **Check file permissions:**
   ```bash
   ls -la ~/.config/gscribe/config.yaml
   # Should show: -rw------- (600)
   
   # Fix if needed:
   chmod 600 ~/.config/gscribe/config.yaml
   ```

### Configuration Not Applied

**Problem:** Changes to config file don't take effect.

**Solutions:**

1. **Check for repository-specific config:**
   ```bash
   # Repository config overrides global config
   cat .git/gscribe-config.yaml
   ```

2. **Verify configuration precedence:**
   ```
   Priority (highest to lowest):
   1. Command-line flags
   2. Environment variables
   3. Repository config (.git/gscribe-config.yaml)
   4. Global config (~/.config/gscribe/config.yaml)
   ```

3. **Use verbose mode to see active config:**
   ```bash
   gscribe config show -v
   ```

---

## Git-Related Issues

### "not a git repository"

**Problem:**
```
Error: not a git repository
```

**Solutions:**

1. **Ensure you're in a git repository:**
   ```bash
   git status
   # If fails, initialize:
   git init
   ```

2. **Check current directory:**
   ```bash
   pwd
   # Make sure you're in the project root
   ```

### "no staged changes to commit"

**Problem:**
```
Error: no staged changes to commit
```

**Solutions:**

1. **Stage your changes:**
   ```bash
   git add <files>
   # or
   git add .
   ```

2. **Verify staged changes:**
   ```bash
   git status
   # Should show files under "Changes to be committed:"
   ```

3. **Check if working tree is clean:**
   ```bash
   git diff --staged
   # Should show differences
   ```

### "failed to create commit"

**Problem:** Commit creation fails after message generation.

**Solutions:**

1. **Check for git hooks:**
   ```bash
   ls .git/hooks/
   # Check for pre-commit or commit-msg hooks that might be failing
   ```

2. **Verify git user config:**
   ```bash
   git config user.name
   git config user.email
   
   # If not set:
   git config --global user.name "Your Name"
   git config --global user.email "your.email@example.com"
   ```

3. **Check for index lock:**
   ```bash
   # If .git/index.lock exists, remove it:
   rm .git/index.lock
   ```

---

## API Issues

### "API request failed: 401 Unauthorized"

**Problem:** Authentication with AI provider fails.

**Solutions:**

1. **Verify API key:**
   ```bash
   gscribe config show | grep api_key
   # Check if key is set correctly
   ```

2. **Check API key validity:**
   - Log in to [OpenRouter](https://openrouter.ai/)
   - Verify key is active and not revoked
   - Check usage limits

3. **Update API key:**
   ```bash
   gscribe config set api_key "new-api-key"
   ```

### "API request failed: 429 Too Many Requests"

**Problem:** Rate limit exceeded.

**Solutions:**

1. **Wait and retry:**
   ```bash
   # Wait a few seconds and try again
   sleep 5
   gscribe
   ```

2. **Check rate limit settings:**
   ```yaml
   # In config.yaml
   ai:
     max_retries: 5
     rate_limit_wait: 10s
   ```

3. **Use caching to reduce API calls:**
   ```bash
   gscribe config set cache.enabled true
   ```

4. **Check OpenRouter dashboard:**
   - View your rate limits and usage
   - Consider upgrading plan if needed

### "API request failed: timeout"

**Problem:** Request takes too long and times out.

**Solutions:**

1. **Increase timeout:**
   ```bash
   gscribe config set ai.timeout 60s
   ```

2. **Check network connection:**
   ```bash
   curl -I https://openrouter.ai/
   ```

3. **Reduce number of files:**
   ```bash
   gscribe --max-files 20
   ```

### "Model not found" or "Invalid model"

**Problem:**
```
Error: model not available: invalid-model-name
```

**Solutions:**

1. **Use a valid model:**
   ```bash
   # Recommended models:
   gscribe config set model "anthropic/claude-3.5-sonnet"
   gscribe config set model "openai/gpt-4"
   gscribe config set model "google/gemini-pro"
   ```

2. **Check available models:**
   - Visit [OpenRouter Models](https://openrouter.ai/models)
   - Verify model name and availability

3. **Check account access:**
   - Some models require special access or credits

---

## Performance Issues

### Slow Commit Generation

**Problem:** `gscribe` takes a long time to generate commits.

**Solutions:**

1. **Enable caching:**
   ```bash
   gscribe config set cache.enabled true
   ```

2. **Reduce file limit:**
   ```bash
   gscribe --max-files 30
   ```

3. **Use a faster model:**
   ```bash
   # Some models are faster than others
   gscribe config set model "openai/gpt-3.5-turbo"
   ```

4. **Check diff size:**
   ```bash
   git diff --staged --stat
   # Large diffs take longer to process
   ```

5. **Commit in smaller batches:**
   ```bash
   # Instead of committing everything at once
   git add src/feature1/
   gscribe -y
   git add src/feature2/
   gscribe -y
   ```

### High Memory Usage

**Problem:** `gscribe` uses too much memory.

**Solutions:**

1. **Reduce cache size:**
   ```yaml
   cache:
     max_memory_mb: 50
     max_disk_mb: 500
   ```

2. **Clear cache:**
   ```bash
   gscribe cache clear --all
   ```

3. **Disable memory cache:**
   ```bash
   gscribe config set cache.enabled false
   ```

---

## Generate Message Issues

### Poor Quality Messages

**Problem:** Generated commit messages are not helpful or accurate.

**Solutions:**

1. **Use a better model:**
   ```bash
   gscribe config set model "anthropic/claude-3.5-sonnet"
   ```

2. **Adjust temperature:**
   ```yaml
   # Lower temperature for more focused output
   ai:
     temperature: 0.5
   ```

3. **Review staged changes:**
   ```bash
   git diff --staged
   # Make sure you're committing logical units of work
   ```

4. **Use custom prompts:**
   ```bash
   # See docs/prompts.md for custom prompt configuration
   ```

5. **Specify commit type and scope:**
   ```bash
   gscribe -t feat -s auth
   ```

### Messages Too Long or Too Short

**Problem:** Generated messages don't match desired length.

**Solutions:**

1. **Use one-line mode:**
   ```bash
   gscribe --one-line
   ```

2. **Adjust description length:**
   ```bash
   gscribe --description 100
   ```

3. **Configure defaults:**
   ```yaml
   commit:
     description_length: 72
     one_line: false
   ```

### Incorrect Commit Type

**Problem:** AI chooses the wrong commit type.

**Solutions:**

1. **Specify type explicitly:**
   ```bash
   gscribe -t fix
   ```

2. **Review your changes:**
   ```bash
   git diff --staged
   # Make sure changes match a single type
   ```

3. **Customize prompts:**
   - See `docs/prompts.md` for guidance
   - Add examples to your custom prompt

### Breaking Changes Not Detected

**Problem:** Breaking changes not identified automatically.

**Solutions:**

1. **Use breaking flag explicitly:**
   ```bash
   gscribe --breaking
   ```

2. **Review generated message:**
   ```bash
   gscribe -d  # Dry run first
   # Then confirm if correct
   ```

---

## Platform-Specific Issues

### macOS: "cannot be opened because the developer cannot be verified"

**Problem:** Security warning when running binary.

**Solutions:**

1. **Allow the application:**
   ```bash
   xattr -d com.apple.quarantine $(which gscribe)
   ```

2. **Or install via Homebrew (if available):**
   ```bash
   brew install git-commit-scribe
   ```

### Windows: Path Issues

**Problem:** Command not found on Windows.

**Solutions:**

1. **Add to PATH manually:**
   - Open System Properties > Environment Variables
   - Add `%USERPROFILE%\go\bin` to PATH

2. **Use PowerShell:**
   ```powershell
   $env:Path += ";$env:USERPROFILE\go\bin"
   ```

### Linux: Permission Denied

**Problem:**
```
bash: /path/to/gscribe: Permission denied
```

**Solutions:**

1. **Make executable:**
   ```bash
   chmod +x $(which gscribe)
   ```

2. **Check installation directory permissions:**
   ```bash
   ls -la $(go env GOPATH)/bin/gscribe
   ```

---

## Debug Mode

### Enabling Verbose Output

Get detailed information about what's happening:

```bash
# Via flag
gscribe -v

# Via environment variable
export GSCRIBE_VERBOSE=true
gscribe

# Via config
gscribe config set verbose true
```

**Verbose output includes:**
- Configuration loaded
- Files being analyzed
- API requests and responses
- Cache hits/misses
- Timing information

### Checking Cache State

View cache statistics:

```bash
gscribe cache stats
```

Output:
```
Cache Statistics:
  Total entries: 15
  Memory usage: 2.3 MB
  Disk usage: 12.5 MB
  Hit rate: 67%
  Oldest entry: 2024-01-15 10:30:00
```

Clear cache if needed:

```bash
gscribe cache clear --all
```

### Testing Configuration

Test without making commits:

```bash
# Dry run to see what would be generated
gscribe -d -v

# Show effective configuration
gscribe config show

# Test API connectivity
echo "test" > test.txt
git add test.txt
gscribe -d -v
git reset HEAD test.txt
rm test.txt
```

---

## Getting Help

If you're still experiencing issues:

1. **Check existing issues:**
   - Visit [GitHub Issues](https://github.com/tenxprotocols/git-commit-scribe/issues)

2. **Create a bug report:**
   ```bash
   # Include this information:
   gscribe --version
   go version
   git --version
   uname -a  # or systeminfo on Windows
   
   # Include configuration (redact API key):
   gscribe config show
   
   # Include error output with verbose mode:
   gscribe -v [your command]
   ```

3. **Community support:**
   - Check documentation: https://github.com/tenxprotocols/git-commit-scribe
   - Search discussions
   - Ask in community forums

---

## Quick Diagnostics Checklist

Run through this checklist to diagnose issues:

- [ ] Is `gscribe` installed? (`which gscribe`)
- [ ] Is it in PATH? (`echo $PATH`)
- [ ] Is configuration valid? (`gscribe config show`)
- [ ] Is API key set? (`gscribe config show | grep api_key`)
- [ ] Is it a git repository? (`git status`)
- [ ] Are there staged changes? (`git diff --staged`)
- [ ] Is network accessible? (`curl -I https://openrouter.ai/`)
- [ ] Does dry run work? (`gscribe -d`)
- [ ] What does verbose show? (`gscribe -v`)

If all checks pass but issues persist, use verbose mode and review the output carefully for specific error messages.
