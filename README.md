# gitcomm

A CLI tool that automates git commit message creation following the [Conventional Commits](https://www.conventionalcommits.org/) specification. It supports both manual input and AI-assisted generation.

## Features

- ✅ **Manual Commit Messages**: Interactive prompts for creating Conventional Commits compliant messages
- ✅ **AI-Assisted Generation**: Support for OpenAI, Anthropic, Mistral, and local models (using official SDKs)
- ✅ **Unified AI Prompts**: All AI providers use identical prompts with validation rules extracted dynamically from the validator, ensuring consistent commit message quality
- ✅ **AI Message Acceptance Options**: When an AI-generated message is displayed, choose from three options:
  - **Accept and commit directly**: Commit immediately with the AI message (fastest path)
  - **Accept and edit**: Pre-fill all commit message fields with AI values for quick editing
  - **Reject**: Generate a new AI message or proceed with manual input
- ✅ **Token Calculation**: Estimate AI token usage before generating messages
- ✅ **Diff Computation**: Automatically computes unified diffs for staged files to provide AI models with actual code changes (optimized for token usage with 0 context lines and 5000 character limit)
- ✅ **Format Validation**: Automatic validation against Conventional Commits specification
- ✅ **Auto-Staging**: Automatically stage modified files on launch (or all files with `-a` flag)
- ✅ **State Restoration**: Automatically restore staging state if you cancel or exit without committing
- ✅ **Signal Handling**: Graceful interruption handling (Ctrl+C) with state restoration and timeout protection (exits within 5 seconds)
- ✅ **CLI Options**: Auto-stage files (`-a`), disable signoff (`-s`), disable signing (`--no-sign`), provider selection, debug logging (`-d`, `--debug`)
- ✅ **Git Config Integration**: Automatically uses `user.name` and `user.email` from git configuration for commit author
- ✅ **SSH Commit Signing**: Automatically signs commits with SSH keys when configured in git config (`gpg.format = ssh`, `user.signingkey`)
- ✅ **Error Handling**: Graceful fallback to manual input on AI failures
- ✅ **Debug Logging**: Optional debug mode with raw text format (no timestamps) for troubleshooting
- ✅ **Interactive UI**: Interactive select list for commit type selection with visual feedback (checkmarks, highlighting)
- ✅ **Type Selection Confirmation**: Displays confirmation line after commit type selection showing the chosen type
- ✅ **Multiline Input**: Support for multiline body and footer input with double-Enter completion

## Installation

```bash
# Build from source
git clone https://github.com/golgoth31/gitcomm.git
cd gitcomm
go build -o gitcomm ./cmd/gitcomm

# Install globally
go install ./cmd/gitcomm
```

## Quick Start

### Basic Usage

```bash
# Navigate to your git repository
cd /path/to/your/repo

# Run gitcomm
gitcomm

# Follow the interactive prompts to create a commit message
```

### Auto-Staging Behavior

```bash
# Modified files are automatically staged on launch
gitcomm

# With -a flag: also stage untracked files
gitcomm -a

# If you cancel or exit without committing, staging state is automatically restored
# (Press Ctrl+C or reject the commit message to see restoration in action)
```

### Without Signoff

```bash
# Create commit without "Signed-off-by" line
gitcomm -s
```

### Git Configuration

GitComm automatically reads your git configuration (`.git/config` and `~/.gitconfig`) to:
- Use `user.name` and `user.email` for commit author
- Sign commits with SSH keys when configured

**Commit Author**: GitComm uses `user.name` and `user.email` from your git config. If not configured, defaults to "gitcomm <gitcomm@local>".

**SSH Commit Signing**: If your git config has:
```ini
[user]
    signingkey = ~/.ssh/id_ed25519.pub
[gpg]
    format = ssh
[commit]
    gpgsign = true
```

GitComm will automatically sign commits with your SSH key. Use `--no-sign` to disable signing for a specific commit.

```bash
# Disable commit signing for this commit
gitcomm --no-sign
```

## AI Configuration

GitComm uses official Go SDKs for AI providers:
- **OpenAI**: [github.com/openai/openai-go/v3](https://github.com/openai/openai-go) (SDK v3, Responses API)
- **Anthropic**: [github.com/anthropics/anthropic-sdk-go](https://github.com/anthropics/anthropic-sdk-go)
- **Mistral**: [github.com/gage-technologies/mistral-go](https://github.com/gage-technologies/mistral-go)

1. Configuration file is automatically created at `~/.gitcomm/config.yaml` when you first run GitComm. If the file doesn't exist, it will be created as an empty file with secure permissions (0600). You can also specify a custom path using the `--config` flag.

   The config file will be created automatically with the following structure:

```yaml
ai:
  default_provider: openai
  providers:
    openai:
      api_key: ${OPENAI_API_KEY}
      model: gpt-4
    anthropic:
      api_key: ${ANTHROPIC_API_KEY}
      model: claude-3-opus
    mistral:
      api_key: ${MISTRAL_API_KEY}
      model: mistral-large-latest
```

   **Note**: The config file and parent directories (`~/.gitcomm/`) are automatically created if they don't exist. The file is created with restrictive permissions (0600) to protect your API keys.

2. Set environment variables:

```bash
export OPENAI_API_KEY="your-api-key-here"
# or
export ANTHROPIC_API_KEY="your-api-key-here"
# or
export MISTRAL_API_KEY="your-api-key-here"
```

3. Use AI generation:

```bash
# Use default provider
gitcomm

# Use specific provider
gitcomm --provider openai

# Skip AI and use manual input
gitcomm --skip-ai
```

### AI Message Acceptance Options

When GitComm displays an AI-generated commit message, you'll see three options:

```
--- AI Generated Message ---
feat(auth): add user authentication

Implement JWT-based authentication with refresh tokens.
Add login and logout endpoints.

Closes #123
---

Options:
  1. Accept and commit directly
  2. Accept and edit
  3. Reject
Choose option (1/2/3):
```

- **Option 1 - Accept and commit directly**: Creates the commit immediately with the AI message. Fastest path for messages you're satisfied with.
- **Option 2 - Accept and edit**: Pre-fills all commit message fields (type, scope, subject, body, footer) with AI values. You can then modify any field before committing.
- **Option 3 - Reject**: Choose to generate a new AI message or proceed with manual input (empty fields).

**Pre-filling**: When you choose "Accept and edit", all fields are pre-filled:
- Commit type is automatically selected in the interactive list (if it matches)
- Scope, subject, body, and footer are pre-populated with AI values
- You can modify any field or accept the defaults by pressing Enter

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

## CLI Options

- `-a, --add-all`: Automatically stage all files (modified + untracked). Without this flag, only modified files are auto-staged
- `-s, --no-signoff`: Disable commit signoff (omit Signed-off-by line)
- `--no-sign`: Disable commit signing (overrides git config `commit.gpgsign` setting)
- `--provider <name>`: Override default AI provider (openai, anthropic, mistral, local)
- `--skip-ai`: Skip AI generation and proceed directly to manual input
- `--config <path>`: Path to configuration file (default: ~/.gitcomm/config.yaml)
- `-d, --debug`: Enable debug logging (raw text format, no timestamps). When enabled, all DEBUG-level log messages are displayed. By default, the CLI runs silently with no log output.
- `-v, --verbose`: Verbose flag (no-op when debug flag is not set). Debug flag takes precedence.
- `-h, --help`: Display help information

## Auto-Staging and State Restoration

**Auto-Staging**: When you run `gitcomm`, all modified files are automatically staged before any prompts are shown. This ensures AI analysis has access to all changes. Use the `-a` flag to also include untracked files.

**State Restoration**: If you cancel the CLI (Ctrl+C), reject the commit message, or encounter an error, the staging state is automatically restored to what it was before you ran `gitcomm`. This prevents accidental staging of files you didn't intend to commit.

**Timeout Protection**: When you press Ctrl+C, the CLI will restore the staging state and exit within 5 seconds. If restoration takes longer than 3 seconds, it will timeout and exit immediately with a warning message, ensuring the CLI never hangs indefinitely.

**Safety**: Files that were already staged before running `gitcomm` are preserved - only files staged by the CLI are restored.

## Requirements

- Go 1.25.0 or later
- Git installed and configured
- (Optional) AI provider API keys for AI-assisted generation

## Development

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Build
make build

# Lint
make lint

# Format code
make format
```

## License

MIT
