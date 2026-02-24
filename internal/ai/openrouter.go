package ai

import (
	"context"
	"fmt"

	openrouter "github.com/revrost/go-openrouter"
)

// OpenRouterProvider implements the Provider interface using OpenRouter
type OpenRouterProvider struct {
	client *openrouter.Client
	model  string
}

// NewOpenRouterProvider creates a new OpenRouter provider
func NewOpenRouterProvider(apiKey, model string) (*OpenRouterProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	client := openrouter.NewClient(apiKey)

	return &OpenRouterProvider{
		client: client,
		model:  model,
	}, nil
}

// GenerateCommitMessage generates a commit message using OpenRouter
func (p *OpenRouterProvider) GenerateCommitMessage(ctx context.Context, opts GenerateOptions) (*CommitResult, error) {
	prompt := buildPrompt(opts)

	req := openrouter.ChatCompletionRequest{
		Model: p.model,
		Messages: []openrouter.ChatCompletionMessage{
			openrouter.UserMessage(prompt),
		},
		Temperature: float32(opts.Temperature),
	}

	resp, err := p.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate commit message: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from AI")
	}

	content := resp.Choices[0].Message.Content.Text

	// Parse the response
	result, err := parseResponse(content, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	return result, nil
}
