# Quick Start: Auto-Stage Modified Files and State Restoration

**Feature**: 002-auto-stage-restore
**Date**: 2025-01-27

## Overview

This feature automatically stages modified files when the CLI launches and restores the staging state if you cancel or exit without committing. This ensures AI analysis has access to all changes and prevents accidental staging of files.

## Key Features

- ✅ **Auto-stage modified files** on CLI launch
- ✅ **Auto-stage untracked files** with `-a` flag
- ✅ **Automatic state restoration** on cancellation
- ✅ **Signal handling** for graceful interruption
- ✅ **Preserves original staging** (only restores files staged by CLI)

## Basic Usage

### Auto-Stage Modified Files

```bash
# Navigate to your git repository
cd /path/to/your/repo

# Make some changes to files
echo "new content" >> existing-file.txt

# Run gitcomm - files are automatically staged
gitcomm

# Follow prompts to create commit
# If you cancel, staging state is automatically restored
```

### Auto-Stage All Files (Including Untracked)

```bash
# Create a new file
echo "new file content" > new-file.txt

# Modify an existing file
echo "modified" >> existing-file.txt

# Run gitcomm with -a flag to stage everything
gitcomm -a

# Both modified and untracked files are staged
# If you cancel, all auto-staged files are unstaged
```

## Workflow Examples

### Example 1: Normal Commit with Auto-Staging

```bash
$ gitcomm
# CLI automatically stages modified files
# Shows token count
# Prompts for AI usage
# ... (rest of workflow)
# Commit created successfully
# Staging state NOT restored (commit succeeded)
```

### Example 2: Cancel and Restore

```bash
$ gitcomm
# CLI automatically stages modified files
# Shows token count
# Prompts for AI usage
# User cancels (Ctrl+C or rejects commit)
# CLI automatically restores staging state
# Exit with message: "Staging state restored."
```

### Example 3: Error During Staging

```bash
$ gitcomm
# CLI attempts to stage modified files
# File permission error occurs
# CLI aborts staging, restores any partially staged files
# Exit with error: "Error: failed to stage files: file.txt"
```

### Example 4: Interruption During Staging

```bash
$ gitcomm
# CLI starts staging files
# User presses Ctrl+C during staging
# CLI catches signal, aborts staging, restores partial state
# Exit with code 130 (SIGINT)
```

## Edge Cases

### No Modified Files

```bash
$ gitcomm
# No modified files to stage
# Auto-staging is a no-op
# CLI proceeds with existing workflow (empty commit handling)
```

### Files Already Staged

```bash
# Manually stage some files
$ git add file1.txt

# Run gitcomm
$ gitcomm
# CLI stages additional modified files
# If you cancel, only files staged by CLI are unstaged
# file1.txt remains staged (preserved)
```

### External Changes During CLI

```bash
$ gitcomm
# CLI captures staging state
# Another process stages/unstages files
# CLI continues with original plan
# On restoration, CLI restores to captured state
# Warning logged if current state differs
```

## Troubleshooting

### Staging State Not Restored

**Symptom**: After canceling CLI, files remain staged

**Solution**:
- Check if commit was actually created (check git log)
- Check for restoration errors in logs (use `-v` flag)
- Manually unstage: `git reset HEAD <file>`

### Restoration Failed

**Symptom**: Error message: "Warning: failed to restore staging state"

**Solution**:
- Check git status: `git status`
- Manually restore: `git reset HEAD <file>` for each auto-staged file
- Check for repository lock: `.git/index.lock` (delete if exists)

### Files Not Auto-Staged

**Symptom**: Modified files not staged when CLI runs

**Solution**:
- Verify you're in a git repository: `git status`
- Check file permissions (must be readable)
- Check for git lock: `.git/index.lock` (delete if exists)
- Use `-v` flag for verbose logging

## Best Practices

1. **Review Before Committing**: Even though files are auto-staged, review the diff before committing
2. **Use `-a` Flag Sparingly**: Only use `-a` when you want to include untracked files
3. **Check Status After Cancellation**: Run `git status` after canceling to verify restoration
4. **Handle Errors Promptly**: If staging fails, fix issues before retrying

## Integration with Existing Features

This feature integrates seamlessly with existing gitcomm features:

- **AI Analysis**: Auto-staged files are available for AI analysis
- **Manual Input**: Works with manual commit message creation
- **CLI Options**: `-a` flag extends auto-staging to untracked files
- **Error Handling**: Staging failures trigger restoration and clear error messages

## Performance

- Auto-staging: < 2 seconds for typical repositories
- State restoration: < 1 second
- Signal handling: < 100ms (non-blocking)

## Security

- No secrets involved in staging operations
- File paths validated (no path traversal)
- Git operations use existing security model
- No external network calls

## Next Steps

After using this feature:

1. Review the staged changes: `git diff --cached`
2. Complete the commit workflow
3. Verify commit was created: `git log -1`

For more information, see the [full specification](./spec.md).
