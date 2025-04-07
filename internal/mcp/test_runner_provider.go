package mcp

import (
	"fmt"
	"git-tidy/internal/watcher"
	"time"
)

// TestRunnerProvider provides context from the test runner
type TestRunnerProvider struct {
	workDir string
	testRunner *watcher.TestRunner
}

// NewTestRunnerProvider creates a new test runner provider
func NewTestRunnerProvider(workDir string) *TestRunnerProvider {
	return &TestRunnerProvider{
		workDir: workDir,
		testRunner: watcher.New(workDir),
	}
}

// GetContext implements the ContextProvider interface
func (p *TestRunnerProvider) GetContext() (Context, error) {
	// Create context
	context := Context{
		Custom: make(map[string]interface{}),
	}

	// Add custom fields
	context.Custom["test_runner_work_dir"] = p.workDir
	context.Custom["test_runner_timestamp"] = time.Now()

	// Add test cache information
	context.Custom["test_cache"] = p.testRunner.GetTestCache()

	// Add watcher information
	context.Custom["watcher"] = map[string]interface{}{
		"is_running": p.testRunner.IsRunning(),
		"work_dir": p.testRunner.GetWorkDir(),
	}

	return context, nil
} 
