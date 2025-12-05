# UI Contract: Display Commit Type Selection Confirmation

**Feature**: 006-type-confirmation-display
**Date**: 2025-01-27

## Function Contract

### PromptCommitType

**Signature**: `func PromptCommitType(reader *bufio.Reader) (string, error)`

**Preconditions**:
- Valid git repository context
- Terminal supports Unicode (for checkmark symbol)
- Bubbletea alt-screen mode available

**Postconditions**:
- If selection successful: Confirmation line displayed, selected type returned
- If selection cancelled: No confirmation line, error returned
- Alt-screen cleared before confirmation display

**Behavior**:
1. Display interactive select list (existing behavior)
2. User selects type and presses Enter OR presses Escape
3. If Enter pressed:
   - Alt-screen clears
   - **NEW**: Display confirmation line: `"✔ Choose a type(<scope>): <chosen type>\n"`
   - Return selected type string
4. If Escape pressed:
   - Alt-screen clears
   - **NO confirmation line displayed**
   - Return error: `"commit type selection cancelled"`

**Output Format**:
- Confirmation line: `"✔ Choose a type(<scope>): %s\n"` where `%s` is the selected commit type
- Checkmark symbol: Unicode ✔ (U+2714)
- Format text: Literal "Choose a type(<scope>):" (scope is not filled, it's literal text)
- Chosen type: Exact value from `SelectListModel.GetSelectedType()`

**Error Cases**:
- Selection cancelled → Returns error, no confirmation displayed
- Invalid model type → Returns error, no confirmation displayed
- GetSelectedType fails → Returns error, no confirmation displayed

**Side Effects**:
- Writes to standard output (stdout) for confirmation line
- Clears alt-screen (bubbletea behavior)
- No state changes, no file I/O, no network calls

## Display Timing Contract

**Timing Requirements**:
- Confirmation line MUST appear after alt-screen clears
- Confirmation line MUST appear before function returns
- Confirmation line MUST appear before next prompt (scope) in calling code
- Display latency: < 100ms (perceived as immediate)

**Visual Layout**:
```
[Alt-screen with select list]
  ↓ (User presses Enter)
[Alt-screen clears]
[✔ Choose a type(<scope>): feat]  ← Confirmation line (new line)
[Scope (optional, press Enter to skip): ]  ← Next prompt
```

## Integration Points

**Called By**:
- `internal/service/commit_service.go` → `promptCommitMessage()` function

**Calls**:
- `NewSelectListModel()` - Creates select list model
- `tea.NewProgram()` - Creates bubbletea program
- `SelectListModel.GetSelectedType()` - Gets selected type value
- `fmt.Printf()` - Displays confirmation line

**No Breaking Changes**:
- Function signature unchanged
- Return values unchanged
- Error handling unchanged
- Only addition: confirmation line display
