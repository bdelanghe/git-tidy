package git

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// Executor implements the Git interface using os/exec
type Executor struct {
	workDir string
}

// New creates a new Git executor
func New(workDir string) *Executor {
	return &Executor{
		workDir: workDir,
	}
}

// execute runs a Git command and returns the result
func (e *Executor) execute(args ...string) (*Result, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = e.workDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return &Result{
		Stdout: strings.TrimSpace(stdout.String()),
		Stderr: strings.TrimSpace(stderr.String()),
		Error:  err,
	}, err
}

// ListBranches returns all local branches
func (e *Executor) ListBranches() ([]string, error) {
	result, err := e.execute("branch", "--format=%(refname:short)")
	if err != nil {
		return nil, fmt.Errorf("failed to list branches: %w", err)
	}

	var branches []string
	for _, line := range strings.Split(result.Stdout, "\n") {
		branch := strings.TrimSpace(line)
		if branch != "" {
			branches = append(branches, branch)
		}
	}
	return branches, nil
}

// GetCurrentBranch returns the current branch name
func (e *Executor) GetCurrentBranch() (string, error) {
	result, err := e.execute("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}
	return result.Stdout, nil
}

// GetBranchCommit returns the commit hash for a branch
func (e *Executor) GetBranchCommit(branch string) (string, error) {
	result, err := e.execute("rev-parse", branch)
	if err != nil {
		return "", fmt.Errorf("failed to get commit for branch %s: %w", branch, err)
	}
	return result.Stdout, nil
}

// GetPRMergeCommit returns the merge commit hash for a PR
func (e *Executor) GetPRMergeCommit(branch string) (string, error) {
	result, err := e.execute("gh", "pr", "view", "--head", branch, "--json", "mergeCommit")
	if err != nil {
		return "", fmt.Errorf("failed to get PR merge commit for branch %s: %w", branch, err)
	}
	
	// Parse the JSON response
	var prData struct {
		MergeCommit struct {
			Oid string `json:"oid"`
		} `json:"mergeCommit"`
	}
	
	if err := json.Unmarshal([]byte(result.Stdout), &prData); err != nil {
		return "", fmt.Errorf("failed to parse PR data: %w", err)
	}
	
	return prData.MergeCommit.Oid, nil
}

// IsGitRepository checks if the current directory is a Git repository
func (e *Executor) IsGitRepository() bool {
	_, err := e.execute("rev-parse", "--git-dir")
	return err == nil
}

// GetRepositoryRoot returns the root directory of the Git repository
func (e *Executor) GetRepositoryRoot() (string, error) {
	result, err := e.execute("rev-parse", "--show-toplevel")
	if err != nil {
		return "", fmt.Errorf("failed to get repository root: %w", err)
	}
	return result.Stdout, nil
} 
