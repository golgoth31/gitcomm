# Tasks: Auto-Stage Modified Files and State Restoration

**Input**: Design documents from `/specs/002-auto-stage-restore/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: TDD approach is mandatory per gitcomm constitution. Tests must be written first and fail before implementation.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project structure verification and preparation

- [x] T001 Verify existing project structure matches plan.md requirements
- [x] T002 [P] Review existing GitRepository interface in internal/repository/git_repository.go
- [x] T003 [P] Review existing CommitService in internal/service/commit_service.go

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T004 [P] Add staging/restoration error types to internal/utils/errors.go (ErrStagingFailed, ErrRestorationFailed, ErrStagingStateInvalid, ErrRestorationPlanInvalid, ErrInterruptedDuringStaging)
- [x] T005 [P] Create StagingState model in internal/model/staging_state.go (StagedFiles, CapturedAt, RepositoryPath fields, IsEmpty, Contains, Diff methods)
- [x] T006 [P] Create AutoStagingResult model in internal/model/staging_state.go (StagedFiles, FailedFiles, Success, Duration fields, HasFailures, GetFailedFilePaths methods)
- [x] T007 [P] Create StagingFailure model in internal/model/staging_state.go (FilePath, Error, ErrorType fields)
- [x] T008 [P] Create RestorationPlan model in internal/model/staging_state.go (FilesToUnstage, PreCLIState, CurrentState fields, IsEmpty, Validate, Execute methods)

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Auto-Stage Modified Files on Launch (Priority: P1) 🎯 MVP

**Goal**: Automatically stage all modified files when CLI launches, enabling AI analysis of complete diff without manual staging steps.

**Independent Test**: Run CLI in git repository with modified files, verify all modified files are automatically staged before prompts, confirm staged diff is available for AI analysis.

### Tests for User Story 1 (TDD - Write First) ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T009 [P] [US1] Write unit test for StagingState model in internal/model/staging_state_test.go (IsEmpty, Contains, Diff methods)
- [ ] T010 [P] [US1] Write unit test for CaptureStagingState in internal/repository/git_repository_impl_test.go
- [ ] T011 [P] [US1] Write unit test for StageModifiedFiles in internal/repository/git_repository_impl_test.go
- [ ] T012 [P] [US1] Write integration test for auto-staging workflow in test/integration/auto_stage_test.go (modified files auto-staged, staging occurs before prompts)

### Implementation for User Story 1

- [x] T013 [US1] Extend GitRepository interface with CaptureStagingState method in internal/repository/git_repository.go
- [x] T014 [US1] Extend GitRepository interface with StageModifiedFiles method in internal/repository/git_repository.go
- [x] T015 [US1] Implement CaptureStagingState in internal/repository/git_repository_impl.go (capture pre-CLI staging state using go-git worktree.Status)
- [x] T016 [US1] Implement StageModifiedFiles in internal/repository/git_repository_impl.go (stage only modified files, return AutoStagingResult)
- [x] T017 [US1] Implement staging state management logic in internal/repository/staging_state.go (helper functions for state capture/restoration)
- [x] T018 [US1] Integrate auto-staging into CommitService.CreateCommit in internal/service/commit_service.go (capture state, stage modified files before workflow)
- [x] T019 [US1] Add error handling for staging failures in internal/service/commit_service.go (abort on failure, restore state, exit with error)
- [x] T020 [US1] Add logging for auto-staging operations in internal/service/commit_service.go (log staging start, success, failures)

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently - modified files are auto-staged on CLI launch

---

## Phase 4: User Story 2 - Auto-Stage Unmanaged Files with -a Flag (Priority: P2)

**Goal**: When `-a` flag is used, automatically stage both modified and untracked files, providing complete repository state for AI analysis.

**Independent Test**: Run CLI with `-a` flag in repository with both modified and untracked files, verify all files (modified + untracked) are staged, confirm complete diff is available for AI analysis.

### Tests for User Story 2 (TDD - Write First) ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T021 [P] [US2] Write unit test for StageAllFilesIncludingUntracked in internal/repository/git_repository_impl_test.go
- [ ] T022 [P] [US2] Write integration test for -a flag workflow in test/integration/auto_stage_test.go (modified + untracked files staged with -a flag)

### Implementation for User Story 2

- [x] T023 [US2] Extend GitRepository interface with StageAllFilesIncludingUntracked method in internal/repository/git_repository.go
- [x] T024 [US2] Implement StageAllFilesIncludingUntracked in internal/repository/git_repository_impl.go (stage modified + untracked files, return AutoStagingResult)
- [x] T025 [US2] Update CommitService.CreateCommit to use StageAllFilesIncludingUntracked when -a flag is set in internal/service/commit_service.go
- [x] T026 [US2] Add logging for -a flag staging operations in internal/service/commit_service.go

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently - auto-staging works with and without -a flag

---

## Phase 5: User Story 3 - Restore Staging State on Cancellation (Priority: P1)

**Goal**: When CLI exits without committing (cancellation, error, interruption), restore staging state to exactly what it was before CLI launch, ensuring no unintended changes persist.

**Independent Test**: Run CLI, allow auto-staging, cancel/exiting without committing, verify staging state matches pre-CLI state. Test all exit scenarios (Ctrl+C, rejection, errors).

### Tests for User Story 3 (TDD - Write First) ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T027 [P] [US3] Write unit test for RestorationPlan model in internal/model/staging_state_test.go (IsEmpty, Validate, Execute methods)
- [ ] T028 [P] [US3] Write unit test for UnstageFiles in internal/repository/git_repository_impl_test.go
- [ ] T029 [P] [US3] Write integration test for state restoration on cancellation in test/integration/staging_restore_test.go (cancel CLI, verify state restored)
- [ ] T030 [P] [US3] Write integration test for state restoration on error in test/integration/staging_restore_test.go (error prevents commit, verify state restored)
- [ ] T031 [P] [US3] Write integration test for preserving originally staged files in test/integration/staging_restore_test.go (files staged before CLI, only CLI-staged files restored)

### Implementation for User Story 3

- [x] T032 [US3] Extend GitRepository interface with UnstageFiles method in internal/repository/git_repository.go
- [x] T033 [US3] Implement UnstageFiles in internal/repository/git_repository_impl.go (unstage specified files, best-effort on failures)
- [x] T034 [US3] Implement RestorationPlan.Execute in internal/model/staging_state.go (validate plan, call UnstageFiles)
- [x] T035 [US3] Add restoration logic to CommitService.CreateCommit on cancellation in internal/service/commit_service.go (create RestorationPlan, execute on exit without commit)
- [x] T036 [US3] Add restoration logic to CommitService.CreateCommit on error in internal/service/commit_service.go (restore state if commit creation fails)
- [x] T037 [US3] Add signal handling for interruption in cmd/gitcomm/main.go (register SIGINT/SIGTERM handlers, restore state on interruption)
- [x] T038 [US3] Add error handling for restoration failures in internal/service/commit_service.go (log error, display warning, exit with error)
- [x] T039 [US3] Add logging for restoration operations in internal/service/commit_service.go (log restoration start, success, failures)

**Checkpoint**: At this point, all user stories should be independently functional - auto-staging works and state is restored on cancellation/error/interruption

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [x] T040 [P] Add comprehensive error messages for staging/restoration errors in internal/utils/errors.go
- [x] T041 [P] Update README.md with auto-staging and state restoration documentation
- [x] T042 [P] Update CHANGELOG.md entry for this feature
- [x] T043 Code cleanup and refactoring (review all files for consistency)
- [x] T044 [P] Add additional unit tests for edge cases in test/unit/staging_edge_cases_test.go (empty state, external changes, partial failures)
- [x] T045 [P] Add integration tests for signal interruption scenarios in test/integration/signal_interrupt_test.go (interrupt during staging, interrupt during restoration)
- [x] T046 Run quickstart.md validation scenarios
- [x] T047 Performance validation (verify SC-001, SC-003, SC-006 timing requirements)
- [x] T048 Security review (verify no secrets in logs, validate file path handling)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (US1 → US2 → US3, but note US1 and US3 are both P1)
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - Depends on US1 (uses same staging infrastructure)
- **User Story 3 (P1)**: Can start after Foundational (Phase 2) - Depends on US1 (needs staging state capture from US1)

### Within Each User Story

- Tests (TDD) MUST be written and FAIL before implementation
- Models before services
- Services before integration
- Interface extensions before implementations
- Core implementation before error handling/logging
- Story complete before moving to next priority

### Parallel Opportunities

- All Foundational tasks marked [P] can run in parallel (T004-T008)
- All tests for a user story marked [P] can run in parallel
- Models within a story marked [P] can run in parallel
- Different user stories can be worked on in parallel by different team members (after Foundational)
- Polish tasks marked [P] can run in parallel

---

## Parallel Example: User Story 1

```bash
# Launch all tests for User Story 1 together:
Task: "Write unit test for StagingState model in internal/model/staging_state_test.go"
Task: "Write unit test for CaptureStagingState in internal/repository/git_repository_impl_test.go"
Task: "Write unit test for StageModifiedFiles in internal/repository/git_repository_impl_test.go"
Task: "Write integration test for auto-staging workflow in test/integration/auto_stage_test.go"

# Launch interface extensions together:
Task: "Extend GitRepository interface with CaptureStagingState method in internal/repository/git_repository.go"
Task: "Extend GitRepository interface with StageModifiedFiles method in internal/repository/git_repository.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Test User Story 1 independently
5. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Deploy/Demo (MVP!)
3. Add User Story 3 → Test independently → Deploy/Demo (Complete P1 stories)
4. Add User Story 2 → Test independently → Deploy/Demo
5. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 (P1)
   - Developer B: User Story 3 (P1) - can start after US1 captures staging state
   - Developer C: User Story 2 (P2) - can start after US1 completes
3. Stories complete and integrate independently

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing (TDD)
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence
- US1 and US3 are both P1 - implement US1 first as it provides foundation for US3
- US2 depends on US1 infrastructure but can be implemented in parallel after US1 core is done
