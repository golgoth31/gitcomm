<!--
Sync Impact Report:
Version change: N/A → 1.0.0 (initial constitution)
Modified principles: N/A (new constitution)
Added sections:
  - Core Principles (6 principles)
  - Technical Constraints
  - Regulatory/Operational Constraints
  - Development Workflow
  - Governance
Templates requiring updates:
  ✅ plan-template.md - Constitution Check section aligns with principles
  ✅ spec-template.md - Requirements section aligns with constraints
  ✅ tasks-template.md - Task organization aligns with TDD principle
  ✅ checklist-template.md - No constitution-specific references
Follow-up TODOs: None
-->

# gitcomm Constitution

## Core Principles

### I. Clean Architecture & Repository Pattern

All code MUST follow Clean Architecture principles with clear layer separation. Structure code into layers: `cmd/` for application entrypoints, `internal/` for private application code (service, repository, model, config, utils), `pkg/` for shared utilities, `configs/` for configuration schemas, `test/` for test utilities and integration tests, and `docs/` for documentation. The Repository Pattern MUST be used to separate data access logic from business logic. High-level modules MUST NOT depend on low-level modules (Dependency Inversion Principle).

**Rationale**: Clean Architecture ensures maintainability, testability, and independence from frameworks. The Repository Pattern decouples business logic from data access, enabling easier testing and future storage mechanism changes.

### II. Interface-Driven Development & Dependency Injection

All public functions MUST interact with interfaces, not concrete types. Prefer small, purpose-specific interfaces over large ones. Use explicit dependency injection via constructors; avoid global state. All dependencies MUST be injected through constructor functions.

**Rationale**: Interfaces enable mocking and testing, reduce coupling, and make code more flexible. Dependency injection eliminates hidden dependencies and makes testing straightforward.

### III. Test-First Development (NON-NEGOTIABLE)

TDD is mandatory for core business logic: Tests written → User approved → Tests fail → Then implement. Follow the Red-Green-Refactor cycle strictly. Use standard Go testing framework (`testing` package) for all tests. Organize tests alongside source files with `_test.go` suffix for unit tests, and integration tests under `test/integration/`. Include table-driven tests for functions with many input variants. Write unit tests using table-driven patterns and parallel execution (`t.Parallel()`). Ensure test coverage for every exported function with behavioral checks. Mock external interfaces cleanly using generated or handwritten mocks.

**Rationale**: Test-first development ensures code correctness, prevents regressions, and drives better design by forcing consideration of usage before implementation.

### IV. Idiomatic Go Code Style

Follow idiomatic Go conventions as defined in Effective Go and Google's Go Style Guide. Use named functions over long anonymous ones. Organize logic into small, composable functions with single responsibility. Use tabs for indentation, double quotes for strings, omit semicolons (unless required for disambiguation). Limit line length to 100-120 characters where practical. Use `gofmt` or `goimports` to enforce formatting. Enforce naming consistency with `golangci-lint`. Use `gofumpt` for stricter formatting if configured.

**Naming Conventions**: PascalCase for exported identifiers, camelCase for unexported, UPPERCASE for exported constants, snake_case for environment variables, kebab-case for directory/file names. Prefix boolean variables with verbs (`isReady`, `hasError`, `canConnect`). Use complete words over abbreviations (except common ones: `err`, `ctx`, `req`, `res`, `id`, `msg`). Interface names use noun or `-er` suffix. Error variables use `Err` prefix for exported errors.

**Rationale**: Consistent code style improves readability, maintainability, and team collaboration. Idiomatic Go ensures the codebase is familiar to Go developers and aligns with community standards.

### V. Explicit Error Handling & Resource Management

Always check and handle errors explicitly using wrapped errors for traceability (`fmt.Errorf("context: %w", err)`). Use custom error types for wrapping and handling business logic errors. Defer closing resources and handle them carefully to avoid leaks. No panics in library code; return errors instead. Must handle errors explicitly—no panics in library code.

**Rationale**: Explicit error handling makes failures visible and traceable. Proper resource management prevents leaks and ensures graceful degradation.

### VI. Context Propagation & Thread Safety

Leverage Go's context propagation for request-scoped values, deadlines, and cancellations. Must use `context.Context` for request-scoped values and cancellation. Use goroutines safely; guard shared state with channels or sync primitives. Must be thread-safe when using goroutines—use channels or sync primitives for concurrent access.

**Rationale**: Context propagation enables proper cancellation and timeout handling. Thread safety prevents race conditions and data corruption in concurrent code.

## Technical Constraints

- Must use `context.Context` for request-scoped values and cancellation
- No global state—all dependencies must be injected
- Must handle errors explicitly—no panics in library code
- Must support graceful shutdown and resource cleanup
  - Gracefully close client connections
  - Complete in-flight message processing
  - Persist session state if applicable
- Must be thread-safe when using goroutines
  - Use channels or sync primitives for concurrent access
- Must minimize allocations and avoid premature optimization
- Must validate input using Go structs and validation tags
- Must use dependency injection via constructors (avoid global state)
- Keep logic decoupled from framework-specific code
- Write short, focused functions with a single responsibility

## Regulatory/Operational Constraints

- Must not expose secrets in logs or error messages
- Must support namespace isolation for multi-tenant scenarios (if applicable)
- Must be observable for operations teams
- Logging via `zerolog` with JSON formatting (raw, json, or json with ECS fields)
- Don't hardcode config—use environment variables or config files
- Don't expose secrets—use `.env` or secret managers

## Development Workflow

- Use semantic versioning for releases
- Follow conventional commit messages
- Maintain a `CHANGELOG.md` for tracking changes
- Use feature branches for development
- Require code review before merging to main branch
- Separate fast unit tests from slower integration and E2E tests
- Use test helpers to reduce test code duplication
- Test files should be in the same package (for white-box testing) or `_test` package (for black-box testing)

## Patterns to Avoid

- Don't use global state unless absolutely required
- Don't hardcode config—use environment variables or config files
- Don't panic or exit in library code; return errors instead
- Don't expose secrets—use `.env` or secret managers
- Avoid embedding business logic in HTTP handlers
- Avoid unnecessary abstraction; keep things simple and readable

## Governance

This constitution supersedes all other practices and conventions. All PRs and reviews MUST verify compliance with these principles. Amendments require documentation, approval, and a migration plan if breaking changes are introduced. Complexity must be justified when deviating from these principles.

**Amendment Procedure**:
1. Propose amendment with rationale and impact analysis
2. Review by maintainers
3. Update constitution version according to semantic versioning:
   - MAJOR: Backward incompatible governance/principle removals or redefinitions
   - MINOR: New principle/section added or materially expanded guidance
   - PATCH: Clarifications, wording, typo fixes, non-semantic refinements
4. Update dependent templates and documentation
5. Communicate changes to all contributors

**Compliance Review**: All code contributions must pass constitution checks before merge. Automated checks via `golangci-lint` and test coverage requirements enforce technical principles. Manual review ensures architectural and design principle compliance.

**Version**: 1.0.0 | **Ratified**: 2025-01-27 | **Last Amended**: 2025-01-27
