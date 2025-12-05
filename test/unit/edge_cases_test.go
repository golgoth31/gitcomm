package unit

import (
	"strings"
	"testing"

	"github.com/golgoth31/gitcomm/internal/model"
	"github.com/golgoth31/gitcomm/internal/service"
	"github.com/golgoth31/gitcomm/pkg/conventional"
)

func TestValidator_EdgeCases(t *testing.T) {
	validator := conventional.NewValidator()

	tests := []struct {
		name      string
		message   *model.CommitMessage
		wantValid bool
	}{
		{
			name: "subject exactly 72 characters",
			message: &model.CommitMessage{
				Type:    "feat",
				Subject: string(make([]byte, 72)),
			},
			wantValid: true,
		},
		{
			name: "subject 73 characters",
			message: &model.CommitMessage{
				Type:    "feat",
				Subject: string(make([]byte, 73)),
			},
			wantValid: false,
		},
		{
			name: "body exactly 320 characters",
			message: &model.CommitMessage{
				Type:    "feat",
				Subject: "test",
				Body:    string(make([]byte, 320)),
			},
			wantValid: true,
		},
		{
			name: "body 321 characters",
			message: &model.CommitMessage{
				Type:    "feat",
				Subject: "test",
				Body:    string(make([]byte, 321)),
			},
			wantValid: false,
		},
		{
			name: "scope with valid characters",
			message: &model.CommitMessage{
				Type:    "feat",
				Scope:   "auth-api",
				Subject: "test",
			},
			wantValid: true,
		},
		{
			name: "scope with underscores",
			message: &model.CommitMessage{
				Type:    "feat",
				Scope:   "auth_api",
				Subject: "test",
			},
			wantValid: true,
		},
		{
			name: "scope with numbers",
			message: &model.CommitMessage{
				Type:    "feat",
				Scope:   "api-v2",
				Subject: "test",
			},
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, _ := validator.Validate(tt.message)
			if valid != tt.wantValid {
				t.Errorf("Validator.Validate() valid = %v, want %v", valid, tt.wantValid)
			}
		})
	}
}

func TestFormattingService_EdgeCases(t *testing.T) {
	formatter := &service.FormattingService{}

	tests := []struct {
		name         string
		message      *model.CommitMessage
		wantContains []string
	}{
		{
			name: "message with all components",
			message: &model.CommitMessage{
				Type:    "feat",
				Scope:   "api",
				Subject: "add endpoint",
				Body:    "Detailed description",
				Footer:  "Closes #123",
			},
			wantContains: []string{"feat(api): add endpoint", "Detailed description", "Closes #123"},
		},
		{
			name: "message with only type and subject",
			message: &model.CommitMessage{
				Type:    "fix",
				Subject: "resolve bug",
			},
			wantContains: []string{"fix: resolve bug"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatted := formatter.Format(tt.message)
			for _, want := range tt.wantContains {
				if !strings.Contains(formatted, want) {
					t.Errorf("Format() = %q, want to contain %q", formatted, want)
				}
			}
		})
	}
}
