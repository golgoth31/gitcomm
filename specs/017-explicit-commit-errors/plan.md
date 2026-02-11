# Implementation Plan: Explicit Commit Error Messages

**Branch**: `017-explicit-commit-errors` | **Date**: 2026-02-10 | **Spec**: [spec.md](./spec.md)  
**Input**: Feature specification from `/specs/017-explicit-commit-errors/spec.md`

## Summary

Display explicit, actionable error information when git operations (commit, stage, unstage, init) fail. The current message "Error creating commit: failed to create commit: git commit failed (exit 1):" provides no diagnostic detail. This plan adds a `FormatErrorForDisplay` helper that extracts git stderr, applies truncation (1500 chars + suffix), and produces a "Brief. Details: [content]" format at all user-facing error display points.

## Technical Context

**Language/Version**: Go 1.25.0  
**Primary Dependencies**: zerolog (logging), cobra (CLI), gcfg (git config)  
**Storage**: N/A  
**Testing**: Go `testing` package, table-driven tests  
**Target Platform**: Linux, macOS, Windows  
**Project Type**: Single CLI application  
**Performance Goals**: Negligible (error formatting is one-time per failure)  
**Constraints**: No secrets in error output; max 1500 visible chars per spec  
**Scale/Scope**: 3–4 display call sites; 1 new helper file; ~100 LOC

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Clean Architecture & Repository Pattern | ✓ PASS | Changes in repository (error display) and service/cmd (call sites); no layer violations |
| II. Interface-Driven Development & DI | ✓ PASS | No new interfaces; formatter is a pure function |
| III. Test-First Development | ✓ PASS | Unit tests for FormatErrorForDisplay and ErrGitCommandFailed.Error() |
| IV. Idiomatic Go Code Style | ✓ PASS | Standard Go, gofmt |
| V. Explicit Error Handling & Resource Management | ✓ PASS | Wrapped errors preserved; formatter reads error chain |
| VI. Context Propagation & Thread Safety | ✓ PASS | Not applicable (no goroutines) |
| Technical Constraints | ✓ PASS | No panics; no global state |
| Regulatory: No secrets in logs/errors | ✓ PASS | Git stderr assumed non-sensitive per spec |
| Logging via zerolog | ✓ PASS | Unchanged; no new logging requirements |

**Gate Result**: PASS — no violations

## Project Structure

### Documentation (this feature)

```text
specs/017-explicit-commit-errors/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
│   └── error-display-contract.md
└── tasks.md             # Phase 2 output (/speckit.tasks - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
internal/
├── repository/
│   ├── git_errors.go          # UPDATE: ErrGitCommandFailed.Error() for empty Stderr
│   ├── error_display.go       # CREATE: FormatErrorForDisplay
│   └── error_display_test.go  # CREATE: unit tests
├── cmd/
│   └── root.go               # UPDATE: use FormatErrorForDisplay at lines 81, 158
└── service/
    └── commit_service.go     # UPDATE: use FormatErrorForDisplay at line 556

test/
└── integration/               # Optional: commit failure integration test
```

**Structure Decision**: Single project. Changes are localized to `internal/repository`, `internal/cmd`, and `internal/service`. No new packages.

## Complexity Tracking

> No constitution violations. Section left empty.
