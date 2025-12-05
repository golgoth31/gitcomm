package ui

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
)

// AIMessageAcceptance represents the user's choice when presented with an AI-generated commit message
type AIMessageAcceptance int

const (
	// AcceptAndCommit indicates the user wants to commit immediately with the AI message
	AcceptAndCommit AIMessageAcceptance = iota
	// AcceptAndEdit indicates the user wants to edit the AI message before committing
	AcceptAndEdit
	// Reject indicates the user wants to reject the AI message and start over
	Reject
)

// String returns a human-readable string representation of the acceptance value
func (a AIMessageAcceptance) String() string {
	switch a {
	case AcceptAndCommit:
		return "accept and commit"
	case AcceptAndEdit:
		return "accept and edit"
	case Reject:
		return "reject"
	default:
		return "unknown"
	}
}

// PrefilledCommitMessage represents a commit message structure where fields are populated
// with values from an AI-generated message, ready for user editing
type PrefilledCommitMessage struct {
	Type    string // Pre-filled commit type from AI message
	Scope   string // Pre-filled scope from AI message (may be empty)
	Subject string // Pre-filled subject from AI message
	Body    string // Pre-filled body from AI message (may be empty)
	Footer  string // Pre-filled footer from AI message (may be empty)
}

// PromptScope prompts the user for commit scope (optional)
func PromptScope(reader *bufio.Reader) (string, error) {
	var scope string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Scope (optional, press Enter to skip)").
				Value(&scope),
		),
	)

	if err := form.Run(); err != nil {
		return "", fmt.Errorf("scope input cancelled: %w", err)
	}

	// Print post-validation summary line
	printPostValidationSummary("Scope (optional)", scope)

	return scope, nil
}

// PromptScopeWithDefault prompts the user for commit scope with a default value
func PromptScopeWithDefault(reader *bufio.Reader, defaultValue string) (string, error) {
	scope := defaultValue

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Scope").
				Value(&scope),
		),
	)

	if err := form.Run(); err != nil {
		return "", fmt.Errorf("scope input cancelled: %w", err)
	}

	// If empty and default exists, return default
	if scope == "" && defaultValue != "" {
		scope = defaultValue
	}

	// Print post-validation summary line
	printPostValidationSummary("Scope", scope)

	return scope, nil
}

// PromptSubject prompts the user for commit subject (required)
func PromptSubject(reader *bufio.Reader) (string, error) {
	var subject string

	validator := func(value string) error {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("subject cannot be empty")
		}
		return nil
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Subject (required)").
				Value(&subject).
				Validate(validator),
		),
	)

	if err := form.Run(); err != nil {
		return "", fmt.Errorf("subject input cancelled: %w", err)
	}

	subject = strings.TrimSpace(subject)

	// Print post-validation summary line
	printPostValidationSummary("Subject (required)", subject)

	return subject, nil
}

// PromptSubjectWithDefault prompts the user for commit subject with a default value
func PromptSubjectWithDefault(reader *bufio.Reader, defaultValue string) (string, error) {
	subject := defaultValue

	validator := func(value string) error {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" && defaultValue == "" {
			return fmt.Errorf("subject cannot be empty")
		}
		return nil
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Subject").
				Value(&subject).
				Validate(validator),
		),
	)

	if err := form.Run(); err != nil {
		return "", fmt.Errorf("subject input cancelled: %w", err)
	}

	subject = strings.TrimSpace(subject)
	// If empty and default exists, return default
	if subject == "" && defaultValue != "" {
		subject = defaultValue
	}

	// Print post-validation summary line
	printPostValidationSummary("Subject", subject)

	return subject, nil
}

// PromptBody prompts the user for commit body (optional) using multiline input
func PromptBody(reader *bufio.Reader) (string, error) {
	var body string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title("Body").
				Value(&body),
		),
	)

	if err := form.Run(); err != nil {
		return "", fmt.Errorf("body input cancelled: %w", err)
	}

	// Warn if body is too long (optional validation)
	if len(body) > 320 {
		fmt.Printf("Warning: Body is %d characters (recommended: ≤320). Continue? (y/n): ", len(body))
		confirm, _ := reader.ReadString('\n')
		if strings.TrimSpace(strings.ToLower(confirm)) != "y" {
			return "", fmt.Errorf("body too long, user cancelled")
		}
	}

	// Print post-validation summary line (truncated for multiline)
	printPostValidationSummary("Body", body)

	return body, nil
}

// PromptBodyWithDefault prompts the user for commit body with a default value pre-populated
func PromptBodyWithDefault(reader *bufio.Reader, defaultValue string) (string, error) {
	body := defaultValue

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title("Body").
				Value(&body),
		),
	)

	if err := form.Run(); err != nil {
		return "", fmt.Errorf("body input cancelled: %w", err)
	}

	// Warn if body is too long (optional validation)
	if len(body) > 320 {
		fmt.Printf("Warning: Body is %d characters (recommended: ≤320). Continue? (y/n): ", len(body))
		confirm, _ := reader.ReadString('\n')
		if strings.TrimSpace(strings.ToLower(confirm)) != "y" {
			return "", fmt.Errorf("body too long, user cancelled")
		}
	}

	// Print post-validation summary line (truncated for multiline)
	printPostValidationSummary("Body", body)

	return body, nil
}

