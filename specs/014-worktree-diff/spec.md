# Feature Specification: Compute Worktree Diff in GetRepositoryState

**Feature Branch**: `014-worktree-diff`
**Created**: 2025-01-27
**Status**: Draft
**Input**: User description: "in the GetRepositoryState function, using the diff capability of the plumbing part of the git library, compute the full diff (patch type) between current git worktree (with staged modified or added files) and clean worktree (without staged modified or added files)"

## User Scenarios & Testing

### User Story 1 - Diff Computation for Staged Files (Priority: P1)

Users want gitcomm to compute and include full diff content (patch format) for staged files when retrieving repository state, enabling AI commit message generation to have access to the actual code changes.

**Why this priority**: This is fundamental to accurate commit message generation. Without diff content, AI models can only see file names and status, not the actual changes, leading to less accurate commit messages.

**Independent Test**: Stage a modified file with code changes, call GetRepositoryState, and verify the FileChange.Diff field contains the unified diff patch for that file.

**Acceptance Scenarios**:

1. **Given** a repository with staged modified files, **When** GetRepositoryState is called, **Then** each staged FileChange contains a Diff field with unified diff patch content showing changes between the staged version and the clean worktree version (metadata only if diff exceeds 5000 characters)
2. **Given** a repository with staged added files (new files), **When** GetRepositoryState is called, **Then** each staged FileChange contains a Diff field with file metadata (file size, line count) for files exceeding 5000 characters, or full diff content for files below 5000 characters
3. **Given** a repository with staged deleted files, **When** GetRepositoryState is called, **Then** each staged FileChange contains a Diff field with unified diff patch content showing the file removal
4. **Given** a repository with staged renamed files, **When** GetRepositoryState is called, **Then** each staged FileChange contains a Diff field with rename diff showing similarity percentage (e.g., "rename from old.txt, rename to new.txt, similarity 95%")
5. **Given** a repository with staged copied files, **When** GetRepositoryState is called, **Then** each staged FileChange contains a Diff field with copy diff showing similarity percentage (e.g., "copy from old.txt, copy to new.txt, similarity 100%")
6. **Given** a repository with no staged changes, **When** GetRepositoryState is called, **Then** all FileChange entries have empty Diff fields

---

### User Story 2 - Diff Computation Accuracy (Priority: P2)

Users want the computed diff to accurately represent the difference between the current worktree state (including staged changes) and the clean worktree state (HEAD), ensuring AI models receive correct change information.

**Why this priority**: Incorrect diff content leads to incorrect commit messages, which defeats the purpose of the feature and reduces user trust in the tool.

**Independent Test**: Stage specific changes to a file, call GetRepositoryState, and verify the diff content matches what `git diff --cached` would show for that file.

**Acceptance Scenarios**:

1. **Given** a file with staged modifications, **When** GetRepositoryState is called, **Then** the diff content matches the output of `git diff --cached` for that file
2. **Given** multiple files with staged changes, **When** GetRepositoryState is called, **Then** each file's diff is computed independently and accurately
3. **Given** a file with staged changes that also has unstaged modifications, **When** GetRepositoryState is called, **Then** the diff for the staged file only includes staged changes, not unstaged modifications
4. **Given** binary files that are staged, **When** GetRepositoryState is called, **Then** the diff field is empty (binary files have empty Diff content)

---

### Edge Cases

- What happens when a file is staged but the file content cannot be read? **Answer**: System logs an error, sets Diff to empty string, and continues processing other files (FR-005)
- What happens when computing diff for a very large file (e.g., >10MB) or very large diff? **Answer**: For files or diffs exceeding 5000 characters, system shows only metadata (file size, line count, change summary) instead of full content to limit token usage (FR-006, FR-016)
- What happens when a staged file is deleted from the worktree after staging? **Answer**: System computes diff based on the staged state (what was staged), not the current worktree state (FR-007)
- What happens when there are merge conflicts or unmerged files? **Answer**: System attempts to compute diff ignoring conflict markers (show diff of one side only), or sets Diff to empty if computation fails (FR-008)
- What happens when the repository HEAD is empty (initial commit)? **Answer**: System treats all staged files as new file additions and generates appropriate diff content (FR-009)
- What happens when diff computation fails for a specific file? **Answer**: System logs the error, sets Diff to empty string for that file, and continues processing other files (FR-010)
- What happens when a file is staged as renamed? **Answer**: System computes diff showing rename information with similarity percentage (e.g., "rename from old.txt, rename to new.txt, similarity 95%") (FR-014)
- What happens when a file is staged as copied? **Answer**: System computes diff showing copy information with similarity percentage (e.g., "copy from old.txt, copy to new.txt, similarity 100%") (FR-015)

## Requirements

### Functional Requirements

