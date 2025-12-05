# Tasks: Improve Commit Generation with Git Config

**Input**: Design documents from `/specs/012-git-config-commit/`
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

- [X] T001 Create directory structure for `pkg/git/config/` per implementation plan
- [X] T002 [P] Create error types in `pkg/git/config/errors.go` for config extraction errors

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T003 Define `ConfigExtractor` interface in `pkg/git/config/extractor.go` with `Extract(repoPath string) *GitConfig` method
- [X] T004 Define `GitConfig` struct in `pkg/git/config/extractor.go` with fields: `UserName`, `UserEmail`, `SigningKey`, `GPGFormat`, `CommitGPGSign`
- [X] T005 [P] Write unit test for `ConfigExtractor.Extract()` with local config only in `pkg/git/config/extractor_test.go` - test should fail initially
- [X] T006 [P] Write unit test for `ConfigExtractor.Extract()` with global config only in `pkg/git/config/extractor_test.go` - test should fail initially
- [X] T007 [P] Write unit test for `ConfigExtractor.Extract()` with local taking precedence in `pkg/git/config/extractor_test.go` - test should fail initially
- [X] T008 [P] Write unit test for `ConfigExtractor.Extract()` with missing files (uses defaults) in `pkg/git/config/extractor_test.go` - test should fail initially
- [X] T009 [P] Write unit test for `ConfigExtractor.Extract()` with unreadable/corrupted files in `pkg/git/config/extractor_test.go` - test should fail initially
- [X] T010 [P] Write unit test for `ConfigExtractor.Extract()` with partial values (user.name but no user.email) in `pkg/git/config/extractor_test.go` - test should fail initially
- [X] T011 [P] Write unit test for `ConfigExtractor.Extract()` performance (<50ms) in `pkg/git/config/extractor_test.go` - test should fail initially
- [X] T012 Implement `FileConfigExtractor` struct in `pkg/git/config/extractor.go` implementing `ConfigExtractor` interface
- [X] T013 Implement `NewFileConfigExtractor()` constructor function in `pkg/git/config/extractor.go` returning `ConfigExtractor`
- [X] T014 Implement `Extract()` method in `FileConfigExtractor` in `pkg/git/config/extractor.go` - read and parse `.git/config` using gcfg
- [X] T015 Implement `Extract()` method fallback to `~/.gitconfig` in `FileConfigExtractor` in `pkg/git/config/extractor.go` - expand `~` to home directory
- [X] T016 Implement local config precedence logic in `Extract()` method in `pkg/git/config/extractor.go` - local values override global
- [X] T017 Implement default values logic in `Extract()` method in `pkg/git/config/extractor.go` - use "gitcomm" and "gitcomm@local" for missing user.name/user.email
- [X] T018 Implement silent error handling in `Extract()` method in `pkg/git/config/extractor.go` - log debug messages, never return errors
- [X] T019 Implement debug logging for missing/unreadable config files in `Extract()` method in `pkg/git/config/extractor.go`

**Checkpoint**: Foundation ready - ConfigExtractor complete, user story implementation can now begin

---

## Phase 3: User Story 1 - Git Config Extraction and Commit Author Configuration (Priority: P1) 🎯 MVP

**Goal**: Commits automatically use user.name and user.email from git configuration, ensuring proper commit attribution without manual configuration.

**Independent Test**: Configure git with user.name and user.email, run gitcomm to create a commit, and verify the commit author matches the git configuration.

### Tests for User Story 1 ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [X] T020 [P] [US1] Write integration test for commit author from local config in `test/integration/git_config_test.go` - test should fail initially
- [X] T021 [P] [US1] Write integration test for commit author from global config in `test/integration/git_config_test.go` - test should fail initially
- [X] T022 [P] [US1] Write integration test for local config precedence over global in `test/integration/git_config_test.go` - test should fail initially
- [X] T023 [P] [US1] Write integration test for default author values when config missing in `test/integration/git_config_test.go` - test should fail initially
- [X] T024 [P] [US1] Write unit test for `NewGitRepository()` extracts config before opening repository in `internal/repository/git_repository_impl_test.go` - test should fail initially
- [X] T025 [P] [US1] Write unit test for `CreateCommit()` uses extracted user.name and user.email for author in `internal/repository/git_repository_impl_test.go` - test should fail initially

