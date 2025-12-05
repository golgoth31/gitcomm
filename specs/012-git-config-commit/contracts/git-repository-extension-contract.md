# Contract: GitRepository Interface Extension

**Feature**: 012-git-config-commit
**Date**: 2025-01-27
**Component**: `internal/repository`

## Interface Changes

### Existing Interface

```go
package repository

type GitRepository interface {
    GetRepositoryState(ctx context.Context) (*model.RepositoryState, error)
    CreateCommit(ctx context.Context, message *model.CommitMessage) error
    StageAllFiles(ctx context.Context) error
    CaptureStagingState(ctx context.Context) (*model.StagingState, error)
    StageModifiedFiles(ctx context.Context) (*model.AutoStagingResult, error)
    StageAllFilesIncludingUntracked(ctx context.Context) (*model.AutoStagingResult, error)
    UnstageFiles(ctx context.Context, files []string) error
}
```

### No Interface Changes Required

**Rationale**: The `GitRepository` interface remains unchanged. Config extraction and signing are internal implementation details of `gitRepositoryImpl`.

---

## Implementation Changes

### gitRepositoryImpl Struct

**New Fields**:
```go
type gitRepositoryImpl struct {
    repo     *git.Repository
    path     string
    config   *config.GitConfig  // NEW: Extracted git config
    signer   *config.CommitSigner // NEW: Prepared commit signer (if applicable)
}
```

### NewGitRepository Function

**Changes**:
- Extract git config BEFORE opening repository with `git.PlainOpen()`
- Store extracted config in repository struct
- Prepare commit signer if SSH signing is configured

**New Implementation Flow**:
```go
func NewGitRepository(repoPath string) (GitRepository, error) {
    // 1. Find repository root (existing logic)

    // 2. Extract git config BEFORE opening repository
    extractor := config.NewFileConfigExtractor()
    gitConfig := extractor.Extract(path)

    // 3. Prepare commit signer if SSH signing configured
    signer := prepareCommitSigner(gitConfig)

    // 4. Open repository (existing logic)
    repo, err := git.PlainOpen(path)

    // 5. Return repository with config and signer
    return &gitRepositoryImpl{
        repo:   repo,
        path:   path,
        config: gitConfig,
        signer: signer,
    }, nil
}
```

### CreateCommit Method

**Changes**:
- Use `gitConfig.UserName` and `gitConfig.UserEmail` for commit author
- Use `signer` in `CommitOptions.SignKey` if signing is enabled
- Handle signing failures gracefully (create unsigned commit)

**New Implementation**:
```go
func (r *gitRepositoryImpl) CreateCommit(ctx context.Context, message *model.CommitMessage) error {
    // ... existing formatting logic ...

    // Use extracted config for author
    author := &object.Signature{
        Name:  r.config.UserName,
        Email: r.config.UserEmail,
        When:  time.Now(),
    }

    // Prepare commit options
    opts := &git.CommitOptions{
        Author: author,
    }

    // Add signer if configured and enabled
    if r.signer != nil && r.signer.Enabled {
        opts.SignKey = r.signer.Signer // SSH signer from go-git
    }

    // Create commit (signing failures handled by go-git, we catch and log)
    _, err := worktree.Commit(commitMsg, opts)
    if err != nil {
        // Check if error is signing-related
        if isSigningError(err) {
            // Retry without signing (per FR-013)
            opts.SignKey = nil
            _, err = worktree.Commit(commitMsg, opts)
            if err != nil {
                return fmt.Errorf("failed to create commit: %w", err)
            }
            utils.Logger.Debug().Err(err).Msg("SSH signing failed, created unsigned commit")
            return nil
        }
        return fmt.Errorf("failed to create commit: %w", err)
    }

    return nil
}
```

---

## Behavior Specification

### Config Extraction

**Timing**: MUST happen before `git.PlainOpen()` (FR-001, FR-002)

**Precedence**: Local config (`.git/config`) takes precedence over global (`~/.gitconfig`) (FR-005)

**Error Handling**: Silent ignore missing/unreadable files, use defaults (FR-009, FR-010)

### Commit Author

**Source**: `gitConfig.UserName` and `gitConfig.UserEmail` (FR-003, FR-004)

**Defaults**: "gitcomm" and "gitcomm@local" if config values missing (FR-012)

### Commit Signing

**Configuration Check**:
- `gpg.format == "ssh"` (FR-006)
- `user.signingkey != ""` (FR-006)
- `commit.gpgsign != false` (FR-007)
- `--no-sign` flag not set (FR-008)

**Signing Process**:
- Derive private key path from `user.signingkey` (remove `.pub`)
- Load private key and create SSH signer
- Pass signer to `CommitOptions.SignKey`
- Handle signing failures gracefully (FR-013)

**Failure Handling**:
- If signing fails: Create unsigned commit, log error, proceed (FR-013)
- Do not fail commit creation due to signing errors

---

## Test Requirements

### Unit Tests

**File**: `internal/repository/git_repository_impl_test.go`

**New Test Cases**:
1. NewGitRepository extracts config before opening repository
2. CreateCommit uses extracted user.name and user.email for author
3. CreateCommit signs commit when SSH signing configured
4. CreateCommit does not sign when commit.gpgsign = false
5. CreateCommit does not sign when --no-sign flag set
6. CreateCommit handles signing failure gracefully (creates unsigned commit)
7. CreateCommit uses local config values over global config
8. CreateCommit uses defaults when config files missing

### Integration Tests

**File**: `test/integration/git_config_test.go`

**Test Cases**:
1. End-to-end: Extract config, create signed commit, verify signature
2. End-to-end: Extract config, create commit with correct author
3. Config precedence: Local config overrides global config
4. Signing failure: SSH key missing, commit still created (unsigned)

---

## Breaking Changes

**None**: Interface unchanged, backward compatible implementation.

---

## Migration Notes

**N/A**: Internal implementation change, no API changes.
