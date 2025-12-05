# Data Model: Improve Commit Generation with Git Config

**Feature**: 012-git-config-commit
**Date**: 2025-01-27

## Entities

### GitConfig

Represents extracted git configuration values from `.git/config` and `~/.gitconfig` files.

**Fields**:
- `UserName` (string): Value of `user.name` from git config
- `UserEmail` (string): Value of `user.email` from git config
- `SigningKey` (string): Value of `user.signingkey` from git config (SSH public key file path)
- `GPGFormat` (string): Value of `gpg.format` from git config (should be "ssh" for SSH signing)
- `CommitGPGSign` (bool): Value of `commit.gpgsign` from git config (true/false)

**Relationships**:
- Extracted by `ConfigExtractor` from git config files
- Used by `GitRepository` for commit author and signing configuration

**Validation Rules**:
- `UserName` and `UserEmail` can be empty (will use defaults)
- `SigningKey` must be a valid file path if provided
- `GPGFormat` must be "ssh" for SSH signing to be enabled
- `CommitGPGSign` is optional (defaults to false if not set)

**State Transitions**:
- Created: When config is extracted from files
- Used: When values are applied to commit creation
- No state changes (immutable after creation)

---

### CommitSigner

Represents the configured commit signer extracted from git config and prepared for use with go-git.

**Fields**:
- `PrivateKeyPath` (string): Path to SSH private key file (derived from `user.signingkey`)
- `PublicKeyPath` (string): Path to SSH public key file (from `user.signingkey`)
- `Format` (string): Signing format ("ssh")
- `Enabled` (bool): Whether signing should be performed (based on config and flags)

**Relationships**:
- Created from `GitConfig` values
- Used by `GitRepository.CreateCommit()` via go-git's `CommitOptions.SignKey`

**Validation Rules**:
- `PrivateKeyPath` must point to an existing, readable file if `Enabled` is true
- `PublicKeyPath` must point to an existing, readable file if `Enabled` is true
- `Format` must be "ssh" for SSH signing
- `Enabled` is true only if: `GPGFormat == "ssh"` AND `SigningKey != ""` AND `CommitGPGSign != false` AND `--no-sign` flag not set

**State Transitions**:
- Created: When config is extracted and signing is configured
- Validated: When private key file is checked for existence
- Used: When passed to go-git for commit signing
- Failed: When signing fails (creates unsigned commit)

---

## Data Flow

```
1. NewGitRepository() called
   ↓
2. ConfigExtractor.Extract() reads .git/config and ~/.gitconfig
   ↓
3. GitConfig entity created with extracted values
   ↓
4. Repository struct stores GitConfig
   ↓
5. CreateCommit() called
   ↓
6. GitConfig values used for:
   - Author name/email (from UserName, UserEmail)
   - CommitSigner creation (from SigningKey, GPGFormat, CommitGPGSign)
   ↓
7. CommitSigner validated and used in CommitOptions
```

---

## Default Values

When git config values are missing:

- `UserName`: Defaults to "gitcomm"
- `UserEmail`: Defaults to "gitcomm@local"
- `SigningKey`: Empty (no signing)
- `GPGFormat`: Empty (no signing)
- `CommitGPGSign`: Defaults to false (no signing)

---

## Error Handling

**Config Extraction Errors**:
- Missing files: Silent ignore, use defaults (FR-009)
- Unreadable files: Silent ignore, log debug, use defaults (FR-009, FR-010)
- Corrupted files: Silent ignore, log debug, use defaults (FR-009, FR-010)

**Signing Errors**:
- Private key not found: Create unsigned commit, log error (FR-013)
- SSH agent not running: Create unsigned commit, log error (FR-013)
- Signing failure: Create unsigned commit, log error (FR-013)

---

## Relationships Diagram

```
ConfigExtractor
    │
    ├──> Extracts from
    │    ├──> .git/config (local)
    │    └──> ~/.gitconfig (global)
    │
    └──> Creates
         └──> GitConfig
              │
              ├──> Used by
              │    └──> GitRepository
              │         │
              │         ├──> For Author (UserName, UserEmail)
              │         │
              │         └──> For Signing
              │              └──> CommitSigner
              │                   │
              │                   └──> Used by
              │                        └──> go-git CommitOptions.SignKey
```
