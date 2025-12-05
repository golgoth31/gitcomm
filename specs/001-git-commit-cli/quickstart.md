# Quick Start: Git Commit Message Automation CLI

**Feature**: 001-git-commit-cli
**Date**: 2025-01-27

## Prerequisites

- Go 1.25.0 or later
- Git installed and configured
- (Optional) AI provider API keys (OpenAI, Anthropic, or local model endpoint)

## Installation

```bash
# Clone the repository
git clone https://github.com/golgoth31/gitcomm.git
cd gitcomm

# Build the CLI
go build -o gitcomm cmd/gitcomm/main.go

# Install globally (optional)
go install ./cmd/gitcomm
```

## Basic Usage

### Manual Commit Message Creation

```bash
# Navigate to your git repository
cd /path/to/your/repo

# Run the CLI
gitcomm

# Follow the interactive prompts:
# 1. Decline AI assistance (or accept if configured)
# 2. Enter scope (optional, press Enter to skip)
# 3. Enter subject (required)
# 4. Enter body (optional, press Enter to skip)
# 5. Enter footer (optional, press Enter to skip)
# 6. Review and confirm the commit message
```

### With Auto-Staging

```bash
# Automatically stage all unstaged files before committing
gitcomm -a
```

### Without Signoff

```bash
# Create commit without "Signed-off-by" line
gitcomm -s
```

## AI-Assisted Usage

### Setup AI Provider

1. Create configuration file:
```bash
mkdir -p ~/.gitcomm
cat > ~/.gitcomm/config.yaml <<EOF
ai:
  default_provider: openai
  providers:
    openai:
      api_key: ${OPENAI_API_KEY}
      model: gpt-4
    anthropic:
      api_key: ${ANTHROPIC_API_KEY}
      model: claude-3-opus
EOF
```

2. Set environment variables:
```bash
export OPENAI_API_KEY="your-api-key-here"
# or
export ANTHROPIC_API_KEY="your-api-key-here"
```

### Using AI Generation

```bash
# Run CLI (will prompt for AI usage)
gitcomm

# Or skip AI and go directly to manual input
gitcomm --skip-ai

# Or use specific provider
gitcomm --provider openai
```

## Example Workflow

```bash
# 1. Make some changes
echo "New feature" > feature.txt
git add feature.txt

# 2. Run gitcomm
gitcomm

# 3. Interactive session:
#    - AI usage? (y/n): y
#    - Review AI-generated message
#    - Accept or reject
#    - If rejected, enter manual input
#    - Review final message
#    - Confirm commit

# 4. Commit created with Conventional Commits format
```

## Example Commit Messages

The CLI generates commit messages following Conventional Commits format:

```
feat(auth): add user authentication

Implement JWT-based authentication with refresh tokens.
Supports login, logout, and token refresh endpoints.

Closes #123
```

```
fix(api): resolve timeout issue in request handler

The handler was not properly handling context cancellation,
causing requests to hang indefinitely.

Fixes #456
```

## Configuration Options

### Custom Config Location

```bash
gitcomm --config /path/to/custom/config.yaml
```

### Verbose Logging

```bash
gitcomm -v
```

## Troubleshooting

### "Error: not a git repository"

**Solution**: Run `gitcomm` from within a git repository directory.

### "No changes to commit"

**Solution**:
- Stage some files: `git add <files>`
- Or confirm empty commit when prompted

### "AI provider unavailable"

**Solution**:
- Check API key configuration
- Verify network connectivity
- CLI will automatically fallback to manual input

### "Invalid commit message format"

**Solution**:
- Review Conventional Commits specification
- Ensure subject is non-empty and ≤72 characters
- Ensure body is ≤320 characters if provided

## Next Steps

- Read the full [specification](./spec.md) for detailed requirements
- Review the [implementation plan](./plan.md) for technical details
- Check [data model](./data-model.md) for domain entities
- See [CLI contract](./contracts/cli-contract.md) for API details
