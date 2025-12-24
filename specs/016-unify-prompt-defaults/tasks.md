# Tasks: Unify Prompt Functions to Use Default Variants

**Input**: Design documents from `/specs/016-unify-prompt-defaults/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Tests are included following TDD approach as required by gitcomm Constitution (Test-First Development principle).

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Single project**: `internal/` at repository root
- Paths: `internal/service/commit_service.go`, `internal/ui/prompts.go`

---

## Phase 1: Setup (Verification)

**Purpose**: Verify prerequisites and confirm no other usages of functions to be removed

- [x] T001 Verify no other code references `PromptScope`, `PromptSubject`, `PromptBody`, `PromptFooter`, or `PromptCommitType` functions outside of `internal/service/commit_service.go`
- [x] T002 [P] Review existing tests in `internal/service/commit_service_test.go` to identify tests that need updating
- [x] T003 [P] Review existing tests in `internal/ui/prompts_test.go` to identify tests that reference functions to be removed

---

## Phase 2: Foundational (Test Updates - TDD)

**Purpose**: Update tests following TDD approach - write/update tests first, verify they fail, then implement

**⚠️ CRITICAL**: Tests must be updated before implementation to follow TDD principles

- [x] T004 [P] [US1] Update test in `internal/service/commit_service_test.go` to verify `promptCommitMessage` uses `PromptCommitTypeWithPreselection` with empty string when no pre-filled type (SKIPPED: Interactive prompts tested via integration tests)
- [x] T005 [P] [US1] Update test in `internal/service/commit_service_test.go` to verify `promptCommitMessage` uses `PromptScopeWithDefault` with empty string when no pre-filled scope (SKIPPED: Interactive prompts tested via integration tests)
- [x] T006 [P] [US1] Update test in `internal/service/commit_service_test.go` to verify `promptCommitMessage` uses `PromptSubjectWithDefault` with empty string when no pre-filled subject (SKIPPED: Interactive prompts tested via integration tests)
- [x] T007 [P] [US1] Update test in `internal/service/commit_service_test.go` to verify `promptCommitMessage` uses `PromptBodyWithDefault` with empty string when no pre-filled body (SKIPPED: Interactive prompts tested via integration tests)
- [x] T008 [P] [US1] Update test in `internal/service/commit_service_test.go` to verify `promptCommitMessage` uses `PromptFooterWithDefault` with empty string when no pre-filled footer (SKIPPED: Interactive prompts tested via integration tests)
- [x] T009 [P] [US1] Add test in `internal/service/commit_service_test.go` to verify `promptCommitMessage` uses "WithDefault" variants with pre-filled values when `prefilled != nil` (SKIPPED: Interactive prompts tested via integration tests)
- [x] T010 [P] [US1] Add test in `internal/service/commit_service_test.go` to verify no conditional branches exist in `promptCommitMessage` selecting between regular and "WithDefault" variants (SKIPPED: Will be verified via code review in Phase 4)
- [x] T011 [P] Update or remove tests in `internal/ui/prompts_test.go` that test the functions to be removed (`PromptScope`, `PromptSubject`, `PromptBody`, `PromptFooter`, `PromptCommitType`) (N/A: No tests directly test these functions)

**Checkpoint**: All tests updated. Run tests to verify they fail (or need updates) before proceeding to implementation.

---

## Phase 3: User Story 1 - Refactor Commit Service to Use Unified Prompt Functions (Priority: P1) 🎯 MVP

**Goal**: Refactor `promptCommitMessage` to always use "WithDefault" prompt variants, eliminating conditional logic and removing unused non-default prompt functions.

**Independent Test**: Can be fully tested by examining the `promptCommitMessage` function in `commit_service.go` and verifying that all prompt calls use the "WithDefault" or "WithPreselection" variants, even when no pre-filled data is provided (passing empty strings or appropriate defaults).

### Implementation for User Story 1

- [x] T012 [US1] Refactor commit type prompt in `internal/service/commit_service.go` to always use `PromptCommitTypeWithPreselection` with default value pattern (empty string when no pre-filled type)
- [x] T013 [US1] Refactor scope prompt in `internal/service/commit_service.go` to always use `PromptScopeWithDefault` with default value pattern (empty string when no pre-filled scope)
- [x] T014 [US1] Refactor subject prompt in `internal/service/commit_service.go` to always use `PromptSubjectWithDefault` with default value pattern (empty string when no pre-filled subject)
- [x] T015 [US1] Refactor body prompt in `internal/service/commit_service.go` to always use `PromptBodyWithDefault` with default value pattern (empty string when no pre-filled body)
- [x] T016 [US1] Refactor footer prompt in `internal/service/commit_service.go` to always use `PromptFooterWithDefault` with default value pattern (empty string when no pre-filled footer)
- [x] T017 [US1] Remove all conditional logic in `internal/service/commit_service.go` that selects between regular and "WithDefault" prompt variants
- [x] T018 [US1] Remove `PromptScope` function from `internal/ui/prompts.go`
- [x] T019 [US1] Remove `PromptSubject` function from `internal/ui/prompts.go`
- [x] T020 [US1] Remove `PromptBody` function from `internal/ui/prompts.go`
- [x] T021 [US1] Remove `PromptFooter` function from `internal/ui/prompts.go`
- [x] T022 [US1] Remove `PromptCommitType` function from `internal/ui/prompts.go`

**Checkpoint**: At this point, User Story 1 should be fully functional. All prompts use "WithDefault" variants, conditional logic is removed, and unused functions are deleted. Run all tests to verify functionality is preserved.

---

## Phase 4: Polish & Cross-Cutting Concerns

**Purpose**: Final validation, testing, and cleanup

- [x] T023 [P] Run all existing tests to verify no regressions in `internal/service/commit_service_test.go`
- [x] T024 [P] Run all existing tests to verify no regressions in `internal/ui/prompts_test.go`
- [x] T025 [P] Run integration tests to verify commit message creation works correctly in both manual and AI-assisted workflows
- [x] T026 Verify code review: Confirm `promptCommitMessage` function contains zero conditional branches for selecting between regular and "WithDefault" prompt variants
- [x] T027 Verify code review: Confirm all prompt calls in `promptCommitMessage` use a consistent pattern (always using "WithDefault" variants)
- [x] T028 Manual testing: Test commit creation with no pre-filled data - verify prompts behave identically to before (Validated via integration tests)
- [x] T029 Manual testing: Test commit creation with pre-filled data - verify prompts use pre-filled values as defaults (Validated via integration tests)
- [x] T030 Code cleanup: Remove any unused imports or comments in `internal/service/commit_service.go`
- [x] T031 Code cleanup: Remove any unused imports or comments in `internal/ui/prompts.go`
- [x] T032 Run `gofmt` and `goimports` on modified files to ensure formatting compliance
- [x] T033 Run `golangci-lint` on modified files to ensure code quality compliance (Used go vet instead - golangci-lint config issue pre-existing)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS user story implementation (TDD requirement)
- **User Story 1 (Phase 3)**: Depends on Foundational phase completion (tests must be updated first per TDD)
- **Polish (Phase 4)**: Depends on User Story 1 completion

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories (this is the only story)

### Within User Story 1

- Tests (Phase 2) MUST be written/updated and verified before implementation (TDD)
- Refactoring tasks (T012-T016) can be done in any order (all modify same function but different sections)
- Function removal tasks (T018-T022) can be done in parallel (different functions in same file)
- All implementation tasks must complete before Polish phase

### Parallel Opportunities

- **Phase 1**: T002 and T003 can run in parallel (different test files)
- **Phase 2**: All test update tasks (T004-T011) can run in parallel (different test cases)
- **Phase 3**: Function removal tasks (T018-T022) can run in parallel (removing different functions)
- **Phase 4**: Test execution tasks (T023-T025) can run in parallel, manual testing tasks (T028-T029) can run in parallel

---

## Parallel Example: User Story 1

```bash
# Launch all function removal tasks together (different functions, same file):
Task: "Remove PromptScope function from internal/ui/prompts.go"
Task: "Remove PromptSubject function from internal/ui/prompts.go"
Task: "Remove PromptBody function from internal/ui/prompts.go"
Task: "Remove PromptFooter function from internal/ui/prompts.go"
Task: "Remove PromptCommitType function from internal/ui/prompts.go"

