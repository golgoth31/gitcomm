# Contract: GitRepository Diff Computation Extension

**Feature**: 014-worktree-diff
**Date**: 2025-01-27
**Interface**: `GitRepository`

## Overview

This contract extends the existing `GitRepository` interface's `GetRepositoryState` method to compute and populate diff content for staged files. The interface itself does not change - only the behavior and return value structure.

## Method Contract

### GetRepositoryState

**Signature**: (No changes to interface)
```go
GetRepositoryState(ctx context.Context) (*model.RepositoryState, error)
```

**Preconditions**:
- Repository must be initialized and valid
- Context must be valid (not cancelled)

**Postconditions**:
- Returns `RepositoryState` with `StagedFiles` containing `FileChange` entries
- Each staged `FileChange.Diff` field is populated according to rules below
- Each unstaged `FileChange.Diff` field is empty (FR-011)
- No errors returned for individual file failures (FR-010)

**Behavior**:

1. **Diff Computation**:
   - Computes diff for all staged files (FR-001)
   - Uses go-git plumbing diff API (FR-002)
   - Compares staged index to HEAD (clean worktree state) (FR-002)
   - Uses unified diff format with 0 lines of context (FR-012)

2. **Diff Population Rules**:
   - **Staged files**: Diff field populated (unless conditions below apply)
   - **Unstaged files**: Diff field always empty (FR-011)
   - **Binary files**: Diff field empty (FR-013)
   - **Large files/diffs (>5000 chars)**: Diff field contains metadata only (FR-016)
   - **Unmerged files**: Attempt diff computation, fallback to empty if fails (FR-008)
   - **Failed computation**: Diff field empty, error logged (FR-010)

3. **File Status Handling**:
   - **Modified**: Full diff content (or metadata if >5000 chars)
   - **Added**: Full diff content (or metadata if >5000 chars)
   - **Deleted**: Full diff content showing deletion
   - **Renamed**: Rename diff with similarity percentage (FR-014)
   - **Copied**: Copy diff with similarity percentage (FR-015)
   - **Unmerged**: Attempt diff, fallback to empty (FR-008)

4. **Error Handling**:
   - Individual file failures don't block other files (FR-010)
   - Errors logged via existing logging infrastructure
   - Failed files have empty Diff field
   - Method returns error only if repository access fails (not per-file failures)

**Performance Requirements**:
- Completes for up to 100 staged files in under 2 seconds (SC-003)
- Error rate < 1% of files (SC-004)

**Token Optimization**:
- 0 lines of context (FR-012)
- 5000 character limit per file/diff (FR-016)
- Metadata-only for large content (FR-006)

## Error Conditions

| Condition | Behavior | Error Returned |
|-----------|----------|----------------|
| Repository not initialized | Return error | Yes |
| Context cancelled | Return error | Yes |
| HEAD not found (empty repo) | Treat as empty tree, continue | No (FR-009) |
| File read error | Log error, Diff="", continue | No (FR-005) |
| Diff computation failure (single file) | Log error, Diff="", continue | No (FR-010) |
| Binary file | Diff="", continue | No (FR-013) |
| Unmerged file | Attempt diff, fallback to "", continue | No (FR-008) |

## Examples

### Example 1: Standard Modified File

**Input**: File `src/main.go` staged with modifications

**Output**:
```go
FileChange{
    Path: "src/main.go",
    Status: "modified",
    Diff: "diff --git a/src/main.go b/src/main.go\nindex abc123..def456 100644\n--- a/src/main.go\n+++ b/src/main.go\n@@ -10,0 +11,3 @@\n+func newFunction() {\n+    // implementation\n+}",
}
```

### Example 2: Large File (Exceeds 5000 Characters)

**Input**: File `large.go` staged with >5000 character diff

**Output**:
```go
FileChange{
    Path: "large.go",
    Status: "modified",
    Diff: "file: large.go\nsize: 10240 bytes\nlines: 250\nchanges: +150 -100 (50 modified)",
}
```

### Example 3: Binary File

**Input**: File `image.png` staged

**Output**:
```go
FileChange{
    Path: "image.png",
    Status: "added",
    Diff: "", // Empty per FR-013
}
```

### Example 4: Renamed File

**Input**: File renamed from `old.go` to `new.go`

**Output**:
```go
FileChange{
    Path: "new.go",
    Status: "renamed",
    Diff: "rename from old.go\nrename to new.go\nsimilarity 95%",
}
```

## Testing Requirements

**Unit Tests**:
- Diff computation for each file status type
- Size limit enforcement (5000 char threshold)
- Binary file detection
- Error handling (individual file failures)
- Empty repository handling
- Unmerged file handling

**Integration Tests**:
- End-to-end GetRepositoryState with diff computation
- Multiple staged files with various statuses
- Performance test (100 files in <2 seconds)
- Error rate test (<1% failure rate)

## Backward Compatibility

**Breaking Changes**: None
- Interface signature unchanged
- Existing callers continue to work
- Diff field was already optional (now populated)

**Migration**: None required
- Feature is additive enhancement
- No API changes
- No data model changes (Diff field already exists)
