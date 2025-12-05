# Tasks: Display Commit Type Selection Confirmation

**Input**: Design documents from `/specs/006-type-confirmation-display/`
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

- [X] T001 Verify existing PromptCommitType function in internal/ui/prompts.go
- [X] T002 Verify existing SelectListModel and GetSelectedType method in internal/ui/select_list.go
- [X] T003 [P] Verify existing test structure in internal/ui/prompts_test.go

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

### Tests for Confirmation Display (TDD - Write First)

- [X] T004 [P] Write unit test for confirmation line display in internal/ui/prompts_test.go
- [X] T005 [P] Write unit test for confirmation line format validation in internal/ui/prompts_test.go
- [X] T006 [P] Write unit test for no confirmation on cancellation in internal/ui/prompts_test.go
- [X] T007 [P] Write integration test for confirmation display workflow in test/integration/ui_confirmation_test.go

**Checkpoint**: Foundation ready - tests written and failing, user story implementation can now begin

---

## Phase 3: User Story 1 - Display Commit Type Confirmation (Priority: P1) 🎯 MVP

**Goal**: Display a confirmation line after commit type selection with format "✔ Choose a type(<scope>): <chosen type>" on a new line after alt-screen clears.

**Independent Test**: Run `gitcomm`, select a commit type from the interactive list, verify that a confirmation line appears with the format "✔ Choose a type(<scope>): <chosen type>" where <chosen type> is the actual selected type (e.g., "feat", "fix") before the next prompt appears.

### Implementation for User Story 1

- [X] T008 [US1] Add confirmation line display using fmt.Printf in internal/ui/prompts.go after GetSelectedType call
- [X] T009 [US1] Ensure confirmation line only displays when selection is successful (not cancelled) in internal/ui/prompts.go
- [X] T010 [US1] Verify confirmation line appears before function returns in internal/ui/prompts.go
- [X] T011 [US1] Verify confirmation line format matches specification exactly in internal/ui/prompts.go

**Checkpoint**: At this point, User Story 1 should be fully functional - confirmation line displays correctly for all valid selections, no confirmation on cancellation, and format matches specification

---

## Phase 4: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect the entire feature

- [X] T012 [P] Update README.md with confirmation display feature documentation
- [X] T013 [P] Update CHANGELOG.md with feature entry
- [X] T014 Run all unit tests to verify no regressions
- [X] T015 Run all integration tests to verify no regressions
- [X] T016 [P] Verify backward compatibility (existing commit workflow still works)
- [X] T017 [P] Verify confirmation works with all commit types (feat, fix, docs, style, refactor, test, chore, version)
- [X] T018 Code cleanup and formatting (gofmt, goimports)
- [X] T019 Run golangci-lint to verify code quality
- [X] T020 [P] Run quickstart.md validation scenarios

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
# However, some polish tasks can be done in parallel:

# Can run in parallel (different files):
Task: "Update README.md with confirmation display feature documentation"
Task: "Update CHANGELOG.md with feature entry"
Task: "Verify confirmation works with all commit types"
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
   - Developer A: Implementation in prompts.go
   - Developer B: Documentation updates (README, CHANGELOG)
   - Developer C: Integration testing and validation
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
- Total tasks: 20
- Tasks per story: US1 (4 tasks), Setup (3 tasks), Foundational (4 tests), Polish (9 tasks)
