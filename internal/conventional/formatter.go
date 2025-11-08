package conventional

import (
	"fmt"
	"strings"
)

// CommitMessage represents a conventional commit message
type CommitMessage struct {
	Type        string
	Scope       string
	Breaking    bool
	Description string
	Body        string
	Footer      string
	Emoji       bool
}

// Format formats the commit message according to Conventional Commits v1.0.0
func (c *CommitMessage) Format() string {
	var parts []string

	// Build the header
	header := c.buildHeader()
	parts = append(parts, header)

	// Add body if present
	if c.Body != "" {
		parts = append(parts, "", c.Body)
	}

	// Add footer for breaking changes if needed
	if c.Breaking && c.Footer != "" {
		parts = append(parts, "", c.Footer)
	} else if c.Breaking {
		parts = append(parts, "", "BREAKING CHANGE: "+c.Description)
	} else if c.Footer != "" {
		parts = append(parts, "", c.Footer)
	}

	return strings.Join(parts, "\n")
}

// buildHeader builds the commit header (type, scope, description)
func (c *CommitMessage) buildHeader() string {
	var header strings.Builder

	// Add emoji if requested
	if c.Emoji {
		if emoji := GetTypeEmoji(c.Type); emoji != "" {
			header.WriteString(emoji)
			header.WriteString(" ")
		}
	}

	// Add type
	header.WriteString(c.Type)

	// Add scope if present
	if c.Scope != "" {
		header.WriteString("(")
		header.WriteString(c.Scope)
		header.WriteString(")")
	}

	// Add breaking change indicator
	if c.Breaking {
		header.WriteString("!")
	}

	// Add separator
	header.WriteString(": ")

	// Add description
	header.WriteString(c.Description)

	return header.String()
}

// FormatOneLine formats a one-line commit message
func (c *CommitMessage) FormatOneLine() string {
	return c.buildHeader()
}

// Validate validates the commit message
func (c *CommitMessage) Validate(validTypes map[string]string) error {
	if c.Type == "" {
		return fmt.Errorf("commit type is required")
	}

	if validTypes != nil {
		if _, ok := validTypes[c.Type]; !ok {
			return fmt.Errorf("invalid commit type: %s", c.Type)
		}
	}

	if c.Description == "" {
		return fmt.Errorf("commit description is required")
	}

	return nil
}

// TruncateDescription truncates the description to the specified length
func (c *CommitMessage) TruncateDescription(maxLength int) {
	if len(c.Description) > maxLength {
		c.Description = c.Description[:maxLength-3] + "..."
	}
}

// GetTypeEmoji returns the emoji for a given commit type
func GetTypeEmoji(commitType string) string {
	emojis := map[string]string{
		"feat":     "✨",
		"fix":      "🐛",
		"docs":     "📚",
		"style":    "💎",
		"refactor": "♻️",
		"perf":     "⚡",
		"test":     "✅",
		"build":    "📦",
		"ci":       "👷",
		"chore":    "🔧",
		"revert":   "⏪",
	}

	return emojis[commitType]
}
