# Tasks: Rewrite All CLI Prompts

**Input**: Design documents from `/specs/001-rewrite-cli-prompts/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: TDD approach required per constitution - tests written first, ensure they fail before implementation.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Setup

**Purpose**: Verify project structure and dependencies

- [X] T001 Verify bubbletea and lipgloss dependencies in go.mod
- [X] T002 Verify existing UI package structure in internal/ui/

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T003 [P] Create PromptState enum in internal/ui/prompt_state.go
- [X] T004 [P] Create VisualIndicator helper function in internal/ui/display.go
- [X] T005 [P] Write tests for PromptState in internal/ui/prompt_state_test.go
- [X] T006 [P] Write tests for VisualIndicator in internal/ui/display_test.go

**Checkpoint**: Foundation ready - user story implementation can now begin

---

## Phase 3: User Story 1 - Consistent Prompt Visual Design (Priority: P1) 🎯 MVP

**Goal**: Establish visual foundation with blue '?' and green '✓' indicators for all prompts

**Independent Test**: User runs gitcomm and interacts with any prompt. All prompts display blue '?' at start, and after answering, the '?' is replaced with a green checkmark.

### Tests for User Story 1

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [X] T007 [P] [US1] Write test for visual indicator rendering in prompt titles in internal/ui/display_test.go
- [X] T008 [P] [US1] Write test for state transition from Active to Completed in internal/ui/prompt_state_test.go
- [X] T009 [P] [US1] Write integration test for visual indicator in text input prompt in test/integration/visual_indicator_test.go

### Implementation for User Story 1

- [X] T010 [US1] Implement PromptState enum with all states (Pending, Active, Completed, Cancelled, Error) in internal/ui/prompt_state.go
- [X] T011 [US1] Implement GetVisualIndicator function with state-based rendering in internal/ui/display.go
- [X] T012 [US1] Add state field to SelectListModel in internal/ui/select_list.go
- [X] T013 [US1] Add state field to MultilineInputModel in internal/ui/multiline_input.go
- [X] T014 [US1] Update SelectListModel.View() to use visual indicators in internal/ui/select_list.go
- [X] T015 [US1] Update MultilineInputModel.View() to use visual indicators in internal/ui/multiline_input.go

**Checkpoint**: At this point, User Story 1 should be fully functional - all prompts show blue '?' and green '✓' indicators

---

## Phase 4: User Story 2 - Bubbletea for All Prompts (Priority: P1)

**Goal**: Migrate all prompts to use bubbletea for rendering and interaction

**Independent Test**: User runs gitcomm and interacts with all prompt types (text input, select lists, multiline input, yes/no choices). All prompts use bubbletea for rendering and interaction.

### Tests for User Story 2

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [X] T016 [P] [US2] Write tests for TextInputModel in internal/ui/text_input_test.go
- [X] T017 [P] [US2] Write tests for YesNoChoiceModel in internal/ui/yes_no_choice_test.go
- [X] T018 [P] [US2] Write tests for TextInputModel with validation in internal/ui/text_input_test.go
- [X] T019 [P] [US2] Write tests for TextInputModel with default values in internal/ui/text_input_test.go
- [X] T020 [P] [US2] Write tests for YesNoChoiceModel keyboard handling in internal/ui/yes_no_choice_test.go
- [X] T021 [P] [US2] Write integration test for text input prompt workflow in test/integration/text_input_test.go
- [X] T022 [P] [US2] Write integration test for yes/no choice prompt workflow in test/integration/yes_no_choice_test.go

### Implementation for User Story 2

- [X] T023 [P] [US2] Create TextInputModel struct implementing tea.Model in internal/ui/text_input.go
- [X] T024 [P] [US2] Create YesNoChoiceModel struct implementing tea.Model in internal/ui/yes_no_choice.go
- [X] T025 [US2] Implement TextInputModel.Init() method in internal/ui/text_input.go
- [X] T026 [US2] Implement TextInputModel.Update() method with keyboard handling in internal/ui/text_input.go
- [X] T027 [US2] Implement TextInputModel.View() method with visual indicators in internal/ui/text_input.go
- [X] T028 [US2] Implement TextInputModel validation logic in internal/ui/text_input.go
- [X] T029 [US2] Implement YesNoChoiceModel.Init() method in internal/ui/yes_no_choice.go
- [X] T030 [US2] Implement YesNoChoiceModel.Update() method with y/n key handling in internal/ui/yes_no_choice.go
- [X] T031 [US2] Implement YesNoChoiceModel.View() method with visual indicators in internal/ui/yes_no_choice.go
- [X] T032 [US2] Refactor PromptScope to use TextInputModel in internal/ui/prompts.go
- [X] T033 [US2] Refactor PromptScopeWithDefault to use TextInputModel in internal/ui/prompts.go
- [X] T034 [US2] Refactor PromptSubject to use TextInputModel with validation in internal/ui/prompts.go
- [X] T035 [US2] Refactor PromptSubjectWithDefault to use TextInputModel with validation in internal/ui/prompts.go
- [X] T036 [US2] Refactor PromptEmptyCommit to use YesNoChoiceModel in internal/ui/prompts.go
- [X] T037 [US2] Refactor PromptConfirm to use YesNoChoiceModel in internal/ui/prompts.go
- [X] T038 [US2] Refactor PromptAIUsage to use YesNoChoiceModel in internal/ui/prompts.go
- [X] T039 [US2] Refactor PromptAIMessageEdit to use YesNoChoiceModel in internal/ui/prompts.go
- [X] T040 [US2] Refactor PromptRejectChoice to use YesNoChoiceModel in internal/ui/prompts.go

**Checkpoint**: At this point, User Stories 1 AND 2 should both work - all prompts use bubbletea with visual indicators

---

## Phase 5: User Story 3 - No AltScreen Usage (Priority: P2)

**Goal**: Remove alt screen mode from all prompts, render inline

**Independent Test**: User runs gitcomm and interacts with prompts. All prompts render inline without switching to alt screen mode. Terminal history remains visible throughout the interaction.

### Tests for User Story 3

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [X] T041 [P] [US3] Write test to verify no alt screen in SelectListModel in internal/ui/select_list_test.go
- [X] T042 [P] [US3] Write test to verify no alt screen in MultilineInputModel in internal/ui/multiline_input_test.go
- [X] T043 [P] [US3] Write test to verify no alt screen in TextInputModel in internal/ui/text_input_test.go
- [X] T044 [P] [US3] Write test to verify no alt screen in YesNoChoiceModel in internal/ui/yes_no_choice_test.go
- [X] T045 [P] [US3] Write integration test for inline rendering in test/integration/inline_rendering_test.go

### Implementation for User Story 3

- [X] T046 [US3] Remove tea.WithAltScreen() from PromptCommitType in internal/ui/prompts.go
- [X] T047 [US3] Remove tea.WithAltScreen() from PromptCommitTypeWithPreselection in internal/ui/prompts.go
- [X] T048 [US3] Remove tea.WithAltScreen() from PromptBody in internal/ui/prompts.go
- [X] T049 [US3] Remove tea.WithAltScreen() from PromptBodyWithDefault in internal/ui/prompts.go
- [X] T050 [US3] Remove tea.WithAltScreen() from PromptFooter in internal/ui/prompts.go
- [X] T051 [US3] Remove tea.WithAltScreen() from PromptFooterWithDefault in internal/ui/prompts.go
- [X] T052 [US3] Verify all prompt functions use tea.NewProgram() without alt screen option in internal/ui/prompts.go

**Checkpoint**: At this point, User Stories 1, 2, AND 3 should work - all prompts render inline without alt screen

---

## Phase 6: User Story 4 - Select List Value Display (Priority: P2)

**Goal**: Display selected value after prompt title in select lists after confirmation

**Independent Test**: User runs gitcomm and selects a commit type from the select list. After selection, the selected value (e.g., "feat") is displayed immediately after the prompt title.

### Tests for User Story 4

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [X] T053 [P] [US4] Write test for confirmedValue field in SelectListModel in internal/ui/select_list_test.go
- [X] T054 [P] [US4] Write test for selected value display after confirmation in internal/ui/select_list_test.go
- [X] T055 [P] [US4] Write test for preselection value display in SelectListModel in internal/ui/select_list_test.go
- [X] T056 [P] [US4] Write integration test for select list value display in test/integration/select_list_display_test.go

### Implementation for User Story 4

- [X] T057 [US4] Add confirmedValue field to SelectListModel in internal/ui/select_list.go
- [X] T058 [US4] Update SelectListModel.Update() to set confirmedValue on Enter key in internal/ui/select_list.go
- [X] T059 [US4] Update SelectListModel.View() to show selected value after title only after confirmation in internal/ui/select_list.go
- [X] T060 [US4] Update SelectListModel.View() to show preselected value before confirmation in internal/ui/select_list.go
- [X] T061 [US4] Update PromptCommitType to display confirmed value in internal/ui/prompts.go
- [X] T062 [US4] Update PromptCommitTypeWithPreselection to display confirmed value in internal/ui/prompts.go

**Checkpoint**: At this point, User Stories 1-4 should work - select lists show selected value after title

---

## Phase 7: User Story 5 - Multiline Input Result Display (Priority: P2)

**Goal**: Display multiline input result under prompt title with line wrapping after validation

**Independent Test**: User runs gitcomm and enters text in the commit body multiline prompt. After validation (double Enter on empty line), the entered text is displayed under the prompt title.

### Tests for User Story 5

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [X] T063 [P] [US5] Write test for result display after completion in MultilineInputModel in internal/ui/multiline_input_test.go
- [X] T064 [P] [US5] Write test for line wrapping in result display in internal/ui/multiline_input_test.go
- [X] T065 [P] [US5] Write test for empty multiline input result display in internal/ui/multiline_input_test.go
- [X] T066 [P] [US5] Write integration test for multiline result display in test/integration/multiline_display_test.go

### Implementation for User Story 5

- [X] T067 [US5] Update MultilineInputModel.View() to show result under title after completion in internal/ui/multiline_input.go
- [X] T068 [US5] Implement line wrapping for result text using lipgloss in internal/ui/multiline_input.go
- [X] T069 [US5] Update MultilineInputModel to track completion state for result display in internal/ui/multiline_input.go
- [X] T070 [US5] Update PromptBody to display result after completion in internal/ui/prompts.go
- [X] T071 [US5] Update PromptBodyWithDefault to display result after completion in internal/ui/prompts.go
- [X] T072 [US5] Update PromptFooter to display result after completion in internal/ui/prompts.go
- [X] T073 [US5] Update PromptFooterWithDefault to display result after completion in internal/ui/prompts.go

**Checkpoint**: At this point, User Stories 1-5 should work - multiline inputs show result under title with wrapping

---

## Phase 8: Additional Prompt Functions

**Purpose**: Migrate remaining prompt functions that use select lists or custom logic

- [ ] T074 [P] Refactor PromptAIMessageAcceptanceOptions to use SelectListModel or custom model in internal/ui/prompts.go
- [ ] T075 [P] Refactor PromptCommitFailureChoice to use SelectListModel or custom model in internal/ui/prompts.go
- [ ] T076 [P] Update PromptAIMessageAcceptanceOptions to use visual indicators in internal/ui/prompts.go
- [ ] T077 [P] Update PromptCommitFailureChoice to use visual indicators in internal/ui/prompts.go

---

## Phase 9: Error Handling and Edge Cases

**Purpose**: Implement error states, cancellation, and validation error handling

- [ ] T078 [P] Implement error state handling in TextInputModel for validation errors in internal/ui/text_input.go
- [ ] T079 [P] Implement cancellation state handling in all models (Escape key) in internal/ui/
- [ ] T080 [P] Implement Ctrl+C interruption handling in all prompt functions in internal/ui/prompts.go
- [ ] T081 [P] Add error message display with yellow warning indicator in TextInputModel.View() in internal/ui/text_input.go
- [ ] T082 [P] Add cancellation indicator (red 'X') display in all models on Escape in internal/ui/
- [ ] T083 [P] Write tests for validation error display in internal/ui/text_input_test.go
- [ ] T084 [P] Write tests for cancellation handling in all model tests in internal/ui/
- [ ] T085 [P] Write integration tests for error scenarios in test/integration/error_handling_test.go

---

## Phase 10: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [ ] T086 [P] Update README.md with new prompt visual design documentation
- [ ] T087 [P] Update CHANGELOG.md with feature description
- [ ] T088 [P] Verify backward compatibility of all prompt function signatures in internal/ui/prompts.go
- [ ] T089 [P] Run gofmt and golangci-lint on all modified files
- [ ] T090 [P] Add comprehensive integration tests for all prompt types in test/integration/
- [ ] T091 [P] Test all prompt functions with default values
- [ ] T092 [P] Test all prompt functions with preselection
- [ ] T093 [P] Test terminal resize handling in all models
- [ ] T094 [P] Test narrow terminal width handling in all models
- [ ] T095 [P] Test long multiline input wrapping in MultilineInputModel
- [ ] T096 [P] Run quickstart.md validation
- [ ] T097 [P] Performance validation (prompt rendering <100ms)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-7)**: All depend on Foundational phase completion
  - User stories can then proceed sequentially in priority order (P1 → P2)
  - US1 and US2 are both P1 but US1 should complete first (visual foundation)
- **Additional Functions (Phase 8)**: Depends on US2 completion (bubbletea models)
- **Error Handling (Phase 9)**: Depends on US2 completion (models exist)
- **Polish (Phase 10)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P1)**: Can start after Foundational (Phase 2) - Depends on US1 (visual indicators)
- **User Story 3 (P2)**: Can start after US2 completion - Depends on bubbletea models existing
- **User Story 4 (P2)**: Can start after US2 completion - Depends on SelectListModel existing
- **User Story 5 (P2)**: Can start after US2 completion - Depends on MultilineInputModel existing

### Within Each User Story

- Tests (TDD) MUST be written and FAIL before implementation
- Models before functions
- Core implementation before integration
- Story complete before moving to next priority

### Parallel Opportunities

- All Foundational tasks marked [P] can run in parallel (T003-T006)
- All test tasks for a user story marked [P] can run in parallel
- Model creation tasks marked [P] can run in parallel (T023-T024)
- Refactoring tasks for different prompt functions marked [P] can run in parallel (T032-T040, T046-T051, T061-T062, T070-T073)
- Error handling tasks marked [P] can run in parallel (T078-T085)
- Polish tasks marked [P] can run in parallel (T086-T097)

---

## Parallel Example: User Story 2

```bash
# Launch all tests for User Story 2 together:
Task: "Write tests for TextInputModel in internal/ui/text_input_test.go"
Task: "Write tests for YesNoChoiceModel in internal/ui/yes_no_choice_test.go"
Task: "Write tests for TextInputModel with validation in internal/ui/text_input_test.go"
Task: "Write tests for TextInputModel with default values in internal/ui/text_input_test.go"
Task: "Write tests for YesNoChoiceModel keyboard handling in internal/ui/yes_no_choice_test.go"

