# Data Model: Explicit Commit Error Messages

**Feature**: 017-explicit-commit-errors  
**Date**: 2026-02-10

## Entities

### Existing (unchanged)

- **ErrGitCommandFailed** (`internal/repository/git_errors.go`): Struct with `Command`, `Args`, `ExitCode`, `Stderr` — no schema change
- **ErrGitSigningFailed**, **ErrGitPermissionDenied**, etc.: Sentinel errors — unchanged
- **CommitService**, **GitRepository**: Interfaces unchanged

### Modified Behavior

- **ErrGitCommandFailed.Error()**: When `Stderr` is empty, MUST NOT return a string ending with `": "` and nothing after. Instead, include the generic hint: `"No additional details from git. Check repository state or run the command manually."` as the trailing part.

### New

- **FormatErrorForDisplay** (function in `internal/repository/error_display.go`): Accepts `error`, returns `string` formatted for user display. See contracts.

## Entity Relationships

```text
User sees error
    │
    ▼
root.go / commit_service.go
    │
    ▼
FormatErrorForDisplay(err)
    │
    ├── errors.As(err, &ErrGitCommandFailed) → extract Stderr, truncate, format
    ├── errors.As(err, &ErrGitSigningFailed)  → use as Details
    └── fallback → err.Error()
    │
    ▼
"Brief. Details: [content]" or "Brief. [hint if empty]"
```

## Validation Rules

- Truncation: When content length > 1500, visible part ≤ 1500 characters; suffix = "… (N additional characters)" where N = len(content) - 1500
- Empty stderr: Display MUST include exit code and generic hint
- Format: "Brief. Details: [content]" when content exists; "Brief. [hint]" when empty
