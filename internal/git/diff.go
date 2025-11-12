package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// DiffOptions configures diff retrieval
type DiffOptions struct {
	MaxFiles         int
	IgnoreGenerated  bool
	IgnoreWhitespace bool
}

// GetStagedDiff retrieves the staged git diff
func GetStagedDiff(opts DiffOptions) (string, error) {
	// Check if we're in a git repository
	if !isGitRepo() {
		return "", fmt.Errorf("not a git repository")
	}

	// Check if there are staged changes
	if !hasStagedChanges() {
		return "", fmt.Errorf("no staged changes to commit")
	}

	// Get the diff
	args := []string{"diff", "--cached"}
	if opts.IgnoreWhitespace {
		args = append(args, "--ignore-all-space", "--ignore-blank-lines")
	}

	cmd := exec.Command("git", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to get diff: %w (stderr: %s)", err, stderr.String())
	}

	diff := stdout.String()
	if diff == "" {
		return "", fmt.Errorf("no diff content found")
	}

	// Filter generated files if requested
	if opts.IgnoreGenerated {
		diff = filterGeneratedFiles(diff)
	}

	// Check file count
	fileCount := countFilesInDiff(diff)
	if fileCount > opts.MaxFiles {
		return "", fmt.Errorf("too many files changed (%d), max is %d", fileCount, opts.MaxFiles)
	}

	return diff, nil
}

// GetStagedFiles returns a list of staged file paths
func GetStagedFiles() ([]string, error) {
	cmd := exec.Command("git", "diff", "--cached", "--name-only")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to get staged files: %w", err)
	}

	files := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	if len(files) == 1 && files[0] == "" {
		return []string{}, nil
	}

	return files, nil
}

// ValidateRepository validates that we're in a git repository with staged changes
func ValidateRepository() error {
	if !isGitRepo() {
		return fmt.Errorf("not a git repository (or any parent up to mount point /)")
	}

	if !hasStagedChanges() {
		return fmt.Errorf("no changes added to commit (use \"git add\" to stage changes)")
	}

	return nil
}

// isGitRepo checks if current directory is in a git repository
func isGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	return cmd.Run() == nil
}

// hasStagedChanges checks if there are any staged changes
func hasStagedChanges() bool {
	cmd := exec.Command("git", "diff", "--cached", "--quiet")
	// Returns non-zero if there are differences
	return cmd.Run() != nil
}

// filterGeneratedFiles removes generated files from diff
func filterGeneratedFiles(diff string) string {
	// Split diff into individual file diffs
	fileDiffs := SplitDiffByFiles(diff)
	
	// Common patterns for generated files (matching the "diff --git" line)
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`^diff --git a/.*\.lock b/.*\.lock$`),
		regexp.MustCompile(`^diff --git a/.*\.min\.(js|css) b/.*\.min\.(js|css)$`),
		regexp.MustCompile(`^diff --git a/.*\.generated\. b/.*\.generated\.$`),
		regexp.MustCompile(`^diff --git a/.*package-lock\.json b/.*package-lock\.json$`),
		regexp.MustCompile(`^diff --git a/.*yarn\.lock b/.*yarn\.lock$`),
		regexp.MustCompile(`^diff --git a/.*go\.sum b/.*go\.sum$`),
		regexp.MustCompile(`^diff --git a/.*Cargo\.lock b/.*Cargo\.lock$`),
	}

	var filtered []string
	for _, fileDiff := range fileDiffs {
		isGenerated := false
		firstLine := strings.Split(fileDiff, "\n")[0]
		
		for _, pattern := range patterns {
			if pattern.MatchString(firstLine) {
				isGenerated = true
				break
			}
		}
		
		if !isGenerated {
			filtered = append(filtered, fileDiff)
		}
	}

	return strings.Join(filtered, "")
}

// countFilesInDiff counts the number of files in a diff
func countFilesInDiff(diff string) int {
	pattern := regexp.MustCompile(`(?m)^diff --git`)
	matches := pattern.FindAllString(diff, -1)
	return len(matches)
}

// SplitDiffByFiles splits a large diff into per-file diffs
func SplitDiffByFiles(diff string) []string {
	pattern := regexp.MustCompile(`(?m)^diff --git`)
	indices := pattern.FindAllStringIndex(diff, -1)

	if len(indices) == 0 {
		return []string{}
	}

	var diffs []string
	for i, idx := range indices {
		start := idx[0]
		var end int
		if i+1 < len(indices) {
			end = indices[i+1][0]
		} else {
			end = len(diff)
		}
		diffs = append(diffs, diff[start:end])
	}

	return diffs
}
