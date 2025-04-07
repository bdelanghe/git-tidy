package mcp

import (
	"fmt"
	"git-tidy/internal/watcher"
	"time"
)

// WatcherContextProvider provides context from the watcher's context integration
type WatcherContextProvider struct {
	workDir string
	testRunner *watcher.TestRunner
}

// NewWatcherContextProvider creates a new watcher context provider
func NewWatcherContextProvider(workDir string) *WatcherContextProvider {
	return &WatcherContextProvider{
		workDir: workDir,
		testRunner: watcher.New(workDir),
	}
}

// GetContext implements the ContextProvider interface
func (p *WatcherContextProvider) GetContext() (Context, error) {
	// Create context
	context := Context{
		Custom: make(map[string]interface{}),
	}

	// Add custom fields
	context.Custom["watcher_context_work_dir"] = p.workDir
	context.Custom["watcher_context_timestamp"] = time.Now()

	// Add context information
	context.Custom["context"] = map[string]interface{}{
		"is_running": p.testRunner.IsContextRunning(),
		"work_dir": p.testRunner.GetWorkDir(),
	}

	// Add context results
	context.Custom["context_results"] = p.testRunner.GetContextResults()

	return context, nil
} 
