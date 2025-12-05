package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestManualCommitWorkflow tests the complete manual commit workflow
// This is an integration test that requires a git repository
func TestManualCommitWorkflow(t *testing.T) {
	// Skip if not in CI or if git is not available
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create a temporary directory for the test repository
	tmpDir, err := os.MkdirTemp("", "gitcomm-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Initialize git repository
	cmd := exec.Command("git", "init", tmpDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Configure git user (required for commits)
	cmd = exec.Command("git", "-C", tmpDir, "config", "user.name", "Test User")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to configure git user.name: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "config", "user.email", "test@example.com")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to configure git user.email: %v", err)
	}

	// Create a test file
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Stage the file
	cmd = exec.Command("git", "-C", tmpDir, "add", testFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	// TODO: Once CLI is implemented, test the actual workflow:
	// 1. Run gitcomm CLI
	// 2. Decline AI assistance
	// 3. Enter commit message components
	// 4. Verify commit is created with correct format

	t.Logf("Test repository created at: %s", tmpDir)
	t.Log("Integration test structure ready - will be completed when CLI workflow is implemented")
}
