# Tasks: Use Official SDKs for AI Providers

**Input**: Design documents from `/specs/008-official-sdk-integration/`
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

**Purpose**: Verify existing code structure, add SDK dependencies, and identify integration points

- [X] T001 Verify existing AIProvider interface in internal/ai/provider.go
- [X] T002 Verify existing provider implementations (OpenAI, Anthropic, Mistral) in internal/ai/
- [X] T003 [P] Add OpenAI SDK dependency (github.com/openai/openai-go) to go.mod
- [X] T004 [P] Add Anthropic SDK dependency (github.com/anthropics/anthropic-sdk-go) to go.mod
- [X] T005 [P] Add Mistral SDK dependency (github.com/Gage-Technologies/mistral-go) to go.mod
- [X] T006 [P] Verify existing provider selection mechanism in internal/service/commit_service.go
- [X] T007 [P] Verify existing AIProviderConfig structure in internal/model/config.go

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

### Tests for OpenAI Provider SDK Integration (TDD - Write First)

- [X] T008 [P] Write unit test for NewOpenAIProvider with SDK client initialization in internal/ai/openai_provider_test.go
- [X] T009 [P] Write unit test for GenerateCommitMessage with successful SDK API response in internal/ai/openai_provider_test.go
- [X] T010 [P] Write unit test for GenerateCommitMessage error handling (SDK errors mapped to existing types) in internal/ai/openai_provider_test.go
- [X] T011 [P] Write unit test for SDK initialization failure handling in internal/ai/openai_provider_test.go

### Tests for Anthropic Provider SDK Integration (TDD - Write First)

- [X] T012 [P] Write unit test for NewAnthropicProvider with SDK client initialization in internal/ai/anthropic_provider_test.go
- [X] T013 [P] Write unit test for GenerateCommitMessage with successful SDK API response in internal/ai/anthropic_provider_test.go
- [X] T014 [P] Write unit test for GenerateCommitMessage error handling (SDK errors mapped to existing types) in internal/ai/anthropic_provider_test.go
- [X] T015 [P] Write unit test for SDK initialization failure handling in internal/ai/anthropic_provider_test.go

### Tests for Mistral Provider SDK Integration (TDD - Write First)

- [X] T016 [P] Write unit test for NewMistralProvider with SDK client initialization in internal/ai/mistral_provider_test.go
- [X] T017 [P] Write unit test for GenerateCommitMessage with successful SDK API response in internal/ai/mistral_provider_test.go
- [X] T018 [P] Write unit test for GenerateCommitMessage error handling (SDK errors mapped to existing types) in internal/ai/mistral_provider_test.go
- [X] T019 [P] Write unit test for SDK initialization failure handling in internal/ai/mistral_provider_test.go

**Checkpoint**: Foundation ready - tests written and failing, user story implementation can now begin

---

## Phase 3: User Story 1 - Replace OpenAI HTTP Client with Official SDK (Priority: P1) 🎯 MVP

**Goal**: Replace OpenAI HTTP client implementation with official OpenAI Go SDK, maintaining 100% backward compatibility with existing functionality, interfaces, and configuration.

**Independent Test**: Configure OpenAI API credentials, run gitcomm with OpenAI provider, and verify that commit messages are generated successfully using the official OpenAI SDK. The behavior should be identical to the current HTTP client implementation.

### Implementation for User Story 1

- [X] T020 [US1] Replace HTTP client with OpenAI SDK client initialization in NewOpenAIProvider in internal/ai/openai_provider.go
- [X] T021 [US1] Replace HTTP request/response handling with OpenAI SDK API calls in GenerateCommitMessage in internal/ai/openai_provider.go
- [X] T022 [US1] Implement SDK error mapping to existing error types (preserve user-facing messages) in internal/ai/openai_provider.go
- [X] T023 [US1] Add SDK initialization failure handling (fail fast with clear error) in internal/ai/openai_provider.go
- [X] T024 [US1] Ensure context cancellation and timeout handling works with SDK in internal/ai/openai_provider.go
- [X] T025 [US1] Verify prompt building logic remains unchanged in internal/ai/openai_provider.go
- [X] T026 [US1] Verify response parsing logic extracts commit message correctly from SDK response in internal/ai/openai_provider.go

**Checkpoint**: At this point, User Story 1 should be fully functional - OpenAI provider uses official SDK with 100% functional parity to HTTP client implementation

---

## Phase 4: User Story 2 - Replace Anthropic HTTP Client with Official SDK (Priority: P2)

**Goal**: Replace Anthropic HTTP client implementation with official Anthropic Go SDK, maintaining 100% backward compatibility with existing functionality, interfaces, and configuration.

**Independent Test**: Configure Anthropic API credentials, run gitcomm with Anthropic provider, and verify that commit messages are generated successfully using the official Anthropic SDK. The behavior should be identical to the current HTTP client implementation.

### Implementation for User Story 2

- [X] T027 [US2] Replace HTTP client with Anthropic SDK client initialization in NewAnthropicProvider in internal/ai/anthropic_provider.go
- [X] T028 [US2] Replace HTTP request/response handling with Anthropic SDK API calls in GenerateCommitMessage in internal/ai/anthropic_provider.go
- [X] T029 [US2] Implement SDK error mapping to existing error types (preserve user-facing messages) in internal/ai/anthropic_provider.go
- [X] T030 [US2] Add SDK initialization failure handling (fail fast with clear error) in internal/ai/anthropic_provider.go
- [X] T031 [US2] Ensure context cancellation and timeout handling works with SDK in internal/ai/anthropic_provider.go
- [X] T032 [US2] Verify prompt building logic remains unchanged in internal/ai/anthropic_provider.go
- [X] T033 [US2] Verify response parsing logic extracts commit message correctly from SDK response in internal/ai/anthropic_provider.go

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently - both providers use official SDKs with 100% functional parity

