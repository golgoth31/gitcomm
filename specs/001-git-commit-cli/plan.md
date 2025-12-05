# Implementation Plan: Git Commit Message Automation CLI

**Branch**: `001-git-commit-cli` | **Date**: 2025-01-27 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-git-commit-cli/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Build a CLI tool that automates git commit message creation following the Conventional Commits specification. The tool guides users through creating properly formatted commit messages either manually or with AI assistance. It supports multiple AI providers (OpenAI, Anthropic, local models), calculates token usage, validates message format, and integrates seamlessly with git workflows. The implementation uses Go with Clean Architecture, TUI for interactive prompts, and follows the gitcomm constitution principles.

## Technical Context

**Language/Version**: Go 1.25.0+
**Primary Dependencies**:
  - `github.com/spf13/cobra` - CLI framework for command parsing and options
  - `github.com/charmbracelet/bubbletea` - TUI framework for interactive prompts
  - `github.com/charmbracelet/lipgloss` - TUI styling
  - `github.com/spf13/viper` - Configuration management
  - `github.com/go-git/go-git/v5` - Git operations (porcelain commands)
  - `github.com/rs/zerolog` - Structured JSON logging
  - `github.com/onsi/ginkgo` - Testing framework
  - `github.com/onsi/gomega` - Testing assertions
  - Tokenization libraries: tiktoken (Go bindings) for OpenAI, custom for Anthropic

**Storage**: Configuration file (YAML/JSON) at `~/.gitcomm/config.yaml` for AI provider credentials and settings. No persistent data storage required.

**Testing**: Standard Go testing framework (`testing` package) with Ginkgo/Gomega for BDD-style tests. Integration tests verify git operations and AI provider interactions.

**Target Platform**: Cross-platform CLI (Linux, macOS, Windows) - single binary distribution

**Project Type**: Single CLI application

**Performance Goals**:
  - CLI responds to user input within 100ms for all interactive prompts (SC-005)
  - Token calculation completes in <500ms for typical repository states
  - AI provider API calls with timeout handling (default 30s)

**Constraints**:
  - Must work offline (manual mode) without AI provider connectivity
  - Must handle git repository state changes gracefully
  - Must validate Conventional Commits format before committing
  - No global state - all dependencies injected via constructors
  - Graceful error handling - no panics in library code

**Scale/Scope**:
  - Single-user CLI tool
  - Handles typical git repositories (up to thousands of changed files)
  - Supports multiple AI providers with extensible architecture

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Verify compliance with gitcomm Constitution principles:

- ✅ **Clean Architecture**: Structure follows layer separation:
  - `cmd/gitcomm/` - CLI entrypoint
  - `internal/service/` - Business logic (commit message formatting, validation)
  - `internal/repository/` - Git operations abstraction (Repository Pattern)
  - `internal/model/` - Domain models (CommitMessage, RepositoryState, etc.)
  - `internal/config/` - Configuration management
  - `internal/ai/` - AI provider interfaces and implementations
  - `pkg/` - Shared utilities (tokenization, validation)
  - `test/` - Test utilities and integration tests

- ✅ **Interface-Driven Development**: All dependencies defined as interfaces:
  - `GitRepository` interface for git operations
  - `AIProvider` interface for AI integrations
  - `TokenCalculator` interface for token estimation
  - `MessageValidator` interface for format validation
  - Dependency injection via constructors throughout

- ✅ **Test-First Development**: TDD mandatory for core business logic:
  - Tests written before implementation (Red-Green-Refactor)
  - Table-driven tests for validation and formatting functions
  - Integration tests for git operations and AI provider interactions
  - Mock interfaces for external dependencies

- ✅ **Idiomatic Go**: Design follows Go conventions:
  - PascalCase for exported, camelCase for unexported
  - Interface names with `-er` suffix (e.g., `GitRepository`, `AIProvider`)
  - Error variables with `Err` prefix
  - Tabs for indentation, double quotes for strings

- ✅ **Error Handling**: Explicit error handling strategy:
  - Custom error types for business logic errors (`ErrInvalidFormat`, `ErrNoChanges`, etc.)
  - Wrapped errors with context (`fmt.Errorf("context: %w", err)`)
  - No panics in library code
  - Clear error messages for users

- ✅ **Context & Thread Safety**: Context propagation planned:
  - `context.Context` for cancellation and timeouts (AI provider calls)
  - No shared mutable state requiring synchronization
  - Sequential workflow (no concurrent operations needed)

- ✅ **Technical Constraints**: All constraints met:
  - No global state - all dependencies injected
  - Graceful shutdown for long-running operations (AI calls)
  - Resource cleanup (close git repository handles, HTTP clients)

- ✅ **Operational Constraints**: Logging and secrets management:
  - Structured logging via `zerolog` with JSON formatting
  - Configuration file for AI provider credentials (not in CLI args)
  - No secrets exposed in logs or error messages

**Violations**: None - design fully complies with constitution principles.

## Project Structure

### Documentation (this feature)

```text
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
cmd/
└── gitcomm/
    └── main.go                    # CLI entrypoint

internal/
├── service/
│   ├── commit_service.go         # Core commit message creation logic
│   ├── validation_service.go     # Conventional Commits validation
│   └── formatting_service.go     # Message formatting
├── repository/
│   ├── git_repository.go         # Git operations interface
│   └── git_repository_impl.go    # go-git implementation
├── ai/
│   ├── provider.go               # AIProvider interface
│   ├── openai_provider.go        # OpenAI implementation
│   ├── anthropic_provider.go     # Anthropic implementation
│   ├── local_provider.go         # Local model implementation
│   └── token_calculator.go       # Token calculation interface and implementations
├── model/
│   ├── commit_message.go         # CommitMessage domain model
│   ├── repository_state.go       # RepositoryState domain model
│   └── config.go                  # Configuration models
├── config/
│   └── config.go                  # Configuration loading and management
├── ui/
│   ├── prompts.go                # Interactive prompt handlers
│   └── display.go                # Message display formatting
└── utils/
    └── validation.go              # Validation utilities

pkg/
├── conventional/
│   └── validator.go               # Conventional Commits validation logic
└── tokenization/
    ├── tiktoken.go                # OpenAI tokenization
    └── fallback.go                # Character-based fallback

test/
├── integration/
│   ├── git_operations_test.go   # Git integration tests
│   └── ai_provider_test.go       # AI provider integration tests
└── mocks/
    └── [generated mocks]          # Mock implementations for testing

configs/
└── config.yaml.example            # Example configuration file
```

**Structure Decision**: Single CLI project following Clean Architecture. The structure separates concerns into layers: `cmd/` for entrypoint, `internal/` for application-specific code (not importable by other projects), and `pkg/` for shared utilities that could be reused. The `internal/` directory is organized by feature/domain (service, repository, ai, model, config, ui) with clear interfaces between layers. This enables testability, maintainability, and future extensibility (e.g., adding new AI providers).

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

No violations - design fully complies with constitution principles.
