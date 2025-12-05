# Contract: ConfigExtractor Interface

**Feature**: 012-git-config-commit
**Date**: 2025-01-27
**Component**: `pkg/git/config`

## Interface Definition

```go
package config

// GitConfig represents extracted git configuration values
type GitConfig struct {
    UserName      string
    UserEmail     string
    SigningKey    string
    GPGFormat     string
    CommitGPGSign bool
}

// ConfigExtractor defines the interface for extracting git configuration
type ConfigExtractor interface {
    // Extract reads git configuration from .git/config and ~/.gitconfig
    // Returns extracted config values, with local config taking precedence
    // Errors are logged but not returned (silent ignore per FR-009)
    Extract(repoPath string) *GitConfig
}
```

## Behavior Specification

### Extract(repoPath string) *GitConfig

**Purpose**: Extract git configuration from local and global config files.

**Parameters**:
- `repoPath` (string): Path to git repository root (used to locate `.git/config`)

**Returns**:
- `*GitConfig`: Extracted configuration values (never nil, uses defaults if extraction fails)

**Behavior**:
1. Attempts to read `.git/config` from repository path
2. If local config unavailable or missing values, attempts to read `~/.gitconfig`
3. Local config values take precedence over global config (FR-005)
4. Missing or unreadable files are silently ignored (FR-009)
5. Debug logs are written for missing/unreadable files (FR-010)
6. Returns config with defaults for missing values (FR-012)

**Error Handling**:
- File not found: Log debug, use defaults
- File unreadable: Log debug, use defaults
- File corrupted: Log debug, use defaults
- No errors returned (all errors handled internally)

**Performance**:
- Must complete in <50ms (SC-003, SC-004)

**Thread Safety**:
- Must be thread-safe (stateless operation)

---

## Implementation Requirements

### FileConfigExtractor

**Location**: `pkg/git/config/extractor.go`

**Dependencies**:
- `github.com/go-git/gcfg` for INI parsing
- `os` for file operations
- `path/filepath` for path manipulation
- `github.com/rs/zerolog` for debug logging

**Implementation Notes**:
- Use `gcfg.ReadFileInto()` for parsing git config files
- Expand `~` in `~/.gitconfig` to user home directory
- Handle both absolute and relative paths
- Cache parsed config per repository path (optional optimization)

---

## Test Requirements

### Unit Tests

**File**: `pkg/git/config/extractor_test.go`

**Test Cases**:
1. Extract from local config only (`.git/config` present)
2. Extract from global config only (`~/.gitconfig` present, no local)
3. Extract with local taking precedence (both present, different values)
4. Extract with missing files (both missing, uses defaults)
5. Extract with unreadable local config (falls back to global)
6. Extract with unreadable global config (uses defaults)
7. Extract with corrupted config file (silent ignore, uses defaults)
8. Extract with partial values (user.name but no user.email)
9. Extract SSH signing configuration (gpg.format = ssh, user.signingkey)
10. Extract with commit.gpgsign = false (signing disabled)

**Performance Tests**:
- Extract completes in <50ms for local config
- Extract completes in <50ms for global config

---

## Integration Points

### Used By

- `internal/repository/git_repository_impl.go`: `NewGitRepository()` calls `Extract()` before opening repository

### Dependencies

- None (self-contained utility)

---

## Breaking Changes

**None**: New interface, no existing code affected.

---

## Migration Notes

**N/A**: New feature, no migration required.
