# Quickstart: Respect addAll Flag for New Files

**Feature**: 005-addall-new-files
**Date**: 2025-01-27

## Overview

This feature ensures that the `-a` (add-all) flag controls whether new/untracked files are included in the commit workflow. When the flag is false, new files are excluded from repository state, allowing users to commit only changes to existing files.

## User Experience

### Before (Current Behavior)

```bash
# Create a new file
echo "new content" > newfile.txt

# Run gitcomm without -a flag
gitcomm

# Problem: newfile.txt is included in commit even though -a flag was not used
```

### After (New Behavior)

```bash
# Create a new file
echo "new content" > newfile.txt

# Run gitcomm without -a flag
gitcomm

# Result: newfile.txt is excluded from commit, only modified tracked files included

# Run gitcomm with -a flag
gitcomm -a

# Result: newfile.txt is included in commit along with all other files
```

## Implementation Flow

### 1. Service Layer (CommitService)

**Location**: `internal/service/commit_service.go`

**Change**: Pass `addAll` flag via context to repository layer

```go
// Extract AutoStage value from options
useAllFiles := s.options != nil && s.options.AutoStage

// Set context value for repository filtering
ctx = context.WithValue(ctx, includeNewFilesKey, useAllFiles)

// Call GetRepositoryState with context
state, err := s.gitRepo.GetRepositoryState(ctx)
```

### 2. Repository Layer (GitRepository)

**Location**: `internal/repository/git_repository_impl.go`

**Change**: Filter new files based on context value

```go
// Extract includeNewFiles from context (default: true)
includeNewFiles := true
if val := ctx.Value(includeNewFilesKey); val != nil {
    includeNewFiles = val.(bool)
}

// Filter during iteration
for file, fileStatus := range status {
    if fileStatus.Staging != git.Unmodified {
        // Skip new files when addAll is false
        if fileStatus.Staging == git.Added && !includeNewFiles {
            continue // Skip this new file
        }

        // Process file normally (compute diff, etc.)
        // ...
    }
}
```

## Key Changes

### Files Modified

1. **`internal/service/commit_service.go`**:
   - Add context value setting before `GetRepositoryState` call
   - Pass `AutoStage` flag via context

2. **`internal/repository/git_repository_impl.go`**:
   - Add context key definition
   - Extract context value in `GetRepositoryState`
   - Add filtering logic for new files

### Files Unchanged

- `internal/model/repository_state.go` - No model changes
- `internal/repository/git_repository.go` - Interface unchanged
- `cmd/gitcomm/main.go` - CLI unchanged (flag already exists)

## Testing

### Unit Tests

**Location**: `internal/repository/git_repository_impl_test.go`

**Test Cases**:
1. `TestGetRepositoryState_ExcludesNewFilesWhenAddAllFalse` - Verify new files excluded
2. `TestGetRepositoryState_IncludesNewFilesWhenAddAllTrue` - Verify new files included
3. `TestGetRepositoryState_IncludesModifiedFilesAlways` - Verify modified files always included
4. `TestGetRepositoryState_DefaultBehaviorIncludesAll` - Verify backward compatibility

### Integration Tests

**Location**: `test/integration/` (if applicable)

**Test Cases**:
1. End-to-end: `gitcomm` without `-a` excludes new files
2. End-to-end: `gitcomm -a` includes new files
3. Edge case: Manually staged new file excluded when `-a` not used

## Verification

### Manual Testing Steps

1. **Setup**:
   ```bash
   # Create test repository
   mkdir test-repo && cd test-repo
   git init
   echo "existing" > existing.txt
   git add existing.txt
   git commit -m "Initial commit"
   ```

2. **Test Exclude New Files**:
   ```bash
   # Modify existing file
   echo "modified" >> existing.txt

   # Create new file
   echo "new" > newfile.txt

   # Run without -a flag
   gitcomm

   # Verify: Only existing.txt should be in commit, newfile.txt excluded
   ```

3. **Test Include New Files**:
   ```bash
   # Run with -a flag
   gitcomm -a

   # Verify: Both existing.txt and newfile.txt should be in commit
   ```

## Success Criteria

- ✅ New files excluded when `-a` flag not used
- ✅ New files included when `-a` flag used
- ✅ Modified files always included regardless of flag
- ✅ No performance regression
- ✅ Backward compatibility maintained (default behavior unchanged)
