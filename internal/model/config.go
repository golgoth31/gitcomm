package model

import "time"

// CommitOptions represents CLI options for commit creation
type CommitOptions struct {
	// AutoStage automatically stages all unstaged files (-a flag)
	AutoStage bool

	// NoSignoff disables commit signoff (-s flag)
	NoSignoff bool

	// AIProvider overrides the default AI provider
	AIProvider string

	// SkipAI skips AI generation and goes directly to manual input
	SkipAI bool
}

// AIProviderConfig represents configuration for an AI provider
type AIProviderConfig struct {
	// Name is the provider name (openai, anthropic, local)
	Name string

	// APIKey is the API key or authentication token
	APIKey string

	// Model is the optional model identifier (e.g., "gpt-4", "claude-3-opus")
	Model string

	// Endpoint is the optional custom API endpoint (for local models)
	Endpoint string

	// Timeout is the optional request timeout (default: 30s)
	Timeout time.Duration

	// MaxTokens is the optional maximum tokens for response (default: 500)
	MaxTokens int
}
