package mcp

import (
	"fmt"
	"git-tidy/internal/watcher"
	"time"
)

// WatcherTestProvider provides context from the watcher's test integration
type WatcherTestProvider struct {
	workDir string
	testRunner *watcher.TestRunner
}

// NewWatcherTestProvider creates a new watcher test provider
func NewWatcherTestProvider(workDir string) *WatcherTestProvider {
	return &WatcherTestProvider{
		workDir: workDir,
		testRunner: watcher.New(workDir),
	}
}

// GetContext implements the ContextProvider interface
func (p *WatcherTestProvider) GetContext() (Context, error) {
	// Create context
	context := Context{
		Custom: make(map[string]interface{}),
	}

	// Add custom fields
	context.Custom["watcher_test_work_dir"] = p.workDir
	context.Custom["watcher_test_timestamp"] = time.Now()

	// Add test information
	context.Custom["test"] = map[string]interface{}{
		"is_running": p.testRunner.IsTestRunning(),
		"work_dir": p.testRunner.GetWorkDir(),
	}

	// Add test results
	context.Custom["test_results"] = p.testRunner.GetTestResults()

	return context, nil
} 
