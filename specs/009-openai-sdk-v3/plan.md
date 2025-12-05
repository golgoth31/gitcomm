# Implementation Plan: Upgrade OpenAI Provider to SDK v3

**Branch**: `009-openai-sdk-v3` | **Date**: 2025-01-27 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/009-openai-sdk-v3/spec.md`

## Summary

This feature upgrades the OpenAI provider implementation from SDK v1 (`github.com/openai/openai-go`) to SDK v3 (`github.com/openai/openai-go/v3`). The upgrade maintains 100% backward compatibility with existing functionality, interfaces, and configuration while leveraging the latest SDK features, bug fixes, and improvements. The technical approach involves updating the import path, adapting to any SDK v3 API changes, and ensuring error handling maps correctly to existing patterns.

## Technical Context

**Language/Version**: Go 1.25.0+
**Primary Dependencies**:
- `github.com/openai/openai-go/v3` - OpenAI official Go SDK v3 (replaces v1.12.0)
- Standard Go libraries (`context`, `fmt`, `strings`, `errors`)
- Existing dependencies: `github.com/golgoth31/gitcomm/internal/model`, `github.com/golgoth31/gitcomm/internal/utils`

**Storage**: N/A (no data persistence, API calls only)

**Testing**:
- Standard Go testing framework (`testing` package)
- `github.com/onsi/ginkgo/v2` and `github.com/onsi/gomega` for BDD-style tests (existing)
- Unit tests for SDK v3 integration (mocking SDK clients)
- Integration tests for SDK v3 API calls (with test API keys or mocked responses)
- Existing provider tests must pass with minimal updates to match SDK v3 API

**Target Platform**: Linux, macOS, Windows (CLI application)

**Project Type**: CLI tool (single binary)

**Performance Goals**:
- SDK v3 initialization completes in < 10ms (same as SDK v1)
- API calls complete within configured timeout (default 30 seconds, same as current)
- No measurable performance regression compared to SDK v1 implementation
- SDK v3 automatic retries should improve reliability without significant latency increase

**Constraints**:
- Must maintain backward compatibility with existing AIProvider interface (no breaking changes)
- Must maintain backward compatibility with existing AIProviderConfig structure (no changes)
- Must preserve all existing error handling behavior (same error types and messages)
- Must not expose SDK-specific implementation details in error messages or logs
- Must handle SDK v3 initialization failures gracefully (fail fast with clear error, fallback to manual input)
- Must map SDK v3 error types to existing error handling patterns (unmappable errors wrapped generically)
- Must respect existing timeout configurations from user config
- Must maintain same prompt building and response parsing logic

**Scale/Scope**:
- Single repository per CLI invocation
- Handles typical repository sizes (hundreds to thousands of files)
- No concurrent API calls (single request per commit generation)
- Supports OpenAI provider only (other providers unchanged)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Verify compliance with gitcomm Constitution principles:

- **Clean Architecture**: ✅ The OpenAI provider is in `internal/ai/` following Clean Architecture. No repository pattern needed (no data persistence). Layer separation maintained.
- **Interface-Driven Development**: ✅ OpenAI provider implements `AIProvider` interface. Dependencies injected via constructor (`NewOpenAIProvider`). No global state.
- **Test-First Development**: ✅ TDD approach required - tests must be written before implementation. Unit and integration tests planned.
- **Idiomatic Go**: ✅ Follows Go conventions. Naming conventions appropriate (PascalCase for exported, camelCase for unexported). Uses `gofmt`/`goimports`.
- **Error Handling**: ✅ Explicit error handling with wrapped errors. Custom error types (`ErrAIProviderUnavailable`) used. No panics in library code.
- **Context & Thread Safety**: ✅ Uses `context.Context` for cancellation/timeout. No goroutines needed (single-threaded API calls).
- **Technical Constraints**: ✅ No global state. Graceful error handling. Resource cleanup via context cancellation.
- **Operational Constraints**: ✅ Logging via `zerolog` (debug mode). Secrets (API keys) not exposed in logs.

**Violations**: None - all principles satisfied.

**Post-Design Re-check** (after Phase 1):

- **Clean Architecture**: ✅ Maintained - provider in `internal/ai/`, no structural changes
- **Interface-Driven Development**: ✅ Maintained - `AIProvider` interface unchanged, dependency injection via constructor
- **Test-First Development**: ✅ Maintained - TDD approach required, tests planned before implementation
- **Idiomatic Go**: ✅ Maintained - Go conventions followed, naming appropriate
- **Error Handling**: ✅ Maintained - Explicit error handling, custom error types, error mapping strategy defined
- **Context & Thread Safety**: ✅ Maintained - Context usage for cancellation/timeout, thread-safe SDK client
- **Technical Constraints**: ✅ Maintained - No global state, graceful error handling, resource cleanup
- **Operational Constraints**: ✅ Maintained - Logging via zerolog, secrets not exposed

**Post-Design Violations**: None - all principles continue to be satisfied after design phase.

## Project Structure

### Documentation (this feature)

```text
specs/009-openai-sdk-v3/
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
    ├── openai_provider.go         # OpenAI provider implementation (upgrade to SDK v3)
    ├── openai_provider_test.go   # Unit tests for OpenAI provider
    ├── provider.go                # AIProvider interface (unchanged)
    └── ...

test/
└── integration/
    └── ai_commit_test.go          # Integration tests (may need updates for SDK v3)
```

**Structure Decision**: Single project structure. The OpenAI provider upgrade is isolated to `internal/ai/openai_provider.go` and its tests. No new directories needed. Existing structure maintained.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

No violations - all constitution principles satisfied.
