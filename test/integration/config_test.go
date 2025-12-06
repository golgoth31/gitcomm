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
