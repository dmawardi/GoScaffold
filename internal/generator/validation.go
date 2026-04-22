package generator

import "strings"

// IsTextFile determines if a file should be processed for text replacement
func IsTextFile(filename string) bool {
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
func ShouldSkipFile(filename string) bool {
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
