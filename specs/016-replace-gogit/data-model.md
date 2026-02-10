# Data Model: Replace Go-Git Library with External Git Commands

**Feature**: 016-replace-gogit
**Date**: 2026-02-10

## Entities

### Existing Entities (UNCHANGED)

These entities are defined in `internal/model/` and remain unchanged by this migration. The new CLI-based implementation produces identical data.

#### RepositoryState (`internal/model/repository_state.go`)

```go
type RepositoryState struct {
    StagedFiles   []FileChange  // Files in staging area with diffs
    UnstagedFiles []FileChange  // Files with worktree changes (no diffs)
}

type FileChange struct {
    Path   string  // File path relative to repository root
    Status string  // "added", "modified", "deleted", "renamed", "copied", "unmerged"
    Diff   string  // Unified diff content (empty for unstaged, binary, or oversized)
}
```

#### CommitMessage (`internal/model/commit_message.go`)

```go
type CommitMessage struct {
    Type     string  // feat, fix, docs, style, refactor, test, chore, version
    Scope    string  // Optional scope
    Subject  string  // Short description, ≤72 characters
    Body     string  // Optional detailed explanation, ≤320 characters
    Footer   string  // Optional footer (issue refs, breaking changes)
    Signoff  bool    // Include "Signed-off-by" line
}
```

#### StagingState (`internal/model/staging_state.go`)

```go
type StagingState struct {
    StagedFiles    []string      // File paths that are staged
    CapturedAt     time.Time     // Timestamp of capture
    RepositoryPath string        // Repository root path
}

type AutoStagingResult struct {
    StagedFiles []string          // Successfully staged files
    FailedFiles []StagingFailure  // Files that failed to stage
    Success     bool              // Overall success
    Duration    time.Duration     // Operation duration
}

type StagingFailure struct {
    FilePath  string  // File that failed
    Error     error   // Error encountered
    ErrorType string  // "permission", "locked", "conflict", "other"
}
```

#### GitConfig (`pkg/git/config/extractor.go`)

```go
type GitConfig struct {
    UserName      string  // user.name
    UserEmail     string  // user.email
    SigningKey    string  // user.signingkey
    GPGFormat     string  // gpg.format (gpg, ssh)
    CommitGPGSign bool    // commit.gpgsign
}
```

### Modified Entity

#### CommitSigner (`pkg/git/config/extractor.go`)

The `Signer` field currently holds a `git.Signer` interface from go-git. After migration, this field is no longer needed since signing is delegated to git CLI.

**Before** (current):
```go
type CommitSigner struct {
    PrivateKeyPath string
    PublicKeyPath  string
    Format         string
    Signer         interface{}  // go-git git.Signer
}
```

**After** (migrated):
```go
type CommitSigner struct {
    PrivateKeyPath string  // Path to private key (for env var setup)
    PublicKeyPath  string  // Path to public key (user.signingkey)
    Format         string  // Signing format ("ssh", "gpg")
    Enabled        bool    // Whether signing is enabled
}
```

Changes:
- Removed `Signer interface{}` field (no longer need go-git Signer)
- Added `Enabled bool` field to track signing enablement (replaces nil check on Signer)

### New Entity

#### gitRepositoryImpl (REWRITTEN) (`internal/repository/git_repository_impl.go`)

**Before** (current):
```go
type gitRepositoryImpl struct {
    repo   *git.Repository           // go-git Repository object
    path   string                    // Repository root path
    config *gitconfig.GitConfig      // Git configuration
    signer *gitconfig.CommitSigner   // Commit signer
}
```

**After** (migrated):
```go
type gitRepositoryImpl struct {
    path    string                   // Repository root path
    gitBin  string                   // Resolved path to git executable
    config  *gitconfig.GitConfig     // Git configuration
    signer  *gitconfig.CommitSigner  // Commit signer (no Signer interface)
}
```

Changes:
- Removed `repo *git.Repository` (no longer using go-git)
- Added `gitBin string` for resolved git executable path (from `exec.LookPath`)

### New Entity: Git Error Types (`internal/repository/git_errors.go`)