---

## Phase 5: User Story 3 - Replace Mistral HTTP Client with Official SDK (Priority: P2)

**Goal**: Replace Mistral HTTP client implementation with official Mistral Go SDK, maintaining 100% backward compatibility with existing functionality, interfaces, and configuration.

**Independent Test**: Configure Mistral API credentials, run gitcomm with Mistral provider, and verify that commit messages are generated successfully using the official Mistral SDK. The behavior should be identical to the current HTTP client implementation.

### Implementation for User Story 3

- [X] T034 [US3] Replace HTTP client with Mistral SDK client initialization in NewMistralProvider in internal/ai/mistral_provider.go
- [X] T035 [US3] Replace HTTP request/response handling with Mistral SDK API calls in GenerateCommitMessage in internal/ai/mistral_provider.go
- [X] T036 [US3] Implement SDK error mapping to existing error types (preserve user-facing messages) in internal/ai/mistral_provider.go
- [X] T037 [US3] Add SDK initialization failure handling (fail fast with clear error) in internal/ai/mistral_provider.go
- [X] T038 [US3] Ensure context cancellation and timeout handling works with SDK in internal/ai/mistral_provider.go
- [X] T039 [US3] Verify prompt building logic remains unchanged in internal/ai/mistral_provider.go
- [X] T040 [US3] Verify response parsing logic extracts commit message correctly from SDK response in internal/ai/mistral_provider.go

**Checkpoint**: At this point, all three user stories should be fully functional - all providers use official SDKs with 100% functional parity

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect the entire feature

- [X] T041 [P] Extend AIProviderConfig with optional SDK-specific fields if needed in internal/model/config.go
- [X] T042 [P] Update config.yaml.example with SDK-specific configuration examples if needed in configs/config.yaml.example
- [X] T043 Run all unit tests to verify no regressions (existing tests must pass)
- [X] T044 Run all integration tests to verify no regressions
- [X] T045 [P] Verify backward compatibility (existing providers still work, existing configs still work)
- [X] T046 [P] Verify error handling behavior remains identical (same error messages, same fallback behavior)
- [X] T047 [P] Verify SDK automatic retries work correctly (if available)
- [X] T048 [P] Verify timeout handling works correctly with SDKs
- [X] T049 Code cleanup and formatting (gofmt, goimports)
- [X] T050 Run golangci-lint to verify code quality
- [X] T051 [P] Update README.md with SDK integration information
- [X] T052 [P] Update CHANGELOG.md with feature entry
- [X] T053 [P] Run quickstart.md validation scenarios

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can proceed sequentially (P1 → P2 → P2) or in parallel if team capacity allows
  - Each story is independently testable
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - No dependencies on other stories (can run in parallel with US1)
- **User Story 3 (P2)**: Can start after Foundational (Phase 2) - No dependencies on other stories (can run in parallel with US1/US2)

### Within Each User Story

- Tests (if included) MUST be written and FAIL before implementation
- SDK client initialization before API calls
- Error mapping implementation before integration
- Story complete before moving to polish phase

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational test tasks marked [P] can run in parallel (within Phase 2)
- User Stories 1, 2, and 3 can run in parallel after Foundational phase (different files, no dependencies)
- Polish tasks marked [P] can run in parallel

---

## Parallel Example: User Stories 1, 2, 3

```bash
# After Foundational phase, all three user stories can run in parallel:

# Developer A: User Story 1 (OpenAI)
Task: "Replace HTTP client with OpenAI SDK client initialization in internal/ai/openai_provider.go"
Task: "Replace HTTP request/response handling with OpenAI SDK API calls in internal/ai/openai_provider.go"

# Developer B: User Story 2 (Anthropic)
Task: "Replace HTTP client with Anthropic SDK client initialization in internal/ai/anthropic_provider.go"
Task: "Replace HTTP request/response handling with Anthropic SDK API calls in internal/ai/anthropic_provider.go"

# Developer C: User Story 3 (Mistral)
Task: "Replace HTTP client with Mistral SDK client initialization in internal/ai/mistral_provider.go"
Task: "Replace HTTP request/response handling with Mistral SDK API calls in internal/ai/mistral_provider.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1 (OpenAI SDK)
4. **STOP and VALIDATE**: Test User Story 1 independently
5. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 (OpenAI) → Test independently → Deploy/Demo (MVP!)
3. Add User Story 2 (Anthropic) → Test independently → Deploy/Demo
4. Add User Story 3 (Mistral) → Test independently → Deploy/Demo
5. Polish phase → Final validation → Deploy

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 (OpenAI SDK)
   - Developer B: User Story 2 (Anthropic SDK)
   - Developer C: User Story 3 (Mistral SDK)
3. Stories complete and integrate independently
4. Polish phase together

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence
- Total tasks: 53
- Tasks per story: US1 (7 tasks), US2 (7 tasks), US3 (7 tasks), Setup (7 tasks), Foundational (12 tests), Polish (13 tasks)
- All existing provider tests must pass without modification (backward compatibility requirement)
