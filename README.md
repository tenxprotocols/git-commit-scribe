# git-commit-scribe

AI-powered Git commit message generator using Conventional Commits.

## Installation

```bash
go install github.com/tenxprotocols/git-commit-scribe/cmd/gscribe@latest
```

## Quick Start

1. Get an API key from [OpenRouter](https://openrouter.ai/)
2. Initialize: `gscribe config init`
3. Use: `git add . && gscribe`

## Usage

```bash
gscribe              # Generate commit
gscribe -d           # Dry-run/Preview only
gscribe -y           # Yes (skip confirmation)
gscribe -yp          # Generate commit and push
```

## Documentation

- [Guide](docs/GUIDE.md) - Setup, usage, workflows
- [Reference](docs/REFERENCE.md) - CLI flags, config options
- [Troubleshooting](docs/troubleshooting.md) - Common issues

## Features

- AI-powered commit message generation
- Includes Conventional Commits v1.0.0 support
- Highly configurable with smart caching for speed

## License

MIT
