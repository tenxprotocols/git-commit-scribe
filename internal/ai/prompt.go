package ai

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/template"
)

//go:embed prompts/default.txt
var defaultPromptFS embed.FS

// buildPrompt builds the prompt for the AI using templates
func buildPrompt(opts GenerateOptions) string {
	// Get the prompt template
	promptTemplate, err := loadPromptTemplate(opts.CustomPrompt)
	if err != nil {
		// Fallback to hardcoded default if template loading fails
		return buildDefaultPrompt(opts)
	}

	// Parse the template
	tmpl, err := template.New("prompt").Parse(promptTemplate)
	if err != nil {
		return buildDefaultPrompt(opts)
	}

	// Execute the template with options
	var sb strings.Builder
	if err := tmpl.Execute(&sb, opts); err != nil {
		return buildDefaultPrompt(opts)
	}

	return sb.String()
}

// loadPromptTemplate loads a prompt template from custom source or default
func loadPromptTemplate(customPrompt string) (string, error) {
	// If custom prompt provided, use it
	if customPrompt != "" {
		// Check if it's a file path
		if _, err := os.Stat(customPrompt); err == nil {
			content, err := os.ReadFile(customPrompt)
			if err != nil {
				return "", fmt.Errorf("failed to read custom prompt file: %w", err)
			}
			return string(content), nil
		}
		// Otherwise treat it as inline prompt text
		return customPrompt, nil
	}

	// Load default embedded prompt
	content, err := defaultPromptFS.ReadFile("prompts/default.txt")
	if err != nil {
		return "", fmt.Errorf("failed to read default prompt: %w", err)
	}

	return string(content), nil
}

// buildDefaultPrompt is a fallback that builds the prompt without templates
func buildDefaultPrompt(opts GenerateOptions) string {
	var sb strings.Builder

	sb.WriteString("You are an expert at writing conventional commit messages. ")
	sb.WriteString("Analyze the following git diff and generate a conventional commit message.\n\n")

	// Add available types
	if len(opts.AvailableTypes) > 0 {
		sb.WriteString("Available commit types:\n")
		for t, desc := range opts.AvailableTypes {
			sb.WriteString(fmt.Sprintf("- %s: %s\n", t, desc))
		}
		sb.WriteString("\n")
	}

	// Add constraints
	sb.WriteString("Requirements:\n")
	sb.WriteString("- Follow Conventional Commits v1.0.0 specification\n")
	sb.WriteString("- Choose the most appropriate type from the list above\n")

	if opts.Type != "" {
		sb.WriteString(fmt.Sprintf("- Use commit type: %s\n", opts.Type))
	}

	if opts.Scope != "" {
		sb.WriteString(fmt.Sprintf("- Use scope: %s\n", opts.Scope))
	} else {
		sb.WriteString("- Include a scope if appropriate (optional)\n")
	}

	if opts.Breaking {
		sb.WriteString("- Mark as BREAKING CHANGE\n")
	}

	if opts.OneLine {
		sb.WriteString("- Generate ONLY a one-line commit message (no body)\n")
	} else {
		sb.WriteString("- Include a body if the changes are complex\n")
	}

	if opts.MaxLength > 0 {
		sb.WriteString(fmt.Sprintf("- Keep the description under %d characters\n", opts.MaxLength))
	}

	sb.WriteString("- Be concise and specific\n")
	sb.WriteString("- Use imperative mood (e.g., 'add' not 'added')\n\n")

	sb.WriteString("Return your response as JSON in this exact format:\n")
	sb.WriteString("{\n")
	sb.WriteString("  \"type\": \"feat\",\n")
	sb.WriteString("  \"scope\": \"api\" (or empty string if no scope),\n")
	sb.WriteString("  \"description\": \"add user authentication\",\n")
	sb.WriteString("  \"body\": \"Detailed explanation...\" (or empty string if one-line),\n")
	sb.WriteString("  \"breaking\": false\n")
	sb.WriteString("}\n\n")

	sb.WriteString("Git diff:\n")
	sb.WriteString("```\n")
	sb.WriteString(opts.Diff)
	sb.WriteString("\n```\n")

	// Add additional context if provided
	if opts.AdditionalContext != "" {
		sb.WriteString("\nAdditional context from user:\n")
		sb.WriteString(opts.AdditionalContext)
		sb.WriteString("\n")
	}

	return sb.String()
}

// parseResponse parses the AI response into a CommitResult
func parseResponse(content string, opts GenerateOptions) (*CommitResult, error) {
	// Try to extract JSON from the response
	jsonStr := extractJSON(content)
	if jsonStr == "" {
		return nil, fmt.Errorf("no JSON found in response")
	}

	var result CommitResult
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Apply overrides from opts
	if opts.Type != "" {
		result.Type = opts.Type
	}
	if opts.Scope != "" {
		result.Scope = opts.Scope
	}
	if opts.Breaking {
		result.Breaking = true
	}

	// Validate
	if result.Type == "" {
		return nil, fmt.Errorf("AI did not provide a commit type")
	}
	if result.Description == "" {
		return nil, fmt.Errorf("AI did not provide a description")
	}

	return &result, nil
}

// extractJSON extracts JSON from markdown code blocks or plain text
func extractJSON(content string) string {
	// Try to find JSON in code block
	if idx := strings.Index(content, "```json"); idx != -1 {
		start := idx + 7
		if end := strings.Index(content[start:], "```"); end != -1 {
			return strings.TrimSpace(content[start : start+end])
		}
	}

	// Try to find JSON in regular code block
	if idx := strings.Index(content, "```"); idx != -1 {
		start := idx + 3
		if end := strings.Index(content[start:], "```"); end != -1 {
			return strings.TrimSpace(content[start : start+end])
		}
	}

	// Try to find JSON object directly
	if idx := strings.Index(content, "{"); idx != -1 {
		if end := strings.LastIndex(content, "}"); end != -1 && end > idx {
			return strings.TrimSpace(content[idx : end+1])
		}
	}

	return ""
}
