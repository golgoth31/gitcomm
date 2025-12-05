# Research: Environment Variable Placeholder Substitution in Config Files

**Feature**: 014-env-var-placeholders
**Date**: 2025-01-27

## Technology Decisions

### 1. Placeholder Identification Pattern

**Decision**: Use regex pattern `\$\{([A-Za-z_][A-Za-z0-9_]*)\}` to identify and extract placeholders.

**Rationale**:
- Regex provides efficient pattern matching for `${VAR_NAME}` syntax
- Capturing group extracts the variable name for validation and lookup
- Pattern enforces valid environment variable naming (starts with letter/underscore, followed by alphanumeric/underscore)
- Standard Go `regexp` package is well-tested and performant
- Pattern can be compiled once and reused for multiple matches

**Alternatives Considered**:
- String search with manual parsing: More error-prone, harder to validate syntax
- Third-party template libraries (e.g., `text/template`): Overkill for simple substitution, adds dependency
- Simple string replacement: Doesn't validate syntax or handle edge cases

**Implementation Pattern**:
```go
var placeholderRegex = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}`)

func findPlaceholders(content string) []string {
    matches := placeholderRegex.FindAllStringSubmatch(content, -1)
    var vars []string
    for _, match := range matches {
        vars = append(vars, match[1]) // match[1] is the captured group
    }
    return vars
}
```

### 2. Substitution Strategy

**Decision**: Process config file content as text before YAML parsing, performing all substitutions in memory, then pass substituted content to viper.

**Rationale**:
- Ensures placeholders are replaced before YAML parsing (viper reads clean YAML)
- Maintains backward compatibility (files without placeholders work unchanged)
- In-memory processing avoids file system modifications
- Single-pass substitution is efficient
- Allows validation of all placeholders before any substitution

**Alternatives Considered**:
- Post-processing after viper reads: Complex, may break YAML structure parsing
- Write temporary file: Unnecessary I/O, cleanup complexity, security concerns
- Custom viper extension: Overly complex, viper doesn't support this pattern

**Implementation Pattern**:
```go
func substitutePlaceholders(content string) (string, error) {
    // 1. Find all placeholders
    // 2. Validate syntax
    // 3. Check all environment variables exist
    // 4. Perform substitutions
    // 5. Return substituted content
}
```

### 3. Environment Variable Lookup

**Decision**: Use `os.LookupEnv()` to check for variable existence, then `os.Getenv()` for value retrieval.

**Rationale**:
- `os.LookupEnv()` distinguishes between unset variables and empty string values (required by FR-007)
- Standard library, no dependencies
- Thread-safe (read-only operations)
- Efficient for single lookups

**Alternatives Considered**:
- `os.Getenv()` only: Cannot distinguish unset from empty (violates FR-007)
- `os.Environ()` + manual parsing: Less efficient, unnecessary complexity
- Third-party env libraries: Unnecessary dependency for simple use case

**Implementation Pattern**:
```go
func getEnvVar(name string) (string, bool) {
    value, exists := os.LookupEnv(name)
    return value, exists
}
```

### 4. Error Handling for Missing Variables

**Decision**: Collect all missing variables during validation phase, then exit with a single error message listing all missing variables.

**Rationale**:
- Better user experience (user sees all missing variables at once, not one at a time)
- Efficient (single validation pass)
- Aligns with FR-005 and FR-010 (exit immediately, clear error messages)
- Prevents partial substitution (FR-006: process all before proceeding)

**Alternatives Considered**:
- Fail on first missing variable: Less helpful, requires multiple attempts to fix
- Continue with partial substitution: Violates spec (must exit immediately)
- Warn and continue: Violates spec (must exit)

**Implementation Pattern**:
```go
func validatePlaceholders(vars []string) error {
    var missing []string
    for _, v := range vars {
        if _, exists := os.LookupEnv(v); !exists {
            missing = append(missing, v)
        }
    }
    if len(missing) > 0 {
        return fmt.Errorf("missing environment variables: %s", strings.Join(missing, ", "))
    }
    return nil
}
```

### 5. Placeholder Substitution Order

**Decision**: Process placeholders in order of appearance, allowing same variable to appear multiple times (all occurrences replaced with same value).

**Rationale**:
- Simple and predictable behavior
- Efficient (can use `strings.ReplaceAll()` for each variable)
- Handles multiple occurrences correctly (FR-006)
- No ambiguity about substitution order

**Alternatives Considered**:
- Deduplicate variables first: Unnecessary optimization, adds complexity
- Process in sorted order: Unpredictable, doesn't match user expectations
- One-pass replacement: More complex regex, harder to validate

**Implementation Pattern**:
```go
func substituteAll(content string, vars []string) string {
    result := content
    for _, v := range vars {
        value := os.Getenv(v) // Already validated to exist
        placeholder := fmt.Sprintf("${%s}", v)
        result = strings.ReplaceAll(result, placeholder, value)
    }
    return result
}
```

### 6. YAML Comment Handling

**Decision**: Process config file content as raw text, but skip substitution in comment lines (lines starting with `#`).

