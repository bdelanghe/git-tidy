package mcp

import (
	"fmt"
	"git-tidy/internal/watcher"
	"time"
)

// WatcherContextFileProvider provides context from the watcher's context file integration
type WatcherContextFileProvider struct {
	workDir string
	testRunner *watcher.TestRunner
}

// NewWatcherContextFileProvider creates a new watcher context file provider
func NewWatcherContextFileProvider(workDir string) *WatcherContextFileProvider {
	return &WatcherContextFileProvider{
		workDir: workDir,
		testRunner: watcher.New(workDir),
	}
}

// GetContext implements the ContextProvider interface
func (p *WatcherContextFileProvider) GetContext() (Context, error) {
	// Create context
	context := Context{
		Custom: make(map[string]interface{}),
	}

	// Add custom fields
	context.Custom["watcher_context_file_work_dir"] = p.workDir
	context.Custom["watcher_context_file_timestamp"] = time.Now()

	// Add context file information
	context.Custom["context_file"] = map[string]interface{}{
		"is_running": p.testRunner.IsContextFileRunning(),
		"work_dir": p.testRunner.GetWorkDir(),
	}

	// Add context file results
	context.Custom["context_file_results"] = p.testRunner.GetContextFileResults()

	return context, nil
} 
