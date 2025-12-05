# Tasks: Unify AI Provider Prompts with Validation Rules

**Input**: Design documents from `/specs/011-unify-ai-prompts/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: TDD is mandatory per constitution - all core business logic must have tests written first.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Single project**: Repository root structure per plan.md
- Paths use absolute structure: `pkg/`, `internal/`, `test/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [x] T001 Create directory structure for `pkg/ai/prompt/` per implementation plan
- [x] T002 [P] Create error types in `pkg/ai/prompt/errors.go` for prompt generation errors (ErrNilValidator, ErrNilRepositoryState, ErrRuleExtractionFailed)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T003 Extend `MessageValidator` interface in `pkg/conventional/validator.go` with `GetValidTypes() []string` method
- [x] T004 Extend `MessageValidator` interface in `pkg/conventional/validator.go` with `GetSubjectMaxLength() int` method
- [x] T005 Extend `MessageValidator` interface in `pkg/conventional/validator.go` with `GetBodyMaxLength() int` method
- [x] T006 Extend `MessageValidator` interface in `pkg/conventional/validator.go` with `GetScopeFormatDescription() string` method
- [x] T007 [P] Implement `GetValidTypes()` method in `Validator` struct in `pkg/conventional/validator.go` returning `["feat", "fix", "docs", "style", "refactor", "test", "chore", "version"]`
- [x] T008 [P] Implement `GetSubjectMaxLength()` method in `Validator` struct in `pkg/conventional/validator.go` returning `72`
- [x] T009 [P] Implement `GetBodyMaxLength()` method in `Validator` struct in `pkg/conventional/validator.go` returning `320`
- [x] T010 [P] Implement `GetScopeFormatDescription()` method in `Validator` struct in `pkg/conventional/validator.go` returning `"alphanumeric, hyphens, underscores only"`
- [x] T011 [P] Add unit tests for `GetValidTypes()` in `pkg/conventional/validator_test.go`
- [x] T012 [P] Add unit tests for `GetSubjectMaxLength()` in `pkg/conventional/validator_test.go`
- [x] T013 [P] Add unit tests for `GetBodyMaxLength()` in `pkg/conventional/validator_test.go`
- [x] T014 [P] Add unit tests for `GetScopeFormatDescription()` in `pkg/conventional/validator_test.go`

**Checkpoint**: Foundation ready - MessageValidator extension complete, user story implementation can now begin

---

## Phase 3: User Story 1 - Unified Prompt Across All AI Providers (Priority: P1) 🎯 MVP

**Goal**: All AI providers (OpenAI, Anthropic, Mistral, local) use identical system and user messages that include validation rules extracted dynamically from MessageValidator. This ensures consistent commit message generation across all providers and guarantees that AI-generated messages pass validation.

**Independent Test**: Configure gitcomm with different AI providers, generate commit messages for the same repository state, and verify that all providers generate commit messages that:
1. Follow the same format and structure
2. Pass MessageValidator validation
3. Include all required validation constraints in their generation

### Tests for User Story 1 ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T015 [P] [US1] Write unit test for `PromptGenerator.GenerateSystemMessage()` with valid validator in `pkg/ai/prompt/generator_test.go` - test should fail initially
- [x] T016 [P] [US1] Write unit test for `PromptGenerator.GenerateSystemMessage()` with nil validator (error case) in `pkg/ai/prompt/generator_test.go` - test should fail initially
- [x] T017 [P] [US1] Write unit test for `PromptGenerator.GenerateUserMessage()` with valid repository state in `pkg/ai/prompt/generator_test.go` - test should fail initially
- [x] T018 [P] [US1] Write unit test for `PromptGenerator.GenerateUserMessage()` with nil repository state (error case) in `pkg/ai/prompt/generator_test.go` - test should fail initially
- [x] T019 [P] [US1] Write unit test for `PromptGenerator.GenerateUserMessage()` with empty repository state in `pkg/ai/prompt/generator_test.go` - test should fail initially
- [x] T020 [P] [US1] Write unit test for prompt consistency (same inputs produce same outputs) in `pkg/ai/prompt/generator_test.go` - test should fail initially
- [x] T021 [P] [US1] Write integration test for prompt consistency across all providers in `test/integration/prompt_unification_test.go` - test should fail initially
- [x] T022 [P] [US1] Write integration test for Anthropic system/user message combination in `test/integration/prompt_unification_test.go` - test should fail initially

### Implementation for User Story 1