### Implementation for User Story 1

- [X] T026 [US1] Update `gitRepositoryImpl` struct in `internal/repository/git_repository_impl.go` to add `config *config.GitConfig` field
- [X] T027 [US1] Update `NewGitRepository()` function in `internal/repository/git_repository_impl.go` to call `ConfigExtractor.Extract()` BEFORE `git.PlainOpen()`
- [X] T028 [US1] Update `NewGitRepository()` function in `internal/repository/git_repository_impl.go` to store extracted config in repository struct
- [X] T029 [US1] Update `CreateCommit()` method in `internal/repository/git_repository_impl.go` to use `r.config.UserName` for commit author name
- [X] T030 [US1] Update `CreateCommit()` method in `internal/repository/git_repository_impl.go` to use `r.config.UserEmail` for commit author email
- [X] T031 [US1] Update `CreateCommit()` method in `internal/repository/git_repository_impl.go` to use defaults ("gitcomm", "gitcomm@local") when config values are empty

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently. Commits use git config for author attribution.

---

## Phase 4: User Story 2 - Commit Signing with Git Config (Priority: P2)

**Goal**: Commits are automatically signed using SSH signing configuration from git config when available, ensuring commit authenticity and integrity.

**Independent Test**: Configure git with SSH signing settings (`gpg.format = ssh` and `user.signingkey` pointing to SSH public key file), run gitcomm to create a commit, and verify the commit is signed with the configured SSH key.

### Tests for User Story 2 ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T032 [P] [US2] Write unit test for `prepareCommitSigner()` with SSH signing configured in `internal/repository/git_repository_impl_test.go` - test should fail initially
- [ ] T033 [P] [US2] Write unit test for `prepareCommitSigner()` with commit.gpgsign=false in `internal/repository/git_repository_impl_test.go` - test should fail initially
- [ ] T034 [P] [US2] Write unit test for `prepareCommitSigner()` with missing gpg.format in `internal/repository/git_repository_impl_test.go` - test should fail initially
- [ ] T035 [P] [US2] Write unit test for `CreateCommit()` signs commit when SSH signing configured in `internal/repository/git_repository_impl_test.go` - test should fail initially
- [ ] T036 [P] [US2] Write unit test for `CreateCommit()` handles signing failure gracefully in `internal/repository/git_repository_impl_test.go` - test should fail initially
- [ ] T037 [P] [US2] Write integration test for SSH signed commit creation in `test/integration/git_config_test.go` - test should fail initially
- [ ] T038 [P] [US2] Write integration test for --no-sign flag disables signing in `test/integration/git_config_test.go` - test should fail initially

### Implementation for User Story 2

- [X] T039 [US2] Define `CommitSigner` struct in `pkg/git/config/extractor.go` with fields: `PrivateKeyPath`, `PublicKeyPath`, `Format`, `Enabled`, `Signer`
- [X] T040 [US2] Implement `prepareCommitSigner()` helper function in `internal/repository/git_repository_impl.go` to create CommitSigner from GitConfig
- [X] T041 [US2] Implement SSH signing configuration check in `prepareCommitSigner()` in `internal/repository/git_repository_impl.go` - verify `gpg.format == "ssh"` and `user.signingkey != ""`
- [X] T042 [US2] Implement private key path derivation in `prepareCommitSigner()` in `internal/repository/git_repository_impl.go` - remove `.pub` from `user.signingkey` to get private key path
- [X] T043 [US2] Implement SSH signer creation in `prepareCommitSigner()` in `internal/repository/git_repository_impl.go` - load private key and create go-git SSH signer
- [X] T044 [US2] Update `gitRepositoryImpl` struct in `internal/repository/git_repository_impl.go` to add `signer *config.CommitSigner` field
- [X] T045 [US2] Update `NewGitRepository()` function in `internal/repository/git_repository_impl.go` to call `prepareCommitSigner()` and store signer
- [X] T046 [US2] Update `CreateCommit()` method in `internal/repository/git_repository_impl.go` to add `SignKey` to `CommitOptions` when signer is enabled
- [X] T047 [US2] Implement signing failure handling in `CreateCommit()` method in `internal/repository/git_repository_impl.go` - retry without signing, log error, proceed
- [X] T048 [US2] Add `--no-sign` CLI flag definition in `cmd/gitcomm/main.go` - boolean flag with default false
- [X] T049 [US2] Pass `--no-sign` flag value to `CommitService` in `cmd/gitcomm/main.go` - add to `CommitOptions` or service constructor
- [X] T050 [US2] Update `CommitService` struct in `internal/service/commit_service.go` to store `noSign` flag value
- [X] T051 [US2] Update `prepareCommitSigner()` to respect `--no-sign` flag in `internal/repository/git_repository_impl.go` - disable signing if flag set
- [X] T052 [US2] Update `prepareCommitSigner()` to respect `commit.gpgsign = false` in `internal/repository/git_repository_impl.go` - disable signing if config opt-out

