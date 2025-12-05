package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/golgoth31/gitcomm/internal/model"
	"github.com/golgoth31/gitcomm/internal/repository"
	"github.com/golgoth31/gitcomm/internal/utils"
)

func TestCommitAuthor_FromLocalConfig(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	utils.InitLogger(true)

	// Create temporary directory
	tmpDir := t.TempDir()

	// Initialize git repository
	cmd := exec.Command("git", "init", tmpDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Configure git user in local .git/config
	cmd = exec.Command("git", "-C", tmpDir, "config", "user.name", "Local User")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to configure git user.name: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "config", "user.email", "local@example.com")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to configure git user.email: %v", err)
	}

	// Create repository
	repo, err := repository.NewGitRepository(tmpDir, false)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	// Create and stage a test file
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", testFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	// Create commit
	commitMsg := &model.CommitMessage{
		Type:    "test",
		Subject: "test commit",
	}
	if err := repo.CreateCommit(nil, commitMsg); err != nil {
		t.Fatalf("Failed to create commit: %v", err)
	}

	// Verify commit author
	cmd = exec.Command("git", "-C", tmpDir, "log", "-1", "--format=%an <%ae>")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to get commit author: %v", err)
	}

	author := strings.TrimSpace(string(output))
	if author != "Local User <local@example.com>" {
		t.Errorf("Expected author 'Local User <local@example.com>', got '%s'", author)
	}
}

func TestCommitAuthor_FromGlobalConfig(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	utils.InitLogger(true)

	// Create temporary directory without .git/config
	tmpDir := t.TempDir()

	// Initialize git repository
	cmd := exec.Command("git", "init", tmpDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Save original HOME
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// Create temporary home directory with .gitconfig
	homeDir := t.TempDir()
	os.Setenv("HOME", homeDir)

	globalConfigPath := filepath.Join(homeDir, ".gitconfig")
	configContent := `[user]
	name = Global User
	email = global@example.com
`
	if err := os.WriteFile(globalConfigPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create global config: %v", err)
	}

	// Create repository
	repo, err := repository.NewGitRepository(tmpDir, false)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	// Create and stage a test file
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", testFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	// Create commit
	commitMsg := &model.CommitMessage{
		Type:    "test",
		Subject: "test commit",
	}
	if err := repo.CreateCommit(nil, commitMsg); err != nil {
		t.Fatalf("Failed to create commit: %v", err)
	}

	// Verify commit author
	cmd = exec.Command("git", "-C", tmpDir, "log", "-1", "--format=%an <%ae>")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to get commit author: %v", err)
	}

	author := strings.TrimSpace(string(output))
	if author != "Global User <global@example.com>" {
		t.Errorf("Expected author 'Global User <global@example.com>', got '%s'", author)
	}
}

func TestCommitAuthor_LocalTakesPrecedence(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	utils.InitLogger(true)

	// Create temporary directory
	tmpDir := t.TempDir()

	// Initialize git repository
	cmd := exec.Command("git", "init", tmpDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Configure local config
	cmd = exec.Command("git", "-C", tmpDir, "config", "user.name", "Local User")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to configure git user.name: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "config", "user.email", "local@example.com")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to configure git user.email: %v", err)
	}

	// Save original HOME
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// Create temporary home directory with different global config
	homeDir := t.TempDir()
	os.Setenv("HOME", homeDir)

	globalConfigPath := filepath.Join(homeDir, ".gitconfig")
	configContent := `[user]
	name = Global User
	email = global@example.com
`
	if err := os.WriteFile(globalConfigPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create global config: %v", err)
	}

	// Create repository
	repo, err := repository.NewGitRepository(tmpDir, false)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	// Create and stage a test file
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", testFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	// Create commit
	commitMsg := &model.CommitMessage{
		Type:    "test",
		Subject: "test commit",
	}
	if err := repo.CreateCommit(nil, commitMsg); err != nil {
		t.Fatalf("Failed to create commit: %v", err)
	}

	// Verify commit author uses local config (precedence)
	cmd = exec.Command("git", "-C", tmpDir, "log", "-1", "--format=%an <%ae>")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to get commit author: %v", err)
	}

	author := strings.TrimSpace(string(output))
	if author != "Local User <local@example.com>" {
		t.Errorf("Expected author 'Local User <local@example.com>' (local precedence), got '%s'", author)
	}
}

func TestCommitAuthor_DefaultsWhenConfigMissing(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	utils.InitLogger(true)

	// Create temporary directory
	tmpDir := t.TempDir()

	// Initialize git repository
	cmd := exec.Command("git", "init", tmpDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Don't configure git user - should use defaults

	// Save original HOME
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// Set HOME to non-existent directory to simulate missing global config
	homeDir := filepath.Join(t.TempDir(), "nonexistent")
	os.Setenv("HOME", homeDir)

	// Create repository
	repo, err := repository.NewGitRepository(tmpDir, false)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	// Create and stage a test file
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", testFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	// Create commit
	commitMsg := &model.CommitMessage{
		Type:    "test",
		Subject: "test commit",
	}
	if err := repo.CreateCommit(nil, commitMsg); err != nil {
		t.Fatalf("Failed to create commit: %v", err)
	}

	// Verify commit author uses defaults
	cmd = exec.Command("git", "-C", tmpDir, "log", "-1", "--format=%an <%ae>")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to get commit author: %v", err)
	}

	author := strings.TrimSpace(string(output))
	if author != "gitcomm <gitcomm@local>" {
		t.Errorf("Expected default author 'gitcomm <gitcomm@local>', got '%s'", author)
	}
}
