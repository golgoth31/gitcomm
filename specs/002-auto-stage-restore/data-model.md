# Data Model: Auto-Stage Modified Files and State Restoration

**Feature**: 002-auto-stage-restore
**Date**: 2025-01-27

## Domain Entities

### StagingState

Represents a snapshot of the git repository staging state at a specific point in time.

**Fields**:
- `StagedFiles []string` - List of file paths that are staged
- `CapturedAt time.Time` - Timestamp when state was captured
- `RepositoryPath string` - Path to the git repository root

**Methods**:
- `IsEmpty() bool` - Returns true if no files are staged
- `Contains(file string) bool` - Returns true if file is staged
- `Diff(other *StagingState) []string` - Returns files in this state but not in other (for restoration)

**Validation Rules**:
- `StagedFiles` must contain valid file paths (relative to repository root)
- `RepositoryPath` must be an absolute path to a valid git repository
- `CapturedAt` must be a valid timestamp

**Lifecycle**:
1. Created at CLI launch (capture pre-CLI state)
2. Used during auto-staging to track which files were staged by CLI
3. Used during restoration to determine which files to unstage
4. Discarded after successful commit or restoration

---

### AutoStagingResult

Represents the result of an automatic staging operation.

**Fields**:
- `StagedFiles []string` - List of file paths successfully staged
- `FailedFiles []StagingFailure` - List of files that failed to stage
- `Success bool` - Overall success status (true if all files staged)
- `Duration time.Duration` - Time taken for staging operation

**Methods**:
- `HasFailures() bool` - Returns true if any files failed to stage
- `GetFailedFilePaths() []string` - Returns list of failed file paths

**Validation Rules**:
- If `Success` is true, `FailedFiles` must be empty
- If `Success` is false, `FailedFiles` must not be empty
- `StagedFiles` and `FailedFiles` must not overlap

**Lifecycle**:
1. Created when auto-staging operation completes
2. Used to determine if restoration is needed
3. Discarded after restoration or successful commit

---

### StagingFailure

Represents a single file staging failure.

**Fields**:
- `FilePath string` - Path to the file that failed to stage
- `Error error` - The error that occurred during staging
- `ErrorType string` - Type of error (permission, locked, conflict, etc.)

**Validation Rules**:
- `FilePath` must be a valid file path
- `Error` must not be nil
- `ErrorType` must be one of: "permission", "locked", "conflict", "not_found", "other"

---

### RestorationPlan

Represents the plan for restoring staging state to pre-CLI state.

**Fields**:
- `FilesToUnstage []string` - List of file paths to unstage (files staged by CLI)
- `PreCLIState *StagingState` - The captured pre-CLI staging state
- `CurrentState *StagingState` - The current staging state (for validation)

**Methods**:
- `IsEmpty() bool` - Returns true if no restoration needed
- `Validate() error` - Validates that restoration plan is valid
- `Execute(ctx context.Context, repo GitRepository) error` - Executes the restoration

**Validation Rules**:
- `FilesToUnstage` must only contain files that were staged by CLI (not pre-existing)
- `PreCLIState` must not be nil
- `CurrentState` must not be nil
- `FilesToUnstage` must be a subset of `CurrentState.StagedFiles`

**Lifecycle**:
1. Created when restoration is needed (cancellation, error, interruption)
2. Validated before execution
3. Executed to restore state
4. Discarded after execution

---

## Relationships

- `StagingState` is captured from `GitRepository` at CLI launch
- `AutoStagingResult` is produced by `GitRepository.StageModifiedFiles()` or `GitRepository.StageAllFiles()`
- `StagingFailure` is part of `AutoStagingResult` when staging fails
- `RestorationPlan` is created from `StagingState` (pre-CLI) and current `StagingState`
- `RestorationPlan` uses `GitRepository` to execute restoration

---

## State Transitions

### Auto-Staging Workflow

```
[CLI Launch]
  ↓
[Capture Pre-CLI StagingState]
  ↓
[Auto-Stage Modified Files]
  ↓
[Create AutoStagingResult]
  ↓
[Success?] → No → [Create RestorationPlan] → [Restore State] → [Exit with Error]
  ↓ Yes
[Continue Workflow]
```

### Restoration Workflow

```
[Restoration Triggered]
  ↓
[Create RestorationPlan from Pre-CLI State]
  ↓
[Validate RestorationPlan]
  ↓
[Execute Restoration]
  ↓
[Success?] → No → [Log Error] → [Display Warning] → [Exit with Error]
  ↓ Yes
[Exit Successfully]
```

### Signal Interruption Workflow

```
[Signal Received]
  ↓
[Check if Staging in Progress]
  ↓
[If Staging Complete] → [Create RestorationPlan] → [Execute Restoration] → [Exit]
[If Staging In Progress] → [Abort Staging] → [Restore Partial State] → [Exit]
```

---

## Data Flow

1. **Staging State Capture**:
   - `GitRepository.GetRepositoryState()` → `StagingState` (pre-CLI)

2. **Auto-Staging**:
   - `GitRepository.StageModifiedFiles()` or `GitRepository.StageAllFiles()` → `AutoStagingResult`

3. **Restoration Planning**:
   - `StagingState` (pre-CLI) + `StagingState` (current) → `RestorationPlan`

4. **Restoration Execution**:
   - `RestorationPlan.Execute()` → `GitRepository.UnstageFiles()` → Repository state restored

---

## Persistence

- **No Persistent Storage**: All entities are in-memory during CLI execution
- **Staging State**: Captured at CLI launch, stored in memory until restoration or commit
- **No Database**: No persistent storage required (stateless CLI)

---

## Error Types

- `ErrStagingFailed`: Auto-staging operation failed (partial or complete failure)
- `ErrRestorationFailed`: State restoration operation failed
- `ErrStagingStateInvalid`: Captured staging state is invalid or corrupted
- `ErrRestorationPlanInvalid`: Restoration plan is invalid (e.g., files don't exist)
- `ErrInterruptedDuringStaging`: CLI was interrupted while staging was in progress

---

## Integration with Existing Models

### RepositoryState (existing)
- Extended to support staging state capture
- Used to determine which files need staging
- Used to validate restoration state

### FileChange (existing)
- Used to represent files in staging state
- Used to track which files were staged by CLI

### CommitOptions (existing)
- `AutoStage` flag already exists
- No changes needed for this feature
