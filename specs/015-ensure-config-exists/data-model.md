# Data Model: Ensure Config File Exists Before Reading

**Feature**: 015-ensure-config-exists
**Date**: 2025-01-27

## Overview

This feature does not introduce new data entities. It modifies the behavior of the existing `LoadConfig` function to ensure the config file exists before reading. The data model remains unchanged from the existing implementation.

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

## File System Entities

### Config File

**Type**: File system entity (not a Go struct)

**Description**: YAML configuration file stored on disk.

**Location**:
- Default: `~/.gitcomm/config.yaml`
- Custom: User-specified path via `configPath` parameter

**Properties**:
- **Format**: YAML
- **Content**: Empty (0 bytes) when newly created
- **Permissions**: 0600 (owner read/write only)
- **Parent Directory Permissions**: 0755 (default file system permissions)

**State Transitions**:
- Does not exist → Created (empty, 0 bytes)
- Exists → Read (no modification)

**Validation Rules**:
- Must be a file (not a directory) - FR-005
- Must be readable by owner
- Must be writable by owner
- Must be valid YAML (when non-empty)

### Config Directory

**Type**: File system entity (not a Go struct)

**Description**: Parent directory containing the config file.

**Location**:
- Default: `~/.gitcomm/`
- Custom: Parent directory of user-specified `configPath`

**Properties**:
- **Permissions**: 0755 (default file system permissions)
- **Created**: Recursively if missing (FR-004)

**State Transitions**:
- Does not exist → Created (with all parent directories)
- Exists → Used (no modification)

## No New Data Structures

This feature does not introduce new Go data structures. It only modifies the file system state (creating files/directories) before loading the existing `Config` structure.

## Data Flow

1. **Input**: `configPath string` (optional, defaults to `~/.gitcomm/config.yaml`)
2. **File System Check**: Check if config file exists
3. **File Creation** (if missing):
   - Create parent directories (0755)
   - Create empty file (0 bytes, 0600)
4. **Config Loading**: Read config file using viper (existing behavior)
5. **Output**: `*Config, error` (existing return type)

## Edge Cases

- **File exists as directory**: Return error (FR-005)
- **Parent directory missing**: Create recursively (FR-004)
- **Permission denied**: Return error with context (FR-002, FR-010)
- **Race condition**: Handle gracefully, treat as success if file exists (FR-007)
- **File deleted between creation and read**: Handle as missing file, continue with defaults (FR-008)
