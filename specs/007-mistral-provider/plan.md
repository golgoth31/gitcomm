# Implementation Plan: Add Mistral as AI Provider

**Branch**: `007-mistral-provider` | **Date**: 2025-01-27 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/007-mistral-provider/spec.md`

## Summary

This feature adds Mistral AI as a new provider option alongside OpenAI and Anthropic. Users can configure Mistral API credentials, select it as their provider, and use it to generate commit messages following the same patterns as existing providers.

The technical approach follows the existing provider implementation pattern: create a `MistralProvider` struct that implements the `AIProvider` interface, integrate it into the provider selection mechanism in `CommitService`, add Mistral to token calculation infrastructure, and update configuration examples. This is a straightforward extension that follows established patterns with no architectural changes.

## Technical Context

**Language/Version**: Go 1.25.0+
**Primary Dependencies**:
- `github.com/go-git/go-git/v5` - Git operations (existing)
- `net/http` - HTTP client for Mistral API (standard library)
- `encoding/json` - JSON encoding/decoding (standard library)
- `context` - Context propagation (standard library)
- `time` - Timeout management (standard library)
- Standard Go libraries only (no new external dependencies)

**Storage**: N/A (no data persistence, API calls only)

**Testing**:
- Standard Go testing framework (`testing` package)
- `github.com/onsi/ginkgo/v2` and `github.com/onsi/gomega` for BDD-style tests (existing)
- Unit tests for MistralProvider implementation
- Integration tests for Mistral API integration (with mocked HTTP responses)
- No external testing dependencies

**Target Platform**: Linux, macOS, Windows (CLI application)

**Project Type**: CLI tool (single binary)

**Performance Goals**:
- Mistral API calls complete within configured timeout (default 30 seconds)
- Token calculation completes in < 10ms (same as existing providers)
- No measurable performance impact on CLI startup or workflow

**Constraints**:
- Must maintain backward compatibility with existing providers (OpenAI, Anthropic, local)
- Must follow existing provider patterns (no architectural changes)
- Must support configurable API endpoint (default to standard Mistral API)
- Must support configurable model selection (default to appropriate Mistral model)
- Must handle API errors gracefully with fallback to manual input
- Must not expose API keys in logs or error messages

**Scale/Scope**:
- Single repository per CLI invocation
- Handles typical repository sizes (hundreds to thousands of files)
- No concurrent API calls (single request per commit generation)
- Supports all Mistral API models that follow chat completion pattern

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Verify compliance with gitcomm Constitution principles:

- **Clean Architecture**: ✅ Follows existing layer separation - new provider in `internal/ai/` layer, no new layers needed
- **Interface-Driven Development**: ✅ Implements existing `AIProvider` interface, uses dependency injection via constructor
- **Test-First Development**: ✅ Tests will be written before implementation (TDD approach for provider logic)
- **Idiomatic Go**: ✅ Uses standard library HTTP client, follows Go naming conventions, matches existing provider patterns
- **Error Handling**: ✅ Uses existing error handling patterns, wraps errors for traceability, no panics
- **Context & Thread Safety**: ✅ Uses `context.Context` for cancellation/timeout, HTTP client respects context, no concurrency needed
- **Technical Constraints**: ✅ No global state, HTTP client properly closed, API keys not exposed in logs
- **Operational Constraints**: ✅ Uses existing logging infrastructure, API keys from config (not hardcoded), no secrets in logs

**Violations**: None - this feature fully complies with all constitution principles.

### Post-Design Constitution Check

*Re-evaluated after Phase 1 design artifacts created.*

After completing research, data model, and contracts:

- **Clean Architecture**: ✅ Confirmed - new provider in `internal/ai/` layer, follows existing patterns, no architectural changes
- **Interface-Driven Development**: ✅ Confirmed - implements existing `AIProvider` interface, uses dependency injection
- **Test-First Development**: ✅ Confirmed - test strategy defined in contracts, TDD approach maintained
- **Idiomatic Go**: ✅ Confirmed - uses standard library HTTP client, follows Go conventions, matches existing provider code
- **Error Handling**: ✅ Confirmed - uses existing error patterns, wraps errors, no panics, graceful fallback
- **Context & Thread Safety**: ✅ Confirmed - uses `context.Context` for cancellation/timeout, HTTP client respects context, no concurrency
- **Technical Constraints**: ✅ Confirmed - no global state, HTTP client properly managed, API keys not exposed
- **Operational Constraints**: ✅ Confirmed - uses existing logging, API keys from config, no secrets in logs

**Post-Design Violations**: None - design artifacts confirm full compliance with all principles.

## Project Structure

### Documentation (this feature)

```text
specs/007-mistral-provider/
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
├── openai_provider.go   # OpenAI implementation (existing, no changes)
├── anthropic_provider.go # Anthropic implementation (existing, no changes)
├── mistral_provider.go  # NEW: Mistral implementation
└── mistral_provider_test.go # NEW: Mistral unit tests

internal/service/
└── commit_service.go    # Modify: Add "mistral" case to provider switch

pkg/tokenization/
└── token_calculator.go  # Modify: Add "mistral" case (use fallback or Mistral tokenization)

configs/
└── config.yaml.example  # Modify: Add mistral provider example

test/integration/
└── ai_mistral_test.go   # NEW: Integration tests for Mistral provider
```

**Structure Decision**: This is a straightforward extension to existing AI provider infrastructure. No new packages or modules needed. The change adds a new provider implementation following the same pattern as OpenAI and Anthropic, with modifications to provider selection and token calculation.

## Complexity Tracking

> **No violations - feature fully complies with constitution principles**
