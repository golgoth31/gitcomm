# Research: Auto-Stage Modified Files and State Restoration

**Feature**: 002-auto-stage-restore
**Date**: 2025-01-27

## Technology Decisions

### 1. Git Staging State Capture

**Decision**: Use `go-git/v5` worktree status to capture staging state snapshot at CLI launch.

**Rationale**:
- `go-git/v5` is already a dependency and provides `worktree.Status()` method
- Status returns a map of file paths to `FileStatus` with staging/worktree information
- Can capture which files are staged before auto-staging begins
- No need for additional git operations or external commands

**Alternatives Considered**:
- `git diff --cached --name-only` via exec: More complex, requires external git command
- Storing file list in memory: Simpler but less robust if external changes occur
- Git index inspection: Lower-level, more complex, unnecessary for this use case

**Implementation Pattern**:
```go
// Capture staging state
status, err := worktree.Status()
preCLIStaged := make(map[string]bool)
for file, fileStatus := range status {
    if fileStatus.Staging != git.Unmodified {
        preCLIStaged[file] = true
    }
}
```

### 2. Signal Handling for Interruption

**Decision**: Use Go's `os/signal` package with buffered channels for graceful interruption handling.

**Rationale**:
- Standard Go approach for signal handling
- Thread-safe via channels
- Allows cleanup before exit
- Works cross-platform (SIGINT, SIGTERM)

**Alternatives Considered**:
- No signal handling: Leaves repository in inconsistent state on interruption
- Panic recovery: Not appropriate for normal interruption flow
- External signal handlers: Unnecessary complexity

**Implementation Pattern**:
```go
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
go func() {
    <-sigChan
    // Restore staging state
    // Exit gracefully
}()
```

### 3. Staging State Restoration Strategy

**Decision**: Track files staged by CLI (not pre-existing staged files) and unstage only those files during restoration.

**Rationale**:
- Preserves user's original staging state (FR-011)
- Only affects files modified by CLI
- Simpler restoration logic
- Aligns with user expectations

**Alternatives Considered**:
- Unstage all files: Would lose user's original staging, violates FR-011
- Full state snapshot/restore: More complex, unnecessary for use case
- Git reset: Too destructive, affects worktree changes

**Implementation Pattern**:
```go
// Track files staged by CLI
cliStagedFiles := []string{}
// After auto-staging, compare with pre-CLI state
// Only unstage files that were staged by CLI
```

### 4. Error Handling for Staging Failures

**Decision**: Abort staging operation on any failure, restore all staged files, exit with error.

**Rationale**:
- Prevents partial staging (clarification Q1)
- Ensures repository consistency
- Clear failure mode for users
- Aligns with FR-009

**Alternatives Considered**:
- Continue with partial staging: Leaves repository in inconsistent state
- Prompt user on failure: Adds complexity, delays failure detection
- Retry with backoff: Unnecessary for file system operations

### 5. Empty Repository State Handling

**Decision**: Treat auto-staging as no-op when no files to stage, proceed with existing workflow.

**Rationale**:
- Aligns with clarification Q2
- Maintains existing empty commit handling
- No special case needed
- Simple implementation

**Alternatives Considered**:
- Exit immediately: Breaks existing workflow
- Error on empty: Unnecessary, existing workflow handles this

### 6. External Repository State Changes

**Decision**: Restore to captured pre-CLI state, ignoring external changes, log warning if mismatch detected.

**Rationale**:
- Predictable behavior (clarification Q5)
- Prevents external changes from affecting restoration
- Warning provides visibility without blocking
- Aligns with FR-003

**Alternatives Considered**:
- Detect and abort: Too strict, external changes may be intentional
- Restore relative to current: Complex, unpredictable behavior

## Best Practices

### Git Operations
- Use `context.Context` for cancellation
- Handle errors explicitly with wrapped errors
- Validate repository state before operations
- Use transactions where possible (git operations are atomic)

### Signal Handling
- Use buffered channels (size 1) for signal channel
- Handle cleanup in signal handler
- Don't block signal handler (use goroutines for cleanup)
- Exit with appropriate exit code

### State Management
- Capture state early (at CLI launch)
- Store minimal necessary information
- Validate state before restoration
- Log state changes for debugging

## Integration Points

### Existing Code
- `GitRepository` interface: Extend with staging state methods
- `CommitService`: Integrate auto-staging before workflow
- `main.go`: Add signal handling for interruption
- Error types: Add staging/restoration error types

### Dependencies
- No new external dependencies required
- Uses existing `go-git/v5` capabilities
- Uses standard library `os/signal`

## Performance Considerations

- Staging state capture: O(n) where n is number of files (acceptable for typical repos)
- Restoration: O(m) where m is number of files staged by CLI (typically small)
- Signal handling: Minimal overhead (channel operations)
- All operations complete within success criteria time limits

## Security Considerations

- No secrets involved in staging operations
- File paths validated (no path traversal)
- Git operations use existing security model
- No external network calls

## Testing Strategy

- Unit tests: Staging state capture/restoration logic
- Integration tests: Full workflow with git repository
- Signal handling tests: Mock signals, verify restoration
- Edge case tests: Empty state, partial failures, external changes