- [x] T023 [US1] Define `PromptGenerator` interface in `pkg/ai/prompt/generator.go` with `GenerateSystemMessage(validator conventional.MessageValidator) (string, error)` and `GenerateUserMessage(repoState *model.RepositoryState) (string, error)` methods
- [x] T024 [US1] Implement `UnifiedPromptGenerator` struct in `pkg/ai/prompt/generator.go` implementing `PromptGenerator` interface
- [x] T025 [US1] Implement `NewUnifiedPromptGenerator()` constructor function in `pkg/ai/prompt/generator.go` returning `PromptGenerator`
- [x] T026 [US1] Implement `GenerateSystemMessage()` method in `UnifiedPromptGenerator` in `pkg/ai/prompt/generator.go` - extract validation rules from validator and format as structured bullet points
- [x] T027 [US1] Implement `GenerateUserMessage()` method in `UnifiedPromptGenerator` in `pkg/ai/prompt/generator.go` - format repository state (staged/unstaged files, diffs) as user message
- [x] T028 [US1] Update `OpenAIProvider.GenerateCommitMessage()` in `internal/ai/openai_provider.go` to use `PromptGenerator` instead of `buildPrompt()` method
- [x] T029 [US1] Update `AnthropicProvider.GenerateCommitMessage()` in `internal/ai/anthropic_provider.go` to use `PromptGenerator` and prepend system message to user message
- [x] T030 [US1] Update `MistralProvider.GenerateCommitMessage()` in `internal/ai/mistral_provider.go` to use `PromptGenerator` instead of `buildPrompt()` method
- [x] T031 [US1] Update `LocalProvider.GenerateCommitMessage()` in `internal/ai/local_provider.go` to use `PromptGenerator` instead of `buildPrompt()` method
- [x] T032 [US1] Remove `buildPrompt()` method from `OpenAIProvider` in `internal/ai/openai_provider.go`
- [x] T033 [US1] Remove `buildPrompt()` method from `AnthropicProvider` in `internal/ai/anthropic_provider.go`
- [x] T034 [US1] Remove `buildPrompt()` method from `MistralProvider` in `internal/ai/mistral_provider.go`
- [x] T035 [US1] Remove `buildPrompt()` method from `LocalProvider` in `internal/ai/local_provider.go`
- [x] T036 [US1] Update provider constructors to inject `PromptGenerator` dependency (or create it internally) in `internal/ai/*_provider.go` files

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently. All providers use unified prompts with validation rules.

---

## Phase 4: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [x] T037 [P] Update `README.md` to document unified prompt generation feature
- [x] T038 [P] Update `CHANGELOG.md` with feature description and breaking changes (none)
- [x] T039 [P] Verify all existing provider tests still pass after prompt unification
- [x] T040 [P] Run quickstart.md validation to ensure implementation matches design
- [x] T041 [P] Code cleanup: remove any unused imports or dead code from provider files
- [x] T042 [P] Add code comments documenting prompt generator usage in provider implementations

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

- Tests (T015-T022) MUST be written and FAIL before implementation
- Interface definition (T023) before implementation (T024-T027)
- PromptGenerator implementation (T024-T027) before provider updates (T028-T031)
- Provider updates (T028-T031) before removing old methods (T032-T035)
- Core implementation complete before integration tests

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational tasks marked [P] can run in parallel (within Phase 2)
- All test tasks for User Story 1 marked [P] can run in parallel
- Provider update tasks (T028-T031) can run in parallel (different files)
- Provider cleanup tasks (T032-T035) can run in parallel (different files)
- All Polish tasks marked [P] can run in parallel

---

## Parallel Example: User Story 1

```bash
# Launch all tests for User Story 1 together:
Task: "Write unit test for GenerateSystemMessage() with valid validator"
Task: "Write unit test for GenerateSystemMessage() with nil validator"
Task: "Write unit test for GenerateUserMessage() with valid repository state"
Task: "Write unit test for GenerateUserMessage() with nil repository state"
Task: "Write unit test for GenerateUserMessage() with empty repository state"
Task: "Write unit test for prompt consistency"
Task: "Write integration test for prompt consistency across all providers"
Task: "Write integration test for Anthropic system/user message combination"

# Launch provider updates in parallel (after PromptGenerator is implemented):
Task: "Update OpenAIProvider to use PromptGenerator"
Task: "Update AnthropicProvider to use PromptGenerator"
Task: "Update MistralProvider to use PromptGenerator"
Task: "Update LocalProvider to use PromptGenerator"

# Launch provider cleanup in parallel:
Task: "Remove buildPrompt() from OpenAIProvider"
Task: "Remove buildPrompt() from AnthropicProvider"
Task: "Remove buildPrompt() from MistralProvider"
Task: "Remove buildPrompt() from LocalProvider"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Test User Story 1 independently
   - Verify all providers use identical prompts
   - Verify validation rules are included
   - Verify Anthropic prepends system to user correctly
   - Verify all existing tests still pass
5. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Deploy/Demo (MVP!)
3. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: PromptGenerator implementation (T023-T027)
   - Developer B: Write all tests (T015-T022)
3. Once PromptGenerator is complete:
   - Developer A: Update OpenAIProvider and MistralProvider (T028, T030)
   - Developer B: Update AnthropicProvider and LocalProvider (T029, T031)
4. All developers: Remove old buildPrompt() methods in parallel (T032-T035)

---

## Notes

- [P] tasks = different files, no dependencies
- [US1] label maps task to User Story 1 for traceability
- User Story 1 should be independently completable and testable
- Verify tests fail before implementing (TDD approach)
- Commit after each task or logical group
- Stop at checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence
- All provider updates maintain backward compatibility (AIProvider interface unchanged)
