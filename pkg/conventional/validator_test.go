package conventional

import (
	"testing"

	"github.com/golgoth31/gitcomm/internal/model"
)

func TestValidator_Validate(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name       string
		message    *model.CommitMessage
		wantValid  bool
		wantErrors int
	}{
		{
			name: "valid commit message",
			message: &model.CommitMessage{
				Type:    "feat",
				Scope:   "auth",
				Subject: "add user authentication",
				Body:    "Implement JWT-based authentication",
			},
			wantValid:  true,
			wantErrors: 0,
		},
		{
			name: "invalid type",
			message: &model.CommitMessage{
				Type:    "invalid",
				Subject: "add feature",
			},
			wantValid:  false,
			wantErrors: 1,
		},
		{
			name: "empty subject",
			message: &model.CommitMessage{
				Type:    "feat",
				Subject: "",
			},
			wantValid:  false,
			wantErrors: 1,
		},
		{
			name: "subject too long",
			message: &model.CommitMessage{
				Type:    "feat",
				Subject: "this is a very long subject that exceeds the 72 character limit for Conventional Commits specification",
			},
			wantValid:  false,
			wantErrors: 1,
		},
		{
			name: "body too long",
			message: &model.CommitMessage{
				Type:    "feat",
				Subject: "add feature",
				Body:    string(make([]byte, 321)), // 321 characters
			},
			wantValid:  false,
			wantErrors: 1,
		},
		{
			name: "invalid scope",
			message: &model.CommitMessage{
				Type:    "feat",
				Scope:   "invalid scope!",
				Subject: "add feature",
			},
			wantValid:  false,
			wantErrors: 1,
		},
		{
			name: "valid with empty scope",
			message: &model.CommitMessage{
				Type:    "feat",
				Scope:   "",
				Subject: "add feature",
			},
			wantValid:  true,
			wantErrors: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, errors := validator.Validate(tt.message)
			if valid != tt.wantValid {
				t.Errorf("Validator.Validate() valid = %v, want %v", valid, tt.wantValid)
			}
			if len(errors) != tt.wantErrors {
				t.Errorf("Validator.Validate() errors = %v, want %d errors", errors, tt.wantErrors)
			}
		})
	}
}

func TestValidator_Validate_AllTypes(t *testing.T) {
	validator := NewValidator()
	validTypes := []string{"feat", "fix", "docs", "style", "refactor", "test", "chore", "version"}

	for _, typ := range validTypes {
		t.Run(typ, func(t *testing.T) {
			message := &model.CommitMessage{
				Type:    typ,
				Subject: "test subject",
			}
			valid, errors := validator.Validate(message)
			if !valid {
				t.Errorf("Type %s should be valid, got errors: %v", typ, errors)
			}
		})
	}
}

func TestValidator_GetValidTypes(t *testing.T) {
	validator := NewValidator()
	types := validator.GetValidTypes()

	expectedTypes := []string{"feat", "fix", "docs", "style", "refactor", "test", "chore", "version"}
	if len(types) != len(expectedTypes) {
		t.Errorf("GetValidTypes() returned %d types, want %d", len(types), len(expectedTypes))
	}

	for _, expectedType := range expectedTypes {
		found := false
		for _, typ := range types {
			if typ == expectedType {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("GetValidTypes() missing type: %s", expectedType)
		}
	}
}

func TestValidator_GetSubjectMaxLength(t *testing.T) {
	validator := NewValidator()
	maxLength := validator.GetSubjectMaxLength()

	if maxLength != 72 {
		t.Errorf("GetSubjectMaxLength() = %d, want 72", maxLength)
	}
}

func TestValidator_GetBodyMaxLength(t *testing.T) {
	validator := NewValidator()
	maxLength := validator.GetBodyMaxLength()

	if maxLength != 320 {
		t.Errorf("GetBodyMaxLength() = %d, want 320", maxLength)
	}
}

func TestValidator_GetScopeFormatDescription(t *testing.T) {
	validator := NewValidator()
	description := validator.GetScopeFormatDescription()

	expected := "alphanumeric, hyphens, underscores only"
	if description != expected {
		t.Errorf("GetScopeFormatDescription() = %q, want %q", description, expected)
	}
}
