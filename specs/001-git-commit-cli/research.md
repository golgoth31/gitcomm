# Research: Git Commit Message Automation CLI

**Date**: 2025-01-27
**Feature**: 001-git-commit-cli

## Research Decisions

### 1. AI Provider Integration Architecture

**Decision**: Use interface-based abstraction with provider-specific implementations

**Rationale**:
- Enables support for multiple AI providers (OpenAI, Anthropic, local models) without tight coupling
- Follows Interface-Driven Development principle from constitution
- Allows easy addition of new providers without modifying core logic
- Enables mocking for testing

**Alternatives Considered**:
- Single provider implementation: Too limiting, doesn't meet requirement for multiple providers
- Plugin system: Over-engineered for initial implementation, adds complexity
- Direct API calls in service layer: Violates Clean Architecture, hard to test

**Implementation Approach**:
- Define `AIProvider` interface with `GenerateCommitMessage(ctx context.Context, repoState RepositoryState) (string, error)`
- Implement provider-specific structs: `OpenAIProvider`, `AnthropicProvider`, `LocalProvider`
- Use dependency injection to select provider at runtime
- Store provider configuration (API keys, endpoints) in config file

**References**:
- Go interface best practices: https://go.dev/doc/effective_go#interfaces
- Multiple provider pattern: Common in Go ecosystem (e.g., database drivers)

---

### 2. Token Calculation Strategy

**Decision**: Provider-specific tokenization libraries with character-based fallback

**Rationale**:
- Ensures accurate token estimates for cost calculation and user decision-making
- Different providers use different tokenization methods (OpenAI uses BPE, Anthropic uses different encoding)
- Character-based fallback provides reasonable estimate for unknown/local providers
- Meets success criteria SC-006 (10% accuracy for 90% of cases)

**Alternatives Considered**:
- Simple character count: Inaccurate for OpenAI/Anthropic (BPE tokenization varies significantly)
- Word count: Less accurate than character count, doesn't account for tokenization nuances
- No calculation: Poor UX, users can't make informed decisions about AI usage

**Implementation Approach**:
- Use tiktoken Go bindings (or reimplement) for OpenAI tokenization
- Research Anthropic tokenization method and implement custom calculator
- Fallback: `tokens ≈ characters / 4` (rough average for most tokenizers)
- Cache tokenization results for same repository state to avoid recalculation

**References**:
- OpenAI tiktoken: https://github.com/openai/tiktoken
- Token counting best practices: Provider-specific documentation

---

### 3. Conventional Commits Validation

**Decision**: Implement validation library in `pkg/conventional/validator.go`

**Rationale**:
- Must ensure 100% compliance (SC-003) before committing
- Validation logic is reusable and testable in isolation
- Can be used for both AI-generated and manually-entered messages
- Follows Clean Architecture (validation in pkg/, business logic in internal/)

**Alternatives Considered**:
- External library: Risk of dependency issues, may not match exact spec requirements
- Regex-based validation: Fragile, hard to maintain, doesn't handle edge cases well
- No validation: Violates requirement FR-017, risks non-compliant commits

**Implementation Approach**:
- Parse commit message into structured components (type, scope, subject, body, footer)
- Validate each component against Conventional Commits spec:
  - Type: Must be one of: feat, fix, docs, style, refactor, test, chore, version
  - Scope: Optional, must be valid identifier
  - Subject: Required, imperative mood, no period, ≤72 chars
  - Body: Optional, wrapped at 72 chars, ≤320 chars
  - Footer: Optional, key-value pairs or issue references
- Return structured validation errors with specific failure reasons

**References**:
- Conventional Commits spec: https://www.conventionalcommits.org/en/v1.0.0/#specification
- Go parsing libraries: `strings`, `regexp`, or custom parser

---

### 4. TUI Framework Selection

**Decision**: Use `github.com/charmbracelet/bubbletea` for interactive prompts

**Rationale**:
- Mature, well-maintained TUI framework for Go
- Supports complex interactive workflows (multi-step prompts, validation)
- Good documentation and community support
- Integrates well with `lipgloss` for styling
- Follows Go idioms and patterns

**Alternatives Considered**:
- `github.com/manifoldco/promptui`: Simpler but less flexible for complex workflows
- `github.com/AlecAivazis/survey`: Good for simple forms, but limited customization
- Plain `fmt.Scan`: Too low-level, requires manual input handling and validation

**Implementation Approach**:
- Use bubbletea models for each prompt step (scope, subject, body, footer)
- Implement validation in model Update methods
- Chain models together for multi-step workflow
- Use lipgloss for consistent styling and formatting

**References**:
- Bubbletea: https://github.com/charmbracelet/bubbletea
- Lipgloss: https://github.com/charmbracelet/lipgloss
- TUI best practices: Charm ecosystem documentation

---

### 5. Git Operations Library

**Decision**: Use `github.com/go-git/go-git/v5` for git operations

