# Contract: CLI Flag Extension

**Feature**: 012-git-config-commit
**Date**: 2025-01-27
**Component**: `cmd/gitcomm`

## Interface Changes

### New CLI Flag

**Flag Name**: `--no-sign`
**Type**: Boolean flag
**Default**: `false` (signing enabled by default when config available)
**Purpose**: Disable commit signing regardless of git config settings (FR-008)

### Flag Definition

**Location**: `cmd/gitcomm/main.go`

```go
var (
    // ... existing flags ...
    noSign bool  // NEW: Disable commit signing
)

func init() {
    // ... existing flag definitions ...
    rootCmd.Flags().BoolVar(&noSign, "no-sign", false, "Disable commit signing")
}
```

---

## Behavior Specification

### Flag Precedence

**Priority Order** (highest to lowest):
1. `--no-sign` flag (explicit user override)
2. `commit.gpgsign = false` in git config (explicit opt-out)
3. SSH signing configuration (if `gpg.format = ssh` and `user.signingkey` set)

**Logic**:
- If `--no-sign` is true: Do not sign, regardless of git config
- If `--no-sign` is false: Check git config for signing settings
- If `commit.gpgsign = false`: Do not sign (even if signing key configured)
- If `gpg.format = ssh` and `user.signingkey` set and `commit.gpgsign != false`: Sign commits

### Flag Propagation

**Flow**:
```
main.go (flag parsing)
    ↓
CommitService (receives flag value)
    ↓
GitRepository.CreateCommit() (checks flag via service)
    ↓
CommitSigner.Enabled (false if --no-sign set)
```

**Implementation**:
- Pass `noSign` flag value to `CommitService` constructor or method
- `CommitService` passes flag to `GitRepository` if needed
- `GitRepository` respects flag when preparing `CommitSigner`

---

## Test Requirements

### Unit Tests

**File**: `cmd/gitcomm/main_test.go` (if exists) or integration tests

**Test Cases**:
1. `--no-sign` flag disables signing even when git config has signing enabled
2. `--no-sign` flag works with local config signing
3. `--no-sign` flag works with global config signing
4. Default behavior (no flag): Signing enabled when config available
5. Flag parsing: `--no-sign` sets flag to true

### Integration Tests

**File**: `test/integration/git_config_test.go`

**Test Cases**:
1. End-to-end: `gitcomm --no-sign` creates unsigned commit despite config
2. End-to-end: `gitcomm` (no flag) creates signed commit when config available

---

## Breaking Changes

**None**: New optional flag, backward compatible.

---

## Migration Notes

**N/A**: New feature, no migration required.

**Usage Examples**:
```bash
# Sign commits (default when config available)
gitcomm

# Disable signing explicitly
gitcomm --no-sign

# Combine with other flags
gitcomm -a --no-sign
```
