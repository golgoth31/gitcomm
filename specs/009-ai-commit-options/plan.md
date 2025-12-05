# Implementation Plan: AI Commit Message Acceptance Options

**Branch**: `009-ai-commit-options` | **Date**: 2025-01-27 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/009-ai-commit-options/spec.md`

## Summary

This feature enhances the AI commit message acceptance workflow by providing three distinct options when an AI-generated message is displayed: "accept and commit directly" (immediate commit), "accept and edit" (pre-filled editing), and "reject" (start over). The implementation modifies the `PromptAIMessageAcceptance` function to return an enum-like response instead of a boolean, updates the commit service workflow to handle the three paths, and extends the manual prompt system to support pre-filling fields from AI messages.

## Technical Context

**Language/Version**: Go 1.25.0
**Primary Dependencies**:
- `github.com/charmbracelet/bubbletea` v1.3.10 (interactive UI for commit type selection)
- `github.com/charmbracelet/lipgloss` v1.1.0 (styling for interactive UI)
- Existing internal packages: `internal/ui`, `internal/service`, `internal/model`
**Storage**: N/A (in-memory state only)
**Testing**: Go standard `testing` package with table-driven tests
**Target Platform**: Linux/macOS/Windows CLI
**Project Type**: Single CLI application
**Performance Goals**:
- "Accept and commit directly" completes in <5 seconds (SC-001)
- "Accept and edit" completes in <30 seconds including editing (SC-002)
**Constraints**:
- Must maintain backward compatibility with existing commit workflow
- Must preserve existing staging state restoration behavior
- Must not break existing AI message generation or validation
**Scale/Scope**: Single-user CLI tool, no concurrency requirements

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Verify compliance with gitcomm Constitution principles:

- **Clean Architecture**: ✅ Follows existing layer separation (cmd/, internal/, pkg/). Repository Pattern already used for git operations. Changes are isolated to `internal/ui` and `internal/service` layers.

- **Interface-Driven Development**: ✅ UI functions use `*bufio.Reader` interface. Service layer uses `GitRepository` interface. No new interfaces needed - extending existing ones.

- **Test-First Development**: ✅ TDD approach required. Tests must be written before implementation for:
  - New `PromptAIMessageAcceptanceOptions` function
  - Modified `promptCommitMessage` with pre-filling
  - Updated `generateWithAI` workflow handling three options

- **Idiomatic Go**: ✅ Follows Go conventions. Uses existing patterns from codebase (error handling, naming, structure).

- **Error Handling**: ✅ Explicit error handling required. Uses existing error wrapping patterns. Custom errors may be needed for invalid acceptance responses.

- **Context & Thread Safety**: ✅ Uses `context.Context` for cancellation (already in place). No new concurrency patterns needed - single-threaded CLI workflow.

- **Technical Constraints**: ✅ No global state. All dependencies injected. Graceful cancellation already handled via context. Resource cleanup via defer patterns.

- **Operational Constraints**: ✅ Uses existing `zerolog` logging via `utils.Logger`. No secrets involved. Debug logging for acceptance flow decisions.

**Violations**: None - all changes align with existing architecture and patterns.

## Project Structure

### Documentation (this feature)

```text
specs/009-ai-commit-options/
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
│   ├── prompts.go                    # MODIFY: Add PromptAIMessageAcceptanceOptions function
│   ├── prompts_test.go               # MODIFY: Add tests for new acceptance options
│   ├── select_list.go                # MODIFY: Add support for pre-selecting commit type
│   └── select_list_test.go           # MODIFY: Add tests for pre-selection
├── service/
│   ├── commit_service.go             # MODIFY: Update generateWithAI to handle three options
│   └── commit_service_test.go         # MODIFY: Add tests for new workflow paths
└── model/
    └── commit_message.go              # NO CHANGE (existing structure sufficient)

test/
└── integration/
    └── ai_acceptance_options_test.go  # NEW: Integration tests for three-option workflow
```

**Structure Decision**: Single project structure. Changes are isolated to UI and service layers. No new packages needed. Follows existing patterns from previous features (e.g., 004-improve-commit-ui, 006-type-confirmation-display).

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

No violations - all changes comply with constitution principles.
