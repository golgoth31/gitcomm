# Implementation Plan: Rewrite All CLI Prompts

**Branch**: `001-rewrite-cli-prompts` | **Date**: 2025-01-27 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-rewrite-cli-prompts/spec.md`

## Summary

Rewrite all CLI prompts to use bubbletea consistently with unified visual design. All prompts must start with a blue '?' indicator, show a green checkmark on completion, red 'X' on cancellation, and yellow/orange warning on validation errors. Prompts must render inline without alt screen mode. Select lists display selected value after title, multiline inputs display result under title with wrapping. All existing prompt functionality (validation, defaults, pre-selection) must be preserved while maintaining backward-compatible function signatures.

## Technical Context

**Language/Version**: Go 1.25.0
**Primary Dependencies**:
- `github.com/charmbracelet/bubbletea v1.3.10` (TUI framework)
- `github.com/charmbracelet/lipgloss v1.1.0` (styling)
**Storage**: N/A (in-memory prompt state)
**Testing**: Go standard testing framework (`testing` package), table-driven tests
**Target Platform**: Linux/Unix/macOS terminals (CLI application)
**Project Type**: Single CLI application
**Performance Goals**: Prompt rendering should be responsive (<100ms for state updates)
**Constraints**:
- Must not use alt screen mode (inline rendering only)
- Must maintain backward-compatible function signatures
- Must preserve all existing validation and business logic
**Scale/Scope**:
- 18 prompt functions to migrate
- 4 prompt types: text input, select lists, multiline input, yes/no choices
- All prompts must support visual indicators and state transitions

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Verify compliance with gitcomm Constitution principles:

- **Clean Architecture**: ✅ Structure follows layer separation (cmd/, internal/, pkg/). UI components in `internal/ui/` are properly separated from business logic. No repository pattern needed (no data persistence).

- **Interface-Driven Development**: ✅ Prompt functions already use interfaces implicitly (return types). Bubbletea models implement `tea.Model` interface. Dependency injection via constructors (New* functions) is used.

- **Test-First Development**: ✅ TDD approach required. Tests must be written before implementation for each prompt type migration. Table-driven tests for validation scenarios.

- **Idiomatic Go**: ✅ Follows Go conventions. Naming: PascalCase for exported, camelCase for unexported. Functions are small and focused. Uses standard Go testing patterns.

- **Error Handling**: ✅ Explicit error handling required. Custom error types for prompt cancellation, validation errors. Errors wrapped with context using `fmt.Errorf`.

- **Context & Thread Safety**: ✅ No concurrency required for prompts (single-threaded TUI). Context propagation not needed for prompt interactions.

- **Technical Constraints**: ✅ No global state (all state in models). Graceful shutdown handled by bubbletea. Resource cleanup via defer in prompt functions.

- **Operational Constraints**: ✅ Logging via zerolog for debug messages (prompt state transitions). No secrets involved in prompts.

**Violations**: None. All principles are satisfied.

## Project Structure

### Documentation (this feature)

```text
specs/001-rewrite-cli-prompts/
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
├── prompts.go           # Main prompt functions (to be refactored)
├── prompts_test.go     # Tests for prompt functions
├── select_list.go      # Select list model (to be updated - remove alt screen)
├── select_list_test.go # Tests for select list
├── multiline_input.go  # Multiline input model (to be updated - remove alt screen)
├── multiline_input_test.go # Tests for multiline input
├── text_input.go       # NEW: Text input model for bubbletea
├── text_input_test.go  # NEW: Tests for text input
├── yes_no_choice.go    # NEW: Yes/no choice model for bubbletea
├── yes_no_choice_test.go # NEW: Tests for yes/no choice
└── display.go          # Display utilities (may need updates for visual indicators)
```

**Structure Decision**: Single project structure. All UI components in `internal/ui/` package. New models for text input and yes/no choices. Existing select list and multiline input models updated to remove alt screen and add visual indicators.

## Complexity Tracking

> **No violations - all Constitution principles satisfied**
