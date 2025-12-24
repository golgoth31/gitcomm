# Research: Unify Prompt Functions to Use Default Variants

**Date**: 2025-01-27
**Feature**: 016-unify-prompt-defaults
**Purpose**: Verify behavior of "WithDefault" prompt functions when empty strings are passed as defaults

## Research Questions

### Q1: How do "WithDefault" prompt functions behave with empty string defaults?

**Decision**: All "WithDefault" prompt functions accept empty strings as valid defaults and maintain existing validation behavior.

**Rationale**:
- `PromptSubjectWithDefault`: When `defaultValue` is empty string, the validator checks `if trimmed == "" && defaultValue == ""` and returns an error, requiring user input. This matches the behavior of `PromptSubject`.
- `PromptScopeWithDefault`: When `defaultValue` is empty, the function allows empty input (scope is optional), matching `PromptScope` behavior.
- `PromptBodyWithDefault`: When `defaultValue` is empty, the function allows empty input (body is optional), matching `PromptBody` behavior.
- `PromptFooterWithDefault`: When `defaultValue` is empty, the function allows empty input (footer is optional), matching `PromptFooter` behavior.
- `PromptCommitTypeWithPreselection`: When `preselectedType` is empty string, the form initializes with empty value and user must select a type, matching `PromptCommitType` behavior.

**Alternatives considered**:
- Creating wrapper functions that handle empty defaults differently - **Rejected**: Unnecessary complexity. The existing "WithDefault" functions already handle empty strings correctly.

### Q2: Are non-default prompt functions used elsewhere in the codebase?

**Decision**: Non-default prompt functions (`PromptScope`, `PromptSubject`, `PromptBody`, `PromptFooter`, `PromptCommitType`) are only used in `commit_service.go` and can be safely removed.

**Rationale**:
- Code search confirms these functions are only called from `promptCommitMessage` in `internal/service/commit_service.go`.
- No test files directly import or test these functions independently.
- No other service or package references these functions.

**Alternatives considered**:
- Deprecating functions with comments - **Rejected**: Spec clarification confirmed removal is preferred.
- Keeping functions for backward compatibility - **Rejected**: No external callers exist.

### Q3: What is the correct default value for PromptCommitTypeWithPreselection when no type is pre-filled?

**Decision**: Pass empty string (`""`) as the `preselectedType` parameter when no pre-filled type is available.

**Rationale**:
- `PromptCommitTypeWithPreselection` accepts a string parameter and initializes `commitType := preselectedType`.
- When `preselectedType` is empty, the form starts with an empty value and the user must select a type from the options.
- The function iterates through options to mark a match as selected, but if no match is found (empty string), no option is pre-selected, which is the desired behavior.
- This matches the behavior of `PromptCommitType` which has no pre-selection.

**Alternatives considered**:
- Passing first option as default - **Rejected**: Would pre-select a type when user may want a different one, changing behavior.
- Creating a new function - **Rejected**: Unnecessary when existing function handles empty string correctly.

## Implementation Notes

1. **Refactoring Strategy**:
   - Replace all conditional branches in `promptCommitMessage` with direct calls to "WithDefault" variants.
   - When `prefilled == nil` or field is empty, pass empty string as default.
   - When `prefilled != nil` and field has value, pass that value as default.

2. **Function Removal**:
   - Remove `PromptScope`, `PromptSubject`, `PromptBody`, `PromptFooter`, `PromptCommitType` from `prompts.go`.
   - Verify no other references exist before removal (already confirmed via code search).

3. **Testing Strategy**:
   - Update existing tests that may reference removed functions.
   - Add tests verifying behavior with empty string defaults matches original behavior.
   - Ensure all existing integration tests pass.

## References

- `internal/ui/prompts.go`: Implementation of all prompt functions
- `internal/service/commit_service.go`: Current usage of prompt functions
- Spec clarification: Non-default functions should be removed entirely
