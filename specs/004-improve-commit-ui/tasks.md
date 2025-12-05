# Tasks: Improved Commit Message UI

**Input**: Design documents from `/specs/004-improve-commit-ui/`
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

**Purpose**: Project initialization and dependency verification

- [X] T001 Verify existing project structure and dependencies
- [X] T002 [P] Verify bubbletea dependency availability (add to go.mod if missing)
- [X] T003 [P] Verify lipgloss dependency availability (add to go.mod if missing)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core UI infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

### Tests for UI Components (TDD - Write First)

- [X] T004 [P] Write unit test for CommitTypeItem model in internal/ui/select_list_test.go
- [X] T005 [P] Write unit test for SelectListModel state transitions in internal/ui/select_list_test.go
- [X] T006 [P] Write unit test for MultilineInputModel completion detection in internal/ui/multiline_input_test.go
- [X] T007 [P] Write unit test for MultilineInputModel whitespace handling in internal/ui/multiline_input_test.go

### Implementation for UI Components

- [X] T008 Create CommitTypeItem struct in internal/ui/select_list.go
- [X] T009 Create SelectListModel struct with initial state in internal/ui/select_list.go
- [X] T010 Implement SelectListModel arrow key navigation (MoveUp, MoveDown) in internal/ui/select_list.go
- [X] T011 Implement SelectListModel letter-based navigation (JumpToLetter) in internal/ui/select_list.go
- [X] T012 Implement SelectListModel visual rendering with checkmarks and highlighting in internal/ui/select_list.go
- [X] T013 Create MultilineInputModel struct with initial state in internal/ui/multiline_input.go
- [X] T014 Implement MultilineInputModel Enter key handling (newline vs completion) in internal/ui/multiline_input.go
- [X] T015 Implement MultilineInputModel double-Enter completion detection in internal/ui/multiline_input.go
- [X] T016 Implement MultilineInputModel whitespace-only detection and trimming in internal/ui/multiline_input.go

**Checkpoint**: Foundation ready - UI components complete, user story implementation can now begin

---

## Phase 3: User Story 1 - Interactive Commit Type Selection (Priority: P1) 🎯 MVP

**Goal**: Replace numbered text-based commit type selection with an interactive select list with visual feedback (checkmarks, highlighting), arrow key navigation, and letter-based navigation.

**Independent Test**: Run `gitcomm` and verify that commit type selection displays as an interactive select list with checkmarks, highlighting, arrow key navigation works, letter navigation works, and first option (feat) is pre-selected.

### Tests for User Story 1 (TDD - Write First)

- [X] T017 [P] [US1] Write integration test for interactive commit type selection display in test/integration/ui_interactive_test.go
- [X] T018 [P] [US1] Write integration test for arrow key navigation in commit type selection in test/integration/ui_interactive_test.go
- [X] T019 [P] [US1] Write integration test for letter-based navigation in commit type selection in test/integration/ui_interactive_test.go
- [X] T020 [P] [US1] Write integration test for Enter key confirmation in commit type selection in test/integration/ui_interactive_test.go
- [X] T021 [P] [US1] Write integration test for Escape key cancellation in commit type selection in test/integration/ui_interactive_test.go

### Implementation for User Story 1

- [X] T022 [US1] Refactor PromptCommitType to use SelectListModel in internal/ui/prompts.go
- [X] T023 [US1] Integrate bubbletea program initialization for commit type selection in internal/ui/prompts.go
- [X] T024 [US1] Implement bubbletea model Update method for commit type selection in internal/ui/select_list.go
- [X] T025 [US1] Implement bubbletea model View method with lipgloss styling for commit type selection in internal/ui/select_list.go
- [X] T026 [US1] Configure SelectListModel to pre-select first option (feat) in internal/ui/select_list.go
- [X] T027 [US1] Add error handling for Escape key cancellation in internal/ui/prompts.go
- [X] T028 [US1] Update CommitService to handle new PromptCommitType return value in internal/service/commit_service.go

**Checkpoint**: At this point, User Story 1 should be fully functional - commit type selection uses interactive select list with all navigation features

---

## Phase 4: User Story 2 - Multiline Body Input (Priority: P1)

**Goal**: Enable proper multiline input for commit body field with double-Enter completion, blank line preservation, and whitespace handling.

**Independent Test**: Run `gitcomm` and verify that body prompt accepts multiple lines, Enter creates new lines, double Enter on empty lines completes input, single blank lines are preserved, and whitespace-only input is treated as empty.

### Tests for User Story 2 (TDD - Write First)

- [X] T029 [P] [US2] Write integration test for multiline body input acceptance in test/integration/ui_interactive_test.go
- [X] T030 [P] [US2] Write integration test for double-Enter completion in body input in test/integration/ui_interactive_test.go
- [X] T031 [P] [US2] Write integration test for blank line preservation in body input in test/integration/ui_interactive_test.go
- [X] T032 [P] [US2] Write integration test for whitespace-only body input handling in test/integration/ui_interactive_test.go

