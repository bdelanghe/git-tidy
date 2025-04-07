package mcp

import (
	"fmt"
	"time"
)

// WatcherProvider provides context from the watcher
type WatcherProvider struct {
	workDir string
	startTime time.Time
}

// NewWatcherProvider creates a new watcher provider
func NewWatcherProvider(workDir string) *WatcherProvider {
	return &WatcherProvider{
		workDir: workDir,
		startTime: time.Now(),
	}
}

// GetContext implements the ContextProvider interface
func (p *WatcherProvider) GetContext() (Context, error) {
	// Create context
	context := Context{
		Custom: make(map[string]interface{}),
	}

	// Add custom fields
	context.Custom["watcher_start_time"] = p.startTime
	context.Custom["watcher_uptime"] = time.Since(p.startTime).String()
	context.Custom["watcher_work_dir"] = p.workDir

	return context, nil
} 
