package ui

import (
	"testing"
)

func TestPromptState_String(t *testing.T) {
	tests := []struct {
		name     string
		state    PromptState
		expected string
	}{
		{
			name:     "StatePending",
			state:    StatePending,
			expected: "pending",
		},
		{
			name:     "StateActive",
			state:    StateActive,
			expected: "active",
		},
		{
			name:     "StateCompleted",
			state:    StateCompleted,
			expected: "completed",
		},
		{
			name:     "StateCancelled",
			state:    StateCancelled,
			expected: "cancelled",
		},
		{
			name:     "StateError",
			state:    StateError,
			expected: "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.state.String(); got != tt.expected {
				t.Errorf("PromptState.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPromptState_IsTerminal(t *testing.T) {
	tests := []struct {
		name     string
		state    PromptState
		expected bool
	}{
		{
			name:     "StatePending is not terminal",
			state:    StatePending,
			expected: false,
		},
		{
			name:     "StateActive is not terminal",
			state:    StateActive,
			expected: false,
		},
		{
			name:     "StateCompleted is terminal",
			state:    StateCompleted,
			expected: true,
		},
		{
			name:     "StateCancelled is terminal",
			state:    StateCancelled,
			expected: true,
		},
		{
			name:     "StateError is not terminal",
			state:    StateError,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.state.IsTerminal(); got != tt.expected {
				t.Errorf("PromptState.IsTerminal() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPromptState_TransitionFromActiveToCompleted(t *testing.T) {
	// Test state transition from Active to Completed
	state := StateActive

	// Simulate user confirming input
	state = StateCompleted

	if state != StateCompleted {
		t.Errorf("Expected state to be StateCompleted after transition, got %v", state)
	}

	if !state.IsTerminal() {
		t.Error("StateCompleted should be a terminal state")
	}
}

func TestPromptState_TransitionFromActiveToCancelled(t *testing.T) {
	// Test state transition from Active to Cancelled
	state := StateActive

	// Simulate user cancelling
	state = StateCancelled

	if state != StateCancelled {
		t.Errorf("Expected state to be StateCancelled after transition, got %v", state)
	}

	if !state.IsTerminal() {
		t.Error("StateCancelled should be a terminal state")
	}
}

func TestPromptState_TransitionFromErrorToActive(t *testing.T) {
	// Test state transition from Error back to Active (user correcting input)
	state := StateError

	// Simulate user continuing to edit after error
	state = StateActive

	if state != StateActive {
		t.Errorf("Expected state to be StateActive after correction, got %v", state)
	}

	if state.IsTerminal() {
		t.Error("StateActive should not be a terminal state")
	}
}
