# Quickstart: Explicit Commit Error Messages

**Feature**: 017-explicit-commit-errors  
**Date**: 2026-02-10

## Overview

Improve user-facing error messages when git operations (commit, stage, unstage, init) fail. Display git stderr with truncation and a consistent "Brief. Details:" format.

## Prerequisites

- Go 1.25.0+
- Existing tests passing

## Key Files

| File | Action |
|------|--------|
| `internal/repository/error_display.go` | CREATE — `FormatErrorForDisplay(err error) string` |
| `internal/repository/git_errors.go` | UPDATE — `ErrGitCommandFailed.Error()` handle empty Stderr |
| `internal/cmd/root.go` | UPDATE — use `repository.FormatErrorForDisplay` at lines 81, 158 |
| `internal/service/commit_service.go` | UPDATE — use `repository.FormatErrorForDisplay` at line 556 |
| `internal/repository/error_display_test.go` | CREATE — unit tests for formatter |

## Implementation Order

### Step 1: Fix ErrGitCommandFailed for empty Stderr

In `git_errors.go`, update `Error()` so that when `Stderr` is empty, the message ends with the generic hint instead of ": ".

### Step 2: Create FormatErrorForDisplay

In `internal/repository/error_display.go`:

1. Define truncation constant (1500).
2. Implement `FormatErrorForDisplay(err error) string`:
   - Use `errors.As` to detect `*repository.ErrGitCommandFailed`
   - Extract Stderr; if non-empty, truncate and append suffix when needed
   - If empty, return message with generic hint (or rely on updated ErrGitCommandFailed)
   - For other git errors, format with "Details:" prefix where applicable
   - Fallback: `return err.Error()`
3. Handle nil: Return "" or treat as caller responsibility (document that nil should not be passed).

### Step 3: Update display call sites

- root.go:81 — `repository.FormatErrorForDisplay(err)` for repo init
- root.go:158 — `repository.FormatErrorForDisplay(commitErr)` for commit failure
- commit_service.go:556 — `repository.FormatErrorForDisplay(commitErr)` for handleCommitFailure

### Step 4: Add unit tests

- Test FormatErrorForDisplay with various error types
- Test ErrGitCommandFailed.Error() with empty Stderr
- Test truncation at 1500 chars

### Step 5: Manual verification

- Create a repo with a failing pre-commit hook; run gitcomm; confirm output includes hook stderr.
- Simulate empty stderr (e.g., mock) and confirm generic hint appears.

## Critical Implementation Notes

1. **Truncation**: First 1500 characters only; suffix "… (N additional characters)" when truncated.
2. **Empty stderr**: Both ErrGitCommandFailed.Error() and FormatErrorForDisplay must handle it.
3. **Format**: Output should read naturally: "Error: commit failed: git commit failed (exit 1). Details: [hook output]"
4. **Import cycle**: `internal/utils` must not import `internal/repository` if that creates a cycle. Check: repository imports utils (logger); utils importing repository for ErrGitCommandFailed type would create cycle. **Resolution**: Use `errors.As` with interface or define the struct shape in utils. Prefer: keep formatter in utils and use `errors.As(err, &repoErr)` — but that requires importing repository. Alternative: move formatter to a `pkg/display` or `internal/display` package that imports both. Simpler: put `FormatGitErrorForDisplay` in `internal/repository` (it formats repo errors) and have cmd/service import it. Or: use a small interface in utils — `type GitError interface { Command() string; ExitCode() int; Stderr() string }` — but that changes the repository package. Easiest: place the formatter in `internal/repository` as `FormatErrorForDisplay(err error) string` in a new file `error_display.go` in the repository package. Then cmd and service import repository and call `repository.FormatErrorForDisplay(err)`. No cycle.

**Revised**: Create `internal/repository/error_display.go` with `FormatErrorForDisplay(err error) string`. Repository package already has git error types; no cycle. Cmd and service already import repository.
