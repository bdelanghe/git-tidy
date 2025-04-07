package mcp

import "time"

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
	LastModified time.Time `json:"last_modified,omitempty"`
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

// TestEvent represents a test event from the Go test JSON output
type TestEvent struct {
	Time    string `json:"Time"`
	Action  string `json:"Action"`
	Package string `json:"Package"`
	Test    string `json:"Test,omitempty"`
	Output  string `json:"Output,omitempty"`
	Elapsed string `json:"Elapsed,omitempty"`
}

// Rule represents a Cursor rule
type Rule struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Pattern     string   `json:"pattern,omitempty"`
	Content     string   `json:"content"`
	References  []string `json:"references,omitempty"`
} 
