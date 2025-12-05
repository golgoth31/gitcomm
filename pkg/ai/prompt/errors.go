package prompt

import "errors"

var (
	// ErrNilValidator is returned when a nil validator is passed to GenerateSystemMessage
	ErrNilValidator = errors.New("validator cannot be nil")

	// ErrNilRepositoryState is returned when a nil repository state is passed to GenerateUserMessage
	ErrNilRepositoryState = errors.New("repository state cannot be nil")

	// ErrRuleExtractionFailed is returned when validation rules cannot be extracted from validator
	ErrRuleExtractionFailed = errors.New("failed to extract validation rules from validator")
)
