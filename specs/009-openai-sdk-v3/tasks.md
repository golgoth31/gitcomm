# Tasks: Upgrade OpenAI Provider to SDK v3

**Input**: Design documents from `/specs/009-openai-sdk-v3/`
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

**Purpose**: Verify existing code structure, add SDK v3 dependency, and identify integration points

- [X] T001 Verify existing OpenAI provider implementation in internal/ai/openai_provider.go
- [X] T002 Verify existing AIProvider interface in internal/ai/provider.go
- [X] T003 Verify existing AIProviderConfig structure in internal/model/config.go
- [X] T004 [P] Add OpenAI SDK v3 dependency (github.com/openai/openai-go/v3) to go.mod
- [X] T005 [P] Verify existing provider selection mechanism in internal/service/commit_service.go
- [X] T006 [P] Review existing OpenAI provider tests in internal/ai/openai_provider_test.go

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

### Tests for OpenAI Provider SDK v3 Integration (TDD - Write First)

- [X] T007 [P] Write unit test for NewOpenAIProvider with SDK v3 client initialization in internal/ai/openai_provider_test.go
- [X] T008 [P] Write unit test for GenerateCommitMessage with successful SDK v3 API response in internal/ai/openai_provider_test.go
- [X] T009 [P] Write unit test for GenerateCommitMessage error handling (SDK v3 errors mapped to existing types) in internal/ai/openai_provider_test.go
- [X] T010 [P] Write unit test for SDK v3 initialization failure handling in internal/ai/openai_provider_test.go
- [X] T011 [P] Write unit test for context cancellation with SDK v3 in internal/ai/openai_provider_test.go
- [X] T012 [P] Write unit test for unmappable SDK v3 error handling (generic wrapping) in internal/ai/openai_provider_test.go

**Checkpoint**: Foundation ready - tests written and failing, user story implementation can now begin

---

## Phase 3: User Story 1 - Upgrade OpenAI Provider to SDK v3 (Priority: P1) 🎯 MVP

**Goal**: Upgrade OpenAI provider from SDK v1 to SDK v3 while maintaining 100% backward compatibility with existing functionality, interfaces, and configuration.

**Independent Test**: Configure gitcomm with OpenAI API credentials, run gitcomm with OpenAI provider, and verify that commit messages are generated successfully using SDK v3. The behavior should be identical to the current SDK v1 implementation, with no breaking changes to the user experience.

### Implementation for User Story 1

- [X] T013 [US1] Update import paths from SDK v1 to SDK v3 in internal/ai/openai_provider.go
- [X] T014 [US1] Update OpenAI SDK client initialization to use SDK v3 in NewOpenAIProvider in internal/ai/openai_provider.go
- [X] T015 [US1] Verify SDK v3 ChatCompletionNewParams structure matches current usage in internal/ai/openai_provider.go
- [X] T016 [US1] Update chat completion API call to use SDK v3 API in GenerateCommitMessage in internal/ai/openai_provider.go
- [X] T017 [US1] Verify SDK v3 response structure and update content extraction if needed in internal/ai/openai_provider.go
- [X] T018 [US1] Update mapSDKError function to handle SDK v3 error types in internal/ai/openai_provider.go
- [X] T019 [US1] Verify context cancellation and timeout handling works with SDK v3 in internal/ai/openai_provider.go
- [X] T020 [US1] Verify prompt building logic remains unchanged in internal/ai/openai_provider.go
- [X] T021 [US1] Update unit tests to match SDK v3 API structure in internal/ai/openai_provider_test.go
- [X] T022 [US1] Verify all existing unit tests pass with SDK v3 implementation in internal/ai/openai_provider_test.go
- [X] T023 [US1] Update integration tests if needed for SDK v3 in test/integration/ai_commit_test.go
- [X] T024 [US1] Verify all existing integration tests pass with SDK v3 implementation in test/integration/ai_commit_test.go

**Checkpoint**: At this point, User Story 1 should be fully functional - OpenAI provider uses SDK v3 with 100% functional parity to SDK v1 implementation

---

## Phase 4: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect the entire feature

- [X] T025 [P] Verify backward compatibility (existing providers still work, existing configs still work)
- [X] T026 [P] Verify error handling behavior remains identical (same error types, same user-facing messages, same fallback behavior)
- [X] T027 [P] Verify SDK v3 automatic retries work correctly (if available)
- [X] T028 [P] Verify timeout handling works correctly with SDK v3
- [X] T029 Code cleanup and formatting (gofmt, goimports) in internal/ai/openai_provider.go
- [X] T030 Run golangci-lint to verify code quality
- [X] T031 [P] Update README.md with SDK v3 information if needed
- [X] T032 [P] Update CHANGELOG.md with SDK v3 upgrade entry
- [X] T033 [P] Run quickstart.md validation scenarios
- [X] T034 Verify no regressions in other providers (Anthropic, Mistral, local)
- [X] T035 Run all unit tests to verify no regressions
- [X] T036 Run all integration tests to verify no regressions

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Foundational phase completion
- **Polish (Phase 4)**: Depends on User Story 1 completion

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories

### Within User Story 1

- Tests (T007-T012) MUST be written and FAIL before implementation
- Update imports (T013) before API calls (T014-T018)
- Verify API compatibility (T015) before updating calls (T016)
- Update error handling (T018) after API calls are updated
- Update tests (T021-T024) after implementation is complete

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel (T004, T005, T006)
- All Foundational test tasks marked [P] can run in parallel (T007-T012)
- Polish tasks marked [P] can run in parallel (T025-T028, T031-T033)

---

## Parallel Example: User Story 1

```bash
# Launch all foundational tests together:
Task: "Write unit test for NewOpenAIProvider with SDK v3 client initialization"
Task: "Write unit test for GenerateCommitMessage with successful SDK v3 API response"
Task: "Write unit test for GenerateCommitMessage error handling"
Task: "Write unit test for SDK v3 initialization failure handling"
Task: "Write unit test for context cancellation with SDK v3"
Task: "Write unit test for unmappable SDK v3 error handling"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Test User Story 1 independently
5. Complete Phase 4: Polish
6. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Deploy/Demo (MVP!)
3. Each story adds value without breaking previous stories

### TDD Approach

1. Write tests first (Phase 2) - ensure they fail
2. Implement SDK v3 upgrade (Phase 3)
3. Verify tests pass
4. Update tests if SDK v3 API structure changed
5. Verify all existing tests still pass

---

## Notes

- [P] tasks = different files, no dependencies
- [US1] label maps task to User Story 1 for traceability
- User Story 1 should be independently completable and testable
- Verify tests fail before implementing
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Maintain 100% backward compatibility throughout
- No changes to other providers (Anthropic, Mistral, local)
