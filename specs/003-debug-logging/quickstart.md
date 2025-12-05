# Quick Start: Debug Logging Configuration

**Feature**: 003-debug-logging
**Date**: 2025-01-27

## Overview

This feature adds a debug flag to enable detailed logging when troubleshooting issues. By default, the CLI runs silently (no log output). When the debug flag is enabled, all DEBUG-level log messages are displayed in human-readable structured text format without timestamps.

## Key Features

- ✅ **Debug flag** (`--debug` or `-d`) to enable logging
- ✅ **Silent by default** (no log output unless debug enabled)
- ✅ **Raw text format** (human-readable, not JSON)
- ✅ **No timestamps** in log output
- ✅ **DEBUG-only** log level throughout codebase

## Basic Usage

### Enable Debug Logging

```bash
# Use long form
gitcomm --debug

# Use short form
gitcomm -d

# Combine with other flags
gitcomm --debug -a
gitcomm -d --provider openai
```

### Normal Usage (Silent)

```bash
# No debug flag - no log output
gitcomm

# Verbose flag alone - no log output (verbose is no-op)
gitcomm --verbose
```

## Log Output Format

When debug mode is enabled, log messages appear in this format:

```
[DEBUG] Starting commit creation workflow
[DEBUG] Auto-staging modified files staged_count=3
[DEBUG] Files auto-staged successfully
[DEBUG] Commit created successfully
```

**Characteristics**:
- `[DEBUG]` prefix indicates log level
- Message text follows the prefix
- Structured fields appear as `key=value` pairs
- No timestamps
- Output goes to stderr

## Examples

### Example 1: Debugging Commit Creation

```bash
$ gitcomm --debug
[DEBUG] Starting commit creation workflow
[DEBUG] Auto-staging modified files
[DEBUG] Files auto-staged successfully staged_count=2
[DEBUG] Estimated AI tokens: 150 (input) + 50 (output) = 200 (total)
Do you want to use AI to generate the commit message? (y/n):
```

### Example 2: Debugging with Auto-Stage

```bash
$ gitcomm --debug -a
[DEBUG] Starting commit creation workflow
[DEBUG] Auto-staging all files (including untracked)
[DEBUG] Files auto-staged successfully staged_count=5 failed_count=0
[DEBUG] Estimated AI tokens: 300 (input) + 100 (output) = 400 (total)
...
```

### Example 3: Silent Operation (Default)

```bash
$ gitcomm
Do you want to use AI to generate the commit message? (y/n):
# No log messages appear
```

## Error Messages vs Log Messages

**Important**: Error messages are separate from log messages and are always displayed:

```bash
$ gitcomm
Error: failed to initialize git repository: not a git repository
# Error message shown even without debug flag
```

```bash
$ gitcomm --debug
[DEBUG] Starting commit creation workflow
Error: failed to initialize git repository: not a git repository
[DEBUG] Workflow cancelled by signal
# Both error message and debug logs shown
```

## Troubleshooting

### No Log Output When Debug Flag Used

**Symptom**: Using `--debug` flag but no log messages appear

**Possible Causes**:
- No DEBUG-level log statements executed
- Logger not properly initialized
- Check that flag is spelled correctly: `--debug` or `-d`

**Solution**:
- Verify flag is set: `gitcomm --debug --help` should show debug flag
- Check that code has DEBUG-level log statements
- Ensure logger initialization receives debug flag

### Log Output in JSON Format

**Symptom**: Log messages appear in JSON format instead of raw text

**Solution**:
- Verify logger is configured with `ConsoleWriter` for raw text format
- Check that `TimeFormat` is set to empty string
- Ensure debug flag is properly passed to `InitLogger`

### Timestamps Appear in Logs

**Symptom**: Log messages include timestamps when they shouldn't

**Solution**:
- Verify `ConsoleWriter.TimeFormat` is set to empty string
- Check logger configuration in `InitLogger` function

## Best Practices

1. **Use Debug Flag Sparingly**: Only enable when troubleshooting issues
2. **Check Log Format**: Verify output matches expected format
3. **Separate Errors from Logs**: Remember error messages are always shown
4. **Review Log Content**: Debug logs may contain sensitive information

## Integration with Existing Features

This feature integrates seamlessly with existing gitcomm features:

- **Auto-Staging**: Debug logs show staging operations when enabled
- **AI Generation**: Debug logs show token calculations and AI calls
- **State Restoration**: Debug logs show restoration operations
- **Error Handling**: Error messages remain separate from logging

## Performance

- Logger initialization: <10ms (negligible)
- Log message output: <1ms per message
- No performance impact when debug disabled
- Minimal overhead when debug enabled

## Security

- Debug logs may contain sensitive information (file paths, repository state)
- Users should be aware of log content when sharing
- No secrets are logged (existing constraint maintained)
- Error messages remain separate (not affected by logging)

## Next Steps

After using debug logging:

1. Review the log output to identify issues
2. Fix any problems identified
3. Run CLI again without debug flag for normal operation

For more information, see the [full specification](./spec.md).
