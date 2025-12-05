# Feature Specification: Auto-Stage Modified Files and State Restoration

**Feature Branch**: `002-auto-stage-restore`
**Created**: 2025-01-27
**Status**: Draft
**Input**: User description: "when launching the cli, all modified files should be staged. if the \"a\" flag is used, also add unmanaged files. the staged diff is to be analised by AI. if the cli is ended in any manner other then commiting with message, the stage state MUST be reverted to the state before the cli startup"

## Clarifications

### Session 2025-01-27

- Q: If some files stage successfully but others fail (partial staging), what should happen? → A: Abort and restore state (unstage all files, exit with error message, require user to fix issues)
- Q: If there are no modified files to stage (empty repository state), what should the CLI do? → A: Proceed with existing workflow (auto-staging is no-op, continue to AI/manual input, existing empty commit handling applies)
- Q: If the CLI is interrupted (e.g., Ctrl+C) while auto-staging is in progress, what should happen? → A: Attempt restoration on interruption (register signal handlers, catch interruption, restore staging state if staging was in progress, then exit)
- Q: If state restoration fails (e.g., git error, repository locked), what should the CLI do? → A: Exit with clear error and instructions (log detailed error, display warning message explaining state mismatch, provide manual recovery instructions, exit with non-zero code)
- Q: If repository state changes externally while CLI is running, how should restoration handle this? → A: Restore to captured pre-CLI state (ignore external changes, restore to the exact state captured at CLI launch, log warning if current state differs from expected)

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Auto-Stage Modified Files on Launch (Priority: P1)

A developer runs the CLI and expects all modified files to be automatically staged at startup, enabling AI analysis of the complete diff without manual staging steps.

**Why this priority**: This is the core workflow enhancement that improves user experience by eliminating the need to manually stage files before using the CLI. It ensures AI analysis has access to all changes, making the feature immediately valuable.

**Independent Test**: Can be fully tested by running the CLI in a git repository with modified files, verifying that all modified files are automatically staged, and confirming the staged diff is available for AI analysis. The test verifies staging occurs before any prompts are shown.

**Acceptance Scenarios**:

1. **Given** a developer has modified files in a git repository, **When** they run the CLI without the `-a` flag, **Then** all modified files are automatically staged before the CLI proceeds
2. **Given** a developer has modified files, **When** the CLI stages them automatically, **Then** the staged diff is available for AI analysis
3. **Given** a developer has both modified and untracked files, **When** they run the CLI without the `-a` flag, **Then** only modified files are staged (untracked files remain untracked)
4. **Given** a developer runs the CLI, **When** files are auto-staged, **Then** the staging happens before any user prompts are displayed

---

### User Story 2 - Auto-Stage Unmanaged Files with -a Flag (Priority: P2)

A developer runs the CLI with the `-a` flag and expects both modified files and unmanaged (untracked) files to be automatically staged, providing complete repository state for AI analysis.

**Why this priority**: This extends the auto-staging behavior to include untracked files when explicitly requested via the `-a` flag. It provides flexibility for users who want to include new files in their commit analysis.

**Independent Test**: Can be fully tested by running the CLI with the `-a` flag in a repository with both modified and untracked files, verifying that all files (modified and untracked) are staged, and confirming the complete diff is available for AI analysis.

**Acceptance Scenarios**:

1. **Given** a developer has both modified and untracked files, **When** they run the CLI with the `-a` flag, **Then** both modified and untracked files are automatically staged
2. **Given** a developer runs the CLI with `-a` flag, **When** files are auto-staged, **Then** the complete repository state (modified + untracked) is available for AI analysis
3. **Given** a developer runs the CLI with `-a` flag, **When** staging completes, **Then** the CLI proceeds with the normal workflow using the staged diff

---

### User Story 3 - Restore Staging State on Cancellation (Priority: P1)

A developer cancels or exits the CLI without committing, and expects the staging state to be restored to exactly what it was before the CLI was launched, ensuring no unintended changes persist.

**Why this priority**: This is critical for user trust and safety. Users must be confident that canceling the CLI won't leave their repository in an unexpected state. This prevents accidental staging of files that the user didn't intend to commit.

**Independent Test**: Can be fully tested by running the CLI, allowing it to auto-stage files, then canceling/exiting without committing, and verifying that the staging state matches the pre-CLI state. The test verifies restoration works for all exit scenarios (Ctrl+C, rejection of commit message, etc.).

