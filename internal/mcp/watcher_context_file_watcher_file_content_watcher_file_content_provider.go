package mcp

import (
	"fmt"
	"git-tidy/internal/watcher"
	"time"
)

// WatcherContextFileWatcherFileContentWatcherFileContentProvider provides context from the watcher's context file watcher file content watcher file content integration
type WatcherContextFileWatcherFileContentWatcherFileContentProvider struct {
	workDir string
	testRunner *watcher.TestRunner
}

// NewWatcherContextFileWatcherFileContentWatcherFileContentProvider creates a new watcher context file watcher file content watcher file content provider
func NewWatcherContextFileWatcherFileContentWatcherFileContentProvider(workDir string) *WatcherContextFileWatcherFileContentWatcherFileContentProvider {
	return &WatcherContextFileWatcherFileContentWatcherFileContentProvider{
		workDir: workDir,
		testRunner: watcher.New(workDir),
	}
}

// GetContext implements the ContextProvider interface
func (p *WatcherContextFileWatcherFileContentWatcherFileContentProvider) GetContext() (Context, error) {
	// Create context
	context := Context{
		Custom: make(map[string]interface{}),
	}

	// Add custom fields
	context.Custom["watcher_context_file_watcher_file_content_watcher_file_content_work_dir"] = p.workDir
	context.Custom["watcher_context_file_watcher_file_content_watcher_file_content_timestamp"] = time.Now()

	// Add context file watcher file content watcher file content information
	context.Custom["context_file_watcher_file_content_watcher_file_content"] = map[string]interface{}{
		"is_running": p.testRunner.IsContextFileWatcherFileContentWatcherFileContentRunning(),
		"work_dir": p.testRunner.GetWorkDir(),
	}

	// Add context file watcher file content watcher file content results
	context.Custom["context_file_watcher_file_content_watcher_file_content_results"] = p.testRunner.GetContextFileWatcherFileContentWatcherFileContentResults()

	return context, nil
} 
