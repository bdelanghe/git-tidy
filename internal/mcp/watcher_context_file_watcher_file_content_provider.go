package mcp

import (
	"fmt"
	"git-tidy/internal/watcher"
	"time"
)

// WatcherContextFileWatcherFileContentProvider provides context from the watcher's context file watcher file content integration
type WatcherContextFileWatcherFileContentProvider struct {
	workDir string
	testRunner *watcher.TestRunner
}

// NewWatcherContextFileWatcherFileContentProvider creates a new watcher context file watcher file content provider
func NewWatcherContextFileWatcherFileContentProvider(workDir string) *WatcherContextFileWatcherFileContentProvider {
	return &WatcherContextFileWatcherFileContentProvider{
		workDir: workDir,
		testRunner: watcher.New(workDir),
	}
}

// GetContext implements the ContextProvider interface
func (p *WatcherContextFileWatcherFileContentProvider) GetContext() (Context, error) {
	// Create context
	context := Context{
		Custom: make(map[string]interface{}),
	}

	// Add custom fields
	context.Custom["watcher_context_file_watcher_file_content_work_dir"] = p.workDir
	context.Custom["watcher_context_file_watcher_file_content_timestamp"] = time.Now()

	// Add context file watcher file content information
	context.Custom["context_file_watcher_file_content"] = map[string]interface{}{
		"is_running": p.testRunner.IsContextFileWatcherFileContentRunning(),
		"work_dir": p.testRunner.GetWorkDir(),
	}

	// Add context file watcher file content results
	context.Custom["context_file_watcher_file_content_results"] = p.testRunner.GetContextFileWatcherFileContentResults()

	return context, nil
} 
