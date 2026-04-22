package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dmawardi/goScaffold/internal/config"
)

type ModuleGenerator struct {
	config *config.ModuleConfig
}

// NewModuleGenerator creates a new ModuleGenerator instance
func NewModuleGenerator(cfg *config.ModuleConfig) *ModuleGenerator {
	return &ModuleGenerator{
		config: cfg,
	}
}

// Generate creates a new project from the template
func (g *ModuleGenerator) Generate() error {
	// Generate the output path from current directory (project root) to the internal folder to generate a module inside the current project
	outputPath := filepath.Join(".", "internal", g.config.ModuleName)

	// Create output directory with 0755 permissions (rwxr-xr-x):
	// owner gets full read/write/execute; group and others get read/execute only.
	// Execute permission on a directory allows traversal (cd into it), which is
	// required for any tools (go build, editors, etc.) to access files inside.
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	return nil

}
