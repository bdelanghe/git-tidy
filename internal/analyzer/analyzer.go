package analyzer

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// AnalysisResult represents the results of analyzing a Go codebase
type AnalysisResult struct {
	PackageCount    int
	InterfaceCount  int
	StructCount     int
	MethodCount     int
	FunctionCount   int
	ImportCount     int
	TestFileCount   int
	PackageMetrics  map[string]PackageMetrics
	Violations      []Violation
}

// PackageMetrics contains metrics for a specific package
type PackageMetrics struct {
	Files          []string
	Interfaces     []string
	Structs        []string
	Methods        []string
	Functions      []string
	Imports        []string
	Dependencies   []string
	TestFiles      []string
	Complexity     float64
	TestCoverage   float64
}

// Violation represents a code quality violation
type Violation struct {
	Type        string
	File        string
	Line        int
	Message     string
	Severity    string
	Rule        string
}

// Analyzer is the main analyzer type
type Analyzer struct {
	fset *token.FileSet
}

// New creates a new Analyzer
func New() *Analyzer {
	return &Analyzer{
		fset: token.NewFileSet(),
	}
}

// AnalyzeDirectory analyzes a directory of Go code
func (a *Analyzer) AnalyzeDirectory(dir string) (*AnalysisResult, error) {
	result := &AnalysisResult{
		PackageMetrics: make(map[string]PackageMetrics),
	}

	// Walk through the directory
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-Go files
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		// Parse the file
		f, err := parser.ParseFile(a.fset, path, nil, parser.ParseComments)
		if err != nil {
			return err
		}

		// Analyze the file
		a.analyzeFile(f, path, result)

		return nil
	})

	return result, err
}

// analyzeFile analyzes a single Go file
func (a *Analyzer) analyzeFile(f *ast.File, path string, result *AnalysisResult) {
	pkgName := f.Name.Name
	metrics, exists := result.PackageMetrics[pkgName]
	if !exists {
		metrics = PackageMetrics{}
		result.PackageMetrics[pkgName] = metrics
		result.PackageCount++
	}

	metrics.Files = append(metrics.Files, path)

	// Analyze imports
	for _, imp := range f.Imports {
		metrics.Imports = append(metrics.Imports, imp.Path.Value)
		result.ImportCount++
	}

	// Analyze declarations
	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			switch d.Tok {
			case token.TYPE:
				for _, spec := range d.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						switch typeSpec.Type.(type) {
						case *ast.InterfaceType:
							metrics.Interfaces = append(metrics.Interfaces, typeSpec.Name.Name)
							result.InterfaceCount++
						case *ast.StructType:
							metrics.Structs = append(metrics.Structs, typeSpec.Name.Name)
							result.StructCount++
						}
					}
				}
			}
		case *ast.FuncDecl:
			if d.Recv != nil {
				metrics.Methods = append(metrics.Methods, d.Name.Name)
				result.MethodCount++
			} else {
				metrics.Functions = append(metrics.Functions, d.Name.Name)
				result.FunctionCount++
			}
		}
	}

	// Check for test files
	if strings.HasSuffix(path, "_test.go") {
		metrics.TestFiles = append(metrics.TestFiles, path)
		result.TestFileCount++
	}

	result.PackageMetrics[pkgName] = metrics
}

// CheckSOLIDPrinciples checks if the code follows SOLID principles
func (a *Analyzer) CheckSOLIDPrinciples(result *AnalysisResult) {
	// Single Responsibility Principle
	for pkgName, metrics := range result.PackageMetrics {
		if len(metrics.Structs) > 5 {
			result.Violations = append(result.Violations, Violation{
				Type:     "SRP",
				File:     pkgName,
				Message:  "Package has too many structs, might violate Single Responsibility Principle",
				Severity: "warning",
				Rule:     "single-responsibility",
			})
		}
	}

	// Open/Closed Principle
	for pkgName, metrics := range result.PackageMetrics {
		if len(metrics.Interfaces) > 0 && len(metrics.Structs) == 0 {
			result.Violations = append(result.Violations, Violation{
				Type:     "OCP",
				File:     pkgName,
				Message:  "Interface without implementing structs might violate Open/Closed Principle",
				Severity: "warning",
				Rule:     "open-closed",
			})
		}
	}

	// Interface Segregation Principle
	for pkgName, metrics := range result.PackageMetrics {
		if len(metrics.Methods) > 10 {
			result.Violations = append(result.Violations, Violation{
				Type:     "ISP",
				File:     pkgName,
				Message:  "Interface with too many methods might violate Interface Segregation Principle",
				Severity: "warning",
				Rule:     "interface-segregation",
			})
		}
	}

	// Dependency Inversion Principle
	for pkgName, metrics := range result.PackageMetrics {
		if len(metrics.Imports) > 10 {
			result.Violations = append(result.Violations, Violation{
				Type:     "DIP",
				File:     pkgName,
				Message:  "Package with too many imports might violate Dependency Inversion Principle",
				Severity: "warning",
				Rule:     "dependency-inversion",
			})
		}
	}
} 
