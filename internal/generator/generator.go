package generator

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dmawardi/goScaffold/internal/config"
)

// BaseGenerator handles project scaffolding
type BaseGenerator struct {
	config *config.Config
}

// NewBaseGenerator creates a new BaseGenerator instance
func NewBaseGenerator(cfg *config.Config) *BaseGenerator {
	return &BaseGenerator{
		config: cfg,
	}
}

// Generate creates a new project from the template
func (g *BaseGenerator) Generate() error {
	outputPath := filepath.Join(g.config.OutputDir, g.config.ProjectName)

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

	// Initialise go.mod with the user-supplied module path
	if err := g.runGoModInit(outputPath); err != nil {
		return fmt.Errorf("failed to run go mod init: %w", err)
	}

	// Run go mod tidy
	if err := g.runGoModTidy(outputPath); err != nil {
		return fmt.Errorf("failed to run go mod tidy: %w", err)
	}

	return nil
}

// templateRoot is the path prefix used by the embedded TemplateFS.
// All entries in the FS live under this prefix, so it must be stripped
// before constructing the destination path in the output directory.
const templateRoot = "templates/goTemplate"

// processTemplate walks the embedded TemplateFS and mirrors its structure
// into outputDir, applying all placeholder substitutions along the way.
func (g *BaseGenerator) processTemplate(outputDir string) error {
	return fs.WalkDir(g.config.TemplateFS, templateRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Strip the embedded root prefix to obtain the path relative to the
		// template root (e.g. "templates/goTemplate/cmd/main.go" -> "cmd/main.go").
		relPath, err := filepath.Rel(templateRoot, path)
		if err != nil {
			return err
		}

		// The root entry itself ("." after stripping the prefix) has nothing to
		// create — the output directory was already created in Generate.
		if relPath == "." {
			return nil
		}

		// Skip files and directories matched by the configured ignore list.
		if ShouldSkipFile(relPath) {
			if g.config.Verbose {
				fmt.Printf("Skipping: %s\n", relPath)
			}
			return nil
		}

		// Substitute placeholder tokens in each path component
		// (e.g. "goTemplate" -> "myproject") before joining with the output root.
		destPath := g.replacePath(relPath)
		fullDestPath := filepath.Join(outputDir, destPath)

		// If verbose output is enabled, print the file being processed
		if g.config.Verbose {
			fmt.Printf("Processing: %s -> %s\n", relPath, destPath)
		}

		if d.IsDir() {
			// Recreate the directory in the output tree with standard permissions.
			return os.MkdirAll(fullDestPath, 0755)
		}

		// Delegate file copying and content transformation to processFile,
		// passing the full embedded FS path so it can be opened from the FS.
		return g.processFile(path, fullDestPath)
	})
}

// replacePath replaces template placeholders in file/directory paths
func (g *BaseGenerator) replacePath(path string) string {
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

// processFile reads a single file from the embedded TemplateFS (using its full
// FS path), applies content substitution if it is a text file, then writes the
// result to destPath on the real filesystem.
func (g *BaseGenerator) processFile(srcPath, destPath string) error {
	// Ensure the parent directory exists before creating the file.
	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	// Open the source file from the embedded filesystem.
	srcFile, err := g.config.TemplateFS.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create the destination file with standard read/write permissions.
	// Embedded FS files are always read-only (0444), so we use 0644 explicitly
	// to ensure the generated files are editable by the project owner.
	destFile, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Text files have their placeholder tokens substituted before writing.
	// Binary files are streamed directly without modification.
	if IsTextFile(filepath.Base(srcPath)) {
		return g.processTextFile(srcFile, destFile)
	}

	// For binary files, just copy
	_, err = io.Copy(destFile, srcFile)
	return err
}

// processTextFile reads the file contents and applies replacements
func (g *BaseGenerator) processTextFile(src io.Reader, dest io.Writer) error {
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
func (g *BaseGenerator) applyReplacements(content string) string {
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

	// Iterate through ordered keys above and apply replacements within the file
	for _, key := range orderedKeys {
		if replacement, exists := replacements[key]; exists {
			content = strings.ReplaceAll(content, key, replacement)
		}
	}

	return content
}

// runGoModInit runs `go mod init <modulePath>` inside the generated project directory,
// creating the go.mod file with the correct module declaration. The template does not
// ship with a go.mod because the module path is only known at generation time.
func (g *BaseGenerator) runGoModInit(projectDir string) error {
	if g.config.Verbose {
		fmt.Printf("Running 'go mod init %s'...\n", g.config.ModulePath)
	}

	cmd := exec.Command("go", "mod", "init", g.config.ModulePath)
	cmd.Dir = projectDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go mod init failed: %w", err)
	}

	return nil
}

// runGoModTidy runs 'go mod tidy' in the generated project directory
func (g *BaseGenerator) runGoModTidy(projectDir string) error {
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
