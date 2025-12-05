# Quick Start: Fix CLI Hang on Ctrl+C During State Restoration

**Feature**: 005-fix-ctrl-c-hang
**Date**: 2025-01-27

## Overview

This bug fix ensures that when you press Ctrl+C to interrupt the CLI, it exits promptly (within 5 seconds) instead of hanging indefinitely. The fix adds a 3-second timeout to restoration operations and ensures proper synchronization between signal handling and the main process.

## Problem

**Before the fix**: When pressing Ctrl+C, the CLI would display "Interrupted. Restoring staging state..." but then hang indefinitely, requiring force-kill (kill -9) to exit.

**After the fix**: When pressing Ctrl+C, the CLI displays the message, attempts restoration with a 3-second timeout, and exits within 5 seconds total, even if restoration is incomplete.

## Usage

### Normal Interruption

```bash
# Run CLI
gitcomm

# Press Ctrl+C at any point
# Output:
# Interrupted. Restoring staging state...
# (CLI exits within 5 seconds with code 130)
```

### Timeout Scenario

```bash
# Run CLI with many files or slow git operations
gitcomm

# Press Ctrl+C
# Output:
# Interrupted. Restoring staging state...
# Warning: Restoration timed out. Repository may be in unexpected state.
# Please check git status and manually restore if needed.
# (CLI exits within 5 seconds)
```

### Verification

After interruption, verify repository state:

```bash
# Check git status
git status

# If restoration timed out, you may need to manually unstage files
git reset HEAD <file>
```

## Testing

### Manual Test: Normal Interruption

1. Create a test repository with modified files
2. Run `gitcomm`
3. Press Ctrl+C during any phase
4. Verify CLI exits within 5 seconds
5. Verify exit code is 130
6. Verify git status shows correct staging state

### Manual Test: Timeout Simulation

1. Create a test repository with many files (100+)
2. Run `gitcomm` with `-a` flag to stage all files
3. Press Ctrl+C immediately
4. Verify timeout warning appears after 3 seconds
5. Verify CLI exits within 5 seconds total
6. Check git status to see if restoration completed

### Automated Tests

```bash
# Run integration tests
go test ./test/integration -v -run TestSignalTimeout

# Run all timeout-related tests
go test ./test/integration -v -run Timeout
```

## Troubleshooting

### Symptom: CLI still hangs after fix

**Possible causes**:
- Timeout context not being used in all git operations
- Signal handler not properly synchronized with main process
- Git operations not respecting context cancellation

**Solution**: Check that all restoration operations use the timeout context and that git operations check context cancellation.

### Symptom: Warning message not displayed

**Possible causes**:
- Timeout detection not working correctly
- Error handling not checking for `context.DeadlineExceeded`

**Solution**: Verify timeout error detection using `errors.Is(err, context.DeadlineExceeded)`.

### Symptom: Exit takes longer than 5 seconds

**Possible causes**:
- Overall timeout not enforced
- Main process not waiting for restoration properly

**Solution**: Ensure main process has overall 5-second timeout in addition to 3-second restoration timeout.

## Success Indicators

- ✅ CLI exits within 5 seconds of Ctrl+C in all test cases
- ✅ Restoration operations timeout after 3 seconds if not complete
- ✅ Warning messages displayed when timeout occurs
- ✅ No test case hangs indefinitely
- ✅ Users can exit without force-kill
