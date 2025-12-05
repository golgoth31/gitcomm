# Contract: Commit Service Extension for AI Acceptance Options

**Feature**: 009-ai-commit-options
**Date**: 2025-01-27
**Component**: `internal/service/commit_service.go`

## Overview

This contract defines the extension to the `CommitService` to handle the three-option AI message acceptance workflow.

## Method: generateWithAI (Modified)

### Current Signature

```go
func (s *CommitService) generateWithAI(ctx context.Context, stagedFiles []string) (*model.CommitMessage, error)
```

### Changes

The method is modified to:
1. Call `PromptAIMessageAcceptanceOptions()` instead of `PromptAIMessageAcceptance()`
2. Handle three acceptance paths: AcceptAndCommit, AcceptAndEdit, Reject
3. Support pre-filling commit message fields when AcceptAndEdit is selected
4. Handle rejection with option to generate new AI message or use manual input

### Behavior

#### Step 1: Generate AI Message

```go
aiMessage, err := s.generateAIMessage(ctx, stagedFiles)
if err != nil {
    return nil, fmt.Errorf("failed to generate AI message: %w", err)
}
```

#### Step 2: Prompt Acceptance Options

```go
acceptance, err := ui.PromptAIMessageAcceptanceOptions(reader)
if err != nil {
    // Handle cancellation (restore staging state)
    return nil, err
}
```

#### Step 3: Handle Acceptance Path

**If AcceptAndCommit**:
```go
// Parse AI message into CommitMessage
commitMsg, err := s.parseAIMessage(aiMessage)
if err != nil {
    return nil, fmt.Errorf("failed to parse AI message: %w", err)
}

// Validate message
if err := s.validateMessage(commitMsg); err != nil {
    // Should not happen if validation was checked before showing options
    return nil, fmt.Errorf("AI message validation failed: %w", err)
}

// Create commit immediately
err = s.gitRepo.CreateCommit(ctx, commitMsg, s.options)
if err != nil {
    // Restore staging state, show error, prompt for retry/edit/cancel
    return s.handleCommitFailure(ctx, commitMsg, err)
}

return commitMsg, nil
```

**If AcceptAndEdit**:
```go
// Parse AI message into PrefilledCommitMessage
prefilled := s.parseAIMessageToPrefilled(aiMessage)

// Call promptCommitMessage with pre-filled values
commitMsg, err := s.promptCommitMessage(ctx, &prefilled)
if err != nil {
    // Handle cancellation (restore staging state)
    return nil, err
}

// Create commit
err = s.gitRepo.CreateCommit(ctx, commitMsg, s.options)
if err != nil {
    return s.handleCommitFailure(ctx, commitMsg, err)
}

return commitMsg, nil
```

**If Reject**:
```go
// Prompt for choice: new AI or manual
useNewAI, err := ui.PromptRejectChoice(reader)
if err != nil {
    return nil, err
}

if useNewAI {
    // Recursively call generateWithAI (with retry limit to prevent infinite loop)
    return s.generateWithAI(ctx, stagedFiles)
} else {
    // Fall back to manual input with empty fields
    return s.promptCommitMessage(ctx, nil)
}
```

### Error Handling

#### Commit Failure After AcceptAndCommit

```go
func (s *CommitService) handleCommitFailure(ctx context.Context, commitMsg *model.CommitMessage, commitErr error) (*model.CommitMessage, error) {
    // Restore staging state (already handled by defer, but ensure it's done)
    // Display error message
    fmt.Printf("Error creating commit: %v\n", commitErr)

    // Prompt for retry/edit/cancel
    choice, err := ui.PromptCommitFailureChoice(reader)
    if err != nil {
        return nil, err
    }

    switch choice {
    case ui.RetryCommit:
        // Retry with same message
        err = s.gitRepo.CreateCommit(ctx, commitMsg, s.options)
        if err != nil {
            return s.handleCommitFailure(ctx, commitMsg, err) // Recursive retry
        }
        return commitMsg, nil
    case ui.EditMessage:
        // Fall back to accept and edit flow
        prefilled := s.commitMessageToPrefilled(commitMsg)
        return s.promptCommitMessage(ctx, &prefilled)
    case ui.CancelCommit:
        return nil, fmt.Errorf("commit cancelled by user")
    }
}
```

#### AI Generation Failure After Reject

```go
if useNewAI {
    newAIMessage, err := s.generateAIMessage(ctx, stagedFiles)
    if err != nil {
        // Display error and fall back to manual input
        fmt.Printf("Error generating new AI message: %v\n", err)
        fmt.Println("Falling back to manual input...")
        return s.promptCommitMessage(ctx, nil)
    }
    // Continue with new AI message (recursive call)
    return s.generateWithAI(ctx, stagedFiles)
}
```

