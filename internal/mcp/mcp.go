package mcp

import (
	"fmt"
	"time"
)

// MCPServer represents the Model Context Protocol server
type MCPServer struct {
	server *Server
}

// NewMCPServer creates a new MCP server
func NewMCPServer(workDir string) (*MCPServer, error) {
	server := NewServer(".cursor/context.json")

	// Add providers
	fsProvider, err := NewFileSystemProvider(workDir)
	if err != nil {
		return nil, fmt.Errorf("error creating file system provider: %v", err)
	}
	server.AddProvider(fsProvider)

	watcherProvider, err := NewWatcherProvider(workDir)
	if err != nil {
		return nil, fmt.Errorf("error creating watcher provider: %v", err)
	}
	server.AddProvider(watcherProvider)

	contentProvider, err := NewContentProvider(workDir)
	if err != nil {
		return nil, fmt.Errorf("error creating content provider: %v", err)
	}
	server.AddProvider(contentProvider)

	analyzerProvider, err := NewAnalyzerProvider(workDir)
	if err != nil {
		return nil, fmt.Errorf("error creating analyzer provider: %v", err)
	}
	server.AddProvider(analyzerProvider)

	lintProvider, err := NewLintProvider(workDir)
	if err != nil {
		return nil, fmt.Errorf("error creating lint provider: %v", err)
	}
	server.AddProvider(lintProvider)

	return &MCPServer{
		server: server,
	}, nil
}

// Start starts the MCP server
func (s *MCPServer) Start() error {
	fmt.Printf("Starting MCP server with output path: %s\n", s.server.outputPath)
	fmt.Printf("Updating context every 5 seconds\n")

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := s.server.UpdateContext(); err != nil {
				return fmt.Errorf("error updating context: %v", err)
			}
		}
	}
} 
