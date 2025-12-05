# Data Model: Respect addAll Flag for New Files

**Feature**: 005-addall-new-files
**Date**: 2025-01-27

## Overview

This feature does not introduce new data models. It modifies the behavior of existing models to respect the `addAll` flag when filtering repository state. The existing `RepositoryState` and `FileChange` models remain unchanged.

## Existing Entities

### RepositoryState

**Location**: `internal/model/repository_state.go`

**Description**: Represents the current state of the git repository for commit message generation.

**Fields**:
- `StagedFiles []FileChange` - List of staged file changes (filtered based on `addAll` flag)
- `UnstagedFiles []FileChange` - List of unstaged file changes (unchanged behavior)

**Behavior Changes**:
- `StagedFiles` now excludes new files (status `git.Added`) when `addAll` flag is false
- `StagedFiles` includes all files (including new files) when `addAll` flag is true
- `UnstagedFiles` behavior unchanged (always includes all unstaged files)

**Methods**:
- `IsEmpty() bool` - Returns true if no staged or unstaged changes (behavior unchanged)
- `HasChanges() bool` - Returns true if there are staged or unstaged changes (behavior unchanged)

### FileChange

**Location**: `internal/model/repository_state.go`

**Description**: Represents a single file change in the repository.

**Fields**:
- `Path string` - File path relative to repository root
- `Status string` - Change status (added, modified, deleted, renamed)
- `Diff string` - Optional unified diff content for the change

**Behavior Changes**:
- Files with `Status == "added"` are filtered from `RepositoryState.StagedFiles` when `addAll` is false
- All other statuses (modified, deleted, renamed) are always included regardless of `addAll` flag

**Status Values**:
- `"added"` - New file (filtered when `addAll` is false)
- `"modified"` - Modified tracked file (always included)
- `"deleted"` - Deleted tracked file (always included)
- `"renamed"` - Renamed tracked file (always included)
- `"copied"` - Copied tracked file (always included)
- `"unmerged"` - Unmerged file (always included)

## Context Values

### includeNewFilesKey

**Type**: `contextKey` (string type)

**Location**: `internal/repository/git_repository_impl.go`

**Description**: Context key used to pass `addAll` flag value from service layer to repository layer.

**Usage**:
```go
type contextKey string
const includeNewFilesKey contextKey = "includeNewFiles"

// In service layer:
ctx = context.WithValue(ctx, includeNewFilesKey, s.options.AutoStage)

// In repository layer:
includeNewFiles := true // default
if val := ctx.Value(includeNewFilesKey); val != nil {
    includeNewFiles = val.(bool)
}
```

**Default Value**: `true` (include new files by default for backward compatibility)

## State Transitions

### File Status to Inclusion Decision

```
File Status (git.StatusCode) → Included in RepositoryState.StagedFiles?
├── git.Added → Depends on addAll flag
│   ├── addAll = true → ✅ Included
│   └── addAll = false → ❌ Excluded
├── git.Modified → ✅ Always included
├── git.Deleted → ✅ Always included
├── git.Renamed → ✅ Always included
├── git.Copied → ✅ Always included
└── git.UpdatedButUnmerged → ✅ Always included
```

## Validation Rules

1. **New File Filtering**: Files with `git.Added` status are only included when `addAll` flag is true
2. **Tracked File Inclusion**: Files with any status other than `git.Added` are always included
3. **Unstaged Files**: Unstaged files are always included regardless of `addAll` flag (unchanged behavior)
4. **Empty State**: If all staged files are filtered out, `RepositoryState.IsEmpty()` returns true

## Relationships

- `RepositoryState` contains multiple `FileChange` entities
- `FileChange.Status` determines filtering behavior based on `addAll` flag
- Context value `includeNewFilesKey` connects service layer (`CommitOptions.AutoStage`) to repository layer filtering logic
