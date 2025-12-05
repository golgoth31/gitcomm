# Tasks: Compute Worktree Diff in GetRepositoryState

**Input**: Design documents from `/specs/014-worktree-diff/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: TDD is mandatory per constitution - all core business logic must have tests written first.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2)
- Include exact file paths in descriptions

## Path Conventions

- **Single project**: Repository root structure per plan.md
- Paths use absolute structure: `internal/`, `test/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and constants

- [X] T001 Define diff computation constants in `internal/repository/git_repository_impl.go`: `maxDiffSize = 5000` (characters) and `diffContext = 0` (lines)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [ ] T002 [P] Write unit test for `computeFileDiff()` helper function with modified file in `internal/repository/git_repository_impl_test.go` - test should fail initially
- [ ] T003 [P] Write unit test for `computeFileDiff()` helper function with added file in `internal/repository/git_repository_impl_test.go` - test should fail initially
- [ ] T004 [P] Write unit test for `computeFileDiff()` helper function with deleted file in `internal/repository/git_repository_impl_test.go` - test should fail initially
- [ ] T005 [P] Write unit test for `computeFileDiff()` helper function with binary file in `internal/repository/git_repository_impl_test.go` - test should fail initially
- [ ] T006 [P] Write unit test for `computeFileDiff()` helper function with empty repository (no HEAD) in `internal/repository/git_repository_impl_test.go` - test should fail initially
- [ ] T007 [P] Write unit test for `generateMetadata()` helper function for large files in `internal/repository/git_repository_impl_test.go` - test should fail initially
- [ ] T008 [P] Write unit test for `formatRenameDiff()` helper function with similarity percentage in `internal/repository/git_repository_impl_test.go` - test should fail initially
- [ ] T009 [P] Write unit test for `formatCopyDiff()` helper function with similarity percentage in `internal/repository/git_repository_impl_test.go` - test should fail initially
- [X] T010 Implement `getHEADTree()` helper function in `internal/repository/git_repository_impl.go` - returns HEAD tree or empty tree if no HEAD
- [X] T011 Implement `getStagedIndexTree()` helper function in `internal/repository/git_repository_impl.go` - converts staged index to tree for diff computation
- [X] T012 Implement `computeFileDiff()` helper function in `internal/repository/git_repository_impl.go` - computes diff between HEAD and staged state for a single file using go-git diff API
- [X] T013 Implement `formatUnifiedDiff()` helper function in `internal/repository/git_repository_impl.go` - formats diff as unified patch with 0 context lines
- [X] T014 Implement `isBinaryFile()` helper function in `internal/repository/git_repository_impl.go` - detects binary files using go-git's IsBinary() method
- [X] T015 Implement `generateMetadata()` helper function in `internal/repository/git_repository_impl.go` - generates metadata string (file size, line count, change summary) for large files/diffs
- [X] T016 Implement `formatRenameDiff()` helper function in `internal/repository/git_repository_impl.go` - formats rename diff with similarity percentage (e.g., "rename from X, rename to Y, similarity Z%")
- [X] T017 Implement `formatCopyDiff()` helper function in `internal/repository/git_repository_impl.go` - formats copy diff with similarity percentage (e.g., "copy from X, copy to Y, similarity Z%")
- [X] T018 Implement `applySizeLimit()` helper function in `internal/repository/git_repository_impl.go` - checks if diff exceeds 5000 characters and replaces with metadata if needed

**Checkpoint**: Foundation ready - helper functions complete, user story implementation can now begin

---

## Phase 3: User Story 1 - Diff Computation for Staged Files (Priority: P1) 🎯 MVP

**Goal**: Compute and include full diff content (patch format) for staged files when retrieving repository state, enabling AI commit message generation to have access to the actual code changes.

**Independent Test**: Stage a modified file with code changes, call GetRepositoryState, and verify the FileChange.Diff field contains the unified diff patch for that file.

