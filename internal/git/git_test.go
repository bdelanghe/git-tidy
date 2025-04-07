package git

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	workDir := "/test/dir"
	g := New(workDir)
	if g.workDir != workDir {
		t.Errorf("Expected workDir to be %s, got %s", workDir, g.workDir)
	}
}

func TestIsGitRepository(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "git-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test with non-git directory
	g := New(tmpDir)
	if g.IsGitRepository() {
		t.Error("Expected IsGitRepository to return false for non-git directory")
	}

	// Initialize git repository
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Test with git directory
	if !g.IsGitRepository() {
		t.Error("Expected IsGitRepository to return true for git directory")
	}
}

func TestGetRepositoryRoot(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "git-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Initialize git repository
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	g := New(tmpDir)
	root, err := g.GetRepositoryRoot()
	if err != nil {
		t.Fatalf("Failed to get repository root: %v", err)
	}

	// The root should be the same as tmpDir
	if root != tmpDir {
		t.Errorf("Expected repository root to be %s, got %s", tmpDir, root)
	}
}

func TestListBranches(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "git-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Initialize git repository
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Create a test file and commit it
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd = exec.Command("git", "add", "test.txt")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to add test file: %v", err)
	}

	cmd = exec.Command("git", "commit", "-m", "Initial commit")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to commit test file: %v", err)
	}

	// Create a test branch
	cmd = exec.Command("git", "checkout", "-b", "test-branch")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create test branch: %v", err)
	}

	g := New(tmpDir)
	branches, err := g.ListBranches()
	if err != nil {
		t.Fatalf("Failed to list branches: %v", err)
	}

	// We should have at least main and test-branch
	if len(branches) < 2 {
		t.Errorf("Expected at least 2 branches, got %d", len(branches))
	}

	// Check if test-branch is in the list
	found := false
	for _, branch := range branches {
		if branch == "test-branch" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected to find test-branch in the list of branches")
	}
} 
