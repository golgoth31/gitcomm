package config

// Error types for git config extraction
// These are used internally for error handling but are not returned to callers
// per FR-009 (silent ignore of errors)

// ConfigError represents an error during config extraction
type ConfigError struct {
	Message string
	Err     error
}

func (e *ConfigError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *ConfigError) Unwrap() error {
	return e.Err
}
