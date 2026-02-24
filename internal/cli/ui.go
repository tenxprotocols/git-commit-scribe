package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/fatih/color"
)

// Color definitions for consistent styling
var (
	// Status colors
	successColor = color.New(color.FgGreen, color.Bold)
	errorColor   = color.New(color.FgRed, color.Bold)
	warningColor = color.New(color.FgYellow, color.Bold)
	infoColor    = color.New(color.FgCyan)

	// Commit message syntax highlighting
	typeColor        = color.New(color.FgMagenta, color.Bold)
	scopeColor       = color.New(color.FgBlue)
	breakingColor    = color.New(color.FgRed, color.Bold)
	descriptionColor = color.New(color.FgWhite, color.Bold)
	bodyColor        = color.New(color.FgWhite)
	emojiColor       = color.New(color.FgYellow)

	// UI elements
	promptColor    = color.New(color.FgCyan, color.Bold)
	dimColor       = color.New(color.FgHiBlack)
	highlightColor = color.New(color.FgYellow, color.Bold)
)

// Spinner represents a progress indicator
type Spinner struct {
	message string
	frames  []string
	stop    chan bool
	done    chan bool
}

// NewSpinner creates a new spinner with the given message
func NewSpinner(message string) *Spinner {
	return &Spinner{
		message: message,
		frames:  []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		stop:    make(chan bool),
		done:    make(chan bool),
	}
}

// Start begins the spinner animation
func (s *Spinner) Start() {
	go func() {
		i := 0
		for {
			select {
			case <-s.stop:
				s.done <- true
				return
			default:
				frame := s.frames[i%len(s.frames)]
				fmt.Printf("\r%s %s", infoColor.Sprint(frame), s.message)
				i++
				time.Sleep(80 * time.Millisecond)
			}
		}
	}()
}

// Stop stops the spinner and clears the line
func (s *Spinner) Stop() {
	s.stop <- true
	<-s.done
	fmt.Print("\r\033[K") // Clear the line
}

// Success stops the spinner and shows a success message
func (s *Spinner) Success(message string) {
	s.Stop()
	fmt.Printf("%s %s\n", successColor.Sprint("✓"), message)
}

// Error stops the spinner and shows an error message
func (s *Spinner) Error(message string) {
	s.Stop()
	fmt.Printf("%s %s\n", errorColor.Sprint("✗"), message)
}

// PrintSuccess prints a success message
func PrintSuccess(message string) {
	fmt.Printf("%s %s\n", successColor.Sprint("✓"), message)
}

// PrintError prints an error message
func PrintError(message string) {
	fmt.Printf("%s %s\n", errorColor.Sprint("✗"), message)
}

// PrintWarning prints a warning message
func PrintWarning(message string) {
	fmt.Printf("%s %s\n", warningColor.Sprint("⚠"), message)
}

// PrintInfo prints an info message
func PrintInfo(message string) {
	fmt.Println(infoColor.Sprint(message))
}

// PrintDim prints dimmed text
func PrintDim(message string) {
	fmt.Println(dimColor.Sprint(message))
}

// PrintCommitMessage prints a commit message with syntax highlighting
func PrintCommitMessage(message string) {
	lines := strings.Split(message, "\n")

	for i, line := range lines {
		switch {
		case i == 0:
			// First line: type(scope): description or type: description or description with emoji
			highlightFirstLine(line)
		case strings.TrimSpace(line) == "":
			// Empty line
			fmt.Println()
		default:
			// Body text
			fmt.Println(bodyColor.Sprint(line))
		}
	}
}

// highlightFirstLine highlights the first line of a commit message
func highlightFirstLine(line string) {
	// Check for emoji at the start
	hasEmoji := false
	emojiEnd := 0

	// Simple emoji detection (any character that's not ASCII)
	for i, r := range line {
		if r > 127 {
			hasEmoji = true
			emojiEnd = i + len(string(r))
		} else {
			break
		}
	}

	if hasEmoji && emojiEnd > 0 {
		fmt.Print(emojiColor.Sprint(line[:emojiEnd]))
		line = strings.TrimSpace(line[emojiEnd:])
		if len(line) > 0 {
			fmt.Print(" ")
		}
	}

	// Check for breaking change marker
	hasBreaking := strings.HasSuffix(strings.Split(line, ":")[0], "!")

	// Parse type(scope): description or type: description
	if idx := strings.Index(line, ":"); idx != -1 {
		typeScope := line[:idx]
		description := strings.TrimSpace(line[idx+1:])

		// Split type and scope
		if scopeIdx := strings.Index(typeScope, "("); scopeIdx != -1 {
			// type(scope) format
			typ := typeScope[:scopeIdx]
			scope := typeScope[scopeIdx:]

			if hasBreaking && strings.HasSuffix(typ, "!") {
				fmt.Print(breakingColor.Sprint(typ[:len(typ)-1]))
				fmt.Print(breakingColor.Sprint("!"))
			} else {
				fmt.Print(typeColor.Sprint(typ))
			}
			fmt.Print(scopeColor.Sprint(scope))
		} else {
			// type only format
			if hasBreaking && strings.HasSuffix(typeScope, "!") {
				fmt.Print(breakingColor.Sprint(typeScope[:len(typeScope)-1]))
				fmt.Print(breakingColor.Sprint("!"))
			} else {
				fmt.Print(typeColor.Sprint(typeScope))
			}
		}

		fmt.Print(dimColor.Sprint(":"))
		fmt.Print(" ")
		fmt.Println(descriptionColor.Sprint(description))
	} else {
		// No type/scope, just description
		fmt.Println(descriptionColor.Sprint(line))
	}
}

