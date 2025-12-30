# Feature Specification: Replace Go-Git Library with External Git Commands

**Feature Branch**: `016-replace-gogit`  
**Created**: 2026-02-10  
**Status**: Draft  
**Input**: User description: "replace gogit library usage by using exec external git command. first retreive git usage, create git cli equivalent commands"

## Clarifications

### Session 2026-02-10

- Q: Should the system log git command executions for debugging? → A: Log all git command executions with command, arguments, exit codes, and execution time
- Q: Should go-git dependency be completely removed or kept as fallback? → A: Completely remove go-git dependency (no fallback)
- Q: What minimum git version should be enforced? → A: Enforce minimum git 2.34.0 (required for SSH signing)
- Q: How should SSH key passphrase be handled? → A: Support passphrase entry via environment variable or config
- Q: Should git command failures be categorized into distinct error types? → A: Categorize into distinct error types (not found, permission denied, invalid repo, etc.)

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Document Current Git Operations (Priority: P1)

As a developer working on this migration, I need to understand all git operations currently performed by the go-git library so I can identify equivalent git CLI commands.

**Why this priority**: This is the foundation for the entire migration - without understanding current usage, we cannot create accurate CLI equivalents.

**Independent Test**: Can be fully tested by analyzing the codebase and documenting all go-git API calls, their parameters, and their purposes. Delivers a comprehensive inventory of git operations.

**Acceptance Scenarios**:

1. **Given** the codebase uses go-git library, **When** analyzing all git operations, **Then** all go-git API calls are documented with their purpose and context
2. **Given** documented git operations, **When** reviewing the documentation, **Then** each operation includes the file location, method name, and what git concept it represents

---

### User Story 2 - Map Git Operations to CLI Commands (Priority: P1)

As a developer implementing the migration, I need a mapping document that shows which git CLI commands replace each go-git operation.

**Why this priority**: This mapping is essential for implementation - developers need to know exactly which git commands to execute for each operation.

**Independent Test**: Can be fully tested by verifying each go-git operation has a corresponding git CLI command documented. Delivers a complete translation guide from library calls to CLI commands.

**Acceptance Scenarios**:

1. **Given** documented git operations, **When** mapping to CLI commands, **Then** each operation has an equivalent git CLI command with required flags and arguments
2. **Given** the mapping document, **When** reviewing it, **Then** all git CLI commands include examples of expected input/output formats
3. **Given** operations that require parsing git output, **When** documenting CLI commands, **Then** the expected output format and parsing requirements are specified

---

### User Story 3 - Replace Go-Git Implementation with CLI Execution (Priority: P2)

As a developer, I need the GitRepository implementation to use external git commands instead of the go-git library while maintaining the same interface and behavior. The go-git dependency must be completely removed from the project.

**Why this priority**: This is the core implementation work that delivers the actual migration. However, it depends on the previous two stories being complete.

**Independent Test**: Can be fully tested by running existing tests and verifying all git operations work identically to the go-git implementation. Delivers a working implementation that maintains backward compatibility.

**Acceptance Scenarios**:

1. **Given** the GitRepository interface, **When** replacing go-git calls with CLI execution, **Then** all interface methods continue to work without changes to callers
2. **Given** existing tests, **When** running them against the new implementation, **Then** all tests pass with identical behavior
3. **Given** git operations that require context cancellation, **When** executing git commands, **Then** context cancellation properly terminates git processes
4. **Given** git operations that produce errors, **When** executing git commands, **Then** errors are properly wrapped and returned in the same format as go-git errors

---

### Edge Cases

