# Tasks: Replace Go-Git Library with External Git Commands

**Input**: Design documents from `/specs/016-replace-gogit/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: Tests are OPTIONAL. Existing tests use git CLI for setup and should work with the new implementation. No new test tasks are included unless explicitly needed for new error types.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [X] T001 Verify git 2.34.0+ is available in PATH (prerequisite check)
- [X] T002 Backup current `internal/repository/git_repository_impl.go` (safety measure)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T003 Create new error types file `internal/repository/git_errors.go` with ErrGitNotFound, ErrGitVersionTooOld, ErrGitPermissionDenied, ErrGitSigningFailed, ErrGitFileNotFound, and ErrGitCommandFailed struct
- [X] T004 Update `pkg/git/config/extractor.go` to modify CommitSigner struct: remove `Signer interface{}` field, add `Enabled bool` field
- [X] T005 Update `pkg/git/config/extractor.go` function `prepareCommitSigner` to set `Enabled` flag instead of creating go-git Signer interface

**Checkpoint**: Foundation ready - error types and config changes complete. User story implementation can now begin.

---

## Phase 3: User Story 1 - Document Current Git Operations (Priority: P1) 🎯 MVP

**Goal**: Document all go-git API calls with their purposes and contexts (already completed in research.md)

**Independent Test**: Review research.md to verify all go-git API calls are documented with file location, method name, and git concept

**Status**: ✅ COMPLETE - Documented in `research.md` (RQ1: Go-Git API → Git CLI Command Mapping)

**Checkpoint**: User Story 1 documentation is complete and available for reference

---

## Phase 4: User Story 2 - Map Git Operations to CLI Commands (Priority: P1)

**Goal**: Create mapping document showing git CLI commands that replace each go-git operation (already completed in research.md)

**Independent Test**: Review research.md to verify each go-git operation has equivalent git CLI command with flags and arguments

**Status**: ✅ COMPLETE - Documented in `research.md` (RQ1: Go-Git API → Git CLI Command Mapping table)

**Checkpoint**: User Story 2 mapping is complete and available for implementation reference

---

## Phase 5: User Story 3 - Replace Go-Git Implementation with CLI Execution (Priority: P2)

**Goal**: Rewrite GitRepository implementation to use external git commands instead of go-git library while maintaining the same interface and behavior

**Independent Test**: Run existing tests (`go test ./internal/repository/... -v`) and verify all tests pass with identical behavior to go-git implementation

### Implementation for User Story 3

- [X] T006 [US3] Rewrite `gitRepositoryImpl` struct in `internal/repository/git_repository_impl.go`: remove `repo *git.Repository` field, add `gitBin string` field
- [X] T007 [US3] Implement `validateGitVersion()` helper function in `internal/repository/git_repository_impl.go` that runs `git --version` and validates >= 2.34.0, returns ErrGitVersionTooOld if invalid
- [X] T008 [US3] Update `NewGitRepository()` function in `internal/repository/git_repository_impl.go`: add git executable lookup via `exec.LookPath("git")`, add git version validation, remove `git.PlainOpen()` call
- [X] T009 [US3] Implement `execGit()` helper function in `internal/repository/git_repository_impl.go` that executes `git -C <path> <args...>` with context cancellation, captures stdout/stderr, logs execution (FR-018), and categorizes errors (FR-006)
- [X] T010 [US3] Implement `categorizeError()` helper function in `internal/repository/git_repository_impl.go` that parses stderr and maps to appropriate error types (ErrGitPermissionDenied, ErrGitSigningFailed, ErrGitFileNotFound, ErrGitCommandFailed, utils.ErrNotGitRepository)
- [X] T011 [US3] Implement `parseStatus()` helper function in `internal/repository/git_repository_impl.go` that parses `git status --porcelain=v1` output and returns staged/unstaged FileChange arrays
- [X] T012 [US3] Implement `parseDiff()` helper function in `internal/repository/git_repository_impl.go` that parses `git diff --cached --unified=0` output, splits on `diff --git` boundaries, detects binary files, and returns map[filepath]string
- [X] T013 [US3] Rewrite `GetRepositoryState()` method in `internal/repository/git_repository_impl.go`: use `execGit("status", "--porcelain=v1")` and `execGit("diff", "--cached", "--unified=0")`, call parseStatus and parseDiff helpers, apply size limits and includeNewFiles filtering
- [X] T014 [US3] Rewrite `CaptureStagingState()` method in `internal/repository/git_repository_impl.go`: use `execGit("status", "--porcelain=v1")` and parseStatus helper to extract staged files
- [X] T015 [US3] Rewrite `StageAllFiles()` method in `internal/repository/git_repository_impl.go`: use `execGit("add", "-A")`
- [X] T016 [US3] Rewrite `StageModifiedFiles()` method in `internal/repository/git_repository_impl.go`: use `execGit("status", "--porcelain=v1")` to filter modified files (not untracked), loop with `execGit("add", "--", filepath)`, rollback with `execGit("reset", "HEAD", "--", staged...)` on failure
- [X] T017 [US3] Rewrite `StageAllFilesIncludingUntracked()` method in `internal/repository/git_repository_impl.go`: use `execGit("status", "--porcelain=v1")` to filter all changed files, loop with `execGit("add", "--", filepath)`, rollback with `execGit("reset", "HEAD", "--", staged...)` on failure
- [X] T018 [US3] Rewrite `UnstageFiles()` method in `internal/repository/git_repository_impl.go`: use `execGit("reset", "HEAD", "--", files...)` (remove go-git fallback)
- [X] T019 [US3] Rewrite `CreateCommit()` method in `internal/repository/git_repository_impl.go`: format commit message, set GIT_AUTHOR_NAME/EMAIL env vars on exec.Cmd, if signer.Enabled add `-c gpg.format=ssh -c user.signingkey=<key> -c commit.gpgsign=true` flags and `-S` flag, retry without `-S` on signing failure
- [X] T020 [US3] Remove all go-git imports from `internal/repository/git_repository_impl.go`: remove `github.com/go-git/go-git/v5`, `github.com/go-git/go-git/v5/plumbing`, `github.com/go-git/go-git/v5/plumbing/object`, `github.com/go-git/go-git/v5/utils/diff`, `github.com/hiddeco/sshsig`, `github.com/sergi/go-diff/diffmatchpatch`, `golang.org/x/crypto/ssh`
- [X] T021 [US3] Remove unused helper functions from `internal/repository/git_repository_impl.go`: remove `getHEADTree()`, `getStagedIndexTree()`, `formatUnifiedDiff()`, `formatNewFileDiff()`, `formatDeletedFileDiff()`, `formatModifiedFileDiff()`, `formatRenameDiff()`, `formatCopyDiff()`, `isBinaryFile()` (replaced by git CLI), `generateMetadata()` (keep if still used for size limits), `applySizeLimit()` (keep), `statusCodeToString()` (keep), `sshCommitSignerWrapper` type and methods
- [X] T022 [US3] Update `internal/repository/git_repository_impl_test.go`: remove any go-git test imports if present, verify tests still use git CLI for setup (they should already)

**Checkpoint**: At this point, User Story 3 should be fully functional. All interface methods work with git CLI. Run tests to verify.

---

## Phase 6: Dependency Cleanup

**Purpose**: Remove go-git and related dependencies from go.mod

- [X] T023 Remove `github.com/go-git/go-git/v5` dependency from `go.mod` using `go get github.com/go-git/go-git/v5@none`
- [X] T024 Remove `github.com/hiddeco/sshsig` dependency from `go.mod` using `go get github.com/hiddeco/sshsig@none`
- [X] T025 Remove `github.com/sergi/go-diff` dependency from `go.mod` using `go get github.com/sergi/go-diff@none`
- [X] T026 Run `go mod tidy` to clean up transitive dependencies (go-billy, go-git/gcfg v1 indirect, etc.)
- [X] T027 Verify `github.com/go-git/gcfg/v2` is retained in `go.mod` (required for git config parsing)

**Checkpoint**: Dependencies cleaned up. Verify build still works: `go build ./...`

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Final validation and cleanup

- [X] T028 Run all tests: `go test ./internal/repository/... -v -count=1`
- [X] T029 Run full test suite: `go test ./... -v -count=1`
- [X] T030 Verify no go-git imports remain: `grep -r "go-git" internal/ pkg/` (should only find gcfg/v2)
- [X] T031 Verify git CLI logging works: run gitcomm with debug logging and verify git command executions are logged (FR-018)
- [X] T032 Verify error categorization: test with missing git, old git version, permission errors, signing failures
- [X] T033 Update CHANGELOG.md with migration notes
- [X] T034 Run quickstart.md validation: verify all implementation steps were followed

**Checkpoint**: Migration complete. All tests pass, dependencies removed, logging and error handling verified.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Story 1 (Phase 3)**: ✅ Already complete (documentation in research.md)
- **User Story 2 (Phase 4)**: ✅ Already complete (mapping in research.md)
- **User Story 3 (Phase 5)**: Depends on Foundational (Phase 2) completion
- **Dependency Cleanup (Phase 6)**: Depends on User Story 3 completion
- **Polish (Phase 7)**: Depends on Dependency Cleanup completion

### User Story Dependencies

- **User Story 1 (P1)**: ✅ Complete - Documentation phase
- **User Story 2 (P1)**: ✅ Complete - Mapping phase
- **User Story 3 (P2)**: Depends on Foundational phase (error types, config changes)

### Within User Story 3

- Error types (T003) → Config changes (T004, T005) → Struct rewrite (T006) → Validation (T007, T008) → Core helpers (T009, T010, T011, T012) → Interface methods (T013-T019) → Cleanup (T020, T021, T022)

### Parallel Opportunities

- T003, T004, T005 can run in parallel (different files: git_errors.go, extractor.go)
- T011 and T012 can run in parallel (different helper functions)
- T013-T019 can run sequentially (each implements one interface method, but they share helpers)
- T023-T025 can run sequentially (each removes one dependency)
- T028-T032 can run in parallel (different validation tasks)

---

## Parallel Example: Foundational Phase

```bash
# Launch foundational tasks in parallel:
Task: "Create new error types file internal/repository/git_errors.go" (T003)
Task: "Update pkg/git/config/extractor.go CommitSigner struct" (T004)
Task: "Update prepareCommitSigner function" (T005)
```

---

## Parallel Example: User Story 3 Helpers

```bash
# After T006-T010 complete, these helpers can be developed/tested independently:
Task: "Implement parseStatus() helper" (T011)
Task: "Implement parseDiff() helper" (T012)
```

---

## Implementation Strategy

### MVP First (User Stories 1 & 2 Already Complete)

1. ✅ Phase 3: User Story 1 - Documentation (COMPLETE)
2. ✅ Phase 4: User Story 2 - Mapping (COMPLETE)
3. Complete Phase 1: Setup
4. Complete Phase 2: Foundational (CRITICAL - blocks US3)
5. Complete Phase 5: User Story 3 (main implementation)
6. **STOP and VALIDATE**: Run tests, verify all interface methods work
7. Complete Phase 6: Dependency Cleanup
8. Complete Phase 7: Polish & Validation

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Implement core helpers (execGit, parseStatus, parseDiff) → Test helpers independently
3. Implement simple methods first (StageAllFiles, UnstageFiles) → Test independently
4. Implement complex methods (GetRepositoryState, CreateCommit) → Test independently
5. Implement staging methods (StageModifiedFiles, StageAllFilesIncludingUntracked) → Test independently
6. Clean up dependencies → Verify build
7. Final validation → Deploy

### Critical Path

**Sequential dependencies** (cannot parallelize):
- T003 → T006 → T007 → T008 → T009 → T010 → T011 → T012 → T013-T019 → T020-T022 → T023-T027 → T028-T034

**Parallel opportunities**:
- T003, T004, T005 (foundational, different files)
- T011, T012 (helpers, independent)
- T028-T032 (validation, independent checks)

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- User Stories 1 & 2 are documentation/mapping tasks already complete
- User Story 3 is the main implementation work
- Existing tests should work without modification (they use git CLI for setup)
- Verify each interface method works before moving to the next
- Commit after each logical group (foundational, helpers, methods, cleanup)
- Stop at checkpoints to validate independently
- Avoid: skipping error categorization, forgetting logging, missing size limits, breaking interface contract
