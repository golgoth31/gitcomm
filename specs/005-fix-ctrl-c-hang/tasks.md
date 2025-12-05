# Tasks: Fix CLI Hang on Ctrl+C During State Restoration

**Input**: Design documents from `/specs/005-fix-ctrl-c-hang/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Tests are included following TDD approach as required by the constitution.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **CLI project**: `cmd/gitcomm/`, `internal/` at repository root
- Paths shown below follow existing project structure

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Verify existing code structure and identify modification points

- [X] T001 Verify existing signal handling implementation in cmd/gitcomm/main.go
- [X] T002 Verify existing restoration logic in internal/service/commit_service.go
- [X] T003 [P] Verify git operations support context cancellation in internal/repository/git_repository_impl.go

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

### Tests for Timeout Behavior (TDD - Write First)

- [X] T004 [P] Write unit test for timeout context creation in internal/service/commit_service_test.go
- [X] T005 [P] Write unit test for timeout error detection in internal/service/commit_service_test.go
- [X] T006 [P] Write integration test for signal handling with timeout in test/integration/signal_timeout_test.go
- [X] T007 [P] Write integration test for restoration timeout scenario in test/integration/signal_timeout_test.go
- [X] T008 [P] Write integration test for multiple Ctrl+C handling in test/integration/signal_timeout_test.go

**Checkpoint**: Foundation ready - tests written and failing, user story implementation can now begin

---

## Phase 3: User Story 1 - Graceful Interruption with Timely Exit (Priority: P1) 🎯 MVP

**Goal**: Fix CLI hang on Ctrl+C by adding 3-second timeout to restoration operations and ensuring main process waits for completion or timeout before exiting.

**Independent Test**: Run `gitcomm`, press Ctrl+C during any phase, verify CLI exits within 5 seconds with exit code 130, and verify timeout warning appears if restoration exceeds 3 seconds.

### Implementation for User Story 1

- [X] T009 [US1] Modify restoreStagingState to accept timeout context parameter in internal/service/commit_service.go
- [X] T010 [US1] Create timeout context (3 seconds) when restoration is triggered by Ctrl+C in internal/service/commit_service.go
- [X] T011 [US1] Replace context.Background() with timeout context in defer function in internal/service/commit_service.go
- [X] T012 [US1] Add timeout error detection using errors.Is(err, context.DeadlineExceeded) in internal/service/commit_service.go
- [X] T013 [US1] Add warning message for timeout scenario in internal/service/commit_service.go
- [X] T014 [US1] Ensure git operations respect timeout context in internal/repository/git_repository_impl.go
- [X] T015 [US1] Add channel-based synchronization between signal handler and main process in cmd/gitcomm/main.go
- [X] T016 [US1] Ensure main process waits for restoration completion or timeout before exiting in cmd/gitcomm/main.go
- [X] T017 [US1] Add overall 5-second timeout enforcement in cmd/gitcomm/main.go
- [X] T018 [US1] Handle multiple Ctrl+C presses gracefully (ignore subsequent presses) in cmd/gitcomm/main.go

**Checkpoint**: At this point, User Story 1 should be fully functional - CLI exits within 5 seconds of Ctrl+C, restoration times out after 3 seconds if not complete, and warning messages are displayed appropriately

---

## Phase 4: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect the entire feature

- [X] T019 [P] Update README.md with timeout behavior documentation
- [X] T020 [P] Update CHANGELOG.md with bug fix entry
- [X] T021 Run all unit tests to verify no regressions
- [X] T022 Run all integration tests to verify no regressions
- [X] T023 [P] Verify backward compatibility (restoration without interruption still works)
- [X] T024 [P] Verify TUI interruption handling works correctly with timeout
- [X] T025 Code cleanup and formatting (gofmt, goimports)
- [X] T026 Run golangci-lint to verify code quality
- [X] T027 [P] Run quickstart.md validation scenarios

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can then proceed sequentially (only one story in this feature)
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories

### Within Each User Story

- Tests (if included) MUST be written and FAIL before implementation
- Core implementation before integration
- Story complete before moving to polish phase

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational test tasks marked [P] can run in parallel (within Phase 2)
- Polish tasks marked [P] can run in parallel

---

## Parallel Example: User Story 1

```bash
# Note: User Story 1 implementation tasks are mostly sequential (same file modifications)
# However, some tasks can be done in parallel if they modify different files:

# Can run in parallel (different files):
Task: "Ensure git operations respect timeout context in internal/repository/git_repository_impl.go"
Task: "Add channel-based synchronization between signal handler and main process in cmd/gitcomm/main.go"
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
3. Polish phase → Final validation → Deploy

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: Timeout context implementation in commit_service.go
   - Developer B: Signal handler synchronization in main.go
   - Developer C: Git operations context propagation in git_repository_impl.go
3. Integration and testing together

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence
- Total tasks: 27
- Tasks per story: US1 (10 tasks), Setup (3 tasks), Foundational (5 tests), Polish (9 tasks)
