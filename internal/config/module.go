package config

import (
	"embed"
	"strings"
)

// ModuleConfig holds the configuration for module scaffolding
type ModuleConfig struct {
	ProjectName      string   // Needs to be obtained from the current project directory name
	ProjectPath      string   // Needs to be obtained from the current go.mod file
	ModuleName       string   // Obtained from cli flag
	OutputDir        string   // Directory where the module will be created
	ModuleTemplateFS embed.FS // Optional: embedded module template filesystem (if using embed)
	Force            bool     // Whether to overwrite existing directories
	Verbose          bool     // Whether to show verbose output
}

// GetReplacements returns a map of template-specific replacements
func (c *ModuleConfig) GetReplacements() map[string]string {
	return map[string]string{
		"moduleName":                     strings.ToLower(c.ModuleName[:1]) + c.ModuleName[1:], // e.g., "myModule"
		"github.com/dmawardi/goTemplate": c.ProjectPath,                                        // Note: lowercase 't' to match template
		"ModuleName":                     strings.Title(c.ModuleName),
	}
}
