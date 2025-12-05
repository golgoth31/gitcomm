package model

// CommitMessage represents a structured commit message conforming to Conventional Commits specification
type CommitMessage struct {
	// Type is the commit type (feat, fix, docs, style, refactor, test, chore, version)
	Type string

	// Scope is the optional scope of the change (e.g., "auth", "api", "cli")
	Scope string

	// Subject is the short description in imperative mood, no period, ≤72 characters
	Subject string

	// Body is the optional detailed explanation, wrapped at 72 chars, ≤320 characters
	Body string

	// Footer is the optional footer lines (issue references, breaking changes, etc.)
	Footer string

	// Signoff indicates whether to include "Signed-off-by" line (default: true)
	Signoff bool
}

// IsEmpty returns true if the commit message has no meaningful content
// A message is considered empty if it lacks both type and subject, or has type but no subject
func (m *CommitMessage) IsEmpty() bool {
	return (m.Type == "" && m.Subject == "") || (m.Type != "" && m.Subject == "")
}
