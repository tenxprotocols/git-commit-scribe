package ai

import (
	"context"
)

// Provider is the interface for AI providers
type Provider interface {
	// GenerateCommitMessage generates a commit message from a diff
	GenerateCommitMessage(ctx context.Context, opts GenerateOptions) (*CommitResult, error)
}

// GenerateOptions configures commit message generation
type GenerateOptions struct {
	Diff             string
	Type             string
	Scope            string
	Breaking         bool
	OneLine          bool
	MaxLength        int
	AvailableTypes   map[string]string
	Temperature      float64
}

// CommitResult contains the generated commit message
type CommitResult struct {
	Type        string
	Scope       string
	Description string
	Body        string
	Breaking    bool
}
