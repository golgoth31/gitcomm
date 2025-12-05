# Contract: GitRepository Filter Extension

**Feature**: 005-addall-new-files
**Date**: 2025-01-27
**Interface**: `GitRepository`

## Overview

This contract extends the existing `GitRepository` interface's `GetRepositoryState` method to filter new files based on the `addAll` flag passed via context. The interface signature remains unchanged - only the behavior is modified.

## Method Contract

### GetRepositoryState

**Signature**: (No changes to interface)
```go
GetRepositoryState(ctx context.Context) (*model.RepositoryState, error)
```

**Preconditions**:
- Repository must be initialized and valid
- Context must be valid (not cancelled)
- Context may contain `includeNewFilesKey` value (optional, defaults to `true`)

**Postconditions**:
- Returns `RepositoryState` with `StagedFiles` filtered based on `addAll` flag
- New files (`git.Added` status) are excluded when `addAll` is false
- Modified, deleted, renamed files are always included
- `UnstagedFiles` behavior unchanged (always includes all unstaged files)
- No errors returned for individual file filtering (filtering is silent)

**Behavior**:

1. **Context Value Extraction**:
   - Reads `includeNewFilesKey` from context (defaults to `true` if not present)
   - Value type: `bool`
   - Default behavior (when not present): Include all files (backward compatible)

2. **File Filtering Rules**:
   - **New files (`git.Added`)**:
     - Included when `includeNewFiles == true`
     - Excluded when `includeNewFiles == false`
   - **Modified files (`git.Modified`)**: Always included
   - **Deleted files (`git.Deleted`)**: Always included
   - **Renamed files (`git.Renamed`)**: Always included
   - **Copied files (`git.Copied`)**: Always included
   - **Unmerged files (`git.UpdatedButUnmerged`)**: Always included

3. **Filtering Implementation**:
   - Filtering occurs during iteration over `worktree.Status()` results
   - Files with `fileStatus.Staging == git.Added && !includeNewFiles` are skipped
   - All other files are processed normally (diff computation, etc.)

4. **Error Handling**:
   - Filtering failures are silent (no errors returned)
   - If context value is invalid type, default to `true` (include all files)
   - Existing error handling for git operations unchanged

## Context Key Definition

**Key**: `includeNewFilesKey` (type: `contextKey`)

**Location**: `internal/repository/git_repository_impl.go`

**Type**:
```go
type contextKey string
const includeNewFilesKey contextKey = "includeNewFiles"
```

**Value Type**: `bool`

**Default**: `true` (include new files) - ensures backward compatibility

## Service Layer Contract

### CommitService Integration

**Method**: `CreateCommit(ctx context.Context) error`

**Behavior**:
- Extracts `AutoStage` value from `CommitOptions`
- Sets context value before calling `GetRepositoryState`:
  ```go
  ctx = context.WithValue(ctx, includeNewFilesKey, s.options.AutoStage)
  state, err := s.gitRepo.GetRepositoryState(ctx)
  ```

**Preconditions**:
- `CommitService` must have valid `options` with `AutoStage` field set
- Context must be valid

**Postconditions**:
- Context contains `includeNewFilesKey` with value matching `AutoStage`
- `GetRepositoryState` receives context with filtering flag
- Repository state respects `addAll` flag

## Backward Compatibility

**Compatibility Guarantee**:
- When context value is not present, behavior defaults to including all files
- Existing callers of `GetRepositoryState` without context value continue to work
- No interface signature changes required
- No breaking changes to existing code

**Migration Path**:
- Existing code: No changes required (defaults to include all files)
- New code: Set context value to control filtering behavior

## Testing Requirements

**Unit Tests**:
- Test filtering with `includeNewFiles = true` (include all files)
- Test filtering with `includeNewFiles = false` (exclude new files)
- Test default behavior when context value not present (include all files)
- Test with various file status combinations (Added, Modified, Deleted, Renamed)
- Test with binary files (should follow same filtering rules)
- Test with manually staged new files

**Integration Tests**:
- End-to-end test: `gitcomm` without `-a` flag excludes new files
- End-to-end test: `gitcomm -a` includes new files
- Test with real git repositories containing new and modified files
