# Feature Specification: Environment Variable Placeholder Substitution in Config Files

**Feature Branch**: `014-env-var-placeholders`
**Created**: 2025-01-27
**Status**: Draft
**Input**: User description: "config file can contain placeholders in the form ${...}. those place holders contains environment variables names. update config loading to identify those placeholders and read the corresponding environment variables. if the environment variable is not available, exit immediatly."

## Clarifications

### Session 2025-01-27

- Q: When the config file contains invalid placeholder syntax (spaces in variable name, nested placeholders, malformed syntax), what should the system do? → A: Exit with error immediately (fail fast on malformed syntax)
- Q: Should placeholders be processed when they appear in YAML comments or in string values that are meant to be literal text? → A: Process placeholders only in YAML values (ignore comments, process all string values)
- Q: How should the system handle placeholders that span multiple lines? → A: Reject multiline placeholders (exit with error - invalid syntax)

## User Scenarios & Testing

### User Story 1 - Basic Environment Variable Substitution (Priority: P1)

Users can reference environment variables in their config file using `${ENV_VAR_NAME}` syntax. When the config file is loaded, these placeholders are automatically replaced with the actual values from the environment. This allows users to keep sensitive information like API keys out of version control while still using a structured config file.

**Why this priority**: This is the core functionality - without this, the feature provides no value. It enables secure configuration management by separating secrets from config files.

**Independent Test**: Create a config file with `${API_KEY}` placeholder, set the environment variable, load the config, and verify the placeholder is replaced with the environment variable value.

**Acceptance Scenarios**:

1. **Given** a config file contains `${OPENAI_API_KEY}` placeholder, **When** the environment variable `OPENAI_API_KEY` is set to `sk-12345`, **Then** the config loader replaces `${OPENAI_API_KEY}` with `sk-12345` in the loaded configuration
2. **Given** a config file contains multiple placeholders like `${OPENAI_API_KEY}` and `${ANTHROPIC_API_KEY}`, **When** both environment variables are set, **Then** all placeholders are replaced with their corresponding environment variable values
3. **Given** a config file contains `${API_KEY}` in a nested YAML structure (e.g., `providers.openai.api_key: ${OPENAI_API_KEY}`), **When** the config is loaded, **Then** the placeholder is replaced with the environment variable value in the correct location

---

### User Story 2 - Missing Environment Variable Error Handling (Priority: P1)

When a config file contains a placeholder for an environment variable that is not set, the application must exit immediately with a clear error message indicating which environment variable is missing. This prevents the application from running with incomplete or incorrect configuration.

**Why this priority**: Security and correctness - running with missing required configuration could lead to runtime errors, security issues, or incorrect behavior. Immediate failure is safer than silent failures.

**Independent Test**: Create a config file with `${MISSING_VAR}` placeholder, ensure the environment variable is not set, attempt to load the config, and verify the application exits with an error message identifying the missing variable.

**Acceptance Scenarios**:

1. **Given** a config file contains `${MISSING_API_KEY}` placeholder, **When** the environment variable `MISSING_API_KEY` is not set, **Then** the application exits immediately with an error message that clearly identifies `MISSING_API_KEY` as the missing variable
2. **Given** a config file contains multiple placeholders where some are set and some are missing, **When** the config is loaded, **Then** the application exits with an error message listing all missing environment variables
3. **Given** a config file contains `${EMPTY_VAR}` placeholder, **When** the environment variable `EMPTY_VAR` is set to an empty string, **Then** the application treats this as a valid value (empty string substitution) and does not exit

---

### Edge Cases

- **Invalid placeholder syntax** (e.g., `${VAR NAME}` with spaces, `${${NESTED}}` nested placeholders): System exits immediately with an error message identifying the invalid syntax
- **Placeholders in comments**: Placeholders in YAML comments are ignored (comments are not processed for substitution)
- **Placeholders in string values**: All placeholders in YAML string values are processed (no escape mechanism for literal placeholders - if a string contains `${VAR}`, it will be substituted)
- What happens when environment variable names are case-sensitive vs case-insensitive?
- **Multiline placeholders**: Placeholders spanning multiple lines are rejected (exit with error - invalid syntax, as environment variable names cannot contain newlines)
- What happens when a placeholder appears multiple times in the config file?
- What happens when a config file contains both `${VAR}` and `$VAR` syntax (should only `${VAR}` be processed)?

## Requirements

### Functional Requirements

- **FR-001**: System MUST identify placeholders in the form `${ENV_VAR_NAME}` within config file content
- **FR-002**: System MUST extract environment variable names from placeholder syntax (text between `${` and `}`)
- **FR-011**: System MUST validate placeholder syntax (alphanumeric characters and underscores only, no spaces, no nested placeholders, no newlines) and exit immediately with an error if invalid syntax is detected
- **FR-003**: System MUST look up each extracted environment variable name in the current environment
- **FR-004**: System MUST replace each placeholder with the corresponding environment variable value when the variable exists
- **FR-005**: System MUST exit immediately with a clear error message when any required environment variable is not found
- **FR-006**: System MUST process all placeholders in the config file before proceeding with config loading
- **FR-007**: System MUST handle empty string values from environment variables as valid substitutions
- **FR-008**: System MUST preserve non-placeholder content in the config file exactly as written
- **FR-009**: System MUST process placeholders in all YAML values (strings, regardless of nesting level)
- **FR-010**: System MUST provide error messages that clearly identify which environment variable(s) are missing
- **FR-012**: System MUST ignore placeholders in YAML comments (comments are not processed for substitution)
- **FR-013**: System MUST process placeholders in all YAML string values (no escape mechanism for literal placeholders)

### Key Entities

- **Config File**: YAML file containing configuration with optional `${ENV_VAR_NAME}` placeholders
- **Environment Variable**: System environment variable that provides the value to substitute for a placeholder
- **Placeholder**: Text pattern in the form `${ENV_VAR_NAME}` that represents a reference to an environment variable

## Success Criteria

### Measurable Outcomes

- **SC-001**: Users can successfully load config files with environment variable placeholders when all required variables are set (100% success rate for valid configurations)
- **SC-002**: Application exits within 1 second when a required environment variable is missing (immediate failure requirement)
- **SC-003**: Error messages clearly identify missing environment variables (100% of error cases include the variable name)
- **SC-004**: All placeholders in a config file are processed and replaced correctly (100% placeholder substitution rate when variables exist)
- **SC-005**: Config files without placeholders continue to work exactly as before (zero regression for existing functionality)

## Constraints

- Must maintain backward compatibility with existing config files that do not use placeholders
- Must work with the existing YAML config file structure and format
- Must not modify the config file on disk (substitution happens in memory during loading)
- Must handle all valid YAML structures (nested objects, arrays, strings)
- Error handling must be immediate and clear (no silent failures or partial substitutions)

## Dependencies

- Existing config loading mechanism (`LoadConfig` function)
- Access to system environment variables
- YAML parsing capability (already exists via viper)

## Assumptions

- Environment variable names follow standard naming conventions (alphanumeric and underscores, typically uppercase)
- Placeholder syntax `${VAR}` is the only format to be processed (not `$VAR` or other variations)
- Environment variables are case-sensitive (matching the system's environment variable behavior)
- Empty string values from environment variables are intentional and valid (not treated as missing)
- Placeholders only appear in YAML values, not in keys or structural elements
- Placeholders in YAML comments are ignored (comments are not processed)
- All placeholders in YAML string values are processed (no escape mechanism for literal placeholders)
