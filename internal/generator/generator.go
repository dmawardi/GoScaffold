package generator

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dmawardi/goScaffold/internal/config"
)

// Generator handles project scaffolding
type Generator struct {
	config *config.Config
}

// New creates a new Generator instance
func New(cfg *config.Config) *Generator {
	return &Generator{
		config: cfg,
	}
}

// Generate creates a new project from the template
func (g *Generator) Generate() error {
	outputPath := filepath.Join(g.config.OutputDir, g.config.ProjectName)

	// Create output directory
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Copy and process template
	if err := g.processTemplate(g.config.TemplateDir, outputPath); err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	// Run go mod tidy
	if err := g.runGoModTidy(outputPath); err != nil {
		return fmt.Errorf("failed to run go mod tidy: %w", err)
	}

	return nil
}

// processTemplate walks through the template directory and processes each file
func (g *Generator) processTemplate(templateDir, outputDir string) error {
	return filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate relative path from template directory
		relPath, err := filepath.Rel(templateDir, path)
		if err != nil {
			return err
		}

		// Skip if this file should be ignored
		if g.config.ShouldSkipFile(relPath) {
			if g.config.Verbose {
				fmt.Printf("Skipping: %s\n", relPath)
			}
			return nil
		}

		// Replace path components
		destPath := g.replacePath(relPath)
		fullDestPath := filepath.Join(outputDir, destPath)

		if g.config.Verbose {
			fmt.Printf("Processing: %s -> %s\n", relPath, destPath)
		}

		if info.IsDir() {
			// Create directory
			return os.MkdirAll(fullDestPath, info.Mode())
		}

		// Process file
		return g.processFile(path, fullDestPath, info)
	})
}

// replacePath replaces template placeholders in file/directory paths
func (g *Generator) replacePath(path string) string {
	replacements := g.config.GetPathReplacements()

	// Split path and replace each component
	parts := strings.Split(path, string(filepath.Separator))
	for i, part := range parts {
		for old, new := range replacements {
			parts[i] = strings.ReplaceAll(part, old, new)
		}
	}

	return strings.Join(parts, string(filepath.Separator))
}

// processFile copies and processes a single file
func (g *Generator) processFile(srcPath, destPath string, info os.FileInfo) error {
	// Create destination directory if it doesn't exist
	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	// Open source file
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create destination file
	destFile, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Check if this is a text file that needs content replacement
	if g.config.IsTextFile(filepath.Base(srcPath)) {
		return g.processTextFile(srcFile, destFile)
	}

	// For binary files, just copy
	_, err = io.Copy(destFile, srcFile)
	return err
}

// processTextFile reads the file content and applies replacements
func (g *Generator) processTextFile(src io.Reader, dest io.Writer) error {
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
func (g *Generator) applyReplacements(content string) string {
	replacements := g.config.GetReplacements()

	// Apply replacements in a specific order to avoid partial matches
	// Apply longer patterns first to prevent conflicts
	orderedKeys := []string{
		"github.com/dmawardi/goTemplate", // Apply module path first (note: lowercase 't')
		"goTemplate",                     // Then shorter matches
		"GoTemplate",
		"go-template",
		"GO_TEMPLATE",
		"{{.ProjectName}}",
		"{{.ModulePath}}",
		"{{.ProjectNameTitle}}",
		"{{.ProjectNameUpper}}",
	}

	for _, key := range orderedKeys {
		if replacement, exists := replacements[key]; exists {
			content = strings.ReplaceAll(content, key, replacement)
		}
	}

	return content
}

// runGoModTidy runs 'go mod tidy' in the generated project directory
func (g *Generator) runGoModTidy(projectDir string) error {
	if g.config.Verbose {
		fmt.Println("Running 'go mod tidy'...")
	}

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = projectDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		// Don't fail the entire generation if go mod tidy fails
		// Just warn the user
		fmt.Printf("Warning: 'go mod tidy' failed: %v\n", err)
		fmt.Println("You may need to run 'go mod tidy' manually in the project directory.")
	}

	return nil
}
