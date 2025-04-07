package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// PRData holds the PR's merge commit data
type PRData struct {
	MergeCommit struct {
		Oid string `json:"oid"`
	} `json:"mergeCommit"`
}

func listBranches() ([]string, error) {
	output, err := exec.Command("git", "branch", "--format=%(refname:short)").Output()
	if err != nil {
		return nil, err
	}
	var branches []string
	for _, line := range strings.Split(string(output), "\n") {
		branch := strings.TrimSpace(line)
		if branch != "" {
			branches = append(branches, branch)
		}
	}
	return branches, nil
}

func getLocalCommit(branch string) (string, error) {
	output, err := exec.Command("git", "rev-parse", branch).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func getPRMergeCommit(branch string) (string, error) {
	output, err := exec.Command("gh", "pr", "view", "--head", branch, "--json", "mergeCommit").Output()
	if err != nil {
		return "", err
	}
	var prData PRData
	if err := json.Unmarshal(output, &prData); err != nil {
		return "", err
	}
	return strings.TrimSpace(prData.MergeCommit.Oid), nil
}

func compareCommits(branch, local, merged string) {
	if merged == local {
		fmt.Printf("Branch '%s' matches its merged PR commit: %s\n", branch, merged)
	} else {
		fmt.Printf("Branch '%s' has a merged PR, but the commits differ.\n", branch)
		fmt.Printf("  Local Commit: %s\n", local)
		fmt.Printf("  Merged Commit: %s\n", merged)
	}
}

func main() {
	branches, err := listBranches()
	if err != nil {
		fmt.Printf("Error listing branches: %v\n", err)
		return
	}

	for _, branch := range branches {
		localCommit, err := getLocalCommit(branch)
		if err != nil {
			fmt.Printf("Error getting local commit for branch '%s': %v\n", branch, err)
			continue
		}

		mergeCommit, err := getPRMergeCommit(branch)
		if err != nil {
			fmt.Printf("No merged PR found for branch '%s'\n", branch)
			continue
		}

		compareCommits(branch, localCommit, mergeCommit)
	}
}