**Rationale**:
- YAML comments are documentation, not config data
- Aligns with FR-012 (ignore placeholders in comments)
- Simple to implement (check if line starts with `#` before processing)
- Preserves user intent (comments remain unchanged)

**Alternatives Considered**:
- Process all content: Would substitute in comments (violates FR-012)
- Full YAML parsing to identify comments: Overly complex, viper doesn't expose comment locations
- Regex to skip comment lines: Simple and effective

**Implementation Pattern**:
```go
func processLine(line string) string {
    trimmed := strings.TrimSpace(line)
    if strings.HasPrefix(trimmed, "#") {
        return line // Skip comments
    }
    return substitutePlaceholdersInLine(line)
}
```

### 7. Invalid Syntax Detection

**Decision**: Validate placeholder syntax using regex, then perform additional checks for nested placeholders and newlines.

**Rationale**:
- Regex catches most invalid patterns (spaces, invalid characters)
- Additional checks handle edge cases (nested `${${VAR}}`, multiline)
- Fail-fast approach (FR-011: exit immediately on invalid syntax)
- Clear error messages identify the invalid pattern

**Alternatives Considered**:
- Only regex validation: May miss edge cases (nested, multiline)
- Full AST parsing: Overkill for simple syntax validation
- Allow invalid syntax: Violates spec (must exit on invalid syntax)

**Implementation Pattern**:
```go
func validatePlaceholderSyntax(placeholder string) error {
    // Check regex match
    if !placeholderRegex.MatchString(placeholder) {
        return fmt.Errorf("invalid placeholder syntax: %s", placeholder)
    }
    // Check for nested placeholders
    if strings.Contains(placeholder, "${${") {
        return fmt.Errorf("nested placeholders not allowed: %s", placeholder)
    }
    // Check for newlines
    if strings.Contains(placeholder, "\n") {
        return fmt.Errorf("multiline placeholders not allowed: %s", placeholder)
    }
    return nil
}
```

### 8. Integration with Viper

**Decision**: Read config file content, perform substitution, write substituted content to temporary in-memory buffer, then use viper to read from buffer.

**Rationale**:
- Viper supports reading from `io.Reader`, allowing in-memory processing
- No file system modifications (maintains constraint)
- Clean separation: substitution logic independent of YAML parsing
- Backward compatible (viper behavior unchanged)

**Alternatives Considered**:
- Write to temporary file: Unnecessary I/O, cleanup complexity
- Modify viper source: Not feasible, viper is external dependency
- Custom YAML parser: Overly complex, viper already handles YAML well

**Implementation Pattern**:
```go
func LoadConfig(configPath string) (*Config, error) {
    // Read file content
    content, err := os.ReadFile(configPath)
    if err != nil {
        return nil, err
    }

    // Substitute placeholders
    substituted, err := substitutePlaceholders(string(content))
    if err != nil {
        return nil, err
    }

    // Create viper instance and read from substituted content
    v := viper.New()
    v.SetConfigType("yaml")
    if err := v.ReadConfig(strings.NewReader(substituted)); err != nil {
        return nil, err
    }

    // Continue with existing config loading logic...
}
```

## Best Practices Summary

1. **Validate before substitute**: Check all placeholders exist before performing any substitutions (fail-fast)
2. **Use regex for pattern matching**: Efficient and reliable for `${VAR}` syntax
3. **Distinguish unset from empty**: Use `os.LookupEnv()` to handle empty string values correctly
4. **Process as text before YAML**: Ensures clean YAML structure for viper
5. **Skip comments**: Preserve YAML comments unchanged
6. **Collect all errors**: Report all missing variables in single error message
7. **Use standard library**: `os`, `regexp`, `strings` are sufficient, no external dependencies needed
8. **In-memory processing**: Avoid file system modifications, use `io.Reader` for viper