**Acceptance Scenarios**:

1. **Given** a developer runs the CLI which auto-stages files, **When** they cancel the CLI (Ctrl+C or reject commit), **Then** the staging state is restored to the state before CLI launch
2. **Given** a developer runs the CLI with `-a` flag which stages untracked files, **When** they exit without committing, **Then** all auto-staged files (modified and untracked) are unstaged, restoring original state
3. **Given** a developer runs the CLI which auto-stages files, **When** they successfully commit, **Then** the staging state is NOT restored (commit is created with staged files)
4. **Given** a developer runs the CLI, **When** an error occurs that prevents commit creation, **Then** the staging state is restored to the pre-CLI state
5. **Given** a developer had some files already staged before running CLI, **When** CLI auto-stages additional files and then is canceled, **Then** only the files auto-staged by CLI are unstaged, preserving the originally staged files

---

### Edge Cases

- What happens when there are no modified files to stage? → Auto-staging is a no-op, CLI proceeds with existing workflow (existing empty commit handling applies)
- What happens when there are no untracked files but `-a` flag is used?
- What happens when staging fails (e.g., file permissions, repository locked)?
- What happens when state restoration fails (e.g., git operation error)? → System logs detailed error, displays warning message explaining state mismatch, provides manual recovery instructions, exits with non-zero code
- What happens when the CLI is interrupted during the staging process? → System registers signal handlers, catches interruption, restores staging state if staging was in progress, then exits
- What happens when the CLI is interrupted during the restoration process? → System completes restoration attempt (best-effort), then exits
- How does the system handle partial staging (some files staged successfully, others fail)? → System aborts staging, restores all files to pre-CLI state, exits with error message
- What happens if the repository state changes externally while CLI is running (another process modifies files)? → System restores to captured pre-CLI state (ignores external changes, restores to exact state captured at CLI launch, logs warning if current state differs from expected)
- What happens when the user has conflicts that prevent staging?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST automatically stage all modified files when CLI launches (before any prompts or AI analysis). If no modified files exist, auto-staging is a no-op and CLI proceeds with existing workflow
- **FR-002**: System MUST stage untracked (unmanaged) files when the `-a` flag is used
- **FR-003**: System MUST capture the staging state before CLI startup for restoration purposes (snapshot is used as restoration target, ignoring external changes during CLI execution)
- **FR-004**: System MUST use the staged diff (after auto-staging) for AI analysis
- **FR-005**: System MUST restore staging state to pre-CLI state if CLI exits without committing
- **FR-006**: System MUST NOT restore staging state if commit is successfully created
- **FR-007**: System MUST restore staging state on user cancellation (Ctrl+C, rejecting commit message, etc.)
- **FR-008**: System MUST restore staging state on errors that prevent commit creation
- **FR-009**: System MUST handle staging failures gracefully (if any file fails to stage, abort staging operation, restore all staged files to pre-CLI state, exit with error message, require user to fix issues before retrying)
- **FR-010**: System MUST handle restoration failures gracefully (log detailed error, display warning message explaining state mismatch, provide manual recovery instructions, exit with non-zero code)
- **FR-011**: System MUST preserve files that were already staged before CLI launch (only restore files staged by CLI)
- **FR-012**: System MUST restore state even if CLI is interrupted during execution (register signal handlers, catch interruptions, restore staging state if staging was in progress, then exit)

### Key Entities *(include if feature involves data)*

- **Pre-CLI Staging State**: Represents the snapshot of staging state (which files were staged) before CLI launch. Used to restore state if CLI exits without committing.

- **Auto-Staging Result**: Represents the result of automatic staging operation, including which files were successfully staged and which failed (if any).

- **Staging Restoration Plan**: Represents the plan for restoring staging state, including which files to unstage to return to pre-CLI state.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of modified files are automatically staged within 2 seconds of CLI launch
- **SC-002**: Staging state is restored correctly in 100% of cancellation scenarios
- **SC-003**: Users can cancel CLI without committing and have staging state restored in under 1 second
- **SC-004**: No staging state changes persist after CLI cancellation in 100% of cases
- **SC-005**: AI analysis uses staged diff in 100% of cases when auto-staging succeeds
- **SC-006**: Staging failures are detected and reported to user within 1 second
