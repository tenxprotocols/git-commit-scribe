# Workflows and Best Practices

This guide provides practical workflows and best practices for using git-commit-scribe effectively in different scenarios.

## Table of Contents

- [Individual Developer Workflows](#individual-developer-workflows)
- [Team Workflows](#team-workflows)
- [CI/CD Integration](#cicd-integration)
- [Monorepo Workflows](#monorepo-workflows)
- [Advanced Workflows](#advanced-workflows)
- [Best Practices](#best-practices)
- [Team Conventions](#team-conventions)

---

## Individual Developer Workflows

### Quick Daily Workflow

The simplest workflow for individual developers:

```bash
# 1. Make your changes
vim src/feature.js

# 2. Stage and commit in one go
git add .
gscribe -y

# 3. Push to remote
git push
```

**Optimization:**
```bash
# Set up auto-approve and auto-push globally
gscribe config set commit.confirm false
gscribe config set commit.auto_push true

# Now commits and pushes automatically
git add .
gscribe
```

### Feature Development Workflow

Structured workflow for developing features:

```bash
# 1. Create feature branch
git checkout -b feat/user-authentication

# 2. Work on specific components
vim src/auth/login.js
vim src/auth/middleware.js

# 3. Commit logical units
git add src/auth/login.js
gscribe -t feat -s auth -e
# Result: ✨ feat(auth): implement login functionality

git add src/auth/middleware.js
gscribe -t feat -s auth
# Result: feat(auth): add authentication middleware

# 4. Push feature branch
git push -u origin feat/user-authentication

# 5. Create pull request
gh pr create --title "Add user authentication" --body "..."
```

### Bug Fix Workflow

Workflow optimized for bug fixes:

```bash
# 1. Create bug fix branch
git checkout -b fix/memory-leak

# 2. Make the fix
vim src/cache/manager.js

# 3. Commit with fix type
git add src/cache/manager.js
gscribe -t fix -s cache
# Result: fix(cache): resolve memory leak in cache manager

# 4. Add tests
vim tests/cache/manager.test.js
git add tests/cache/manager.test.js
gscribe -t test -s cache

# 5. Push and create PR
git push -u origin fix/memory-leak
```

### Documentation Workflow

Efficient workflow for documentation updates:

```bash
# Configure for docs
gscribe config set commit.one_line true --repo
gscribe config set commit.confirm false --repo

# Update docs
vim README.md docs/api.md
git add README.md docs/api.md
gscribe -t docs

# Result: docs: update README and API documentation
```

### Refactoring Workflow

Managing refactoring commits:

```bash
# 1. Small, focused refactors
git add src/utils/parser.js
gscribe -t refactor -s utils -o
# Result: refactor(utils): simplify parser logic

# 2. For larger refactors, commit by file
for file in src/models/*.js; do
  git add "$file"
  gscribe -t refactor -s models -y
done

# 3. Final integration commit
git add tests/integration/
gscribe -t test -s integration
```

---

## Team Workflows

### Standardized Team Setup

Create consistent configuration across team:

```bash
# 1. Create .gscribe-config.yaml in project root
cat > .gscribe-config.yaml <<EOF
commit:
  emoji: true
  confirm: true  # Always review in team setting
  description_length: 72
  max_files: 30

ai:
  temperature: 0.5  # More consistent output

types:
  feat: "New feature"
  fix: "Bug fix"
  docs: "Documentation"
  refactor: "Code refactoring"
  test: "Tests"
  chore: "Maintenance"
EOF

# 2. Add to git (without API keys)
git add .gscribe-config.yaml

# 3. Document in README
cat >> README.md <<EOF

## Commit Conventions

We use git-commit-scribe for consistent commit messages:

1. Install: \`go install github.com/tenxprotocols/git-commit-scribe/cmd/gscribe@latest\`
2. Configure: \`gscribe config init\`
3. Commit: \`git add .; gscribe\`

Project configuration is in \`.gscribe-config.yaml\`.
EOF
```

### Pull Request Workflow

Workflow optimized for PR-based development:

```bash
# 1. Create feature branch
git checkout -b feat/new-dashboard

# 2. Make incremental commits
git add src/components/Dashboard.js
gscribe -t feat -s ui -e
# ✨ feat(ui): add dashboard component

git add src/styles/dashboard.css
gscribe -t style -s ui
# style(ui): add dashboard styling

git add tests/Dashboard.test.js
gscribe -t test -s ui

# 3. Before PR, squash if needed
git rebase -i main

# 4. Push and create PR
git push -u origin feat/new-dashboard
gh pr create
```

### Code Review Workflow

Using gscribe during code review:

```bash
# Reviewer checks out PR
gh pr checkout 123

# Suggests changes
vim src/auth/login.js

# Creates suggestion commit
git add src/auth/login.js
gscribe -t refactor -s auth --dry-run
# Review message, then commit
gscribe -t refactor -s auth

# Push suggestion
git push
```

### Release Workflow

Managing releases with conventional commits:

```bash
# 1. Update version
vim package.json
git add package.json
gscribe -t chore -s release

# 2. Update changelog
vim CHANGELOG.md
git add CHANGELOG.md
gscribe -t docs -s changelog

# 3. Create release commit
git add .
gscribe -t chore -s release -y
# chore(release): prepare version 2.0.0

# 4. Tag release
git tag -a v2.0.0 -m "Release version 2.0.0"
git push origin main --tags
```

---

## CI/CD Integration

### GitHub Actions

Using gscribe in GitHub Actions:

```yaml
# .github/workflows/commit-check.yml
name: Validate Commits

on:
  pull_request:
    types: [opened, synchronize]

jobs:
  validate-commits:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Install gscribe
        run: go install github.com/tenxprotocols/git-commit-scribe/cmd/gscribe@latest
      
      - name: Validate commit messages
        run: |
          # Check that commits follow conventional format
          for commit in $(git rev-list origin/main..HEAD); do
            msg=$(git log --format=%s -n 1 $commit)
            if ! echo "$msg" | grep -qE '^(feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert)(\(.+\))?: .+'; then
              echo "Invalid commit message: $msg"
              exit 1
            fi
          done
```

### Auto-commit in CI

Automatically commit changes in CI:

```yaml
# .github/workflows/auto-format.yml
name: Auto Format

on:
  pull_request:
    types: [opened, synchronize]

jobs:
  format:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Install gscribe
        run: go install github.com/tenxprotocols/git-commit-scribe/cmd/gscribe@latest
      
      - name: Format code
        run: |
          go fmt ./...
          
      - name: Commit if changed
        env:
          GSCRIBE_API_KEY: ${{ secrets.OPENROUTER_API_KEY }}
        run: |
          git config user.name "github-actions[bot]"
          git config user.email "github-actions[bot]@users.noreply.github.com"
          
          if [[ -n $(git status -s) ]]; then
            git add .
            gscribe -t style -s format -y
            git push
          fi
```

### GitLab CI

Using gscribe in GitLab CI:

```yaml
# .gitlab-ci.yml
stages:
  - validate
  - commit

validate-commits:
  stage: validate
  image: golang:1.21
  script:
    - go install github.com/tenxprotocols/git-commit-scribe/cmd/gscribe@latest
    - |
      for commit in $(git rev-list origin/main..HEAD); do
        msg=$(git log --format=%s -n 1 $commit)
        if ! echo "$msg" | grep -qE '^(feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert)(\(.+\))?: .+'; then
          echo "Invalid commit message: $msg"
          exit 1
        fi
      done
  only:
    - merge_requests

auto-format:
  stage: commit
  image: golang:1.21
  variables:
    GSCRIBE_API_KEY: $OPENROUTER_API_KEY
  script:
    - go install github.com/tenxprotocols/git-commit-scribe/cmd/gscribe@latest
    - go fmt ./...
    - |
      if [[ -n $(git status -s) ]]; then
        git config user.name "gitlab-ci"
        git config user.email "gitlab-ci@example.com"
        git add .
        gscribe -t style -s format -y
        git push
      fi
  only:
    - merge_requests
```

---

## Monorepo Workflows

### Multi-Package Commits

Managing commits in monorepos:

```bash
# Directory structure:
# packages/
#   frontend/
#   backend/
#   shared/

# 1. Commit per package
cd packages/frontend
git add .
gscribe -t feat -s frontend -y

cd ../backend
git add .
gscribe -t feat -s backend -y

# 2. Or use scopes to indicate packages
git add packages/frontend/
gscribe -t feat -s "frontend/auth"

git add packages/backend/
gscribe -t feat -s "backend/auth"

# 3. Shared libraries
git add packages/shared/
gscribe -t feat -s shared
```

### Workspace-Specific Configuration

Configure per workspace:

```bash
# Root config
cat > .gscribe-config.yaml <<EOF
commit:
  emoji: true
  max_files: 100  # Larger for monorepo
EOF

# Frontend-specific (optional)
cat > packages/frontend/.git/gscribe-config.yaml <<EOF
commit:
  description_length: 80
types:
  feat: "New UI feature"
  fix: "UI bug fix"
  style: "Styling changes"
EOF
```

---

## Advanced Workflows

### Custom Prompt Workflow

Using custom prompts for specific projects:

```bash
# 1. Create custom prompt
cat > .gscribe/prompts/detailed.txt <<EOF
Generate a detailed git commit message following these rules:
- Use conventional commits format
- Include technical details
- Reference issue numbers if present
- Maximum 100 characters for description
- Include implementation notes in body
EOF

# 2. Configure to use custom prompt
gscribe config set ai.custom_prompt_file ".gscribe/prompts/detailed.txt" --repo

# 3. Use normally
git add .
gscribe
```

### Interactive Staging Workflow

Commit parts of files:

```bash
# 1. Interactive staging
git add -p src/component.js
# Select hunks to stage

# 2. Generate commit for staged changes
gscribe -t feat -s component

# 3. Stage more changes
git add -p src/component.js
gscribe -t refactor -s component

# Result: Multiple focused commits from one file
```

### Stash and Commit Workflow

Managing work-in-progress:

```bash
# 1. Stash current work
git stash save "WIP: new feature"

# 2. Fix urgent bug
git checkout main
git checkout -b fix/urgent-bug
vim src/critical.js
git add .
gscribe -t fix -s critical -yp

# 3. Return to feature work
git checkout feat/new-feature
git stash pop
```

### Changelog Generation Workflow

Automated changelog from commits:

```bash
# 1. Commit with conventional format
git add .
gscribe -y

# 2. Generate changelog (using conventional-changelog)
npx conventional-changelog -p angular -i CHANGELOG.md -s

# 3. Commit changelog
git add CHANGELOG.md
gscribe -t docs -s changelog -y
```

---

## Best Practices

### Commit Granularity

**DO:**
```bash
# Small, focused commits
git add src/auth/login.js
gscribe -t feat -s auth
# feat(auth): implement login form

git add src/auth/validation.js
gscribe -t feat -s auth
# feat(auth): add input validation

git add tests/auth/login.test.js
gscribe -t test -s auth
# test(auth): add login form tests
```

**DON'T:**
```bash
# Large, unfocused commits
git add .
gscribe
# feat: add auth, update styles, fix tests, refactor utils
```

### Meaningful Scopes

**DO:**
```bash
gscribe -t feat -s auth      # Clear component
gscribe -t fix -s api        # Clear layer
gscribe -t refactor -s utils # Clear module
```

**DON'T:**
```bash
gscribe -t feat -s stuff     # Too vague
gscribe -t fix -s file       # Not meaningful
```

### Breaking Changes

**DO:**
```bash
# Explicit breaking change
gscribe -t feat -s api --breaking
# Result includes: BREAKING CHANGE: ...

# Alternative
git add .
gscribe -t feat -s api
# Review message, add BREAKING CHANGE manually if needed
```

**DON'T:**
```bash
# Hide breaking changes in regular commits
gscribe -t refactor -s api
# (when it actually breaks the API)
```

### Review Before Committing

**DO:**
```bash
# Always review (keep confirm: true)
gscribe config set commit.confirm true

# Use dry run to preview
gscribe -d

# Review staged changes
git diff --staged
```

**DON'T:**
```bash
# Blind auto-commit everything
gscribe config set commit.confirm false
git add .; gscribe  # Without checking what's staged
```

### Cache Management

**DO:**
```bash
# Enable caching for performance
gscribe config set cache.enabled true

# Clear cache periodically
gscribe cache clear --all

# Check cache stats
gscribe cache stats
```

**DON'T:**
```bash
# Use cache for sensitive repos without understanding
# (cache stores diff content)
```

### API Key Security

**DO:**
```bash
# Use environment variables in CI
export GSCRIBE_API_KEY="${SECRET_API_KEY}"

# Use config file with proper permissions
chmod 600 ~/.config/gscribe/config.yaml

# Add to .gitignore
echo ".gscribe-config.yaml" >> .gitignore  # If it contains secrets
```

**DON'T:**
```bash
# Commit API keys
git add .gscribe-config.yaml  # (if it contains your API key)

# Share config files with API keys
```

---

## Team Conventions

### Establish Team Standards

Document team conventions in project README:

```markdown
## Commit Conventions

### Installation
```bash
go install github.com/tenxprotocols/git-commit-scribe/cmd/gscribe@latest
gscribe config init
```

### Usage
- Always use gscribe for commits
- Review generated messages before confirming
- Use appropriate scopes: `auth`, `api`, `ui`, `db`
- Mark breaking changes with `--breaking` flag

### Commit Types
- `feat`: New features
- `fix`: Bug fixes
- `docs`: Documentation
- `test`: Tests
- `refactor`: Code refactoring
- `chore`: Maintenance

### Examples
```bash
gscribe -t feat -s auth -e  # New auth feature with emoji
gscribe -t fix -s api       # API bug fix
gscribe -t docs            # Documentation update
```
```

### Pre-commit Hooks

Enforce conventions with git hooks:

```bash
# .git/hooks/pre-commit
#!/bin/bash

# Ensure commit message follows conventional format
# (This example shows how you might validate)
# The actual commit message validation happens post-commit

# Example: Ensure gscribe is available
if ! command -v gscribe &> /dev/null; then
    echo "Error: gscribe is not installed"
    echo "Install: go install github.com/tenxprotocol/git-commit-scribe/cmd/gscribe@latest"
    exit 1
fi

# Run linting/formatting before commit
npm run lint
go fmt ./...
```

### Team Templates

Share prompt templates:

```bash
# Create team prompts directory
mkdir -p .gscribe/prompts

# Create team prompt
cat > .gscribe/prompts/team.txt <<EOF
Generate commit message following our team conventions:
- Use present tense ("add feature" not "added feature")
- Reference Jira ticket in footer if applicable
- Keep description under 72 characters
- Include technical details in body for complex changes
EOF

# Commit template to repo
git add .gscribe/
gscribe -t chore -s config -m "Add team gscribe prompt template"
```

### Onboarding Checklist

Create onboarding documentation:

```markdown
## New Developer Setup

### git-commit-scribe Setup

1. Install gscribe:
   ```bash
   go install github.com/tenxprotocols/git-commit-scribe/cmd/gscribe@latest
   ```

2. Get OpenRouter API key:
   - Visit https://openrouter.ai/
   - Create account and generate API key

3. Configure gscribe:
   ```bash
   gscribe config init
   # Enter your API key when prompted
   # Select: anthropic/claude-3.5-sonnet
   ```

4. Test it:
   ```bash
   echo "test" > test.txt
   git add test.txt
   gscribe -d  # Dry run
   git reset HEAD test.txt
   rm test.txt
   ```

5. Set team preferences:
   ```bash
   gscribe config set commit.emoji true
   ```

Your first commit should look like:
```
✨ feat(onboarding): add [your name] to team
```
```

---

## Summary

**Key Takeaways:**

1. **Start Simple**: Begin with basic `git add . && gscribe` workflow
2. **Iterate**: Gradually adopt more advanced workflows as needed
3. **Be Consistent**: Use scopes and types consistently across team
4. **Review**: Always review generated messages before committing
5. **Configure**: Set up team-wide configuration for consistency
6. **Automate**: Integrate into CI/CD for validation
7. **Document**: Maintain clear team conventions and examples

**Next Steps:**

- Review [Usage Guide](usage.md) for detailed examples
- Check [Troubleshooting](troubleshooting.md) if you encounter issues
- See [Configuration](configuration.md) for customization options
- Read [CLI Reference](cli-reference.md) for all available commands
