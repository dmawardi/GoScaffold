package generator

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/dmawardi/goScaffold/internal/config"
)

// FileProcessor handles advanced file processing operations
type FileProcessor struct {
	config *config.Config
}

// NewFileProcessor creates a new FileProcessor instance
func NewFileProcessor(cfg *config.Config) *FileProcessor {
	return &FileProcessor{
		config: cfg,
	}
}

// ProcessGoFiles specifically handles Go files with import path replacement
func (fp *FileProcessor) ProcessGoFiles(srcPath, destPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	return fp.processGoFileContent(srcFile, destFile)
}

// processGoFileContent processes Go file content line by line for better import handling
func (fp *FileProcessor) processGoFileContent(src io.Reader, dest io.Writer) error {
	scanner := bufio.NewScanner(src)
	writer := bufio.NewWriter(dest)
	defer writer.Flush()

	for scanner.Scan() {
		line := scanner.Text()
		processedLine := fp.processGoLine(line)

		if _, err := writer.WriteString(processedLine + "\n"); err != nil {
			return err
		}
	}

	return scanner.Err()
}

// processGoLine processes a single line of Go code
func (fp *FileProcessor) processGoLine(line string) string {
	// Handle import statements
	if strings.Contains(line, "import") && strings.Contains(line, "github.com/dmawardi/goTemplate") {
		return strings.ReplaceAll(line, "github.com/dmawardi/goTemplate", fp.config.ModulePath)
	}

	// Handle other template replacements
	replacements := fp.config.GetReplacements()
	for old, new := range replacements {
		line = strings.ReplaceAll(line, old, new)
	}

	return line
}

// ProcessGoModFile specifically handles go.mod files
func (fp *FileProcessor) ProcessGoModFile(srcPath, destPath string) error {
	content, err := os.ReadFile(srcPath)
	if err != nil {
		return err
	}

	// Replace module declaration
	processedContent := string(content)
	processedContent = strings.ReplaceAll(processedContent, "module github.com/dmawardi/goTemplate", fmt.Sprintf("module %s", fp.config.ModulePath))

	// Apply other replacements
	replacements := fp.config.GetReplacements()
	for old, new := range replacements {
		processedContent = strings.ReplaceAll(processedContent, old, new)
	}

	return os.WriteFile(destPath, []byte(processedContent), 0644)
}

// ProcessDockerfile specifically handles Docker files
func (fp *FileProcessor) ProcessDockerfile(srcPath, destPath string) error {
	content, err := os.ReadFile(srcPath)
	if err != nil {
		return err
	}

	processedContent := string(content)

	// Replace binary names and paths in Dockerfile
	processedContent = strings.ReplaceAll(processedContent, "goTemplate", fp.config.ProjectName)

	// Apply other replacements
	replacements := fp.config.GetReplacements()
	for old, new := range replacements {
		processedContent = strings.ReplaceAll(processedContent, old, new)
	}

	return os.WriteFile(destPath, []byte(processedContent), 0644)
}

// CopyBinaryFile copies binary files without modification
func (fp *FileProcessor) CopyBinaryFile(srcPath, destPath string, mode os.FileMode) error {
	// Create destination directory if it doesn't exist
	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	return err
}

// GetProcessingStrategy returns the appropriate processing strategy for a file
func (fp *FileProcessor) GetProcessingStrategy(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	base := strings.ToLower(filepath.Base(filename))

	switch {
	case ext == ".go":
		return "go_file"
	case base == "go.mod":
		return "go_mod"
	case base == "dockerfile" || strings.HasSuffix(base, ".dockerfile"):
		return "dockerfile"
	case IsTextFile(filename):
		return "text_file"
	default:
		return "binary_file"
	}
}
