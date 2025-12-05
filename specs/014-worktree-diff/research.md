# Research: Compute Worktree Diff in GetRepositoryState

**Feature**: 014-worktree-diff
**Date**: 2025-01-27

## Technology Decisions

### 1. Go-Git Diff Computation API

**Decision**: Use `go-git/v5` plumbing diff API to compute diffs between staged index and HEAD.

**Rationale**:
- `go-git/v5` is already a dependency and provides diff computation capabilities
- Plumbing API (`plumbing/object.DiffTree`, `plumbing/object.Patch`) supports computing diffs between trees
- Can compute diff between index (staged state) and HEAD (clean state)
- Supports unified diff format with configurable context lines
- Handles renamed/copied files with similarity detection

**Alternatives Considered**:
- `git diff --cached` via exec: More complex, requires external git command, less control
- Manual diff computation: Too complex, error-prone, reinventing the wheel
- Different git library: Would require new dependency, go-git is already established

**Implementation Pattern**:
```go
// Get HEAD tree
head, err := r.repo.Head()
headCommit, err := r.repo.CommitObject(head.Hash())
headTree, err := headCommit.Tree()

// Get index (staged) tree
index, err := r.repo.Storer.Index()
// Convert index to tree for diff computation

// Compute diff
diffs, err := object.DiffTree(headTree, stagedTree)
// Generate patch for each diff
for _, diff := range diffs {
    patch, err := diff.Patch()
    // Format as unified diff with 0 context lines
}
```

### 2. Diff Format and Context Lines

**Decision**: Use unified diff format with 0 lines of context to minimize token usage.

**Rationale**:
- Unified diff format is standard and matches `git diff --cached` output
- 0 context lines minimizes token usage while preserving essential change information
- AI models can often infer context from the changes themselves
- Consistent with token reduction requirements (FR-012)

**Alternatives Considered**:
- 3 lines of context (standard): Uses more tokens, not optimal for AI models
- Custom format: Non-standard, harder to parse, not compatible with standard tools
- Context-free (only changed lines): Already chosen (0 lines = context-free)

**Implementation Pattern**:
```go
// Use go-git's diff formatter with 0 context
formatter := &diff.UnifiedFormatter{
    Context: 0, // No context lines
}
// Format diff as unified patch
```

### 3. Large File/Diff Handling

**Decision**: Use 5000 character threshold - show metadata only for files/diffs exceeding threshold.

**Rationale**:
- 5000 characters provides reasonable balance between information and token usage
- Large files/diffs would consume excessive tokens without providing proportional value
- Metadata (file size, line count, change summary) still provides useful information to AI models
- Per-file limit (no aggregate) simplifies implementation and reasoning

**Alternatives Considered**:
- No size limits: Would allow token explosion for large files
- Smaller threshold (e.g., 1000 chars): Too restrictive, loses useful information
- Larger threshold (e.g., 10000 chars): Still allows token explosion
- Aggregate limit: More complex, unpredictable file exclusion

**Implementation Pattern**:
```go
const maxDiffSize = 5000 // characters

func computeDiffWithLimit(diff string) string {
    if len(diff) > maxDiffSize {
        // Return metadata only
        return generateMetadata(file, diff)
    }
    return diff
}
```

### 4. Binary File Detection

**Decision**: Use go-git's file type detection to identify binary files and set Diff to empty.

**Rationale**:
- go-git provides file type detection in diff results
- Binary files cannot be meaningfully diffed in text format
- Empty diff for binary files is consistent with FR-013
- No need for custom binary detection logic

**Alternatives Considered**:
- Show binary file indicators: Adds tokens without value, not needed per FR-013
- Attempt text diff of binary: Would produce garbage output
- Skip binary files entirely: Still need to include in FileChange list, just with empty diff

**Implementation Pattern**:
```go
if diff.IsBinary() {
    // Set Diff to empty string per FR-013
    fileChange.Diff = ""
    continue
}
```

### 5. Renamed/Copied File Handling

**Decision**: Use go-git's rename/copy detection with similarity percentage in diff output.

**Rationale**:
- go-git's diff API supports rename/copy detection with similarity calculation
- Rename/copy information is valuable for AI commit message generation
- Similarity percentage helps AI models understand the nature of the change
- Matches standard git behavior (`git diff --cached --find-renames`)

**Alternatives Considered**:
- Treat as delete+add: Loses rename/copy information
- Show diff content only: Loses rename/copy context
- Skip rename/copy detection: Less informative for AI models

**Implementation Pattern**:
```go
// go-git diff includes rename/copy information
if diff.RenameFrom != "" {
    // Format as "rename from X, rename to Y, similarity Z%"
    diffText := fmt.Sprintf("rename from %s, rename to %s, similarity %.0f%%",
        diff.RenameFrom, diff.RenameTo, diff.Similarity*100)
}
```

### 6. Unmerged Files Handling

**Decision**: Attempt to compute diff ignoring conflict markers, fallback to empty diff if computation fails.

**Rationale**:
- Unmerged files represent inconsistent state that shouldn't be committed
- Attempting diff computation provides best effort to extract useful information
- Fallback to empty diff prevents errors from blocking entire operation
- Consistent with FR-008 requirement

**Alternatives Considered**:
- Skip unmerged files entirely: Loses potential information
- Show conflict markers: Would confuse AI models, not useful
- Fail entire operation: Too strict, prevents processing other files

**Implementation Pattern**:
```go
if fileStatus.Staging == git.UpdatedButUnmerged {
    // Attempt diff computation (go-git may handle conflict markers)
    diff, err := computeDiff(file)
    if err != nil {
        // Set Diff to empty per FR-008
        fileChange.Diff = ""
        continue
    }
}
```

### 7. Empty Repository Handling

**Decision**: Treat all staged files as new file additions when HEAD doesn't exist.

**Rationale**:
- Empty repository (no HEAD) means no previous commit to diff against
- All staged files are effectively new file additions
- go-git diff API handles this case (diffing against empty tree)
- Consistent with FR-009 requirement

**Alternatives Considered**:
- Fail with error: Too strict, prevents initial commits
- Skip diff computation: Loses information about new files
- Use empty tree as base: Already handled by go-git

**Implementation Pattern**:
```go
head, err := r.repo.Head()
if err == plumbing.ErrReferenceNotFound {
    // No HEAD - treat as empty repository
    // All staged files are new additions
    // go-git diff will diff against empty tree
}
```

## Implementation Notes

- go-git's diff API requires converting index to tree for comparison
- Unified diff formatting with 0 context lines is supported via formatter options
- Character count includes the entire diff string (headers, file paths, change lines)
- Metadata generation should include: file size (bytes), line count, change summary (lines added/removed)
- Error handling must be per-file to prevent one failure from blocking others (FR-010)

## References

- go-git v5 documentation: https://pkg.go.dev/github.com/go-git/go-git/v5
- go-git diff API: https://pkg.go.dev/github.com/go-git/go-git/v5/plumbing/object#DiffTree
- Unified diff format: https://en.wikipedia.org/wiki/Diff#Unified_format