### Tests for User Story 1 ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [X] T019 [P] [US1] Write integration test for GetRepositoryState with staged modified file in `test/integration/git_repository_diff_test.go` - verify Diff field contains unified diff patch, test should fail initially
- [X] T020 [P] [US1] Write integration test for GetRepositoryState with staged added file (new file) in `test/integration/git_repository_diff_test.go` - verify Diff field contains full diff content, test should fail initially
- [X] T021 [P] [US1] Write integration test for GetRepositoryState with staged deleted file in `test/integration/git_repository_diff_test.go` - verify Diff field contains deletion diff, test should fail initially
- [X] T022 [P] [US1] Write integration test for GetRepositoryState with staged renamed file in `test/integration/git_repository_diff_test.go` - verify Diff field contains rename diff with similarity, test should fail initially
- [X] T023 [P] [US1] Write integration test for GetRepositoryState with staged copied file in `test/integration/git_repository_diff_test.go` - verify Diff field contains copy diff with similarity, test should fail initially
- [X] T024 [P] [US1] Write integration test for GetRepositoryState with no staged changes in `test/integration/git_repository_diff_test.go` - verify all Diff fields are empty, test should fail initially
- [X] T025 [P] [US1] Write integration test for GetRepositoryState with unstaged files in `test/integration/git_repository_diff_test.go` - verify unstaged files have empty Diff fields, test should fail initially
- [X] T026 [P] [US1] Write unit test for GetRepositoryState populates Diff for staged files in `internal/repository/git_repository_impl_test.go` - test should fail initially
- [X] T027 [P] [US1] Write unit test for GetRepositoryState leaves Diff empty for unstaged files in `internal/repository/git_repository_impl_test.go` - test should fail initially

### Implementation for User Story 1

- [X] T028 [US1] Update `GetRepositoryState()` method in `internal/repository/git_repository_impl.go` to get HEAD tree using `getHEADTree()` helper
- [X] T029 [US1] Update `GetRepositoryState()` method in `internal/repository/git_repository_impl.go` to get staged index tree using `getStagedIndexTree()` helper
- [X] T030 [US1] Update `GetRepositoryState()` method in `internal/repository/git_repository_impl.go` to iterate through staged files and compute diff for each using `computeFileDiff()`
- [X] T031 [US1] Update `GetRepositoryState()` method in `internal/repository/git_repository_impl.go` to populate `FileChange.Diff` field for staged files with computed diff
- [X] T032 [US1] Update `GetRepositoryState()` method in `internal/repository/git_repository_impl.go` to ensure unstaged files have empty `Diff` field (FR-011)
- [X] T033 [US1] Update `GetRepositoryState()` method in `internal/repository/git_repository_impl.go` to handle binary files by setting `Diff = ""` (FR-013)
- [X] T034 [US1] Update `GetRepositoryState()` method in `internal/repository/git_repository_impl.go` to format renamed files using `formatRenameDiff()` helper (FR-014)
- [X] T035 [US1] Update `GetRepositoryState()` method in `internal/repository/git_repository_impl.go` to format copied files using `formatCopyDiff()` helper (FR-015)
- [X] T036 [US1] Update `GetRepositoryState()` method in `internal/repository/git_repository_impl.go` to apply size limit using `applySizeLimit()` helper (FR-016)
- [X] T037 [US1] Update `GetRepositoryState()` method in `internal/repository/git_repository_impl.go` to format unified diff with 0 context lines using `formatUnifiedDiff()` helper (FR-012)

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently. GetRepositoryState computes and populates diff content for staged files.

---

## Phase 4: User Story 2 - Diff Computation Accuracy (Priority: P2)

**Goal**: The computed diff accurately represents the difference between the current worktree state (including staged changes) and the clean worktree state (HEAD), ensuring AI models receive correct change information.

**Independent Test**: Stage specific changes to a file, call GetRepositoryState, and verify the diff content matches what `git diff --cached` would show for that file.

