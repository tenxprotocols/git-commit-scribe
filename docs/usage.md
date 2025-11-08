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

## Troubleshooting

### No Staged Changes

```
Error: no staged changes to commit
```

Solution: Stage your changes first:
```bash
git add <files>
```

### API Key Not Configured

```
Error: API key not configured
```

Solution: Initialize your configuration:
```bash
gscribe config init
```

### Too Many Files

```
Error: too many files changed (75), max is 50
```

Solution: Either increase the limit or commit in batches:
```bash
gscribe --max-files 100
```

### Rate Limiting

If you hit API rate limits, the tool will automatically retry with exponential backoff.

You can adjust retry behavior in the config:
```yaml
ai:
  max_retries: 5
  rate_limit_wait: 10s
