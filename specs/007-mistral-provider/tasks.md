# Tasks: Add Mistral as AI Provider

**Input**: Design documents from `/specs/007-mistral-provider/`
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

**Purpose**: Verify existing code structure and identify integration points

- [X] T001 Verify existing AIProvider interface in internal/ai/provider.go
- [X] T002 Verify existing provider implementations (OpenAI, Anthropic) in internal/ai/
- [X] T003 [P] Verify existing provider selection mechanism in internal/service/commit_service.go
- [X] T004 [P] Verify existing token calculation infrastructure in pkg/tokenization/token_calculator.go

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

### Tests for Mistral Provider (TDD - Write First)

- [X] T005 [P] Write unit test for NewMistralProvider constructor in internal/ai/mistral_provider_test.go
- [X] T006 [P] Write unit test for GenerateCommitMessage with successful API response in internal/ai/mistral_provider_test.go
- [X] T007 [P] Write unit test for GenerateCommitMessage error handling (missing API key) in internal/ai/mistral_provider_test.go
- [X] T008 [P] Write unit test for GenerateCommitMessage error handling (API errors, timeouts) in internal/ai/mistral_provider_test.go
- [X] T009 [P] Write unit test for buildPrompt method in internal/ai/mistral_provider_test.go
- [X] T010 [P] Write integration test for Mistral provider workflow in test/integration/ai_mistral_test.go

**Checkpoint**: Foundation ready - tests written and failing, user story implementation can now begin

---

## Phase 3: User Story 1 - Use Mistral for Commit Message Generation (Priority: P1) 🎯 MVP

**Goal**: Implement MistralProvider that generates commit messages using Mistral AI API, integrate it into provider selection mechanism, and enable users to configure and use Mistral as their AI provider.

**Independent Test**: Configure Mistral API credentials, select Mistral as the provider, run the CLI with AI generation enabled, and verify that Mistral generates a commit message based on repository state.

### Implementation for User Story 1

- [X] T011 [US1] Create MistralProvider struct implementing AIProvider interface in internal/ai/mistral_provider.go
- [X] T012 [US1] Implement NewMistralProvider constructor with HTTP client initialization in internal/ai/mistral_provider.go
- [X] T013 [US1] Implement buildPrompt method to create prompt from repository state in internal/ai/mistral_provider.go
- [X] T014 [US1] Implement GenerateCommitMessage method with Mistral API integration in internal/ai/mistral_provider.go
- [X] T015 [US1] Add error handling for missing API key, API errors, and timeouts in internal/ai/mistral_provider.go
- [X] T016 [US1] Add "mistral" case to provider switch in internal/service/commit_service.go
- [X] T017 [US1] Add "mistral" case to token calculator switch in pkg/tokenization/token_calculator.go
- [X] T018 [US1] Update config.yaml.example with Mistral provider configuration in configs/config.yaml.example

**Checkpoint**: At this point, User Story 1 should be fully functional - Mistral provider can be configured, selected, and used to generate commit messages with proper error handling and fallback

---

## Phase 4: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect the entire feature

- [X] T019 [P] Update README.md with Mistral provider documentation
- [X] T020 [P] Update CHANGELOG.md with feature entry
- [X] T021 Run all unit tests to verify no regressions
- [X] T022 Run all integration tests to verify no regressions
- [X] T023 [P] Verify backward compatibility (existing providers still work)
- [X] T024 [P] Verify Mistral provider works with all supported models (mistral-tiny, mistral-small, mistral-medium, mistral-large-latest)
- [X] T025 [P] Verify token calculation works correctly for Mistral (character-based fallback)
- [X] T026 Code cleanup and formatting (gofmt, goimports)
- [X] T027 Run golangci-lint to verify code quality (golangci-lint not available, but code compiles and all tests pass)
- [X] T028 [P] Run quickstart.md validation scenarios

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
# However, some tasks can be done in parallel:

# Can run in parallel (different files):
Task: "Add 'mistral' case to provider switch in internal/service/commit_service.go"
Task: "Add 'mistral' case to token calculator switch in pkg/tokenization/token_calculator.go"
Task: "Update config.yaml.example with Mistral provider configuration in configs/config.yaml.example"
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
   - Developer A: MistralProvider implementation (mistral_provider.go)
   - Developer B: Integration tasks (commit_service.go, token_calculator.go)
   - Developer C: Configuration and documentation updates
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
- Total tasks: 28
- Tasks per story: US1 (8 tasks), Setup (4 tasks), Foundational (6 tests), Polish (10 tasks)
