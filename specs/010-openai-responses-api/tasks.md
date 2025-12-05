# Tasks: Migrate OpenAI Provider to Responses API

**Input**: Design documents from `/specs/010-openai-responses-api/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Tests are included following TDD approach as required by the constitution.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1)
- Include exact file paths in descriptions

## Path Conventions

- **CLI project**: `cmd/gitcomm/`, `internal/` at repository root
- Paths shown below follow existing project structure

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Verify existing code structure, check SDK v3 Responses API support, and identify integration points

- [X] T001 Verify existing OpenAI provider implementation in internal/ai/openai_provider.go
- [X] T002 Verify existing AIProvider interface in internal/ai/provider.go
- [X] T003 Verify existing AIProviderConfig structure in internal/model/config.go
- [X] T004 [P] Check OpenAI SDK v3 for Responses API support (github.com/openai/openai-go/v3)
- [X] T005 [P] Review existing OpenAI provider tests in internal/ai/openai_provider_test.go
- [X] T006 [P] Review Responses API migration guide and documentation

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

### Tests for OpenAI Provider Responses API Integration (TDD - Write First)

- [X] T007 [P] Write unit test for NewOpenAIProvider with Responses API client initialization in internal/ai/openai_provider_test.go
- [X] T008 [P] Write unit test for GenerateCommitMessage with successful Responses API response in internal/ai/openai_provider_test.go
- [X] T009 [P] Write unit test for GenerateCommitMessage error handling (Responses API errors mapped to existing types) in internal/ai/openai_provider_test.go
- [X] T010 [P] Write unit test for Responses API initialization failure handling in internal/ai/openai_provider_test.go
- [X] T011 [P] Write unit test for context cancellation with Responses API in internal/ai/openai_provider_test.go
- [X] T012 [P] Write unit test for empty response handling from Responses API in internal/ai/openai_provider_test.go
- [X] T013 [P] Write unit test for stateless mode configuration (store: false) in internal/ai/openai_provider_test.go

**Checkpoint**: Foundation ready - tests written and failing, user story implementation can now begin

---

## Phase 3: User Story 1 - Migrate to Responses API (Priority: P1) 🎯 MVP

**Goal**: Migrate OpenAI provider from Chat Completions API to Responses API while maintaining 100% backward compatibility with existing functionality, interfaces, and configuration.

**Independent Test**: Configure gitcomm with OpenAI API credentials, run gitcomm with OpenAI provider, and verify that commit messages are generated successfully using the Responses API. The behavior should be identical to the current Chat Completions implementation, with no breaking changes to the user experience.

### Implementation for User Story 1

- [X] T014 [US1] Update import paths if needed for Responses API in internal/ai/openai_provider.go
- [X] T015 [US1] Check if SDK v3 supports Responses API and determine implementation approach (SDK vs custom HTTP client) in internal/ai/openai_provider.go
- [X] T016 [US1] Update API endpoint from `/v1/chat/completions` to `/v1/responses` in internal/ai/openai_provider.go
- [X] T017 [US1] Convert `messages` array to `input` parameter format in GenerateCommitMessage in internal/ai/openai_provider.go
- [X] T018 [US1] Add stateless mode configuration (store: false) to Responses API request in internal/ai/openai_provider.go
- [X] T019 [US1] Update response extraction logic to handle Responses API response structure (content/text field) in internal/ai/openai_provider.go
- [X] T020 [US1] Update mapSDKError function to handle Responses API error types in internal/ai/openai_provider.go
- [X] T021 [US1] Verify context cancellation and timeout handling works with Responses API in internal/ai/openai_provider.go
- [X] T022 [US1] Verify prompt building logic remains unchanged in internal/ai/openai_provider.go
- [X] T023 [US1] Update unit tests to match Responses API structure in internal/ai/openai_provider_test.go
- [X] T024 [US1] Verify all existing unit tests pass with Responses API implementation in internal/ai/openai_provider_test.go
- [X] T025 [US1] Update integration tests if needed for Responses API in test/integration/ai_commit_test.go
- [X] T026 [US1] Verify all existing integration tests pass with Responses API implementation in test/integration/ai_commit_test.go

**Checkpoint**: At this point, User Story 1 should be fully functional - OpenAI provider uses Responses API with 100% functional parity to Chat Completions implementation

---

## Phase 4: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect the entire feature

- [X] T027 [P] Verify backward compatibility (existing providers still work, existing configs still work)
- [X] T028 [P] Verify error handling behavior remains identical (same error types, same user-facing messages, same fallback behavior)
- [X] T029 [P] Verify stateless mode is correctly configured (no state persistence between calls)
- [X] T030 [P] Verify timeout handling works correctly with Responses API
- [X] T031 Code cleanup and formatting (gofmt, goimports) in internal/ai/openai_provider.go
- [X] T032 Run golangci-lint to verify code quality
- [X] T033 [P] Update README.md with Responses API information if needed
- [X] T034 [P] Update CHANGELOG.md with Responses API migration entry
- [X] T035 [P] Run quickstart.md validation scenarios
- [X] T036 Verify no regressions in other providers (Anthropic, Mistral, local)
- [X] T037 Run all unit tests to verify no regressions
- [X] T038 Run all integration tests to verify no regressions

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

- Tests (T007-T013) MUST be written and FAIL before implementation
- Check SDK support (T015) before API calls (T016-T021)
- Update endpoint (T016) before request conversion (T017)
- Convert request format (T017) before response extraction (T019)
- Update error handling (T020) after API calls are updated
- Update tests (T023-T026) after implementation is complete

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel (T004, T005, T006)
- All Foundational test tasks marked [P] can run in parallel (T007-T013)
- Polish tasks marked [P] can run in parallel (T027-T030, T033-T035)

---

## Parallel Example: User Story 1

```bash
# Launch all foundational tests together:
Task: "Write unit test for NewOpenAIProvider with Responses API client initialization"
Task: "Write unit test for GenerateCommitMessage with successful Responses API response"
Task: "Write unit test for GenerateCommitMessage error handling"
Task: "Write unit test for Responses API initialization failure handling"
Task: "Write unit test for context cancellation with Responses API"
Task: "Write unit test for empty response handling from Responses API"
Task: "Write unit test for stateless mode configuration"
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
2. Implement Responses API migration (Phase 3)
3. Verify tests pass
4. Update tests if Responses API structure differs from expectations
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
- SDK v3 support check determines implementation approach (SDK vs custom HTTP client)
