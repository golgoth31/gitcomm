package service

import (
	"strings"
	"testing"

	"github.com/golgoth31/gitcomm/internal/model"
)

func TestFormatCommitMessage(t *testing.T) {
	tests := []struct {
		name    string
		message *model.CommitMessage
		want    string
	}{
		{
			name: "simple message with type and subject",
			message: &model.CommitMessage{
				Type:    "feat",
				Subject: "add new feature",
			},
			want: "feat: add new feature",
		},
		{
			name: "message with scope",
			message: &model.CommitMessage{
				Type:    "feat",
				Scope:   "auth",
				Subject: "add user authentication",
			},
			want: "feat(auth): add user authentication",
		},
		{
			name: "message with body",
			message: &model.CommitMessage{
				Type:    "fix",
				Subject: "resolve bug",
				Body:    "Fixed the issue with authentication",
			},
			want: "fix: resolve bug\n\nFixed the issue with authentication",
		},
		{
			name: "message with footer",
			message: &model.CommitMessage{
				Type:    "feat",
				Subject: "add feature",
				Footer:  "Closes #123",
			},
			want: "feat: add feature\n\nCloses #123",
		},
		{
			name: "complete message",
			message: &model.CommitMessage{
				Type:    "feat",
				Scope:   "api",
				Subject: "add endpoint",
				Body:    "Add new REST endpoint for users",
				Footer:  "Closes #456",
			},
			want: "feat(api): add endpoint\n\nAdd new REST endpoint for users\n\nCloses #456",
		},
		{
			name: "message with signoff",
			message: &model.CommitMessage{
				Type:    "feat",
				Subject: "add feature",
				Signoff: true,
			},
			want: "feat: add feature",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewFormattingService()
			got := formatter.Format(tt.message)

			// Normalize line endings for comparison
			got = strings.ReplaceAll(got, "\r\n", "\n")
			tt.want = strings.ReplaceAll(tt.want, "\r\n", "\n")

			if got != tt.want {
				t.Errorf("FormatCommitMessage() = %q, want %q", got, tt.want)
			}
		})
	}
}
