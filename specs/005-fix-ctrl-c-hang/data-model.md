# Data Model: Fix CLI Hang on Ctrl+C During State Restoration

**Feature**: 005-fix-ctrl-c-hang
**Date**: 2025-01-27

## Domain Entities

### RestorationContext

Represents a timeout-aware context used for restoration operations when triggered by Ctrl+C.

**Fields**:
- `Context context.Context` - The timeout context (3 seconds)
- `Cancel context.CancelFunc` - Function to cancel the context
- `TimeoutDuration time.Duration` - Duration of timeout (3 seconds)
- `IsInterrupted bool` - Whether this restoration was triggered by Ctrl+C

**Methods**:
- `WithTimeout(parent context.Context, timeout time.Duration) (RestorationContext, error)` - Creates a new timeout context
- `IsTimeout(err error) bool` - Checks if error is due to timeout
- `Cancel()` - Cancels the context

**Validation Rules**:
- `TimeoutDuration` must be between 1 and 5 seconds (3 seconds per FR-002)
- `Context` must not be nil
- `Cancel` must be called (deferred) to prevent context leaks

**Lifecycle**:
1. Created when Ctrl+C signal is received
2. Used for all restoration operations
3. Automatically cancelled after timeout or when operations complete
4. Discarded after restoration completes or times out

---

## State Transitions

### Restoration with Timeout Workflow

```
[Ctrl+C Signal Received]
  ↓
[Create RestorationContext with 3-second timeout]
  ↓
[Display "Interrupted. Restoring staging state..."]
  ↓
[Start restoration operations with timeout context]
  ↓
[Operations complete within timeout?]
  ├─ Yes → [Exit immediately with code 130]
  └─ No → [Timeout occurs]
           ↓
           [Display warning message]
           ↓
           [Exit immediately with code 130]
```

### Timeout Detection

```
[Restoration operation in progress]
  ↓
[Check context deadline]
  ├─ Deadline exceeded → [Return timeout error]
  └─ Not exceeded → [Continue operation]
                      ↓
                      [Check context cancellation]
                      ├─ Cancelled → [Return cancellation error]
                      └─ Active → [Continue operation]
```

---

## Relationships

- `RestorationContext` is used by `CommitService.restoreStagingState` method
- `RestorationContext` wraps standard `context.Context` for timeout management
- `RestorationContext` is created in signal handler and passed to restoration operations

---

## Error Types

- `ErrRestorationTimeout` - Restoration operation exceeded timeout (3 seconds)
- `ErrRestorationCancelled` - Restoration operation was cancelled before timeout
- `ErrRestorationFailed` - Restoration operation failed for other reasons (existing)

---

## Integration with Existing Models

### StagingState (existing)
- No changes needed
- Used as input to restoration operations

### RestorationPlan (existing)
- No changes needed
- Used for calculating files to unstage

### CommitService (existing)
- Modified to use `RestorationContext` instead of `context.Background()`
- Modified to handle timeout errors gracefully

---

## Persistence

- **No Persistent Storage**: All context state is in-memory only during restoration
- **No Configuration File**: Timeout duration is hardcoded (3 seconds per FR-002)
- **No Database**: No persistent storage required
