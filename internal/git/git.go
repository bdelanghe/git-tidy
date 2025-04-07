package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// Executor represents a Git repository executor
type Executor struct {
	workDir string
}

// New creates a new Git executor
func New(workDir string) *Executor {
	return &Executor{
		workDir: workDir,
	}
}

// IsGitRepository checks if the current directory is a Git repository
func (g *Executor) IsGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = g.workDir
	return cmd.Run() == nil
}

// GetRepositoryRoot returns the root directory of the Git repository
func (g *Executor) GetRepositoryRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = g.workDir
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// GetCurrentBranch returns the name of the current branch
func (g *Executor) GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = g.workDir
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// ListBranches returns a list of all local branches
func (g *Executor) ListBranches() ([]string, error) {
	cmd := exec.Command("git", "branch", "--format=%(refname:short)")
	cmd.Dir = g.workDir
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	branches := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(branches) == 1 && branches[0] == "" {
		return []string{}, nil
	}
	return branches, nil
}

// GetBranchCommit returns the commit hash of the given branch
func (g *Executor) GetBranchCommit(branch string) (string, error) {
	cmd := exec.Command("git", "rev-parse", branch)
	cmd.Dir = g.workDir
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// GetPRMergeCommit returns the merge commit hash of the PR for the given branch
func (g *Executor) GetPRMergeCommit(branch string) (string, error) {
	// First, try to find the merge commit using git log
	cmd := exec.Command("git", "log", "--merges", "--first-parent", "--grep", fmt.Sprintf("Merge pull request.*%s", branch), "--format=%H")
	cmd.Dir = g.workDir
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	
	commits := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(commits) == 0 || (len(commits) == 1 && commits[0] == "") {
		return "", fmt.Errorf("no merge commit found")
	}
	
	// Return the most recent merge commit
	return commits[0], nil
} 
