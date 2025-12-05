# Tasks: Git Commit Message Automation CLI

**Input**: Design documents from `/specs/001-git-commit-cli/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: TDD is mandatory per constitution - all core business logic must have tests written first.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Single project**: Repository root structure per plan.md
- Paths use absolute structure: `cmd/`, `internal/`, `pkg/`, `test/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [ ] T001 Create project structure per implementation plan (cmd/, internal/, pkg/, test/, configs/)
- [ ] T002 Initialize Go module with go.mod in repository root
- [ ] T003 [P] Add dependencies: github.com/spf13/cobra, github.com/charmbracelet/bubbletea, github.com/charmbracelet/lipgloss, github.com/spf13/viper, github.com/go-git/go-git/v5, github.com/rs/zerolog, github.com/onsi/ginkgo, github.com/onsi/gomega
- [ ] T004 [P] Configure golangci-lint for code quality checks
- [ ] T005 [P] Setup goimports/gofumpt for code formatting
- [ ] T006 [P] Create example configuration file at configs/config.yaml.example

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [ ] T007 Create error types in internal/utils/errors.go (ErrNotGitRepository, ErrNoChanges, ErrInvalidFormat, ErrAIProviderUnavailable, ErrEmptySubject, ErrTokenCalculationFailed)
- [ ] T008 [P] Create logging infrastructure using zerolog in internal/utils/logger.go
- [ ] T009 [P] Create configuration loading infrastructure in internal/config/config.go (using viper, support ~/.gitcomm/config.yaml)
- [ ] T010 Create base domain models in internal/model/ (commit_message.go, repository_state.go, config.go per data-model.md)
- [ ] T011 Create GitRepository interface in internal/repository/git_repository.go
- [ ] T012 Create MessageValidator interface in pkg/conventional/validator.go
- [ ] T013 Create CLI entrypoint structure in cmd/gitcomm/main.go with Cobra root command

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Manual Commit Message Creation (Priority: P1) 🎯 MVP

**Goal**: Enable developers to create Conventional Commits compliant commit messages manually through interactive prompts, without AI assistance.

**Independent Test**: Run CLI in git repository, decline AI assistance, manually enter scope/subject/body/footer, verify commit message follows Conventional Commits format and commit is created successfully.

### Tests for User Story 1 ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T014 [P] [US1] Create unit test for CommitMessage model validation in internal/model/commit_message_test.go
- [ ] T015 [P] [US1] Create unit test for Conventional Commits validator in pkg/conventional/validator_test.go
- [ ] T016 [P] [US1] Create unit test for message formatting service in internal/service/formatting_service_test.go
- [ ] T017 [US1] Create integration test for manual commit workflow in test/integration/manual_commit_test.go

### Implementation for User Story 1

- [ ] T018 [P] [US1] Implement CommitMessage model in internal/model/commit_message.go (Type, Scope, Subject, Body, Footer, Signoff fields)
- [ ] T019 [P] [US1] Implement RepositoryState model in internal/model/repository_state.go (StagedFiles, UnstagedFiles, IsEmpty, HasChanges)
- [ ] T020 [P] [US1] Implement FileChange struct in internal/model/repository_state.go (Path, Status, Diff)
- [ ] T021 [US1] Implement GitRepository interface using go-git in internal/repository/git_repository_impl.go (GetRepositoryState, CreateCommit methods)
- [ ] T022 [US1] Implement Conventional Commits validator in pkg/conventional/validator.go (validate type, scope, subject, body, footer)
- [ ] T023 [US1] Implement formatting service in internal/service/formatting_service.go (format CommitMessage to Conventional Commits string)
- [ ] T024 [US1] Implement validation service in internal/service/validation_service.go (validate CommitMessage using validator)
- [ ] T025 [US1] Implement interactive prompts using bubbletea in internal/ui/prompts.go (scope, subject, body, footer prompts with validation)
- [ ] T026 [US1] Implement message display formatter in internal/ui/display.go (format message for review)
- [ ] T027 [US1] Implement commit service in internal/service/commit_service.go (orchestrate manual workflow: prompts → validation → commit)
- [ ] T028 [US1] Integrate commit service into CLI command in cmd/gitcomm/main.go (handle manual workflow, error cases)
- [ ] T029 [US1] Add error handling for empty subject rejection in internal/ui/prompts.go (re-prompt until non-empty)
- [ ] T030 [US1] Add empty commit confirmation prompt in internal/ui/prompts.go (when no changes detected)
- [ ] T031 [US1] Add logging for commit operations in internal/service/commit_service.go

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently - users can create commit messages manually without AI.

