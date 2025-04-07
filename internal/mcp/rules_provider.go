package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Rule represents a Cursor rule
type Rule struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Pattern     string   `json:"pattern,omitempty"`
	Content     string   `json:"content"`
	References  []string `json:"references,omitempty"`
}

// RulesProvider provides context from Cursor rules
type RulesProvider struct {
	workDir string
}

// NewRulesProvider creates a new rules provider
func NewRulesProvider(workDir string) *RulesProvider {
	return &RulesProvider{
		workDir: workDir,
	}
}

// GetContext implements the ContextProvider interface
func (p *RulesProvider) GetContext() (Context, error) {
	// Create context
	context := Context{
		Custom: make(map[string]interface{}),
	}

	// Check if .cursor/rules directory exists
	rulesDir := filepath.Join(p.workDir, ".cursor", "rules")
	if _, err := os.Stat(rulesDir); os.IsNotExist(err) {
		// Create the directory if it doesn't exist
		if err := os.MkdirAll(rulesDir, 0755); err != nil {
			return context, fmt.Errorf("error creating rules directory: %v", err)
		}

		// Create a default rule
		defaultRule := Rule{
			Name:        "default",
			Description: "Default rule for the project",
			Content:     "This is a default rule for the project. You can customize it to provide specific instructions to the AI.",
		}

		// Write the default rule
		defaultRulePath := filepath.Join(rulesDir, "default.json")
		if err := p.writeRule(defaultRulePath, defaultRule); err != nil {
			return context, fmt.Errorf("error writing default rule: %v", err)
		}

		// Add the default rule to the context
		context.Custom["rules"] = []Rule{defaultRule}
		return context, nil
	}

	// Read all rule files
	rules := make([]Rule, 0)
	err := filepath.Walk(rulesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Skip non-JSON files
		if !strings.HasSuffix(path, ".json") {
			return nil
		}

		// Read the rule file
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading rule file %s: %v", path, err)
		}

		// Parse the rule
		var rule Rule
		if err := json.Unmarshal(data, &rule); err != nil {
			return fmt.Errorf("error parsing rule file %s: %v", path, err)
		}

		// Add the rule to the list
		rules = append(rules, rule)
		return nil
	})

	if err != nil {
		return context, fmt.Errorf("error reading rules: %v", err)
	}

	// Add the rules to the context
	context.Custom["rules"] = rules

	// Add a folder context for the rules directory
	context.Folders = append(context.Folders, FolderContext{
		Path:        ".cursor/rules",
		Description: "Directory containing Cursor rules for the project",
		Files:       p.getRuleFiles(rulesDir),
	})

	return context, nil
}

// writeRule writes a rule to a file
func (p *RulesProvider) writeRule(path string, rule Rule) error {
	data, err := json.MarshalIndent(rule, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling rule: %v", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("error writing rule file: %v", err)
	}

	return nil
}

// getRuleFiles gets the list of rule files in the rules directory
func (p *RulesProvider) getRuleFiles(rulesDir string) []string {
	files := make([]string, 0)
	filepath.Walk(rulesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Skip non-JSON files
		if !strings.HasSuffix(path, ".json") {
			return nil
		}

		// Add the file to the list
		relPath, err := filepath.Rel(p.workDir, path)
		if err != nil {
			return err
		}
		files = append(files, relPath)
		return nil
	})
	return files
} 
