# git-commit-scribe

AI-powered Git commit message generator using Conventional Commits.

## Installation

```bash
go install github.com/tenxprotocols/git-commit-scribe/cmd/gscribe@latest
```

Or build from source:

```bash
git clone https://github.com/tenxprotocols/git-commit-scribe
cd git-commit-scribe
./scripts/build.sh
```

## Quick Start

1. Get an API key from [OpenRouter](https://openrouter.ai/)

2. Initialize configuration:
```bash
gscribe config init
```

3. Stage your changes and generate a commit:
```bash
git add .
gscribe
```

## Basic Usage

```bash
# Generate and create commit (interactive)
gscribe

# Dry run (preview message only)
gscribe -d

# Skip confirmation
gscribe -y

# Specify commit type
gscribe -t feat

# Add emoji
gscribe -e

# Push after commit
gscribe -p
```

## Documentation

- [Configuration](docs/configuration.md) - Configure API keys, models, and settings
- [Usage Guide](docs/usage.md) - Detailed usage examples and workflows
- [CLI Reference](docs/cli-reference.md) - Complete command-line options

## Features

✨ AI-powered commit message generation  
📝 Full Conventional Commits v1.0.0 support  
🎨 Optional emoji in commit messages  
⚙️ Highly configurable via CLI flags or config files  
💾 Smart caching for faster repeated operations  
🔒 Secure API key storage  

## License

MIT
