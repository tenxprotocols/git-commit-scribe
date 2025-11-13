# CI/CD Documentation

This document describes the continuous integration and deployment workflows for git-commit-scribe.

## Overview

The project uses GitHub Actions for automated testing, building, and releases with the following workflows:

- **CI (Continuous Integration)**: Automated testing, linting, and security scanning
- **Build**: Multi-platform binary builds with artifact uploads
- **Release**: Automated versioned releases with changelog generation

## Workflows

### CI Workflow (`.github/workflows/ci.yml`)

**Triggers:**
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop` branches

**Jobs:**

#### 1. Test Matrix
- **Platforms**: Ubuntu, macOS
- **Go Versions**: 1.22, 1.23
- **Steps**:
  - Checkout code
  - Set up Go with caching
  - Download and verify dependencies
  - Run `go vet` for static analysis
  - Check code formatting with `gofmt`
  - Run tests with race detection and coverage
  - Upload coverage to Codecov (Ubuntu + Go 1.23 only)

#### 2. Lint
- Runs `golangci-lint` with comprehensive linter suite
- Ensures code quality and consistency
- Configuration in `.golangci.yml`

#### 3. Security Scan
- Runs Gosec security scanner
- Uploads SARIF results for GitHub Security tab
- Checks for common security issues

**Status**: [![CI](https://github.com/TenXProtocols/git-commit-scribe/actions/workflows/ci.yml/badge.svg)](https://github.com/TenXProtocols/git-commit-scribe/actions/workflows/ci.yml)

### Build Workflow (`.github/workflows/build.yml`)

**Triggers:**
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop` branches

**Build Matrix:**
- **Linux**: amd64, arm64
- **macOS**: amd64 (Intel), arm64 (Apple Silicon)

**Features:**
- Version information embedding (version, commit, build date)
- Binary artifact uploads (30-day retention)
- Basic verification tests for each platform
- Build summary job to ensure all builds succeeded

**Artifacts:**
- `gscribe-linux-amd64`
- `gscribe-linux-arm64`
- `gscribe-darwin-amd64`
- `gscribe-darwin-arm64`

