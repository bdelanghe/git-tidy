package mcp

import (
	"fmt"
	"git-tidy/internal/git"
	"time"
)

// GitProvider provides context from the Git repository
type GitProvider struct {
	workDir string
	gitExecutor *git.Executor
}

// NewGitProvider creates a new Git provider
func NewGitProvider(workDir string) *GitProvider {
	return &GitProvider{
		workDir: workDir,
		gitExecutor: git.NewExecutor(workDir),
	}
}

// GetContext implements the ContextProvider interface
func (p *GitProvider) GetContext() (Context, error) {
	// Create context
	context := Context{
		Custom: make(map[string]interface{}),
	}

	// Check if the directory is a Git repository
	if !p.gitExecutor.IsGitRepository() {
		return context, fmt.Errorf("not a Git repository: %s", p.workDir)
	}

	// Get the repository root
	root, err := p.gitExecutor.GetRepositoryRoot()
	if err != nil {
		return context, fmt.Errorf("error getting repository root: %v", err)
	}

	// Get the current branch
	branch, err := p.gitExecutor.GetCurrentBranch()
	if err != nil {
		return context, fmt.Errorf("error getting current branch: %v", err)
	}

	// Get the current commit
	commit, err := p.gitExecutor.GetBranchCommit(branch)
	if err != nil {
		return context, fmt.Errorf("error getting current commit: %v", err)
	}

	// Get the modified files
	modified, err := p.gitExecutor.GetModifiedFiles()
	if err != nil {
		return context, fmt.Errorf("error getting modified files: %v", err)
	}

	// Get the untracked files
	untracked, err := p.gitExecutor.GetUntrackedFiles()
	if err != nil {
		return context, fmt.Errorf("error getting untracked files: %v", err)
	}

	// Get the deleted files
	deleted, err := p.gitExecutor.GetDeletedFiles()
	if err != nil {
		return context, fmt.Errorf("error getting deleted files: %v", err)
	}

	// Get the staged files
	staged, err := p.gitExecutor.GetStagedFiles()
	if err != nil {
		return context, fmt.Errorf("error getting staged files: %v", err)
	}

	// Get the unstaged files
	unstaged, err := p.gitExecutor.GetUnstagedFiles()
	if err != nil {
		return context, fmt.Errorf("error getting unstaged files: %v", err)
	}

	// Get the recent commits
	commits, err := p.gitExecutor.GetRecentCommits(10)
	if err != nil {
		return context, fmt.Errorf("error getting recent commits: %v", err)
	}

	// Add the Git context
	context.Git = GitContext{
		Root:       root,
		Branch:     branch,
		Commit:     commit,
		Modified:   modified,
		Untracked:  untracked,
		Deleted:    deleted,
		Staged:     staged,
		Unstaged:   unstaged,
		Commits:    commits,
		LastUpdate: time.Now(),
	}

	// Add custom fields
	context.Custom["git_repository"] = root
	context.Custom["git_branch"] = branch
	context.Custom["git_commit"] = commit
	context.Custom["git_modified_count"] = len(modified)
	context.Custom["git_untracked_count"] = len(untracked)
	context.Custom["git_deleted_count"] = len(deleted)
	context.Custom["git_staged_count"] = len(staged)
	context.Custom["git_unstaged_count"] = len(unstaged)
	context.Custom["git_commit_count"] = len(commits)

	return context, nil
} 
