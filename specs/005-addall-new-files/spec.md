# Feature Specification: Respect addAll Flag for New Files

**Feature Branch**: `005-addall-new-files`
**Created**: 2025-01-27
**Status**: Draft
**Input**: User description: "the addAll flag is not taken into account, new files are always taken into account by the git repository implementation. change this behaviour to add or not the new files depending on the flag state"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Exclude New Files When addAll is False (Priority: P1)

When a user runs gitcomm without the `-a` (add-all) flag, the system should exclude new/untracked files from commit consideration, even if they exist in the worktree. Only modified, deleted, or renamed files that are already tracked by git should be included in the commit workflow.

**Why this priority**: This is the core functionality requested - the addAll flag must control whether new files are considered. Without this, users cannot selectively commit only changes to existing files.

**Independent Test**: Can be fully tested by creating a new untracked file in a git repository, running `gitcomm` without the `-a` flag, and verifying the new file is not included in the repository state or commit process.

**Acceptance Scenarios**:

1. **Given** a git repository with modified tracked files and new untracked files, **When** user runs `gitcomm` without `-a` flag, **Then** only modified tracked files are included in the commit workflow
2. **Given** a git repository with only new untracked files (no modified files), **When** user runs `gitcomm` without `-a` flag, **Then** system reports no changes to commit
3. **Given** a git repository with new untracked files, **When** user runs `gitcomm -a`, **Then** new files are included in the commit workflow

---

### User Story 2 - Include New Files When addAll is True (Priority: P1)

When a user runs gitcomm with the `-a` (add-all) flag, the system should include all files (modified, deleted, renamed, and new/untracked) in the commit workflow.

**Why this priority**: This maintains existing behavior when the flag is set, ensuring backward compatibility and the expected "add all" functionality.

**Independent Test**: Can be fully tested by creating new untracked files in a git repository, running `gitcomm -a`, and verifying the new files are included in the repository state and commit process.

**Acceptance Scenarios**:

1. **Given** a git repository with modified tracked files and new untracked files, **When** user runs `gitcomm -a`, **Then** both modified and new files are included in the commit workflow
2. **Given** a git repository with only new untracked files, **When** user runs `gitcomm -a`, **Then** new files are included and commit workflow proceeds

---

### Edge Cases

- What happens when a new file is manually staged before running gitcomm without `-a` flag? (Should be excluded from consideration)
- How does system handle new files that are binary? (Should follow same exclusion rules as text files)
- What happens when repository has both staged new files and unstaged new files? (Only consider based on addAll flag, not staging status)
- How does system handle renamed files that are effectively new? (Renamed files should be treated as tracked, not new)

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST exclude new/untracked files from repository state when `addAll` flag is false
- **FR-002**: System MUST include new/untracked files in repository state when `addAll` flag is true
- **FR-003**: System MUST determine file status (new vs. modified) based on git tracking status, not staging status
- **FR-004**: System MUST filter repository state results based on `addAll` flag before processing commit workflow
- **FR-005**: System MUST maintain existing behavior for modified, deleted, and renamed tracked files regardless of `addAll` flag state

### Key Entities

- **RepositoryState**: Represents the current state of files in the repository, must respect `addAll` flag when including new files
- **FileChange**: Represents a single file change, includes status indicating whether file is new/untracked or modified/tracked

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: When `addAll` is false, 100% of new/untracked files are excluded from repository state
- **SC-002**: When `addAll` is true, 100% of new/untracked files are included in repository state
- **SC-003**: Modified tracked files are included in repository state regardless of `addAll` flag state (100% consistency)
- **SC-004**: Users can successfully create commits excluding new files when `addAll` is false, completing the workflow in the same time as before
- **SC-005**: No regression in commit workflow performance - commit creation time remains within 10% of baseline

## Assumptions

- New/untracked files are identified by git status as `Untracked` status
- Modified tracked files are identified by git status as `Modified`, `Deleted`, or `Renamed` status
- The `addAll` flag state is available in the commit service when determining repository state
- Users understand that `-a` flag controls whether new files are included, similar to `git add -A` behavior
