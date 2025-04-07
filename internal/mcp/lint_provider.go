package mcp

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// LintProvider provides context from lint errors
type LintProvider struct {
	BaseProvider
}

// NewLintProvider creates a new lint provider
func NewLintProvider(workDir string) (*LintProvider, error) {
	base, err := NewBaseProvider(workDir)
	if err != nil {
		return nil, err
	}
	return &LintProvider{
		BaseProvider: base,
	}, nil
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
	context.Lint = issues

	// Add a custom field with the total number of issues
	context.Custom["lint_issues_count"] = len(issues)

	return context, nil
} 
