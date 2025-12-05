# Research: Respect addAll Flag for New Files

**Feature**: 005-addall-new-files
**Date**: 2025-01-27

## Technology Decisions

### 1. Identifying New vs Modified Files

**Decision**: Use `fileStatus.Staging == git.Added` to identify new files that should be filtered when `addAll` is false.

**Rationale**:
- `git.Added` status in go-git indicates a file is staged and doesn't exist in HEAD (new/untracked file)
- `git.Modified` status indicates a file exists in HEAD and has changes (tracked file)
- `git.Deleted` and `git.Renamed` are tracked file operations
- This approach leverages existing git status information without additional HEAD lookups

**Alternatives Considered**:
- **Check if file exists in HEAD tree**: More accurate but requires additional tree lookups for each file, adding overhead
- **Use worktree status**: `fileStatus.Worktree == git.Untracked` indicates untracked files, but after staging, worktree status may change
- **Track staging method used**: Remember which staging method was called, but this doesn't handle manually staged files

**Implementation Pattern**:
```go
// In GetRepositoryState, when building StagedFiles:
if fileStatus.Staging == git.Added && !includeNewFiles {
    // Skip this file - it's a new file and addAll is false
    continue
}
```

### 2. Filtering Location

**Decision**: Filter new files in `GetRepositoryState` implementation by checking staging status, maintaining interface compatibility.

**Rationale**:
- `GetRepositoryState` already has access to `fileStatus.Staging` which indicates file type
- Keeps filtering logic in repository layer where file status knowledge exists
- Maintains `GitRepository` interface compatibility (no signature changes)
- Service layer can pass `includeNewFiles` flag via context or new optional parameter

**Alternatives Considered**:
- **Filter in service layer**: Would require exposing filtering logic outside repository, violates separation of concerns
- **Extend interface with parameter**: `GetRepositoryState(ctx, includeNewFiles bool)` - breaks interface compatibility, requires all callers to update
- **Use context values**: Pass flag via `context.Context` - maintains interface but uses context for configuration (acceptable pattern)

**Implementation Pattern**:
```go
// Option A: Context value (preferred for interface compatibility)
type contextKey string
const includeNewFilesKey contextKey = "includeNewFiles"

// In service:
ctx = context.WithValue(ctx, includeNewFilesKey, s.options.AutoStage)
state, err := s.gitRepo.GetRepositoryState(ctx)

// In repository:
includeNewFiles := true // default
if val := ctx.Value(includeNewFilesKey); val != nil {
    includeNewFiles = val.(bool)
}
```

### 3. Handling Edge Cases

**Decision**: Apply consistent filtering rules for all edge cases:
- Manually staged new files: Filtered when `addAll` is false (FR-001)
- Binary new files: Filtered same as text files (FR-001)
- Renamed files: Never filtered (treated as tracked, FR-005)
- Files staged then unstaged: Not in staging area, so not included anyway

**Rationale**:
- Consistent behavior regardless of how file was staged
- Renamed files have `git.Renamed` status, not `git.Added`, so automatically excluded from filtering
- Binary vs text distinction doesn't affect new vs modified distinction

**Edge Case Handling**:
1. **Manually staged new file + `addAll` false**: File has `git.Added` status, filtered out ✅
2. **New binary file + `addAll` false**: File has `git.Added` status, filtered out ✅
3. **Renamed file**: Has `git.Renamed` status, not filtered (not `git.Added`) ✅
4. **New file staged then unstaged**: Not in `fileStatus.Staging`, not included in results ✅

### 4. Performance Considerations

**Decision**: Filter during iteration over status results, adding minimal overhead (single boolean check per file).

**Rationale**:
- Filtering is O(1) per file - just a status comparison
- No additional git operations or tree lookups needed
- Maintains existing performance characteristics
- Typical repositories have <100 files, so overhead is negligible

**Performance Impact**:
- Additional check: `if fileStatus.Staging == git.Added && !includeNewFiles { continue }`
- Time complexity: O(1) per file, O(n) total where n = number of files
- Memory: No additional allocations
- Expected overhead: <1ms for typical repositories

## Integration Points

### Repository Layer
- `GetRepositoryState` method in `git_repository_impl.go`
- Uses existing `fileStatus.Staging` from `worktree.Status()`
- No changes to `GitRepository` interface

### Service Layer
- `CommitService.CreateCommit` method
- Passes `addAll` flag via context to repository
- Maintains existing service interface

### Model Layer
- No changes to `RepositoryState` or `FileChange` models
- Existing models already support the required data

## Dependencies

- **Existing**: `github.com/go-git/go-git/v5` - Already used for git operations
- **Existing**: `context.Context` - Already used for cancellation and request-scoped values
- **No new dependencies required**

## Testing Strategy

- **Unit tests**: Test filtering logic with various file status combinations
- **Table-driven tests**: Cover all combinations of:
  - `addAll` true/false
  - File statuses: Added, Modified, Deleted, Renamed
  - Binary vs text files
  - Manually staged vs auto-staged files
- **Integration tests**: Verify end-to-end behavior with real git repositories

## Open Questions Resolved

1. **Q: How to identify new files?** → A: Use `fileStatus.Staging == git.Added`
2. **Q: Where to implement filtering?** → A: In `GetRepositoryState` using context value
3. **Q: How to pass flag to repository?** → A: Via `context.Context` to maintain interface compatibility
4. **Q: What about renamed files?** → A: They have `git.Renamed` status, not `git.Added`, so automatically excluded

All questions resolved - no NEEDS CLARIFICATION markers remain.
