package mcp

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"git-tidy/internal/analyzer"
)

// AnalyzerProvider provides context from the analyzer
type AnalyzerProvider struct {
	workDir  string
	analyzer *analyzer.Analyzer
}

// NewAnalyzerProvider creates a new analyzer provider
func NewAnalyzerProvider(workDir string) *AnalyzerProvider {
	return &AnalyzerProvider{
		workDir:  workDir,
		analyzer: analyzer.New(),
	}
}

// GetContext implements the ContextProvider interface
func (p *AnalyzerProvider) GetContext() (Context, error) {
	// Create context
	context := Context{
		Custom: make(map[string]interface{}),
	}

	// Analyze the codebase
	results, err := p.analyzer.AnalyzeDirectory(p.workDir)
	if err != nil {
		return context, fmt.Errorf("error analyzing codebase: %v", err)
	}

	// Add the results to the context
	context.Custom["analyzer_results"] = results

	// Add custom fields
	context.Custom["analyzer_timestamp"] = time.Now()
	context.Custom["analyzer_work_dir"] = p.workDir

	// Add metrics
	context.Custom["package_count"] = len(results.Packages)
	context.Custom["interface_count"] = results.InterfaceCount
	context.Custom["struct_count"] = results.StructCount
	context.Custom["method_count"] = results.MethodCount
	context.Custom["function_count"] = results.FunctionCount
	context.Custom["import_count"] = results.ImportCount
	context.Custom["test_file_count"] = results.TestFileCount

	// Add SOLID violations
	context.Custom["solid_violations"] = results.SOLIDViolations

	// Add suggested improvements
	context.Custom["suggested_improvements"] = results.SuggestedImprovements

	// Check SOLID principles
	p.analyzer.CheckSOLIDPrinciples(results)

	// Create folder context for each package
	for pkgName, metrics := range results.PackageMetrics {
		folderContext := FolderContext{
			Path:        pkgName,
			Description: fmt.Sprintf("Package with %d files, %d interfaces, %d structs, %d methods, %d functions, %d imports, and %d test files",
				len(metrics.Files), len(metrics.Interfaces), len(metrics.Structs),
				len(metrics.Methods), len(metrics.Functions), len(metrics.Imports),
				len(metrics.TestFiles)),
			Files: metrics.Files,
		}
		context.Folders = append(context.Folders, folderContext)

		// Add file contexts for each file in the package
		for _, file := range metrics.Files {
			// Read file content
			content, err := os.ReadFile(filepath.Join(p.workDir, file))
			if err != nil {
				continue // Skip files that can't be read
			}

			fileContext := FileContext{
				Path:    file,
				Content: string(content),
			}
			context.Files = append(context.Files, fileContext)
		}
	}

	// Add improvements
	if len(results.Improvements) > 0 {
		improvements := make([]map[string]interface{}, 0, len(results.Improvements))
		for _, imp := range results.Improvements {
			improvement := map[string]interface{}{
				"type":      imp.Type,
				"file":      imp.File,
				"priority":  imp.Priority,
				"category":  imp.Category,
				"message":   imp.Message,
				"how_to_fix": imp.HowToFix,
			}
			improvements = append(improvements, improvement)
		}
		context.Custom["improvements"] = improvements
	}

	return context, nil
} 
