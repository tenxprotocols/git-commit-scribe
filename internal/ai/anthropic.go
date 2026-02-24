package ai

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// AnthropicProvider implements the Provider interface using the Anthropic API
type AnthropicProvider struct {
	client anthropic.Client
	model  string
}

// NewAnthropicProvider creates a new Anthropic provider
func NewAnthropicProvider(apiKey, model string) (*AnthropicProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	client := anthropic.NewClient(option.WithAPIKey(apiKey))

	return &AnthropicProvider{
		client: client,
		model:  model,
	}, nil
}

// GenerateCommitMessage generates a commit message using the Anthropic Messages API
func (p *AnthropicProvider) GenerateCommitMessage(ctx context.Context, opts GenerateOptions) (*CommitResult, error) {
	prompt := buildPrompt(opts)

	message, err := p.client.Messages.New(ctx, anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
		Model:       anthropic.Model(p.model),
		Temperature: anthropic.Float(opts.Temperature),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate commit message: %w", err)
	}

	// Extract text from response content blocks
	var content string
	for _, block := range message.Content {
		switch b := block.AsAny().(type) {
		case anthropic.TextBlock:
			content += b.Text
		}
	}

	if content == "" {
		return nil, fmt.Errorf("no text content in response")
	}

	// Parse the response using shared logic
	result, err := parseResponse(content, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	return result, nil
}
