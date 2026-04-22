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
