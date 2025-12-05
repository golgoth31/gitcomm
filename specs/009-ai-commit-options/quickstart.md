# Quickstart: AI Commit Message Acceptance Options

**Feature**: 009-ai-commit-options
**Date**: 2025-01-27

## Overview

This feature adds three options when an AI-generated commit message is displayed:
1. **Accept and commit directly**: Commit immediately with the AI message
2. **Accept and edit**: Pre-fill commit message fields with AI values for editing
3. **Reject**: Start over (new AI generation or manual input)

## User Experience

### Scenario 1: Accept and Commit Directly

```
$ gitcomm
? Use AI to generate commit message? (Y/n): y
Generating AI commit message...

--- AI Generated Message ---
feat(auth): add user authentication

Implement JWT-based authentication with refresh tokens.
Add login and logout endpoints.

Closes #123
---

Options:
  1. Accept and commit directly
  2. Accept and edit
  3. Reject
Choose option (1/2/3): 1

✔ Commit created successfully
```

### Scenario 2: Accept and Edit

```
$ gitcomm
? Use AI to generate commit message? (Y/n): y
Generating AI commit message...

--- AI Generated Message ---
feat(auth): add user authentication

Implement JWT-based authentication with refresh tokens.
Add login and logout endpoints.

Closes #123
---

Options:
  1. Accept and commit directly
  2. Accept and edit
  3. Reject
Choose option (1/2/3): 2

? Choose a type(<scope>): feat ✓ [pre-selected]
? Scope (default: auth): api
? Subject (default: add user authentication): add API authentication
? Body (default: Implement JWT-based authentication...): [pre-filled, user can edit]
? Footer (default: Closes #123): [pre-filled, user can edit]

✔ Commit created successfully
```

### Scenario 3: Reject and New AI

```
$ gitcomm
? Use AI to generate commit message? (Y/n): y
Generating AI commit message...

--- AI Generated Message ---
fix: fix bug

---

Options:
  1. Accept and commit directly
  2. Accept and edit
  3. Reject
Choose option (1/2/3): 3

Options:
  1. Generate new AI message
  2. Manual input
Choose option (1/2): 1

Generating new AI commit message...

--- AI Generated Message ---
fix(api): resolve authentication token expiration

Fix issue where JWT tokens expire prematurely.
Update token refresh logic.

---

Options:
  1. Accept and commit directly
  2. Accept and edit
  3. Reject
Choose option (1/2/3): 1

✔ Commit created successfully
```

### Scenario 4: Reject and Manual Input

```
$ gitcomm
? Use AI to generate commit message? (Y/n): y
Generating AI commit message...

--- AI Generated Message ---
fix: fix bug

---

Options:
  1. Accept and commit directly
  2. Accept and edit
  3. Reject
Choose option (1/2/3): 3

Options:
  1. Generate new AI message
  2. Manual input
Choose option (1/2): 2

? Choose a type(<scope>): [interactive list, no pre-selection]
...
```

## Implementation Steps

### Step 1: Add AIMessageAcceptance Type

Create the enum-like type in `internal/ui/prompts.go`:

```go
type AIMessageAcceptance int

const (
    AcceptAndCommit AIMessageAcceptance = iota
    AcceptAndEdit
    Reject
)

func (a AIMessageAcceptance) String() string {
    switch a {
    case AcceptAndCommit:
        return "accept and commit"
    case AcceptAndEdit:
        return "accept and edit"
    case Reject:
        return "reject"
    default:
        return "unknown"
    }
}
```

### Step 2: Implement PromptAIMessageAcceptanceOptions

```go
func PromptAIMessageAcceptanceOptions(reader *bufio.Reader) (AIMessageAcceptance, error) {
    fmt.Println("\n--- AI Generated Message ---")
    fmt.Println(message)
    fmt.Println("---")
    fmt.Println("Options:")
    fmt.Println("  1. Accept and commit directly")
    fmt.Println("  2. Accept and edit")
    fmt.Println("  3. Reject")
    fmt.Print("Choose option (1/2/3): ")

    input, err := reader.ReadString('\n')
    if err != nil {
        return 0, fmt.Errorf("failed to read input: %w", err)
    }

    response := strings.TrimSpace(input)
    switch response {
    case "1":
        return AcceptAndCommit, nil
    case "2":
        return AcceptAndEdit, nil
    case "3":
        return Reject, nil
    default:
        fmt.Println("Invalid option. Please enter 1, 2, or 3.")
        return PromptAIMessageAcceptanceOptions(reader) // Recursive retry
    }
}
```

### Step 3: Extend SelectListModel for Pre-selection

Add constructor in `internal/ui/select_list.go`:

```go
func NewSelectListModelWithPreselection(preselectedType string) SelectListModel {
    model := NewSelectListModel()

    if preselectedType != "" {
        for i, item := range model.Items {
            if item.Type == preselectedType {
                model.SelectedIndex = i
                break
            }
        }
    }

    return model
}
```

### Step 4: Extend MultilineInputModel for Pre-filling

Add constructor in `internal/ui/multiline_input.go`:

```go
func NewMultilineInputModelWithValue(fieldName, initialValue string) MultilineInputModel {
    model := NewMultilineInputModel(fieldName)
    model.Value = initialValue
    return model
}
```

### Step 5: Update CommitService.generateWithAI

Modify the method to:
1. Call `PromptAIMessageAcceptanceOptions()` instead of `PromptAIMessageAcceptance()`
2. Handle three acceptance paths
3. Support pre-filling when AcceptAndEdit is selected
4. Handle rejection with choice prompt

### Step 6: Update promptCommitMessage

Modify to accept `*PrefilledCommitMessage` parameter and use pre-filled values when provided.

## Testing

### Unit Tests

```go
func TestPromptAIMessageAcceptanceOptions_AcceptAndCommit(t *testing.T) {
    reader := bufio.NewReader(strings.NewReader("1\n"))
    acceptance, err := ui.PromptAIMessageAcceptanceOptions(reader)
    assert.NoError(t, err)
    assert.Equal(t, ui.AcceptAndCommit, acceptance)
}
```

### Integration Tests

```go
func TestAIAcceptanceWorkflow_AcceptAndCommit(t *testing.T) {
    // Setup test repo, stage files, etc.
    // Run gitcomm with AI enabled
    // Verify commit is created with AI message
}
```

## Key Files Modified

- `internal/ui/prompts.go`: Add `PromptAIMessageAcceptanceOptions()` and related functions
- `internal/ui/select_list.go`: Add `NewSelectListModelWithPreselection()`
- `internal/ui/multiline_input.go`: Add `NewMultilineInputModelWithValue()`
- `internal/service/commit_service.go`: Update `generateWithAI()` and `promptCommitMessage()`

## Key Files Created

- `internal/model/prefilled_commit_message.go`: Define `PrefilledCommitMessage` struct (if needed, or use inline)

## Dependencies

- No new external dependencies
- Uses existing `bubbletea` and `lipgloss` for interactive UI

## Success Criteria

- Users can commit with AI message in <5 seconds when selecting "accept and commit directly"
- Users can edit and commit pre-filled AI message in <30 seconds
- All three acceptance options work correctly
- Pre-filling works for all commit message fields
- Error handling works correctly (commit failures, cancellations, etc.)
