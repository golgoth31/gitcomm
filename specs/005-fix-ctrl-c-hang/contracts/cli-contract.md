# CLI Contract: Fix CLI Hang on Ctrl+C During State Restoration

**Feature**: 005-fix-ctrl-c-hang
**Date**: 2025-01-27

## Signal Handling Contract

### Interruption Signals

**Signals Handled**: `SIGINT` (Ctrl+C), `SIGTERM`

**Behavior**:
1. Register signal handlers at CLI startup (existing)
2. On signal receipt:
   - Display "Interrupted. Restoring staging state..." message immediately (FR-005)
   - Create timeout context (3 seconds) for restoration operations (FR-002)
   - Cancel main workflow context to stop ongoing operations
   - Start restoration with timeout context
   - Wait for restoration to complete or timeout (maximum 5 seconds total) (FR-001)
   - Exit with code 130 (SIGINT) or 143 (SIGTERM) (FR-006)

**Timeout Behavior**:
- Restoration operations timeout after 3 seconds (FR-002)
- If timeout occurs, display warning and exit immediately (FR-003, FR-007)
- No partial restoration attempts (per clarification)
- No retry on timeout (FR-003)

**Error Cases**:
- Restoration timeout → Display warning, exit with code 130
- Restoration failure (non-timeout) → Display warning, exit with code 130
- Overall timeout (5 seconds) → Display warning, exit with code 130

---

## Exit Codes

**Existing Codes** (from 001-git-commit-cli):
- `0`: Success - commit created
- `1`: Error - general error
- `2`: Configuration error
- `3`: AI provider error
- `130`: Interrupted by SIGINT (Ctrl+C) (existing, from 002-auto-stage-restore)
- `143`: Interrupted by SIGTERM (existing, from 002-auto-stage-restore)

**No New Exit Codes**: Uses existing exit code 130 for Ctrl+C interruption

---

## User-Facing Messages

### Interruption Messages

- **On Ctrl+C**: "Interrupted. Restoring staging state..." (displayed immediately)
- **On timeout**: "Warning: Restoration timed out. Repository may be in unexpected state. Please check git status and manually restore if needed."
- **On failure**: "Warning: Failed to restore staging state. Repository may be in unexpected state. Please check git status and manually restore if needed."

### Message Timing

- Interruption message: Displayed immediately upon Ctrl+C (FR-005)
- Warning message: Displayed after timeout or failure, before exit (FR-007)
- All messages displayed before exit (FR-008)

---

## Performance Contract

- Signal handler response: < 100ms (non-blocking)
- Timeout context creation: < 1ms
- Restoration timeout: 3 seconds maximum (FR-002)
- Overall exit time: 5 seconds maximum from Ctrl+C (FR-001)
- Context cancellation check: < 1ms per operation

---

## Testing Contract

### Manual Testing Scenarios

1. **Normal Interruption**:
   - Run CLI with modified files
   - Press Ctrl+C during commit workflow
   - Verify "Interrupted. Restoring staging state..." appears
   - Verify CLI exits within 5 seconds
   - Verify exit code 130

2. **Timeout Scenario**:
   - Run CLI with modified files
   - Simulate slow git operations (or use test repository with many files)
   - Press Ctrl+C
   - Verify timeout occurs after 3 seconds
   - Verify warning message displayed
   - Verify CLI exits within 5 seconds total

3. **Multiple Ctrl+C**:
   - Run CLI
   - Press Ctrl+C multiple times rapidly
   - Verify only one restoration attempt
   - Verify CLI exits after first attempt completes or times out

4. **TUI Interruption**:
   - Run CLI, start interactive prompts
   - Press Ctrl+C during TUI interaction
   - Verify TUI cancels immediately
   - Verify restoration proceeds with timeout
   - Verify CLI exits within 5 seconds

### Automated Testing

- Unit tests for timeout context creation
- Integration tests for signal handling with timeout
- Timeout simulation tests (short timeouts in test environment)
- Race condition tests (signal handler vs main process)
- Multiple Ctrl+C handling tests

---

## Backward Compatibility

- **Existing restoration behavior preserved**: When not interrupted, restoration uses normal context (no timeout)
- **Existing signal handling preserved**: Other signal scenarios unchanged
- **Existing exit codes preserved**: No new exit codes introduced
- **Existing error messages preserved**: Additional timeout-specific messages added, existing messages unchanged