### Implementation for User Story 2

- [X] T033 [US2] Refactor PromptBody to use MultilineInputModel in internal/ui/prompts.go
- [X] T034 [US2] Integrate bubbletea textarea component for body input in internal/ui/multiline_input.go
- [X] T035 [US2] Implement double-Enter completion detection in MultilineInputModel for body in internal/ui/multiline_input.go
- [X] T036 [US2] Implement whitespace-only detection and trimming for body in internal/ui/multiline_input.go
- [X] T037 [US2] Add error handling for body input cancellation in internal/ui/prompts.go
- [X] T038 [US2] Update CommitService to handle multiline body input in internal/service/commit_service.go

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently - interactive commit type selection and multiline body input

---

## Phase 5: User Story 3 - Multiline Footer Input (Priority: P1)

**Goal**: Enable proper multiline input for commit footer field with double-Enter completion, blank line preservation, and whitespace handling.

**Independent Test**: Run `gitcomm` and verify that footer prompt accepts multiple lines, Enter creates new lines, double Enter on empty lines completes input, single blank lines are preserved, and whitespace-only input is treated as empty.

### Tests for User Story 3 (TDD - Write First)

- [X] T039 [P] [US3] Write integration test for multiline footer input acceptance in test/integration/ui_interactive_test.go
- [X] T040 [P] [US3] Write integration test for double-Enter completion in footer input in test/integration/ui_interactive_test.go
- [X] T041 [P] [US3] Write integration test for blank line preservation in footer input in test/integration/ui_interactive_test.go
- [X] T042 [P] [US3] Write integration test for whitespace-only footer input handling in test/integration/ui_interactive_test.go

### Implementation for User Story 3

- [X] T043 [US3] Refactor PromptFooter to use MultilineInputModel in internal/ui/prompts.go
- [X] T044 [US3] Reuse MultilineInputModel component for footer input in internal/ui/prompts.go
- [X] T045 [US3] Implement double-Enter completion detection in MultilineInputModel for footer in internal/ui/multiline_input.go
- [X] T046 [US3] Implement whitespace-only detection and trimming for footer in internal/ui/multiline_input.go
- [X] T047 [US3] Add error handling for footer input cancellation in internal/ui/prompts.go
- [X] T048 [US3] Update CommitService to handle multiline footer input in internal/service/commit_service.go

**Checkpoint**: At this point, all three user stories should work independently - interactive commit type selection, multiline body input, and multiline footer input

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [X] T049 [P] Update README.md with improved UI documentation
- [X] T050 [P] Update CHANGELOG.md with improved commit UI feature entry
- [X] T051 Run all unit tests to verify no regressions
- [X] T052 Run all integration tests to verify no regressions
- [X] T053 [P] Verify terminal resize handling works correctly for all UI components
- [X] T054 [P] Verify Escape key cancellation restores staging state correctly
- [X] T055 Code cleanup and formatting (gofmt, goimports)
- [X] T056 Run golangci-lint to verify code quality
- [X] T057 [P] Run quickstart.md validation scenarios

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (US1 → US2 → US3)
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories (reuses MultilineInputModel from foundational)
- **User Story 3 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories (reuses MultilineInputModel from foundational)

### Within Each User Story

- Tests (if included) MUST be written and FAIL before implementation
- UI component models before prompt integration
- Prompt integration before service layer updates
- Core implementation before integration
- Story complete before moving to next priority

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational test tasks marked [P] can run in parallel (within Phase 2)
- Foundational implementation tasks are sequential (they build on each other)
- Once Foundational phase completes, user stories can start
- All tests for a user story marked [P] can run in parallel
- Different user stories can be worked on in parallel by different team members

---

## Parallel Example: User Story 1

```bash
# Launch all tests for User Story 1 together:
Task: "Write integration test for interactive commit type selection display in test/integration/ui_interactive_test.go"
Task: "Write integration test for arrow key navigation in commit type selection in test/integration/ui_interactive_test.go"
Task: "Write integration test for letter-based navigation in commit type selection in test/integration/ui_interactive_test.go"
Task: "Write integration test for Enter key confirmation in commit type selection in test/integration/ui_interactive_test.go"
Task: "Write integration test for Escape key cancellation in commit type selection in test/integration/ui_interactive_test.go"
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
3. Add User Story 2 → Test independently → Deploy/Demo
4. Add User Story 3 → Test independently → Deploy/Demo
5. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 (interactive select list)
   - Developer B: User Story 2 (multiline body)
   - Developer C: User Story 3 (multiline footer)
3. Stories complete and integrate independently

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence
- Total tasks: 57
- Tasks per story: US1 (12 tasks), US2 (10 tasks), US3 (10 tasks), Setup (3 tasks), Foundational (13 tasks), Polish (9 tasks)