### Tests for User Story 2 ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [X] T038 [P] [US2] Write integration test for diff content matches git diff --cached output in `test/integration/git_repository_diff_test.go` - test should fail initially
- [X] T039 [P] [US2] Write integration test for multiple files with staged changes computed independently in `test/integration/git_repository_diff_test.go` - test should fail initially
- [X] T040 [P] [US2] Write integration test for staged file with unstaged modifications in `test/integration/git_repository_diff_test.go` - verify diff only includes staged changes, not unstaged, test should fail initially
- [X] T041 [P] [US2] Write integration test for binary files have empty Diff field in `test/integration/git_repository_diff_test.go` - test should fail initially
- [X] T042 [P] [US2] Write unit test for diff computation accuracy with specific changes in `internal/repository/git_repository_impl_test.go` - test should fail initially
- [X] T043 [P] [US2] Write unit test for large diff (>5000 chars) shows metadata only in `internal/repository/git_repository_impl_test.go` - test should fail initially
- [X] T044 [P] [US2] Write unit test for large new file (>5000 chars) shows metadata only in `internal/repository/git_repository_impl_test.go` - test should fail initially

### Implementation for User Story 2

- [X] T045 [US2] Update `computeFileDiff()` function in `internal/repository/git_repository_impl.go` to ensure diff matches git diff --cached format exactly
- [X] T046 [US2] Update `GetRepositoryState()` method in `internal/repository/git_repository_impl.go` to ensure each file's diff is computed independently (no cross-file dependencies)
- [X] T047 [US2] Update `GetRepositoryState()` method in `internal/repository/git_repository_impl.go` to ensure diff only includes staged changes (not unstaged modifications) (FR-007)
- [X] T048 [US2] Update `formatUnifiedDiff()` function in `internal/repository/git_repository_impl.go` to ensure output matches standard git diff format with 0 context lines
- [X] T049 [US2] Update `generateMetadata()` function in `internal/repository/git_repository_impl.go` to include accurate change summary (lines added, removed, modified)

**Checkpoint**: At this point, User Story 2 should be fully functional and testable independently. Diff computation is accurate and matches git diff --cached output.

---

## Phase 5: Error Handling & Edge Cases

**Purpose**: Handle error conditions and edge cases gracefully

- [X] T050 [P] Write unit test for file read error handling in `internal/repository/git_repository_impl_test.go` - verify logs error, sets Diff="", continues processing, test should fail initially
- [X] T051 [P] Write unit test for diff computation failure handling in `internal/repository/git_repository_impl_test.go` - verify logs error, sets Diff="", continues processing, test should fail initially
- [X] T052 [P] Write unit test for unmerged file handling in `internal/repository/git_repository_impl_test.go` - verify attempts diff, fallback to empty if fails, test should fail initially
- [X] T053 [P] Write integration test for unmerged files in `test/integration/git_repository_diff_test.go` - verify graceful handling, test should fail initially
- [X] T054 [P] Write integration test for empty repository (no HEAD) in `test/integration/git_repository_diff_test.go` - verify treats all staged files as new additions, test should fail initially
- [X] T055 Implement error handling in `GetRepositoryState()` method in `internal/repository/git_repository_impl.go` for file read errors - log error, set Diff="", continue (FR-005)
- [X] T056 Implement error handling in `GetRepositoryState()` method in `internal/repository/git_repository_impl.go` for diff computation failures - log error, set Diff="", continue (FR-010)
- [X] T057 Implement unmerged file handling in `GetRepositoryState()` method in `internal/repository/git_repository_impl.go` - attempt diff ignoring conflict markers, fallback to empty if fails (FR-008)
- [X] T058 Implement empty repository handling in `getHEADTree()` helper function in `internal/repository/git_repository_impl.go` - return empty tree if HEAD not found (FR-009)

---

## Phase 6: Performance & Polish

**Purpose**: Performance optimization and cross-cutting concerns