**Status**: [![Build](https://github.com/TenXProtocols/git-commit-scribe/actions/workflows/build.yml/badge.svg)](https://github.com/TenXProtocols/git-commit-scribe/actions/workflows/build.yml)

### Release Workflow (`.github/workflows/release.yml`)

**Triggers:**
- Push of version tags matching pattern `v*.*.*` (e.g., `v1.0.0`)

**Jobs:**

#### 1. Pre-Release Testing
- Runs full test suite with race detection
- Runs `go vet` for static analysis
- Must pass before release creation

#### 2. Release Creation (GoReleaser)
- Builds binaries for all platforms
- Creates GitHub release with automated changelog
- Uploads binary artifacts and checksums
- Updates Homebrew tap (if configured)

**Changelog Groups** (in `.goreleaser.yaml`):
- Features (commits starting with `feat:`)
- Bug Fixes (commits starting with `fix:`)
- Performance Improvements (commits starting with `perf:`)
- Refactorings (commits starting with `refactor:`)
- Documentation (commits starting with `docs:`)
- Other Changes

**Excluded from changelog:**
- Test commits (`test:`)
- Chore commits (`chore:`)
- CI commits (`ci:`)
- Build commits (`build:`)
- Style commits (`style:`)
- Merge commits

#### 3. Installation Verification
- Downloads released binaries for Ubuntu and macOS
- Verifies binary execution (`--version`, `--help`)
- Tests basic functionality with git repositories
- Ensures binaries work correctly on each platform

**Status**: [![Release](https://github.com/TenXProtocols/git-commit-scribe/actions/workflows/release.yml/badge.svg)](https://github.com/TenXProtocols/git-commit-scribe/actions/workflows/release.yml)

## Configuration Files

### `.goreleaser.yaml`

GoReleaser configuration for building and releasing:

```yaml
builds:
  - Builds for linux, darwin
  - Architectures: amd64, arm64
  - Embeds version, commit, and build date
  - CGO disabled for static binaries

archives:
  - tar.gz for Unix systems
  - Includes LICENSE and README

changelog:
  - Automatically groups commits by type
  - Filters out non-relevant commits
  - Uses GitHub commit information

brews:
  - Homebrew tap integration
  - Automatic formula updates
```

### `.golangci.yml`

Linter configuration enabling:
- Static analysis (staticcheck, govet)
- Code simplification (gosimple)
- Security scanning (gosec)
- Style checking (gofmt, goimports, revive)
- Bug detection (errcheck, ineffassign)
- Performance optimization hints (gocritic)
- Code duplication detection (dupl)
- Cyclomatic complexity (gocyclo)

## Creating a Release

### Step 1: Prepare the Release
```bash
# Ensure you're on main branch with latest changes
git checkout main
git pull origin main

# Run tests locally
go test ./...
go vet ./...
```

### Step 2: Create and Push Tag
```bash
# Create a new version tag (follow semantic versioning)
git tag -a v1.0.0 -m "Release v1.0.0"

# Push the tag to trigger release workflow
git push origin v1.0.0
```

### Step 3: Monitor Release
1. Go to GitHub Actions tab
2. Watch the Release workflow execution
3. Verify all jobs complete successfully:
   - Pre-release tests pass
   - GoReleaser creates release
   - Installation verification succeeds on all platforms

### Step 4: Verify Release
1. Check the Releases page for the new release
2. Verify changelog is properly generated
3. Download and test binaries for your platform
4. Verify Homebrew tap is updated (if configured)

## Secrets Configuration

The following GitHub secrets should be configured:

- `GITHUB_TOKEN`: Automatically provided by GitHub Actions
- `HOMEBREW_TAP_GITHUB_TOKEN`: Personal access token for updating Homebrew tap (optional)
  - Required permissions: `repo` scope
  - Create at: https://github.com/settings/tokens

## Monitoring

### Viewing Workflow Runs
- Go to the **Actions** tab in the GitHub repository
- Select the workflow you want to view
- Click on a specific run to see detailed logs

### Downloading Artifacts
- Go to a completed workflow run
- Scroll to the **Artifacts** section
- Download the binary for your platform

### Checking Coverage
- Coverage reports are uploaded to Codecov
- View at: https://codecov.io/gh/TenXProtocols/git-commit-scribe

## Troubleshooting

### Build Fails on Specific Platform
1. Check the workflow logs for error messages
2. Test the build locally:
   ```bash
   GOOS=linux GOARCH=amd64 go build -o gscribe ./cmd/gscribe
   ```
3. Ensure all dependencies are properly declared in `go.mod`

### Tests Fail in CI
1. Run tests locally with race detection:
   ```bash
   go test -race ./...
   ```
2. Check for platform-specific issues
3. Review test logs in the workflow run

### Release Creation Fails
1. Verify the tag format matches `v*.*.*`
2. Check GoReleaser configuration syntax
3. Ensure all required secrets are configured
4. Review the workflow logs for specific errors

### Linting Errors
1. Run golangci-lint locally:
   ```bash
   golangci-lint run
   ```
2. Fix issues or adjust `.golangci.yml` configuration
3. Run `gofmt -s -w .` to auto-format code

## Best Practices

### Commit Messages
Use conventional commit format for better changelog generation:
- `feat: add new feature` - New features
- `fix: resolve bug` - Bug fixes
- `perf: improve performance` - Performance improvements
- `refactor: restructure code` - Code refactoring
- `docs: update documentation` - Documentation changes
- `test: add tests` - Test additions
- `chore: update dependencies` - Maintenance tasks

### Version Numbering
Follow semantic versioning (SemVer):
- **MAJOR** version (v2.0.0): Incompatible API changes
- **MINOR** version (v1.1.0): Backwards-compatible functionality
- **PATCH** version (v1.0.1): Backwards-compatible bug fixes

### Testing Before Release
Always run the full test suite locally before creating a release:
```bash
go test -race -coverprofile=coverage.out ./...
go vet ./...
golangci-lint run
```

## Badges for README

Add these badges to your README.md to display workflow status:

```markdown
[![CI](https://github.com/TenXProtocols/git-commit-scribe/actions/workflows/ci.yml/badge.svg)](https://github.com/TenXProtocols/git-commit-scribe/actions/workflows/ci.yml)
[![Build](https://github.com/TenXProtocols/git-commit-scribe/actions/workflows/build.yml/badge.svg)](https://github.com/TenXProtocols/git-commit-scribe/actions/workflows/build.yml)
[![Release](https://github.com/TenXProtocols/git-commit-scribe/actions/workflows/release.yml/badge.svg)](https://github.com/TenXProtocols/git-commit-scribe/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/TenXProtocols/git-commit-scribe)](https://goreportcard.com/report/github.com/TenXProtocols/git-commit-scribe)
