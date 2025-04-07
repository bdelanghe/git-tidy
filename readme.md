# git-tidy

A Git extension to check if your branches are fully integrated with their corresponding merged pull requests.

## Purpose

This tool helps you identify branches that have been merged via pull requests but may have diverged from their merged state. It's useful for keeping your local repository clean and ensuring all branches are in sync with their merged states.

## Prerequisites

- Go 1.16 or higher
- Git
- GitHub CLI (`gh`) installed and authenticated

## Installation

There are two ways to install `git-tidy`:

### 1. Using `go install` (Recommended)

```bash
go install github.com/bdelanghe/git-tidy@latest
```

This will automatically install the binary in your `$GOPATH/bin` directory. Make sure your `$GOPATH/bin` is in your PATH.

### 2. Building from source

1. Clone this repository:
   ```bash
   git clone https://github.com/bdelanghe/git-tidy.git
   cd git-tidy
   ```

2. Build and install:
   ```bash
   # Option A: Install to $GOPATH/bin
   go install

   # Option B: Install to /usr/local/bin (requires sudo)
   sudo GOBIN=/usr/local/bin go install
   ```

## Usage

Once installed, you can use it either way:

```bash
git tidy    # As a Git subcommand
# or
git-tidy    # As a standalone command
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

## Troubleshooting

1. If `git tidy` is not recognized:
   - Make sure the binary is in your PATH
   - Try using `git-tidy` directly
   - Check installation with `which git-tidy`

2. If GitHub CLI errors occur:
   - Ensure `gh` is installed: `gh --version`
   - Make sure you're authenticated: `gh auth status`
   - Run `gh auth login` if needed

## License

MIT
