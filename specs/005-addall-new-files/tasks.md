# Tasks: Respect addAll Flag for New Files

**Input**: Design documents from `/specs/005-addall-new-files/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: TDD approach is mandatory per gitcomm constitution. Tests must be written first and fail before implementation.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project structure verification and preparation

- [x] T001 Verify existing project structure matches plan.md requirements
- [x] T002 [P] Review existing GitRepository interface in internal/repository/git_repository.go
- [x] T003 [P] Review existing CommitService in internal/service/commit_service.go
- [x] T004 [P] Review existing GetRepositoryState implementation in internal/repository/git_repository_impl.go

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T005 [P] Define context key type and constant for includeNewFilesKey in internal/repository/git_repository_impl.go (type contextKey string, const includeNewFilesKey contextKey = "includeNewFiles")

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Exclude New Files When addAll is False (Priority: P1) 🎯 MVP

**Goal**: When a user runs gitcomm without the `-a` (add-all) flag, the system should exclude new/untracked files from commit consideration, even if they exist in the worktree. Only modified, deleted, or renamed files that are already tracked by git should be included in the commit workflow.

**Independent Test**: Can be fully tested by creating a new untracked file in a git repository, running `gitcomm` without the `-a` flag, and verifying the new file is not included in the repository state or commit process.

**Acceptance Scenarios**:
1. **Given** a git repository with modified tracked files and new untracked files, **When** user runs `gitcomm` without `-a` flag, **Then** only modified tracked files are included in the commit workflow
2. **Given** a git repository with only new untracked files (no modified files), **When** user runs `gitcomm` without `-a` flag, **Then** system reports no changes to commit

### Tests for User Story 1 (TDD - Write First) ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T006 [P] [US1] Write unit test TestGetRepositoryState_ExcludesNewFilesWhenAddAllFalse in internal/repository/git_repository_impl_test.go (verify new files with git.Added status are excluded when includeNewFiles is false)
- [x] T007 [P] [US1] Write unit test TestGetRepositoryState_IncludesModifiedFilesWhenAddAllFalse in internal/repository/git_repository_impl_test.go (verify modified files are always included regardless of flag)
- [x] T008 [P] [US1] Write unit test TestGetRepositoryState_IncludesDeletedFilesWhenAddAllFalse in internal/repository/git_repository_impl_test.go (verify deleted files are always included regardless of flag)
- [x] T009 [P] [US1] Write unit test TestGetRepositoryState_IncludesRenamedFilesWhenAddAllFalse in internal/repository/git_repository_impl_test.go (verify renamed files are always included regardless of flag)
- [x] T010 [P] [US1] Write unit test TestGetRepositoryState_ExcludesManuallyStagedNewFiles in internal/repository/git_repository_impl_test.go (verify manually staged new files are excluded when addAll is false)
- [x] T011 [P] [US1] Write unit test TestGetRepositoryState_ExcludesBinaryNewFiles in internal/repository/git_repository_impl_test.go (verify binary new files follow same exclusion rules as text files)
- [x] T012 [P] [US1] Write integration test TestCommitService_ExcludesNewFilesWithoutAddAllFlag in test/integration/commit_service_test.go (end-to-end: gitcomm without -a excludes new files)

### Implementation for User Story 1

- [x] T013 [US1] Extract includeNewFiles from context in GetRepositoryState method in internal/repository/git_repository_impl.go (default to true for backward compatibility)
- [x] T014 [US1] Add filtering logic to skip new files (git.Added status) when includeNewFiles is false in GetRepositoryState method in internal/repository/git_repository_impl.go
- [x] T015 [US1] Pass AutoStage flag via context in CommitService.CreateCommit method in internal/service/commit_service.go (set includeNewFilesKey context value before calling GetRepositoryState)

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently - new files are excluded when addAll flag is false

---

## Phase 4: User Story 2 - Include New Files When addAll is True (Priority: P1)

**Goal**: When a user runs gitcomm with the `-a` (add-all) flag, the system should include all files (modified, deleted, renamed, and new/untracked) in the commit workflow, maintaining existing behavior and backward compatibility.

**Independent Test**: Can be fully tested by creating new untracked files in a git repository, running `gitcomm -a`, and verifying the new files are included in the repository state and commit process.

**Acceptance Scenarios**:
1. **Given** a git repository with modified tracked files and new untracked files, **When** user runs `gitcomm -a`, **Then** both modified and new files are included in the commit workflow
2. **Given** a git repository with only new untracked files, **When** user runs `gitcomm -a`, **Then** new files are included and commit workflow proceeds

### Tests for User Story 2 (TDD - Write First) ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T016 [P] [US2] Write unit test TestGetRepositoryState_IncludesNewFilesWhenAddAllTrue in internal/repository/git_repository_impl_test.go (verify new files with git.Added status are included when includeNewFiles is true)
- [x] T017 [P] [US2] Write unit test TestGetRepositoryState_DefaultBehaviorIncludesAll in internal/repository/git_repository_impl_test.go (verify backward compatibility - when context value not present, all files included)
- [x] T018 [P] [US2] Write integration test TestCommitService_IncludesNewFilesWithAddAllFlag in test/integration/commit_service_test.go (end-to-end: gitcomm -a includes new files)

### Implementation for User Story 2

- [x] T019 [US2] Verify filtering logic correctly includes new files when includeNewFiles is true in GetRepositoryState method in internal/repository/git_repository_impl.go (ensure default behavior and explicit true both work)
- [x] T020 [US2] Verify context value is correctly set when AutoStage is true in CommitService.CreateCommit method in internal/service/commit_service.go (ensure flag passing works correctly)

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently - filtering respects addAll flag in both directions

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Final validation, edge cases, performance, and documentation

### Edge Case Handling

- [x] T021 [P] Verify edge case: manually staged new file excluded when addAll is false (covered by T010)
- [x] T022 [P] Verify edge case: binary new files follow same exclusion rules (covered by T011)
- [x] T023 [P] Verify edge case: renamed files are never filtered (covered by T009)
- [x] T024 [P] Verify edge case: files staged then unstaged are not included (existing behavior, no change needed)

### Performance & Validation

- [x] T025 [P] Run performance benchmarks to verify no regression (filtering overhead <1ms for typical repositories)
- [x] T026 [P] Verify backward compatibility: existing callers without context value continue to work (all files included by default)
- [x] T027 [P] Run all existing tests to ensure no regressions in internal/repository/git_repository_impl_test.go
- [x] T028 [P] Run all existing tests to ensure no regressions in internal/service/commit_service_test.go

### Documentation

- [x] T029 [P] Update code comments in GetRepositoryState method in internal/repository/git_repository_impl.go to document filtering behavior
- [x] T030 [P] Update code comments in CommitService.CreateCommit method in internal/service/commit_service.go to document context value setting

---

## Dependencies

### User Story Completion Order

Both User Stories 1 and 2 are P1 priority and can be implemented in parallel since they test opposite behaviors of the same filtering logic. However, recommended order:

1. **User Story 1** (Exclude new files) - Core functionality, enables selective commits
2. **User Story 2** (Include new files) - Backward compatibility verification

### Task Dependencies

- **T005** (Context key definition) must complete before T013 (context extraction)
- **T013-T014** (Filtering logic) must complete before T015 (service integration)
- **T006-T012** (Tests for US1) should complete before T013-T015 (Implementation for US1)
- **T016-T018** (Tests for US2) should complete before T019-T020 (Implementation for US2)

### Parallel Execution Opportunities

**Phase 3 (US1 Tests)**: T006-T012 can all run in parallel (different test cases, same file)

**Phase 4 (US2 Tests)**: T016-T018 can all run in parallel (different test cases, same file)

**Phase 5 (Polish)**: T021-T030 can all run in parallel (different concerns, no dependencies)

**Cross-Phase**: Tests (T006-T012, T016-T018) can be written in parallel with foundational work (T005)

---

## Implementation Strategy

### MVP Scope

**Minimum Viable Product**: User Story 1 only
- Excludes new files when addAll is false
- Maintains existing behavior for modified/deleted/renamed files
- Enables selective commits without new files

**Full Feature**: User Story 1 + User Story 2
- Complete addAll flag behavior in both directions
- Full backward compatibility verification
- All edge cases handled

### Incremental Delivery

1. **Increment 1**: Foundation (T005) + US1 Tests (T006-T012) + US1 Implementation (T013-T015)
   - Delivers core functionality: exclude new files when flag is false
   - Independently testable and deployable

2. **Increment 2**: US2 Tests (T016-T018) + US2 Implementation (T019-T020)
   - Verifies backward compatibility: include new files when flag is true
   - Completes feature functionality

3. **Increment 3**: Polish (T021-T030)
   - Edge case validation
   - Performance verification
   - Documentation updates

### Testing Strategy

- **Unit Tests**: Test filtering logic in isolation with various file status combinations
- **Integration Tests**: Verify end-to-end behavior with real git repositories
- **Regression Tests**: Ensure no existing functionality breaks

---

## Summary

**Total Tasks**: 30
- **Phase 1 (Setup)**: 4 tasks
- **Phase 2 (Foundational)**: 1 task
- **Phase 3 (US1)**: 12 tasks (7 tests + 3 implementation)
- **Phase 4 (US2)**: 5 tasks (3 tests + 2 implementation)
- **Phase 5 (Polish)**: 8 tasks

**Task Count per User Story**:
- **User Story 1**: 12 tasks (7 tests, 3 implementation, 2 verification)
- **User Story 2**: 5 tasks (3 tests, 2 implementation)

**Parallel Opportunities**:
- All test tasks can run in parallel within their phase
- Polish tasks can all run in parallel
- Tests can be written in parallel with foundational work

**Independent Test Criteria**:
- **US1**: Create new file, run gitcomm without -a, verify file excluded
- **US2**: Create new file, run gitcomm -a, verify file included

**Suggested MVP Scope**: User Story 1 (exclude new files when addAll is false) - 12 tasks total

**Format Validation**: ✅ All tasks follow checklist format with checkbox, ID, optional [P] marker, optional [Story] label, and file paths
