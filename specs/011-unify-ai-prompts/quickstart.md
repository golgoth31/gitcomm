# Quickstart: Unify AI Provider Prompts with Validation Rules

**Date**: 2025-01-27
**Feature**: 011-unify-ai-prompts

## Overview

This feature unifies all AI provider prompts to use identical system and user messages that include validation rules extracted dynamically from MessageValidator. This ensures consistent commit message generation across all providers.

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    PromptGenerator                          │
│                  (pkg/ai/prompt/)                           │
│                                                             │
│  GenerateSystemMessage(validator) → System Message         │
│  GenerateUserMessage(repoState) → User Message              │
└─────────────────────────────────────────────────────────────┘
                            │
                            │ uses
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                  MessageValidator                           │
│              (pkg/conventional/)                            │
│                                                             │
│  GetValidTypes() → ["feat", "fix", ...]                     │
│  GetSubjectMaxLength() → 72                                │
│  GetBodyMaxLength() → 320                                   │
│  GetScopeFormatDescription() → "alphanumeric, ..."          │
└─────────────────────────────────────────────────────────────┘
                            │
                            │ used by
                            ▼
┌─────────────────────────────────────────────────────────────┐
│              AI Provider Implementations                    │
│                  (internal/ai/)                             │
│                                                             │
│  OpenAIProvider    → Uses system + user separately         │
│  AnthropicProvider → Prepends system to user                │
│  MistralProvider   → Uses system + user separately          │
│  LocalProvider     → Uses system + user separately          │
└─────────────────────────────────────────────────────────────┘
```

## Key Components

### 1. PromptGenerator (`pkg/ai/prompt/generator.go`)

**Purpose**: Generate unified prompts with validation rules

**Usage**:
```go
generator := prompt.NewUnifiedPromptGenerator()
validator := conventional.NewValidator()

// Generate system message with validation rules
systemMsg, err := generator.GenerateSystemMessage(validator)
if err != nil {
    return err
}

// Generate user message with repository state
userMsg, err := generator.GenerateUserMessage(repoState)
if err != nil {
    return err
}
```

### 2. MessageValidator Extension (`pkg/conventional/validator.go`)

**Purpose**: Expose validation rules programmatically

**New Methods**:
```go
validator := conventional.NewValidator()

types := validator.GetValidTypes()
// Returns: ["feat", "fix", "docs", "style", "refactor", "test", "chore", "version"]

maxSubjectLen := validator.GetSubjectMaxLength()
// Returns: 72

maxBodyLen := validator.GetBodyMaxLength()
// Returns: 320

scopeDesc := validator.GetScopeFormatDescription()
// Returns: "alphanumeric, hyphens, underscores only"
```

### 3. Provider Updates (`internal/ai/*_provider.go`)

**OpenAI/Mistral/Local Providers**:
```go
systemMsg, _ := generator.GenerateSystemMessage(validator)
userMsg, _ := generator.GenerateUserMessage(repoState)

// Use as separate messages in API call
inputItems := []responses.ResponseInputItemUnionParam{
    {
        OfMessage: &responses.EasyInputMessageParam{
            Role: responses.EasyInputMessageRoleSystem,
            Content: responses.EasyInputMessageContentUnionParam{
                OfString: openai.String(systemMsg),
            },
        },
    },
    {
        OfMessage: &responses.EasyInputMessageParam{
            Role: responses.EasyInputMessageRoleUser,
            Content: responses.EasyInputMessageContentUnionParam{
                OfString: openai.String(userMsg),
            },
        },
    },
}
```

**Anthropic Provider**:
```go
systemMsg, _ := generator.GenerateSystemMessage(validator)
userMsg, _ := generator.GenerateUserMessage(repoState)

// Combine system and user messages (Anthropic doesn't support system messages)
combinedMsg := systemMsg + "\n\n" + userMsg

req := anthropic.MessageNewParams{
    Model: anthropic.Model(modelName),
    Messages: []anthropic.MessageParam{
        {
            Role: anthropic.MessageParamRoleUser,
            Content: []anthropic.ContentBlockParamUnion{
                {
                    OfText: &anthropic.TextBlockParam{
                        Text: combinedMsg,
                    },
                },
            },
        },
    },
    MaxTokens: int64(maxTokens),
}
```

## Implementation Steps

1. **Extend MessageValidator Interface**
   - Add getter methods to `MessageValidator` interface
   - Implement methods in `Validator` struct
   - Add unit tests

2. **Create PromptGenerator**
   - Create `pkg/ai/prompt/generator.go`
   - Implement `PromptGenerator` interface
   - Implement `GenerateSystemMessage()` and `GenerateUserMessage()`
   - Add unit tests

3. **Update Provider Implementations**
   - Update `OpenAIProvider` to use unified generator
   - Update `AnthropicProvider` to use unified generator (prepend system to user)
   - Update `MistralProvider` to use unified generator
   - Update `LocalProvider` to use unified generator
   - Remove old `buildPrompt()` methods

4. **Add Integration Tests**
   - Test prompt consistency across all providers
   - Test Anthropic system/user combination
   - Test validation rule extraction

## Testing

### Unit Tests

```go
// Test PromptGenerator
func TestGenerateSystemMessage(t *testing.T) {
    generator := prompt.NewUnifiedPromptGenerator()
    validator := conventional.NewValidator()

    systemMsg, err := generator.GenerateSystemMessage(validator)
    assert.NoError(t, err)
    assert.Contains(t, systemMsg, "feat, fix, docs")
    assert.Contains(t, systemMsg, "≤72 characters")
}

func TestGenerateUserMessage(t *testing.T) {
    generator := prompt.NewUnifiedPromptGenerator()
    repoState := &model.RepositoryState{
        StagedFiles: []model.FileChange{
            {Path: "test.go", Status: "modified", Diff: "..."},
        },
    }

    userMsg, err := generator.GenerateUserMessage(repoState)
    assert.NoError(t, err)
    assert.Contains(t, userMsg, "test.go")
}
```

### Integration Tests

```go
// Test prompt consistency
func TestPromptConsistency(t *testing.T) {
    generator := prompt.NewUnifiedPromptGenerator()
    validator := conventional.NewValidator()
    repoState := createTestRepoState()

    systemMsg1, _ := generator.GenerateSystemMessage(validator)
    systemMsg2, _ := generator.GenerateSystemMessage(validator)
    assert.Equal(t, systemMsg1, systemMsg2)

    userMsg1, _ := generator.GenerateUserMessage(repoState)
    userMsg2, _ := generator.GenerateUserMessage(repoState)
    assert.Equal(t, userMsg1, userMsg2)
}
```

## Verification

After implementation, verify:

1. ✅ All providers use identical system messages
2. ✅ All providers use identical user messages
3. ✅ Anthropic prepends system to user correctly
4. ✅ Validation rules are extracted correctly
5. ✅ Generated prompts include all validation constraints
6. ✅ Backward compatibility maintained (AIProvider interface unchanged)

## Migration Notes

- **No Breaking Changes**: AIProvider interface remains unchanged
- **Backward Compatible**: Existing code continues to work
- **Gradual Migration**: Providers can be updated one at a time
- **Testing**: All existing tests should continue to pass
