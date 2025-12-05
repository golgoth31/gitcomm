# Implementation Plan: Improved Commit Message UI

**Branch**: `004-improve-commit-ui` | **Date**: 2025-01-27 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/004-improve-commit-ui/spec.md`

## Summary

This feature improves the commit message UI by replacing the numbered text-based commit type selection with an interactive select list (with checkmarks and highlighting), and enabling proper multiline input for body and footer fields. The implementation will:

1. **Replace commit type selection** with an interactive TUI select list using bubbletea
2. **Enhance multiline input** for body and footer fields with double-Enter completion
3. **Add letter-based navigation** for quick commit type selection
4. **Improve visual feedback** with checkmarks, highlighting, and clear selection indicators

The technical approach extends the existing `bubbletea` TUI framework usage (already planned in 001-git-commit-cli) to implement interactive select lists and enhanced multiline input fields.

## Technical Context

**Language/Version**: Go 1.25.0+
**Primary Dependencies**:
- `github.com/charmbracelet/bubbletea` - TUI framework for interactive prompts (existing, needs extension)
- `github.com/charmbracelet/lipgloss` - Styling for TUI components (existing, needs extension)
- `bufio` - Standard library for multiline input handling (existing)

**Storage**: N/A (UI state is in-memory only during prompt session)

**Testing**:
- Standard Go testing framework (`testing` package)
- `github.com/onsi/ginkgo/v2` and `github.com/onsi/gomega` for BDD-style tests (existing)
- Unit tests for UI components
- Integration tests for interactive workflows

**Target Platform**: Linux, macOS, Windows (CLI application with TUI)

**Project Type**: CLI tool (single binary)

**Performance Goals**:
- Commit type selection responds to arrow key input within 50ms
- Multiline input handles keystrokes within 10ms
- No noticeable delay in UI updates

**Constraints**:
- Must maintain backward compatibility with existing commit workflow
- Must work in terminals with limited color support (graceful degradation)
- Must handle terminal resizing gracefully
- Must support escape key cancellation with state restoration

**Scale/Scope**:
- Single CLI binary
- 8 commit types in select list
- Multiline input for body and footer (no practical limit, but should handle typical commit message lengths efficiently)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Verify compliance with gitcomm Constitution principles:

- **Clean Architecture**: ✅ **COMPLIANT**
  - Extends existing `internal/ui` layer with new TUI components
  - No new layers required, fits existing structure
  - UI components are presentation-layer concern

- **Interface-Driven Development**: ✅ **COMPLIANT**
  - UI components can be abstracted via interfaces for testing
  - Dependencies injected via constructors
  - No global state introduced

- **Test-First Development**: ✅ **COMPLIANT**
  - TDD approach: Write tests for UI components first
  - Unit tests for select list and multiline input components
  - Integration tests for interactive workflows
  - Table-driven tests for edge cases

- **Idiomatic Go**: ✅ **COMPLIANT**
  - Follows Go naming conventions
  - Uses standard library where possible
  - No panics in library code

- **Error Handling**: ✅ **COMPLIANT**
  - UI component errors handled explicitly
  - Cancellation (Escape key) returns appropriate errors
  - Error types for UI-specific failures

- **Context & Thread Safety**: ✅ **COMPLIANT**
  - UI components are single-threaded (bubbletea model)
  - Context used for cancellation
  - No shared mutable state in UI components

- **Technical Constraints**: ✅ **COMPLIANT**
  - No global state (UI state is model-local)
  - Graceful shutdown via Escape key handling
  - Resource cleanup handled by bubbletea lifecycle

- **Operational Constraints**: ✅ **COMPLIANT**
  - Logging strategy defined (debug logging from 003-debug-logging)
  - No secrets in UI (user input only)

**Violations**: None. All principles are satisfied.

### Post-Design Constitution Check

After Phase 1 design completion, all principles remain satisfied:

- **Clean Architecture**: ✅ UI components remain in `internal/ui` layer, no architectural changes
- **Interface-Driven Development**: ✅ UI components can be abstracted via interfaces for testing
- **Test-First Development**: ✅ Test strategy defined for UI components and workflows
- **Idiomatic Go**: ✅ Uses standard Go patterns and bubbletea idioms
- **Error Handling**: ✅ Error types defined, cancellation handled explicitly
- **Context & Thread Safety**: ✅ UI components are single-threaded, context used for cancellation
- **Technical Constraints**: ✅ No global state, graceful shutdown via Escape key
- **Operational Constraints**: ✅ Logging strategy defined, no secrets in UI

## Project Structure

### Documentation (this feature)

```text
specs/004-improve-commit-ui/
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
├── ui/
│   ├── prompts.go              # Existing - needs refactoring for TUI
│   ├── select_list.go          # New - interactive select list component
│   ├── multiline_input.go      # New - multiline input component
│   └── (existing files)
└── (other directories unchanged)

test/
├── integration/
│   └── ui_interactive_test.go  # New - integration tests for interactive UI
└── (existing test directories)
```

**Structure Decision**: Minimal changes to existing structure. Primary additions:
- `internal/ui/select_list.go` - Interactive select list component using bubbletea
- `internal/ui/multiline_input.go` - Multiline input component using bubbletea
- Refactor `internal/ui/prompts.go` to use new TUI components
- Integration tests for interactive workflows

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

No violations - all principles satisfied.
