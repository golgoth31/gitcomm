# Data Model: Display Commit Type Selection Confirmation

**Feature**: 006-type-confirmation-display
**Date**: 2025-01-27

## Domain Entities

### No New Entities Required

This feature is a display-only enhancement that does not introduce new domain entities. The feature uses existing data structures:

- **Commit Type String**: The selected commit type (e.g., "feat", "fix") is already returned as a `string` from `SelectListModel.GetSelectedType()`
- **No State Management**: The confirmation is a one-time display with no persistent state
- **No Data Transformation**: The chosen type is displayed exactly as selected, with no modification

## Data Flow

### Selection to Display Flow

```
[User selects commit type in SelectListModel]
  ↓
[SelectListModel.GetSelectedType() returns string]
  ↓
[PromptCommitType displays confirmation: "✔ Choose a type(<scope>): <type>"]
  ↓
[PromptCommitType returns selectedType string]
  ↓
[Calling code receives selectedType and continues workflow]
```

## Validation Rules

- **Type Value**: Must be one of the predefined commit types (feat, fix, docs, style, refactor, test, chore, version)
- **Display Format**: Must match exact format: `"✔ Choose a type(<scope>): %s\n"` where `%s` is the selected type
- **Timing**: Must be displayed after alt-screen clears and before next prompt

## State Transitions

### Display State Machine

```
[Selection Screen Active]
  ↓ (User selects type and presses Enter)
[Alt-Screen Clears]
  ↓
[Confirmation Line Displayed] ← New state in this feature
  ↓
[Return to Calling Code]
  ↓
[Next Prompt (Scope) Appears]
```

**Note**: If user cancels (Escape), the confirmation state is never reached - function returns error immediately.

## Relationships

- **Uses**: `SelectListModel` (existing) - provides selected type value
- **Used By**: `PromptCommitType` function (modified) - displays confirmation
- **No Dependencies**: No new entities depend on this feature
