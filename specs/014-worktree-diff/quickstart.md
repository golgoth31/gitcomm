# Quickstart: Compute Worktree Diff in GetRepositoryState

**Feature**: 014-worktree-diff
**Date**: 2025-01-27

## Overview

This feature enhances the `GetRepositoryState` method to compute and include unified diff content for staged files. The diff represents changes between the staged state and HEAD, optimized for token usage (0 context lines, 5000 character limit per file).

## Quick Implementation Guide

### 1. Understand the Enhancement

The `GetRepositoryState` method in `internal/repository/git_repository_impl.go` is enhanced to:
- Compute diffs for all staged files using go-git plumbing API
- Populate the `Diff` field in `FileChange` for staged files
- Apply token optimization (0 context, 5000 char limit)
- Handle edge cases (binary files, unmerged files, empty repo)

### 2. Key Implementation Steps

#### Step 1: Get HEAD Tree
```go
head, err := r.repo.Head()
if err == plumbing.ErrReferenceNotFound {
    // Empty repository - treat as empty tree
    headTree = &object.Tree{} // or use nil tree
} else {
    headCommit, err := r.repo.CommitObject(head.Hash())
    headTree, err := headCommit.Tree()
}
```

#### Step 2: Get Staged Index Tree
```go
index, err := r.repo.Storer.Index()
// Convert index entries to tree for diff computation
stagedTree, err := buildTreeFromIndex(index)
```

#### Step 3: Compute Diff for Each Staged File
```go
for file, fileStatus := range status {
    if fileStatus.Staging == git.Unmodified {
        continue // Skip unstaged files
    }

    // Compute diff
    diff, err := computeFileDiff(headTree, stagedTree, file)
    if err != nil {
        // Log error, set Diff="", continue (FR-010)
        continue
    }

    // Check size limit
    if len(diff) > 5000 {
        diff = generateMetadata(file, fileStatus)
    }

    // Populate Diff field
    change.Diff = diff
}
```

#### Step 4: Handle Special Cases
- **Binary files**: Set `Diff = ""` (FR-013)
- **Unmerged files**: Attempt diff, fallback to `""` (FR-008)
- **Renamed/Copied**: Format with similarity percentage (FR-014, FR-015)
- **Empty repo**: Treat as empty tree (FR-009)

### 3. Testing Strategy

#### Unit Tests
```go
func TestGetRepositoryState_WithDiff(t *testing.T) {
    // Test diff computation for each file status
    // Test size limit enforcement
    // Test binary file handling
    // Test error handling
}
```

#### Integration Tests
```go
func TestGetRepositoryState_DiffComputation(t *testing.T) {
    // Test end-to-end with real repository
    // Test performance (100 files <2s)
    // Test error rate (<1%)
}
```

### 4. Key Constants

```go
const (
    maxDiffSize = 5000 // characters (FR-016)
    diffContext = 0    // lines (FR-012)
)
```

### 5. Error Handling Pattern

```go
// Per-file error handling (FR-010)
for file, fileStatus := range status {
    diff, err := computeDiff(file)
    if err != nil {
        utils.Logger.Debug().
            Err(err).
            Str("file", file).
            Msg("Failed to compute diff, continuing")
        change.Diff = "" // Continue with empty diff
        continue
    }
    change.Diff = diff
}
```

## Verification Checklist

- [ ] Diffs computed for all staged files
- [ ] Unstaged files have empty Diff fields
- [ ] Binary files have empty Diff fields
- [ ] Large files/diffs (>5000 chars) show metadata only
- [ ] 0 lines of context in unified diff format
- [ ] Renamed files show rename diff with similarity
- [ ] Copied files show copy diff with similarity
- [ ] Unmerged files handled gracefully
- [ ] Empty repository handled correctly
- [ ] Individual file failures don't block others
- [ ] Performance: 100 files in <2 seconds
- [ ] Error rate: <1% of files

## Common Pitfalls

1. **Forgetting to check file size**: Always check if diff exceeds 5000 chars before assigning
2. **Not handling empty repository**: Check for `plumbing.ErrReferenceNotFound` when getting HEAD
3. **Binary file detection**: Use go-git's `IsBinary()` method, don't guess
4. **Error handling**: Don't return error for individual file failures, log and continue
5. **Unmerged files**: Attempt diff computation first, don't skip immediately

## Next Steps

1. Review research.md for go-git API details
2. Review data-model.md for data structure constraints
3. Review contracts/ for interface requirements
4. Write tests first (TDD approach)
5. Implement diff computation logic
6. Add error handling and edge cases
7. Verify performance requirements
