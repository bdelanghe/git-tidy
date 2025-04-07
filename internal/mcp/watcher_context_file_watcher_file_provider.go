package mcp

import (
	"fmt"
	"git-tidy/internal/watcher"
	"time"
)

// WatcherContextFileWatcherFileProvider provides context from the watcher's context file watcher file integration
type WatcherContextFileWatcherFileProvider struct {
	workDir string
	testRunner *watcher.TestRunner
}

// NewWatcherContextFileWatcherFileProvider creates a new watcher context file watcher file provider
func NewWatcherContextFileWatcherFileProvider(workDir string) *WatcherContextFileWatcherFileProvider {
	return &WatcherContextFileWatcherFileProvider{
		workDir: workDir,
		testRunner: watcher.New(workDir),
	}
}

// GetContext implements the ContextProvider interface
func (p *WatcherContextFileWatcherFileProvider) GetContext() (Context, error) {
	// Create context
	context := Context{
		Custom: make(map[string]interface{}),
	}

	// Add custom fields
	context.Custom["watcher_context_file_watcher_file_work_dir"] = p.workDir
	context.Custom["watcher_context_file_watcher_file_timestamp"] = time.Now()

	// Add context file watcher file information
	context.Custom["context_file_watcher_file"] = map[string]interface{}{
		"is_running": p.testRunner.IsContextFileWatcherFileRunning(),
		"work_dir": p.testRunner.GetWorkDir(),
	}

	// Add context file watcher file results
	context.Custom["context_file_watcher_file_results"] = p.testRunner.GetContextFileWatcherFileResults()

	return context, nil
} 
