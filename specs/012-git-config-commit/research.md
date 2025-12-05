# Research: Improve Commit Generation with Git Config

**Feature**: 012-git-config-commit
**Date**: 2025-01-27
**Purpose**: Resolve technical unknowns and document design decisions

## Research Questions

### RQ1: INI Parsing Library for Git Config Files

**Question**: What library should be used to parse git config files (INI format) in Go?

**Research**:
- Git config files use INI-style format with sections and key-value pairs
- Existing dependency: `github.com/go-git/gcfg` (v1.5.1-0.20230307220236-3a3c6141e376) already in go.mod
- Alternative: `gopkg.in/ini.v1` - popular INI parser
- Alternative: Standard library - no built-in INI parser

**Decision**: Use `github.com/go-git/gcfg` (already a dependency)

**Rationale**:
- Already present in project dependencies (no new dependency)
- Specifically designed for git config format parsing
- Used by go-git library itself, ensuring compatibility
- Supports git config features (includes, subsections, etc.)
- Lightweight and well-maintained

**Alternatives Considered**:
- `gopkg.in/ini.v1`: More general-purpose, but adds new dependency
- Standard library: Would require custom parser implementation
- Direct file reading + manual parsing: Error-prone and reinvents the wheel

---

### RQ2: go-git SSH Commit Signing Support

**Question**: How does go-git v5 support SSH commit signing?

**Research**:
- go-git v5.16.4 supports SSH signing via `git.CommitOptions.SignKey`
- Requires `github.com/go-git/go-git/v5/plumbing/transport/ssh` package
- SSH signer created from private key file (not public key)
- Git config stores public key path in `user.signingkey`, but signing requires private key
- Need to derive private key path from public key path (remove `.pub` extension)

**Decision**: Use go-git's built-in SSH signing support via `CommitOptions.SignKey`

**Rationale**:
- Native support in go-git library (no external dependencies)
- Follows git's standard SSH signing protocol
- Compatible with git's SSH signing verification
- Private key handling is standard (derive from public key path)

**Implementation Notes**:
- Extract `user.signingkey` from config (public key path)
- Derive private key path (remove `.pub` extension or use standard naming)
- Load private key and create SSH signer
- Pass signer to `CommitOptions.SignKey`
- Handle errors gracefully (create unsigned commit if signing fails)

**Alternatives Considered**:
- External `sshsig` package: Adds dependency, go-git already supports SSH signing
- Custom SSH signing implementation: Complex and error-prone

---

### RQ3: Config Extraction Timing and Location

**Question**: Where and when should config extraction happen relative to repository initialization?

**Research**:
- Requirement: Extract config BEFORE opening repository (FR-001, FR-002)
- Current implementation: Opens repository first, then reads config via `repo.Config()`
- Need to read files directly: `.git/config` and `~/.gitconfig`
- Must handle file path resolution: `~/.gitconfig` needs home directory expansion

**Decision**: Extract config in `NewGitRepository` constructor before `git.PlainOpen()`

**Rationale**:
- Meets requirement to extract before initializing git objects
- Allows using extracted values during repository initialization if needed
- Centralizes config extraction logic
- Maintains single responsibility: repository creation includes config loading

**Implementation Approach**:
1. Create `pkg/git/config/extractor.go` with `ConfigExtractor` interface
2. Implement `FileConfigExtractor` that reads `.git/config` and `~/.gitconfig` directly
3. Call extractor in `NewGitRepository` before `git.PlainOpen()`
4. Store extracted config in repository struct for use in `CreateCommit()`

**Alternatives Considered**:
- Extract in `CreateCommit()`: Violates requirement to extract before initializing
- Separate config service: Adds unnecessary abstraction layer
- Extract on-demand: Doesn't meet "before initializing" requirement

---

### RQ4: SSH Private Key Path Resolution

**Question**: How to resolve private key path from public key path in `user.signingkey`?

**Research**:
- Git config stores public key path: `/path/to/key.pub`
- SSH signing requires private key: `/path/to/key` (without `.pub`)
- Standard SSH key naming: `id_ed25519` (private) and `id_ed25519.pub` (public)
- Edge cases: Custom key names, different extensions, absolute vs relative paths

**Decision**: Derive private key path by removing `.pub` extension from `user.signingkey` value

**Rationale**:
- Standard convention: public keys have `.pub` extension
- Simple and predictable transformation
- Handles most common cases
- If private key not found, signing fails gracefully (per FR-013)

**Implementation Notes**:
- Check if `user.signingkey` ends with `.pub`
- If yes, remove `.pub` to get private key path
- If no, use value as-is (assume it's already private key path)
- Validate private key file exists and is readable before attempting to sign
- Log debug message if private key not found

**Alternatives Considered**:
- Always assume private key path: Doesn't match git's convention of storing public key
- Try both paths: Adds complexity, public key path is authoritative
- Require separate config for private key: Breaks compatibility with git config

---

### RQ5: CLI Flag for Disabling Commit Signing

**Question**: What should the CLI flag name be for disabling commit signing?

**Research**:
- Existing flags: `-a` (add-all), `-s` (no-signoff), `-d` (debug)
- Git command: `git commit --no-gpg-sign` (for GPG, but we use SSH)
- Consistency: Should follow existing flag patterns
- Clarity: Flag name should clearly indicate it disables signing

**Decision**: Use `--no-sign` flag (long form), no short form

**Rationale**:
- Consistent with existing `--no-signoff` pattern
- Clear and unambiguous
- No short form needed (not frequently used)
- Matches user requirement: "add a flag to disable commit signing"

**Implementation Notes**:
- Add `noSign` boolean flag to `cmd/gitcomm/main.go`
- Pass flag value to commit service
- Check flag before attempting to sign commits
- Flag takes precedence over git config `commit.gpgsign` setting

**Alternatives Considered**:
- `--no-gpg-sign`: Too specific (we use SSH, not GPG)
- `-S` (short form): Conflicts with git's `-S` for GPG key ID
- `--disable-signing`: More verbose, less consistent with existing flags

---

## Summary of Decisions

1. **INI Parsing**: Use `github.com/go-git/gcfg` (existing dependency)
2. **SSH Signing**: Use go-git's built-in `CommitOptions.SignKey` with SSH transport
3. **Config Extraction**: Extract in `NewGitRepository` before opening repository
4. **Private Key Path**: Derive by removing `.pub` from `user.signingkey`
5. **CLI Flag**: Use `--no-sign` (long form only)

## Open Questions Resolved

All technical unknowns have been resolved. No blocking issues remain.

## Dependencies

- **Existing**: `github.com/go-git/gcfg` (v1.5.1-0.20230307220236-3a3c6141e376) - for config parsing
- **Existing**: `github.com/go-git/go-git/v5` (v5.16.4) - for SSH signing support
- **No new dependencies required**
