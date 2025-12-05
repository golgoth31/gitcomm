# Implementation Plan: Unify AI Provider Prompts with Validation Rules

**Branch**: `011-unify-ai-prompts` | **Date**: 2025-01-27 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/011-unify-ai-prompts/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Unify all AI provider prompts (OpenAI, Anthropic, Mistral, local) to use identical system and user messages that include validation rules extracted dynamically from MessageValidator. This ensures consistent commit message generation across all providers and guarantees that AI-generated messages pass validation. The unified prompt will be generated programmatically from the validator to maintain sync with validation rules.

## Technical Context

**Language/Version**: Go 1.25.0+
**Primary Dependencies**:
- Existing: `github.com/openai/openai-go/v3`, `github.com/anthropics/anthropic-sdk-go`, `github.com/gage-technologies/mistral-go`
- No new external dependencies required
**Storage**: N/A (in-memory prompt generation)
**Testing**: Go `testing` package with table-driven tests, existing test infrastructure
**Target Platform**: Linux/macOS/Windows (CLI tool)
**Project Type**: Single CLI application
**Performance Goals**: Prompt generation should be <1ms (in-memory string operations)
**Constraints**:
- Must maintain backward compatibility with existing AIProvider interface
- Must not break existing provider implementations
- Prompt generation must be thread-safe
**Scale/Scope**:
- 4 AI providers (OpenAI, Anthropic, Mistral, local)
- Single unified prompt generator shared across all providers
- Dynamic extraction of validation rules from MessageValidator

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Verify compliance with gitcomm Constitution principles:

- **Clean Architecture**: ✅ **COMPLIANT**
  - Unified prompt generator will be placed in `pkg/ai/prompt/` (shared utility)
  - Provider implementations in `internal/ai/` will use the shared generator
  - Clear separation: prompt generation (pkg) → provider usage (internal)

- **Interface-Driven Development**: ✅ **COMPLIANT**
  - Prompt generator will implement an interface for testability
  - Dependencies injected via constructors (MessageValidator, RepositoryState)
  - No global state

- **Test-First Development**: ✅ **COMPLIANT**
  - TDD approach: tests for prompt generator first
  - Table-driven tests for validation rule extraction
  - Integration tests for provider usage

- **Idiomatic Go**: ✅ **COMPLIANT**
  - Follows Go naming conventions
  - Small, focused functions
  - Proper error handling

- **Error Handling**: ✅ **COMPLIANT**
  - Explicit error handling for rule extraction failures
  - Wrapped errors for traceability
  - No panics

- **Context & Thread Safety**: ✅ **COMPLIANT**
  - Prompt generation is stateless (no shared mutable state)
  - Thread-safe string operations
  - No goroutines required

- **Technical Constraints**: ✅ **COMPLIANT**
  - No global state
  - Stateless prompt generation
  - Resource cleanup not applicable (in-memory operations)

- **Operational Constraints**: ✅ **COMPLIANT**
  - No new logging requirements
  - No secrets involved
  - Error messages don't expose internal details

**Violations**: None. All principles are satisfied.

## Project Structure

### Documentation (this feature)

```text
specs/011-unify-ai-prompts/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
pkg/
└── ai/
    └── prompt/
        ├── generator.go          # NEW: Unified prompt generator interface and implementation
        └── generator_test.go     # NEW: Tests for prompt generation

internal/
└── ai/
    ├── openai_provider.go        # MODIFY: Use unified prompt generator
    ├── anthropic_provider.go    # MODIFY: Use unified prompt generator (prepend system to user)
    ├── mistral_provider.go       # MODIFY: Use unified prompt generator
    ├── local_provider.go         # MODIFY: Use unified prompt generator
    └── provider.go               # UNCHANGED: Interface remains the same

pkg/
└── conventional/
    └── validator.go              # MODIFY: Add methods to extract validation rules programmatically

test/
└── integration/
    └── prompt_unification_test.go # NEW: Integration tests for prompt consistency
```

**Structure Decision**: Single project structure. The unified prompt generator is placed in `pkg/ai/prompt/` as a shared utility that can be used by all providers. Provider implementations in `internal/ai/` are modified to use the generator. The MessageValidator in `pkg/conventional/` is extended with methods to extract validation rules programmatically.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

No violations. All principles satisfied.
