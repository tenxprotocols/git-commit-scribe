package git

import (
	"bytes"
	"fmt"
	"os/exec"
)

// CreateCommit creates a git commit with the given message
func CreateCommit(message string) error {
	if !isGitRepo() {
		return fmt.Errorf("not a git repository")
	}

	if !hasStagedChanges() {
		return fmt.Errorf("no staged changes to commit")
	}

	cmd := exec.Command("git", "commit", "-m", message)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create commit: %w (stderr: %s)", err, stderr.String())
	}

	return nil
}

// PushChanges pushes commits to the remote repository
func PushChanges() error {
	if !isGitRepo() {
		return fmt.Errorf("not a git repository")
	}

	cmd := exec.Command("git", "push")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to push: %w (stderr: %s)", err, stderr.String())
	}

	return nil
}

// GetCurrentBranch returns the current git branch name
func GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}

	return stdout.String(), nil
}

// HasRemote checks if the repository has a remote configured
func HasRemote() bool {
	cmd := exec.Command("git", "remote")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return false
	}

	return stdout.Len() > 0
}
