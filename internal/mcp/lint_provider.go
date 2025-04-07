package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// LintProvider provides context from lint errors
type LintProvider struct {
	workDir string
}

// NewLintProvider creates a new lint provider
func NewLintProvider(workDir string) *LintProvider {
	return &LintProvider{
		workDir: workDir,
	}
}

// GetContext implements the ContextProvider interface
func (p *LintProvider) GetContext() (Context, error) {
	// Create context
	context := Context{
		Custom: make(map[string]interface{}),
	}

	// Run golangci-lint
	cmd := exec.Command("golangci-lint", "run", "--out-format=json")
	cmd.Dir = p.workDir
	output, err := cmd.Output()
	if err != nil {
		// If golangci-lint is not installed, return empty context
		if strings.Contains(err.Error(), "executable file not found") {
			return context, nil
		}
		return context, fmt.Errorf("error running golangci-lint: %v", err)
	}

	// Parse the output
	var issues []LintErrorContext
	if err := json.Unmarshal(output, &issues); err != nil {
		return context, fmt.Errorf("error parsing golangci-lint output: %v", err)
	}

	// Add the issues to the context
	context.LintErrors = issues

	// Add a custom field with the total number of issues
	context.Custom["lint_issues_count"] = len(issues)

	return context, nil
} 