**Rationale**:
- Pure Go implementation, no external git binary dependency
- Cross-platform compatibility
- Good API for porcelain commands (status, diff, commit)
- Active maintenance and community support
- Enables testing without actual git repository

**Alternatives Considered**:
- Executing `git` binary directly: Platform-dependent, harder to test, requires git installation
- `github.com/libgit2/git2go`: CGO dependency, more complex, platform-specific builds
- `github.com/src-d/go-git`: Older, less maintained fork

**Implementation Approach**:
- Create `GitRepository` interface abstracting go-git operations
- Implement interface using go-git for:
  - Repository state detection (`git status --porcelain`)
  - Staged/unstaged changes retrieval
  - Commit creation
  - Auto-staging (`git add -A` equivalent)
- Wrap go-git errors in domain-specific error types

**References**:
- go-git: https://github.com/go-git/go-git
- Git porcelain commands: https://git-scm.com/docs/git-status#_porcelain_format

---

### 6. Configuration Management

**Decision**: Use `github.com/spf13/viper` with YAML configuration file

**Rationale**:
- Viper supports multiple formats (YAML, JSON, TOML) with YAML as default
- Integrates well with Cobra for CLI flags
- Supports environment variable overrides
- Configuration file at `~/.gitcomm/config.yaml` for user settings
- Secure credential storage (not in CLI args)

**Alternatives Considered**:
- Environment variables only: Less user-friendly, harder to manage multiple providers
- JSON only: Less readable than YAML for configuration
- No config file: Violates requirement for provider selection and credential management

**Implementation Approach**:
- Default config location: `~/.gitcomm/config.yaml`
- Support config file override via CLI flag
- Structure:
  ```yaml
  ai:
    default_provider: openai
    providers:
      openai:
        api_key: ${OPENAI_API_KEY}
        model: gpt-4
      anthropic:
        api_key: ${ANTHROPIC_API_KEY}
        model: claude-3-opus
  ```
- Use viper for loading and environment variable substitution

**References**:
- Viper: https://github.com/spf13/viper
- Configuration best practices: 12-factor app principles

---

### 7. Error Handling Strategy

**Decision**: Custom error types with error wrapping

**Rationale**:
- Meets constitution requirement for explicit error handling
- Enables type-safe error checking in business logic
- Provides context through error wrapping
- Clear error messages for users

**Implementation Approach**:
- Define domain error types:
  - `ErrNotGitRepository` - CLI run outside git repo
  - `ErrNoChanges` - No staged/unstaged changes
  - `ErrInvalidFormat` - Commit message doesn't conform to spec
  - `ErrAIProviderUnavailable` - AI provider call failed
  - `ErrEmptySubject` - Subject cannot be empty
- Wrap errors with context: `fmt.Errorf("failed to generate commit message: %w", err)`
- Return errors, never panic (except in main.go for unrecoverable errors)

**References**:
- Go error handling: https://go.dev/blog/error-handling-and-go
- Error wrapping: https://go.dev/blog/go1.13-errors

---

### 8. Testing Strategy

**Decision**: Standard Go testing with Ginkgo/Gomega for BDD-style tests

**Rationale**:
- Meets constitution requirement for TDD
- Ginkgo/Gomega provides readable test syntax for complex scenarios
- Table-driven tests for validation and formatting functions
- Integration tests for git operations and AI provider interactions
- Mock interfaces for external dependencies

**Implementation Approach**:
- Unit tests: Table-driven tests for pure functions (validation, formatting)
- Integration tests: Test git operations with temporary repositories
- AI provider tests: Mock AIProvider interface, test error handling
- Test coverage: Aim for >80% coverage on core business logic
- Use `t.Parallel()` for independent tests

**References**:
- Ginkgo: https://onsi.github.io/ginkgo/
- Gomega: https://onsi.github.io/gomega/
- Go testing best practices: https://go.dev/doc/effective_go#testing

---

## Open Questions Resolved

1. **Q**: How to handle AI provider rate limiting?
   **A**: Implement timeout (30s default) and immediate fallback to manual input per FR-020

2. **Q**: Should token calculation be cached?
   **A**: Yes, cache results for same repository state to avoid redundant calculations

3. **Q**: How to handle very long commit messages?
   **A**: Validate against Conventional Commits limits (subject ≤72, body ≤320) and warn user

4. **Q**: Should we support commit message templates?
   **A**: Out of scope for initial implementation - focus on Conventional Commits format

---

## Technology Stack Summary

- **Language**: Go 1.25.0+
- **CLI Framework**: Cobra
- **TUI Framework**: Bubbletea + Lipgloss
- **Git Operations**: go-git
- **Configuration**: Viper
- **Logging**: Zerolog
- **Testing**: Ginkgo + Gomega
- **AI Integration**: Custom interfaces with provider implementations
- **Tokenization**: tiktoken (OpenAI), custom (Anthropic), fallback (others)