---

## Phase 4: User Story 2 - AI-Assisted Commit Message Generation (Priority: P2)

**Goal**: Enable developers to use AI to automatically generate commit messages based on repository state, with token calculation and format validation.

**Independent Test**: Run CLI in git repository, accept AI assistance, verify token calculation displayed, AI provider called, message generated and validated, user can accept/reject AI suggestion.

### Tests for User Story 2 ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T032 [P] [US2] Create unit test for AIProvider interface mock in test/mocks/ai_provider_mock.go
- [ ] T033 [P] [US2] Create unit test for token calculator in pkg/tokenization/token_calculator_test.go
- [ ] T034 [P] [US2] Create unit test for OpenAI provider in internal/ai/openai_provider_test.go
- [ ] T035 [P] [US2] Create unit test for Anthropic provider in internal/ai/anthropic_provider_test.go
- [ ] T036 [US2] Create integration test for AI-assisted commit workflow in test/integration/ai_commit_test.go

### Implementation for User Story 2

- [ ] T037 [P] [US2] Create AIProvider interface in internal/ai/provider.go (GenerateCommitMessage method)
- [ ] T038 [P] [US2] Create TokenCalculator interface in internal/ai/token_calculator.go (CalculateTokens method)
- [ ] T039 [P] [US2] Implement OpenAI tokenization using tiktoken in pkg/tokenization/tiktoken.go
- [ ] T040 [P] [US2] Implement Anthropic tokenization (custom) in pkg/tokenization/anthropic.go
- [ ] T041 [P] [US2] Implement character-based fallback tokenization in pkg/tokenization/fallback.go
- [ ] T042 [US2] Implement token calculator service in internal/ai/token_calculator.go (provider-specific calculation with fallback)
- [ ] T043 [P] [US2] Implement OpenAI provider in internal/ai/openai_provider.go (API integration, error handling, timeout)
- [ ] T044 [P] [US2] Implement Anthropic provider in internal/ai/anthropic_provider.go (API integration, error handling, timeout)
- [ ] T045 [P] [US2] Implement local model provider in internal/ai/local_provider.go (custom endpoint support)
- [ ] T046 [US2] Implement AI provider selection logic in internal/config/config.go (load from config file, support CLI flag override)
- [ ] T047 [US2] Add token count display prompt in internal/ui/prompts.go (show estimated tokens before AI decision)
- [ ] T048 [US2] Add AI usage decision prompt in internal/ui/prompts.go (y/n choice)
- [ ] T049 [US2] Integrate AI provider into commit service in internal/service/commit_service.go (call AI, handle errors, fallback to manual)
- [ ] T050 [US2] Add AI-generated message validation and edit option in internal/ui/prompts.go (validate format, offer edit/use-with-warning)
- [ ] T051 [US2] Add error handling for AI provider failures in internal/service/commit_service.go (network errors, timeouts → fallback to manual with clear message)
- [ ] T052 [US2] Add logging for AI operations in internal/service/commit_service.go (token calculation, provider calls, errors)

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently - users can choose manual or AI-assisted commit message creation.

---

## Phase 5: User Story 3 - CLI Options and Configuration (Priority: P3)

**Goal**: Enable developers to configure CLI behavior via command-line options (auto-stage files, disable signoff).

