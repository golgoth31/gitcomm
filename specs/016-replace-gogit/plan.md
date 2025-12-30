# Implementation Plan: Replace Go-Git Library with External Git Commands

**Branch**: `016-replace-gogit` | **Date**: 2026-02-10 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/016-replace-gogit/spec.md`

## Summary

Replace all go-git library (`github.com/go-git/go-git/v5`) usage in the `internal/repository/git_repository_impl.go` with external `git` CLI commands executed via `os/exec`. The existing `GitRepository` interface remains unchanged. All 7 interface methods (`GetRepositoryState`, `CreateCommit`, `StageAllFiles`, `CaptureStagingState`, `StageModifiedFiles`, `StageAllFilesIncludingUntracked`, `UnstageFiles`) will be reimplemented using `git` CLI equivalents. The go-git dependency and its transitive dependencies (go-billy, go-git/gcfg indirect) will be completely removed from `go.mod`. SSH commit signing will be delegated to git CLI natively. Git version 2.34.0+ is enforced.

## Technical Context

**Language/Version**: Go 1.25.0
**Primary Dependencies**:
- Existing: `os/exec` (stdlib), `github.com/rs/zerolog` (logging), `github.com/spf13/cobra` (CLI)
- Existing: `github.com/go-git/gcfg/v2` (git config parsing - retained for INI parsing)
- Removed: `github.com/go-git/go-git/v5`, `github.com/hiddeco/sshsig`, `github.com/sergi/go-diff`, `golang.org/x/crypto/ssh` (signing handled by git CLI)
**Storage**: N/A (filesystem-based git operations)
**Testing**: Go `testing` package, table-driven tests, `t.Parallel()`, git CLI for test fixture setup (already used in existing tests)
**Target Platform**: Linux, macOS, Windows (wherever git 2.34.0+ is available)
**Project Type**: Single CLI application
**Performance Goals**: Git operations complete within 2x of go-git implementation (accounting for process spawn overhead)
**Constraints**: Requires git 2.34.0+ in PATH; no interactive prompts; SSH passphrase via env var or config
**Scale/Scope**: 1 implementation file rewrite (~1000 lines), 1 test file update, go.mod/go.sum cleanup, config package minor update

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Clean Architecture & Repository Pattern | PASS | GitRepository interface preserved; impl swapped behind same interface |
| II. Interface-Driven Development & DI | PASS | Interface unchanged; implementation injected via `NewGitRepository` constructor |
| III. Test-First Development | PASS | Existing tests use git CLI for setup, will work with new impl; new error types get tests |
| IV. Idiomatic Go Code Style | PASS | Will follow naming conventions, error wrapping, small focused functions |
| V. Explicit Error Handling | PASS | Distinct error types for git failures; wrapped errors with `%w` |
| VI. Context Propagation & Thread Safety | PASS | `exec.CommandContext` used for all git commands; context cancellation kills processes |
| Technical: No global state | PASS | All state in struct fields, injected via constructor |
| Technical: context.Context for cancellation | PASS | All methods already accept context; will use `exec.CommandContext` |
| Operational: No secrets in logs | PASS | Git command logging will exclude sensitive data (SSH keys, passphrases) |
| Operational: Logging via zerolog | PASS | FR-018 requires logging all git executions via existing zerolog logger |

**Gate Result**: PASS - No violations. Proceed to Phase 0.

## Project Structure

### Documentation (this feature)

```text
specs/016-replace-gogit/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
└── tasks.md             # Phase 2 output (created by /speckit.tasks)
```

### Source Code (repository root)

```text
internal/
├── repository/
│   ├── git_repository.go          # Interface (UNCHANGED)
│   ├── git_repository_impl.go     # Implementation (REWRITTEN - go-git → git CLI)
│   ├── git_repository_impl_test.go # Tests (UPDATED - remove go-git test helpers if any)
│   └── git_errors.go              # NEW - categorized git CLI error types
├── model/
│   ├── repository_state.go        # UNCHANGED
│   ├── staging_state.go           # UNCHANGED
│   └── commit_message.go          # UNCHANGED
├── utils/
│   └── errors.go                  # UPDATED - add git version error
└── service/
    └── commit_service.go          # UNCHANGED (uses interface)

pkg/
└── git/
    └── config/
        ├── extractor.go           # UNCHANGED (uses gcfg, not go-git)
        └── errors.go              # UNCHANGED
```

**Structure Decision**: Single project structure (existing). Only `internal/repository/` is modified. The `GitRepository` interface remains unchanged, so all consumers (`internal/service/`) are unaffected.

## Complexity Tracking

No constitution violations to justify.
