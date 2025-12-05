# Implementation Plan: Auto-Stage Modified Files and State Restoration

**Branch**: `002-auto-stage-restore` | **Date**: 2025-01-27 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/002-auto-stage-restore/spec.md`

## Summary

This feature enhances the gitcomm CLI to automatically stage modified files on launch and restore staging state if the CLI exits without committing. The implementation will:

1. **Auto-stage modified files** at CLI startup (before any prompts)
2. **Optionally stage untracked files** when `-a` flag is used
3. **Capture pre-CLI staging state** for restoration purposes
4. **Restore staging state** on cancellation, errors, or interruption
5. **Handle signal interrupts** gracefully with state restoration

The technical approach leverages the existing `GitRepository` interface and extends it with staging state capture/restoration capabilities. Signal handling will be implemented using Go's `os/signal` package for interruption handling.

## Technical Context

**Language/Version**: Go 1.25.0+
**Primary Dependencies**:
- `github.com/go-git/go-git/v5` - Git operations (existing)
- `github.com/spf13/cobra` - CLI framework (existing)
- `github.com/rs/zerolog` - Structured logging (existing)
- `os/signal` - Signal handling for interruption
- `context.Context` - Cancellation and timeout handling

**Storage**: In-memory only (no persistent storage). Staging state snapshot stored in memory during CLI execution.

**Testing**:
- Standard Go testing framework (`testing` package)
- `github.com/onsi/ginkgo/v2` and `github.com/onsi/gomega` for BDD-style tests (existing)
- Integration tests for git operations
- Mock implementations for testing

**Target Platform**: Linux, macOS, Windows (CLI application)

**Project Type**: CLI tool (single binary)

**Performance Goals**:
- Auto-staging completes within 2 seconds (SC-001)
- State restoration completes within 1 second (SC-003)
- Staging failure detection within 1 second (SC-006)

**Constraints**:
- Must not leave repository in inconsistent state
- Must handle interruptions gracefully
- Must preserve files already staged before CLI launch
- Must work with existing gitcomm workflow

**Scale/Scope**:
- Single repository per CLI invocation
- Handles typical repository sizes (hundreds to thousands of files)
- No concurrent CLI instances on same repository (user responsibility)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Verify compliance with gitcomm Constitution principles:

- **Clean Architecture**: ✅ **COMPLIANT**
  - Extends existing `internal/repository` layer with staging state management
  - Uses Repository Pattern (existing `GitRepository` interface)
  - No new layers required, fits existing structure

- **Interface-Driven Development**: ✅ **COMPLIANT**
  - Extends existing `GitRepository` interface with new methods
  - Uses dependency injection via constructors (existing pattern)
  - No global state introduced

- **Test-First Development**: ✅ **COMPLIANT**
  - TDD approach: Write tests for staging state capture/restoration first
  - Unit tests for staging operations
  - Integration tests for signal handling and state restoration
  - Table-driven tests for edge cases

- **Idiomatic Go**: ✅ **COMPLIANT**
  - Follows Go naming conventions
  - Uses `context.Context` for cancellation
  - Error handling with wrapped errors
  - No panics in library code

- **Error Handling**: ✅ **COMPLIANT**
  - Custom error types for staging/restoration failures
  - Wrapped errors for traceability
  - Explicit error handling throughout

- **Context & Thread Safety**: ✅ **COMPLIANT**
  - Uses `context.Context` for cancellation and timeouts
  - Signal handling uses channels (thread-safe)
  - No shared mutable state

- **Technical Constraints**: ✅ **COMPLIANT**
  - No global state
  - Graceful shutdown via signal handlers
  - Resource cleanup (git operations are stateless)
  - Context propagation for cancellation

- **Operational Constraints**: ✅ **COMPLIANT**
  - Logging via `zerolog` (existing)
  - No secrets exposed (no new secrets introduced)
  - Error messages don't expose internal details

**Violations**: None. All principles are satisfied.

## Project Structure

### Documentation (this feature)

```text
specs/002-auto-stage-restore/
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
│   ├── git_repository.go              # Interface (extend with new methods)
│   ├── git_repository_impl.go         # Implementation (extend with staging state management)
│   └── staging_state.go                # NEW: Staging state capture/restoration logic
├── service/
│   └── commit_service.go               # Extend to call auto-staging and handle restoration
├── model/
│   └── staging_state.go                # NEW: StagingState domain model
└── utils/
    └── errors.go                        # Extend with staging/restoration error types

cmd/gitcomm/
└── main.go                              # Extend to handle signal interrupts

test/
├── integration/
│   └── staging_restore_test.go         # NEW: Integration tests for staging/restoration
└── mocks/
    └── git_repository_mock.go           # Extend mock for new methods
```

**Structure Decision**: Extends existing gitcomm structure. New code fits into existing layers:
- `internal/repository/staging_state.go` - Staging state management logic
- `internal/model/staging_state.go` - Domain model for staging state
- Extends existing `GitRepository` interface and implementation
- Signal handling in `cmd/gitcomm/main.go` (application layer)

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

No violations - all principles satisfied.
