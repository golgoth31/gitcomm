# Data Model: Compute Worktree Diff in GetRepositoryState

**Feature**: 014-worktree-diff
**Date**: 2025-01-27

## Entities

### RepositoryState

Represents the current state of the git repository for commit message generation. No changes to structure - only enhancement to populate Diff fields.

**Fields**:
- `StagedFiles []FileChange`: List of staged file changes (Diff field now populated)
- `UnstagedFiles []FileChange`: List of unstaged file changes (Diff field remains empty)

**Relationships**:
- Created by `GitRepository.GetRepositoryState()`
- Used by commit message generation services to provide context to AI models

**Validation Rules**:
- `StagedFiles` contains only files with staging status != Unmodified
- `UnstagedFiles` contains only files with worktree status != Unmodified
- `StagedFiles[].Diff` is populated for staged files (empty for unstaged files)
- `UnstagedFiles[].Diff` is always empty (FR-011)

**State Transitions**:
- Created: When GetRepositoryState() is called
- Populated: When diff computation completes for staged files
- Used: When passed to AI commit message generation

---

### FileChange

Represents a single file change in the repository. The `Diff` field already exists and is now populated for staged files.

**Fields**:
- `Path` (string): File path relative to repository root
- `Status` (string): Change status (added, modified, deleted, renamed, copied, unmerged)
- `Diff` (string): Unified diff content (populated for staged files only, empty for unstaged files)

**Relationships**:
- Contained in `RepositoryState.StagedFiles` or `RepositoryState.UnstagedFiles`
- Diff computed by `GitRepository.GetRepositoryState()` using go-git diff API

**Validation Rules**:
- `Path` must be a valid file path relative to repository root
- `Status` must be one of: "added", "modified", "deleted", "renamed", "copied", "unmerged", "unmodified"
- `Diff` is populated only for staged files (FR-011)
- `Diff` is empty for:
  - Unstaged files (FR-011)
  - Binary files (FR-013)
  - Files/diffs exceeding 5000 characters (shows metadata instead) (FR-016)
  - Files where diff computation failed (FR-010)
- `Diff` format: Unified diff format with 0 lines of context (FR-012)
- `Diff` content for renamed files: "rename from X, rename to Y, similarity Z%" (FR-014)
- `Diff` content for copied files: "copy from X, copy to Y, similarity Z%" (FR-015)

**State Transitions**:
- Created: When file status is detected in git status
- Diff computed: When GetRepositoryState() computes diff for staged files
- Diff populated: When diff content is assigned (or empty if conditions met)
- Used: When RepositoryState is passed to commit message generation

---

## Data Flow

```
1. GetRepositoryState() called
   ↓
2. worktree.Status() retrieves file status
   ↓
3. For each staged file:
   ├──> Check if binary → Diff = "" (FR-013)
   ├──> Check if unmerged → Attempt diff, fallback to "" (FR-008)
   ├──> Compute diff using go-git diff API
   │    ├──> Get HEAD tree
   │    ├──> Get staged index tree
   │    └──> Compute diff with 0 context lines
   ├──> Check diff size → If > 5000 chars, replace with metadata (FR-016)
   └──> Populate FileChange.Diff
   ↓
4. For each unstaged file:
   └──> Diff = "" (FR-011)
   ↓
5. RepositoryState returned with populated Diff fields
```

---

## Diff Content Format

### Standard Modified/Added/Deleted Files

Unified diff format with 0 lines of context:
```
diff --git a/path/to/file.go b/path/to/file.go
index abc123..def456 100644
--- a/path/to/file.go
+++ b/path/to/file.go
@@ -10,0 +11,3 @@
+line 1
+line 2
+line 3
```

### Renamed Files

```
rename from old/path.go
rename to new/path.go
similarity 95%
```

### Copied Files

```
copy from source/path.go
copy to dest/path.go
similarity 100%
```

### Large Files/Diffs (>5000 characters)

Metadata format (exact format TBD in implementation):
```
file: path/to/file.go
size: 10240 bytes
lines: 250
changes: +150 -100 (50 modified)
```

---

## Constraints

**Size Limits**:
- Per-file diff limit: 5000 characters (FR-016)
- No aggregate limit across all files (per clarification)
- Files/diffs exceeding limit show metadata only

**Content Rules**:
- Binary files: Diff = "" (FR-013)
- Unstaged files: Diff = "" (FR-011)
- Context lines: 0 (FR-012)
- Format: Unified diff (FR-012)

**Error Handling**:
- Individual file failures don't block other files (FR-010)
- Failed diffs result in Diff = "" with error logged
- Unmerged files: Attempt diff, fallback to "" if fails (FR-008)

---

## Relationships Diagram

```
GitRepository
    │
    └──> GetRepositoryState()
         │
         ├──> worktree.Status()
         │    └──> Returns file status map
         │
         ├──> go-git diff API
         │    ├──> Get HEAD tree
         │    ├──> Get staged index tree
         │    └──> Compute diff
         │
         └──> Creates
              └──> RepositoryState
                   │
                   ├──> StagedFiles []FileChange
                   │    └──> Diff populated (if conditions met)
                   │
                   └──> UnstagedFiles []FileChange
                        └──> Diff = "" (always)
```
