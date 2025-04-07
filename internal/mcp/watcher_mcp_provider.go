package mcp

import (
	"fmt"
	"git-tidy/internal/watcher"
	"time"
)

// WatcherMCPProvider provides context from the watcher's MCP integration
type WatcherMCPProvider struct {
	workDir string
	testRunner *watcher.TestRunner
}

// NewWatcherMCPProvider creates a new watcher MCP provider
func NewWatcherMCPProvider(workDir string) *WatcherMCPProvider {
	return &WatcherMCPProvider{
		workDir: workDir,
		testRunner: watcher.New(workDir),
	}
}

// GetContext implements the ContextProvider interface
func (p *WatcherMCPProvider) GetContext() (Context, error) {
	// Create context
	context := Context{
		Custom: make(map[string]interface{}),
	}

	// Add custom fields
	context.Custom["watcher_mcp_work_dir"] = p.workDir
	context.Custom["watcher_mcp_timestamp"] = time.Now()

	// Add MCP information
	context.Custom["mcp"] = map[string]interface{}{
		"is_running": p.testRunner.IsMCPRunning(),
		"work_dir": p.testRunner.GetWorkDir(),
	}

	// Add MCP results
	context.Custom["mcp_results"] = p.testRunner.GetMCPResults()

	return context, nil
} 
