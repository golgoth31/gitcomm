# Tasks: Environment Variable Placeholder Substitution in Config Files

**Input**: Design documents from `/specs/014-env-var-placeholders/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/, quickstart.md

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
- [X] T002 [P] Verify os standard library package is available for environment variable access
- [X] T003 [P] Verify regexp standard library package is available for placeholder pattern matching
- [X] T004 [P] Verify strings standard library package is available for string manipulation
- [X] T005 [P] Verify github.com/spf13/viper dependency is available in go.mod
- [X] T006 [P] Verify existing LoadConfig function signature and behavior

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core placeholder processing infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

### Tests for Placeholder Processing Logic (TDD - Write First)

- [X] T007 [P] Write unit test for placeholder identification using regex in internal/config/config_test.go
- [X] T008 [P] Write unit test for placeholder syntax validation (valid patterns) in internal/config/config_test.go
- [X] T009 [P] Write unit test for invalid placeholder syntax detection (spaces, nested, multiline) in internal/config/config_test.go
- [X] T010 [P] Write unit test for environment variable lookup using os.LookupEnv() in internal/config/config_test.go
- [X] T011 [P] Write unit test for placeholder substitution (single placeholder) in internal/config/config_test.go
- [X] T012 [P] Write unit test for multiple placeholder substitution in internal/config/config_test.go
- [X] T013 [P] Write unit test for comment line handling (skip placeholders in comments) in internal/config/config_test.go
- [X] T014 [P] Write unit test for empty string value handling in internal/config/config_test.go

### Implementation for Placeholder Processing Logic

- [X] T015 Add placeholder regex pattern compilation in internal/config/config.go
- [X] T016 Add function to identify all placeholders in config content using regex in internal/config/config.go
- [X] T017 Add function to validate placeholder syntax (check for nested, multiline, invalid chars) in internal/config/config.go
- [X] T018 Add function to extract environment variable names from placeholders in internal/config/config.go
- [X] T019 Add function to validate all environment variables exist using os.LookupEnv() in internal/config/config.go
- [X] T020 Add function to substitute placeholders with environment variable values in internal/config/config.go
- [X] T021 Add function to skip comment lines during processing in internal/config/config.go
- [X] T022 Add error handling for invalid placeholder syntax with clear error messages in internal/config/config.go
- [X] T023 Add error handling for missing environment variables with clear error messages listing all missing vars in internal/config/config.go

**Checkpoint**: Foundation ready - placeholder processing logic complete, user story implementation can now begin

---

## Phase 3: User Story 1 - Basic Environment Variable Substitution (Priority: P1) 🎯 MVP

**Goal**: Automatically replace `${ENV_VAR_NAME}` placeholders in config files with values from environment variables, enabling secure configuration management by separating secrets from config files.

**Independent Test**: Create a config file with `${API_KEY}` placeholder, set the environment variable, call LoadConfig, and verify the placeholder is replaced with the environment variable value in the loaded configuration.

### Tests for User Story 1 (TDD - Write First)

- [X] T024 [P] [US1] Write integration test for single placeholder substitution in test/integration/config_test.go
- [X] T025 [P] [US1] Write integration test for multiple placeholder substitution in test/integration/config_test.go
- [X] T026 [P] [US1] Write integration test for placeholder in nested YAML structure in test/integration/config_test.go
- [X] T027 [P] [US1] Write integration test for backward compatibility (config without placeholders) in test/integration/config_test.go
- [X] T028 [P] [US1] Write integration test for empty string value substitution in test/integration/config_test.go

### Implementation for User Story 1

- [X] T029 [US1] Integrate placeholder identification into LoadConfig function in internal/config/config.go
- [X] T030 [US1] Integrate placeholder substitution into LoadConfig function (before viper.ReadInConfig) in internal/config/config.go
- [X] T031 [US1] Modify LoadConfig to read config file as text, perform substitution, then pass to viper via strings.NewReader() in internal/config/config.go
- [X] T032 [US1] Ensure substitution happens before YAML parsing in LoadConfig function in internal/config/config.go
- [X] T033 [US1] Ensure backward compatibility (config files without placeholders work unchanged) in internal/config/config.go

**Checkpoint**: At this point, User Story 1 should be fully functional - placeholders are automatically replaced with environment variable values

---

## Phase 4: User Story 2 - Missing Environment Variable Error Handling (Priority: P1)

**Goal**: Exit immediately with clear error messages when config file contains placeholders for environment variables that are not set, preventing the application from running with incomplete configuration.

**Independent Test**: Create a config file with `${MISSING_VAR}` placeholder, ensure the environment variable is not set, attempt to load the config, and verify the application exits with an error message identifying the missing variable.

### Tests for User Story 2 (TDD - Write First)

- [X] T034 [P] [US2] Write integration test for missing single environment variable error in test/integration/config_test.go
- [X] T035 [P] [US2] Write integration test for missing multiple environment variables error in test/integration/config_test.go
- [X] T036 [P] [US2] Write integration test for error message clarity (lists all missing variables) in test/integration/config_test.go
- [X] T037 [P] [US2] Write integration test for empty string value treated as valid (does not exit) in test/integration/config_test.go

### Implementation for User Story 2

- [X] T038 [US2] Integrate environment variable validation into LoadConfig function (check all placeholders before substitution) in internal/config/config.go
- [X] T039 [US2] Add error handling to exit immediately when environment variables are missing in internal/config/config.go
- [X] T040 [US2] Ensure error messages clearly identify all missing environment variables in internal/config/config.go
- [X] T041 [US2] Ensure validation happens before any substitution (fail-fast) in internal/config/config.go

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently - placeholders are substituted, missing variables cause immediate exit with clear errors

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [X] T042 [P] Update README.md with environment variable placeholder documentation
- [X] T043 [P] Update CHANGELOG.md with placeholder substitution feature entry
- [X] T044 Run all unit tests to verify no regressions
- [X] T045 Run all integration tests to verify no regressions
- [X] T046 [P] Verify error messages are clear and actionable (code review/audit)
- [X] T047 [P] Run quickstart.md validation scenarios
- [X] T048 Code cleanup and formatting (gofmt, goimports)
- [X] T049 Run golangci-lint to verify code quality
- [X] T050 Verify backward compatibility (config files without placeholders work unchanged)
- [X] T051 Verify performance (placeholder substitution completes in under 10ms for typical config files)

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
- **User Story 2 (P1)**: Can start after Foundational (Phase 2) - Can be implemented independently but tests error scenarios from US1

### Within Each User Story

- Tests (if included) MUST be written and FAIL before implementation
- Placeholder identification before validation
- Validation before substitution
- Substitution before YAML parsing
- Core implementation before integration
- Story complete before moving to next priority

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational test tasks marked [P] can run in parallel (within Phase 2)
- All Foundational implementation tasks can run sequentially (they modify the same file)
- Once Foundational phase completes, user stories can start
- All tests for a user story marked [P] can run in parallel
- Different user stories can be worked on in parallel by different team members

---

## Parallel Example: User Story 1

```bash
# Launch all tests for User Story 1 together:
Task: "Write integration test for single placeholder substitution in test/integration/config_test.go"
Task: "Write integration test for multiple placeholder substitution in test/integration/config_test.go"
Task: "Write integration test for placeholder in nested YAML structure in test/integration/config_test.go"
Task: "Write integration test for backward compatibility (config without placeholders) in test/integration/config_test.go"
Task: "Write integration test for empty string value substitution in test/integration/config_test.go"
```

---

## Parallel Example: User Story 2

```bash
# Launch all tests for User Story 2 together:
Task: "Write integration test for missing single environment variable error in test/integration/config_test.go"
Task: "Write integration test for missing multiple environment variables error in test/integration/config_test.go"
Task: "Write integration test for error message clarity (lists all missing variables) in test/integration/config_test.go"
Task: "Write integration test for empty string value treated as valid (does not exit) in test/integration/config_test.go"
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
   - Developer A: User Story 1 (placeholder substitution)
   - Developer B: User Story 2 (error handling)
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
- Total tasks: 51
- Tasks per story: US1 (10 tasks), US2 (8 tasks), Setup (6 tasks), Foundational (17 tasks), Polish (10 tasks)
