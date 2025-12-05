# Research: Fix CLI Hang on Ctrl+C During State Restoration

**Feature**: 005-fix-ctrl-c-hang
**Date**: 2025-01-27

## Technology Decisions

### 1. Timeout Context for Restoration Operations

**Decision**: Use `context.WithTimeout` to create a 3-second timeout context for restoration operations when triggered by Ctrl+C.

**Rationale**:
- `context.WithTimeout` is the idiomatic Go way to handle timeouts
- Automatically cancels operations after the timeout duration
- Respects cancellation signals (if context is cancelled before timeout)
- Integrates seamlessly with existing `context.Context` usage in git operations
- Standard library solution, no additional dependencies

**Alternatives Considered**:
- Custom timeout wrapper: More complex, reinvents standard library functionality
- Channel-based timeout: More verbose, requires manual goroutine management
- `time.AfterFunc`: Less integrated with context cancellation patterns

**Implementation Pattern**:
```go
// Create timeout context for restoration (3 seconds)
restoreCtx, restoreCancel := context.WithTimeout(context.Background(), 3*time.Second)
defer restoreCancel()

// Use restoreCtx for all restoration operations
if err := s.restoreStagingState(restoreCtx, preCLIState); err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        // Timeout occurred
        fmt.Printf("Warning: Restoration timed out. Repository may be in unexpected state.\n")
    }
    // Handle other errors
}
```

### 2. Synchronization Between Signal Handler and Main Process

**Decision**: Use a channel to signal restoration completion and ensure main process waits before exiting.

**Rationale**:
- Channels are the idiomatic Go way to synchronize goroutines
- Prevents race conditions between signal handler and main process
- Allows main process to wait for restoration completion or timeout
- Simple and clear synchronization pattern
- No need for complex mutex or waitgroup management

**Alternatives Considered**:
- `sync.WaitGroup`: More complex for this use case, requires careful management
- `sync.Mutex`: Overkill for simple completion signaling
- Direct function call: Doesn't work with async signal handling

**Implementation Pattern**:
```go
// In signal handler goroutine
restoreDone := make(chan struct{})
go func() {
    defer close(restoreDone)
    // Perform restoration with timeout
    s.restoreStagingState(restoreCtx, preCLIState)
}()

// Wait for restoration or timeout (with overall 5-second limit)
select {
case <-restoreDone:
    // Restoration completed
case <-time.After(5 * time.Second):
    // Overall timeout exceeded
    fmt.Printf("Warning: Restoration did not complete in time.\n")
}
```

### 3. Context Propagation in Git Operations

**Decision**: Ensure all git operations in `restoreStagingState` respect the timeout context.

**Rationale**:
- Git operations (`CaptureStagingState`, `UnstageFiles`) already accept `context.Context`
- Need to ensure they properly check context cancellation/timeout
- Standard pattern for cancellation-aware operations
- Prevents blocking on slow or stuck git operations

**Alternatives Considered**:
- Ignore context in git operations: Would defeat the purpose of timeout
- Separate timeout per operation: More complex, harder to manage overall timeout

**Implementation Pattern**:
```go
// All git operations use the timeout context
currentState, err := s.gitRepo.CaptureStagingState(restoreCtx)
if err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        // Timeout - exit immediately
        return err
    }
    // Other errors
}

// Unstage operations also respect context
if err := s.gitRepo.UnstageFiles(restoreCtx, files); err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        // Timeout - exit immediately
        return err
    }
    // Other errors
}
```

### 4. Warning Message Content

**Decision**: Display clear, actionable warning messages when restoration times out or fails.

**Rationale**:
- Users need to know the repository state may be inconsistent
- Provides actionable guidance (check git status)
- Consistent with existing error message patterns in the codebase
- Helps users understand what happened and what to do next

**Message Format**:
- Timeout: "Warning: Restoration timed out. Repository may be in unexpected state. Please check git status and manually restore if needed."
- Failure: "Warning: Failed to restore staging state. Repository may be in unexpected state. Please check git status and manually restore if needed."

## Best Practices

### Context Timeout Management
- Always defer cancel() for timeout contexts to prevent leaks
- Check for `context.DeadlineExceeded` to distinguish timeout from other errors
- Use `errors.Is(err, context.DeadlineExceeded)` for error checking

### Signal Handling
- Keep signal handlers minimal (just set flags, don't do heavy work)
- Use channels for communication between signal handler and main process
- Ensure main process waits for cleanup operations before exiting

### Error Handling
- Distinguish between timeout errors and other errors
- Provide clear, actionable error messages
- Log detailed errors for debugging (debug mode)

## Integration Points

### Existing Code
- `cmd/gitcomm/main.go`: Signal handler goroutine (line 108-114)
- `internal/service/commit_service.go`: `restoreStagingState` method (line 218-269)
- `internal/service/commit_service.go`: Defer function for restoration (line 53-67)
- `internal/repository/git_repository_impl.go`: Git operations that accept context

### Dependencies
- Standard library: `context`, `time`, `errors`
- Existing: `github.com/go-git/go-git/v5` (already supports context)

## Performance Considerations

- Timeout context creation: Negligible overhead (<1ms)
- Context cancellation check: Minimal overhead in git operations
- Overall timeout enforcement: Ensures CLI exits within 5 seconds (SC-001)

## Testing Strategy

- Unit tests: Timeout context creation and error handling
- Integration tests: Signal handling with timeout scenarios
- Timeout simulation: Use short timeouts in tests to verify behavior
- Race condition tests: Verify synchronization between signal handler and main process
