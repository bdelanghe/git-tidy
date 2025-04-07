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
func NewBaseProvider(workDir string) (BaseProvider, error) {
	runner, err := watcher.New(workDir)
	if err != nil {
		return BaseProvider{}, err
	}
	return BaseProvider{
		workDir:    workDir,
		testRunner: runner,
	}, nil
}

// FileSystemProvider provides file system context
type FileSystemProvider struct {
	BaseProvider
}

// NewFileSystemProvider creates a new file system provider
func NewFileSystemProvider(workDir string) (*FileSystemProvider, error) {
	base, err := NewBaseProvider(workDir)
	if err != nil {
		return nil, err
	}
	return &FileSystemProvider{
		BaseProvider: base,
	}, nil
}

// GetContext implements the ContextProvider interface
func (p *FileSystemProvider) GetContext() (Context, error) {
	return Context{
		Timestamp: time.Now(),
		Custom: map[string]interface{}{
			"fs_provider": map[string]interface{}{
				"work_dir": p.workDir,
				"status":   "active",
			},
		},
	}, nil
}

// WatcherProvider provides watcher context
type WatcherProvider struct {
	BaseProvider
}

// NewWatcherProvider creates a new watcher provider
func NewWatcherProvider(workDir string) (*WatcherProvider, error) {
	base, err := NewBaseProvider(workDir)
	if err != nil {
		return nil, err
	}
	return &WatcherProvider{
		BaseProvider: base,
	}, nil
}

// GetContext implements the ContextProvider interface
func (p *WatcherProvider) GetContext() (Context, error) {
	return Context{
		Timestamp: time.Now(),
		Custom: map[string]interface{}{
			"watcher": map[string]interface{}{
				"work_dir": p.workDir,
				"status":   "active",
			},
		},
	}, nil
}

// ContentProvider provides file content context
type ContentProvider struct {
	BaseProvider
}

// NewContentProvider creates a new content provider
func NewContentProvider(workDir string) (*ContentProvider, error) {
	base, err := NewBaseProvider(workDir)
	if err != nil {
		return nil, err
	}
	return &ContentProvider{
		BaseProvider: base,
	}, nil
}

// GetContext implements the ContextProvider interface
func (p *ContentProvider) GetContext() (Context, error) {
	return Context{
		Timestamp: time.Now(),
		Custom: map[string]interface{}{
			"content": map[string]interface{}{
				"work_dir": p.workDir,
				"status":   "active",
			},
		},
	}, nil
}

// AnalyzerProvider provides code analysis context
type AnalyzerProvider struct {
	BaseProvider
}

// NewAnalyzerProvider creates a new analyzer provider
func NewAnalyzerProvider(workDir string) (*AnalyzerProvider, error) {
	base, err := NewBaseProvider(workDir)
	if err != nil {
		return nil, err
	}
	return &AnalyzerProvider{
		BaseProvider: base,
	}, nil
}

// GetContext implements the ContextProvider interface
func (p *AnalyzerProvider) GetContext() (Context, error) {
	return Context{
		Timestamp: time.Now(),
		Custom: map[string]interface{}{
			"analyzer": map[string]interface{}{
				"work_dir": p.workDir,
				"status":   "active",
			},
		},
	}, nil
} 
