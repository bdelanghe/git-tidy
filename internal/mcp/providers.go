package mcp

import (
	"git-tidy/internal/watcher"
	"time"
)

// BaseProvider provides common functionality for all providers
type BaseProvider struct {
	workDir    string
	testRunner *watcher.TestRunner
}

// NewBaseProvider creates a new base provider
func NewBaseProvider(workDir string) BaseProvider {
	return BaseProvider{
		workDir:    workDir,
		testRunner: watcher.New(workDir),
	}
}

// FileSystemProvider provides file system context
type FileSystemProvider struct {
	BaseProvider
}

// NewFileSystemProvider creates a new file system provider
func NewFileSystemProvider(workDir string) *FileSystemProvider {
	return &FileSystemProvider{
		BaseProvider: NewBaseProvider(workDir),
	}
}

// GetContext implements the ContextProvider interface
func (p *FileSystemProvider) GetContext() (Context, error) {
	return Context{
		Timestamp: time.Now(),
		Custom: map[string]interface{}{
			"fs_provider": map[string]interface{}{
				"work_dir":   p.workDir,
				"is_running": p.testRunner.IsRunning(),
			},
		},
	}, nil
}

// WatcherProvider provides watcher context
type WatcherProvider struct {
	BaseProvider
}

// NewWatcherProvider creates a new watcher provider
func NewWatcherProvider(workDir string) *WatcherProvider {
	return &WatcherProvider{
		BaseProvider: NewBaseProvider(workDir),
	}
}

// GetContext implements the ContextProvider interface
func (p *WatcherProvider) GetContext() (Context, error) {
	return Context{
		Timestamp: time.Now(),
		Custom: map[string]interface{}{
			"watcher": map[string]interface{}{
				"work_dir":   p.workDir,
				"is_running": p.testRunner.IsRunning(),
				"test_cache": p.testRunner.GetTestCache(),
			},
		},
	}, nil
}

// ContentProvider provides file content context
type ContentProvider struct {
	BaseProvider
}

// NewContentProvider creates a new content provider
func NewContentProvider(workDir string) *ContentProvider {
	return &ContentProvider{
		BaseProvider: NewBaseProvider(workDir),
	}
}

// GetContext implements the ContextProvider interface
func (p *ContentProvider) GetContext() (Context, error) {
	return Context{
		Timestamp: time.Now(),
		Custom: map[string]interface{}{
			"content": map[string]interface{}{
				"work_dir":   p.workDir,
				"is_running": p.testRunner.IsRunning(),
			},
		},
	}, nil
}

// AnalyzerProvider provides code analysis context
type AnalyzerProvider struct {
	BaseProvider
}

// NewAnalyzerProvider creates a new analyzer provider
func NewAnalyzerProvider(workDir string) *AnalyzerProvider {
	return &AnalyzerProvider{
		BaseProvider: NewBaseProvider(workDir),
	}
}

// GetContext implements the ContextProvider interface
func (p *AnalyzerProvider) GetContext() (Context, error) {
	return Context{
		Timestamp: time.Now(),
		Custom: map[string]interface{}{
			"analyzer": map[string]interface{}{
				"work_dir":   p.workDir,
				"is_running": p.testRunner.IsAnalyzerRunning(),
				"results":    p.testRunner.GetAnalysisResults(),
			},
		},
	}, nil
} 
