# Usage Guide

## Basic Workflow

The typical workflow with gscribe:

1. Make changes to your code
2. Stage the changes with git
3. Run gscribe to generate and create the commit
4. Optionally push to remote

```bash
# Make changes
vim myfile.js

# Stage changes
git add myfile.js

# Generate commit message and create commit
gscribe

# Or do all at once with auto-push
gscribe -yp
```

## Interactive Mode

By default, gscribe runs interactively and shows you the generated commit message before creating the commit:

```bash
$ gscribe

Generated commit message:
─────────────────────────────────────────
feat(api): add user authentication

Implement JWT-based authentication for API endpoints.
Includes login, logout, and token refresh functionality.
─────────────────────────────────────────

Create commit with this message? [Y/n]: y
✓ Commit created successfully
```

## Dry Run Mode

Preview the commit message without creating a commit:

```bash
gscribe --dry-run
# or
gscribe -d
```

This is useful for:
- Testing your configuration
- Seeing what message would be generated
- Experimenting with different flags

## Common Scenarios

### Quick Commit

Skip the confirmation prompt:

```bash
gscribe -y
```

### Feature Commit with Scope

```bash
gscribe -t feat -s api
```

### Breaking Change

Mark a commit as breaking:

```bash
gscribe --breaking
# or
gscribe -b
```

### One-Line Commit

Generate a concise one-line message:

```bash
gscribe --one-line
# or
gscribe -o
```

### Commit with Emoji

Add emoji to your commit messages:

```bash
gscribe --emoji
# or
gscribe -e
```

Output: `✨ feat(api): add user authentication`

### Commit and Push

Create the commit and immediately push:

```bash
gscribe --push
# or
gscribe -p
```

### Combine Multiple Flags

```bash
# Feature with emoji, skip confirmation, and push
gscribe -t feat -e -y -p

# Breaking change with verbose output
gscribe -b -v

# Dry run with emoji to preview
gscribe -d -e
```

## Handling Large Changes

### Too Many Files

If you've changed too many files:

```bash
# Increase the limit
gscribe --max-files 100
```

Or commit in smaller batches:

```bash
git add src/feature1/
gscribe -y
git add src/feature2/
gscribe -y
```

### Review Generated Messages

Use verbose mode to see what's happening:

```bash
gscribe -v
```

Output:
```
Getting staged changes...
Analyzing 3 staged file(s)...
Generating commit message using anthropic/claude-3.5-sonnet...
Creating commit...
✓ Commit created successfully
```

## Conventional Commits

gscribe follows [Conventional Commits v1.0.0](https://www.conventionalcommits.org/).

### Commit Types

Default types available:

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting)
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `test`: Test changes
- `build`: Build system changes
- `ci`: CI configuration
- `chore`: Other changes
- `revert`: Revert previous commit

### Message Structure

```
<type>[optional scope][!]: <description>

[optional body]

[optional footer(s)]
```

Examples:

```
feat(auth): add OAuth2 support

fix: resolve memory leak in cache

docs: update API documentation

feat(api)!: redesign authentication flow

BREAKING CHANGE: Authentication now requires OAuth2
```

## Tips and Best Practices

### 1. Stage Related Changes Together

Group related changes in a single commit:

```bash
git add src/auth/*.js src/middleware/auth.js
gscribe -t feat -s auth
```

### 2. Use Meaningful Scopes

Scopes help organize commits by feature or component:

```bash
gscribe -s api      # API changes
gscribe -s ui       # UI changes
gscribe -s docs     # Documentation
gscribe -s tests    # Test changes
```

### 3. Review Before Committing

Always review the generated message. If it's not quite right:
- Answer "n" at the confirmation prompt
- Adjust your staged changes
- Try again with different flags

### 4. Configure Defaults

Set your preferred defaults in the config file instead of using flags every time:

```yaml
commit:
  emoji: true
  confirm: false  # Skip confirmations
```

### 5. Use Verbose Mode While Learning

Enable verbose mode to understand what's happening:

```bash
export GSCRIBE_VERBOSE=true
gscribe
```

## Advanced Usage

### Environment-Specific Configuration

Use different configurations based on environment:

```bash
# Development - more verbose
export GSCRIBE_VERBOSE=true
export GSCRIBE_MODEL="openai/gpt-3.5-turbo"  # Faster, cheaper
gscribe

# Production - more careful
export GSCRIBE_VERBOSE=false
export GSCRIBE_MODEL="anthropic/claude-3.5-sonnet"  # Better quality
gscribe
```

### Custom Workflows

#### Multi-Repository Workflow

Working across multiple repositories:

```bash
# Script to commit changes across repos
for repo in ~/projects/*/; do
  cd "$repo"
  if [[ -n $(git status -s) ]]; then
    echo "Processing $repo..."
    git add .
    gscribe -y
  fi
done
```

#### Feature Branch Workflow

```bash
# Create feature branch
git checkout -b feat/user-profile

# Make incremental commits
git add src/profile/avatar.js
gscribe -t feat -s profile -e
# ✨ feat(profile): add avatar upload

git add src/profile/settings.js
gscribe -t feat -s profile -e
# ✨ feat(profile): add user settings page

git add tests/profile/
gscribe -t test -s profile
# test(profile): add profile component tests

# Push feature
git push -u origin feat/user-profile
```

#### Hotfix Workflow

