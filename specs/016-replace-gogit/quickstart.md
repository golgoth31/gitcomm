# Quickstart Guide: Replace Go-Git Library with External Git Commands

**Feature**: 016-replace-gogit
**Date**: 2026-02-10

## Overview

This migration replaces the `go-git` library with external `git` CLI commands in the `GitRepository` implementation. The interface remains unchanged; only the implementation behind it is rewritten.

## Prerequisites

- Go 1.25.0+
- Git 2.34.0+ installed and in PATH
- Existing tests passing on current branch

## Key Files

| File | Action | Description |
|------|--------|-------------|
| `internal/repository/git_repository_impl.go` | REWRITE | Main implementation (go-git → git CLI) |
| `internal/repository/git_errors.go` | CREATE | New categorized error types |
| `internal/repository/git_repository_impl_test.go` | UPDATE | Remove go-git test imports if any |
| `internal/repository/git_repository.go` | UNCHANGED | Interface stays the same |
| `pkg/git/config/extractor.go` | UPDATE | Remove `Signer` field, add `Enabled` to CommitSigner |
| `go.mod` / `go.sum` | UPDATE | Remove go-git and related deps |

## Implementation Order

### Step 1: Create error types (git_errors.go)

```go
package repository

import "errors"

var (
    ErrGitNotFound       = errors.New("git executable not found in PATH")
    ErrGitVersionTooOld  = errors.New("git version 2.34.0 or higher required")
    ErrGitPermissionDenied = errors.New("permission denied for git operation")
    ErrGitSigningFailed  = errors.New("git commit signing failed")
    ErrGitFileNotFound   = errors.New("file not found in git repository")
)

type ErrGitCommandFailed struct {
    Command  string
    Args     []string
    ExitCode int
    Stderr   string
}

func (e *ErrGitCommandFailed) Error() string {
    return fmt.Sprintf("git %s failed (exit %d): %s", e.Command, e.ExitCode, e.Stderr)
}
```

### Step 2: Rewrite gitRepositoryImpl struct

```go
type gitRepositoryImpl struct {
    path   string
    gitBin string
    config *gitconfig.GitConfig
    signer *gitconfig.CommitSigner
}
```

### Step 3: Implement execGit helper

```go
func (r *gitRepositoryImpl) execGit(ctx context.Context, args ...string) (string, string, error) {
    allArgs := append([]string{"-C", r.path}, args...)
    cmd := exec.CommandContext(ctx, r.gitBin, allArgs...)
    
    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    
    start := time.Now()
    err := cmd.Run()
    duration := time.Since(start)
    
    // Log execution (FR-018)
    utils.Logger.Debug().
        Str("command", args[0]).
        Strs("args", args[1:]).
        Int("exit_code", cmd.ProcessState.ExitCode()).
        Dur("duration", duration).
        Msg("git_exec")
    
    if err != nil {
        return stdout.String(), stderr.String(), r.categorizeError(args[0], args[1:], cmd.ProcessState.ExitCode(), stderr.String())
    }
    
    return stdout.String(), stderr.String(), nil
}
```

### Step 4: Implement each interface method

Implement in this order (dependency order):
1. `NewGitRepository` - constructor with git validation
2. `execGit` - central command execution helper
3. `parseStatus` - status output parser
4. `parseDiff` - diff output parser
5. `GetRepositoryState` - uses parseStatus + parseDiff
6. `CaptureStagingState` - uses parseStatus
7. `StageAllFiles` - simple `git add -A`
8. `StageModifiedFiles` - uses parseStatus + individual adds
9. `StageAllFilesIncludingUntracked` - uses parseStatus + individual adds
10. `UnstageFiles` - `git reset HEAD`
11. `CreateCommit` - with signing support

### Step 5: Update CommitSigner

In `pkg/git/config/extractor.go`, change:
```go
type CommitSigner struct {
    PrivateKeyPath string
    PublicKeyPath  string
    Format         string
    Enabled        bool    // replaces Signer interface{}
}
```

### Step 6: Remove go-git dependencies

```bash
# Remove direct dependencies
go get -d github.com/go-git/go-git/v5@none
go get -d github.com/hiddeco/sshsig@none
go get -d github.com/sergi/go-diff@none

# Tidy to remove transitives
go mod tidy
```

### Step 7: Run tests

```bash
go test ./internal/repository/... -v -count=1
go test ./... -v -count=1
```

## Critical Implementation Notes

1. **Binary file detection**: `git diff --cached` marks binary files with `Binary files differ` - detect this pattern and set empty diff
2. **Empty repository**: `git diff --cached` works on empty repos; `git status --porcelain=v1` also works
3. **Rename detection**: Already included in `git diff --cached` output with `rename from/to` headers
4. **Size limit**: Apply 5000 character limit per file after parsing diff output
5. **Context cancellation**: `exec.CommandContext` handles this - process is killed on context cancel
6. **No interactive prompts**: Do not set `GIT_TERMINAL_PROMPT=0` to allow SSH_ASKPASS
7. **Signing fallback**: If `git commit -S` fails with signing error, retry without `-S`
8. **Status parsing**: First char = staging status, second char = worktree status, `??` = untracked
9. **formattingService**: Keep the existing `formattingService` struct for commit message formatting (it doesn't use go-git)
10. **Metadata generation**: For oversized diffs, generate metadata using `os.Stat` and file reading (same as current)
