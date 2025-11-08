package ai

import (
	"strings"
)

const (
	// Approximate tokens per character for estimation
	// This is a rough estimate - actual tokenization varies by model
	avgCharsPerToken = 4

	// Default max tokens for context (conservative estimate)
	// Actual limits vary by model, but this leaves room for prompt + response
	defaultMaxTokens = 8000

	// Reserve tokens for prompt template and response
	reservedTokens = 2000
)

// TokenCounter provides token counting functionality
type TokenCounter struct {
	maxTokens int
}

// NewTokenCounter creates a new token counter
func NewTokenCounter() *TokenCounter {
	return &TokenCounter{
		maxTokens: defaultMaxTokens - reservedTokens,
	}
}

// EstimateTokens estimates the number of tokens in a string
// This is a rough approximation - real tokenization depends on the model
func (tc *TokenCounter) EstimateTokens(text string) int {
	// Simple estimation: divide character count by average chars per token
	return len(text) / avgCharsPerToken
}

// NeedsSplitting checks if the diff needs to be split
func (tc *TokenCounter) NeedsSplitting(diff string) bool {
	return tc.EstimateTokens(diff) > tc.maxTokens
}

// SplitStrategy determines how to split a large diff
type SplitStrategy struct {
	Chunks    []string
	MaxTokens int
}

// DetermineSplitStrategy analyzes a diff and determines the best split strategy
func (tc *TokenCounter) DetermineSplitStrategy(diff string, fileDiffs []string) *SplitStrategy {
	strategy := &SplitStrategy{
		MaxTokens: tc.maxTokens,
	}

	// If total is under limit, no splitting needed
	if !tc.NeedsSplitting(diff) {
		strategy.Chunks = []string{diff}
		return strategy
	}

	// Try splitting by files first
	if len(fileDiffs) > 1 {
		strategy.Chunks = tc.splitByFiles(fileDiffs)
		return strategy
	}

	// If single large file, split by hunks
	if len(fileDiffs) == 1 {
		strategy.Chunks = tc.splitByHunks(fileDiffs[0])
		return strategy
	}

	// Fallback: just use the original diff (might be too large)
	strategy.Chunks = []string{diff}
	return strategy
}

// splitByFiles splits by grouping files that fit within token limits
func (tc *TokenCounter) splitByFiles(fileDiffs []string) []string {
	var chunks []string
	var currentChunk strings.Builder
	currentTokens := 0

	for _, fileDiff := range fileDiffs {
		fileTokens := tc.EstimateTokens(fileDiff)

		// If this file alone exceeds the limit, split it by hunks
		if fileTokens > tc.maxTokens {
			// Save current chunk if not empty
			if currentChunk.Len() > 0 {
				chunks = append(chunks, currentChunk.String())
				currentChunk.Reset()
				currentTokens = 0
			}

			// Split large file by hunks
			hunkChunks := tc.splitByHunks(fileDiff)
			chunks = append(chunks, hunkChunks...)
			continue
		}

		// If adding this file would exceed limit, start new chunk
		if currentTokens+fileTokens > tc.maxTokens {
			chunks = append(chunks, currentChunk.String())
			currentChunk.Reset()
			currentTokens = 0
		}

		currentChunk.WriteString(fileDiff)
		currentTokens += fileTokens
	}

	// Add remaining chunk
	if currentChunk.Len() > 0 {
		chunks = append(chunks, currentChunk.String())
	}

	return chunks
}

// splitByHunks splits a file diff by its hunks
func (tc *TokenCounter) splitByHunks(fileDiff string) []string {
	// Split by hunk headers (@@)
	lines := strings.Split(fileDiff, "\n")
	var chunks []string
	var currentChunk strings.Builder
	var header strings.Builder
	currentTokens := 0
	inHeader := true

	for _, line := range lines {
		// Detect hunk boundary
		if strings.HasPrefix(line, "@@") {
			// Save previous chunk if it exists
			if currentChunk.Len() > 0 && !inHeader {
				chunks = append(chunks, header.String()+currentChunk.String())
				currentChunk.Reset()
				currentTokens = 0
			}
			inHeader = false
		}

		// Keep file header for all chunks
		if strings.HasPrefix(line, "diff --git") ||
			strings.HasPrefix(line, "index") ||
			strings.HasPrefix(line, "---") ||
			strings.HasPrefix(line, "+++") {
			header.WriteString(line)
			header.WriteString("\n")
			continue
		}

		lineTokens := tc.EstimateTokens(line)

		// If adding this line would exceed limit, start new chunk
		if currentTokens+lineTokens > tc.maxTokens && currentChunk.Len() > 0 {
			chunks = append(chunks, header.String()+currentChunk.String())
			currentChunk.Reset()
			currentTokens = 0
		}

		currentChunk.WriteString(line)
		currentChunk.WriteString("\n")
		currentTokens += lineTokens
	}

	// Add remaining chunk
	if currentChunk.Len() > 0 {
		chunks = append(chunks, header.String()+currentChunk.String())
	}

	return chunks
}

// MergeResults merges multiple commit results into one
func MergeResults(results []*CommitResult) *CommitResult {
	if len(results) == 0 {
		return nil
	}

	if len(results) == 1 {
		return results[0]
	}

	// Use the first result as base
	merged := &CommitResult{
		Type:     results[0].Type,
		Scope:    results[0].Scope,
		Breaking: results[0].Breaking,
	}

	// Merge descriptions
	var descriptions []string
	for _, r := range results {
		if r.Description != "" {
			descriptions = append(descriptions, r.Description)
		}
	}

	if len(descriptions) > 0 {
		// Use the first description as primary
		merged.Description = descriptions[0]
	}

	// Merge bodies
	var bodies []string
	for _, r := range results {
		if r.Body != "" {
			bodies = append(bodies, r.Body)
		}
	}

	if len(bodies) > 0 {
		merged.Body = strings.Join(bodies, "\n\n")
	}

	// If any result is breaking, the merge is breaking
	for _, r := range results {
		if r.Breaking {
			merged.Breaking = true
			break
		}
	}

	return merged
}
