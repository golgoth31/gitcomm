# Tasks: Explicit Commit Error Messages

**Input**: Design documents from `/specs/017-explicit-commit-errors/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Unit tests for FormatErrorForDisplay and ErrGitCommandFailed are included per contract requirements.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Verify environment and branch

- [X] T001 Verify branch `017-explicit-commit-errors` and run `go build ./...` from repo root

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core error formatting infrastructure that MUST be complete before user stories

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T002 [P] Update `ErrGitCommandFailed.Error()` in `internal/repository/git_errors.go`: when Stderr is empty or whitespace, use generic hint "No additional details from git. Check repository state or run the command manually." (add `strings` import)
- [X] T003 Create `FormatErrorForDisplay(err error) string` in `internal/repository/error_display.go`: detect ErrGitCommandFailed via `errors.As`, truncate Stderr at 1500 chars with "… (N additional characters)" suffix, format as "git [Command] failed (exit [ExitCode]). Details: [content]" or generic hint when empty; fallback to `err.Error()` for other errors
- [X] T004 [P] Add unit test for `ErrGitCommandFailed.Error()` with empty Stderr in `internal/repository/git_errors_test.go` (create file if missing)
- [X] T005 [P] Add unit tests for `FormatErrorForDisplay` in `internal/repository/error_display_test.go`: ErrGitCommandFailed (empty stderr), non-empty stderr, stderr >1500 chars, ErrGitSigningFailed, generic error, nil handling

**Checkpoint**: Foundation ready — `FormatErrorForDisplay` and `ErrGitCommandFailed` fix complete. Run `go test ./internal/repository/... -v` to verify.

---

## Phase 3: User Story 1 - See Explicit Error Details When Commit Fails (Priority: P1) 🎯 MVP

**Goal**: When a commit fails, the user sees explicit error details (git stderr or generic hint) with the "Brief. Details:" format.

**Independent Test**: Trigger a commit failure (e.g., failing pre-commit hook) and verify the displayed message contains "Details:" and actionable information (or generic hint when stderr empty).

- [X] T006 [P] [US1] Update `internal/cmd/root.go` line 158: replace `%v` with `repository.FormatErrorForDisplay(commitErr)` in commit failure display
- [X] T007 [P] [US1] Update `internal/service/commit_service.go` line 556: replace `%v` with `repository.FormatErrorForDisplay(commitErr)` in `handleCommitFailure`

**Checkpoint**: Commit failures display explicit errors. Test by running gitcomm in a repo with a failing pre-commit hook.

---

## Phase 4: User Story 2 - Consistent Explicit Errors for All Git Operations (Priority: P2)

**Goal**: Repository init failures and any other git operation failures that reach the user show explicit error details.

**Independent Test**: Run gitcomm outside a git repo or with invalid repo path; verify "Error: failed to initialize git repository:" is followed by explicit details.

- [X] T008 [US2] Update `internal/cmd/root.go` line 81: replace `%v` with `repository.FormatErrorForDisplay(err)` in repository init failure display

**Checkpoint**: All three display points (repo init, commit workflow failure, AcceptAndCommit retry) use explicit formatting. Staging/unstage errors during workflow flow through root.go:158 and get the formatter via T006.

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Final validation and documentation

- [X] T009 Run `go test ./internal/repository/... -v -count=1`
- [X] T010 Run `go test ./... -count=1` (full suite)
- [ ] T011 Manual verification: create repo with failing pre-commit hook, run gitcomm, confirm output contains "Details:" and hook stderr (user to perform)
- [X] T012 Update `CHANGELOG.md` with explicit error display feature
- [X] T013 Run quickstart.md validation: verify all implementation steps were followed

**Checkpoint**: Migration complete. All tests pass, manual verification done, CHANGELOG updated.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — can start immediately
- **Foundational (Phase 2)**: Depends on Setup — BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Foundational
- **User Story 2 (Phase 4)**: Depends on Foundational (can run in parallel with US1 or after)
- **Polish (Phase 5)**: Depends on US1 and US2 completion

### User Story Dependencies

- **User Story 1 (P1)**: T006, T007 — No dependency on US2
- **User Story 2 (P2)**: T008 — No dependency on US1 (both use FormatErrorForDisplay from Foundational)

### Within Phases

- T002 and T003 can run in parallel (different files; T002 fixes Error(), T003 creates formatter)
- T004 depends on T002 (tests the ErrGitCommandFailed fix)
- T005 depends on T003 (tests FormatErrorForDisplay)
- T006 and T007 can run in parallel (different files: root.go vs commit_service.go)
- T008 touches root.go like T006 — can run after T006 or be combined in one edit; kept separate for US2 traceability

### Parallel Opportunities

- T002, T004, T005 (Foundational)
- T006, T007 (US1)

---

## Parallel Example: Foundational Phase

```bash
# Launch in parallel:
Task T002: "Update ErrGitCommandFailed.Error() in git_errors.go"
Task T004: "Add unit test for ErrGitCommandFailed in git_errors_test.go"
Task T005: "Add unit tests for FormatErrorForDisplay in error_display_test.go"
# Then T003 (Create FormatErrorForDisplay) - requires ErrGitCommandFailed fix for consistency
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (T002 → T003, then T004, T005)
3. Complete Phase 3: User Story 1 (T006, T007)
4. **STOP and VALIDATE**: Trigger commit failure, verify explicit error display
5. Proceed to US2 or deploy

### Incremental Delivery

1. Setup + Foundational → FormatErrorForDisplay ready
2. Add US1 (commit display points) → Test with failing hook
3. Add US2 (repo init display) → Test outside git repo
4. Polish → CHANGELOG, quickstart validation

### Critical Path

- T001 → T002 → T003 → T006 → T007 → T009
- T004, T005 can run alongside T002
- T008 after T006 (same file) or independently after Foundational

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to traceability
- FormatErrorForDisplay handles nil by returning ""; callers must not pass nil (document in function)
- Truncation constant: 1500 characters
- Commit failures and staging failures during workflow both surface at root.go:158 — T006 covers both
