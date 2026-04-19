package config

import (
	"embed"
	"strings"
)

// Config holds the configuration for project scaffolding
type Config struct {
	ProjectName string   // Name of the project to be created
	OutputDir   string   // Directory where the project will be created
	ModulePath  string   // Go module path (e.g., github.com/user/project)
	TemplateDir string   // Directory containing the template
	TemplateFS  embed.FS // Optional: embedded template filesystem (if using embed)
	Force       bool     // Whether to overwrite existing directories
	Verbose     bool     // Whether to show verbose output
}

// GetReplacements returns a map of template-specific replacements
func (c *Config) GetReplacements() map[string]string {
	return map[string]string{
		"goTemplate":                     c.ProjectName,
		"github.com/dmawardi/goTemplate": c.ModulePath, // Note: lowercase 't' to match template
		"GoTemplate":                     strings.Title(c.ProjectName),
		"go-template":                    c.ProjectName,
		"GO_TEMPLATE":                    strings.ToUpper(c.ProjectName),
		"{{.ProjectName}}":               c.ProjectName,
		"{{.ModulePath}}":                c.ModulePath,
		"{{.ProjectNameTitle}}":          strings.Title(c.ProjectName),
		"{{.ProjectNameUpper}}":          strings.ToUpper(c.ProjectName),
	}
}

// GetPathReplacements returns a map of path-specific replacements (e.g., for renaming directories/files)
func (c *Config) GetPathReplacements() map[string]string {
	return map[string]string{
		"goTemplate": c.ProjectName,
	}
}

// IsTextFile determines if a file should be processed for text replacement
func (c *Config) IsTextFile(filename string) bool {
	textExtensions := []string{
		".go", ".mod", ".sum", ".md", ".txt", ".yaml", ".yml",
		".json", ".toml", ".conf", ".csv", ".tmpl", ".html",
		".js", ".css", ".sh", ".dockerfile", "Dockerfile",
		".gitignore", ".gitattributes", "README", "LICENSE",
		"Makefile", ".env.example",
	}

	// Check if it's a known text extension
	for _, ext := range textExtensions {
		if strings.HasSuffix(strings.ToLower(filename), ext) {
			return true
		}
	}

	// Check if it's a file without extension (often config files)
	if !strings.Contains(filename, ".") {
		return true
	}

	return false
}

// ShouldSkipFile checks a filename against hardcoded patterns of files/directories to skip during generation
func (c *Config) ShouldSkipFile(filename string) bool {
	skipPatterns := []string{
		".git",
		".DS_Store",
		"node_modules",
		"vendor",
		".idea",
		".vscode",
		"*.tmp",
		"*.temp",
		"*.log",
	}

	// Check if filename contains any of the skip patterns
	for _, pattern := range skipPatterns {
		if strings.Contains(filename, pattern) {
			return true
		}
	}

	return false
}
