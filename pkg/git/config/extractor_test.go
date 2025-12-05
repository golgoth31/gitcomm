package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/golgoth31/gitcomm/internal/utils"
)

func TestFileConfigExtractor_Extract_LocalConfigOnly(t *testing.T) {
	// Setup: Initialize logger for debug messages
	utils.InitLogger(true)

	// Create temporary directory with .git/config
	tmpDir := t.TempDir()
	gitDir := filepath.Join(tmpDir, ".git")
	os.MkdirAll(gitDir, 0755)

	configContent := `[user]
	name = Test User
	email = test@example.com
`
	configPath := filepath.Join(gitDir, "config")
	os.WriteFile(configPath, []byte(configContent), 0644)

	extractor := NewFileConfigExtractor()
	config := extractor.Extract(tmpDir)

	if config.UserName != "Test User" {
		t.Errorf("Expected UserName 'Test User', got '%s'", config.UserName)
	}
	if config.UserEmail != "test@example.com" {
		t.Errorf("Expected UserEmail 'test@example.com', got '%s'", config.UserEmail)
	}
}

func TestFileConfigExtractor_Extract_GlobalConfigOnly(t *testing.T) {
	// Setup: Initialize logger for debug messages
	utils.InitLogger(true)

	// Create temporary directory without .git/config
	tmpDir := t.TempDir()

	// Create temporary global config
	homeDir := t.TempDir()
	globalConfigPath := filepath.Join(homeDir, ".gitconfig")
	configContent := `[user]
	name = Global User
	email = global@example.com
`
	os.WriteFile(globalConfigPath, []byte(configContent), 0644)

	// Mock home directory by setting HOME env var
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", homeDir)
	defer os.Setenv("HOME", originalHome)

	extractor := NewFileConfigExtractor()
	config := extractor.Extract(tmpDir)

	if config.UserName != "Global User" {
		t.Errorf("Expected UserName 'Global User', got '%s'", config.UserName)
	}
	if config.UserEmail != "global@example.com" {
		t.Errorf("Expected UserEmail 'global@example.com', got '%s'", config.UserEmail)
	}
}

func TestFileConfigExtractor_Extract_LocalTakesPrecedence(t *testing.T) {
	// Setup: Initialize logger for debug messages
	utils.InitLogger(true)

	// Create temporary directory with .git/config
	tmpDir := t.TempDir()
	gitDir := filepath.Join(tmpDir, ".git")
	os.MkdirAll(gitDir, 0755)

	localConfigContent := `[user]
	name = Local User
	email = local@example.com
`
	configPath := filepath.Join(gitDir, "config")
	os.WriteFile(configPath, []byte(localConfigContent), 0644)

	// Create temporary global config with different values
	homeDir := t.TempDir()
	globalConfigPath := filepath.Join(homeDir, ".gitconfig")
	globalConfigContent := `[user]
	name = Global User
	email = global@example.com
`
	os.WriteFile(globalConfigPath, []byte(globalConfigContent), 0644)

	// Mock home directory
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", homeDir)
	defer os.Setenv("HOME", originalHome)

	extractor := NewFileConfigExtractor()
	config := extractor.Extract(tmpDir)

	// Local config should take precedence
	if config.UserName != "Local User" {
		t.Errorf("Expected UserName 'Local User' (local precedence), got '%s'", config.UserName)
	}
	if config.UserEmail != "local@example.com" {
		t.Errorf("Expected UserEmail 'local@example.com' (local precedence), got '%s'", config.UserEmail)
	}
}

func TestFileConfigExtractor_Extract_MissingFiles(t *testing.T) {
	// Setup: Initialize logger for debug messages
	utils.InitLogger(true)

	// Create temporary directory without .git/config
	tmpDir := t.TempDir()

	// Set HOME to non-existent directory to simulate missing global config
	homeDir := filepath.Join(t.TempDir(), "nonexistent")
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", homeDir)
	defer os.Setenv("HOME", originalHome)

	extractor := NewFileConfigExtractor()
	config := extractor.Extract(tmpDir)

	// Should use defaults
	if config.UserName != "gitcomm" {
		t.Errorf("Expected default UserName 'gitcomm', got '%s'", config.UserName)
	}
	if config.UserEmail != "gitcomm@local" {
		t.Errorf("Expected default UserEmail 'gitcomm@local', got '%s'", config.UserEmail)
	}
}

