package mcp

import (
	"fmt"
	"git-tidy/internal/watcher"
	"time"
)

// WatcherFSProvider provides context from the watcher's file system integration
type WatcherFSProvider struct {
	workDir string
	testRunner *watcher.TestRunner
}

// NewWatcherFSProvider creates a new watcher file system provider
func NewWatcherFSProvider(workDir string) *WatcherFSProvider {
	return &WatcherFSProvider{
		workDir: workDir,
		testRunner: watcher.New(workDir),
	}
}

// GetContext implements the ContextProvider interface
func (p *WatcherFSProvider) GetContext() (Context, error) {
	// Create context
	context := Context{
		Custom: make(map[string]interface{}),
	}

	// Add custom fields
	context.Custom["watcher_fs_work_dir"] = p.workDir
	context.Custom["watcher_fs_timestamp"] = time.Now()

	// Add file system information
	context.Custom["fs"] = map[string]interface{}{
		"is_running": p.testRunner.IsFSRunning(),
		"work_dir": p.testRunner.GetWorkDir(),
	}

	// Add file system results
	context.Custom["fs_results"] = p.testRunner.GetFSResults()

	return context, nil
} 
