package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Server represents an MCP server
type Server struct {
	contextProviders []ContextProvider
	outputPath      string
}

// NewServer creates a new MCP server
func NewServer(outputPath string) *Server {
	return &Server{
		contextProviders: make([]ContextProvider, 0),
		outputPath:      outputPath,
	}
}

// AddProvider adds a context provider to the server
func (s *Server) AddProvider(provider ContextProvider) {
	s.contextProviders = append(s.contextProviders, provider)
}

// UpdateContext updates the context file with data from all providers
func (s *Server) UpdateContext() error {
	// Create a new context
	context := Context{
		Timestamp: time.Now(),
		Custom:    make(map[string]interface{}),
	}

	// Collect context from all providers
	for _, provider := range s.contextProviders {
		providerContext, err := provider.GetContext()
		if err != nil {
			return fmt.Errorf("error getting context from provider: %v", err)
		}

		// Merge contexts
		context.Files = append(context.Files, providerContext.Files...)
		context.Folders = append(context.Folders, providerContext.Folders...)
		context.Code = append(context.Code, providerContext.Code...)
		if providerContext.Git != nil {
			context.Git = providerContext.Git
		}
		context.Lint = append(context.Lint, providerContext.Lint...)
		
		// Merge custom fields
		for k, v := range providerContext.Custom {
			context.Custom[k] = v
		}
	}

	// Ensure the output directory exists
	if err := os.MkdirAll(filepath.Dir(s.outputPath), 0755); err != nil {
		return fmt.Errorf("error creating output directory: %v", err)
	}

	// Write context to file
	data, err := json.MarshalIndent(context, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling context: %v", err)
	}

	if err := os.WriteFile(s.outputPath, data, 0644); err != nil {
		return fmt.Errorf("error writing context file: %v", err)
	}

	return nil
} 