// PrintSeparator prints a styled separator line
func PrintSeparator() {
	fmt.Println(dimColor.Sprint("─────────────────────────────────────────"))
}

// PromptYesNo prompts the user with a yes/no question
// Returns true for yes, false for no
func PromptYesNo(question string, defaultYes bool) (bool, error) {
	var prompt string
	if defaultYes {
		prompt = fmt.Sprintf("%s [Y/n]: ", promptColor.Sprint(question))
	} else {
		prompt = fmt.Sprintf("%s [y/N]: ", promptColor.Sprint(question))
	}

	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	response = strings.ToLower(strings.TrimSpace(response))

	// Empty response uses default
	if response == "" {
		return defaultYes, nil
	}

	return response == "y" || response == "yes", nil
}

// PromptChoice prompts the user to choose from multiple options
// Returns the selected option (1-indexed) or 0 if cancelled
// If defaultChoice is > 0, pressing Enter will select that option
func PromptChoice(question string, options []string, defaultChoice int) (int, error) {
	fmt.Println(promptColor.Sprint(question))
	fmt.Println()

	for i, option := range options {
		if i+1 == defaultChoice {
			fmt.Printf("  %s %s %s\n", highlightColor.Sprintf("[%d]", i+1), option, dimColor.Sprint("(default)"))
		} else {
			fmt.Printf("  %s %s\n", highlightColor.Sprintf("[%d]", i+1), option)
		}
	}
	fmt.Println()

	var promptText string
	if defaultChoice > 0 && defaultChoice <= len(options) {
		promptText = fmt.Sprintf("Enter choice [%d] or 0 to cancel: ", defaultChoice)
	} else {
		promptText = "Enter choice (or 0 to cancel): "
	}
	fmt.Print(promptColor.Sprint(promptText))

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}

	response = strings.TrimSpace(response)

	// If empty and there's a default, use it
	if response == "" && defaultChoice > 0 && defaultChoice <= len(options) {
		return defaultChoice, nil
	}

	var choice int
	_, err = fmt.Sscanf(response, "%d", &choice)
	if err != nil {
		return 0, fmt.Errorf("invalid choice")
	}

	if choice < 0 || choice > len(options) {
		return 0, fmt.Errorf("choice out of range")
	}

	return choice, nil
}

// EditText opens the user's default editor to edit text
// Returns the edited text or an error
func EditText(initialText string) (string, error) {
	// Determine the editor to use
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		// Try common editors
		for _, e := range []string{"vim", "vi", "nano", "emacs"} {
			if _, err := exec.LookPath(e); err == nil {
				editor = e
				break
			}
		}
	}
	if editor == "" {
		return "", fmt.Errorf("no editor found. Set EDITOR or VISUAL environment variable")
	}

	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "gscribe-commit-*.txt")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	// Write initial text to temp file
	if _, err := tmpFile.WriteString(initialText); err != nil {
		tmpFile.Close()
		return "", fmt.Errorf("failed to write to temp file: %w", err)
	}
	tmpFile.Close()

	// Open editor — split editor string to handle args (e.g. "codium --wait")
	ctx := context.Background()
	parts := strings.Fields(editor)
	editorArgs := make([]string, len(parts)-1, len(parts))
	copy(editorArgs, parts[1:])
	editorArgs = append(editorArgs, tmpPath)
	cmd := exec.CommandContext(ctx, parts[0], editorArgs...) //nolint:gosec // editor command is from user's EDITOR/VISUAL env var
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("editor failed: %w", err)
	}

	// Read edited content
	content, err := os.ReadFile(tmpPath)
	if err != nil {
		return "", fmt.Errorf("failed to read edited file: %w", err)
	}

	return strings.TrimSpace(string(content)), nil
}

// FormatHeader prints a formatted header
func FormatHeader(title string) {
	fmt.Println()
	fmt.Println(highlightColor.Sprint("═══ " + title + " ═══"))
	fmt.Println()
}

// FormatSubHeader prints a formatted sub-header
func FormatSubHeader(title string) {
	fmt.Println()
	fmt.Println(infoColor.Sprint("─── " + title + " ───"))
	fmt.Println()
}

// PromptText prompts the user for text input
func PromptText(question string) (string, error) {
	fmt.Print(promptColor.Sprint(question + " "))
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(response), nil
}

// PromptMultilineText prompts the user for multi-line text input
// User can finish input with an empty line
func PromptMultilineText(question string) (string, error) {
	fmt.Println(promptColor.Sprint(question))
	PrintDim("(Enter an empty line to finish)")

	reader := bufio.NewReader(os.Stdin)
	var lines []string

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		line = strings.TrimRight(line, "\n\r")

		// Empty line signals end of input
		if line == "" {
			break
		}

		lines = append(lines, line)
	}

	return strings.Join(lines, "\n"), nil
}
