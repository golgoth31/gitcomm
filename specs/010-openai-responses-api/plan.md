# Implementation Plan: Migrate OpenAI Provider to Responses API

**Branch**: `010-openai-responses-api` | **Date**: 2025-01-27 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/010-openai-responses-api/spec.md`

## Summary

Migrate the OpenAI provider from Chat Completions API to Responses API while maintaining 100% backward compatibility. The migration involves updating the API endpoint from `/v1/chat/completions` to `/v1/responses`, converting the `messages` array to `input` parameter format, and adapting response extraction logic. The implementation will use OpenAI SDK v3 if available, or fall back to custom HTTP client if SDK support is not yet available. All existing functionality, error handling, and user-facing behavior must remain identical.

## Technical Context

**Language/Version**: Go 1.25.0+
**Primary Dependencies**:
- `github.com/openai/openai-go/v3` (SDK v3) - for Responses API support (or custom HTTP client if SDK doesn't support it yet)
- Existing: `github.com/golgoth31/gitcomm/internal/model`, `github.com/golgoth31/gitcomm/internal/utils`
**Storage**: N/A (stateless API calls)
**Testing**: Go `testing` package, table-driven tests, integration tests in `test/integration/`
**Target Platform**: Linux/macOS/Windows (CLI tool)
**Project Type**: CLI tool (single binary)
**Performance Goals**: Same response time as current Chat Completions implementation (no degradation)
**Constraints**:
- Must maintain 100% backward compatibility with existing configurations
- Must preserve identical error handling behavior
- Must use stateless mode (disable conversation state management)
- Must support same models as current implementation
**Scale/Scope**: Single user, single request per execution (CLI tool)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Verify compliance with gitcomm Constitution principles:

- **Clean Architecture**: ✅ Yes - Changes are isolated to `internal/ai/openai_provider.go`. Repository Pattern not applicable (no data persistence). Layer separation maintained (internal/ai for provider logic).
- **Interface-Driven Development**: ✅ Yes - `OpenAIProvider` implements `AIProvider` interface. Dependencies injected via constructor (`NewOpenAIProvider`). No global state.
- **Test-First Development**: ✅ Yes - TDD approach required. Tests must be written before implementation. Existing test structure in `internal/ai/openai_provider_test.go` will be updated.
- **Idiomatic Go**: ✅ Yes - Follows Go conventions. Naming: `OpenAIProvider`, `GenerateCommitMessage`, `mapSDKError`. Error handling with wrapped errors.
- **Error Handling**: ✅ Yes - Custom error types via `utils.ErrAIProviderUnavailable`. Error mapping function `mapSDKError` wraps SDK errors. Explicit error checking throughout.
- **Context & Thread Safety**: ✅ Yes - Uses `context.Context` for cancellation/timeout. No goroutines needed (single-threaded CLI). Context propagated to API calls.
- **Technical Constraints**: ✅ Yes - No global state. Dependencies injected. Resource cleanup via context cancellation. Error handling explicit.
- **Operational Constraints**: ✅ Yes - Logging via `zerolog` (debug mode). Secrets (API keys) not logged. Configuration via `AIProviderConfig` struct.

**Violations**: None - All principles are followed.

## Project Structure

### Documentation (this feature)

```text
specs/010-openai-responses-api/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
internal/
└── ai/
    ├── openai_provider.go          # Main implementation (modify)
    └── openai_provider_test.go     # Unit tests (modify)

test/
└── integration/
    └── ai_commit_test.go           # Integration tests (verify compatibility)
```

**Structure Decision**: Single project structure. Changes are isolated to the OpenAI provider implementation. No new files needed - only modifications to existing `openai_provider.go` and its tests.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

No violations - all principles followed.