**Independent Test**: Run CLI with -a flag, verify files auto-staged; run with -s flag, verify no signoff; run without flags, verify default behavior.

### Tests for User Story 3 ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T053 [P] [US3] Create unit test for CommitOptions model in internal/model/config_test.go
- [ ] T054 [US3] Create integration test for CLI options in test/integration/cli_options_test.go

### Implementation for User Story 3

- [ ] T055 [US3] Implement CommitOptions model in internal/model/config.go (AutoStage, NoSignoff, AIProvider, SkipAI fields)
- [ ] T056 [US3] Add -a/--add-all flag to Cobra command in cmd/gitcomm/main.go
- [ ] T057 [US3] Add -s/--no-signoff flag to Cobra command in cmd/gitcomm/main.go
- [ ] T058 [US3] Add --provider flag to Cobra command in cmd/gitcomm/main.go
- [ ] T059 [US3] Add --skip-ai flag to Cobra command in cmd/gitcomm/main.go
- [ ] T060 [US3] Implement auto-staging logic in internal/repository/git_repository_impl.go (git add -A equivalent)
- [ ] T061 [US3] Integrate auto-staging into commit workflow in internal/service/commit_service.go (check AutoStage flag, stage before proceeding)
- [ ] T062 [US3] Integrate signoff control into commit creation in internal/repository/git_repository_impl.go (include/exclude Signed-off-by based on NoSignoff flag)
- [ ] T063 [US3] Add logging for CLI options in cmd/gitcomm/main.go

**Checkpoint**: At this point, all user stories should be independently functional - CLI supports manual and AI-assisted workflows with configurable options.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [ ] T064 [P] Add comprehensive error messages for all error cases in internal/utils/errors.go
- [ ] T065 [P] Add help documentation and usage examples in cmd/gitcomm/main.go
- [ ] T066 [P] Add README.md with installation and usage instructions
- [ ] T067 [P] Add CHANGELOG.md entry for this feature
- [ ] T068 Code cleanup and refactoring (review all files for consistency)
- [ ] T069 [P] Add additional unit tests for edge cases in test/unit/
- [ ] T070 Validate quickstart.md scenarios work correctly
- [ ] T071 [P] Add integration tests for error scenarios (not git repo, no changes, AI failures)
- [ ] T072 Performance optimization (ensure <100ms prompt response time per SC-005)
- [ ] T073 Security review (ensure no secrets in logs, config file permissions)

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
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - Uses validation/formatting from US1 but independently testable
- **User Story 3 (P3)**: Can start after Foundational (Phase 2) - Enhances US1/US2 workflows but independently testable

### Within Each User Story

- Tests (TDD) MUST be written and FAIL before implementation
- Models before services
- Services before UI/integration
- Core implementation before integration
- Story complete before moving to next priority

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel (T003-T006)
- All Foundational tasks marked [P] can run in parallel (T008-T009)
- Once Foundational phase completes, all user stories can start in parallel (if team capacity allows)
- All tests for a user story marked [P] can run in parallel
- Models within a story marked [P] can run in parallel
- Different user stories can be worked on in parallel by different team members

---

## Parallel Example: User Story 1

```bash
# Launch all tests for User Story 1 together:
Task: "Create unit test for CommitMessage model validation in internal/model/commit_message_test.go"
Task: "Create unit test for Conventional Commits validator in pkg/conventional/validator_test.go"
Task: "Create unit test for message formatting service in internal/service/formatting_service_test.go"

# Launch all models for User Story 1 together:
Task: "Implement CommitMessage model in internal/model/commit_message.go"
Task: "Implement RepositoryState model in internal/model/repository_state.go"
Task: "Implement FileChange struct in internal/model/repository_state.go"
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
   - Developer A: User Story 1
   - Developer B: User Story 2 (can start in parallel)
   - Developer C: User Story 3 (can start in parallel)
3. Stories complete and integrate independently

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing (TDD)
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence
