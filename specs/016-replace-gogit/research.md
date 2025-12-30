# Research: Replace Go-Git Library with External Git Commands

**Feature**: 016-replace-gogit
**Date**: 2026-02-10
**Purpose**: Resolve technical unknowns and document design decisions for migrating from go-git library to git CLI execution

## Research Questions

### RQ1: Go-Git API → Git CLI Command Mapping

**Question**: What are the exact git CLI commands that replace each go-git API call?

**Research**:

The codebase uses 10 distinct go-git APIs across 20 call sites. Here is the complete mapping:

| go-git API | Git CLI Equivalent | Used In | Notes |
|------------|-------------------|---------|-------|
| `git.PlainOpen(path)` | `git -C <path> rev-parse --git-dir` | `NewGitRepository` | Validates repo exists; returns `.git` path |
| `repo.Worktree()` + `worktree.Status()` | `git -C <path> status --porcelain=v1` | `GetRepositoryState`, `CaptureStagingState`, `StageModifiedFiles`, `StageAllFilesIncludingUntracked` | Porcelain v1 format: `XY PATH` where X=staging, Y=worktree |
| `repo.Head()` | `git -C <path> rev-parse HEAD` | `getHEADTree` | Returns commit hash; fails with exit 128 on empty repo |
| `repo.CommitObject(hash)` + `commit.Tree()` | `git -C <path> rev-parse HEAD^{tree}` | `getHEADTree` | Returns tree hash directly |
| `tree.File(path)` + `file.Contents()` | `git -C <path> show HEAD:<filepath>` | `computeFileDiff`, `generateMetadata` | Returns file content at HEAD; fails if file not in HEAD |
| `diff.Do(old, new)` | `git -C <path> diff --cached --unified=0 -- <filepath>` | `formatUnifiedDiff` | Produces unified diff of staged changes with 0 context |
| `worktree.Commit(msg, opts)` | `git -C <path> commit -m "<message>"` | `CreateCommit` | With signing: `git commit -S -m "<message>"` |
| `worktree.AddGlob(".")` | `git -C <path> add -A` | `StageAllFiles` | Stages all changes including untracked |
| `worktree.Add(filepath)` | `git -C <path> add -- <filepath>` | `StageModifiedFiles`, `StageAllFilesIncludingUntracked` | Stages single file |
| `worktree.Remove(filepath)` | `git -C <path> reset HEAD -- <filepath>` | `StageModifiedFiles` (rollback), `StageAllFilesIncludingUntracked` (rollback), `UnstageFiles` (fallback) | Unstages single file |

**Decision**: Use the git CLI commands listed above as direct replacements.

**Rationale**:
- All commands use stable porcelain/plumbing formats
- `--porcelain=v1` ensures stable, parseable output for status
- `--unified=0` matches existing 0-context-line behavior
- `-C <path>` flag avoids need to change working directory

---

### RQ2: Git Status Porcelain Output Parsing

**Question**: How should `git status --porcelain` output be parsed to match go-git's `Status` map?

**Research**:

`git status --porcelain=v1` output format:
```
XY PATH
XY ORIG_PATH -> PATH   (for renames/copies)
```

Where:
- `X` = staging area status code
- `Y` = worktree status code
- Status codes: `M` (modified), `A` (added), `D` (deleted), `R` (renamed), `C` (copied), `U` (unmerged), `?` (untracked), `!` (ignored), ` ` (unmodified)

Mapping to go-git `StatusCode`:
| Porcelain Code | go-git StatusCode | String output |
|---------------|-------------------|---------------|
| `A` | `git.Added` | `"added"` |
| `M` | `git.Modified` | `"modified"` |
| `D` | `git.Deleted` | `"deleted"` |
| `R` | `git.Renamed` | `"renamed"` |
| `C` | `git.Copied` | `"copied"` |
| `U` | `git.UpdatedButUnmerged` | `"unmerged"` |
| `?` | `git.Untracked` | Mapped to `"added"` for worktree |
| ` ` | `git.Unmodified` | `"unmodified"` |

**Decision**: Parse `git status --porcelain=v1` and map status codes to existing string representations.

**Rationale**:
- Porcelain v1 is explicitly designed for machine parsing
- Format is stable across git versions
- Two-character status code maps cleanly to staging/worktree split

---

### RQ3: Diff Computation Strategy

**Question**: Should diffs be computed via `git diff --cached` or by reading file contents and computing internally?

**Research**:

Two approaches:

**Option A**: Use `git diff --cached --unified=0 -- <filepath>`
- Pros: Native git diff, handles all edge cases (binary, rename, copy), already formatted
- Cons: Need to parse output to apply size limit; one process per file or all at once

**Option B**: Read HEAD content via `git show HEAD:<filepath>`, read staged content from filesystem, compute diff internally
- Pros: More control over diff format, can reuse existing `formatUnifiedDiff` logic
- Cons: Reimplements what git does natively; more complex; still needs binary detection

**Option C**: Use `git diff --cached --unified=0` (all files at once) and parse the combined output
- Pros: Single process invocation for all diffs; most efficient
- Cons: Need to split output per file; more complex parsing

**Decision**: Option C - Use `git diff --cached --unified=0` for all staged files at once, parse per-file.

**Rationale**:
- Single process invocation is most efficient (avoids N process spawns for N files)
- Git handles binary detection, rename detection, and all edge cases natively
- Output format is well-defined and parseable (split on `diff --git` boundaries)
- Size limit can still be applied per-file after parsing
- Falls back to metadata generation for files exceeding 5000 chars

**Implementation notes**:
- For added files not yet in HEAD: `git diff --cached --unified=0 -- <newfile>` still works (diffs against empty)
- Binary files show `Binary files differ` in diff output - detect and set empty diff
- Rename detection included in diff output with `similarity index` header

---

### RQ4: SSH Commit Signing via Git CLI

**Question**: How to handle SSH commit signing through git CLI instead of go-git's programmatic signing?

**Research**:

Git 2.34.0+ supports SSH signing natively via configuration:
```
git config gpg.format ssh
git config user.signingkey ~/.ssh/id_ed25519.pub
git config commit.gpgsign true
```

Then `git commit -S -m "message"` signs automatically.

For passphrase-protected keys:
- **SSH Agent**: If ssh-agent is running and key is loaded, signing works transparently
- **GIT_SSH_COMMAND**: Can configure SSH command with key path
- **SSH_ASKPASS**: Environment variable pointing to a program that provides passphrase
- **GIT_COMMITTER_SIGNING_KEY**: Not directly available; git uses user.signingkey config

**Decision**: Delegate signing entirely to git CLI. The system will:
1. Set `GIT_CONFIG_NOSYSTEM=1` if needed
2. Pass signing configuration via `-c` flags: `git -c gpg.format=ssh -c user.signingkey=<key> -c commit.gpgsign=true commit -S -m "msg"`
3. For passphrase: set `SSH_ASKPASS` env var if configured, or rely on ssh-agent
4. Graceful fallback: if signing fails, retry without `-S` flag (same behavior as current go-git impl)

**Rationale**:
- Git CLI handles SSH signing natively and correctly
- Configuration can be passed via `-c` flags without modifying git config files
- SSH agent integration works transparently
- Passphrase handling via `SSH_ASKPASS` is the standard Unix approach
- Removes need for `golang.org/x/crypto/ssh`, `github.com/hiddeco/sshsig`, and custom `sshCommitSignerWrapper`

---

### RQ5: Git Version Detection and Validation

**Question**: How to detect and validate git version programmatically?

**Research**:

`git --version` outputs: `git version 2.45.1` (or similar)

Parsing approach:
1. Run `git --version`
2. Parse output with regex: `git version (\d+)\.(\d+)\.(\d+)`
3. Compare major.minor.patch against minimum (2.34.0)

**Decision**: Run `git --version` at repository initialization time and validate >= 2.34.0.

**Rationale**:
- Single check at startup, not per-operation
- Clear error message if version too old
- Version string format is stable across platforms

---

### RQ6: Error Categorization Strategy

**Question**: How to categorize git CLI errors into distinct types?

**Research**:

Git exit codes and error patterns:
| Exit Code | Meaning | Error Type |
|-----------|---------|------------|
| 0 | Success | N/A |
| 1 | Generic error / diff found differences | Varies |
| 128 | Fatal error (not a repo, permission denied, etc.) | Parse stderr |
| 129 | Invalid arguments | `ErrGitCommandFailed` |

Stderr patterns for categorization:
| Pattern | Error Type |
|---------|------------|
| `not a git repository` | `ErrNotGitRepository` (existing) |
| `permission denied` | `ErrGitPermissionDenied` (new) |
| `does not exist` / `pathspec` | `ErrGitFileNotFound` (new) |
| `signing failed` / `error: gpg failed` | `ErrGitSigningFailed` (new) |
| `command not found` / exec error | `ErrGitNotFound` (new) |
| `version` check fails | `ErrGitVersionTooOld` (new) |
| Other non-zero exit | `ErrGitCommandFailed` (new, wraps stderr) |

