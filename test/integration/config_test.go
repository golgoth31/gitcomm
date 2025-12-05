package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/golgoth31/gitcomm/internal/config"
	"github.com/golgoth31/gitcomm/internal/utils"
)

func init() {
	// Initialize logger for tests
	utils.InitLogger(true)
}

// T019: Test file creation when missing
func TestLoadConfig_CreatesFileWhenMissing(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Ensure file doesn't exist
	os.Remove(configPath)

	// Load config - should create file
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if cfg == nil {
		t.Fatalf("LoadConfig returned nil config")
	}

	// Verify file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatalf("Config file was not created at %s", configPath)
	}

	// Verify file is empty (0 bytes)
	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("Failed to stat config file: %v", err)
	}
	if info.Size() != 0 {
		t.Fatalf("Expected config file to be empty (0 bytes), got %d bytes", info.Size())
	}
}

// T020: Test existing file not modified
func TestLoadConfig_DoesNotModifyExistingFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Create existing file with content
	existingContent := "ai:\n  default_provider: test\n"
	if err := os.WriteFile(configPath, []byte(existingContent), 0600); err != nil {
		t.Fatalf("Failed to create existing config file: %v", err)
	}

	// Load config
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if cfg == nil {
		t.Fatalf("LoadConfig returned nil config")
	}

	// Verify file content was not modified
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}
	if string(content) != existingContent {
		t.Fatalf("Config file content was modified. Expected:\n%s\nGot:\n%s", existingContent, string(content))
	}
}

// T021: Test parent directory creation
func TestLoadConfig_CreatesParentDirectories(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "subdir", "nested", "config.yaml")

	// Ensure parent directories don't exist
	os.RemoveAll(filepath.Dir(configPath))

	// Load config - should create parent directories
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if cfg == nil {
		t.Fatalf("LoadConfig returned nil config")
	}

	// Verify parent directories were created
	parentDir := filepath.Dir(configPath)
	if _, err := os.Stat(parentDir); os.IsNotExist(err) {
		t.Fatalf("Parent directory was not created at %s", parentDir)
	}

	// Verify file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatalf("Config file was not created at %s", configPath)
	}
}

// T022: Test custom config path
func TestLoadConfig_CustomConfigPath(t *testing.T) {
	tmpDir := t.TempDir()
	customPath := filepath.Join(tmpDir, "custom", "my-config.yaml")

	// Ensure file doesn't exist
	os.RemoveAll(filepath.Dir(customPath))

	// Load config with custom path
	cfg, err := config.LoadConfig(customPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if cfg == nil {
		t.Fatalf("LoadConfig returned nil config")
	}

	// Verify file was created at custom path
	if _, err := os.Stat(customPath); os.IsNotExist(err) {
		t.Fatalf("Config file was not created at custom path %s", customPath)
	}
}

// T023: Test default config path (~/.gitcomm/config.yaml)
func TestLoadConfig_DefaultConfigPath(t *testing.T) {
	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home directory: %v", err)
	}

	defaultPath := filepath.Join(homeDir, ".gitcomm", "config.yaml")

	// Backup existing file if it exists
	backupPath := defaultPath + ".backup"
	if _, err := os.Stat(defaultPath); err == nil {
		if err := os.Rename(defaultPath, backupPath); err != nil {
			t.Fatalf("Failed to backup existing config file: %v", err)
		}
		defer os.Rename(backupPath, defaultPath)
	}

	// Ensure file doesn't exist
	os.Remove(defaultPath)
	os.RemoveAll(filepath.Dir(defaultPath))

	// Load config with empty path (should use default)
	cfg, err := config.LoadConfig("")
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if cfg == nil {
		t.Fatalf("LoadConfig returned nil config")
	}

	// Verify file was created at default path
	if _, err := os.Stat(defaultPath); os.IsNotExist(err) {
		t.Fatalf("Config file was not created at default path %s", defaultPath)
	}

	// Cleanup
	os.Remove(defaultPath)
	os.RemoveAll(filepath.Dir(defaultPath))
}

