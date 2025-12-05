# Tasks: Rewrite CLI Prompts with Huh Library

**Feature**: 013-huh-prompts
**Branch**: `013-huh-prompts`
**Date**: 2025-01-27

## Overview

Migrate all CLI prompts from custom Bubble Tea implementations to the `huh` library. All prompts must render inline (no alt screen) and display post-validation summary lines with green checkmark.

## Implementation Strategy

**MVP Scope**: Complete User Story 1 (Migrate All Prompts to Huh Library) - this delivers the core functionality and inherently covers User Stories 2 and 3.

**Incremental Delivery**:
1. Setup and foundational helpers
2. Migrate text input prompts (scope, subject)
3. Migrate multiline prompts (body, footer)
4. Migrate selection prompts (commit type)
5. Migrate confirmation prompts (yes/no choices)
6. Migrate multi-choice prompts (AI acceptance, commit failure)
7. Remove old models and polish

## Dependencies

### User Story Completion Order

All three user stories (US1, US2, US3) are P1 priority and will be completed together:
- **US1** (Migrate All Prompts): Core implementation - migrating each prompt function
- **US2** (Inline Rendering): Verified during US1 implementation (huh renders inline by default)
- **US3** (Post-Validation Display): Implemented as part of each prompt migration in US1

**Dependency Graph**:
```
Phase 1 (Setup) → Phase 2 (Foundational) → Phase 3 (US1+US2+US3) → Phase 4 (Polish)
```

### Parallel Execution Opportunities

Within Phase 3, prompt migrations can be done in parallel groups:
- **Group A**: Text input prompts (T010-T013) - can be done in parallel
- **Group B**: Multiline prompts (T014-T017) - can be done in parallel
- **Group C**: Selection prompts (T018-T019) - can be done in parallel
- **Group D**: Confirmation prompts (T020-T024) - can be done in parallel
- **Group E**: Multi-choice prompts (T025-T026) - can be done in parallel

## Independent Test Criteria

### User Story 1 - Migrate All Prompts to Huh Library
- **Test**: Run gitcomm and interact with all prompt types
- **Verify**: All prompts use `huh` library (check imports, no custom Bubble Tea models)
- **Verify**: All prompt types work: text input, multiline, select, confirm, multi-choice

### User Story 2 - Inline Rendering Without Alt Screen
- **Test**: Run gitcomm and interact with prompts
- **Verify**: Terminal history remains visible (scroll up to see previous output)
- **Verify**: No alt screen mode is used (check terminal behavior)

### User Story 3 - Post-Validation Display Format
- **Test**: Complete any prompt and verify output
- **Verify**: Summary line appears: `✓ <title>: <value>` with green checkmark
- **Verify**: Prompt UI is cleared before summary line appears

---

## Phase 1: Setup

### Story Goal
Initialize project with `huh` dependency and verify setup.

### Tasks

- [X] T001 Add `huh` library dependency to `go.mod` using `go get github.com/charmbracelet/huh@latest`
- [X] T002 Verify `huh` dependency is added correctly in `go.mod` and `go.sum`
- [X] T003 Run `go mod tidy` to ensure dependency resolution is correct

---

## Phase 2: Foundational

### Story Goal
Create helper functions and utilities needed for prompt migration.

### Tasks

- [X] T004 Create helper function `formatPostValidationSummary(title string, value interface{}) string` in `internal/ui/display.go` that returns formatted string: `✓ <title>: <value>`
- [X] T005 Add unit tests for `formatPostValidationSummary` in `internal/ui/display_test.go` covering text, multiline truncation, and special characters
- [X] T006 Create helper function `printPostValidationSummary(title string, value interface{})` in `internal/ui/display.go` that prints the formatted summary line
- [X] T007 Verify inline rendering capability by creating a test `huh.Form` and confirming it renders inline (no alt screen) - Verified: `huh` forms render inline by default

---

## Phase 3: User Story 1 - Migrate All Prompts to Huh Library

### Story Goal
Migrate all prompt functions from custom Bubble Tea models to `huh` library, ensuring inline rendering (US2) and post-validation display (US3).

