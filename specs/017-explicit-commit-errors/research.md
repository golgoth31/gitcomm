# Research: Explicit Commit Error Messages

**Feature**: 017-explicit-commit-errors  
**Date**: 2026-02-10  
**Purpose**: Resolve technical decisions for improving git error display to users

## Research Questions

### RQ1: Where are git errors displayed to users?

**Findings**:

| Location | Trigger | Current Display |
|----------|---------|-----------------|
| `internal/cmd/root.go:81` | Repository init failure | `Error: failed to initialize git repository: %v` |
| `internal/cmd/root.go:158` | Commit workflow failure | `Error: commit failed: %v` |
| `internal/service/commit_service.go:556` | AcceptAndCommit retry failure | `Error creating commit: %v` |
| `internal/service/commit_service.go:195` | AI generation fallback | `Error: %v` (different context, not git) |
| `internal/service/commit_service.go:89-90, 309-315` | Restoration warnings | Plain messages (not git stderr) |

**Decision**: Focus on root.go:158 and commit_service.go:556 for commit failures. Repository init (root.go:81) should also use explicit format. Staging/unstage errors flow through commit_service and may surface at root.go:158 or in restoration warnings.

**Rationale**: Commit failures are the primary user complaint. Repository init and staging failures that reach the user should use the same format.

---

### RQ2: How does ErrGitCommandFailed flow to the display?

**Findings**:

- `ErrGitCommandFailed` struct has `Command`, `Args`, `ExitCode`, `Stderr`
- `Error()` returns: `"git %s failed (exit %d): %s"` — when Stderr is empty, the trailing `%s` produces nothing → "git commit failed (exit 1):" with nothing after
- Errors are wrapped: `repository.CreateCommit` → `commit_service` → `root.go`
- Display uses `%v` which prints the error chain; the unwrapped message comes from the innermost error's `Error()` method

**Decision**: Two-pronged fix:
1. **Repository layer**: When `ErrGitCommandFailed` has empty Stderr, set a generic hint in the struct or in `Error()` so the message never ends with ": " alone
2. **Display layer**: Add a formatting helper that extracts git error details (ErrGitCommandFailed, ErrGitSigningFailed, etc.), applies truncation, and produces "Brief. Details: [content]" format

**Rationale**: Fixing only at repository layer would improve ErrGitCommandFailed but not the wrapper message structure. A display helper ensures consistent "Brief. Details:" format across all display points and handles truncation in one place.

---

### RQ3: Truncation implementation for long stderr

**Findings**:

- Spec: First 1500 characters + "… (N additional characters)" when truncated
- Go strings: `len(s)` for length; truncate with `s[:1500]` and append suffix
- Multi-byte: Go strings are UTF-8; slicing at 1500 may cut a rune in half. Use `utf8.ValidString` or slice conservatively. `s[:min(1500, len(s))]` is safe for ASCII; for UTF-8, iterate runes or use `utf8.ValidString(s[:i])` to find safe cut point.
- Simple approach: `s[:1500]` — if we hit mid-rune, the display might show a replacement char; acceptable for error output.

**Decision**: Truncate at 1500 bytes (not runes) for simplicity. If the 1500th byte is mid-rune, the rune may display as �; acceptable. Suffix: `"… (" + strconv.Itoa(len(s)-1500) + " additional characters)"`.

**Rationale**: Error output is typically ASCII. UTF-8 edge case is rare and low impact.

---

### RQ4: Staging and unstage error display paths

**Findings**:

- `StageAllFiles`, `StageModifiedFiles`, `StageAllFilesIncludingUntracked` return errors that propagate to `commit_service`
- Commit service calls these before prompting; on failure, `CreateCommit` returns early with wrapped error
- `UnstageFiles` is used in restoration; restoration failures show generic "Restoration timed out" or "failed to restore" — not the underlying git stderr
- `CaptureStagingState` and `GetRepositoryState` failures also propagate

**Decision**: Apply the explicit error formatter wherever a git-related error is displayed:
- root.go:81 (repo init)
- root.go:158 (commit/staging workflow failure)
- commit_service.go:556 (handleCommitFailure)
- commit_service.go:89, 309-315: Enhancement — when restoration fails, include the underlying git error detail if present (optional, lower priority)

**Rationale**: Same formatter used at all display points ensures consistent UX. Restoration messages can be enhanced in a follow-up if time permits.

---

### RQ5: Helper location and interface

**Findings**:

- Formatter is display/UX concern — format an error for user consumption
- `internal/utils` has `errors.go` (domain errors), `logger.go`
- `pkg/` is for shared utilities; `internal/` for app-private
- Formatter could live in `internal/utils` (e.g., `error_display.go`) or `internal/ui` (display package)

**Decision**: Create `internal/repository/error_display.go` with `FormatErrorForDisplay(err error) string`. This function:
- Uses `errors.As` to detect `*repository.ErrGitCommandFailed`
- Extracts Stderr, applies truncation (1500 chars + suffix)
- If Stderr empty, returns "Exit code N. No additional details from git. Check repository state or run the command manually."
- For other git errors (ErrGitSigningFailed, ErrGitPermissionDenied, etc.), returns their message as "Details:"
- Produces "Brief. Details: [content]" when content exists
- For non-git errors, falls back to `err.Error()` (no change in behavior)

**Rationale**: Repository package already defines git error types; placing the formatter there avoids import cycles (utils→repository would create a cycle). Cmd and service already import repository.

---

## Summary of Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Display points | root.go:81, 158; commit_service.go:556 | All user-facing git error paths |
| Fix strategy | Repository layer (empty stderr hint) + Display helper (format + truncation) | Defense in depth; consistent format |
| Truncation | 1500 bytes, suffix "… (N additional characters)" | Spec alignment; simple impl |
| Staging scope | Same formatter for all git errors at display | Consistent UX |
| Helper location | `internal/repository/error_display.go` | Avoids import cycles; cmd/service already import repository |
