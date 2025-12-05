# Implementation Plan: Rewrite CLI Prompts with Huh Library

**Branch**: `013-huh-prompts` | **Date**: 2025-01-27 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/013-huh-prompts/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Migrate all CLI prompts from custom Bubble Tea implementations to the `huh` library (github.com/charmbracelet/huh). All prompts must render inline (no alt screen) and display a post-validation summary line with green checkmark after user validation. Related prompts should be combined into single `huh.Form` instances with multiple fields, showing summary lines progressively as each field is completed.

## Technical Context

**Language/Version**: Go 1.25.0
**Primary Dependencies**:
  - `github.com/charmbracelet/huh` (new dependency - needs research)
  - `github.com/charmbracelet/bubbletea v1.3.10` (existing - will be used indirectly by `huh`)
  - `github.com/charmbracelet/lipgloss v1.1.0` (existing - for styling)
**Storage**: N/A (CLI tool, no persistent storage)
**Testing**: Go `testing` package with table-driven tests, integration tests in `test/integration/`
**Target Platform**: Cross-platform CLI (Linux, macOS, Windows terminals)
**Project Type**: Single CLI application
**Performance Goals**: Prompt rendering should feel responsive (<100ms for field transitions, <500ms for form initialization)
**Constraints**:
  - Must maintain backward compatibility with existing prompt function signatures
  - Must render inline (no alt screen mode)
  - Must support all existing prompt types and validation logic
  - Must handle terminal width/height constraints gracefully
**Scale/Scope**:
  - Replace ~10 prompt functions in `internal/ui/prompts.go`
  - Migrate 4 custom Bubble Tea models (TextInputModel, MultilineInputModel, YesNoChoiceModel, SelectListModel)
  - Update ~16 call sites in `internal/service/commit_service.go`
  - Maintain all existing tests and add new tests for `huh` integration

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Verify compliance with gitcomm Constitution principles:

- **Clean Architecture**: ✅ YES - Changes are within `internal/ui/` layer. No repository pattern needed (no data access). Existing layer separation maintained.
- **Interface-Driven Development**: ✅ YES - Prompt functions already use function signatures (interfaces). `huh` library will be injected via function parameters or wrapped in adapter functions. No global state.
- **Test-First Development**: ✅ YES - TDD required. Tests will be written before implementation for each prompt migration. Table-driven tests for validation scenarios.
- **Idiomatic Go**: ✅ YES - Follow Go conventions. Function names remain PascalCase. Error handling via return values. No panics in library code.
- **Error Handling**: ✅ YES - All prompt functions return `(value, error)`. Errors wrapped with context. Custom error types preserved from existing implementation.
- **Context & Thread Safety**: ✅ YES - Context propagation maintained where applicable. No new concurrency patterns needed (single-threaded CLI interaction).
- **Technical Constraints**: ✅ YES - No global state. Prompt functions are stateless. Resource cleanup handled by `huh` library. Graceful cancellation via existing Ctrl+C handling.
- **Operational Constraints**: ✅ YES - Existing logging via `zerolog` maintained. No secrets management needed (CLI tool).

**Violations**: None identified. All principles maintained.

## Project Structure

### Documentation (this feature)

```text
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
internal/ui/
├── prompts.go              # Main prompt functions (to be migrated to use `huh`)
├── prompts_test.go         # Tests for prompt functions
├── text_input.go           # TextInputModel (to be removed - replaced by `huh.NewInput()`)
├── text_input_test.go      # Tests (to be removed)
├── multiline_input.go      # MultilineInputModel (to be removed - replaced by `huh.NewText()`)
├── multiline_input_test.go  # Tests (to be removed)
├── yes_no_choice.go        # YesNoChoiceModel (to be removed - replaced by `huh.NewConfirm()`)
├── yes_no_choice_test.go   # Tests (to be removed)
├── select_list.go          # SelectListModel (to be removed - replaced by `huh.NewSelect()`)
├── select_list_test.go     # Tests (to be removed)
├── display.go              # Display utilities (kept - may need updates for post-validation format)
├── display_test.go         # Tests (kept)
├── prompt_state.go         # PromptState enum (may be kept for compatibility or removed)
└── prompt_state_test.go    # Tests (may be kept or removed)

internal/service/
└── commit_service.go       # Calls prompt functions (no changes needed - backward compatible)

test/integration/
└── [existing integration tests - may need updates for `huh` behavior]
```

**Structure Decision**: Single project structure. Changes are isolated to `internal/ui/` package. Existing structure maintained. Custom Bubble Tea models will be removed and replaced with `huh` library usage.

## Phase 0: Research Complete ✅

**Status**: Complete
**Output**: [research.md](./research.md)

All technical unknowns resolved:
- ✅ `huh` library capabilities and API patterns
- ✅ Inline rendering configuration
- ✅ Post-validation display implementation approach
- ✅ Progressive field display strategy
- ✅ Validation error display mechanism
- ✅ Backward compatibility strategy
- ✅ Combined forms implementation approach

## Phase 1: Design Complete ✅

**Status**: Complete
**Outputs**:
- [data-model.md](./data-model.md) - Function signatures and data structures
- [contracts/prompt-functions-contract.md](./contracts/prompt-functions-contract.md) - API contracts
- [quickstart.md](./quickstart.md) - Implementation guide

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

No violations identified. All principles maintained.
