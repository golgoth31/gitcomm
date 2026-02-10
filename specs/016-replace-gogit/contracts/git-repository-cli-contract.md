# Contract: GitRepository CLI Implementation

**Feature**: 016-replace-gogit
**Date**: 2026-02-10
**Scope**: `internal/repository/git_repository_impl.go` (rewrite)

## Interface Contract (UNCHANGED)

The `GitRepository` interface remains identical. All consumers continue to work without modification.

```go
type GitRepository interface {
    GetRepositoryState(ctx context.Context) (*model.RepositoryState, error)
    CreateCommit(ctx context.Context, message *model.CommitMessage) error
    StageAllFiles(ctx context.Context) error
    CaptureStagingState(ctx context.Context) (*model.StagingState, error)
    StageModifiedFiles(ctx context.Context) (*model.AutoStagingResult, error)
    StageAllFilesIncludingUntracked(ctx context.Context) (*model.AutoStagingResult, error)
    UnstageFiles(ctx context.Context, files []string) error
}
```

## Constructor Contract

### NewGitRepository

```go
func NewGitRepository(repoPath string, noSign bool) (GitRepository, error)
```

**Behavior changes**:
- Now validates git executable exists via `exec.LookPath("git")`
- Now validates git version >= 2.34.0 via `git --version`
- No longer calls `git.PlainOpen()` (go-git removed)
- Still walks up directory tree to find `.git` directory
- Still extracts git config via `pkg/git/config`
- Still prepares `CommitSigner` (but without go-git Signer interface)

**New error cases**:
- `ErrGitNotFound` if git not in PATH
- `ErrGitVersionTooOld` if version < 2.34.0

**Signature**: Unchanged

## Internal Helper Contract

### execGit

```go
func (r *gitRepositoryImpl) execGit(ctx context.Context, args ...string) (stdout string, stderr string, err error)
```

**Purpose**: Central function for executing all git commands.

**Behavior**:
1. Builds command: `git -C <r.path> <args...>`
2. Uses `exec.CommandContext(ctx, r.gitBin, allArgs...)`
3. Captures both stdout and stderr via `bytes.Buffer`
4. Logs (via zerolog): command, arguments, exit code, execution duration
5. On error: categorizes into appropriate error type based on exit code and stderr
6. On context cancellation: process is killed automatically by `exec.CommandContext`

**Environment**:
- Inherits parent process env
- Does NOT set `GIT_TERMINAL_PROMPT=0` (to support SSH_ASKPASS)
- Sets `GIT_AUTHOR_NAME`, `GIT_AUTHOR_EMAIL`, `GIT_COMMITTER_NAME`, `GIT_COMMITTER_EMAIL` when committing

**Logging** (FR-018):
```
DEBUG git_exec command="status" args=["--porcelain=v1"] exit_code=0 duration_ms=12
DEBUG git_exec command="diff" args=["--cached","--unified=0"] exit_code=0 duration_ms=45
ERROR git_exec command="commit" args=["-S","-m","..."] exit_code=128 duration_ms=200 stderr="error: gpg failed..."
```

## Method Contracts

### GetRepositoryState

**Git commands**:
1. `git -C <path> status --porcelain=v1` → parse staged/unstaged files
2. `git -C <path> diff --cached --unified=0` → parse diffs for staged files

**Status parsing** (`parseStatus`):
```
Input:  "M  src/main.go\nA  src/new.go\n?? untracked.txt\n"
Output: StagedFiles:   [{Path:"src/main.go", Status:"modified"}, {Path:"src/new.go", Status:"added"}]
        UnstagedFiles: [{Path:"untracked.txt", Status:"added", Diff:""}]
```

**Diff parsing** (`parseDiff`):
```
Input:  "diff --git a/src/main.go b/src/main.go\nindex...\n@@ -1,3 +1,4 @@\n+new line\n\ndiff --git a/src/new.go..."
Output: map["src/main.go"] = "diff --git a/src/main.go...\n+new line\n"
        map["src/new.go"] = "diff --git a/src/new.go...\n..."
```

**includeNewFiles context handling**: Same behavior as current - filter `Added` files when `IncludeNewFilesKey` is false in context.

**Size limit**: Per-file diffs exceeding 5000 characters are replaced with metadata (file size, line count, change type).

### CreateCommit

**Git commands**:
1. With signing: `git -C <path> commit -S -m "<formatted_message>"` with env vars
2. Without signing: `git -C <path> commit -m "<formatted_message>"` with env vars
3. Signing fallback: if step 1 fails with signing error, retry as step 2

**Environment variables set on exec.Cmd**:
```go
cmd.Env = append(os.Environ(),
    "GIT_AUTHOR_NAME="+r.config.UserName,
    "GIT_AUTHOR_EMAIL="+r.config.UserEmail,
    "GIT_COMMITTER_NAME="+r.config.UserName,
    "GIT_COMMITTER_EMAIL="+r.config.UserEmail,
)
```

**Signing configuration** (passed via `-c` flags when signing):
```go
args = append(args, "-c", "gpg.format=ssh",
    "-c", "user.signingkey="+r.signer.PublicKeyPath,
    "-c", "commit.gpgsign=true")
```

**Signoff**: Appended to message string (same as current behavior).

### StageAllFiles

**Git command**: `git -C <path> add -A`

### CaptureStagingState

**Git command**: `git -C <path> status --porcelain=v1`
**Parse**: Extract file paths where first character (staging status) is not `?`, `!`, or ` `.

### StageModifiedFiles

**Git commands**:
1. `git -C <path> status --porcelain=v1` → filter modified files (Y != `?` and Y != ` `)
2. For each: `git -C <path> add -- <filepath>`
3. On any failure: `git -C <path> reset HEAD -- <staged_files...>` (rollback)

### StageAllFilesIncludingUntracked

**Git commands**:
1. `git -C <path> status --porcelain=v1` → filter all changed (Y != ` `)
2. For each: `git -C <path> add -- <filepath>`
3. On any failure: `git -C <path> reset HEAD -- <staged_files...>` (rollback)

### UnstageFiles

**Git command**: `git -C <path> reset HEAD -- <file1> <file2> ...`

**Note**: Current implementation already uses `git reset HEAD` as primary path with go-git `worktree.Remove` as fallback. New implementation only uses `git reset HEAD` (no fallback needed).

## Error Contract

All errors are wrapped with context using `fmt.Errorf("context: %w", err)`.

| Error | Condition | Wraps |
|-------|-----------|-------|
| `ErrGitNotFound` | `exec.LookPath("git")` fails | - |
| `ErrGitVersionTooOld` | Version < 2.34.0 | - |
| `utils.ErrNotGitRepository` | `.git` dir not found OR git says "not a git repository" | Existing |
| `ErrGitPermissionDenied` | stderr contains "permission denied" | `ErrGitCommandFailed` |
| `ErrGitSigningFailed` | stderr contains signing-related errors | `ErrGitCommandFailed` |
| `ErrGitFileNotFound` | stderr contains "pathspec" or "does not exist" | `ErrGitCommandFailed` |
| `ErrGitCommandFailed` | Any other non-zero exit code | stderr content |
| `utils.ErrStagingFailed` | Staging operation fails | Existing |
| `utils.ErrRestorationFailed` | Unstaging/rollback fails | Existing |
