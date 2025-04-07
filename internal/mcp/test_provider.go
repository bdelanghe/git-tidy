package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// TestProvider provides context from test results
type TestProvider struct {
	workDir string
}

// NewTestProvider creates a new test provider
func NewTestProvider(workDir string) *TestProvider {
	return &TestProvider{
		workDir: workDir,
	}
}

// GetContext implements the ContextProvider interface
func (p *TestProvider) GetContext() (Context, error) {
	// Create context
	context := Context{
		Custom: make(map[string]interface{}),
	}

	// Run tests with JSON output
	cmd := exec.Command("go", "test", "./...", "-json")
	cmd.Dir = p.workDir
	output, err := cmd.Output()
	if err != nil {
		return context, fmt.Errorf("error running tests: %v", err)
	}

	// Parse the output
	var events []TestEvent
	decoder := json.NewDecoder(strings.NewReader(string(output)))
	for decoder.More() {
		var event TestEvent
		if err := decoder.Decode(&event); err != nil {
			return context, fmt.Errorf("error parsing test output: %v", err)
		}
		events = append(events, event)
	}

	// Add the events to the context
	context.Custom["test_events"] = events

	// Add a custom field with the total number of tests
	context.Custom["test_count"] = len(events)

	// Add a custom field with the number of passed tests
	passedTests := 0
	for _, event := range events {
		if event.Action == "pass" {
			passedTests++
		}
	}
	context.Custom["passed_test_count"] = passedTests

	// Add a custom field with the number of failed tests
	failedTests := 0
	for _, event := range events {
		if event.Action == "fail" {
			failedTests++
		}
	}
	context.Custom["failed_test_count"] = failedTests

	return context, nil
}

// TestEvent represents a test event from the Go test JSON output
type TestEvent struct {
	Time    string `json:"Time"`
	Action  string `json:"Action"`
	Package string `json:"Package"`
	Test    string `json:"Test,omitempty"`
	Output  string `json:"Output,omitempty"`
	Elapsed string `json:"Elapsed,omitempty"`
} 
