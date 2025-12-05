# Tasks: Debug Logging Configuration

**Input**: Design documents from `/specs/003-debug-logging/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Tests are included following TDD approach as required by the constitution.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2)
- Include exact file paths in descriptions

## Path Conventions

- **CLI project**: `cmd/gitcomm/`, `internal/` at repository root
- Paths shown below follow existing project structure

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and verification

- [X] T001 Verify existing project structure and dependencies
- [X] T002 [P] Verify zerolog dependency is available in go.mod
- [X] T003 [P] Verify cobra dependency is available in go.mod

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core logger infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

### Tests for Logger Configuration (TDD - Write First)

- [X] T004 [P] Write unit test for InitLogger with debug=true in internal/utils/logger_test.go
- [X] T005 [P] Write unit test for InitLogger with debug=false in internal/utils/logger_test.go
- [X] T006 [P] Write unit test for InitLogger with debug=true, verbose=true in internal/utils/logger_test.go
- [X] T007 [P] Write unit test for InitLogger with debug=false, verbose=true in internal/utils/logger_test.go
- [X] T008 [P] Write unit test for log output format (raw text, no timestamp) in internal/utils/logger_test.go
- [X] T009 [P] Write unit test for log suppression when debug disabled in internal/utils/logger_test.go

### Implementation for Logger Configuration

- [X] T010 Modify InitLogger function signature to accept debug parameter in internal/utils/logger.go
- [X] T011 Implement debug mode configuration (raw text format, no timestamp, DEBUG level) in internal/utils/logger.go
- [X] T012 Implement silent mode configuration (disabled logger) when debug=false in internal/utils/logger.go
- [X] T013 Implement verbose flag handling (no-op when debug is set) in internal/utils/logger.go
- [X] T014 Configure zerolog.ConsoleWriter for raw text format in internal/utils/logger.go
- [X] T015 Disable timestamps in ConsoleWriter configuration in internal/utils/logger.go

**Checkpoint**: Foundation ready - logger configuration complete, user story implementation can now begin

---

## Phase 3: User Story 1 - Enable Debug Logging with Flag (Priority: P1) 🎯 MVP

**Goal**: Enable debug logging via command-line flag (`--debug` or `-d`), displaying DEBUG-level messages in human-readable structured text format without timestamps.

**Independent Test**: Run `gitcomm --debug` and verify DEBUG-level log messages appear in format `[DEBUG] message key=value` without timestamps.

### Tests for User Story 1 (TDD - Write First)

- [ ] T016 [P] [US1] Write integration test for debug flag parsing in test/integration/debug_flag_test.go
- [ ] T017 [P] [US1] Write integration test for debug flag enabling log output in test/integration/debug_flag_test.go
- [ ] T018 [P] [US1] Write integration test for log format verification (raw text, no timestamp) in test/integration/debug_flag_test.go
- [ ] T019 [P] [US1] Write integration test for short flag form (-d) in test/integration/debug_flag_test.go

### Implementation for User Story 1

- [X] T020 [US1] Add debug flag variable declaration in cmd/gitcomm/main.go
- [X] T021 [US1] Register --debug and -d flags with Cobra in cmd/gitcomm/main.go
- [X] T022 [US1] Update InitLogger call to pass debug flag in cmd/gitcomm/main.go
- [X] T023 [US1] Update all Logger.Info() calls to Logger.Debug() in internal/service/commit_service.go
- [X] T024 [US1] Update all Logger.Warn() calls to Logger.Debug() in internal/service/commit_service.go
- [X] T025 [US1] Update all Logger.Error() calls to Logger.Debug() in internal/service/commit_service.go
- [X] T026 [US1] Update Logger.Warn() call to Logger.Debug() in internal/ai/local_provider.go
- [X] T027 [US1] Update Logger.Warn() call to Logger.Debug() in internal/ai/anthropic_provider.go
- [X] T028 [US1] Update Logger.Warn() call to Logger.Debug() in internal/ai/openai_provider.go
- [X] T029 [US1] Update Logger.Warn() call to Logger.Debug() in internal/ui/prompts.go
- [X] T030 [US1] Update Logger.Warn() and Logger.Info() calls to Logger.Debug() in cmd/gitcomm/main.go

**Checkpoint**: At this point, User Story 1 should be fully functional - debug flag enables logging in raw text format without timestamps

---

## Phase 4: User Story 2 - Default Silent Operation (Priority: P1)

**Goal**: Ensure CLI runs silently by default (no log output) unless debug flag is enabled, maintaining backward compatibility and clean user experience.

**Independent Test**: Run `gitcomm` without flags and verify no log messages appear in output, while error messages are still displayed.

### Tests for User Story 2 (TDD - Write First)

- [ ] T031 [P] [US2] Write integration test for silent operation by default in test/integration/silent_mode_test.go
- [ ] T032 [P] [US2] Write integration test for error messages still displayed when debug disabled in test/integration/silent_mode_test.go
- [ ] T033 [P] [US2] Write integration test for verbose flag being no-op without debug flag in test/integration/silent_mode_test.go

### Implementation for User Story 2

- [X] T034 [US2] Verify default logger configuration is silent (disabled) in internal/utils/logger.go
- [X] T035 [US2] Verify error messages use fmt.Printf (not logger) in cmd/gitcomm/main.go
- [X] T036 [US2] Verify error messages use fmt.Printf (not logger) in internal/service/commit_service.go

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently - debug flag enables logging, default is silent, error messages always shown

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [X] T037 [P] Update README.md with debug flag documentation
- [X] T038 [P] Update CHANGELOG.md with debug logging feature entry
- [X] T039 Run all unit tests to verify no regressions
- [X] T040 Run all integration tests to verify no regressions
- [X] T041 [P] Verify all logging calls use Logger.Debug() (code review/audit)
- [X] T042 [P] Run quickstart.md validation scenarios
- [X] T043 Code cleanup and formatting (gofmt, goimports)
- [X] T044 Run golangci-lint to verify code quality

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (US1 → US2)
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P1)**: Can start after Foundational (Phase 2) - Depends on US1 for verification (but can be tested independently)

### Within Each User Story

- Tests (if included) MUST be written and FAIL before implementation
- Logger configuration before flag implementation
- Flag implementation before logging call updates
- Core implementation before integration
- Story complete before moving to next priority

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational test tasks marked [P] can run in parallel (within Phase 2)
- All Foundational implementation tasks can run sequentially (they modify the same file)
- Once Foundational phase completes, user stories can start
- All tests for a user story marked [P] can run in parallel
- Logging call updates within a story marked [P] can run in parallel (different files)
- Different user stories can be worked on in parallel by different team members

---

## Parallel Example: User Story 1

```bash
# Launch all tests for User Story 1 together:
Task: "Write integration test for debug flag parsing in test/integration/debug_flag_test.go"
Task: "Write integration test for debug flag enabling log output in test/integration/debug_flag_test.go"
Task: "Write integration test for log format verification in test/integration/debug_flag_test.go"
Task: "Write integration test for short flag form (-d) in test/integration/debug_flag_test.go"

# Launch all logging call updates together (different files):
Task: "Update all Logger.Info() calls to Logger.Debug() in internal/service/commit_service.go"
Task: "Update Logger.Warn() call to Logger.Debug() in internal/ai/local_provider.go"
Task: "Update Logger.Warn() call to Logger.Debug() in internal/ai/anthropic_provider.go"
Task: "Update Logger.Warn() call to Logger.Debug() in internal/ai/openai_provider.go"
Task: "Update Logger.Warn() call to Logger.Debug() in internal/ui/prompts.go"
Task: "Update Logger.Warn() and Logger.Info() calls to Logger.Debug() in cmd/gitcomm/main.go"
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
4. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 (flag + logging updates)
   - Developer B: User Story 2 (verification + tests)
3. Stories complete and integrate independently

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence
- Total tasks: 44
- Tasks per story: US1 (15 tasks), US2 (6 tasks), Setup (3 tasks), Foundational (12 tasks), Polish (8 tasks)