// PromptFooter prompts the user for commit footer (optional) using multiline input
func PromptFooter(reader *bufio.Reader) (string, error) {
	var footer string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title("Footer").
				Value(&footer),
		),
	)

	if err := form.Run(); err != nil {
		return "", fmt.Errorf("footer input cancelled: %w", err)
	}

	// Print post-validation summary line (truncated for multiline)
	printPostValidationSummary("Footer", footer)

	return footer, nil
}

// PromptFooterWithDefault prompts the user for commit footer with a default value pre-populated
func PromptFooterWithDefault(reader *bufio.Reader, defaultValue string) (string, error) {
	footer := defaultValue

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title("Footer").
				Value(&footer),
		),
	)

	if err := form.Run(); err != nil {
		return "", fmt.Errorf("footer input cancelled: %w", err)
	}

	// Print post-validation summary line (truncated for multiline)
	printPostValidationSummary("Footer", footer)

	return footer, nil
}

// PromptEmptyCommit prompts the user to confirm creating an empty commit
func PromptEmptyCommit(reader *bufio.Reader) (bool, error) {
	var confirm bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("No changes detected. Create an empty commit?").
				Value(&confirm),
		),
	)

	if err := form.Run(); err != nil {
		return false, fmt.Errorf("empty commit prompt cancelled: %w", err)
	}

	// Print post-validation summary line
	printPostValidationSummary("No changes detected. Create an empty commit?", confirm)

	return confirm, nil
}

// PromptConfirm prompts the user to confirm an action
func PromptConfirm(reader *bufio.Reader, message string, defaultValue bool) (bool, error) {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(message).
				Value(&defaultValue),
		),
	)

	if err := form.Run(); err != nil {
		return false, fmt.Errorf("confirm prompt cancelled: %w", err)
	}

	// Print post-validation summary line
	printPostValidationSummary(message, defaultValue)

	return defaultValue, nil
}

// PromptCommitType prompts the user for commit type using an interactive select list
func PromptCommitType(reader *bufio.Reader) (string, error) {
	var commitType string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose a type").
				Options(
					huh.NewOption("feat", "feat"),
					huh.NewOption("fix", "fix"),
					huh.NewOption("docs", "docs"),
					huh.NewOption("style", "style"),
					huh.NewOption("refactor", "refactor"),
					huh.NewOption("test", "test"),
					huh.NewOption("chore", "chore"),
					huh.NewOption("version", "version"),
				).
				Value(&commitType),
		),
	)

	if err := form.Run(); err != nil {
		return "", fmt.Errorf("commit type selection cancelled: %w", err)
	}

	// Print post-validation summary line
	printPostValidationSummary("Choose a type", commitType)

	return commitType, nil
}

// PromptCommitTypeWithPreselection prompts the user for commit type with a pre-selected type
func PromptCommitTypeWithPreselection(reader *bufio.Reader, preselectedType string) (string, error) {
	commitType := preselectedType

	options := []huh.Option[string]{
		huh.NewOption("feat", "feat"),
		huh.NewOption("fix", "fix"),
		huh.NewOption("docs", "docs"),
		huh.NewOption("style", "style"),
		huh.NewOption("refactor", "refactor"),
		huh.NewOption("test", "test"),
		huh.NewOption("chore", "chore"),
		huh.NewOption("version", "version"),
	}

	// Mark preselected option as selected
	for i := range options {
		if options[i].Value == preselectedType {
			options[i] = options[i].Selected(true)
			break
		}
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose a type").
				Options(options...).
				Value(&commitType),
		),
	)

	if err := form.Run(); err != nil {
		return "", fmt.Errorf("commit type selection cancelled: %w", err)
	}

	// Print post-validation summary line
	printPostValidationSummary("Choose a type", commitType)

	return commitType, nil
}

// PromptAIUsage prompts the user to choose whether to use AI
func PromptAIUsage(reader *bufio.Reader, tokenCount int) (bool, error) {
	var useAI bool = true // Default to "yes" (true) for AI usage

	estimatedTokens := fmt.Sprintf("Estimated tokens: %d", tokenCount)
	message := "Use AI to generate commit message?"

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().Title(estimatedTokens),
			huh.NewConfirm().
				Title(message).
				Value(&useAI),
		),
	)

	if err := form.Run(); err != nil {
		return false, fmt.Errorf("AI usage prompt cancelled: %w", err)
	}

	aiOutputMessage := fmt.Sprintf("Use AI to generate commit message for %d tokens?", tokenCount)
	// Print post-validation summary line
	printPostValidationSummary(aiOutputMessage, useAI)

	return useAI, nil
}

