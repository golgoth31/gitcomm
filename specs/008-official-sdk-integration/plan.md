# Implementation Plan: Use Official SDKs for AI Providers

**Branch**: `008-official-sdk-integration` | **Date**: 2025-01-27 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/008-official-sdk-integration/spec.md`

## Summary

This feature replaces the current HTTP client implementations for OpenAI, Anthropic, and Mistral AI providers with their respective official Go SDKs. The implementation maintains 100% backward compatibility with existing functionality, interfaces, and configuration while leveraging official SDK features like automatic retries and better error handling. The AIProvider interface remains unchanged, ensuring no breaking changes to the public API.

The technical approach involves refactoring each provider implementation to use the official SDK instead of raw HTTP calls, mapping SDK-specific error types to existing error handling patterns, and extending AIProviderConfig with optional SDK-specific fields if needed while maintaining backward compatibility.

## Technical Context

**Language/Version**: Go 1.25.0+
**Primary Dependencies**:
- `github.com/openai/openai-go` - OpenAI official Go SDK
- `github.com/anthropics/anthropic-sdk-go` - Anthropic official Go SDK
- `github.com/Gage-Technologies/mistral-go` - Mistral official Go SDK
- Standard Go libraries (`context`, `time`, `fmt`, `errors`)
- Existing dependencies: `github.com/golgoth31/gitcomm/internal/model`, `github.com/golgoth31/gitcomm/internal/utils`

**Storage**: N/A (no data persistence, API calls only)

**Testing**:
- Standard Go testing framework (`testing` package)
- `github.com/onsi/ginkgo/v2` and `github.com/onsi/gomega` for BDD-style tests (existing)
- Unit tests for SDK integration (mocking SDK clients)
- Integration tests for SDK API calls (with test API keys or mocked responses)
- Existing provider tests must pass without modification

**Target Platform**: Linux, macOS, Windows (CLI application)

**Project Type**: CLI tool (single binary)

**Performance Goals**:
- SDK initialization completes in < 10ms (same as HTTP client initialization)
- API calls complete within configured timeout (default 30 seconds, same as current)
- No measurable performance regression compared to HTTP client implementation
- SDK automatic retries should improve reliability without significant latency increase

**Constraints**:
- Must maintain backward compatibility with existing AIProvider interface (no breaking changes)
- Must maintain backward compatibility with existing AIProviderConfig structure (optional fields only)
- Must preserve all existing error handling behavior (same error types and messages)
- Must not expose SDK-specific implementation details in error messages or logs
- Must handle SDK initialization failures gracefully (fail fast with clear error, fallback to manual input)
- Must map SDK-specific error types to existing error handling patterns
- Must respect existing timeout configurations from user config
- Must maintain same prompt building and response parsing logic

**Scale/Scope**:
- Single repository per CLI invocation
- Handles typical repository sizes (hundreds to thousands of files)
- No concurrent API calls (single request per commit generation)
- Supports all three providers (OpenAI, Anthropic, Mistral)
- Local provider remains unchanged (no SDK replacement)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Verify compliance with gitcomm Constitution principles:

- **Clean Architecture**: ✅ Follows existing layer separation - provider implementations remain in `internal/ai/` layer, no new layers needed. Repository Pattern not applicable (no data persistence).

- **Interface-Driven Development**: ✅ Maintains existing `AIProvider` interface contract, uses dependency injection via constructors. SDK clients are injected/internal to provider implementations, not exposed.

- **Test-First Development**: ✅ Tests will be written before implementation (TDD approach). Existing provider tests must pass without modification, new tests for SDK integration will be added.

- **Idiomatic Go**: ✅ Uses official Go SDKs following Go conventions, matches existing provider patterns. Naming conventions follow Go standards.

- **Error Handling**: ✅ Uses existing error handling patterns, maps SDK errors to existing error types, wraps errors for traceability, no panics.

- **Context & Thread Safety**: ✅ Uses `context.Context` for cancellation/timeout (SDKs support context), no concurrency needed (single request per invocation).

- **Technical Constraints**: ✅ No global state, SDK clients properly initialized and managed, API keys not exposed in logs, graceful error handling with fallback.

- **Operational Constraints**: ✅ Uses existing logging infrastructure (`zerolog`), API keys from config (not hardcoded), no secrets in logs, SDK initialization errors logged appropriately.

**Violations**: None - this feature fully complies with all constitution principles.

### Post-Design Constitution Check

*Re-evaluated after Phase 1 design artifacts created.*

After completing research, data model, and contracts:

- **Clean Architecture**: ✅ Confirmed - provider implementations in `internal/ai/` layer, no architectural changes, maintains layer separation.

- **Interface-Driven Development**: ✅ Confirmed - maintains existing `AIProvider` interface, uses dependency injection, SDK clients are internal implementation details.

- **Test-First Development**: ✅ Confirmed - test strategy defined, TDD approach maintained, existing tests must pass, new SDK integration tests planned.

- **Idiomatic Go**: ✅ Confirmed - uses official Go SDKs, follows Go conventions, matches existing patterns, proper error handling.

- **Error Handling**: ✅ Confirmed - maps SDK errors to existing error types, preserves user-facing messages, wraps errors appropriately, no panics.

- **Context & Thread Safety**: ✅ Confirmed - uses `context.Context` for cancellation/timeout, SDKs support context, no concurrency needed.

- **Technical Constraints**: ✅ Confirmed - no global state, SDK clients properly managed, API keys not exposed, graceful error handling.

- **Operational Constraints**: ✅ Confirmed - uses existing logging, API keys from config, no secrets in logs, appropriate error logging.

**Post-Design Violations**: None - design artifacts confirm full compliance with all principles.

## Project Structure

### Documentation (this feature)

```text
specs/008-official-sdk-integration/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
internal/ai/
├── provider.go          # AIProvider interface (existing, no changes)
├── openai_provider.go   # MODIFY: Replace HTTP client with OpenAI SDK
├── anthropic_provider.go # MODIFY: Replace HTTP client with Anthropic SDK
├── mistral_provider.go  # MODIFY: Replace HTTP client with Mistral SDK
├── local_provider.go    # UNCHANGED: No SDK replacement needed
├── openai_provider_test.go   # MODIFY: Update tests for SDK integration
├── anthropic_provider_test.go # MODIFY: Update tests for SDK integration
└── mistral_provider_test.go  # MODIFY: Update tests for SDK integration

internal/model/
└── config.go           # MODIFY: Extend AIProviderConfig with optional SDK-specific fields if needed

pkg/tokenization/
└── token_calculator.go # UNCHANGED: Token calculation logic remains the same

internal/service/
└── commit_service.go   # UNCHANGED: Provider selection mechanism remains the same
```

**Structure Decision**: This is a refactoring feature that modifies existing provider implementations without changing the overall architecture. The changes are isolated to the `internal/ai/` package and `internal/model/config.go` (if config extension is needed). No new packages or modules are required. The feature maintains the existing Clean Architecture layer separation.

## Complexity Tracking

> **No violations - feature fully complies with constitution principles**
