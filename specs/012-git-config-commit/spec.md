# Feature Specification: Improve Commit Generation with Git Config

**Feature Branch**: `012-git-config-commit`
**Created**: 2025-01-27
**Status**: Draft
**Input**: User description: "improve commit generation: - before initialising git objects, extract git configuration from the first git config file: 1. ./.git/config 2. ~/.gitconfig - use user and useremail to generate commit - use signing method to configur go-git signer and sign commits add a flag to disable commit signing (MUST be done by default. if no git config is available, silently ignore. add debug logs"

## User Scenarios & Testing

### User Story 1 - Git Config Extraction and Commit Author Configuration (Priority: P1)

Users want gitcomm to automatically use their git configuration (user.name and user.email) when creating commits, ensuring commits are attributed correctly without manual configuration.

**Why this priority**: This is fundamental to proper commit attribution. Without this, commits may be incorrectly attributed or fail validation in repositories that require proper author information.

**Independent Test**: Configure git with user.name and user.email, run gitcomm to create a commit, and verify the commit author matches the git configuration.

**Acceptance Scenarios**:

1. **Given** a git repository with `.git/config` containing `user.name = "John Doe"` and `user.email = "john@example.com"`, **When** user creates a commit with gitcomm, **Then** the commit author is "John Doe <john@example.com>"
2. **Given** a git repository without `.git/config` but with `~/.gitconfig` containing user.name and user.email, **When** user creates a commit with gitcomm, **Then** the commit author uses values from `~/.gitconfig`
3. **Given** a git repository with both `.git/config` and `~/.gitconfig` containing different user.name values, **When** user creates a commit with gitcomm, **Then** the commit author uses values from `.git/config` (local config takes precedence)
4. **Given** a git repository with no git config files available, **When** user creates a commit with gitcomm, **Then** the commit is created successfully with default values and debug logs indicate config was not found

---

### User Story 2 - Commit Signing with Git Config (Priority: P2)

Users want commits to be automatically signed using their git SSH signing configuration (`gpg.format = ssh`, `commit.gpgsign`, `user.signingkey`) when available, ensuring commit authenticity and integrity.

**Why this priority**: Commit signing is important for security and compliance in many organizations. Users who have configured SSH signing in git expect it to work automatically.

**Independent Test**: Configure git with SSH signing settings (`gpg.format = ssh` and `user.signingkey` pointing to SSH public key file), run gitcomm to create a commit, and verify the commit is signed with the configured SSH key.

**Acceptance Scenarios**:

1. **Given** a git repository with `.git/config` containing `gpg.format = ssh`, `user.signingkey = "/path/to/key.pub"` and `commit.gpgsign = true`, **When** user creates a commit with gitcomm, **Then** the commit is signed with the specified SSH key
2. **Given** a git repository with `~/.gitconfig` containing SSH signing configuration, **When** user creates a commit with gitcomm, **Then** the commit is signed using the global configuration
3. **Given** a git repository with SSH signing configured, **When** user runs gitcomm with `--no-sign` flag, **Then** the commit is created without signing
4. **Given** a git repository with no SSH signing configuration (missing `gpg.format = ssh` or `user.signingkey`), **When** user creates a commit with gitcomm, **Then** the commit is created without signing and debug logs indicate no signing key was found

---

### Edge Cases

- What happens when `.git/config` exists but is unreadable or corrupted? **Answer**: System silently ignores the file, logs debug message, and falls back to `~/.gitconfig` or defaults (FR-009, FR-010)
- What happens when `~/.gitconfig` exists but is unreadable or corrupted? **Answer**: System silently ignores the file, logs debug message, and uses defaults (FR-009, FR-010)
- What happens when user.name is set but user.email is missing (or vice versa)? **Answer**: System uses the available value from git config and the hardcoded default ("gitcomm" or "gitcomm@local") for the missing value (FR-012)
- What happens when SSH key file is configured but the file is not available or unreadable? **Answer**: System creates an unsigned commit, logs the error, and proceeds with commit creation (same as other signing failures) (FR-013)
- What happens when commit.gpgsign is set to false in git config but user.signingkey is set? **Answer**: System respects `commit.gpgsign = false` and does not sign the commit (explicit opt-out takes precedence over key presence) (FR-007)
- What happens when multiple SSH keys are configured (local vs global)? **Answer**: System uses the first available value found (local config checked first, then global config) (FR-006)
- What happens when the signing process fails (e.g., SSH key file not found, SSH agent not running)? **Answer**: System creates an unsigned commit, logs the error, and proceeds with commit creation (FR-013)
- What happens when `gpg.format` is not set to "ssh" or is missing? **Answer**: System does not sign commits (SSH signing requires `gpg.format = ssh`) (FR-006)