func TestFileConfigExtractor_Extract_UnreadableFiles(t *testing.T) {
	// Setup: Initialize logger for debug messages
	utils.InitLogger(true)

	// Create temporary directory
	tmpDir := t.TempDir()
	gitDir := filepath.Join(tmpDir, ".git")
	os.MkdirAll(gitDir, 0755)

	// Create unreadable config file (on Unix systems)
	configPath := filepath.Join(gitDir, "config")
	os.WriteFile(configPath, []byte("invalid"), 0000) // No read permissions
	defer os.Chmod(configPath, 0644)                  // Restore for cleanup

	// Create readable global config as fallback
	homeDir := t.TempDir()
	globalConfigPath := filepath.Join(homeDir, ".gitconfig")
	configContent := `[user]
	name = Global User
	email = global@example.com
`
	os.WriteFile(globalConfigPath, []byte(configContent), 0644)

	// Mock home directory
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", homeDir)
	defer os.Setenv("HOME", originalHome)

	extractor := NewFileConfigExtractor()
	config := extractor.Extract(tmpDir)

	// Should fall back to global config
	if config.UserName != "Global User" {
		t.Errorf("Expected UserName 'Global User' (fallback), got '%s'", config.UserName)
	}
}

func TestFileConfigExtractor_Extract_PartialValues(t *testing.T) {
	// Setup: Initialize logger for debug messages
	utils.InitLogger(true)

	// Create temporary directory with .git/config
	tmpDir := t.TempDir()
	gitDir := filepath.Join(tmpDir, ".git")
	os.MkdirAll(gitDir, 0755)

	// Config with only user.name, missing user.email
	configContent := `[user]
	name = Test User
`
	configPath := filepath.Join(gitDir, "config")
	os.WriteFile(configPath, []byte(configContent), 0644)

	// Save and clear HOME to prevent reading real global config
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// Set HOME to non-existent directory
	os.Setenv("HOME", filepath.Join(t.TempDir(), "nonexistent"))

	extractor := NewFileConfigExtractor()
	config := extractor.Extract(tmpDir)

	if config.UserName != "Test User" {
		t.Errorf("Expected UserName 'Test User', got '%s'", config.UserName)
	}
	// Missing user.email should use default
	if config.UserEmail != "gitcomm@local" {
		t.Errorf("Expected default UserEmail 'gitcomm@local', got '%s'", config.UserEmail)
	}
}

func TestFileConfigExtractor_Extract_Performance(t *testing.T) {
	// Setup: Initialize logger for debug messages
	utils.InitLogger(true)

	// Create temporary directory with .git/config
	tmpDir := t.TempDir()
	gitDir := filepath.Join(tmpDir, ".git")
	os.MkdirAll(gitDir, 0755)

	configContent := `[user]
	name = Test User
	email = test@example.com
`
	configPath := filepath.Join(gitDir, "config")
	os.WriteFile(configPath, []byte(configContent), 0644)

	extractor := NewFileConfigExtractor()

	start := time.Now()
	config := extractor.Extract(tmpDir)
	duration := time.Since(start)

	if duration > 50*time.Millisecond {
		t.Errorf("Extract took %v, expected <50ms", duration)
	}

	// Verify it still works
	if config.UserName != "Test User" {
		t.Errorf("Expected UserName 'Test User', got '%s'", config.UserName)
	}
}

func TestFileConfigExtractor_Extract_SSHSigningConfiguration(t *testing.T) {
	// Setup: Initialize logger for debug messages
	utils.InitLogger(true)

	// Create temporary directory with .git/config
	tmpDir := t.TempDir()
	gitDir := filepath.Join(tmpDir, ".git")
	os.MkdirAll(gitDir, 0755)

	configContent := `[user]
	name = Test User
	email = test@example.com
	signingkey = /path/to/key.pub
[gpg]
	format = ssh
[commit]
	gpgsign = true
`
	configPath := filepath.Join(gitDir, "config")
	os.WriteFile(configPath, []byte(configContent), 0644)

	extractor := NewFileConfigExtractor()
	config := extractor.Extract(tmpDir)

	if config.SigningKey != "/path/to/key.pub" {
		t.Errorf("Expected SigningKey '/path/to/key.pub', got '%s'", config.SigningKey)
	}
	if config.GPGFormat != "ssh" {
		t.Errorf("Expected GPGFormat 'ssh', got '%s'", config.GPGFormat)
	}
	if !config.CommitGPGSign {
		t.Errorf("Expected CommitGPGSign true, got %v", config.CommitGPGSign)
	}
}

func TestFileConfigExtractor_Extract_CommitGPGSignFalse(t *testing.T) {
	// Setup: Initialize logger for debug messages
	utils.InitLogger(true)

	// Create temporary directory with .git/config
	tmpDir := t.TempDir()
	gitDir := filepath.Join(tmpDir, ".git")
	os.MkdirAll(gitDir, 0755)

	configContent := `[user]
	name = Test User
	email = test@example.com
	signingkey = /path/to/key.pub
[gpg]
	format = ssh
[commit]
	gpgsign = false
`
	configPath := filepath.Join(gitDir, "config")
	os.WriteFile(configPath, []byte(configContent), 0644)

	extractor := NewFileConfigExtractor()
	config := extractor.Extract(tmpDir)

	if config.CommitGPGSign {
		t.Errorf("Expected CommitGPGSign false, got %v", config.CommitGPGSign)
	}
}
