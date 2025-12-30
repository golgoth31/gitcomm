# Data Model: Unify Prompt Functions to Use Default Variants

**Date**: 2025-01-27
**Feature**: 016-unify-prompt-defaults

## Overview

This refactoring does not introduce new data models or modify existing ones. The data structures remain unchanged:

- `model.CommitMessage`: Unchanged
- `ui.PrefilledCommitMessage`: Unchanged

## Existing Entities

### CommitMessage

**Location**: `internal/model/commit_message.go`

Represents a structured commit message following Conventional Commits format.

**Fields**:
- `Type` (string): Commit type (feat, fix, docs, etc.)
- `Scope` (string): Optional scope
- `Subject` (string): Required commit subject
- `Body` (string): Optional commit body
- `Footer` (string): Optional commit footer
- `Signoff` (bool): Whether to include signoff

**Validation Rules**: (unchanged)
- Type must be one of: feat, fix, docs, style, refactor, test, chore, version
- Subject is required and cannot be empty
- Scope, Body, and Footer are optional

### PrefilledCommitMessage

**Location**: `internal/ui/prompts.go`

Represents pre-populated commit message data used as defaults in prompts.

**Fields**:
- `Type` (string): Pre-filled commit type
- `Scope` (string): Pre-filled scope (may be empty)
- `Subject` (string): Pre-filled subject
- `Body` (string): Pre-filled body (may be empty)
- `Footer` (string): Pre-filled footer (may be empty)

**Usage**: Passed to `promptCommitMessage` to pre-populate prompt fields. After refactoring, empty strings are passed when fields are not pre-filled.

## Changes Summary

**No data model changes**. This refactoring only affects:
- Function call patterns in `commit_service.go`
- Removal of unused prompt functions from `prompts.go`

The data structures and their usage remain identical.
