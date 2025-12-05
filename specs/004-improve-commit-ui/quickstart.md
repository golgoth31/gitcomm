# Quick Start: Improved Commit Message UI

**Feature**: 004-improve-commit-ui
**Date**: 2025-01-27

## Overview

This feature enhances the commit message creation experience with an interactive select list for commit types and proper multiline input for body and footer fields. The UI is more intuitive, faster to use, and supports detailed commit messages.

## Key Features

- ✅ **Interactive commit type selection** with checkmarks and highlighting
- ✅ **Arrow key navigation** for quick type selection
- ✅ **Letter-based navigation** (type 'f' to jump to 'feat')
- ✅ **Multiline body input** with double-Enter completion
- ✅ **Multiline footer input** with double-Enter completion
- ✅ **Blank line preservation** within content
- ✅ **Whitespace handling** (whitespace-only = empty)

## Basic Usage

### Interactive Commit Type Selection

When you run `gitcomm`, the commit type selection now appears as an interactive list:

```
? Choose a type(<scope>):
  ✓ feat     [new feature]
    fix      [bug fix]
    docs     [changes to documentation]
    style    [format, missing semi colons, etc; no code change]
    refactor [refactor production code]
```

**Navigation**:
- Use **↑/↓ arrow keys** to move selection
- Type a **letter** (e.g., 'f' for 'feat') to jump to first matching option
- Press **Enter** to confirm selection
- Press **Escape** to cancel

### Multiline Body Input

The body field now supports multiple lines:

```
Body (optional, press Enter twice on empty line when done):
> This commit adds a new feature for user authentication.
>
> It includes:
> - Login functionality
> - Session management
> - Password reset
>
> [Press Enter twice on empty line to complete]
```

**Completion**:
- Press **Enter** to create a new line
- Press **Enter twice** on consecutive empty lines to complete
- Single blank lines are preserved as content
- Whitespace-only input is treated as empty

### Multiline Footer Input

The footer field also supports multiple lines:

```
Footer (optional, press Enter twice on empty line when done):
> Fixes #123
> Closes #456
>
> [Press Enter twice on empty line to complete]
```

**Completion**: Same as body (double Enter on empty lines)

## Examples

### Example 1: Quick Feature Commit

```bash
$ gitcomm
? Choose a type(<scope>):
  ✓ feat     [new feature]  # Pre-selected, just press Enter
    fix      [bug fix]
    ...

Scope (optional): auth
Subject (required): Add user login functionality

Body (optional):
> Implements user authentication with email/password.
>
> Features:
> - Login form
> - Session management
> - Remember me option
>
> [Double Enter to complete]

Footer (optional):
> Fixes #42
>
> [Double Enter to complete]
```

### Example 2: Using Letter Navigation

```bash
$ gitcomm
? Choose a type(<scope>):
    feat     [new feature]
  ✓ fix      [bug fix]  # Typed 'f' to jump here
    docs     [changes to documentation]
    ...

# Continue with commit...
```

### Example 3: Skipping Optional Fields

```bash
$ gitcomm
# ... select type, enter scope, enter subject ...

Body (optional):
> [Press Enter twice immediately to skip]

Footer (optional):
> [Press Enter twice immediately to skip]
```

## Tips & Tricks

### Quick Type Selection

- **Type the first letter** of the commit type to jump directly (e.g., 'd' for 'docs')
- **Use arrow keys** for precise navigation
- **First option (feat) is pre-selected** - just press Enter for features

### Multiline Input

- **Single Enter** = new line (content continues)
- **Double Enter** = complete input (two empty lines)
- **Blank lines are preserved** - use them for paragraph separation
- **Whitespace-only input** = treated as empty (skipped)

### Keyboard Shortcuts

- **↑/↓** - Navigate commit type list
- **Letter key** - Jump to matching commit type
- **Enter** - Confirm selection or add new line
- **Double Enter** - Complete multiline input
- **Escape** - Cancel and restore staging state

## Troubleshooting

### Selection Not Responding

**Symptom**: Arrow keys don't move selection

**Solution**:
- Ensure terminal supports ANSI escape sequences
- Check that terminal is in raw mode (bubbletea handles this)
- Try typing a letter to jump to an option

### Multiline Input Not Completing

**Symptom**: Pressing Enter twice doesn't complete input

**Solution**:
- Ensure you press Enter on an **empty line** (no spaces or text)
- Press Enter **twice consecutively** on empty lines
- Check that you're not in the middle of typing

### Blank Lines Not Preserved

**Symptom**: Blank lines disappear in final commit message

**Solution**:
- Ensure you're using the multiline input (not single-line)
- Single blank lines should be preserved
- Only two consecutive empty lines signal completion

## Integration with Existing Features

This feature integrates seamlessly with existing gitcomm features:

- **Auto-Staging**: Works with auto-staging (002-auto-stage-restore)
- **State Restoration**: Escape key cancellation restores staging state
- **AI Generation**: AI-generated messages still work (uses same UI)
- **Debug Logging**: UI interactions can be logged with --debug flag (003-debug-logging)

## Performance

- Commit type selection: <50ms response time
- Multiline input: <10ms per keystroke
- No noticeable delay in UI updates
- Smooth interaction even on slower terminals

## Accessibility

- **Keyboard-only navigation** - No mouse required
- **Clear visual indicators** - Checkmarks and highlighting
- **Letter-based shortcuts** - Fast navigation for power users
- **Graceful degradation** - Works in terminals without color support

## Next Steps

After using the improved UI:

1. **Familiarize yourself** with arrow key and letter navigation
2. **Practice multiline input** - try double Enter completion
3. **Use blank lines** for better commit message formatting
4. **Take advantage of letter navigation** for faster type selection

For more information, see the [full specification](./spec.md).
