# Feature Specification: Explicit Commit Error Messages

**Feature Branch**: `017-explicit-commit-errors`  
**Created**: 2026-02-10  
**Status**: Draft  
**Input**: User description: "j'ai l'erreur 'Error creating commit: failed to create commit: git commit failed (exit 1):' sans autre message explicite. il faut afficher plus d'informations en cas d'erreur"

## Clarifications

### Session 2026-02-10

- Q: When git error output is very long, what truncation strategy should we use? → A: First segment + size indicator: keep first X characters (e.g., 1500), add "… (N additional characters)" if truncated
- Q: When debug mode is enabled, how should the explicit error interact with logging? → A: Supplement: Always display the explicit error to the user. In debug mode, logs may duplicate or augment it; the explicit display remains the primary user-facing channel
- Q: Should this feature include improving error display for stage/unstage/status operations, or only commits? → A: Include in scope: Improve error display for stage, unstage, and status operations when they reach the user, with same level of detail
- Q: What format should the displayed error use when including git output? → A: Prefixed + git output: Use a structure like "Error creating commit: [brief]. Details: [git stderr content]" so the summary is separated from the actionable detail
- Q: When stderr is empty, what message should we show? → A: Generic: Show exit code plus a hint such as "No additional details from git. Check repository state or run the command manually."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - See Explicit Error Details When Commit Fails (Priority: P1)

When a user creates a commit and the operation fails (e.g., pre-commit hook failure, signing error, empty message, git config issue), the user sees a clear, explicit error message that explains what went wrong. The message includes the relevant output from the git command (typically stderr) so the user can diagnose and fix the issue without having to run git manually to reproduce the error.

**Why this priority**: Without explicit error details, users cannot troubleshoot commit failures and are forced to guess or reproduce the error manually.

**Independent Test**: Can be fully tested by triggering a commit failure (e.g., invalid hook, signing key issue) and verifying the displayed message contains actionable information (e.g., the actual git error output).

**Acceptance Scenarios**:

1. **Given** a repository with a pre-commit hook that fails, **When** the user attempts to create a commit, **Then** the error message displayed includes the hook failure output (stderr from git)
2. **Given** a repository with SSH signing configured but an invalid key, **When** the user attempts to create a signed commit, **Then** the error message includes the signing failure details from git
3. **Given** any git commit failure, **When** the error is displayed, **Then** the message contains more than just "git commit failed (exit N):" and includes the actual error output from git
4. **Given** a commit failure where git produces no stderr but returns non-zero, **When** the error is displayed, **Then** the message still indicates the exit code and suggests checking git configuration or running the command manually

---

### User Story 2 - Consistent Explicit Errors for All Git Operations (Priority: P2)

When any git operation fails (stage, unstage, status, etc.), the user sees explicit error details. This ensures a consistent experience: every git-related failure surfaces the underlying cause.

**Why this priority**: Commit failures are most critical, but staging or restoration failures should also be explicit for a cohesive user experience.

**Independent Test**: Can be tested by triggering failures in StageAllFiles, UnstageFiles, or GetRepositoryState and verifying explicit error details are shown.

**Acceptance Scenarios**:

1. **Given** a staging operation that fails (e.g., permission denied), **When** the error is displayed, **Then** the message uses a prefixed structure (brief summary + "Details:" + git stderr or categorized description)
2. **Given** any git operation failure, **When** the error is surfaced to the user, **Then** the message does not end with an empty or truncated ":" without further details

---

### Edge Cases

- **Empty stderr**: When git fails but produces no stderr, the message MUST include exit code and a generic hint (e.g., "No additional details from git. Check repository state or run the command manually.")
- **Very long stderr**: Git may output many lines (e.g., hook output). Truncate by keeping the first 1500 characters and appending "… (N additional characters)" to indicate how much was omitted.
- **Sensitive data**: Git stderr should not contain secrets by default; if a future feature adds credential handling, ensure errors do not leak tokens. For this feature, assume standard git outputs (hooks, signing, config) which are typically non-sensitive.
- **Non-commit failures**: Errors from staging, unstaging, or status should follow the same explicit-message pattern when surfaced to the user.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: When a commit creation fails, the system MUST display the git command stderr output (or equivalent error details) to the user, using a structure that separates a brief summary from the detail (e.g., "Error creating commit: [brief]. Details: [git stderr content]")
- **FR-002**: When a commit creation fails and git produces no stderr, the system MUST display the exit code and a generic hint (e.g., "No additional details from git. Check repository state or run the command manually.")
- **FR-003**: Error messages for commit failures MUST NOT end with a truncated or empty suffix (e.g., "git commit failed (exit 1):" with nothing after the colon)
- **FR-004**: When displaying git operation failures to the user (commit, stage, unstage, status, or any other git operation), the system MUST preserve and surface the underlying error details through all error wrapping layers
- **FR-005**: For very long error output, the system MUST truncate by keeping the first 1500 characters and appending "… (N additional characters)" when truncated, preserving the most actionable information (error messages are typically at the start)

### Assumptions

- Explicit error display is the primary user-facing channel for commit failures; debug mode supplements with logs but does not replace the explicit error message
- Git CLI outputs failures to stderr; the current implementation captures stderr via `exec.Cmd`
- Users expect plain-language error display in the CLI (no need for structured JSON or machine-parseable format for end users)
- Commit failures are the primary focus; stage, unstage, and status errors that reach the user are in scope and receive the same explicit-error treatment
- The existing error categorization (ErrGitCommandFailed, ErrGitSigningFailed, etc.) remains; the improvement is ensuring these details are not lost when errors are wrapped and displayed

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users encountering a commit failure see at least one line of explicit error content (stderr or equivalent) in 100% of cases where git produces such output
- **SC-002**: No commit failure results in a message that ends with ":" followed by no further text
- **SC-003**: Users can diagnose at least 80% of common commit failures (hooks, signing, empty message, config) from the displayed error without running git manually
- **SC-004**: Error messages for git failures remain human-readable; when truncated, the visible portion does not exceed 1500 characters plus the "… (N additional characters)" suffix (to avoid terminal overflow)
