package mcp

import (
	"fmt"
	"git-tidy/internal/watcher"
	"time"
)

// WatcherAnalyzerProvider provides context from the watcher's analyzer
type WatcherAnalyzerProvider struct {
	workDir string
	testRunner *watcher.TestRunner
}

// NewWatcherAnalyzerProvider creates a new watcher analyzer provider
func NewWatcherAnalyzerProvider(workDir string) *WatcherAnalyzerProvider {
	return &WatcherAnalyzerProvider{
		workDir: workDir,
		testRunner: watcher.New(workDir),
	}
}

// GetContext implements the ContextProvider interface
func (p *WatcherAnalyzerProvider) GetContext() (Context, error) {
	// Create context
	context := Context{
		Custom: make(map[string]interface{}),
	}

	// Add custom fields
	context.Custom["watcher_analyzer_work_dir"] = p.workDir
	context.Custom["watcher_analyzer_timestamp"] = time.Now()

	// Add analyzer information
	context.Custom["analyzer"] = map[string]interface{}{
		"is_running": p.testRunner.IsAnalyzerRunning(),
		"work_dir": p.testRunner.GetWorkDir(),
	}

	// Add analysis results
	context.Custom["analysis_results"] = p.testRunner.GetAnalysisResults()

	return context, nil
} 
