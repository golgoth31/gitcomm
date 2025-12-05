# CLI Contract: Git Commit Message Automation

**Date**: 2025-01-27
**Feature**: 001-git-commit-cli

## Command Interface

### Command: `gitcomm`

CLI tool for automated git commit message creation.

**Usage**:
```bash
gitcomm [flags]
```

**Flags**:
- `-a, --add-all`: Automatically stage all unstaged files before proceeding
- `-s, --no-signoff`: Disable commit signoff (omit Signed-off-by line)
- `--provider <name>`: Override default AI provider (openai, anthropic, local)
- `--skip-ai`: Skip AI generation and proceed directly to manual input
- `--config <path>`: Path to configuration file (default: ~/.gitcomm/config.yaml)
- `-v, --verbose`: Enable verbose logging
- `-h, --help`: Display help information

**Exit Codes**:
- `0`: Success - commit created
- `1`: Error - general error (invalid input, git error, etc.)
- `2`: Configuration error - invalid config file or missing credentials
- `3`: AI provider error - provider unavailable or authentication failed

---

## Interactive Workflow Contract

### Step 1: Repository State Analysis

**Input**: None (reads from current git repository)

**Output**:
- Repository state summary (staged/unstaged files)
- Token count estimate (if AI available)

**Error Cases**:
- Not in git repository → Exit with code 1, message: "Error: not a git repository"
- No changes detected → Prompt for empty commit confirmation

---

### Step 2: AI Usage Decision

**Input**: User choice (y/n)

**Output**:
- Proceed to AI generation, OR
- Proceed to manual input

**Error Cases**:
- Invalid input → Re-prompt

---

### Step 3: AI Generation (if selected)

**Input**: Repository state

**Output**:
- AI-generated commit message, OR
- Error message with fallback to manual input

**Error Cases**:
- Provider unavailable → Fallback to manual with error message
- Network timeout → Fallback to manual with error message
- Invalid response → Fallback to manual with error message

---

### Step 4: Message Validation (AI or Manual)

**Input**: Commit message components or AI-generated message

**Output**:
- Validated message for review, OR
- Validation errors with option to edit

**Error Cases**:
- Invalid format → Display errors, offer edit/use-with-warning options
- Empty subject → Re-prompt for subject

---

### Step 5: Manual Input (if needed)

**Input Sequence**:
1. Scope (optional, can be empty)
2. Subject (required, non-empty)
3. Body (optional, can be empty)
4. Footer (optional, can be empty)

**Output**: Complete commit message

**Validation**:
- Subject must not be empty
- Subject must be ≤72 characters (warn if longer)
- Body must be ≤320 characters (warn if longer)

---

### Step 6: Final Review

**Input**: Formatted commit message

**Output**:
- User accepts → Create commit
- User rejects → Return to editing

**Error Cases**:
- Git commit fails → Display error, exit with code 1

---

## Configuration File Contract

**Location**: `~/.gitcomm/config.yaml` (or path specified via `--config`)

**Schema**:
```yaml
ai:
  default_provider: string  # openai, anthropic, or local
  providers:
    openai:
      api_key: string        # Required, can use ${OPENAI_API_KEY}
      model: string          # Optional, default: gpt-4
      timeout: duration      # Optional, default: 30s
    anthropic:
      api_key: string        # Required, can use ${ANTHROPIC_API_KEY}
      model: string          # Optional, default: claude-3-opus
      timeout: duration      # Optional, default: 30s
    local:
      endpoint: string       # Required for local models
      api_key: string        # Optional
      timeout: duration      # Optional, default: 30s
```

**Validation**:
- YAML must be valid
- `default_provider` must be one of configured providers
- Each provider must have required fields
- Environment variable substitution supported: `${VAR_NAME}`

---

## AI Provider Interface Contract

### Method: GenerateCommitMessage

**Input**:
```go
type GenerateRequest struct {
    RepositoryState RepositoryState
    Model          string  // Optional, override default
}
```

**Output**:
```go
type GenerateResponse struct {
    Message string
    Tokens  TokenUsage
}
```

**Error Cases**:
- Authentication failure → `ErrAuthenticationFailed`
- Network error → `ErrNetworkError`
- Timeout → `ErrTimeout`
- Rate limit → `ErrRateLimit`
- Invalid response → `ErrInvalidResponse`

**Timeout**: 30 seconds (configurable)

---

## Git Operations Contract

### Method: GetRepositoryState

**Input**: None (uses current working directory)

**Output**: `RepositoryState` with staged and unstaged changes

**Error Cases**:
- Not a git repository → `ErrNotGitRepository`
- Git command failure → `ErrGitOperationFailed`

---

### Method: StageAllFiles

**Input**: None

**Output**: Success/failure

**Error Cases**:
- Git command failure → `ErrGitOperationFailed`

---

### Method: CreateCommit

**Input**:
```go
type CommitRequest struct {
    Message  CommitMessage
    Signoff  bool
}
```

**Output**: Commit SHA

**Error Cases**:
- No changes staged → `ErrNoStagedChanges`
- Git command failure → `ErrGitOperationFailed`
- Invalid message format → `ErrInvalidFormat`

---

## Validation Contract

### Method: ValidateCommitMessage

**Input**: `CommitMessage`

**Output**:
```go
type ValidationResult struct {
    Valid   bool
    Errors  []ValidationError
}
```

**Validation Rules**:
- Type must be valid Conventional Commits type
- Subject must be non-empty
- Subject must be ≤72 characters
- Body must be ≤320 characters if provided
- Scope must be valid identifier if provided

**Error Types**:
- `ErrInvalidType`: Type not in allowed list
- `ErrEmptySubject`: Subject is empty
- `ErrSubjectTooLong`: Subject > 72 characters
- `ErrBodyTooLong`: Body > 320 characters
- `ErrInvalidScope`: Scope contains invalid characters

---

## Logging Contract

**Format**: Structured JSON via zerolog

**Log Levels**:
- `DEBUG`: Detailed execution flow (only with `-v` flag)
- `INFO`: Normal operations (repository state, commit created)
- `WARN`: Recoverable issues (format warnings, fallbacks)
- `ERROR`: Errors requiring attention (provider failures, git errors)

**Fields**:
- `timestamp`: ISO 8601 timestamp
- `level`: Log level
- `message`: Human-readable message
- `error`: Error details (if applicable)
- `component`: Component name (service, repository, ai, etc.)

**Secrets**: Never log API keys, tokens, or sensitive data
