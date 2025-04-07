package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ContextProvider defines the interface for providing context to Cursor
type ContextProvider interface {
	GetContext() (Context, error)
}

// Context represents the context data to be provided to Cursor
type Context struct {
	Timestamp time.Time              `json:"timestamp"`
	Files     []FileContext         `json:"files,omitempty"`
	Folders   []FolderContext       `json:"folders,omitempty"`
	Code      []CodeContext         `json:"code,omitempty"`
	Git       *GitContext           `json:"git,omitempty"`
	Lint      []LintErrorContext    `json:"lint,omitempty"`
	Custom    map[string]interface{} `json:"custom,omitempty"`
}

// FileContext represents context for a specific file
type FileContext struct {
	Path        string `json:"path"`
	Content     string `json:"content,omitempty"`
	Description string `json:"description,omitempty"`
}

// FolderContext represents context for a folder
type FolderContext struct {
	Path        string   `json:"path"`
	Description string   `json:"description,omitempty"`
	Files       []string `json:"files,omitempty"`
}

// CodeContext represents context for a code snippet
type CodeContext struct {
	File    string `json:"file"`
	Start   int    `json:"start"`
	End     int    `json:"end"`
	Content string `json:"content"`
}

// GitContext represents Git-related context
type GitContext struct {
	Branch      string   `json:"branch,omitempty"`
	Commit      string   `json:"commit,omitempty"`
	Modified    []string `json:"modified,omitempty"`
	Untracked   []string `json:"untracked,omitempty"`
	Deleted     []string `json:"deleted,omitempty"`
	Staged      []string `json:"staged,omitempty"`
	Unstaged    []string `json:"unstaged,omitempty"`
	RecentCommits []string `json:"recent_commits,omitempty"`
}

// LintErrorContext represents a lint error
type LintErrorContext struct {
	File    string `json:"file"`
	Line    int    `json:"line"`
	Column  int    `json:"column"`
	Message string `json:"message"`
	Level   string `json:"level"`
}

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
