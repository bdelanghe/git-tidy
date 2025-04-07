// Package git provides functionality for interacting with Git repositories
package git

// RepositoryReader provides methods for reading repository information
type RepositoryReader interface {
	// IsGitRepository checks if the current directory is a Git repository
	IsGitRepository() bool
	// GetRepositoryRoot returns the root directory of the Git repository
	GetRepositoryRoot() (string, error)
}

// BranchReader provides methods for reading branch information
type BranchReader interface {
	// GetCurrentBranch returns the name of the current branch
	GetCurrentBranch() (string, error)
	// ListBranches returns a list of all local branches
	ListBranches() ([]string, error)
	// GetBranchCommit returns the commit hash of the given branch
	GetBranchCommit(branch string) (string, error)
}

// PRReader provides methods for reading pull request information
type PRReader interface {
	// GetPRMergeCommit returns the merge commit hash of the PR for the given branch
	GetPRMergeCommit(branch string) (string, error)
}

// GitOperator combines all Git operations
type GitOperator interface {
	RepositoryReader
	BranchReader
	PRReader
} 
