package mcp

import (
	"fmt"
	"git-tidy/internal/watcher"
	"time"
)

// WatcherLintProvider provides context from the watcher's lint integration
type WatcherLintProvider struct {
	workDir string
	testRunner *watcher.TestRunner
}

// NewWatcherLintProvider creates a new watcher lint provider
func NewWatcherLintProvider(workDir string) *WatcherLintProvider {
	return &WatcherLintProvider{
		workDir: workDir,
		testRunner: watcher.New(workDir),
	}
}

// GetContext implements the ContextProvider interface
func (p *WatcherLintProvider) GetContext() (Context, error) {
	// Create context
	context := Context{
		Custom: make(map[string]interface{}),
	}

	// Add custom fields
	context.Custom["watcher_lint_work_dir"] = p.workDir
	context.Custom["watcher_lint_timestamp"] = time.Now()

	// Add lint information
	context.Custom["lint"] = map[string]interface{}{
		"is_running": p.testRunner.IsLintRunning(),
		"work_dir": p.testRunner.GetWorkDir(),
	}

	// Add lint results
	context.Custom["lint_results"] = p.testRunner.GetLintResults()

	return context, nil
} 