```go
// ErrGitNotFound indicates git executable is not in PATH
var ErrGitNotFound = errors.New("git executable not found in PATH")

// ErrGitVersionTooOld indicates git version < 2.34.0
var ErrGitVersionTooOld = errors.New("git version 2.34.0 or higher required")

// ErrGitPermissionDenied indicates a filesystem permission error
var ErrGitPermissionDenied = errors.New("permission denied for git operation")

// ErrGitSigningFailed indicates commit signing failed
var ErrGitSigningFailed = errors.New("git commit signing failed")

// ErrGitFileNotFound indicates a file was not found in the repository
var ErrGitFileNotFound = errors.New("file not found in git repository")

// ErrGitCommandFailed is a generic error for git command failures
type ErrGitCommandFailed struct {
    Command  string   // Git subcommand (e.g., "status", "commit", "add")
    Args     []string // Command arguments
    ExitCode int      // Process exit code
    Stderr   string   // Captured stderr output
}
```

## Entity Relationships

```text
NewGitRepository()
    ├──> validates git executable exists (ErrGitNotFound)
    ├──> validates git version >= 2.34.0 (ErrGitVersionTooOld)
    ├──> extracts GitConfig via pkg/git/config
    ├──> prepares CommitSigner (Enabled flag)
    └──> creates gitRepositoryImpl{path, gitBin, config, signer}

gitRepositoryImpl
    ├──> GetRepositoryState()
    │     ├──> execGit("status", "--porcelain=v1") → parse → FileChange[]
    │     ├──> execGit("diff", "--cached", "--unified=0") → parse → Diff per file
    │     └──> produces RepositoryState{StagedFiles, UnstagedFiles}
    │
    ├──> CreateCommit()
    │     ├──> formats CommitMessage → string
    │     ├──> sets GIT_AUTHOR_NAME/EMAIL env from GitConfig
    │     ├──> if signer.Enabled: execGit("commit", "-S", "-m", msg)
    │     │     └──> on signing failure: retry without -S (ErrGitSigningFailed)
    │     └──> execGit("commit", "-m", msg)
    │
    ├──> StageAllFiles()
    │     └──> execGit("add", "-A")
    │
    ├──> CaptureStagingState()
    │     ├──> execGit("status", "--porcelain=v1") → parse staged files
    │     └──> produces StagingState{StagedFiles, CapturedAt, RepositoryPath}
    │
    ├──> StageModifiedFiles()
    │     ├──> execGit("status", "--porcelain=v1") → filter modified (not untracked)
    │     ├──> for each: execGit("add", "--", filepath)
    │     ├──> on failure: rollback via execGit("reset", "HEAD", "--", staged...)
    │     └──> produces AutoStagingResult
    │
    ├──> StageAllFilesIncludingUntracked()
    │     ├──> execGit("status", "--porcelain=v1") → filter all changed
    │     ├──> for each: execGit("add", "--", filepath)
    │     ├──> on failure: rollback via execGit("reset", "HEAD", "--", staged...)
    │     └──> produces AutoStagingResult
    │
    └──> UnstageFiles()
          └──> execGit("reset", "HEAD", "--", files...)
```

## Data Flow

```text
User Input (CommitMessage)
    │
    ▼
commit_service.go (uses GitRepository interface)
    │
    ▼
gitRepositoryImpl (NEW: git CLI execution)
    │
    ├──> execGit() helper
    │     ├──> builds exec.CommandContext
    │     ├──> sets env vars (author, signing)
    │     ├──> runs git -C <path> <subcommand> <args...>
    │     ├──> logs: command, args, exit code, duration (FR-018)
    │     ├──> captures stdout + stderr
    │     └──> categorizes errors (FR-006)
    │
    ├──> parseStatus() helper
    │     ├──> parses porcelain v1 output
    │     └──> returns []FileChange with status codes
    │
    ├──> parseDiff() helper
    │     ├──> splits diff output on "diff --git" boundaries
    │     ├──> detects binary files
    │     ├──> applies 5000 char size limit per file
    │     └──> returns map[filepath]string
    │
    └──> Output (RepositoryState / StagingState / error)
```

## Validation Rules

- Git path must resolve to valid repository (contains `.git` directory)
- Git version must be >= 2.34.0 (parsed from `git --version`)
- Git executable must be resolvable via `exec.LookPath("git")`
- Status output lines must match `^.{2} .+$` pattern (2 chars + space + path)
- Diff size per file must not exceed 5000 characters (replaced with metadata)
- Commit message must not be empty (validated by caller, but git also rejects)
