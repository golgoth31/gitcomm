# CLI Contract: Auto-Stage and State Restoration

**Feature**: 002-auto-stage-restore
**Date**: 2025-01-27
**Command**: `gitcomm`

## Extended Workflow Contract

### Pre-Existing Workflow (from 001-git-commit-cli)

1. Repository State Analysis
2. AI Usage Decision
3. AI Generation (if selected)
4. Message Validation
5. Commit Creation

### New Workflow Steps

#### Step 0: Staging State Capture (NEW)

**When**: Immediately after CLI launch, before any other operations

**Action**: Capture current staging state

**Input**: None (reads from current git repository)

**Output**:
- Pre-CLI staging state snapshot stored in memory
- Used later for restoration if needed

**Error Cases**:
- Not in git repository → Exit with code 1, message: "Error: not a git repository"
- Git operation failed → Exit with code 1, message: "Error: failed to capture staging state"

---

#### Step 0.5: Auto-Stage Modified Files (NEW)

**When**: After staging state capture, before repository state analysis

**Action**: Automatically stage all modified files

**Input**:
- `-a` flag (optional): If present, also stage untracked files

**Output**:
- All modified files staged (or modified + untracked if `-a` flag)
- Staging result (success/failure, list of staged files)

**Error Cases**:
- Staging fails (partial or complete) → Restore state, exit with code 1, message: "Error: failed to stage files: [details]"
- No files to stage → Continue normally (no-op)

**Behavior**:
- If `-a` flag: Stage all files (modified + untracked)
- If no `-a` flag: Stage only modified files
- If any file fails to stage: Abort, restore all staged files, exit with error
- If no files to stage: Continue with existing workflow

---

#### Step 6: State Restoration (NEW)

**When**: CLI exits without committing (cancellation, error, interruption)

**Action**: Restore staging state to pre-CLI state

**Input**: Pre-CLI staging state snapshot

**Output**:
- Staging state restored to pre-CLI state
- Only files staged by CLI are unstaged (preserves originally staged files)

**Error Cases**:
- Restoration fails → Log error, display warning, exit with code 1, message: "Warning: failed to restore staging state. Repository may be in unexpected state. [recovery instructions]"

**Behavior**:
- Only unstage files that were staged by CLI
- Preserve files that were already staged before CLI launch
- If restoration fails, display clear error message with recovery instructions

---

## Signal Handling Contract

### Interruption Signals

**Signals Handled**: `SIGINT` (Ctrl+C), `SIGTERM`

**Behavior**:
1. Register signal handlers at CLI startup
2. On signal receipt:
   - If staging in progress: Abort staging, restore partial state, exit
   - If staging complete: Restore state, exit
   - If commit in progress: Complete commit if possible, then exit
3. Exit with code 130 (SIGINT) or 143 (SIGTERM)

**Error Cases**:
- Signal handler fails → Best-effort restoration, log error, exit

---

## Exit Codes

**Existing Codes** (from 001-git-commit-cli):
- `0`: Success - commit created
- `1`: Error - general error
- `2`: Configuration error
- `3`: AI provider error

**New Codes**:
- `130`: Interrupted by SIGINT (Ctrl+C)
- `143`: Interrupted by SIGTERM

**Error Code Usage**:
- Staging failure: Exit code 1
- Restoration failure: Exit code 1 (with warning message)
- Interruption: Exit code 130 or 143

---

## User-Facing Messages

### Staging Messages

- Success: (silent, no message needed)
- Failure: "Error: failed to stage files: [file list]. Staging state has been restored."
- Partial failure: "Error: failed to stage some files: [file list]. All staged files have been restored. Please fix issues and try again."

### Restoration Messages

- Success: (silent, no message needed)
- Failure: "Warning: failed to restore staging state. Repository may be in unexpected state. Please check git status and manually restore if needed."

### Interruption Messages

- "Interrupted. Restoring staging state..."
- "Staging state restored."

---

## Performance Contract

- Staging state capture: < 500ms
- Auto-staging: < 2 seconds (SC-001)
- State restoration: < 1 second (SC-003)
- Signal handling: < 100ms (non-blocking)

---

## Testing Contract

### Manual Testing Scenarios

1. **Normal Flow**:
   - Run CLI with modified files
   - Verify files are auto-staged
   - Complete commit
   - Verify staging state not restored

2. **Cancellation Flow**:
   - Run CLI with modified files
   - Verify files are auto-staged
   - Cancel CLI (Ctrl+C or reject commit)
   - Verify staging state restored

3. **Interruption Flow**:
   - Run CLI with modified files
   - Interrupt during staging (Ctrl+C)
   - Verify partial state restored

4. **Error Flow**:
   - Run CLI with file permission issues
   - Verify staging fails
   - Verify state restored
   - Verify error message displayed

### Automated Testing

- Unit tests for staging state capture/restoration
- Integration tests for full workflow
- Signal handling tests (mock signals)
- Error scenario tests