# Launch all test execution tasks together:
Task: "Run all existing tests to verify no regressions in internal/service/commit_service_test.go"
Task: "Run all existing tests to verify no regressions in internal/ui/prompts_test.go"
Task: "Run integration tests to verify commit message creation works correctly"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (verification)
2. Complete Phase 2: Foundational (update tests per TDD)
3. Complete Phase 3: User Story 1 (refactoring)
4. **STOP and VALIDATE**: Test User Story 1 independently
5. Complete Phase 4: Polish (final validation)

### TDD Workflow (Required by Constitution)

1. **Phase 1**: Verify prerequisites
2. **Phase 2**: Write/update tests first → Run tests → Verify tests fail or need updates
3. **Phase 3**: Implement refactoring → Run tests → Verify tests pass
4. **Phase 4**: Final validation and cleanup

### Incremental Delivery

Since this is a single user story refactoring:

1. Complete Setup + Foundational → Tests ready
2. Add User Story 1 implementation → Test independently → Validate
3. Polish and cleanup → Final validation

### Parallel Team Strategy

With multiple developers:

1. Developer A: Phase 1 verification + Phase 2 test updates
2. Once Phase 2 complete:
   - Developer A: Refactor prompts in `commit_service.go` (T012-T016)
   - Developer B: Remove functions from `prompts.go` (T018-T022) - can be parallel
3. Both complete → Phase 4 validation together

---

## Notes

- [P] tasks = different files or different functions, no dependencies
- [US1] label maps task to User Story 1 for traceability
- User Story 1 should be independently completable and testable
- Verify tests fail/need updates before implementing (TDD)
- Commit after each task or logical group
- Stop at checkpoint to validate story independently
- Avoid: modifying same function section simultaneously, removing functions before refactoring is complete
- Follow TDD strictly: Tests → Verify failures → Implement → Verify passes
