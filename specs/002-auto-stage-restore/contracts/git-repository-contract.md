# Git Repository Interface Contract: Staging State Management

**Feature**: 002-auto-stage-restore
**Date**: 2025-01-27
**Interface**: `GitRepository`

## Extended Interface

### Existing Methods (from 001-git-commit-cli)

- `GetRepositoryState(ctx context.Context) (*model.RepositoryState, error)`
- `CreateCommit(ctx context.Context, message *model.CommitMessage) error`
- `StageAllFiles(ctx context.Context) error`

### New Methods

#### CaptureStagingState

Captures the current staging state of the repository for restoration purposes.

```go
CaptureStagingState(ctx context.Context) (*model.StagingState, error)
```

**Input**:
- `ctx context.Context` - Context for cancellation and timeout

**Output**:
- `*model.StagingState` - Snapshot of current staging state
- `error` - Error if capture fails (e.g., not a git repository, git operation failed)

**Behavior**:
- Returns snapshot of which files are currently staged
- Must be called before any auto-staging operations
- Returns error if repository state cannot be determined

**Error Cases**:
- `ErrNotGitRepository`: Not in a git repository
- `ErrGitOperationFailed`: Git operation failed (e.g., repository locked)

---

#### StageModifiedFiles

Stages all modified (but not untracked) files in the repository.

```go
StageModifiedFiles(ctx context.Context) (*model.AutoStagingResult, error)
```

**Input**:
- `ctx context.Context` - Context for cancellation and timeout

**Output**:
- `*model.AutoStagingResult` - Result of staging operation
- `error` - Error if staging fails completely (all files fail)

**Behavior**:
- Stages all modified files (equivalent to `git add` for modified files)
- Does not stage untracked files
- Returns result with list of successfully staged files and any failures
- If any file fails, all staged files should be restored (abort behavior)

**Error Cases**:
- `ErrStagingFailed`: Staging operation failed (partial or complete)
- `ErrNotGitRepository`: Not in a git repository
- `ErrGitOperationFailed`: Git operation failed

---

#### StageAllFilesIncludingUntracked

Stages all modified and untracked files in the repository (equivalent to `git add -A`).

```go
StageAllFilesIncludingUntracked(ctx context.Context) (*model.AutoStagingResult, error)
```

**Input**:
- `ctx context.Context` - Context for cancellation and timeout

**Output**:
- `*model.AutoStagingResult` - Result of staging operation
- `error` - Error if staging fails completely

**Behavior**:
- Stages all modified and untracked files
- Used when `-a` flag is provided
- Returns result with list of successfully staged files and any failures
- If any file fails, all staged files should be restored (abort behavior)

**Error Cases**:
- `ErrStagingFailed`: Staging operation failed
- `ErrNotGitRepository`: Not in a git repository
- `ErrGitOperationFailed`: Git operation failed

---

#### UnstageFiles

Unstages the specified files, restoring them to their pre-staged state.

```go
UnstageFiles(ctx context.Context, files []string) error
```

**Input**:
- `ctx context.Context` - Context for cancellation and timeout
- `files []string` - List of file paths to unstage

**Output**:
- `error` - Error if unstaging fails

**Behavior**:
- Unstages only the specified files
- Files must be currently staged
- If a file is not staged, it is ignored (no error)
- If any file fails to unstage, operation continues for remaining files (best-effort)

**Error Cases**:
- `ErrRestorationFailed`: Restoration operation failed
- `ErrNotGitRepository`: Not in a git repository
- `ErrGitOperationFailed`: Git operation failed

---

## Implementation Contract

### Thread Safety

- All methods must be thread-safe
- Use `context.Context` for cancellation
- No shared mutable state between concurrent calls

### Error Handling

- All errors must be wrapped with context: `fmt.Errorf("operation: %w", err)`
- Use custom error types from `internal/utils/errors.go`
- Never panic or exit in implementation

### Performance

- `CaptureStagingState`: Must complete within 500ms for typical repositories
- `StageModifiedFiles`: Must complete within 2 seconds (SC-001)
- `StageAllFilesIncludingUntracked`: Must complete within 2 seconds
- `UnstageFiles`: Must complete within 1 second (SC-003)

### Resource Management

- All git operations must respect `context.Context` cancellation
- No resource leaks (file handles, network connections)
- Cleanup on error or cancellation

---

## Testing Contract

### Unit Tests

- Test each method with valid inputs
- Test error cases (not git repo, locked repository, etc.)
- Test cancellation via context
- Test empty repository state

### Integration Tests

- Test full workflow: capture → stage → restore
- Test with real git repository
- Test signal interruption during operations
- Test external state changes

### Mock Requirements

- Mock implementation must support all new methods
- Mock must allow injection of errors for testing
- Mock must track state changes for verification