// T024: Test single placeholder substitution
func TestLoadConfig_SinglePlaceholderSubstitution(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	testVar := "TEST_OPENAI_API_KEY"
	testValue := "sk-test12345"
	os.Setenv(testVar, testValue)
	defer os.Unsetenv(testVar)

	// Create config file with placeholder
	configContent := "ai:\n  providers:\n    openai:\n      api_key: ${" + testVar + "}\n"
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Load config
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if cfg == nil {
		t.Fatalf("LoadConfig returned nil config")
	}

	// Verify placeholder was substituted by checking the loaded config
	providerConfig, err := cfg.GetProviderConfig("openai")
	if err != nil {
		t.Fatalf("Failed to get provider config: %v", err)
	}
	if providerConfig.APIKey != testValue {
		t.Fatalf("Expected API key %s, got %s", testValue, providerConfig.APIKey)
	}
}

// T025: Test multiple placeholder substitution
func TestLoadConfig_MultiplePlaceholderSubstitution(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	var1 := "TEST_OPENAI_API_KEY"
	val1 := "sk-openai-123"
	var2 := "TEST_ANTHROPIC_API_KEY"
	val2 := "sk-anthropic-456"

	os.Setenv(var1, val1)
	os.Setenv(var2, val2)
	defer os.Unsetenv(var1)
	defer os.Unsetenv(var2)

	// Create config file with multiple placeholders
	configContent := "ai:\n  providers:\n    openai:\n      api_key: ${" + var1 + "}\n    anthropic:\n      api_key: ${" + var2 + "}\n"
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Load config
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if cfg == nil {
		t.Fatalf("LoadConfig returned nil config")
	}

	// Verify both placeholders were substituted
	openaiConfig, err := cfg.GetProviderConfig("openai")
	if err != nil {
		t.Fatalf("Failed to get openai config: %v", err)
	}
	if openaiConfig.APIKey != val1 {
		t.Fatalf("Expected OpenAI API key %s, got %s", val1, openaiConfig.APIKey)
	}

	anthropicConfig, err := cfg.GetProviderConfig("anthropic")
	if err != nil {
		t.Fatalf("Failed to get anthropic config: %v", err)
	}
	if anthropicConfig.APIKey != val2 {
		t.Fatalf("Expected Anthropic API key %s, got %s", val2, anthropicConfig.APIKey)
	}
}

// T026: Test placeholder in nested YAML structure
func TestLoadConfig_PlaceholderInNestedYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	testVar := "TEST_NESTED_API_KEY"
	testValue := "sk-nested-789"
	os.Setenv(testVar, testValue)
	defer os.Unsetenv(testVar)

	// Create config file with placeholder in nested structure
	configContent := "ai:\n  default_provider: openai\n  providers:\n    openai:\n      api_key: ${" + testVar + "}\n      model: gpt-4\n"
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Load config
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if cfg == nil {
		t.Fatalf("LoadConfig returned nil config")
	}

	// Verify placeholder was substituted in nested structure
	providerConfig, err := cfg.GetProviderConfig("openai")
	if err != nil {
		t.Fatalf("Failed to get provider config: %v", err)
	}
	if providerConfig.APIKey != testValue {
		t.Fatalf("Expected API key %s in nested structure, got %s", testValue, providerConfig.APIKey)
	}
}

// T027: Test backward compatibility (config without placeholders)
func TestLoadConfig_BackwardCompatibility(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Create config file without placeholders (existing format)
	configContent := "ai:\n  default_provider: openai\n  providers:\n    openai:\n      api_key: sk-hardcoded-123\n      model: gpt-4\n"
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Load config
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if cfg == nil {
		t.Fatalf("LoadConfig returned nil config")
	}

	// Verify config loads correctly without placeholders
	providerConfig, err := cfg.GetProviderConfig("openai")
	if err != nil {
		t.Fatalf("Failed to get provider config: %v", err)
	}
	if providerConfig.APIKey != "sk-hardcoded-123" {
		t.Fatalf("Expected API key sk-hardcoded-123, got %s", providerConfig.APIKey)
	}
}

// T028: Test empty string value substitution
func TestLoadConfig_EmptyStringValueSubstitution(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	testVar := "TEST_EMPTY_VAR"
	os.Setenv(testVar, "")
	defer os.Unsetenv(testVar)

	// Create config file with placeholder for empty string variable
	configContent := "ai:\n  providers:\n    openai:\n      api_key: ${" + testVar + "}\n"
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Load config - should succeed (empty string is valid)
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v (empty string should be treated as valid)", err)
	}
	if cfg == nil {
		t.Fatalf("LoadConfig returned nil config")
	}

	// Verify empty string was substituted
	providerConfig, err := cfg.GetProviderConfig("openai")
	if err != nil {
		t.Fatalf("Failed to get provider config: %v", err)
	}
	if providerConfig.APIKey != "" {
		t.Fatalf("Expected empty API key, got %s", providerConfig.APIKey)
	}
}

