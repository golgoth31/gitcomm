# Feature Specification: Fix CLI Hang on Ctrl+C During State Restoration

**Feature Branch**: `005-fix-ctrl-c-hang`
**Created**: 2025-01-27
**Status**: Draft
**Input**: User description: "fix the cli that stay stucked with message "restauring state" when receivinig ctrl+C command"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Graceful Interruption with Timely Exit (Priority: P1)

A developer presses Ctrl+C to interrupt the CLI and expects the CLI to restore staging state and exit promptly without hanging, even if restoration operations take time or encounter issues.

**Why this priority**: This is a critical bug fix that affects user experience. When users interrupt the CLI, they expect immediate feedback and exit, not a hanging process that requires force-killing.

**Independent Test**: Can be fully tested by running the CLI, pressing Ctrl+C during any phase of execution, and verifying that the CLI exits within a reasonable time (under 5 seconds) with appropriate exit code.

**Acceptance Scenarios**:

1. **Given** a developer runs the CLI which auto-stages files, **When** they press Ctrl+C during the commit workflow, **Then** the CLI displays "Interrupted. Restoring staging state..." and exits within 5 seconds with exit code 130
2. **Given** a developer presses Ctrl+C during state restoration, **When** the restoration operation is slow or blocking, **Then** the CLI applies a timeout (maximum 3 seconds) and exits even if restoration is incomplete, displaying appropriate warning
3. **Given** a developer presses Ctrl+C, **When** the restoration operation completes successfully, **Then** the CLI exits immediately after restoration completes (no hanging)
4. **Given** a developer presses Ctrl+C, **When** the restoration operation fails, **Then** the CLI displays a warning message and exits within 3 seconds (does not hang waiting for retry)

---

### Edge Cases

- What happens when Ctrl+C is pressed multiple times rapidly? → System should handle gracefully, exit after first restoration attempt completes or times out
- How does system handle restoration timeout? → System exits immediately with warning message after timeout, does not attempt partial restoration, does not retry or hang
- What happens if git operations are completely blocked (repository locked)? → System detects timeout, displays warning, exits immediately
- How does system handle restoration during interactive prompts (TUI)? → System cancels TUI immediately, performs restoration with timeout, exits

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST exit within 5 seconds of receiving Ctrl+C signal, regardless of restoration operation status
- **FR-002**: System MUST apply a timeout (maximum 3 seconds) to state restoration operations when triggered by Ctrl+C
- **FR-003**: System MUST not block on git operations during restoration - if operations exceed timeout, system MUST exit immediately with warning (no partial restoration attempts)
- **FR-004**: System MUST ensure restoration operations use a context with timeout that respects cancellation
- **FR-005**: System MUST display "Interrupted. Restoring staging state..." message immediately upon Ctrl+C
- **FR-006**: System MUST exit with code 130 (SIGINT) after restoration completes or times out
- **FR-007**: System MUST display a warning message if restoration times out or fails, but still exit promptly
- **FR-008**: System MUST ensure main process waits for restoration to complete or timeout before exiting (no race conditions)

### Key Entities

- **Restoration Context**: A context with timeout (3 seconds) used for all restoration operations when triggered by Ctrl+C. Must respect cancellation and timeout to prevent hanging.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: CLI exits within 5 seconds of Ctrl+C signal in 100% of test cases
- **SC-002**: Restoration operations complete or timeout within 3 seconds in 95% of normal cases
- **SC-003**: No test case results in CLI hanging indefinitely (all tests complete within 10 seconds maximum)
- **SC-004**: Users can successfully interrupt and exit CLI without requiring force-kill (kill -9) in 100% of cases

## Assumptions

- Git operations (unstage, capture state) typically complete in under 1 second in normal conditions
- Maximum acceptable restoration timeout is 3 seconds before forcing exit
- Users expect immediate exit after Ctrl+C, even if restoration is incomplete
- Warning messages are acceptable if restoration times out, as long as CLI exits promptly

## Dependencies

- Existing state restoration functionality (from feature 002-auto-stage-restore)
- Signal handling infrastructure (from feature 002-auto-stage-restore)
- Git repository operations (go-git library)

## Clarifications

### Session 2025-01-27

- Q: If restoration times out partway through (some files unstaged, others not), should the system attempt partial restoration or exit immediately? → A: Exit immediately without attempting to restore remaining files (prioritize fast exit per FR-001 and FR-003)

## Constraints

- Must maintain backward compatibility with existing restoration behavior (when not interrupted)
- Must not break existing signal handling for other scenarios
- Must work correctly with interactive TUI prompts (bubbletea) that may be active during interruption
