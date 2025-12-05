package ui

// PromptState represents the current state of a prompt during its lifecycle
type PromptState int

const (
	// StatePending is the initial state, prompt not yet displayed
	StatePending PromptState = iota
	// StateActive indicates the prompt is displayed and accepting input
	StateActive
	// StateCompleted indicates the user confirmed input/selection successfully
	StateCompleted
	// StateCancelled indicates the user cancelled (Escape key)
	StateCancelled
	// StateError indicates a validation error occurred
	StateError
)

// String returns a human-readable string representation of the prompt state
func (s PromptState) String() string {
	switch s {
	case StatePending:
		return "pending"
	case StateActive:
		return "active"
	case StateCompleted:
		return "completed"
	case StateCancelled:
		return "cancelled"
	case StateError:
		return "error"
	default:
		return "unknown"
	}
}

// IsTerminal returns true if the state is a terminal state (Completed or Cancelled)
func (s PromptState) IsTerminal() bool {
	return s == StateCompleted || s == StateCancelled
}
