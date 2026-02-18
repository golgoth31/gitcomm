# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**gitcomm** is a Go CLI tool that automates git commit message creation following the Conventional Commits specification. It supports interactive manual input and AI-assisted generation via OpenAI, Anthropic, Mistral, and local model providers.

## Build & Development Commands

```bash
make build            # go build -o gitcomm ./cmd/gitcomm
make test             # go test -v ./...
make test-coverage    # Tests with HTML coverage report
make lint             # golangci-lint run
make format           # go fmt + goimports
make install          # go install ./cmd/gitcomm
make clean            # Remove binary and coverage files
```

Run a single test:
```bash
go test -v -run TestFunctionName ./internal/service/
```

Requires: Go 1.25.0+, git CLI, golangci-lint.

## Architecture

### Layer Structure

```
cmd/gitcomm/main.go          → Entry point, calls cmd.Execute()
internal/cmd/root.go          → Cobra root command, signal handling, orchestration
internal/service/              → Business logic orchestration (commit workflow)
internal/ai/                   → AIProvider interface + implementations per provider
internal/repository/           → GitRepository interface, git CLI wrapper
internal/model/                → Data structures (CommitMessage, RepositoryState, etc.)
internal/config/               → YAML config loading with ${ENV_VAR} substitution
internal/ui/                   → Interactive prompts (Charmbracelet Huh)
internal/utils/                → Zerolog logger, custom error types
pkg/ai/prompt/                 → Unified prompt generation for all AI providers
pkg/conventional/              → Conventional Commits validation rules
pkg/git/config/                → Git config reader (user.name, signing keys)
pkg/tokenization/              → Token estimation (TikToken, Anthropic, fallback)
```

### Key Interfaces

- **`AIProvider`** (`internal/ai/provider.go`): `GenerateCommitMessage(ctx, repoState) (string, error)` — implemented by each AI provider
- **`GitRepository`** (`internal/repository/git_repository.go`): Git operations (stage, commit, diff, state capture/restore) — uses external git CLI, not go-git

### Core Workflow

`CommitService.CreateCommit()` in `internal/service/commit_service.go` orchestrates:
1. Capture staging state → auto-stage files → get repo state (diffs)
2. If AI enabled: generate message → validate → prompt accept/edit/reject (max 3 retries)
3. If manual: interactive prompts for type/scope/subject/body/footer
4. Validate against Conventional Commits → commit → restore staging on cancel (3s timeout)

### Configuration

Location: `~/.gitcomm/config.yaml` (see `configs/config.yaml.example`)
- Supports `${ENV_VAR}` placeholder substitution with validation
- Created with 0600 permissions on first run

## Coding Conventions (from project constitution)

- **TDD is mandatory**: Red-Green-Refactor cycle for core business logic
- **Interface-driven**: All public functions use interfaces, not concrete types. Dependency injection via constructors
- **No global state**: All dependencies injected
- **Error handling**: Always wrap errors with context (`fmt.Errorf("context: %w", err)`), custom error types for business logic, no panics in library code
- **Context propagation**: `context.Context` for cancellation/deadlines throughout
- **Naming**: PascalCase exported, camelCase unexported, `Err` prefix for error vars, `-er` suffix for interfaces, verb prefix for booleans (`isReady`, `hasError`)
- **Common abbreviations allowed**: `err`, `ctx`, `req`, `res`, `id`, `msg`
- **Tests**: `_test.go` alongside source, table-driven, parallel execution (`t.Parallel()`), mock external interfaces
- **Conventional Commits**: Valid types are `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`, `version`

## Linting

golangci-lint with: errcheck, gosimple, govet, ineffassign, staticcheck, unused, gofmt, goimports, misspell, revive (exported + var-naming rules).