### Validation Before Showing Options

Before calling `PromptAIMessageAcceptanceOptions()`, the AI message must be validated:

```go
// Parse and validate AI message
commitMsg, err := s.parseAIMessage(aiMessage)
if err != nil {
    // Handle parsing error (should not happen for valid AI output)
    return nil, fmt.Errorf("failed to parse AI message: %w", err)
}

// Validate message
validationErrors := s.validateMessage(commitMsg)
if len(validationErrors) > 0 {
    // Only show "accept and commit" option if validation passes
    // For now, show all options but handle validation in AcceptAndCommit path
    // OR: Only show AcceptAndEdit and Reject if validation fails
}
```

**Decision**: Show all three options regardless of validation status. If user selects "accept and commit" but validation fails, show error and fall back to edit flow.

---

## Method: promptCommitMessage (Modified)

### Current Signature

```go
func (s *CommitService) promptCommitMessage(ctx context.Context) (*model.CommitMessage, error)
```

### New Signature

```go
func (s *CommitService) promptCommitMessage(ctx context.Context, prefilled *PrefilledCommitMessage) (*model.CommitMessage, error)
```

### Parameters

- `ctx` (context.Context): Context for cancellation
- `prefilled` (*PrefilledCommitMessage): Pre-filled values (nil for empty/manual input)

### Behavior

1. **Commit Type**:
   - If `prefilled != nil && prefilled.Type != ""`: Call `PromptCommitTypeWithPreselection(reader, prefilled.Type)`
   - Otherwise: Call `PromptCommitType(reader)` (existing function)

2. **Scope**:
   - If `prefilled != nil`: Call `PromptScopeWithDefault(reader, prefilled.Scope)`
   - Otherwise: Call `PromptScope(reader)` (existing function)

3. **Subject**:
   - If `prefilled != nil && prefilled.Subject != ""`: Call `PromptSubjectWithDefault(reader, prefilled.Subject)`
   - Otherwise: Call `PromptSubject(reader)` (existing function)

4. **Body**:
   - If `prefilled != nil`: Call `PromptBodyWithDefault(reader, prefilled.Body)`
   - Otherwise: Call `PromptBody(reader)` (existing function)

5. **Footer**:
   - If `prefilled != nil`: Call `PromptFooterWithDefault(reader, prefilled.Footer)`
   - Otherwise: Call `PromptFooter(reader)` (existing function)

6. **Return**: Create and return `CommitMessage` from user input

### Error Handling

- **Cancellation**: If user cancels at any point (Ctrl+C, Escape), restore staging state and return error
- **Validation**: Existing validation logic applies to pre-filled values

---

## Helper Methods

### parseAIMessageToPrefilled

```go
func (s *CommitService) parseAIMessageToPrefilled(aiMessage string) PrefilledCommitMessage
```

Parses AI message string into `PrefilledCommitMessage` structure.

### commitMessageToPrefilled

```go
func (s *CommitService) commitMessageToPrefilled(msg *model.CommitMessage) PrefilledCommitMessage
```

Converts `CommitMessage` to `PrefilledCommitMessage` (for retry/edit flow).

---

## Testing Requirements

### Unit Tests

- `TestGenerateWithAI_AcceptAndCommit`: Test accept and commit path
- `TestGenerateWithAI_AcceptAndEdit`: Test accept and edit path
- `TestGenerateWithAI_Reject_NewAI`: Test reject with new AI generation
- `TestGenerateWithAI_Reject_Manual`: Test reject with manual input
- `TestGenerateWithAI_CommitFailure_Retry`: Test commit failure retry
- `TestGenerateWithAI_CommitFailure_Edit`: Test commit failure edit fallback
- `TestGenerateWithAI_CommitFailure_Cancel`: Test commit failure cancel
- `TestGenerateWithAI_AIGenerationFailure`: Test AI generation failure after reject
- `TestPromptCommitMessage_WithPrefilled`: Test promptCommitMessage with pre-filled values
- `TestPromptCommitMessage_Empty`: Test promptCommitMessage with nil prefilled

### Integration Tests

- `TestAIAcceptanceWorkflow_EndToEnd`: Full end-to-end workflow test

---

## Backward Compatibility

- Existing `promptCommitMessage()` callers must be updated to pass `nil` for manual input
- No breaking changes to public interfaces (internal method only)

---

## Dependencies

- `internal/ui`: For prompt functions
- `internal/model`: For `CommitMessage` and `PrefilledCommitMessage`
- `internal/repository`: For `GitRepository.CreateCommit()`
- `internal/service`: For validation and formatting services
