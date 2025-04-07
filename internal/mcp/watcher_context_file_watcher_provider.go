package mcp

import (
	"fmt"
	"git-tidy/internal/watcher"
	"time"
)

// WatcherContextFileWatcherProvider provides context from the watcher's context file watcher integration
type WatcherContextFileWatcherProvider struct {
	workDir string
	testRunner *watcher.TestRunner
}

// NewWatcherContextFileWatcherProvider creates a new watcher context file watcher provider
func NewWatcherContextFileWatcherProvider(workDir string) *WatcherContextFileWatcherProvider {
	return &WatcherContextFileWatcherProvider{
		workDir: workDir,
		testRunner: watcher.New(workDir),
	}
}

// GetContext implements the ContextProvider interface
func (p *WatcherContextFileWatcherProvider) GetContext() (Context, error) {
	// Create context
	context := Context{
		Custom: make(map[string]interface{}),
	}

	// Add custom fields
	context.Custom["watcher_context_file_watcher_work_dir"] = p.workDir
	context.Custom["watcher_context_file_watcher_timestamp"] = time.Now()

	// Add context file watcher information
	context.Custom["context_file_watcher"] = map[string]interface{}{
		"is_running": p.testRunner.IsContextFileWatcherRunning(),
		"work_dir": p.testRunner.GetWorkDir(),
	}

	// Add context file watcher results
	context.Custom["context_file_watcher_results"] = p.testRunner.GetContextFileWatcherResults()

	return context, nil
} 
