// Package main provides the git-tidy command
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bdelanghe/git-tidy/git"
)

func main() {
	// Get the current working directory
	workDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting working directory: %v\n", err)
		os.Exit(1)
	}

	// Create a new Git executor
	g := git.New(workDir)

	// Check if we're in a Git repository
	if !g.IsGitRepository() {
		fmt.Println("Error: Not a Git repository")
		os.Exit(1)
	}

	// Get repository root
	repoRoot, err := g.GetRepositoryRoot()
	if err != nil {
		fmt.Printf("Error getting repository root: %v\n", err)
		os.Exit(1)
	}

	// Create a new executor with the repository root
	g = git.New(repoRoot)

	// Get all branches
	branches, err := g.ListBranches()
	if err != nil {
		fmt.Printf("Error listing branches: %v\n", err)
		os.Exit(1)
	}

	// Check each branch
	for _, branch := range branches {
		// Skip the current branch
		currentBranch, err := g.GetCurrentBranch()
		if err != nil {
			fmt.Printf("Error getting current branch: %v\n", err)
			continue
		}
		if branch == currentBranch {
			continue
		}

		// Get the local commit
		localCommit, err := g.GetBranchCommit(branch)
		if err != nil {
			fmt.Printf("Error getting local commit for branch '%s': %v\n", branch, err)
			continue
		}

		// Get the PR merge commit
		mergeCommit, err := g.GetPRMergeCommit(branch)
		if err != nil {
			fmt.Printf("No merged PR found for branch '%s'\n", branch)
			continue
		}

		// Compare the commits
		if mergeCommit == localCommit {
			fmt.Printf("Branch '%s' matches its merged PR commit: %s\n", branch, mergeCommit)
		} else {
			fmt.Printf("Branch '%s' has a merged PR, but the commits differ.\n", branch)
			fmt.Printf("  Local Commit: %s\n", localCommit)
			fmt.Printf("  Merged Commit: %s\n", mergeCommit)
		}
	}
}
