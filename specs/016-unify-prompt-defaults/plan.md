# Implementation Plan: Unify Prompt Functions to Use Default Variants

**Branch**: `016-unify-prompt-defaults` | **Date**: 2025-01-27 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/016-unify-prompt-defaults/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Refactor the `promptCommitMessage` function in `commit_service.go` to always use prompt functions with default value support (`PromptScopeWithDefault`, `PromptSubjectWithDefault`, etc.), eliminating conditional logic that selects between regular and "WithDefault" variants. Remove the non-default prompt functions (`PromptScope`, `PromptSubject`, `PromptBody`, `PromptFooter`, `PromptCommitType`) from `prompts.go` after confirming they are not used elsewhere. This simplifies code maintenance and ensures consistent behavior across all prompt interactions.

## Technical Context

**Language/Version**: Go 1.25.0
**Primary Dependencies**: github.com/charmbracelet/huh v0.8.0 (existing)
**Storage**: N/A (no storage changes)
**Testing**: Go standard testing framework (`testing` package)
**Target Platform**: CLI tool (Linux/macOS/Windows)
**Project Type**: Single CLI application
**Performance Goals**: N/A (refactoring only, no performance impact expected)
**Constraints**: Must maintain existing functionality - prompts must behave identically when empty strings are passed as defaults
**Scale/Scope**: Single function refactoring (`promptCommitMessage` in `internal/service/commit_service.go`) and removal of 5 functions from `internal/ui/prompts.go`

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Verify compliance with gitcomm Constitution principles:

- **Clean Architecture**: ✅ No changes to architecture layers. Refactoring within existing service layer (`internal/service/`) and UI layer (`internal/ui/`). Repository Pattern not affected.
- **Interface-Driven Development**: ✅ No interface changes. Existing interfaces remain unchanged. Dependency injection via constructors maintained.
- **Test-First Development**: ✅ Tests must be updated/added before implementation. TDD approach: Update tests → Verify failures → Refactor code → Verify passes.
- **Idiomatic Go**: ✅ Refactoring follows Go conventions. Naming conventions maintained. Code simplification improves readability.
- **Error Handling**: ✅ Error handling patterns unchanged. Existing error wrapping and propagation maintained.
- **Context & Thread Safety**: ✅ No context or concurrency changes. Existing patterns preserved.
- **Technical Constraints**: ✅ No global state introduced. No resource management changes. Graceful error handling maintained.
- **Operational Constraints**: ✅ No logging changes required. No secrets management changes.

**Violations**: None. This refactoring aligns with all constitution principles.

## Project Structure

### Documentation (this feature)

```text
specs/016-unify-prompt-defaults/
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
├── service/
│   └── commit_service.go        # Refactor promptCommitMessage function
└── ui/
    └── prompts.go                # Remove non-default prompt functions
```

**Structure Decision**: Single project structure. Refactoring affects only two files:
- `internal/service/commit_service.go`: Update `promptCommitMessage` function
- `internal/ui/prompts.go`: Remove 5 non-default prompt functions

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

No violations. This is a straightforward refactoring that simplifies code without introducing complexity.

## Phase Completion Status

### Phase 0: Outline & Research ✅

**Status**: Complete

**Deliverables**:
- ✅ `research.md`: Verified behavior of "WithDefault" functions with empty string defaults
- ✅ Confirmed non-default functions are only used in `commit_service.go`
- ✅ Determined correct default value pattern for all prompt types

**Key Findings**:
- All "WithDefault" functions correctly handle empty string defaults
- Empty strings maintain existing validation behavior
- No breaking changes to functionality

### Phase 1: Design & Contracts ✅

**Status**: Complete

**Deliverables**:
- ✅ `data-model.md`: Documented that no data model changes are required
- ✅ `contracts/prompt-functions-contract.md`: Documented function contracts and refactoring patterns
- ✅ `quickstart.md`: Developer guide for understanding and implementing the refactoring
- ✅ Agent context updated: Cursor IDE context file updated with Go 1.25.0 and huh v0.8.0

**Design Decisions**:
- Unified pattern: Always use "WithDefault" variants with empty string defaults when no pre-filled data
- Function removal: Remove 5 non-default prompt functions after confirming no other usage
- Error handling: Maintain existing error wrapping and propagation patterns

### Phase 2: Task Breakdown

**Status**: Pending `/speckit.tasks` command

**Next Steps**: Run `/speckit.tasks` to generate task breakdown for implementation.
