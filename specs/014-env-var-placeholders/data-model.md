# Data Model: Environment Variable Placeholder Substitution in Config Files

**Feature**: 014-env-var-placeholders
**Date**: 2025-01-27

## Overview

This feature does not introduce new data entities. It modifies the behavior of the existing `LoadConfig` function to process environment variable placeholders in config file content before YAML parsing. The data model remains unchanged from the existing implementation.

## Existing Entities

### Config

**Location**: `internal/config/config.go`

**Description**: Represents the application configuration loaded from file or environment variables.

**Fields**:
- `AI AIConfig` - AI provider configuration

**Relationships**: None

**State Transitions**: None (immutable after creation)

**Validation Rules**: None (handled by viper)

### AIConfig

**Location**: `internal/config/config.go`

**Description**: Represents AI provider configuration.

**Fields**:
- `DefaultProvider string` - Default AI provider name
- `Providers map[string]model.AIProviderConfig` - Map of provider configurations

**Relationships**: Contains `model.AIProviderConfig` instances

**State Transitions**: None (immutable after creation)

**Validation Rules**: None (handled by viper)

## Processing Entities (Not Data Structures)

### Placeholder

**Type**: Processing entity (not a Go struct)

**Description**: Text pattern in config file content that represents a reference to an environment variable.

**Pattern**: `${ENV_VAR_NAME}` where `ENV_VAR_NAME` matches:
- Starts with letter or underscore: `[A-Za-z_]`
- Followed by alphanumeric characters or underscores: `[A-Za-z0-9_]*`
- No spaces, no nested placeholders, no newlines

**Validation Rules**:
- Must match regex pattern: `\$\{([A-Za-z_][A-Za-z0-9_]*)\}`
- Must not contain nested placeholders (e.g., `${${VAR}}`)
- Must not span multiple lines
- Variable name must correspond to an existing environment variable (or application exits)

**State Transitions**:
- Identified in config content → Extracted variable name
- Variable name extracted → Environment variable looked up
- Variable found → Placeholder replaced with value
- Variable not found → Application exits with error

### Config File Content

**Type**: Processing entity (string, not a Go struct)

**Description**: Raw YAML content from config file before placeholder substitution.

**Properties**:
- **Format**: YAML text
- **Content**: May contain `${ENV_VAR_NAME}` placeholders
- **Comments**: YAML comments (lines starting with `#`) are preserved unchanged

**State Transitions**:
- Raw content → Placeholders identified
- Placeholders identified → Variables validated
- Variables validated → Substitutions performed
- Substituted content → YAML parsed by viper

## No New Data Structures

This feature does not introduce new Go data structures. It only modifies the processing of config file content (string manipulation) before loading the existing `Config` structure.

## Data Flow

1. **Input**: `configPath string` (optional, defaults to `~/.gitcomm/config.yaml`)
2. **File Reading**: Read config file content as raw text
3. **Placeholder Identification**: Find all `${ENV_VAR_NAME}` patterns in content
4. **Syntax Validation**: Validate placeholder syntax (alphanumeric/underscores only, no nested, no multiline)
5. **Variable Lookup**: Check all referenced environment variables exist
6. **Error Handling**: If validation fails, exit immediately with clear error message
7. **Substitution**: Replace all placeholders with environment variable values
8. **YAML Parsing**: Pass substituted content to viper for YAML parsing
9. **Config Loading**: Continue with existing config loading logic
10. **Output**: `*Config, error` (existing return type)

## Edge Cases

- **Invalid placeholder syntax**: Exit immediately with error identifying invalid pattern
- **Missing environment variable**: Exit immediately with error listing all missing variables
- **Empty string value**: Treated as valid substitution (empty string replaces placeholder)
- **Multiple occurrences**: All occurrences of same placeholder replaced with same value
- **Placeholders in comments**: Ignored (comments preserved unchanged)
- **Placeholders in string values**: Processed (all YAML string values are processed)
- **Nested placeholders**: Rejected (exit with error - invalid syntax)
- **Multiline placeholders**: Rejected (exit with error - invalid syntax)
