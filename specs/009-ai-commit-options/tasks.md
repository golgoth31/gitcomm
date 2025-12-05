# Tasks: AI Commit Message Acceptance Options

**Input**: Design documents from `/specs/009-ai-commit-options/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: TDD approach is required per constitution. Tests must be written before implementation.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: No new project structure needed - modifying existing codebase

- [X] T001 Review existing code structure in `internal/ui/prompts.go` and `internal/service/commit_service.go`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core types and structures that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T002 [P] Define `AIMessageAcceptance` type with constants (AcceptAndCommit, AcceptAndEdit, Reject) in `internal/ui/prompts.go`
- [X] T003 [P] Implement `String()` method for `AIMessageAcceptance` type in `internal/ui/prompts.go`
- [X] T004 [P] Define `PrefilledCommitMessage` struct with fields (Type, Scope, Subject, Body, Footer) in `internal/ui/prompts.go` or `internal/model/commit_message.go`
- [X] T005 [P] Write unit tests for `AIMessageAcceptance` type in `internal/ui/prompts_test.go`
- [X] T006 [P] Write unit tests for `PrefilledCommitMessage` struct in `internal/ui/prompts_test.go` or `internal/model/commit_message_test.go`

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Accept and Commit Directly (Priority: P1) 🎯 MVP

**Goal**: Users can accept an AI-generated commit message and have the commit created immediately without additional prompts or manual editing.

**Independent Test**: User runs gitcomm with AI enabled, receives an AI-generated message, selects "accept and commit", and the commit is created immediately with the AI message. This can be fully tested independently and delivers immediate value by reducing steps for satisfied users.

### Tests for User Story 1 ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [X] T007 [P] [US1] Write unit test for `PromptAIMessageAcceptanceOptions` with AcceptAndCommit selection in `internal/ui/prompts_test.go`
- [X] T008 [P] [US1] Write unit test for `PromptAIMessageAcceptanceOptions` with invalid input handling in `internal/ui/prompts_test.go`
- [ ] T009 [P] [US1] Write unit test for `generateWithAI` AcceptAndCommit path in `internal/service/commit_service_test.go`
- [ ] T010 [P] [US1] Write unit test for commit failure handling after AcceptAndCommit in `internal/service/commit_service_test.go`
- [ ] T011 [P] [US1] Write integration test for AcceptAndCommit workflow in `test/integration/ai_acceptance_options_test.go`

### Implementation for User Story 1

- [X] T012 [US1] Implement `PromptAIMessageAcceptanceOptions` function that displays AI message and three options (1/2/3) in `internal/ui/prompts.go`
- [X] T013 [US1] Add input validation and re-prompting logic for invalid options in `PromptAIMessageAcceptanceOptions` in `internal/ui/prompts.go`
- [X] T014 [US1] Modify `generateWithAI` to call `PromptAIMessageAcceptanceOptions` instead of `PromptAIMessageAcceptance` in `internal/service/commit_service.go`
- [X] T015 [US1] Implement AcceptAndCommit path in `generateWithAI` that creates commit immediately in `internal/service/commit_service.go`
- [X] T016 [US1] Implement `handleCommitFailure` helper method for commit failure recovery (retry/edit/cancel) in `internal/service/commit_service.go`
- [X] T017 [US1] Add error handling for commit failures after AcceptAndCommit with staging state restoration in `internal/service/commit_service.go`
- [X] T018 [US1] Add `PromptCommitFailureChoice` function for retry/edit/cancel options in `internal/ui/prompts.go`
- [X] T019 [US1] Write unit tests for `PromptCommitFailureChoice` in `internal/ui/prompts_test.go`

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - Accept and Edit (Priority: P2)

**Goal**: Users can accept an AI-generated commit message and then edit specific fields (type, scope, subject, body, footer) with the AI values pre-filled.

**Independent Test**: User runs gitcomm with AI enabled, receives an AI-generated message, selects "accept and edit", and is presented with the manual commit message prompts where all fields are pre-filled with values from the AI message. User can then modify any field and proceed to commit. This can be fully tested independently and delivers value by streamlining the editing workflow.

### Tests for User Story 2 ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T020 [P] [US2] Write unit test for `PromptAIMessageAcceptanceOptions` with AcceptAndEdit selection in `internal/ui/prompts_test.go`
- [ ] T021 [P] [US2] Write unit test for `NewSelectListModelWithPreselection` with matching type in `internal/ui/select_list_test.go`
- [ ] T022 [P] [US2] Write unit test for `NewSelectListModelWithPreselection` with non-matching type in `internal/ui/select_list_test.go`
- [ ] T023 [P] [US2] Write unit test for `PromptCommitTypeWithPreselection` in `internal/ui/prompts_test.go`
- [ ] T024 [P] [US2] Write unit test for `PromptScopeWithDefault` in `internal/ui/prompts_test.go`
- [ ] T025 [P] [US2] Write unit test for `PromptSubjectWithDefault` in `internal/ui/prompts_test.go`
- [ ] T026 [P] [US2] Write unit test for `NewMultilineInputModelWithValue` in `internal/ui/multiline_input_test.go`
- [ ] T027 [P] [US2] Write unit test for `PromptBodyWithDefault` in `internal/ui/prompts_test.go`
- [ ] T028 [P] [US2] Write unit test for `PromptFooterWithDefault` in `internal/ui/prompts_test.go`
- [ ] T029 [P] [US2] Write unit test for `parseAIMessageToPrefilled` helper method in `internal/service/commit_service_test.go`
- [ ] T030 [P] [US2] Write unit test for `promptCommitMessage` with pre-filled values in `internal/service/commit_service_test.go`
- [ ] T031 [P] [US2] Write unit test for AcceptAndEdit path in `generateWithAI` in `internal/service/commit_service_test.go`
- [ ] T032 [P] [US2] Write integration test for AcceptAndEdit workflow in `test/integration/ai_acceptance_options_test.go`

### Implementation for User Story 2

- [X] T033 [US2] Implement `NewSelectListModelWithPreselection` constructor that finds matching type and sets SelectedIndex in `internal/ui/select_list.go`
- [X] T034 [US2] Implement `PromptCommitTypeWithPreselection` function that uses preselected type in `internal/ui/prompts.go`
- [X] T035 [US2] Implement `PromptScopeWithDefault` function that displays default value in prompt in `internal/ui/prompts.go`
- [X] T036 [US2] Implement `PromptSubjectWithDefault` function that displays default value in prompt in `internal/ui/prompts.go`
- [X] T037 [US2] Implement `NewMultilineInputModelWithValue` constructor that pre-populates model value in `internal/ui/multiline_input.go`
- [X] T038 [US2] Implement `PromptBodyWithDefault` function that uses pre-populated multiline model in `internal/ui/prompts.go`
- [X] T039 [US2] Implement `PromptFooterWithDefault` function that uses pre-populated multiline model in `internal/ui/prompts.go`
- [X] T040 [US2] Implement `parseAIMessageToPrefilled` helper method that converts AI message to PrefilledCommitMessage in `internal/service/commit_service.go`
- [X] T041 [US2] Modify `promptCommitMessage` to accept `*PrefilledCommitMessage` parameter in `internal/service/commit_service.go`
- [X] T042 [US2] Update `promptCommitMessage` to use pre-filled values when provided (type, scope, subject, body, footer) in `internal/service/commit_service.go`
- [X] T043 [US2] Implement AcceptAndEdit path in `generateWithAI` that calls `promptCommitMessage` with pre-filled values in `internal/service/commit_service.go`
- [X] T044 [US2] Add cancellation handling for AcceptAndEdit flow with staging state restoration in `internal/service/commit_service.go`
- [X] T045 [US2] Update all callers of `promptCommitMessage()` to pass `nil` for manual input in `internal/service/commit_service.go`

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - Reject and Start Over (Priority: P3)

**Goal**: Users can reject an AI-generated commit message and start the commit message creation process from scratch (either with a new AI generation or manual input).

**Independent Test**: User runs gitcomm with AI enabled, receives an AI-generated message, selects "reject", and is presented with the option to generate a new AI message or proceed with manual input. This can be fully tested independently and maintains existing functionality.

### Tests for User Story 3 ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T046 [P] [US3] Write unit test for `PromptAIMessageAcceptanceOptions` with Reject selection in `internal/ui/prompts_test.go`
- [ ] T047 [P] [US3] Write unit test for `PromptRejectChoice` function in `internal/ui/prompts_test.go`
- [ ] T048 [P] [US3] Write unit test for Reject path with new AI generation in `generateWithAI` in `internal/service/commit_service_test.go`
- [ ] T049 [P] [US3] Write unit test for Reject path with manual input in `generateWithAI` in `internal/service/commit_service_test.go`
- [ ] T050 [P] [US3] Write unit test for AI generation failure fallback to manual input in `internal/service/commit_service_test.go`
- [ ] T051 [P] [US3] Write integration test for Reject workflow with new AI in `test/integration/ai_acceptance_options_test.go`
- [ ] T052 [P] [US3] Write integration test for Reject workflow with manual input in `test/integration/ai_acceptance_options_test.go`

### Implementation for User Story 3

- [X] T053 [US3] Implement `PromptRejectChoice` function that prompts for new AI or manual input in `internal/ui/prompts.go`
- [X] T054 [US3] Implement Reject path in `generateWithAI` that calls `PromptRejectChoice` in `internal/service/commit_service.go`
- [X] T055 [US3] Implement new AI generation path after rejection (recursive call to `generateWithAI`) in `internal/service/commit_service.go`
- [X] T056 [US3] Add retry limit to prevent infinite recursion in `generateWithAI` in `internal/service/commit_service.go`
- [X] T057 [US3] Implement manual input path after rejection (call `promptCommitMessage` with nil) in `internal/service/commit_service.go`
- [X] T058 [US3] Implement AI generation failure fallback to manual input with error message in `internal/service/commit_service.go`

**Checkpoint**: All user stories should now be independently functional

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [ ] T059 [P] Add comprehensive integration tests for all three acceptance workflows in `test/integration/ai_acceptance_options_test.go`
- [ ] T060 [P] Add edge case tests for partial AI message parsing in `internal/service/commit_service_test.go`
- [ ] T061 [P] Add edge case tests for invalid commit type handling in `internal/ui/prompts_test.go`
- [ ] T062 [P] Add edge case tests for empty scope/body/footer pre-filling in `internal/ui/prompts_test.go`
- [ ] T063 [P] Add error handling tests for cancellation during AcceptAndEdit flow in `internal/service/commit_service_test.go`
- [X] T064 [P] Verify backward compatibility - existing `PromptAIMessageAcceptance` function still works (deprecated but functional)
- [X] T065 [P] Update documentation in `README.md` with new AI acceptance options
- [X] T066 [P] Update `CHANGELOG.md` with feature description
- [ ] T067 [P] Run quickstart.md validation scenarios
- [X] T068 [P] Code cleanup and refactoring - remove any unused code
- [ ] T069 [P] Performance validation - verify SC-001 (<5 seconds for AcceptAndCommit) and SC-002 (<30 seconds for AcceptAndEdit)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2 → P3)
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - Depends on US1 for `PromptAIMessageAcceptanceOptions` but can be implemented independently
- **User Story 3 (P3)**: Can start after Foundational (Phase 2) - Depends on US1 for `PromptAIMessageAcceptanceOptions` but can be implemented independently

### Within Each User Story

- Tests (if included) MUST be written and FAIL before implementation
- UI functions before service layer usage
- Helper methods before main workflow
- Core implementation before integration
- Story complete before moving to next priority

### Parallel Opportunities

- All Foundational tasks marked [P] can run in parallel (T002-T006)
- Once Foundational phase completes, all user stories can start in parallel (if team capacity allows)
- All tests for a user story marked [P] can run in parallel
- UI prompt functions within a story marked [P] can run in parallel (different files)
- Different user stories can be worked on in parallel by different team members

---

## Parallel Example: User Story 1

```bash
# Launch all tests for User Story 1 together:
Task: "Write unit test for PromptAIMessageAcceptanceOptions with AcceptAndCommit selection in internal/ui/prompts_test.go"
Task: "Write unit test for PromptAIMessageAcceptanceOptions with invalid input handling in internal/ui/prompts_test.go"
Task: "Write unit test for generateWithAI AcceptAndCommit path in internal/service/commit_service_test.go"
Task: "Write unit test for commit failure handling after AcceptAndCommit in internal/service/commit_service_test.go"
Task: "Write integration test for AcceptAndCommit workflow in test/integration/ai_acceptance_options_test.go"
```

---

## Parallel Example: User Story 2

```bash
# Launch all UI prompt function implementations together (different files):
Task: "Implement NewSelectListModelWithPreselection constructor in internal/ui/select_list.go"
Task: "Implement PromptScopeWithDefault function in internal/ui/prompts.go"
Task: "Implement PromptSubjectWithDefault function in internal/ui/prompts.go"
Task: "Implement NewMultilineInputModelWithValue constructor in internal/ui/multiline_input.go"
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
   - Developer A: User Story 1 (Accept and Commit)
   - Developer B: User Story 2 (Accept and Edit) - can start after T012-T013 complete
   - Developer C: User Story 3 (Reject) - can start after T012-T013 complete
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
- TDD approach: Write tests first, ensure they fail, then implement
- All file paths use absolute paths from repository root