# Launch model creation together:
Task: "Create TextInputModel struct implementing tea.Model in internal/ui/text_input.go"
Task: "Create YesNoChoiceModel struct implementing tea.Model in internal/ui/yes_no_choice.go"

# Launch prompt function refactoring together (after models complete):
Task: "Refactor PromptScope to use TextInputModel in internal/ui/prompts.go"
Task: "Refactor PromptScopeWithDefault to use TextInputModel in internal/ui/prompts.go"
Task: "Refactor PromptSubject to use TextInputModel with validation in internal/ui/prompts.go"
Task: "Refactor PromptEmptyCommit to use YesNoChoiceModel in internal/ui/prompts.go"
Task: "Refactor PromptConfirm to use YesNoChoiceModel in internal/ui/prompts.go"
```

---

## Implementation Strategy

### MVP First (User Stories 1 & 2 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1 (Visual Design)
4. Complete Phase 4: User Story 2 (Bubbletea Migration)
5. **STOP and VALIDATE**: Test User Stories 1 & 2 independently
6. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Visual indicators working
3. Add User Story 2 → Test independently → All prompts use bubbletea (MVP!)
4. Add User Story 3 → Test independently → Inline rendering
5. Add User Story 4 → Test independently → Select list enhancements
6. Add User Story 5 → Test independently → Multiline enhancements
7. Add Error Handling → Test independently → Complete error scenarios
8. Polish → Final validation

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 (Visual Design)
   - Developer B: Prepares for User Story 2 (studies bubbletea patterns)
3. Once US1 is done:
   - Developer A: User Story 2 (Text Input models)
   - Developer B: User Story 2 (Yes/No Choice models)
   - Developer C: User Story 2 (Prompt function refactoring)
4. Once US2 is done:
   - Developer A: User Story 3 (Remove alt screen)
   - Developer B: User Story 4 (Select list display)
   - Developer C: User Story 5 (Multiline display)
5. All developers: Error handling and polish

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing (TDD)
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- All prompt function signatures must remain backward compatible
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence
- Total tasks: 97
- Tasks per story: US1 (9), US2 (25), US3 (12), US4 (10), US5 (11), Additional (4), Error Handling (8), Polish (12), Setup (2), Foundational (4)
