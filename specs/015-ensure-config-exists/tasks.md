# Tasks: Ensure Config File Exists Before Reading

**Input**: Design documents from `/specs/015-ensure-config-exists/`
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
- [X] T002 [P] Verify os and path/filepath standard library packages are available
- [X] T003 [P] Verify github.com/spf13/viper dependency is available in go.mod
- [X] T004 [P] Verify github.com/rs/zerolog dependency is available in go.mod
- [X] T005 [P] Verify utils.Logger is accessible from internal/config package

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core file creation infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

### Tests for File Creation Logic (TDD - Write First)

- [X] T006 [P] Write unit test for file existence check in internal/config/config_test.go
- [X] T007 [P] Write unit test for empty file creation (0 bytes) in internal/config/config_test.go
- [X] T008 [P] Write unit test for file permissions (0600) in internal/config/config_test.go
- [X] T009 [P] Write unit test for parent directory creation (0755) in internal/config/config_test.go
- [X] T010 [P] Write unit test for path validation (directory check) in internal/config/config_test.go
- [X] T011 [P] Write unit test for race condition handling in internal/config/config_test.go

### Implementation for File Creation Logic

- [X] T012 Add file existence check using os.Stat() in internal/config/config.go
- [X] T013 Add path validation (check if path is directory) in internal/config/config.go
- [X] T014 Add parent directory creation using os.MkdirAll() with 0755 permissions in internal/config/config.go
- [X] T015 Add empty file creation using os.OpenFile() with O_CREATE|O_WRONLY|O_EXCL flags in internal/config/config.go
- [X] T016 Add file permission setting (0600) using os.Chmod() in internal/config/config.go
- [X] T017 Add race condition handling (check os.IsExist) in internal/config/config.go
- [X] T018 Add logging for file creation using utils.Logger.Debug() in internal/config/config.go

**Checkpoint**: Foundation ready - file creation logic complete, user story implementation can now begin

---

## Phase 3: User Story 1 - Config File Auto-Creation (Priority: P1) 🎯 MVP

**Goal**: Automatically create an empty config file if it doesn't exist when loading configuration, ensuring the application can always write configuration settings without manual file creation.

**Independent Test**: Delete the config file (if it exists), call LoadConfig, and verify that an empty config file is created at the expected location before the function returns.

### Tests for User Story 1 (TDD - Write First)

- [X] T019 [P] [US1] Write integration test for file creation when missing in test/integration/config_test.go
- [X] T020 [P] [US1] Write integration test for existing file not modified in test/integration/config_test.go
- [X] T021 [P] [US1] Write integration test for parent directory creation in test/integration/config_test.go
- [X] T022 [P] [US1] Write integration test for custom config path in test/integration/config_test.go
- [X] T023 [P] [US1] Write integration test for default config path (~/.gitcomm/config.yaml) in test/integration/config_test.go

### Implementation for User Story 1

- [X] T024 [US1] Integrate file existence check into LoadConfig function in internal/config/config.go
- [X] T025 [US1] Integrate parent directory creation into LoadConfig function in internal/config/config.go
- [X] T026 [US1] Integrate empty file creation into LoadConfig function in internal/config/config.go
- [X] T027 [US1] Integrate file permission setting into LoadConfig function in internal/config/config.go
- [X] T028 [US1] Integrate logging into LoadConfig function in internal/config/config.go
- [X] T029 [US1] Ensure file creation happens before v.ReadInConfig() call in internal/config/config.go

**Checkpoint**: At this point, User Story 1 should be fully functional - config file is automatically created when missing

---

## Phase 4: User Story 2 - Config File Creation Error Handling (Priority: P2)

**Goal**: Handle file creation errors gracefully, providing clear error messages when config file creation fails due to permissions or other system issues.

**Independent Test**: Attempt to create config file in a read-only directory, call LoadConfig, and verify that an appropriate error is returned explaining the file creation failure.

### Tests for User Story 2 (TDD - Write First)

- [X] T030 [P] [US2] Write unit test for read-only directory error in internal/config/config_test.go
- [X] T031 [P] [US2] Write unit test for disk space exhausted error in internal/config/config_test.go
- [X] T032 [P] [US2] Write unit test for path is directory error in internal/config/config_test.go
- [X] T033 [P] [US2] Write unit test for home directory resolution error in internal/config/config_test.go
- [X] T034 [P] [US2] Write unit test for error message clarity and context in internal/config/config_test.go

### Implementation for User Story 2

- [X] T035 [US2] Add error handling for directory creation failures with clear error messages in internal/config/config.go
- [X] T036 [US2] Add error handling for file creation failures with clear error messages in internal/config/config.go
- [X] T037 [US2] Add error handling for permission setting failures with clear error messages in internal/config/config.go
- [X] T038 [US2] Add error handling for path validation (directory check) with clear error messages in internal/config/config.go
- [X] T039 [US2] Ensure all errors are wrapped with context using fmt.Errorf() in internal/config/config.go

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently - file creation works, errors are handled gracefully with clear messages

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [X] T040 [P] Update README.md with config file auto-creation documentation
- [X] T041 [P] Update CHANGELOG.md with config file auto-creation feature entry
- [X] T042 Run all unit tests to verify no regressions
- [X] T043 Run all integration tests to verify no regressions
- [X] T044 [P] Verify error messages are clear and actionable (code review/audit)
- [X] T045 [P] Run quickstart.md validation scenarios
- [X] T046 Code cleanup and formatting (gofmt, goimports)
- [X] T047 Run golangci-lint to verify code quality
- [X] T048 Verify file permissions are correct (0600 for file, 0755 for directories)
- [X] T049 Verify logging only occurs when file is actually created (not when it exists)

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
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - Can be implemented independently but tests error scenarios from US1

### Within Each User Story

- Tests (if included) MUST be written and FAIL before implementation
- File existence check before file creation
- Directory creation before file creation
- File creation before permission setting
- Permission setting before logging
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
Task: "Write integration test for file creation when missing in test/integration/config_test.go"
Task: "Write integration test for existing file not modified in test/integration/config_test.go"
Task: "Write integration test for parent directory creation in test/integration/config_test.go"
Task: "Write integration test for custom config path in test/integration/config_test.go"
Task: "Write integration test for default config path (~/.gitcomm/config.yaml) in test/integration/config_test.go"
```

---

## Parallel Example: User Story 2

```bash
# Launch all tests for User Story 2 together:
Task: "Write unit test for read-only directory error in internal/config/config_test.go"
Task: "Write unit test for disk space exhausted error in internal/config/config_test.go"
Task: "Write unit test for path is directory error in internal/config/config_test.go"
Task: "Write unit test for home directory resolution error in internal/config/config_test.go"
Task: "Write unit test for error message clarity and context in internal/config/config_test.go"
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
   - Developer A: User Story 1 (file creation logic)
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
- Total tasks: 49
- Tasks per story: US1 (11 tasks), US2 (10 tasks), Setup (5 tasks), Foundational (13 tasks), Polish (10 tasks)
