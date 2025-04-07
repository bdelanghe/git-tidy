package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	a := New()
	if a == nil {
		t.Error("Expected non-nil Analyzer")
	}
	if a.fset == nil {
		t.Error("Expected non-nil FileSet")
	}
}

func TestAnalyzeDirectory(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "analyzer-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test Go file
	testFile := filepath.Join(tmpDir, "test.go")
	code := `package test

import "fmt"

type TestInterface interface {
	DoSomething()
}

type TestStruct struct {
	Field string
}

func (t *TestStruct) DoSomething() {
	fmt.Println(t.Field)
}

func RegularFunction() {
	fmt.Println("Hello")
}
`
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Create a test file
	testTestFile := filepath.Join(tmpDir, "test_test.go")
	testCode := `package test

import "testing"

func TestSomething(t *testing.T) {
	// Test implementation
}
`
	if err := os.WriteFile(testTestFile, []byte(testCode), 0644); err != nil {
		t.Fatalf("Failed to write test test file: %v", err)
	}

	// Analyze the directory
	a := New()
	result, err := a.AnalyzeDirectory(tmpDir)
	if err != nil {
		t.Fatalf("Failed to analyze directory: %v", err)
	}

	// Check the results
	if result.PackageCount != 1 {
		t.Errorf("Expected 1 package, got %d", result.PackageCount)
	}
	if result.InterfaceCount != 1 {
		t.Errorf("Expected 1 interface, got %d", result.InterfaceCount)
	}
	if result.StructCount != 1 {
		t.Errorf("Expected 1 struct, got %d", result.StructCount)
	}
	if result.MethodCount != 1 {
		t.Errorf("Expected 1 method, got %d", result.MethodCount)
	}
	if result.FunctionCount != 2 {
		t.Errorf("Expected 2 functions, got %d", result.FunctionCount)
	}
	if result.TestFileCount != 1 {
		t.Errorf("Expected 1 test file, got %d", result.TestFileCount)
	}
}

func TestCheckSOLIDPrinciples(t *testing.T) {
	// Create a test result with violations
	result := &AnalysisResult{
		PackageMetrics: map[string]PackageMetrics{
			"test": {
				Structs:    make([]string, 6),   // SRP violation
				Methods:    make([]string, 11),  // ISP violation
				Imports:    make([]string, 11),  // DIP violation
				Interfaces: []string{"TestInterface"},
			},
		},
	}

	// Check SOLID principles
	a := New()
	a.CheckSOLIDPrinciples(result)

	// Verify violations
	expectedViolations := 3 // SRP, ISP, and DIP violations
	if len(result.Violations) != expectedViolations {
		t.Errorf("Expected %d violations, got %d", expectedViolations, len(result.Violations))
	}

	// Check for specific violations
	violationTypes := make(map[string]bool)
	for _, v := range result.Violations {
		violationTypes[v.Type] = true
	}

	expectedTypes := []string{"SRP", "ISP", "DIP"}
	for _, expectedType := range expectedTypes {
		if !violationTypes[expectedType] {
			t.Errorf("Expected violation of type %s not found", expectedType)
		}
	}
	
	// Verify improvements
	if len(result.Improvements) == 0 {
		t.Error("Expected improvements to be generated, but none were found")
	}
	
	// Check for specific improvements
	improvementTypes := make(map[string]bool)
	for _, imp := range result.Improvements {
		improvementTypes[imp.Type] = true
	}
	
	// We should have at least the same number of improvements as violations
	if len(result.Improvements) < len(result.Violations) {
		t.Errorf("Expected at least %d improvements, got %d", len(result.Violations), len(result.Improvements))
	}
	
	// Check that each violation has a corresponding improvement
	for _, v := range result.Violations {
		found := false
		for _, imp := range result.Improvements {
			if imp.Type == v.Type && imp.File == v.File {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("No improvement found for violation of type %s in %s", v.Type, v.File)
		}
	}
}

func TestGenerateImprovements(t *testing.T) {
	// Create a test result with violations
	result := &AnalysisResult{
		PackageMetrics: map[string]PackageMetrics{
			"test": {
				Structs:    make([]string, 6),   // SRP violation
				Methods:    make([]string, 11),  // ISP violation
				Imports:    make([]string, 11),  // DIP violation
				Interfaces: []string{"TestInterface"},
				Files:      []string{"test.go"},
				TestFiles:  []string{}, // No test files
			},
		},
		Violations: []Violation{
			{
				Type:     "SRP",
				File:     "test",
				Message:  "Package has too many structs",
				Severity: "warning",
				Rule:     "single-responsibility",
			},
			{
				Type:     "ISP",
				File:     "test",
				Message:  "Interface with too many methods",
				Severity: "warning",
				Rule:     "interface-segregation",
			},
		},
	}
	
	// Generate improvements
	a := New()
	a.GenerateImprovements(result)
	
	// Verify improvements
	if len(result.Improvements) == 0 {
		t.Error("Expected improvements to be generated, but none were found")
	}
	
	// Check for specific improvements
	improvementTypes := make(map[string]bool)
	for _, imp := range result.Improvements {
		improvementTypes[imp.Type] = true
	}
	
	// We should have improvements for each violation plus additional ones
	expectedTypes := []string{"SRP", "ISP", "TestCoverage", "Documentation"}
	for _, expectedType := range expectedTypes {
		if !improvementTypes[expectedType] {
			t.Errorf("Expected improvement of type %s not found", expectedType)
		}
	}
	
	// Check that each improvement has the required fields
	for _, imp := range result.Improvements {
		if imp.Type == "" {
			t.Error("Improvement has empty Type")
		}
		if imp.File == "" {
			t.Error("Improvement has empty File")
		}
		if imp.Message == "" {
			t.Error("Improvement has empty Message")
		}
		if imp.Priority == "" {
			t.Error("Improvement has empty Priority")
		}
		if imp.Category == "" {
			t.Error("Improvement has empty Category")
		}
		if imp.HowToFix == "" {
			t.Error("Improvement has empty HowToFix")
		}
	}
} 
