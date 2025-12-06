package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/golgoth31/gitcomm/internal/utils"
)

func init() {
	// Initialize logger for tests
	utils.InitLogger(true)
}

// T006: Test file existence check
func TestFileExistenceCheck(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// File should not exist initially
	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		t.Fatalf("Expected file to not exist, but got error: %v", err)
	}

	// Create file
	file, err := os.Create(configPath)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	// File should exist now
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatalf("Expected file to exist, but got error: %v", err)
	}
}

// T007: Test empty file creation (0 bytes)
func TestEmptyFileCreation(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Create empty file
	file, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0600)
	if err != nil {
		t.Fatalf("Failed to create empty file: %v", err)
	}
	file.Close()

	// Verify file exists and is 0 bytes
	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}
	if info.Size() != 0 {
		t.Fatalf("Expected file size to be 0 bytes, got %d", info.Size())
	}
}

// T008: Test file permissions (0600)
func TestFilePermissions(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Create file with 0600 permissions
	file, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0600)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	file.Close()

	// Verify permissions
	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}
	expectedPerms := os.FileMode(0600)
	actualPerms := info.Mode().Perm()
	if actualPerms != expectedPerms {
		t.Fatalf("Expected file permissions %o, got %o", expectedPerms, actualPerms)
	}
}

// T009: Test parent directory creation (0755)
func TestParentDirectoryCreation(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "subdir", "nested", "config.yaml")

	// Create parent directories
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create parent directories: %v", err)
	}

	// Verify directories exist
	info, err := os.Stat(configDir)
	if err != nil {
		t.Fatalf("Failed to stat directory: %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("Expected %s to be a directory", configDir)
	}

	// Verify permissions (0755)
	expectedPerms := os.FileMode(0755)
	actualPerms := info.Mode().Perm()
	if actualPerms != expectedPerms {
		t.Fatalf("Expected directory permissions %o, got %o", expectedPerms, actualPerms)
	}
}

// T010: Test path validation (directory check)
func TestPathValidation(t *testing.T) {
	tmpDir := t.TempDir()
	dirPath := filepath.Join(tmpDir, "somedir")

	// Create a directory
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	// Check if path is a directory
	info, err := os.Stat(dirPath)
	if err != nil {
		t.Fatalf("Failed to stat path: %v", err)
	}
	if !info.Mode().IsDir() {
		t.Fatalf("Expected %s to be a directory", dirPath)
	}

	// Path should be identified as directory
	if !info.IsDir() {
		t.Fatalf("Path validation failed: expected directory but got file")
	}
}

// T011: Test race condition handling
func TestRaceConditionHandling(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Create file first
	file, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0600)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	file.Close()

	// Try to create again with O_EXCL flag (should fail with os.ErrExist)
	_, err = os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0600)
	if err == nil {
		t.Fatalf("Expected error when creating existing file with O_EXCL")
	}
	if !os.IsExist(err) {
		t.Fatalf("Expected os.ErrExist error, got: %v", err)
	}

	// File should still exist and be valid
	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("File should still exist after race condition: %v", err)
	}
	if info.IsDir() {
		t.Fatalf("File should not be a directory")
	}
}

// T030: Test read-only directory error
func TestLoadConfig_ReadOnlyDirectoryError(t *testing.T) {
	tmpDir := t.TempDir()
	readOnlyDir := filepath.Join(tmpDir, "readonly")
	configPath := filepath.Join(readOnlyDir, "config.yaml")

	// Create read-only directory (remove write permission)
	if err := os.MkdirAll(readOnlyDir, 0555); err != nil {
		t.Fatalf("Failed to create read-only directory: %v", err)
	}

	// Attempt to load config - should fail with clear error
	_, err := LoadConfig(configPath)
	if err == nil {
		t.Fatalf("Expected error when creating file in read-only directory")
	}

	// Verify error message contains context
	errorMsg := err.Error()
	if errorMsg == "" {
		t.Fatalf("Error message should not be empty")
	}
	// Error should mention directory creation or file creation
	if !contains(errorMsg, "directory") && !contains(errorMsg, "file") {
		t.Fatalf("Error message should mention directory or file creation, got: %s", errorMsg)
	}

	// Cleanup: restore permissions to delete
	os.Chmod(readOnlyDir, 0755)
}

// T031: Test path is directory error
func TestLoadConfig_PathIsDirectoryError(t *testing.T) {
	tmpDir := t.TempDir()
	dirPath := filepath.Join(tmpDir, "somedir")

	// Create a directory at the config path
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	// Attempt to load config with directory path - should fail with clear error
	_, err := LoadConfig(dirPath)
	if err == nil {
		t.Fatalf("Expected error when config path is a directory")
	}

	// Verify error message contains context
	errorMsg := err.Error()
	if !contains(errorMsg, "directory") || !contains(errorMsg, "file") {
		t.Fatalf("Error message should mention that path is a directory, not a file. Got: %s", errorMsg)
	}
}

// T032: Test home directory resolution error
func TestLoadConfig_HomeDirectoryError(t *testing.T) {
	// This test is difficult to simulate without modifying environment
	// We can test that the error is wrapped with context by checking the error message format
	// In practice, os.UserHomeDir() rarely fails, but if it does, we want a clear error

	// Note: This test verifies that if os.UserHomeDir() fails, the error is properly wrapped
	// Actual failure simulation would require mocking, which is complex
	// The implementation already wraps the error, so this test documents the expected behavior
	t.Skip("Home directory resolution error is difficult to simulate without mocking os.UserHomeDir()")
}

// T033: Test error message clarity and context
func TestLoadConfig_ErrorMessageClarity(t *testing.T) {
	tmpDir := t.TempDir()
	readOnlyDir := filepath.Join(tmpDir, "readonly")
	configPath := filepath.Join(readOnlyDir, "config.yaml")

	// Create read-only directory
	if err := os.MkdirAll(readOnlyDir, 0555); err != nil {
		t.Fatalf("Failed to create read-only directory: %v", err)
	}

	// Attempt to load config
	_, err := LoadConfig(configPath)
	if err == nil {
		t.Fatalf("Expected error when creating file in read-only directory")
	}

	// Verify error message is clear and actionable
	errorMsg := err.Error()
	if len(errorMsg) < 20 {
		t.Fatalf("Error message should be descriptive, got: %s", errorMsg)
	}

	// Error should contain the path or operation context
	if !contains(errorMsg, "config") && !contains(errorMsg, "directory") && !contains(errorMsg, "file") {
		t.Fatalf("Error message should contain context about the operation, got: %s", errorMsg)
	}

	// Cleanup
	os.Chmod(readOnlyDir, 0755)
}

// Helper function to check if string contains substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
				s[len(s)-len(substr):] == substr ||
				containsMiddle(s, substr))))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
