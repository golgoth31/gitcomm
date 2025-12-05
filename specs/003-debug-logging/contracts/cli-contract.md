# CLI Contract: Debug Logging Configuration

**Feature**: 003-debug-logging
**Date**: 2025-01-27
**Command**: `gitcomm`

## Extended Command Interface

### Existing Flags (from 001-git-commit-cli)

- `-a, --add-all`: Automatically stage all unstaged files
- `-s, --no-signoff`: Disable commit signoff
- `--provider <name>`: Override default AI provider
- `--skip-ai`: Skip AI generation
- `--config <path>`: Path to configuration file
- `-v, --verbose`: Enable verbose logging (now no-op unless debug is also set)

### New Flags

#### Debug Flag

**Flag**: `-d, --debug`

**Description**: Enable debug logging (raw text format, no timestamps)

**Behavior**:
- When set: Enables debug logging, all DEBUG-level messages are displayed
- When not set: No log messages are displayed (silent mode)
- Takes precedence over verbose flag (if both are set, verbose is ignored)

**Examples**:
```bash
# Enable debug logging
gitcomm --debug
gitcomm -d

# Debug logging with other flags
gitcomm --debug -a
gitcomm -d --provider openai
```

---

## Extended Workflow Contract

### Logger Initialization (Modified)

**When**: At CLI startup, before any other operations

**Action**: Initialize logger based on debug and verbose flags

**Input**:
- `--debug` or `-d` flag (optional)
- `--verbose` or `-v` flag (optional, existing)

**Output**:
- Logger configured for debug mode (if debug flag set)
- Logger configured for silent mode (if debug flag not set)
- Verbose flag is no-op when debug flag is set

**Error Cases**:
- None (logger initialization always succeeds)

---

## Log Output Format Contract

### Debug Mode Enabled

**Format**: Human-readable structured text
**Example**: `[DEBUG] message key=value key2=value2`
**Timestamp**: None
**Level**: DEBUG only

### Debug Mode Disabled

**Output**: None (complete silence)
**Error Messages**: Still displayed as plain text (separate from logging)

---

## Exit Codes

**Existing Codes** (unchanged):
- `0`: Success - commit created
- `1`: Error - general error
- `2`: Configuration error
- `3`: AI provider error
- `130`: Interrupted by SIGINT
- `143`: Interrupted by SIGTERM

**No new exit codes** for logging feature.

---

## User-Facing Messages

### Log Messages (Debug Mode)

- Format: `[DEBUG] message key=value`
- Output: stderr
- Only shown when `--debug` flag is used

### Error Messages (Always Shown)

- Format: `Error: message` (plain text)
- Output: stdout/stderr (as appropriate)
- Shown regardless of debug flag state

---

## Performance Contract

- Logger initialization: <10ms
- Log message output: <1ms per message
- No performance impact when debug disabled

---

## Testing Contract

### Manual Testing Scenarios

1. **Debug Mode Enabled**:
   - Run `gitcomm --debug`
   - Verify debug log messages appear in raw text format
   - Verify no timestamps in log output
   - Verify format: `[DEBUG] message key=value`

2. **Debug Mode Disabled**:
   - Run `gitcomm` (no flags)
   - Verify no log messages appear
   - Verify error messages still shown (if errors occur)

3. **Verbose Flag Interaction**:
   - Run `gitcomm --verbose` (no debug)
   - Verify no log messages (verbose is no-op)
   - Run `gitcomm --verbose --debug`
   - Verify debug messages appear (debug takes precedence)

### Automated Testing

- Unit tests for InitLogger with different flag combinations
- Integration tests for CLI flag parsing
- Format verification tests
