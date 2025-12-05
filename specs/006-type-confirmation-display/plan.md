# Implementation Plan: Display Commit Type Selection Confirmation

**Branch**: `006-type-confirmation-display` | **Date**: 2025-01-27 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/006-type-confirmation-display/spec.md`

## Summary

This feature adds a confirmation line display after a commit type is selected from the interactive list. The confirmation line uses the format "✔ Choose a type(<scope>): <chosen type>" and appears on a new line after the alt-screen clears, before the next prompt (scope) appears. This provides immediate visual feedback to users that their selection was successful.

The technical approach is straightforward: modify the `PromptCommitType` function in `internal/ui/prompts.go` to display the confirmation line using `fmt.Printf` after the bubbletea program exits and before returning the selected type. This is a display-only enhancement that does not change the selection logic or data flow.

## Technical Context

**Language/Version**: Go 1.25.0+
**Primary Dependencies**:
- `github.com/charmbracelet/bubbletea` - TUI framework (existing, used for interactive select list)
- `fmt` - Standard library for formatted output (standard library)
- Standard Go libraries only

**Storage**: N/A (display-only feature, no data persistence)

**Testing**:
- Standard Go testing framework (`testing` package)
- Integration tests for UI behavior verification
- No external testing dependencies

**Target Platform**: Linux, macOS, Windows (CLI application)

**Project Type**: CLI tool (single binary)

**Performance Goals**:
- Confirmation line appears within 100ms of selection confirmation (SC-002)
- No measurable performance impact on commit workflow

**Constraints**:
- Must maintain backward compatibility with existing commit type selection behavior
- Must not break existing signal handling or interruption flows
- Must work correctly with bubbletea alt-screen mode
- Checkmark symbol (✔) must be supported by terminal encoding

**Scale/Scope**:
- Single repository per CLI invocation
- Handles all predefined commit types (feat, fix, docs, style, refactor, test, chore, version)
- No concurrent selection scenarios (single user, single selection)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Verify compliance with gitcomm Constitution principles:

- **Clean Architecture**: ✅ Follows existing layer separation - modification is in `internal/ui/` layer, no new layers needed
- **Interface-Driven Development**: ✅ No new interfaces required - uses existing `PromptCommitType` function signature
- **Test-First Development**: ✅ Tests will be written before implementation (TDD approach for display verification)
- **Idiomatic Go**: ✅ Uses standard `fmt.Printf` for output, follows Go naming conventions
- **Error Handling**: ✅ Error handling already exists in `PromptCommitType` - no changes needed
- **Context & Thread Safety**: ✅ No concurrency involved - single-threaded UI interaction
- **Technical Constraints**: ✅ No global state, no new resources, simple display logic
- **Operational Constraints**: ✅ No logging needed for this simple display feature

**Violations**: None - this feature fully complies with all constitution principles.

### Post-Design Constitution Check

*Re-evaluated after Phase 1 design artifacts created.*

After completing research, data model, and contracts:

- **Clean Architecture**: ✅ Confirmed - modification isolated to UI layer, no architectural changes
- **Interface-Driven Development**: ✅ Confirmed - no new interfaces, uses existing function signature
- **Test-First Development**: ✅ Confirmed - test strategy defined in contracts, TDD approach maintained
- **Idiomatic Go**: ✅ Confirmed - uses standard `fmt.Printf`, follows Go conventions
- **Error Handling**: ✅ Confirmed - existing error handling preserved, no new error cases
- **Context & Thread Safety**: ✅ Confirmed - no concurrency, single-threaded UI interaction
- **Technical Constraints**: ✅ Confirmed - no global state, no new resources, simple display
- **Operational Constraints**: ✅ Confirmed - no logging needed, no secrets, no configuration

**Post-Design Violations**: None - design artifacts confirm full compliance with all principles.

## Project Structure

### Documentation (this feature)

```text
specs/006-type-confirmation-display/
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
├── prompts.go           # Modify PromptCommitType to display confirmation line
└── prompts_test.go      # Add tests for confirmation line display

test/integration/
└── ui_confirmation_test.go  # Integration tests for confirmation display
```

**Structure Decision**: This is a simple modification to existing UI code. No new packages or modules needed. The change is isolated to the `PromptCommitType` function in `internal/ui/prompts.go` with corresponding test updates.

## Complexity Tracking

> **No violations - feature fully complies with constitution principles**
