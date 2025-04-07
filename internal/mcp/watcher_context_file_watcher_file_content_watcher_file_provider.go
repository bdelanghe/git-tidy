package mcp

import (
	"fmt"
	"git-tidy/internal/watcher"
	"time"
)

// WatcherContextFileWatcherFileContentWatcherFileProvider provides context from the watcher's context file watcher file content watcher file integration
type WatcherContextFileWatcherFileContentWatcherFileProvider struct {
	workDir string
	testRunner *watcher.TestRunner
}

// NewWatcherContextFileWatcherFileContentWatcherFileProvider creates a new watcher context file watcher file content watcher file provider
func NewWatcherContextFileWatcherFileContentWatcherFileProvider(workDir string) *WatcherContextFileWatcherFileContentWatcherFileProvider {
	return &WatcherContextFileWatcherFileContentWatcherFileProvider{
		workDir: workDir,
		testRunner: watcher.New(workDir),
	}
}

// GetContext implements the ContextProvider interface
func (p *WatcherContextFileWatcherFileContentWatcherFileProvider) GetContext() (Context, error) {
	// Create context
	context := Context{
		Custom: make(map[string]interface{}),
	}

	// Add custom fields
	context.Custom["watcher_context_file_watcher_file_content_watcher_file_work_dir"] = p.workDir
	context.Custom["watcher_context_file_watcher_file_content_watcher_file_timestamp"] = time.Now()

	// Add context file watcher file content watcher file information
	context.Custom["context_file_watcher_file_content_watcher_file"] = map[string]interface{}{
		"is_running": p.testRunner.IsContextFileWatcherFileContentWatcherFileRunning(),
		"work_dir": p.testRunner.GetWorkDir(),
	}

	// Add context file watcher file content watcher file results
	context.Custom["context_file_watcher_file_content_watcher_file_results"] = p.testRunner.GetContextFileWatcherFileContentWatcherFileResults()

	return context, nil
} 