**Decision**: Create sentinel error variables matching the patterns above. Parse stderr content to categorize errors.

**Rationale**:
- Maintains compatibility with existing `utils.ErrNotGitRepository`
- Enables callers to handle specific failure modes with `errors.Is()`
- stderr parsing is reliable since git's error messages are stable
- Follows Go idiomatic error handling with wrapped sentinel errors

---

### RQ7: Removing go-git Dependencies

**Question**: Which dependencies can be removed from go.mod after migration?

**Research**:

Direct dependencies to remove:
- `github.com/go-git/go-git/v5` v5.16.4
- `github.com/hiddeco/sshsig` v0.2.0 (only used for SSH signing wrapper)
- `github.com/sergi/go-diff` v1.4.0 (only used for diff computation via go-git's diff util)
- `golang.org/x/crypto` v0.45.0 (only used for `ssh.ParsePrivateKey` and `ssh.Signer`)

Dependencies to **retain**:
- `github.com/go-git/gcfg/v2` v2.0.2 - Used by `pkg/git/config/extractor.go` for parsing git config INI files (independent of go-git)

Dependencies that become indirect-only or removable:
- `github.com/go-git/go-billy/v5` - Transitive from go-git, can be removed
- `github.com/go-git/gcfg` v1 (indirect) - Can be removed (only v2 is directly used)
- Multiple transitive deps from go-git (ProtonMail, cloudflare, etc.)

**Decision**: Remove go-git/v5, sshsig, sergi/go-diff, and x/crypto. Retain gcfg/v2. Run `go mod tidy` to clean up transitives.

**Rationale**:
- Significant reduction in dependency tree
- gcfg/v2 is independently useful for INI parsing
- `go mod tidy` will automatically remove unused indirect dependencies
- Reduces binary size and compilation time

**Note**: Need to verify `golang.org/x/crypto` is not used elsewhere. If `charmbracelet/huh` or other deps pull it in transitively, `go mod tidy` will handle it.

---

### RQ8: Commit Message Author Configuration via CLI

**Question**: How to set commit author name/email via git CLI without modifying git config files?

**Research**:

Options:
1. Environment variables: `GIT_AUTHOR_NAME`, `GIT_AUTHOR_EMAIL`, `GIT_COMMITTER_NAME`, `GIT_COMMITTER_EMAIL`
2. Config overrides: `git -c user.name="Name" -c user.email="email" commit ...`

Both approaches work without modifying persistent git config.

**Decision**: Use environment variables (`GIT_AUTHOR_NAME`, `GIT_AUTHOR_EMAIL`, `GIT_COMMITTER_NAME`, `GIT_COMMITTER_EMAIL`) set on the `exec.Cmd.Env` field.

**Rationale**:
- Environment variables are cleaner than multiple `-c` flags
- Can be set per-command on `exec.Cmd` without affecting global state
- Standard git mechanism for overriding author/committer info
- Current code already extracts user.name and user.email from git config via `pkg/git/config`

---

## Summary of Decisions

| Decision | Choice | Key Rationale |
|----------|--------|---------------|
| CLI command mapping | Direct replacement per API | Stable porcelain formats |
| Status parsing | `git status --porcelain=v1` | Machine-parseable, stable format |
| Diff computation | `git diff --cached --unified=0` (batch) | Single process, native handling |
| SSH signing | Delegate to git CLI with `-c` flags | Native support, removes 3 deps |
| Version detection | `git --version` at init | Single check, clear errors |
| Error categorization | Sentinel errors + stderr parsing | Go-idiomatic, actionable errors |
| Dependency cleanup | Remove go-git, sshsig, go-diff, x/crypto | Major dep reduction |
| Author config | Env vars on exec.Cmd | Clean per-command override |

## Dependencies

### Removed
- `github.com/go-git/go-git/v5` (v5.16.4)
- `github.com/hiddeco/sshsig` (v0.2.0)
- `github.com/sergi/go-diff` (v1.4.0)
- `golang.org/x/crypto` (v0.45.0) - if not used elsewhere

### Retained
- `github.com/go-git/gcfg/v2` (v2.0.2) - git config INI parsing

### New
- None (all functionality via stdlib `os/exec`)

## References

- Git status porcelain format: https://git-scm.com/docs/git-status#_porcelain_format_version_1
- Git diff format: https://git-scm.com/docs/git-diff
- Git SSH signing: https://git-scm.com/docs/git-config#Documentation/git-config.txt-gpgformat
- Git environment variables: https://git-scm.com/book/en/v2/Git-Internals-Environment-Variables