- [X] T059 [P] Write performance test for GetRepositoryState with 100 staged files in `test/integration/git_repository_diff_test.go` - verify completes in <2 seconds (SC-003), test should fail initially
- [X] T060 [P] Write error rate test for GetRepositoryState in `test/integration/git_repository_diff_test.go` - verify error rate <1% of files (SC-004), test should fail initially
- [X] T061 [P] Verify all existing GetRepositoryState tests still pass after diff computation changes
- [X] T062 [P] Update `README.md` to document diff computation feature
- [X] T063 [P] Update `CHANGELOG.md` with feature description and breaking changes (none)
- [X] T064 [P] Run quickstart.md validation to ensure implementation matches design
- [X] T065 [P] Code cleanup: remove any unused imports or dead code from modified files
- [X] T066 [P] Add code comments documenting diff computation logic and token optimization in `internal/repository/git_repository_impl.go`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Foundational phase completion
- **User Story 2 (Phase 4)**: Depends on Foundational phase completion, can use User Story 1 components
- **Error Handling (Phase 5)**: Depends on User Story 1 and User Story 2 completion
- **Performance & Polish (Phase 6)**: Depends on all previous phases completion

### Story Completion Order

1. **User Story 1 (P1)** - MVP: Must complete first - provides core diff computation functionality
2. **User Story 2 (P2)**: Can start after User Story 1 - enhances accuracy and edge cases

### Parallel Execution Opportunities

**Within Phase 2 (Foundational)**:
- T002-T009: All test tasks can run in parallel (different test cases)
- T010-T018: Implementation tasks can run sequentially (dependencies between helpers)

**Within Phase 3 (User Story 1)**:
- T019-T027: All test tasks can run in parallel (different test scenarios)
- T028-T037: Implementation tasks mostly sequential (building on previous steps)

**Within Phase 4 (User Story 2)**:
- T038-T044: All test tasks can run in parallel (different test scenarios)
- T045-T049: Implementation tasks mostly sequential (refinements)

**Within Phase 5 (Error Handling)**:
- T050-T054: All test tasks can run in parallel (different error scenarios)
- T055-T058: Implementation tasks can run in parallel (different error types)

**Within Phase 6 (Performance & Polish)**:
- T059-T066: Most tasks can run in parallel (different concerns)

---

## Implementation Strategy

### MVP Scope

**Minimum Viable Product**: User Story 1 (Phase 3) only
- Provides core diff computation for staged files
- Handles basic file types (modified, added, deleted)
- Includes size limiting and token optimization
- Can be delivered independently and tested

### Incremental Delivery

1. **Increment 1 (MVP)**: User Story 1 - Core diff computation
   - Delivers: Diff computation for staged files with token optimization
   - Testable: Independent test criteria met
   - Value: AI models receive diff content for commit message generation

2. **Increment 2**: User Story 2 - Accuracy enhancements
   - Delivers: Accurate diff computation matching git diff --cached
   - Testable: Independent test criteria met
   - Value: Ensures correct change information for AI models

3. **Increment 3**: Error handling and edge cases
   - Delivers: Graceful error handling and edge case support
   - Testable: All edge cases covered
   - Value: Robust production-ready feature

4. **Increment 4**: Performance and polish
   - Delivers: Performance optimization and documentation
   - Testable: Performance requirements met
   - Value: Production-ready with full documentation

---

## Task Summary

- **Total Tasks**: 66
- **Phase 1 (Setup)**: 1 task
- **Phase 2 (Foundational)**: 17 tasks (9 tests, 8 implementation)
- **Phase 3 (User Story 1)**: 19 tasks (9 tests, 10 implementation)
- **Phase 4 (User Story 2)**: 12 tasks (7 tests, 5 implementation)
- **Phase 5 (Error Handling)**: 9 tasks (5 tests, 4 implementation)
- **Phase 6 (Performance & Polish)**: 8 tasks

### Parallel Opportunities

- **Test tasks**: 30 tasks can run in parallel (marked with [P])
- **Implementation tasks**: Limited parallelization due to dependencies
- **Polish tasks**: 8 tasks can run in parallel

### Independent Test Criteria

- **User Story 1**: Stage a modified file, call GetRepositoryState, verify FileChange.Diff contains unified diff patch
- **User Story 2**: Stage specific changes, call GetRepositoryState, verify diff matches `git diff --cached` output

### Suggested MVP Scope

**MVP = User Story 1 (Phase 3) only**
- Provides core functionality
- Independently testable
- Delivers immediate value
- Can be extended with User Story 2 in next increment
