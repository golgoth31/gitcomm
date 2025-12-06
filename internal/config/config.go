package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/golgoth31/gitcomm/internal/model"
	"github.com/golgoth31/gitcomm/internal/utils"
	"github.com/spf13/viper"
)

// T015: Placeholder regex pattern compiled once for reuse
var placeholderRegex = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}`)

// Config represents the application configuration
type Config struct {
	AI AIConfig
}

// AIConfig represents AI provider configuration
type AIConfig struct {
	DefaultProvider string
	Providers       map[string]model.AIProviderConfig
}

// LoadConfig loads configuration from file or environment variables
func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()

	// Set default config path
	if configPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		configPath = filepath.Join(homeDir, ".gitcomm", "config.yaml")
	}

	// T013: Validate path is not a directory
	if info, err := os.Stat(configPath); err == nil && info.Mode().IsDir() {
		return nil, fmt.Errorf("config path is a directory, not a file: %s", configPath)
	}

	// T012: Check if config file exists, create if missing
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// T014: Create parent directories
		configDir := filepath.Dir(configPath)
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create config directory %s: %w", configDir, err)
		}

		// T015: Create empty file with O_CREATE|O_WRONLY|O_EXCL flags
		// T016: Set file permissions to 0600
		file, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0600)
		if err != nil {
			// T017: Handle race condition (file created by another process)
			if os.IsExist(err) {
				// File was created by another process, treat as success
			} else {
				return nil, fmt.Errorf("failed to create config file: %w", err)
			}
		} else {
			file.Close()
			// T018: Log file creation
			utils.Logger.Debug().Str("path", configPath).Msg("Created config file")
		}
	}

	// Configure viper
	v.SetConfigType("yaml")
	v.SetEnvPrefix("GITCOMM")
	v.AutomaticEnv()

	// T029-T032: Read config file content and perform placeholder substitution before YAML parsing
	content, err := os.ReadFile(configPath)
	if err != nil {
		// If file doesn't exist, continue with defaults (backward compatibility)
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// File doesn't exist - viper will use defaults
	} else {
		// T029-T030: Perform placeholder substitution on config content
		substituted, err := substitutePlaceholdersInContent(string(content))
		if err != nil {
			return nil, err
		}

		// T031-T032: Read from substituted content instead of file
		if err := v.ReadConfig(strings.NewReader(substituted)); err != nil {
			// Config file is optional, continue with defaults
		}
	}

	config := &Config{
		AI: AIConfig{
			DefaultProvider: v.GetString("ai.default_provider"),
			Providers:       make(map[string]model.AIProviderConfig),
		},
	}

	// Load provider configurations
	providers := v.GetStringMap("ai.providers")
	for name := range providers {
		providerConfig := model.AIProviderConfig{
			Name:     name,
			APIKey:   v.GetString(fmt.Sprintf("ai.providers.%s.api_key", name)),
			Model:    v.GetString(fmt.Sprintf("ai.providers.%s.model", name)),
			Endpoint: v.GetString(fmt.Sprintf("ai.providers.%s.endpoint", name)),
			Timeout:  30 * time.Second,
		}

		// Override timeout if specified
		if timeoutStr := v.GetString(fmt.Sprintf("ai.providers.%s.timeout", name)); timeoutStr != "" {
			if timeout, err := time.ParseDuration(timeoutStr); err == nil {
				providerConfig.Timeout = timeout
			}
		}

		config.AI.Providers[name] = providerConfig
	}

	return config, nil
}

// GetProviderConfig returns the configuration for a specific provider
func (c *Config) GetProviderConfig(name string) (*model.AIProviderConfig, error) {
	if name == "" {
		name = c.AI.DefaultProvider
	}

	provider, ok := c.AI.Providers[name]
	if !ok {
		return nil, fmt.Errorf("provider %s not configured", name)
	}

	return &provider, nil
}

// T016: findPlaceholders identifies all placeholders in config content using regex
func findPlaceholders(content string) []string {
	matches := placeholderRegex.FindAllStringSubmatch(content, -1)
	var vars []string
	for _, match := range matches {
		if len(match) > 1 {
			vars = append(vars, match[1]) // match[1] is the captured group (variable name)
		}
	}
	return vars
}

// T017: validatePlaceholderSyntax validates placeholder syntax (check for nested, multiline, invalid chars)
func validatePlaceholderSyntax(content string) error {
	// Check for nested placeholders
	if strings.Contains(content, "${${") {
		return fmt.Errorf("nested placeholders not allowed")
	}

	// Check for multiline placeholders (newline within ${...})
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.Contains(line, "${") {
			// Check if placeholder spans multiple lines
			if strings.Contains(line, "${") && !strings.Contains(line, "}") {
				// Check if next line continues the placeholder
				if i+1 < len(lines) && strings.Contains(lines[i+1], "}") {
					return fmt.Errorf("multiline placeholders not allowed")
				}
			}
		}
	}

	// Find all placeholders and validate each one
	matches := placeholderRegex.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 0 {
			placeholder := match[0]
			varName := match[1]

			// Validate variable name contains only valid characters
			for _, r := range varName {
				if !((r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_') {
					return fmt.Errorf("invalid placeholder syntax: %s (invalid character in variable name)", placeholder)
				}
			}

			// Check for spaces in placeholder
			if strings.Contains(placeholder, " ") {
				return fmt.Errorf("invalid placeholder syntax: %s (spaces not allowed)", placeholder)
			}
		}
	}

	return nil
}

// T018: extractEnvVarNames extracts environment variable names from placeholders
func extractEnvVarNames(content string) []string {
	return findPlaceholders(content)
}

// T019: validateEnvVarsExist validates all environment variables exist using os.LookupEnv()
func validateEnvVarsExist(varNames []string) error {
	var missing []string
	for _, v := range varNames {
		if _, exists := os.LookupEnv(v); !exists {
			missing = append(missing, v)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing environment variables: %s", strings.Join(missing, ", "))
	}
	return nil
}

// T020: substitutePlaceholders substitutes placeholders with environment variable values
func substitutePlaceholders(content string, varNames []string) string {
	result := content
	for _, v := range varNames {
		placeholder := fmt.Sprintf("${%s}", v)
		value := os.Getenv(v) // Already validated to exist
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}

// T021: processContentWithComments processes content, skipping comment lines
func processContentWithComments(content string, processFunc func(string) (string, error)) (string, error) {
	lines := strings.Split(content, "\n")
	var processedLines []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			// Skip comments - preserve unchanged
			processedLines = append(processedLines, line)
		} else {
			// Process non-comment lines
			processed, err := processFunc(line)
			if err != nil {
				return "", err
			}
			processedLines = append(processedLines, processed)
		}
	}

	return strings.Join(processedLines, "\n"), nil
}

// substitutePlaceholdersInContent performs full placeholder substitution with validation
func substitutePlaceholdersInContent(content string) (string, error) {
	// T017: Validate placeholder syntax
	if err := validatePlaceholderSyntax(content); err != nil {
		return "", fmt.Errorf("invalid placeholder syntax: %w", err)
	}

	// T016: Find all placeholders
	varNames := findPlaceholders(content)
	if len(varNames) == 0 {
		// No placeholders, return content unchanged
		return content, nil
	}

	// Deduplicate variable names
	uniqueVars := make(map[string]bool)
	var deduplicated []string
	for _, v := range varNames {
		if !uniqueVars[v] {
			uniqueVars[v] = true
			deduplicated = append(deduplicated, v)
		}
	}

	// T019: Validate all environment variables exist
	if err := validateEnvVarsExist(deduplicated); err != nil {
		return "", err
	}

	// T020: Substitute placeholders (skip comments)
	result, err := processContentWithComments(content, func(line string) (string, error) {
		return substitutePlaceholders(line, deduplicated), nil
	})
	if err != nil {
		return "", err
	}

	return result, nil
}