// PromptAIMessageAcceptance prompts the user to accept or reject AI-generated message
// Deprecated: Use PromptAIMessageAcceptanceOptions instead
func PromptAIMessageAcceptance(reader *bufio.Reader, message string) (bool, error) {
	fmt.Println("\n--- AI Generated Message ---")
	fmt.Println(message)
	fmt.Println("---")
	fmt.Print("Accept this message and commit? (Y/n): ")
	input, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("failed to read input: %w", err)
	}
	response := strings.TrimSpace(strings.ToLower(input))
	return response == "" || response == "y" || response == "yes", nil
}

// PromptAIMessageAcceptanceOptions prompts the user to choose from three options when presented with an AI-generated commit message
func PromptAIMessageAcceptanceOptions(reader *bufio.Reader, message string) (AIMessageAcceptance, error) {
	fmt.Println("\n--- AI Generated Message ---")
	fmt.Println(message)
	fmt.Println("---")

	var choice string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Options").
				Options(
					huh.NewOption("Accept and commit directly", "accept-commit"),
					huh.NewOption("Accept and edit", "accept-edit"),
					huh.NewOption("Reject", "reject"),
				).
				Value(&choice),
		),
	)

	if err := form.Run(); err != nil {
		return 0, fmt.Errorf("AI message acceptance prompt cancelled: %w", err)
	}

	var acceptance AIMessageAcceptance
	switch choice {
	case "accept-commit":
		acceptance = AcceptAndCommit
	case "accept-edit":
		acceptance = AcceptAndEdit
	case "reject":
		acceptance = Reject
	default:
		return 0, fmt.Errorf("invalid choice: %s", choice)
	}

	// Print post-validation summary line
	var choiceStr string
	switch acceptance {
	case AcceptAndCommit:
		choiceStr = "Accept and commit directly"
	case AcceptAndEdit:
		choiceStr = "Accept and edit"
	case Reject:
		choiceStr = "Reject"
	}
	printPostValidationSummary("Options", choiceStr)

	return acceptance, nil
}

// PromptAIMessageEdit prompts the user to edit or use AI message with warning
func PromptAIMessageEdit(reader *bufio.Reader, errors []string) (bool, error) {
	var edit bool = true // Default to "yes" (edit) when there are validation errors

	var messageBuilder strings.Builder
	messageBuilder.WriteString("\nValidation errors found:\n")
	for _, e := range errors {
		messageBuilder.WriteString("  - " + e + "\n")
	}
	messageBuilder.WriteString("\nEdit the message? (y=edit, n=use as-is)")

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(messageBuilder.String()).
				Value(&edit),
		),
	)

	if err := form.Run(); err != nil {
		return false, fmt.Errorf("AI message edit prompt cancelled: %w", err)
	}

	// Print post-validation summary line
	printPostValidationSummary("Edit the message?", edit)

	// Return true if user selected "yes" (edit), false if "no" (use as-is)
	return edit, nil
}

// CommitFailureChoice represents the user's choice when commit fails
type CommitFailureChoice int

const (
	// RetryCommit indicates the user wants to retry the commit with the same message
	RetryCommit CommitFailureChoice = iota
	// EditMessage indicates the user wants to edit the message
	EditMessage
	// CancelCommit indicates the user wants to cancel the commit
	CancelCommit
)

// PromptCommitFailureChoice prompts the user to choose an action when commit fails
func PromptCommitFailureChoice(reader *bufio.Reader) (CommitFailureChoice, error) {
	var choice string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Options").
				Options(
					huh.NewOption("Retry commit", "retry"),
					huh.NewOption("Edit message", "edit"),
					huh.NewOption("Cancel", "cancel"),
				).
				Value(&choice),
		),
	)

	if err := form.Run(); err != nil {
		return 0, fmt.Errorf("commit failure choice prompt cancelled: %w", err)
	}

	var failureChoice CommitFailureChoice
	switch choice {
	case "retry":
		failureChoice = RetryCommit
	case "edit":
		failureChoice = EditMessage
	case "cancel":
		failureChoice = CancelCommit
	default:
		return 0, fmt.Errorf("invalid choice: %s", choice)
	}

	// Print post-validation summary line
	var choiceStr string
	switch failureChoice {
	case RetryCommit:
		choiceStr = "Retry commit"
	case EditMessage:
		choiceStr = "Edit message"
	case CancelCommit:
		choiceStr = "Cancel"
	}
	printPostValidationSummary("Options", choiceStr)

	return failureChoice, nil
}

// PromptRejectChoice prompts the user to choose between generating a new AI message or proceeding with manual input
func PromptRejectChoice(reader *bufio.Reader) (bool, error) {
	var generateNew bool = true // Default to "yes" (generate new AI message)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Generate new AI message? (y=new AI, n=manual input)").
				Value(&generateNew),
		),
	)

	if err := form.Run(); err != nil {
		return false, fmt.Errorf("reject choice prompt cancelled: %w", err)
	}

	// Print post-validation summary line
	printPostValidationSummary("Generate new AI message?", generateNew)

	// Return true if user selected "yes" (generate new AI), false if "no" (manual input)
	return generateNew, nil
}
