# Quickstart: Improve Commit Generation with Git Config

**Feature**: 012-git-config-commit
**Date**: 2025-01-27

## Overview

This feature automatically extracts git configuration (user.name, user.email, SSH signing settings) from `.git/config` and `~/.gitconfig` files and uses these values to configure commit author and signing. This ensures commits are properly attributed and signed according to your git configuration.

## Prerequisites

- Git repository initialized
- Git config files configured (optional, defaults used if missing)
- SSH signing keys set up (optional, for commit signing)

## Setup

### 1. Configure Git User Information

**Local repository config** (`.git/config`):
```ini
[user]
    name = John Doe
    email = john@example.com
```

**Global user config** (`~/.gitconfig`):
```ini
[user]
    name = John Doe
    email = john@example.com
```

**Note**: Local config takes precedence over global config.

### 2. Configure SSH Commit Signing (Optional)

**Local repository config** (`.git/config`):
```ini
[user]
    signingkey = /path/to/your/public-key.pub
[gpg]
    format = ssh
[commit]
    gpgsign = true
```

**Global user config** (`~/.gitconfig`):
```ini
[user]
    signingkey = ~/.ssh/id_ed25519.pub
[gpg]
    format = ssh
[commit]
    gpgsign = true
```

**SSH Key Setup**:
```bash
# Generate SSH key for signing (if not already done)
ssh-keygen -t ed25519 -C "your_email@example.com"

# Add public key to git config
git config --global user.signingkey ~/.ssh/id_ed25519.pub
git config --global gpg.format ssh
git config --global commit.gpgsign true
```

## Usage

### Basic Usage

**Create commit with automatic author from git config**:
```bash
gitcomm
```

The commit will use:
- Author name: From `user.name` in git config
- Author email: From `user.email` in git config
- Signing: Automatic if SSH signing configured

### Disable Commit Signing

**Create commit without signing** (overrides git config):
```bash
gitcomm --no-sign
```

### Combine with Other Flags

**Auto-stage files and create signed commit**:
```bash
gitcomm -a
```

**Auto-stage files and create unsigned commit**:
```bash
gitcomm -a --no-sign
```

**Create commit without signoff and with signing**:
```bash
gitcomm -s
```

## Verification

### Verify Commit Author

```bash
# Check last commit author
git log -1 --format="%an <%ae>"

# Should match your git config user.name and user.email
```

### Verify Commit Signing

```bash
# Check if commit is signed
git log --show-signature -1

# Should show SSH signature if signing was enabled
```

### Verify Config Extraction

**Enable debug logging**:
```bash
gitcomm --debug
```

**Look for debug messages**:
- `Extracted git config from .git/config`
- `Extracted git config from ~/.gitconfig`
- `SSH signing configured: key=/path/to/key.pub`
- `SSH signing disabled: commit.gpgsign=false`

## Scenarios

### Scenario 1: Local Config Takes Precedence

**Setup**:
- `.git/config`: `user.name = "Local User"`
- `~/.gitconfig`: `user.name = "Global User"`

**Result**: Commit author is "Local User" (local config used)

### Scenario 2: Global Config Used When Local Missing

**Setup**:
- `.git/config`: Not present or missing `user.name`
- `~/.gitconfig`: `user.name = "Global User"`

**Result**: Commit author is "Global User" (global config used)

### Scenario 3: Defaults Used When Config Missing

**Setup**:
- `.git/config`: Not present
- `~/.gitconfig`: Not present or missing `user.name`

**Result**: Commit author is "gitcomm <gitcomm@local>" (defaults used)

### Scenario 4: SSH Signing Enabled

**Setup**:
- `.git/config`: `gpg.format = ssh`, `user.signingkey = "/path/to/key.pub"`, `commit.gpgsign = true`

**Result**: Commit is signed with SSH key

### Scenario 5: SSH Signing Disabled by Flag

**Setup**:
- `.git/config`: SSH signing configured
- Command: `gitcomm --no-sign`

**Result**: Commit is not signed (flag overrides config)

### Scenario 6: SSH Signing Disabled by Config

**Setup**:
- `.git/config`: `gpg.format = ssh`, `user.signingkey = "/path/to/key.pub"`, `commit.gpgsign = false`

**Result**: Commit is not signed (config opt-out)

### Scenario 7: Signing Failure Handling

**Setup**:
- `.git/config`: SSH signing configured
- SSH private key file missing or unreadable

**Result**: Commit is created without signing, error logged in debug output

## Troubleshooting

### Issue: Commit Author Not Matching Git Config

**Check**:
1. Verify git config files exist and are readable
2. Check debug logs for config extraction messages
3. Verify config values are correct: `git config --list`

**Solution**: Ensure `.git/config` or `~/.gitconfig` contains `user.name` and `user.email`

### Issue: Commits Not Being Signed

**Check**:
1. Verify `gpg.format = ssh` is set
2. Verify `user.signingkey` points to valid public key file
3. Verify `commit.gpgsign` is not set to `false`
4. Check if `--no-sign` flag was used
5. Check debug logs for signing configuration messages

**Solution**: Ensure SSH signing is properly configured in git config

### Issue: Signing Fails But Commit Created

**Expected Behavior**: This is correct behavior (FR-013). Signing failures are handled gracefully.

**Check**:
1. Verify SSH private key file exists (remove `.pub` from `user.signingkey` path)
2. Verify SSH agent is running (if using agent)
3. Check debug logs for signing error messages

**Solution**: Fix SSH key configuration, commit will be signed on next attempt

## Integration with Existing Features

### Auto-Staging

Config extraction works with auto-staging:
```bash
gitcomm -a  # Auto-stage files, use git config for author and signing
```

### No Signoff

Config extraction works with no-signoff:
```bash
gitcomm -s  # No signoff, but still uses git config for author and signing
```

### AI-Assisted Generation

Config extraction works with AI providers:
```bash
gitcomm --provider openai  # AI generates message, git config used for author and signing
```

## Performance

- Config extraction: <50ms per file (SC-003, SC-004)
- No noticeable impact on commit creation time
- Config cached per repository instance (no repeated file reads)

## Debugging

**Enable debug logging**:
```bash
gitcomm --debug
```

**Debug messages include**:
- Config file paths being read
- Extracted config values
- Signing configuration status
- Signing failures (if any)
- Default values used (if config missing)
