# Implementation Plan: Compute Worktree Diff in GetRepositoryState

**Branch**: `014-worktree-diff` | **Date**: 2025-01-27 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/014-worktree-diff/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Enhance the `GetRepositoryState` function to compute unified diff (patch format) for all staged files using go-git plumbing diff capabilities. The diff represents changes between the current worktree state (with staged changes) and the clean worktree state (HEAD). To minimize token usage for AI models, diffs use 0 lines of context, and files/diffs exceeding 5000 characters show only metadata (file size, line count, change summary) instead of full content.

## Technical Context

**Language/Version**: Go 1.25.0+
**Primary Dependencies**:
- Existing: `github.com/go-git/go-git/v5` (v5.16.4) - git repository operations and diff computation
- Existing: `github.com/rs/zerolog` (v1.34.0) - debug logging
- No new external dependencies required
**Storage**: N/A (in-memory diff computation, no persistence)
**Testing**: Go `testing` package with table-driven tests, existing test infrastructure
**Target Platform**: Linux/macOS/Windows (CLI tool)
**Project Type**: Single CLI application
**Performance Goals**:
  - Diff computation completes for up to 100 staged files in under 2 seconds (SC-003)
  - Error rate < 1% of files (SC-004)
**Constraints**:
  - Must compute diff only for staged files (not unstaged) (FR-011)
  - Must use 0 lines of context to minimize token usage (FR-012)
  - Must limit diff size to 5000 characters per file (FR-016)
  - Must handle binary files gracefully (empty diff) (FR-013)
  - Must handle unmerged files (attempt diff ignoring conflict markers) (FR-008)
  - Must handle empty repository (no HEAD) (FR-009)
  - Must maintain backward compatibility with existing GitRepository interface
  - Must not break existing GetRepositoryState workflow
**Scale/Scope**:
  - Single diff computation per GetRepositoryState call
  - Handles typical git repositories (up to 100 staged files)
  - Supports all file status types: modified, added, deleted, renamed, copied
  - Token optimization: 0 context lines, 5000 char threshold per file

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Verify compliance with gitcomm Constitution principles:

- **Clean Architecture**: ✅ **COMPLIANT**
  - Diff computation logic will be added to existing `internal/repository/` layer
  - Uses existing `internal/model/` for RepositoryState and FileChange
  - Clear separation: diff computation (internal/repository) → model usage (internal/model)
  - No new layers required, fits existing structure

- **Interface-Driven Development**: ✅ **COMPLIANT**
  - Extends existing `GitRepository` interface (no changes needed)
  - Uses go-git interfaces for diff computation
  - Dependencies injected via existing constructor pattern
  - No global state

- **Test-First Development**: ✅ **COMPLIANT**
  - TDD approach: tests for diff computation first
  - Table-driven tests for various file status types and edge cases
  - Integration tests for GetRepositoryState with diff computation
  - Unit tests for diff size limiting and metadata generation

- **Idiomatic Go**: ✅ **COMPLIANT**
  - Follows Go naming conventions
  - Small, focused functions for diff computation
  - Proper error handling

- **Error Handling**: ✅ **COMPLIANT**
  - Explicit error handling for diff computation failures
  - Wrapped errors for traceability
  - Graceful degradation: continue processing other files on individual file failures (FR-010)

- **Context & Thread Safety**: ✅ **COMPLIANT**
  - Uses context.Context for cancellation (already in GetRepositoryState signature)
  - Diff computation is stateless (no shared mutable state)
  - Thread-safe file operations
  - No goroutines required

- **Technical Constraints**: ✅ **COMPLIANT**
  - No global state
  - Stateless diff computation
  - Resource cleanup handled by go-git library

- **Operational Constraints**: ✅ **COMPLIANT**
  - Debug logging via existing zerolog infrastructure
  - No secrets involved
  - Error messages don't expose internal details

**Violations**: None. All principles are satisfied.

## Project Structure

### Documentation (this feature)

```text
specs/014-worktree-diff/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
internal/
├── repository/
│   ├── git_repository.go           # GitRepository interface (no changes)
│   └── git_repository_impl.go      # GetRepositoryState enhancement (add diff computation)
└── model/
    └── repository_state.go         # RepositoryState, FileChange models (Diff field exists)

test/
└── integration/
    └── git_repository_test.go      # Integration tests for diff computation
```

**Structure Decision**: Single project structure. Diff computation logic is added to existing `internal/repository/git_repository_impl.go` in the `GetRepositoryState` method. No new packages or modules required. The existing `FileChange` model already has a `Diff` field, so no model changes needed.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

No violations. All principles are satisfied.
