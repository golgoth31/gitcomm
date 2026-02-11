# Contract: Error Display for Git Operations

**Feature**: 017-explicit-commit-errors  
**Date**: 2026-02-10  
**Scope**: `internal/repository/error_display.go`, `internal/repository/git_errors.go`, display call sites

## Interface Contract

### FormatErrorForDisplay

**Package**: `repository`

```go
func FormatErrorForDisplay(err error) string
```

**Purpose**: Format an error for user display, surfacing git stderr and applying truncation per FR-001–FR-005.

**Input**: Any error (typically wrapped chain from repository or service layer).

**Output**: Human-readable string suitable for `fmt.Print`/`fmt.Fprintf`. Never returns empty string for non-nil err.

**Behavior**:

1. **ErrGitCommandFailed** (via `errors.As`):
   - If `Stderr` non-empty: Truncate to 1500 chars; if truncated, append "… (N additional characters)" where N = len(Stderr) - 1500.
   - Format: `"git [Command] failed (exit [ExitCode]). Details: [Stderr or truncated Stderr]"`
   - If `Stderr` empty: `"git [Command] failed (exit [ExitCode]). No additional details from git. Check repository state or run the command manually."`

2. **ErrGitSigningFailed, ErrGitPermissionDenied, ErrGitFileNotFound** (via `errors.As`):
   - Use error message as Details: `"Brief from wrapper. Details: [error message]"` if a brief exists from the chain; otherwise use the error message as-is.

3. **Other errors** (fallback):
   - Return `err.Error()` unchanged (preserves current behavior for non-git errors).

**Truncation constant**: 1500 (package-level or same-file constant).

## ErrGitCommandFailed Contract Update

**File**: `internal/repository/git_errors.go`

**Change**: In `Error()` method, when `e.Stderr` is empty or only whitespace, use the generic hint instead of empty string:

```go
func (e *ErrGitCommandFailed) Error() string {
    detail := strings.TrimSpace(e.Stderr)
    if detail == "" {
        detail = "No additional details from git. Check repository state or run the command manually."
    }
    return fmt.Sprintf("git %s failed (exit %d): %s", e.Command, e.ExitCode, detail)
}
```

**Rationale**: Ensures the raw error (when used with `%v`) never ends with ": " alone. Display helper can still override format, but this fixes the root cause.

## Call Site Contract

**Files to update**:

| File | Line | Current | New |
|------|------|---------|-----|
| `internal/cmd/root.go` | 81 | `fmt.Fprintf(os.Stderr, "Error: failed to initialize git repository: %v\n", err)` | `fmt.Fprintf(os.Stderr, "Error: failed to initialize git repository: %s\n", repository.FormatErrorForDisplay(err))` |
| `internal/cmd/root.go` | 158 | `fmt.Fprintf(os.Stderr, "Error: commit failed: %v\n", commitErr)` | `fmt.Fprintf(os.Stderr, "Error: commit failed: %s\n", repository.FormatErrorForDisplay(commitErr))` |
| `internal/service/commit_service.go` | 556 | `fmt.Printf("\nError creating commit: %v\n", commitErr)` | `fmt.Printf("\nError creating commit: %s\n", repository.FormatErrorForDisplay(commitErr))` |

**Note**: Use `%s` with the formatted string; the formatter already includes the full message. The "Error: commit failed:" / "Error creating commit:" prefix is retained as the brief; the formatter produces the part after that (which may include "Details:").

**Clarification**: The spec says "Error creating commit: [brief]. Details: [content]". The call site keeps "Error creating commit: " and the formatter returns the rest (e.g., "git commit failed (exit 1). Details: [stderr]"). So the full output is: `"Error creating commit: git commit failed (exit 1). Details: [stderr]"`. That satisfies the prefixed structure.

## Test Requirements

- Unit test `FormatErrorForDisplay` with: ErrGitCommandFailed (empty stderr), ErrGitCommandFailed (non-empty stderr), ErrGitCommandFailed (stderr > 1500 chars), ErrGitSigningFailed, generic error, nil (must not be called with nil; caller responsibility).
- Unit test `ErrGitCommandFailed.Error()` with empty Stderr to assert it includes the hint.
- Integration: Trigger a commit failure (e.g., failing pre-commit hook) and assert output contains "Details:" and the hook output (or truncated version).
