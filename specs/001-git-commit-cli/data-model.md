# Data Model: Git Commit Message Automation CLI

**Date**: 2025-01-27
**Feature**: 001-git-commit-cli

## Domain Entities

### CommitMessage

Represents a structured commit message conforming to Conventional Commits specification.

**Fields**:
- `Type` (string, required): Commit type (feat, fix, docs, style, refactor, test, chore, version)
- `Scope` (string, optional): Scope of the change (e.g., "auth", "api", "cli")
- `Subject` (string, required): Short description in imperative mood, no period, ≤72 characters
- `Body` (string, optional): Detailed explanation, wrapped at 72 chars, ≤320 characters
- `Footer` (string, optional): Footer lines (issue references, breaking changes, etc.)
- `Signoff` (bool): Whether to include "Signed-off-by" line (default: true)

**Validation Rules**:
- Type must be one of the allowed values
- Subject cannot be empty
- Subject must be ≤72 characters
- Body must be ≤320 characters if provided
- Scope must be valid identifier if provided (alphanumeric, hyphens, underscores)

**State Transitions**:
- Created → Validated → Committed
- Created → Rejected (user rejects) → Created (restart)
- Created → Invalid (validation fails) → Created (edit required)

---

### RepositoryState

Represents the current state of the git repository for commit message generation.

**Fields**:
- `StagedFiles` ([]FileChange, required): List of staged file changes
- `UnstagedFiles` ([]FileChange, required): List of unstaged file changes
- `IsEmpty` (bool, computed): True if no staged or unstaged changes
- `HasChanges` (bool, computed): True if has staged or unstaged changes

**FileChange**:
- `Path` (string): File path relative to repository root
- `Status` (string): Change status (added, modified, deleted, renamed)
- `Diff` (string, optional): Unified diff content for the change

**Computed Properties**:
- `IsEmpty`: Returns true if both StagedFiles and UnstagedFiles are empty
- `HasChanges`: Returns true if either StagedFiles or UnstagedFiles has items
- `TokenEstimate` (int, computed): Estimated token count for AI generation

**Validation Rules**:
- Must be in a valid git repository
- At least one of StagedFiles or UnstagedFiles must have changes (unless creating empty commit)

---

### AIProviderConfig

Represents configuration for an AI provider.

**Fields**:
- `Name` (string, required): Provider name (openai, anthropic, local)
- `APIKey` (string, required): API key or authentication token
- `Model` (string, optional): Model identifier (e.g., "gpt-4", "claude-3-opus")
- `Endpoint` (string, optional): Custom API endpoint (for local models)
- `Timeout` (time.Duration, optional): Request timeout (default: 30s)
- `MaxTokens` (int, optional): Maximum tokens for response (default: 500)

**Validation Rules**:
- Name must be one of supported providers
- APIKey cannot be empty
- Timeout must be > 0

---

### TokenCalculation

Represents token calculation result for repository state.

**Fields**:
- `Provider` (string, required): AI provider name
- `InputTokens` (int, required): Estimated input tokens
- `OutputTokens` (int, required): Estimated output tokens (for commit message)
- `TotalTokens` (int, computed): Sum of input and output tokens
- `Method` (string, required): Calculation method used (tiktoken, custom, fallback)

**Computed Properties**:
- `TotalTokens`: Returns InputTokens + OutputTokens

**Validation Rules**:
- InputTokens must be ≥ 0
- OutputTokens must be ≥ 0
- Method must be one of: "tiktoken", "custom", "fallback"

---

### CommitOptions

Represents CLI options for commit creation.

**Fields**:
- `AutoStage` (bool): Automatically stage all unstaged files (`-a` flag)
- `NoSignoff` (bool): Disable commit signoff (`-s` flag)
- `AIProvider` (string, optional): Override default AI provider
- `SkipAI` (bool): Skip AI generation and go directly to manual input

**Validation Rules**:
- AIProvider must be configured if specified
- AutoStage and NoSignoff are mutually independent flags

---

## Relationships

- `CommitMessage` is created from `RepositoryState` (via AI or manual input)
- `RepositoryState` is analyzed to generate `TokenCalculation`
- `TokenCalculation` uses `AIProviderConfig` for provider-specific calculation
- `CommitOptions` modifies behavior of commit creation workflow

---

## State Machine: Commit Creation Workflow

```
[Start]
  ↓
[Check Repository State]
  ↓
[Has Changes?] → No → [Prompt Empty Commit?] → No → [Exit]
  ↓ Yes                              ↓ Yes
[Calculate Tokens]                   [Manual Input]
  ↓
[Prompt AI Usage?]
  ↓
[Use AI?] → No → [Manual Input]
  ↓ Yes
[Call AI Provider]
  ↓
[AI Success?] → No → [Fallback to Manual Input]
  ↓ Yes
[Validate Format]
  ↓
[Valid?] → No → [Edit or Use with Warning]
  ↓ Yes
[Display for Review]
  ↓
[User Accepts?] → No → [Edit/Restart]
  ↓ Yes
[Create Commit]
  ↓
[End]
```

---

## Data Flow

1. **Repository State Collection**:
   - Git operations → `RepositoryState` (staged/unstaged files)
   - RepositoryState → Token calculation → `TokenCalculation`

2. **Commit Message Creation**:
   - Option A (AI): RepositoryState + AIProviderConfig → AIProvider → `CommitMessage`
   - Option B (Manual): User input → `CommitMessage` components → `CommitMessage`

3. **Validation**:
   - `CommitMessage` → Validator → Validation result (pass/fail with errors)

4. **Commit Execution**:
   - `CommitMessage` + `CommitOptions` → Git operation → Commit created

---

## Persistence

- **Configuration**: Stored in `~/.gitcomm/config.yaml` (AIProviderConfig instances)
- **No Runtime State**: All entities are in-memory during CLI execution
- **No Database**: No persistent storage required (stateless CLI)

---

## Error Types

- `ErrNotGitRepository`: RepositoryState cannot be created (not in git repo)
- `ErrNoChanges`: RepositoryState has no changes and empty commit not confirmed
- `ErrInvalidFormat`: CommitMessage validation failed
- `ErrAIProviderUnavailable`: AI provider call failed
- `ErrEmptySubject`: CommitMessage subject is empty
- `ErrTokenCalculationFailed`: Token calculation error
