package mcp

import (
	"fmt"
	"git-tidy/internal/watcher"
	"time"
)

// WatcherGitProvider provides context from the watcher's Git integration
type WatcherGitProvider struct {
	workDir string
	testRunner *watcher.TestRunner
}

// NewWatcherGitProvider creates a new watcher Git provider
func NewWatcherGitProvider(workDir string) *WatcherGitProvider {
	return &WatcherGitProvider{
		workDir: workDir,
		testRunner: watcher.New(workDir),
	}
}

// GetContext implements the ContextProvider interface
func (p *WatcherGitProvider) GetContext() (Context, error) {
	// Create context
	context := Context{
		Custom: make(map[string]interface{}),
	}

	// Add custom fields
	context.Custom["watcher_git_work_dir"] = p.workDir
	context.Custom["watcher_git_timestamp"] = time.Now()

	// Add Git information
	context.Custom["git"] = map[string]interface{}{
		"is_running": p.testRunner.IsGitRunning(),
		"work_dir": p.testRunner.GetWorkDir(),
	}

	// Add Git results
	context.Custom["git_results"] = p.testRunner.GetGitResults()

	return context, nil
} 
