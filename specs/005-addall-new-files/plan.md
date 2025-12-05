# Implementation Plan: Respect addAll Flag for New Files

**Branch**: `005-addall-new-files` | **Date**: 2025-01-27 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/005-addall-new-files/spec.md`

## Summary

The `addAll` flag currently controls which files are staged (via `StageModifiedFiles` vs `StageAllFilesIncludingUntracked`), but `GetRepositoryState` doesn't respect this flag when building the repository state. New/untracked files that are staged (either manually or via `-a` flag) are always included in the repository state, regardless of the `addAll` flag value. This feature will filter repository state results to exclude new files when `addAll` is false, ensuring the flag controls both staging and state retrieval behavior.

## Technical Context

**Language/Version**: Go 1.25.0
**Primary Dependencies**:
- `github.com/go-git/go-git/v5` (v5.16.4) - Git repository operations
- `github.com/golgoth31/gitcomm/internal/model` - Domain models
- `github.com/golgoth31/gitcomm/internal/repository` - Repository interface
- `github.com/golgoth31/gitcomm/internal/service` - Commit service

**Storage**: N/A (in-memory state processing)
**Testing**: Go `testing` package with table-driven tests
**Target Platform**: Linux/macOS/Windows (Go CLI application)
**Project Type**: Single CLI application
**Performance Goals**: No performance regression - repository state filtering should complete in <10ms for typical repositories (<100 files)
**Constraints**:
- Must maintain backward compatibility with existing `GetRepositoryState` interface
- Must not break existing tests
- Must handle edge cases (manually staged new files, binary files, etc.)
**Scale/Scope**:
- Single repository per execution
- Typical repositories: 10-1000 files
- Filtering logic must be efficient for large repositories

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Pre-Phase 0 (Initial Check)

Verify compliance with gitcomm Constitution principles:

- **Clean Architecture**: ✅ Changes are isolated to `internal/repository` and `internal/service` layers. Repository pattern is maintained. No changes to `cmd/` or `pkg/` layers.
- **Interface-Driven Development**: ✅ `GitRepository` interface may need extension to pass `addAll` flag, or filtering happens in service layer using existing interface. Dependency injection maintained.
- **Test-First Development**: ✅ TDD required - tests must be written before implementation. Table-driven tests for filtering logic with various file status combinations.
- **Idiomatic Go**: ✅ Follows Go naming conventions. Uses existing patterns from codebase. Error handling with wrapped errors.
- **Error Handling**: ✅ Explicit error handling maintained. No new error types needed - reuse existing error handling patterns.
- **Context & Thread Safety**: ✅ `context.Context` already used in `GetRepositoryState`. No new concurrency patterns needed.
- **Technical Constraints**: ✅ No global state. No new resources to manage. Filtering is pure function.
- **Operational Constraints**: ✅ Existing logging via `zerolog` maintained. No secrets involved.

**Violations**: None - design fully complies with constitution.

### Post-Phase 1 (Design Complete)

After Phase 1 design (data model, contracts, quickstart), re-verification:

- **Clean Architecture**: ✅ Confirmed - filtering logic in repository layer, flag passing in service layer. No layer violations.
- **Interface-Driven Development**: ✅ Confirmed - using context values to pass flag maintains interface compatibility. No interface signature changes required.
- **Test-First Development**: ✅ Confirmed - test strategy defined in contracts. Table-driven tests planned for all file status combinations.
- **Idiomatic Go**: ✅ Confirmed - context value pattern is idiomatic Go. No anti-patterns introduced.
- **Error Handling**: ✅ Confirmed - filtering is silent (no errors), maintains existing error handling patterns.
- **Context & Thread Safety**: ✅ Confirmed - using `context.Context` for request-scoped values is appropriate. No concurrency concerns.
- **Technical Constraints**: ✅ Confirmed - no global state, no new resources, pure filtering function.
- **Operational Constraints**: ✅ Confirmed - no logging changes needed, no secrets involved.

**Violations**: None - Phase 1 design fully complies with constitution. Ready for Phase 2 (task breakdown).

## Project Structure

### Documentation (this feature)

```text
specs/005-addall-new-files/
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
│   ├── git_repository.go           # Interface definition (may extend)
│   └── git_repository_impl.go      # Implementation (filtering logic)
├── service/
│   └── commit_service.go           # Service layer (passes addAll flag)
└── model/
    └── repository_state.go         # Data models (no changes)

internal/repository/
└── git_repository_impl_test.go     # Unit tests for filtering
```

**Structure Decision**: Changes are isolated to existing repository and service layers. No new packages or modules needed. Filtering logic can be implemented either:
1. In `GetRepositoryState` by adding optional parameter (breaks interface)
2. In service layer after calling `GetRepositoryState` (maintains interface)
3. In `GetRepositoryState` by checking worktree status to identify new files (maintains interface, preferred)

Option 3 is preferred as it maintains interface compatibility and keeps filtering logic in the repository layer where file status knowledge exists.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

No violations - all design decisions comply with constitution principles.