- What happens when git command is not available in PATH?
- How does system handle git commands that require user interaction (e.g., SSH key passphrase)? → Handled via environment variable or configuration (no interactive prompts)
- How does system handle git operations in repositories with no commits (empty repository)?
- How does system handle git operations when repository is in a detached HEAD state?
- How does system handle git operations when there are merge conflicts?
- How does system handle git operations when repository is corrupted or .git directory is missing?
- How does system handle git operations when filesystem permissions prevent git operations?
- How does system handle git commands that produce non-zero exit codes?
- How does system handle git commands that produce output to stderr vs stdout?
- How does system handle binary files in diff operations?
- How does system handle very large files or repositories?
- How does system handle git operations when working directory changes during execution?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST document all current go-git library usage, including API calls, their parameters, and their purposes
- **FR-002**: System MUST create a mapping document that shows equivalent git CLI commands for each go-git operation
- **FR-003**: System MUST execute git commands in the correct repository directory (using `-C` flag or `chdir`)
- **FR-004**: System MUST handle context cancellation by terminating git processes when context is cancelled
- **FR-005**: System MUST parse git command output to extract required information (status, diffs, file lists, etc.)
- **FR-006**: System MUST handle git command errors by categorizing them into distinct error types (git not found, repository not found, permission denied, invalid repository, command failed, etc.) while maintaining compatibility with existing go-git error formats
- **FR-007**: System MUST maintain the same GitRepository interface without breaking changes
- **FR-008**: System MUST support all existing git operations: GetRepositoryState, CreateCommit, StageAllFiles, CaptureStagingState, StageModifiedFiles, StageAllFilesIncludingUntracked, UnstageFiles
- **FR-009**: System MUST support commit signing with SSH keys using git CLI equivalent commands, including support for passphrase-protected keys via environment variable or configuration
- **FR-010**: System MUST compute diffs with 0 context lines (equivalent to `git diff --unified=0`)
- **FR-011**: System MUST handle binary file detection and exclude binary files from diff content
- **FR-012**: System MUST apply size limits to diffs (5000 character limit) and generate metadata for large files
- **FR-013**: System MUST support filtering staged files based on context (includeNewFiles flag)
- **FR-014**: System MUST handle empty repositories (no commits) gracefully
- **FR-015**: System MUST handle file status detection (added, modified, deleted, renamed, copied, unmerged)
- **FR-016**: System MUST validate that git executable is available and version 2.34.0 or higher before attempting operations
- **FR-017**: System MUST handle git config extraction (user.name, user.email, signing configuration) using git CLI commands
- **FR-018**: System MUST log all git command executions including command, arguments, exit codes, and execution time for debugging and troubleshooting
- **FR-019**: System MUST completely remove go-git library dependency from the project (no fallback to go-git)

### Key Entities *(include if feature involves data)*

- **Git Operation**: Represents a single git operation currently performed via go-git library, including its purpose, parameters, and expected output
- **CLI Command Mapping**: Maps a git operation to its equivalent git CLI command, including flags, arguments, and output parsing requirements
- **Repository State**: Current state of the git repository including staged and unstaged files with their diffs
- **Staging State**: Snapshot of which files are currently staged in the repository

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: All go-git library calls are documented with 100% coverage of current usage
- **SC-002**: All documented git operations have equivalent CLI commands mapped with complete command syntax
- **SC-003**: All existing tests pass with the new CLI-based implementation without modification
- **SC-004**: Git operations complete within 2x the time of go-git implementation (accounting for process overhead)
- **SC-005**: System correctly handles all edge cases (empty repo, detached HEAD, merge conflicts, missing git) with appropriate error messages
- **SC-006**: Commit signing works identically to go-git implementation for SSH-signed commits
- **SC-007**: Diff computation produces identical output format to go-git implementation (unified diff with 0 context)
- **SC-008**: File staging and unstaging operations work identically to go-git implementation
- **SC-009**: Repository state detection (staged/unstaged files) produces identical results to go-git implementation
- **SC-010**: System provides clear error messages when git executable is not available, version is too old, or git commands fail

## Assumptions

- Git executable is available in system PATH
- Git version is 2.34.0 or higher (required for SSH commit signing support)
- Repository structure follows standard git conventions (.git directory exists)
- System has appropriate filesystem permissions to execute git commands
- Context cancellation is properly propagated to git subprocesses
- Git CLI output format is stable and parseable
- SSH signing can be performed using git CLI with appropriate environment variables or git config, including passphrase-protected keys via environment variable or configuration
