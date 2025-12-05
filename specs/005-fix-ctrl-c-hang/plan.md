# Implementation Plan: Fix CLI Hang on Ctrl+C During State Restoration

**Branch**: `005-fix-ctrl-c-hang` | **Date**: 2025-01-27 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/005-fix-ctrl-c-hang/spec.md`

## Summary

This bug fix addresses a critical issue where the CLI hangs indefinitely when Ctrl+C is pressed during state restoration. The problem occurs because restoration operations use `context.Background()` which doesn't respect cancellation or timeout, causing the CLI to block on git operations.

The implementation will:
1. **Add timeout context** (3 seconds) for restoration operations when triggered by Ctrl+C
2. **Ensure main process waits** for restoration to complete or timeout before exiting
3. **Respect cancellation** in all git operations during restoration
4. **Handle timeout gracefully** with warning messages and immediate exit

The technical approach modifies the existing restoration logic in `CommitService` to use a timeout context and ensures proper synchronization between signal handling and restoration completion.

## Technical Context

**Language/Version**: Go 1.25.0+
**Primary Dependencies**:
- `github.com/go-git/go-git/v5` - Git operations (existing)
- `github.com/spf13/cobra` - CLI framework (existing)
- `github.com/rs/zerolog` - Structured logging (existing)
- `os/signal` - Signal handling (existing)
- `context.Context` - Cancellation and timeout handling (standard library)
- `time` - Timeout duration management (standard library)

**Storage**: N/A (in-memory state only)

**Testing**:
- Standard Go testing framework (`testing` package)
- `github.com/onsi/ginkgo/v2` and `github.com/onsi/gomega` for BDD-style tests (existing)
- Integration tests for signal handling and timeout scenarios
- Timeout simulation tests

**Target Platform**: Linux, macOS, Windows (CLI application)

**Project Type**: CLI tool (single binary)

**Performance Goals**:
- CLI exits within 5 seconds of Ctrl+C signal (SC-001)
- Restoration operations complete or timeout within 3 seconds (SC-002)
- No test case hangs indefinitely (all tests complete within 10 seconds) (SC-003)

**Constraints**:
- Must maintain backward compatibility with existing restoration behavior (when not interrupted)
- Must not break existing signal handling for other scenarios
- Must work correctly with interactive TUI prompts (bubbletea) that may be active during interruption
- Must ensure no race conditions between signal handler and main process

**Scale/Scope**:
- Single repository per CLI invocation
- Handles typical repository sizes (hundreds to thousands of files)
- No concurrent CLI instances on same repository (user responsibility)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Verify compliance with gitcomm Constitution principles:

- **Clean Architecture**: ✅ **COMPLIANT**
  - Modifies existing `internal/service` layer (CommitService)
  - No new layers required, fits existing structure
  - Uses existing Repository Pattern

- **Interface-Driven Development**: ✅ **COMPLIANT**
  - Uses existing `GitRepository` interface
  - No new interfaces required
  - Uses dependency injection via constructors (existing pattern)

- **Test-First Development**: ✅ **COMPLIANT**
  - TDD approach: Write tests for timeout behavior first
  - Unit tests for timeout context creation
  - Integration tests for signal handling with timeout
  - Table-driven tests for timeout scenarios

- **Idiomatic Go**: ✅ **COMPLIANT**
  - Follows Go naming conventions
  - Uses `context.Context` for cancellation and timeout
  - Uses `context.WithTimeout` for timeout management
  - Standard library patterns for signal handling

- **Error Handling**: ✅ **COMPLIANT**
  - Explicit error handling for timeout scenarios
  - Custom error types for timeout detection
  - Wrapped errors for traceability

- **Context & Thread Safety**: ✅ **COMPLIANT**
  - Uses `context.Context` for cancellation and timeout
  - Proper synchronization between signal handler and main process
  - No shared mutable state without synchronization

- **Technical Constraints**: ✅ **COMPLIANT**
  - No global state (uses context and dependency injection)
  - Graceful shutdown with timeout
  - Resource cleanup handled by context cancellation

- **Operational Constraints**: ✅ **COMPLIANT**
  - Logging strategy defined (debug logging from 003-debug-logging)
  - No secrets involved (git operations only)

**Violations**: None. All principles are satisfied.

### Post-Design Constitution Check

After Phase 1 design completion, all principles remain satisfied:

- **Clean Architecture**: ✅ No architectural changes, modifications fit existing structure
- **Interface-Driven Development**: ✅ No new interfaces, uses existing patterns
- **Test-First Development**: ✅ Test strategy defined for timeout scenarios
- **Idiomatic Go**: ✅ Uses standard Go patterns (`context.WithTimeout`, channels)
- **Error Handling**: ✅ Timeout error handling defined, uses `errors.Is` for detection
- **Context & Thread Safety**: ✅ Timeout context usage defined, channel-based synchronization
- **Technical Constraints**: ✅ No global state, graceful shutdown with timeout
- **Operational Constraints**: ✅ Logging strategy defined, no secrets involved

## Project Structure

### Documentation (this feature)

```text
specs/005-fix-ctrl-c-hang/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
cmd/gitcomm/
└── main.go              # Signal handling modifications

internal/service/
└── commit_service.go     # Restoration timeout logic

internal/repository/
└── git_repository_impl.go # Context propagation for git operations

test/integration/
└── signal_timeout_test.go # Integration tests for timeout scenarios
```

**Structure Decision**: Single project structure. Modifications are limited to existing files in `cmd/gitcomm/` and `internal/service/`. No new packages or layers required.

## Complexity Tracking

> **No violations - all principles satisfied**
