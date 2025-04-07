package mcp

import (
	"fmt"
	"git-tidy/internal/watcher"
	"time"
)

// WatcherRulesProvider provides context from the watcher's rules integration
type WatcherRulesProvider struct {
	workDir string
	testRunner *watcher.TestRunner
}

// NewWatcherRulesProvider creates a new watcher rules provider
func NewWatcherRulesProvider(workDir string) *WatcherRulesProvider {
	return &WatcherRulesProvider{
		workDir: workDir,
		testRunner: watcher.New(workDir),
	}
}

// GetContext implements the ContextProvider interface
func (p *WatcherRulesProvider) GetContext() (Context, error) {
	// Create context
	context := Context{
		Custom: make(map[string]interface{}),
	}

	// Add custom fields
	context.Custom["watcher_rules_work_dir"] = p.workDir
	context.Custom["watcher_rules_timestamp"] = time.Now()

	// Add rules information
	context.Custom["rules"] = map[string]interface{}{
		"is_running": p.testRunner.IsRulesRunning(),
		"work_dir": p.testRunner.GetWorkDir(),
	}

	// Add rules results
	context.Custom["rules_results"] = p.testRunner.GetRulesResults()

	return context, nil
} 