- **FR-001**: System MUST compute unified diff (patch format) for all staged files in GetRepositoryState function
- **FR-002**: System MUST use git library diff capabilities to compute diffs between current worktree state (with staged changes) and clean worktree state (HEAD)
- **FR-003**: System MUST populate the Diff field in FileChange for each staged file with the computed patch content
- **FR-004**: System MUST compute diff for staged modified files showing changes between staged version and HEAD version
- **FR-005**: System MUST handle file read errors gracefully by logging the error, setting Diff to empty string, and continuing with other files
- **FR-006**: System MUST compute diff for files of any size, but for large files (new additions or modified files with large diffs exceeding 5000 characters), show only metadata (file size, line count, change summary) instead of full content to limit token usage
- **FR-016**: System MUST use a character count threshold of 5000 characters for "large" files (files or diffs exceeding 5000 characters show metadata only, smaller files/diffs show full content)
- **FR-007**: System MUST compute diff based on staged state, not current worktree state (diff should reflect what is staged, not what is currently in the worktree)
- **FR-008**: System MUST handle unmerged files by attempting to compute diff ignoring conflict markers (show diff of one side only), or set Diff to empty if computation fails
- **FR-009**: System MUST handle empty repository (no HEAD) by treating all staged files as new file additions
- **FR-010**: System MUST handle diff computation failures for individual files by logging errors, setting Diff to empty, and continuing with other files
- **FR-011**: System MUST NOT compute diff for unstaged files (only staged files should have diff content)
- **FR-012**: System MUST use patch format (unified diff format) consistent with git diff output with 0 lines of context around each change (minimal token usage)
- **FR-013**: System MUST set Diff to empty string for binary files (no binary file indicators or content)
- **FR-014**: System MUST compute diff for staged renamed files showing rename information with similarity percentage (e.g., "rename from old.txt, rename to new.txt, similarity 95%")
- **FR-015**: System MUST compute diff for staged copied files showing copy information with similarity percentage (e.g., "copy from old.txt, copy to new.txt, similarity 100%")

### Key Entities

- **RepositoryState**: Represents the current state of the git repository, containing lists of staged and unstaged FileChange entries
- **FileChange**: Represents a single file change, containing Path, Status, and Diff fields (Diff field is populated for staged files)

## Success Criteria

### Measurable Outcomes

- **SC-001**: 100% of staged files in GetRepositoryState response have Diff field populated with accurate patch content
- **SC-002**: Diff content for staged files matches git diff --cached output with 100% accuracy for text files
- **SC-003**: GetRepositoryState completes diff computation for repositories with up to 100 staged files in under 2 seconds
- **SC-004**: System handles diff computation errors gracefully without failing the entire GetRepositoryState operation (error rate < 1% of files)
- **SC-005**: Diff content is in unified patch format and can be parsed by standard diff tools
- **SC-006**: Unstaged files have empty Diff fields (diff computation only occurs for staged files)

## Assumptions

- Git library diff capabilities support computing diffs between worktree and HEAD
- Staged file content is accessible through git library worktree API
- Diff computation performance is acceptable for typical repository sizes (up to 100 files)
- Binary files have empty Diff fields (no binary content or indicators)
- Users expect diff content to match standard git diff output format
- Empty diff fields for unstaged files is acceptable behavior (only staged files need diff content)

## Dependencies

- Existing GetRepositoryState function implementation
- Git library diff computation capabilities
- RepositoryState and FileChange data models (Diff field already exists in FileChange)
- Existing error handling and logging infrastructure

## Clarifications

### Session 2025-01-27

- Q: Should diff be computed for unstaged files as well, or only staged files? → A: Only staged files should have diff content computed (user requirement specifies "staged modified or added files")
- Q: What format should the diff be in? → A: Unified patch format (standard git diff format) (user requirement specifies "patch type")
- Q: Should binary files be included in diff computation? → A: System should handle binary files gracefully, either skipping or marking appropriately (informed guess based on standard git behavior)
- Q: What should happen if diff computation fails for a file? → A: Log error, set Diff to empty, continue with other files (informed guess based on error handling patterns)
- Q: When a file has merge conflicts (unmerged state), should the diff computation skip that file entirely, or should it include conflict markers in the diff output? → A: Attempt to compute diff ignoring conflict markers (show diff of one side only) (Option C)
- Q: For staged binary files (images, executables, etc.), should the Diff field be empty, or should it contain a binary file indicator message? → A: Set Diff to empty string for binary files (Option A)
- Q: When a file is staged as renamed (e.g., file1.txt → file2.txt), should the diff show a rename diff with similarity, separate delete/add diffs, or diff content only for the new file? → A: Show rename diff with similarity percentage (Option A)
- Q: How many lines of context should be included around each change in the unified diff? → A: Use 3 lines of context (standard unified diff format) (Option A)
- Q: When a file is staged as copied (e.g., file1.txt copied to file2.txt), should the diff show a copy diff with similarity, only the new file as addition, or both old and new file paths? → A: Show copy diff with similarity percentage (Option A)
- Q: To minimize tokens, should context lines be reduced from 3 to a smaller number, or eliminated entirely? → A: Use 0 lines of context (minimal, only changed lines shown) (Option B)
- Q: For staged new files (additions), should the diff include the entire file content, or should it be truncated/summarized to limit token usage? → A: Show only file metadata for large additions (file size, line count, no content) (Option C)
- Q: What size threshold should determine when a new file addition shows only metadata instead of full content? → A: Use character count threshold (user specified: "use char number")
- Q: What character count threshold value should be used? → A: 5000 characters (suggested and accepted)
- Q: Should large diffs for modified files also be truncated/summarized to limit token usage, or should all modified file diffs be included in full regardless of size? → A: Use same 5000 character threshold for modified file diffs as for new files (Option C)
- Q: When a diff exceeds 5000 characters, should it show first N lines with summary, only metadata, or summary only? → A: Show only metadata (same as large new files: file size, line count, change summary) (Option B)
- Q: Should there be a total aggregate limit across all staged files, or only per-file limits? → A: Per-file limits only (no aggregate limit, each file independently checked against 5000 character threshold) (Option B)

## Out of Scope

- Computing diff for unstaged files
- Computing diff between different commits or branches
- Diff formatting options or customization
- Diff size limits or truncation
- Real-time diff updates or caching
- Diff compression or optimization
- Support for custom diff algorithms
- Integration with external diff tools
