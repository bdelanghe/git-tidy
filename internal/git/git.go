// Package git provides functionality for interacting with Git repositories.
// It offers a clean interface for common Git operations, focusing on branch
// and pull request management.
package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// Executor implements the GitOperator interface and provides methods
// for interacting with a Git repository. It uses the git command-line
// tool under the hood.
type Executor struct {
	workDir string // The working directory for Git operations
}

// New creates a new Git executor for the specified working directory.
// The working directory should be within a Git repository.
//
// Example:
//
//	g := git.New("/path/to/repo")
//	if g.IsGitRepository() {
//	    // Perform Git operations
//	}
func New(workDir string) *Executor {
	return &Executor{
		workDir: workDir,
	}
}

// IsGitRepository checks if the current directory is a Git repository
// by attempting to run a git command that requires a repository.
func (g *Executor) IsGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = g.workDir
	return cmd.Run() == nil
}

// GetRepositoryRoot returns the root directory of the Git repository.
// This is useful when the working directory is in a subdirectory of
// the repository.
func (g *Executor) GetRepositoryRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = g.workDir
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get repository root: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GetCurrentBranch returns the name of the current branch.
// If HEAD is detached, it returns the commit hash instead.
func (g *Executor) GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = g.workDir
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// ListBranches returns a list of all local branches in the repository.
// The branches are returned in alphabetical order.
func (g *Executor) ListBranches() ([]string, error) {
	cmd := exec.Command("git", "branch", "--format=%(refname:short)")
	cmd.Dir = g.workDir
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list branches: %w", err)
	}
	branches := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(branches) == 1 && branches[0] == "" {
		return []string{}, nil
	}
	return branches, nil
}

// GetBranchCommit returns the commit hash of the given branch.
// The branch must exist in the local repository.
func (g *Executor) GetBranchCommit(branch string) (string, error) {
	cmd := exec.Command("git", "rev-parse", branch)
	cmd.Dir = g.workDir
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get commit for branch '%s': %w", branch, err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GetPRMergeCommit returns the merge commit hash of the PR for the given branch.
// It searches through the Git history for merge commits that mention the branch
// name in their commit message. This assumes that the PR merge commit follows
// the standard format used by GitHub and similar platforms.
func (g *Executor) GetPRMergeCommit(branch string) (string, error) {
	// Search for merge commits that mention the branch
	cmd := exec.Command("git", "log", "--merges", "--first-parent", "--grep", fmt.Sprintf("Merge pull request.*%s", branch), "--format=%H")
	cmd.Dir = g.workDir
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to find PR merge commit for branch '%s': %w", branch, err)
	}
	
	commits := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(commits) == 0 || (len(commits) == 1 && commits[0] == "") {
		return "", fmt.Errorf("no merge commit found for branch '%s'", branch)
	}
	
	// Return the most recent merge commit
	return commits[0], nil
} 