### Independent Test
Run gitcomm and interact with all prompt types. Verify all prompts use `huh` library and function correctly.

### Text Input Prompts

- [X] T008 [P] [US1] Migrate `PromptScope` function to use `huh.NewInput()` in `internal/ui/prompts.go`, maintain signature, add post-validation summary display
- [X] T009 [P] [US1] Migrate `PromptScopeWithDefault` function to use `huh.NewInput()` with default value in `internal/ui/prompts.go`, maintain signature, add post-validation summary display
- [X] T010 [P] [US1] Migrate `PromptSubject` function to use `huh.NewInput()` with validation in `internal/ui/prompts.go`, maintain signature, add post-validation summary display
- [X] T011 [P] [US1] Migrate `PromptSubjectWithDefault` function to use `huh.NewInput()` with default and validation in `internal/ui/prompts.go`, maintain signature, add post-validation summary display

### Multiline Input Prompts

- [X] T012 [P] [US1] Migrate `PromptBody` function to use `huh.NewText()` in `internal/ui/prompts.go`, maintain signature, add post-validation summary display with truncation for long content
- [X] T013 [P] [US1] Migrate `PromptBodyWithDefault` function to use `huh.NewText()` with default value in `internal/ui/prompts.go`, maintain signature, add post-validation summary display
- [X] T014 [P] [US1] Migrate `PromptFooter` function to use `huh.NewText()` in `internal/ui/prompts.go`, maintain signature, add post-validation summary display
- [X] T015 [P] [US1] Migrate `PromptFooterWithDefault` function to use `huh.NewText()` with default value in `internal/ui/prompts.go`, maintain signature, add post-validation summary display

### Selection Prompts

- [X] T016 [P] [US1] Migrate `PromptCommitType` function to use `huh.NewSelect()` with commit type options in `internal/ui/prompts.go`, maintain signature, add post-validation summary display
- [X] T017 [P] [US1] Migrate `PromptCommitTypeWithPreselection` function to use `huh.NewSelect()` with pre-selected value in `internal/ui/prompts.go`, maintain signature, add post-validation summary display

### Confirmation Prompts

- [X] T018 [P] [US1] Migrate `PromptEmptyCommit` function to use `huh.NewConfirm()` in `internal/ui/prompts.go`, maintain signature, add post-validation summary display
- [X] T019 [P] [US1] Migrate `PromptConfirm` function to use `huh.NewConfirm()` with custom message in `internal/ui/prompts.go`, maintain signature, add post-validation summary display
- [X] T020 [P] [US1] Migrate `PromptAIUsage` function to use `huh.NewConfirm()` with default true in `internal/ui/prompts.go`, maintain signature, add post-validation summary display
- [X] T021 [P] [US1] Migrate `PromptAIMessageEdit` function to use `huh.NewConfirm()` with default true and error messages in `internal/ui/prompts.go`, maintain signature, add post-validation summary display
- [X] T022 [P] [US1] Migrate `PromptRejectChoice` function to use `huh.NewConfirm()` with default true in `internal/ui/prompts.go`, maintain signature, add post-validation summary display

### Multi-Choice Prompts

- [X] T023 [P] [US1] Migrate `PromptAIMessageAcceptanceOptions` function to use `huh.NewSelect()` with three options (AcceptAndCommit, AcceptAndEdit, Reject) in `internal/ui/prompts.go`, maintain signature, add post-validation summary display
- [X] T024 [P] [US1] Migrate `PromptCommitFailureChoice` function to use `huh.NewSelect()` with three options (RetryCommit, EditMessage, CancelCommit) in `internal/ui/prompts.go`, maintain signature, add post-validation summary display

### Testing

- [X] T025 [US1] Update unit tests in `internal/ui/prompts_test.go` for all migrated prompt functions, verify `huh` library usage and post-validation display
- [X] T026 [US1] Add integration tests in `test/integration/` to verify all prompt types work end-to-end with `huh` library

---