**Checkpoint**: At this point, User Story 2 should be fully functional and testable independently. Commits are signed with SSH keys when configured.

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [X] T053 [P] Update `README.md` to document git config extraction and SSH signing features
- [X] T054 [P] Update `CHANGELOG.md` with feature description and breaking changes (none)
- [X] T055 [P] Verify all existing repository tests still pass after config extraction changes
- [X] T056 [P] Run quickstart.md validation to ensure implementation matches design
- [X] T057 [P] Code cleanup: remove any unused imports or dead code from modified files
- [X] T058 [P] Add code comments documenting config extraction and signing usage in repository implementation

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Foundational phase completion
- **User Story 2 (Phase 4)**: Depends on Foundational phase completion, can use User Story 1 components
- **Polish (Phase 5)**: Depends on User Story 1 and User Story 2 completion

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - Uses ConfigExtractor from Foundational, can integrate with User Story 1

### Within User Story 1

- Tests (T020-T025) MUST be written and FAIL before implementation
- ConfigExtractor interface (T003-T004) before implementation (T012-T019)
- ConfigExtractor implementation (T012-T019) before repository integration (T026-T031)
- Repository integration (T026-T031) before integration tests

### Within User Story 2

- Tests (T032-T038) MUST be written and FAIL before implementation
- CommitSigner struct (T039) before prepareCommitSigner (T040-T043)
- prepareCommitSigner implementation (T040-T043) before repository integration (T044-T047)
- CLI flag (T048-T049) before service integration (T050-T052)
- Core signing implementation complete before integration tests

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational test tasks marked [P] can run in parallel (within Phase 2)
- All User Story 1 test tasks marked [P] can run in parallel
- All User Story 2 test tasks marked [P] can run in parallel
- Repository struct updates (T026, T044) can run in parallel (different fields)
- CLI flag tasks (T048-T049) can run in parallel with repository signing tasks

---

## Parallel Example: User Story 1

```bash
# Launch all tests for User Story 1 together:
Task: "Write integration test for commit author from local config"
Task: "Write integration test for commit author from global config"
Task: "Write integration test for local config precedence over global"
Task: "Write integration test for default author values when config missing"
Task: "Write unit test for NewGitRepository() extracts config before opening repository"
Task: "Write unit test for CreateCommit() uses extracted user.name and user.email for author"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Test User Story 1 independently
   - Verify commits use git config for author
   - Verify local config takes precedence
   - Verify defaults used when config missing
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
   - Developer A: User Story 1 (author configuration)
   - Developer B: User Story 2 (SSH signing) - can start in parallel
3. Stories complete and integrate independently

---

## Notes

- [P] tasks = different files, no dependencies
- [US1] label maps task to User Story 1 for traceability
- [US2] label maps task to User Story 2 for traceability
- User Story 1 should be independently completable and testable
- User Story 2 should be independently completable and testable
- Verify tests fail before implementing (TDD approach)
- Commit after each task or logical group
- Stop at checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence
- All repository updates maintain backward compatibility (GitRepository interface unchanged)