```bash
# Create hotfix branch from main
git checkout main
git checkout -b hotfix/security-patch

# Make fix
git add src/auth/validator.js
gscribe -t fix -s security --breaking
# fix(security)!: patch authentication vulnerability
# BREAKING CHANGE: Updated token validation algorithm

# Push and merge immediately
git push -u origin hotfix/security-patch
```

### Integration with Git Aliases

Add gscribe to your git aliases:

```bash
# In ~/.gitconfig
[alias]
  cm = !git add . && gscribe -y
  cmf = !git add . && gscribe -t feat -y
  cmx = !git add . && gscribe -t fix -y
  cmd = !git add . && gscribe -t docs -y
  cmp = !git add . && gscribe -yp
```

Usage:
```bash
git cm    # Quick commit
git cmf   # Feature commit
git cmx   # Fix commit
git cmd   # Docs commit  
git cmp   # Commit and push
```

### Batch Processing

#### Commit Multiple Directories

```bash
# Commit each directory separately with appropriate scope
for dir in src/components/*/; do
  dirname=$(basename "$dir")
  git add "$dir"
  gscribe -t feat -s "components/$dirname" -y
done
```

#### Process Files by Type

```bash
# Commit all documentation
git add *.md docs/
gscribe -t docs -o -y

# Commit all tests
git add **/*.test.js
gscribe -t test -y

# Commit all style changes
git add **/*.css **/*.scss
gscribe -t style -y
```

### Using with Git Hooks

#### Prepare Commit Message Hook

Create `.git/hooks/prepare-commit-msg`:

```bash
#!/bin/bash
# Auto-generate commit message if not provided

COMMIT_MSG_FILE=$1
COMMIT_SOURCE=$2

# Only generate for normal commits (not merge, squash, etc.)
if [ -z "$COMMIT_SOURCE" ]; then
  # Generate message with gscribe
  gscribe -d > "$COMMIT_MSG_FILE"
fi
```

Make executable:
```bash
chmod +x .git/hooks/prepare-commit-msg
```

### Performance Optimization

#### Enable Caching

Significantly speeds up repeated operations:

```bash
gscribe config set cache.enabled true
gscribe config set cache.ttl 86400  # 24 hours
```

#### Optimize for Large Repositories

```bash
# Increase file limit
gscribe config set commit.max_files 100

# Use faster model
gscribe config set model "openai/gpt-3.5-turbo"

# Disable detailed analysis for speed
gscribe config set ai.temperature 0.3
```

#### Reduce API Calls

```bash
# Use one-line mode (shorter prompts)
gscribe config set commit.one_line true

# Use cache aggressively
gscribe config set cache.ttl 604800  # 7 days
```

## Integration Examples

### VS Code Integration

Create a VS Code task (`.vscode/tasks.json`):

```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "gscribe: commit",
      "type": "shell",
      "command": "gscribe",
      "problemMatcher": []
    },
    {
      "label": "gscribe: quick commit",
      "type": "shell",
      "command": "gscribe -y",
      "problemMatcher": []
    },
    {
      "label": "gscribe: commit and push",
      "type": "shell",
      "command": "gscribe -yp",
      "problemMatcher": []
    }
  ]
}
```

### Shell Functions

Add to `.bashrc` or `.zshrc`:

```bash
# Quick commit function
gc() {
  git add .
  gscribe -y
}

# Feature commit
gcf() {
  git add .
  gscribe -t feat -s "$1" -y
}

# Fix commit
gcx() {
  git add .
  gscribe -t fix -s "$1" -y
}

# Commit with emoji and push
gcp() {
  git add .
  gscribe -e -yp
}
```

Usage:
```bash
gc           # Quick commit
gcf api      # Feature commit with scope
gcx auth     # Fix commit with scope
gcp          # Commit with emoji and push
```

### Makefile Integration

Add to `Makefile`:

```makefile
.PHONY: commit commit-push commit-feat commit-fix

commit:
	@git add .
	@gscribe

commit-push:
	@git add .
	@gscribe -yp

commit-feat:
	@git add .
	@gscribe -t feat -y

commit-fix:
	@git add .
	@gscribe -t fix -y
```

Usage:
```bash
make commit        # Interactive commit
make commit-push   # Commit and push
make commit-feat   # Feature commit
make commit-fix    # Fix commit
```

## Troubleshooting

For detailed troubleshooting information, see the [Troubleshooting Guide](troubleshooting.md).

### Quick Fixes

#### No Staged Changes

```
Error: no staged changes to commit
```

Solution: Stage your changes first:
```bash
git add <files>
```

#### API Key Not Configured

```
Error: API key not configured
```

Solution: Initialize your configuration:
```bash
gscribe config init
```

#### Too Many Files

```
Error: too many files changed (75), max is 50
```

Solution: Either increase the limit or commit in batches:
```bash
gscribe --max-files 100
```

#### Rate Limiting

If you hit API rate limits, the tool will automatically retry with exponential backoff.

You can adjust retry behavior in the config:
```yaml
ai:
  max_retries: 5
  rate_limit_wait: 10s
```

## See Also

- [Configuration Guide](configuration.md) - Detailed configuration options
- [CLI Reference](cli-reference.md) - Complete command reference
- [Troubleshooting Guide](troubleshooting.md) - Common issues and solutions
- [Workflows and Best Practices](workflows.md) - Advanced workflows and team practices