## Phase 4: User Story 2 - Inline Rendering Verification

### Story Goal
Verify that all prompts render inline without alt screen mode.

### Independent Test
Run gitcomm and interact with prompts. Verify terminal history remains visible and no alt screen is used.

### Tasks

- [X] T027 [US2] Add integration test in `test/integration/inline_rendering_test.go` to verify prompts render inline (terminal history visible)
- [X] T028 [US2] Verify all prompt functions use `huh` forms that render inline by default (no alt screen configuration needed)
- [X] T029 [US2] Test with narrow terminal width to ensure inline rendering handles constraints gracefully

---

## Phase 5: User Story 3 - Post-Validation Display Verification

### Story Goal
Verify that all prompts display post-validation summary lines correctly.

### Independent Test
Complete any prompt and verify output shows `✓ <title>: <value>` with green checkmark.

### Tasks

- [X] T030 [US3] Add integration test in `test/integration/post_validation_display_test.go` to verify summary lines appear for all prompt types
- [X] T031 [US3] Verify post-validation summary format is consistent across all prompt types: `✓ <title>: <value>`
- [X] T032 [US3] Test multiline content truncation in post-validation display (body, footer prompts)
- [X] T033 [US3] Verify prompt UI is cleared before summary line appears for all prompt types

---

## Phase 6: Polish & Cross-Cutting Concerns

### Story Goal
Remove old implementations, update documentation, and ensure all tests pass.

### Tasks

- [X] T034 Remove `TextInputModel` implementation from `internal/ui/text_input.go`
- [X] T035 Remove `TextInputModel` tests from `internal/ui/text_input_test.go`
- [X] T036 Remove `MultilineInputModel` implementation from `internal/ui/multiline_input.go`
- [X] T037 Remove `MultilineInputModel` tests from `internal/ui/multiline_input_test.go`
- [X] T038 Remove `YesNoChoiceModel` implementation from `internal/ui/yes_no_choice.go`
- [X] T039 Remove `YesNoChoiceModel` tests from `internal/ui/yes_no_choice_test.go`
- [X] T040 Remove `SelectListModel` implementation from `internal/ui/select_list.go`
- [X] T041 Remove `SelectListModel` tests from `internal/ui/select_list_test.go`
- [X] T042 Review `internal/ui/prompt_state.go` and remove if no longer needed (check for remaining usages) - Kept for GetVisualIndicator which may be useful for future use
- [X] T043 Update all integration tests in `test/integration/` to reflect `huh` library behavior
- [X] T044 Run full test suite and verify all tests pass: `go test ./...`
- [X] T045 Verify backward compatibility: all existing callers in `internal/service/commit_service.go` work without changes
- [X] T046 Update CHANGELOG.md with migration details
- [X] T047 Verify no references to old Bubble Tea models remain in codebase (grep for TextInputModel, MultilineInputModel, etc.)

---

## Summary

**Total Tasks**: 47

**Tasks by Phase**:
- Phase 1 (Setup): 3 tasks
- Phase 2 (Foundational): 4 tasks
- Phase 3 (US1 - Migrate Prompts): 19 tasks (17 migrations + 2 testing)
- Phase 4 (US2 - Inline Rendering): 3 tasks
- Phase 5 (US3 - Post-Validation Display): 4 tasks
- Phase 6 (Polish): 14 tasks

**Parallel Opportunities**:
- 17 prompt migration tasks (T008-T024) can be done in parallel groups
- 4 foundational helper tasks (T004-T007) can be done in parallel
- 14 polish/cleanup tasks (T034-T047) can be done in parallel after migrations complete

**MVP Scope**: Phases 1-3 (Setup, Foundational, Migrate All Prompts) deliver complete functionality covering all three user stories.

**Suggested Implementation Order**:
1. Setup (Phase 1) - sequential
2. Foundational (Phase 2) - can parallelize T004-T007
3. Migrate Prompts (Phase 3) - can parallelize by prompt type groups
4. Verification (Phases 4-5) - sequential verification
5. Polish (Phase 6) - can parallelize cleanup tasks
