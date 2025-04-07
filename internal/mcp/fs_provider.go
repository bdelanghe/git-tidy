package mcp

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FSProvider provides context from file system changes
type FSProvider struct {
	workDir string
}

// NewFSProvider creates a new file system provider
func NewFSProvider(workDir string) *FSProvider {
	return &FSProvider{
		workDir: workDir,
	}
}

// GetContext implements the ContextProvider interface
func (p *FSProvider) GetContext() (Context, error) {
	// Create context
	context := Context{
		Custom: make(map[string]interface{}),
	}

	// Walk the directory tree
	var files []FileContext
	var folders []FolderContext
	err := filepath.Walk(p.workDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip .git directory
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}

		// Skip .cursor directory
		if info.IsDir() && info.Name() == ".cursor" {
			return filepath.SkipDir
		}

		// Get the relative path
		relPath, err := filepath.Rel(p.workDir, path)
		if err != nil {
			return err
		}

		// Skip hidden files and directories
		if strings.HasPrefix(relPath, ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Add the file or folder to the context
		if info.IsDir() {
			folders = append(folders, FolderContext{
				Path:        relPath,
				Description: fmt.Sprintf("Directory: %s", relPath),
			})
		} else {
			// Read the file content
			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("error reading file %s: %v", path, err)
			}

			files = append(files, FileContext{
				Path:        relPath,
				Content:     string(content),
				LastModified: info.ModTime(),
			})
		}

		return nil
	})

	if err != nil {
		return context, fmt.Errorf("error walking directory: %v", err)
	}

	// Add the files and folders to the context
	context.Files = files
	context.Folders = folders

	// Add custom fields
	context.Custom["file_count"] = len(files)
	context.Custom["folder_count"] = len(folders)
	context.Custom["last_scan"] = time.Now()

	return context, nil
} 