## Requirements

### Functional Requirements

- **FR-001**: System MUST extract git configuration by directly reading and parsing `.git/config` (local repository config) as an INI-format file before opening the repository with go-git
- **FR-002**: System MUST extract git configuration by directly reading and parsing `~/.gitconfig` (global user config) as an INI-format file if local config is not available or missing values
- **FR-003**: System MUST use `user.name` from git config as the commit author name
- **FR-004**: System MUST use `user.email` from git config as the commit author email
- **FR-005**: System MUST use local repository config (`.git/config`) values when both local and global configs contain the same settings (local takes precedence)
- **FR-006**: System MUST configure commit signing using git config values (`gpg.format = ssh`, `user.signingkey`, `commit.gpgsign`) when available, using the first available value found (check local config first, then global config). System MUST only sign commits when `gpg.format = ssh` is set
- **FR-007**: System MUST sign commits by default when SSH signing configuration is available (both `gpg.format = ssh` is set AND `user.signingkey` is set AND `commit.gpgsign` is not explicitly set to false)
- **FR-008**: System MUST provide a CLI flag (e.g., `--no-sign`) to disable commit signing
- **FR-009**: System MUST silently ignore missing or unreadable git config files (no user-facing errors)
- **FR-010**: System MUST log debug messages when git config files are not found or cannot be read
- **FR-011**: System MUST log debug messages when SSH signing configuration is not available (missing `gpg.format = ssh` or `user.signingkey`)
- **FR-012**: System MUST handle missing user.name or user.email gracefully by using hardcoded defaults: "gitcomm" for name and "gitcomm@local" for email when git config values are not available
- **FR-013**: System MUST handle SSH signing failures gracefully by creating an unsigned commit, logging the error, and proceeding with commit creation (do not fail commit creation due to signing errors)

### Key Entities

- **GitConfig**: Represents extracted git configuration values (user.name, user.email, signing configuration)
- **CommitSigner**: Represents the configured commit signer (SSH key file path, format, etc.) extracted from git config

## Success Criteria

### Measurable Outcomes

- **SC-001**: 100% of commits created with gitcomm use the correct author name and email from git configuration when config is available
- **SC-002**: Commits are signed automatically when git signing configuration is present in 100% of cases (unless `--no-sign` flag is used)
- **SC-003**: System successfully extracts git config from local repository config in under 50ms
- **SC-004**: System successfully extracts git config from global user config in under 50ms when local config is unavailable
- **SC-005**: Debug logs provide sufficient information to diagnose git config extraction issues (config file paths, missing values, signing key availability)
- **SC-006**: Users can disable commit signing via CLI flag in 100% of cases, regardless of git config settings

## Assumptions

- Git config files follow standard git config format (INI-style)
- Users have proper permissions to read `.git/config` and `~/.gitconfig`
- SSH public key files referenced in `user.signingkey` are available and readable when signing is enabled
- Default behavior (signing enabled when config available) is acceptable for most users
- Silent failure (no user-facing errors) for missing config is acceptable behavior
- Debug logging is sufficient for troubleshooting config extraction issues

## Dependencies

- Existing git repository implementation (`internal/repository/git_repository_impl.go`)
- go-git library support for SSH commit signing
- CLI flag parsing infrastructure (existing `-a`, `-s` flags)
- Debug logging infrastructure (existing debug logging system)

## Clarifications

### Session 2025-01-27

- Q: When SSH signing fails (e.g., key file not found, SSH agent not running), should the system fail the commit or create an unsigned commit? → A: Create unsigned commit - proceed without signing, log error (Option B)
- Q: Should the system support GPG or SSH signing? → A: Support only SSH signing method, NOT GPG (user requirement)
- Q: What default values should be used when user.name or user.email is missing from git config? → A: Use hardcoded defaults: "gitcomm" for name, "gitcomm@local" for email (Option A)
- Q: How should the system extract git config - directly from files or through go-git? → A: Read config files directly (parse as INI) before opening repository (Option A)
- Q: When commit.gpgsign=false but user.signingkey is set, should the system sign or not? → A: Respect commit.gpgsign=false, do not sign (opt-out takes precedence) (Option A)
- Q: When both local and global config contain different user.signingkey values, which should be used? → A: Use first available value (local first, then global) (Option C)

## Out of Scope

- Support for system-wide git config (`/etc/gitconfig`)
- Support for custom git config file locations
- Automatic SSH key generation or setup
- SSH key management (import, export, list keys)
- Support for GPG signing methods (only SSH signing is supported)
- Interactive prompts for missing git config values
- Validation or verification of SSH keys before use
