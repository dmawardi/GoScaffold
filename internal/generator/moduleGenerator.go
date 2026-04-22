package generator

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

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

	// Copy and process template files from the embedded filesystem
	if err := g.processTemplate(outputPath); err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	return nil
}

// moduleTemplateRoot is the path prefix for the module template within the embedded FS.
const moduleTemplateRoot = "templates/moduleName"

func (g *ModuleGenerator) processTemplate(outputPath string) error {
	// Walk through the embedded template filesystem and copy files to the output directory
	return fs.WalkDir(g.config.ModuleTemplateFS, moduleTemplateRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories (they will be created as needed when copying files)
		if d.IsDir() {
			return nil
		}

		// Construct the destination path by stripping the template root prefix
		relPath := strings.TrimPrefix(path, moduleTemplateRoot+"/")
		destPath := filepath.Join(outputPath, relPath)

		// Ensure the destination directory exists
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return fmt.Errorf("failed to create destination directory: %w", err)
		}

		// Open the source file from the embedded filesystem
		srcFile, err := g.config.ModuleTemplateFS.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open source file: %w", err)
		}
		defer srcFile.Close()

		// Create the destination file
		destFile, err := os.Create(destPath)
		if err != nil {
			return fmt.Errorf("failed to create destination file: %w", err)
		}
		defer destFile.Close()

		// Text files have their placeholder tokens substituted before writing.
		// Binary files are streamed directly without modification.
		if IsTextFile(filepath.Base(path)) {
			if err := g.processTextFile(srcFile, destFile); err != nil {
				return fmt.Errorf("failed to process text file: %w", err)
			}
		} else {
			// Else, copy binary files directly without modification
			if _, err := io.Copy(destFile, srcFile); err != nil {
				return fmt.Errorf("failed to copy binary file: %w", err)
			}
		}

		return nil
	})
}

// processTextFile reads the file contents and applies replacements
func (g *ModuleGenerator) processTextFile(src io.Reader, dest io.Writer) error {
	// Read entire file content
	content, err := io.ReadAll(src)
	if err != nil {
		return err
	}

	// Apply replacements
	processedContent := g.applyReplacements(string(content))

	// Write processed content
	_, err = dest.Write([]byte(processedContent))
	return err
}

// applyReplacements applies all configured text replacements
func (g *ModuleGenerator) applyReplacements(content string) string {
	replacements := g.config.GetReplacements()

	// Apply replacements in a specific order to avoid partial matches
	// Apply longer patterns first to prevent conflicts
	orderedKeys := []string{
		"github.com/dmawardi/goTemplate", // Apply module path first (note: lowercase 't')
		"moduleName",                     // Then shorter matches
		"ModuleName",
	}

	// Iterate through ordered keys above and apply replacements within the file
	for _, key := range orderedKeys {
		if replacement, exists := replacements[key]; exists {
			content = strings.ReplaceAll(content, key, replacement)
		}
	}

	return content
}
