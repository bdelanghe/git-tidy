package git

// Interface defines the contract for Git operations
type Interface interface {
	// Branch operations
	ListBranches() ([]string, error)
	GetCurrentBranch() (string, error)
	GetBranchCommit(branch string) (string, error)
	
	// PR operations
	GetPRMergeCommit(branch string) (string, error)
	
	// Repository operations
	IsGitRepository() bool
	GetRepositoryRoot() (string, error)
}

// Command represents a Git command execution
type Command struct {
	Args []string
	Dir  string
}

// Result represents the output of a Git command
type Result struct {
	Stdout string
	Stderr string
	Error  error
} 
