package mcp

import (
	"fmt"
	"git-tidy/internal/watcher"
	"time"
)

// WatcherContextFileWatcherFileContentWatcherProvider provides context from the watcher's context file watcher file content watcher integration
type WatcherContextFileWatcherFileContentWatcherProvider struct {
	workDir string
	testRunner *watcher.TestRunner
}

// NewWatcherContextFileWatcherFileContentWatcherProvider creates a new watcher context file watcher file content watcher provider
func NewWatcherContextFileWatcherFileContentWatcherProvider(workDir string) *WatcherContextFileWatcherFileContentWatcherProvider {
	return &WatcherContextFileWatcherFileContentWatcherProvider{
		workDir: workDir,
		testRunner: watcher.New(workDir),
	}
}

// GetContext implements the ContextProvider interface
func (p *WatcherContextFileWatcherFileContentWatcherProvider) GetContext() (Context, error) {
	// Create context
	context := Context{
		Custom: make(map[string]interface{}),
	}

	// Add custom fields
	context.Custom["watcher_context_file_watcher_file_content_watcher_work_dir"] = p.workDir
	context.Custom["watcher_context_file_watcher_file_content_watcher_timestamp"] = time.Now()

	// Add context file watcher file content watcher information
	context.Custom["context_file_watcher_file_content_watcher"] = map[string]interface{}{
		"is_running": p.testRunner.IsContextFileWatcherFileContentWatcherRunning(),
		"work_dir": p.testRunner.GetWorkDir(),
	}

	// Add context file watcher file content watcher results
	context.Custom["context_file_watcher_file_content_watcher_results"] = p.testRunner.GetContextFileWatcherFileContentWatcherResults()

	return context, nil
} 
