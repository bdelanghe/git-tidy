# Clean Desk

A tool to check if your Git branches are fully integrated with their corresponding merged pull requests.

## Purpose

This tool helps you identify branches that have been merged via pull requests but may have diverged from their merged state. It's useful for keeping your local repository clean and ensuring all branches are in sync with their merged states.

## Prerequisites

- Go 1.16 or higher
- Git
- GitHub CLI (`gh`) installed and authenticated

## Installation

1. Clone this repository
2. Run `go mod tidy` to download dependencies
3. Build the project: `go build`

## Usage

Simply run the compiled binary:

```bash
./clean-desk
```

The tool will:
1. List all local branches
2. For each branch, check if it has a merged PR
3. Compare the local branch commit with the merged PR commit
4. Report any discrepancies

## Output

The tool will output one of the following for each branch:
- If the branch has no merged PR: "No merged PR found for branch 'branch-name'"
- If the branch matches its merged PR: "Branch 'branch-name' matches its merged PR commit: commit-hash"
- If the branch differs from its merged PR: It will show both commit hashes for comparison

## License

MIT