// T034: Test missing single environment variable error
func TestLoadConfig_MissingSingleEnvironmentVariable(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	missingVar := "MISSING_VAR_12345"
	// Ensure variable is not set
	os.Unsetenv(missingVar)

	// Create config file with placeholder for missing variable
	configContent := "ai:\n  providers:\n    openai:\n      api_key: ${" + missingVar + "}\n"
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Load config - should fail with error
	_, err := config.LoadConfig(configPath)
	if err == nil {
		t.Fatalf("Expected error when environment variable is missing")
	}

	// Verify error message identifies the missing variable
	errorMsg := err.Error()
	if !contains(errorMsg, missingVar) {
		t.Fatalf("Error message should identify missing variable %s, got: %s", missingVar, errorMsg)
	}
}

// T035: Test missing multiple environment variables error
func TestLoadConfig_MissingMultipleEnvironmentVariables(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	missingVar1 := "MISSING_VAR_1"
	missingVar2 := "MISSING_VAR_2"
	// Ensure variables are not set
	os.Unsetenv(missingVar1)
	os.Unsetenv(missingVar2)

	// Create config file with placeholders for missing variables
	configContent := "ai:\n  providers:\n    openai:\n      api_key: ${" + missingVar1 + "}\n    anthropic:\n      api_key: ${" + missingVar2 + "}\n"
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Load config - should fail with error
	_, err := config.LoadConfig(configPath)
	if err == nil {
		t.Fatalf("Expected error when environment variables are missing")
	}

	// Verify error message lists all missing variables
	errorMsg := err.Error()
	if !contains(errorMsg, missingVar1) {
		t.Fatalf("Error message should identify missing variable %s, got: %s", missingVar1, errorMsg)
	}
	if !contains(errorMsg, missingVar2) {
		t.Fatalf("Error message should identify missing variable %s, got: %s", missingVar2, errorMsg)
	}
}

// T036: Test error message clarity (lists all missing variables)
func TestLoadConfig_ErrorMessageClarity(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	missingVar1 := "MISSING_VAR_A"
	missingVar2 := "MISSING_VAR_B"
	missingVar3 := "MISSING_VAR_C"

	os.Unsetenv(missingVar1)
	os.Unsetenv(missingVar2)
	os.Unsetenv(missingVar3)

	// Create config file with multiple missing variables
	configContent := "ai:\n  providers:\n    openai:\n      api_key: ${" + missingVar1 + "}\n    anthropic:\n      api_key: ${" + missingVar2 + "}\n    mistral:\n      api_key: ${" + missingVar3 + "}\n"
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Load config - should fail
	_, err := config.LoadConfig(configPath)
	if err == nil {
		t.Fatalf("Expected error when environment variables are missing")
	}

	// Verify error message is clear and lists all missing variables
	errorMsg := err.Error()
	if len(errorMsg) < 20 {
		t.Fatalf("Error message should be descriptive, got: %s", errorMsg)
	}

	// Check that all missing variables are mentioned
	if !contains(errorMsg, missingVar1) || !contains(errorMsg, missingVar2) || !contains(errorMsg, missingVar3) {
		t.Fatalf("Error message should list all missing variables. Got: %s", errorMsg)
	}

	// Check that error message mentions "missing environment variables"
	if !contains(errorMsg, "missing") {
		t.Fatalf("Error message should mention 'missing', got: %s", errorMsg)
	}
}

// T037: Test empty string value treated as valid (does not exit)
func TestLoadConfig_EmptyStringValueTreatedAsValid(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	testVar := "EMPTY_VAR_TEST"
	os.Setenv(testVar, "")
	defer os.Unsetenv(testVar)

	// Create config file with placeholder for empty string variable
	configContent := "ai:\n  providers:\n    openai:\n      api_key: ${" + testVar + "}\n"
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Load config - should succeed (empty string is valid, not missing)
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig should not fail for empty string value, got error: %v", err)
	}
	if cfg == nil {
		t.Fatalf("LoadConfig returned nil config")
	}

	// Verify config loaded successfully
	providerConfig, err := cfg.GetProviderConfig("openai")
	if err != nil {
		t.Fatalf("Failed to get provider config: %v", err)
	}
	if providerConfig.APIKey != "" {
		t.Fatalf("Expected empty API key, got %s", providerConfig.APIKey)
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
