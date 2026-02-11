package repository

import (
	"fmt"
	"strings"
	"testing"
)

func TestFormatErrorForDisplay_Nil(t *testing.T) {
	got := FormatErrorForDisplay(nil)
	if got != "" {
		t.Errorf("FormatErrorForDisplay(nil) = %q, want %q", got, "")
	}
}

func TestFormatErrorForDisplay_ErrGitCommandFailed_EmptyStderr(t *testing.T) {
	e := &ErrGitCommandFailed{Command: "commit", ExitCode: 1, Stderr: ""}
	got := FormatErrorForDisplay(e)
	const wantContains = "Details: No additional details from git"
	if !strings.Contains(got, wantContains) {
		t.Errorf("FormatErrorForDisplay(...) = %q, want to contain %q", got, wantContains)
	}
	if !strings.Contains(got, "git commit failed (exit 1)") {
		t.Errorf("FormatErrorForDisplay(...) = %q, want to contain command and exit", got)
	}
}

func TestFormatErrorForDisplay_ErrGitCommandFailed_NonEmptyStderr(t *testing.T) {
	stderr := "fatal: not a git repository"
	e := &ErrGitCommandFailed{Command: "status", ExitCode: 128, Stderr: stderr}
	got := FormatErrorForDisplay(e)
	if !strings.Contains(got, "git status failed (exit 128)") {
		t.Errorf("FormatErrorForDisplay(...) = %q, want to contain command and exit", got)
	}
	if !strings.Contains(got, "Details: "+stderr) {
		t.Errorf("FormatErrorForDisplay(...) = %q, want to contain stderr %q", got, stderr)
	}
}

func TestFormatErrorForDisplay_ErrGitCommandFailed_LongStderr(t *testing.T) {
	stderr := strings.Repeat("x", 2000)
	e := &ErrGitCommandFailed{Command: "commit", ExitCode: 1, Stderr: stderr}
	got := FormatErrorForDisplay(e)
	if len(got) > 1600 {
		t.Errorf("FormatErrorForDisplay(...) length = %d, expected truncated output", len(got))
	}
	if !strings.Contains(got, "… (500 additional characters)") {
		t.Errorf("FormatErrorForDisplay(...) = %q, want to contain truncation suffix", got)
	}
	if !strings.Contains(got, "git commit failed (exit 1)") {
		t.Errorf("FormatErrorForDisplay(...) = %q, want to contain command and exit", got)
	}
}

func TestFormatErrorForDisplay_ErrGitSigningFailed(t *testing.T) {
	err := ErrGitSigningFailed
	got := FormatErrorForDisplay(err)
	want := err.Error()
	if got != want {
		t.Errorf("FormatErrorForDisplay(ErrGitSigningFailed) = %q, want %q", got, want)
	}
}

func TestFormatErrorForDisplay_GenericError(t *testing.T) {
	err := fmt.Errorf("something went wrong")
	got := FormatErrorForDisplay(err)
	want := "something went wrong"
	if got != want {
		t.Errorf("FormatErrorForDisplay(...) = %q, want %q", got, want)
	}
}

func TestFormatErrorForDisplay_WrappedErrGitCommandFailed(t *testing.T) {
	inner := &ErrGitCommandFailed{Command: "add", ExitCode: 1, Stderr: "path not found"}
	err := fmt.Errorf("stage failed: %w", inner)
	got := FormatErrorForDisplay(err)
	if !strings.Contains(got, "Details: path not found") {
		t.Errorf("FormatErrorForDisplay(...) = %q, want to unwrap and show stderr", got)
	}
	if !strings.Contains(got, "git add failed (exit 1)") {
		t.Errorf("FormatErrorForDisplay(...) = %q, want to contain command and exit", got)
	}
}
