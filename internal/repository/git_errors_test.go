package repository

import (
	"strings"
	"testing"
)

func TestErrGitCommandFailed_Error_EmptyStderr(t *testing.T) {
	const hint = "No additional details from git. Check repository state or run the command manually."
	tests := []struct {
		name   string
		stderr string
	}{
		{"empty", ""},
		{"whitespace only", "   "},
		{"newlines only", "\n\t\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ErrGitCommandFailed{Command: "commit", ExitCode: 1, Stderr: tt.stderr}
			got := e.Error()
			if !strings.Contains(got, hint) {
				t.Errorf("Error() = %q, expected to contain hint %q", got, hint)
			}
			if !strings.Contains(got, "git commit failed (exit 1)") {
				t.Errorf("Error() = %q, expected to contain command and exit", got)
			}
		})
	}
}
