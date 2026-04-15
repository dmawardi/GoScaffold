package cli

import (
	"fmt"
	"regexp"
	"strings"
)

// validateModulePath validates the Go module path format
func validateModulePath(path string) error {
	if len(path) == 0 {
		return fmt.Errorf("module path cannot be empty")
	}

	// Basic validation for module path format
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		return fmt.Errorf("module path '%s' should be in format 'domain.com/user/project'", path)
	}

	// Check for valid characters
	validPath := regexp.MustCompile(`^[a-zA-Z0-9._/-]+$`)
	if !validPath.MatchString(path) {
		return fmt.Errorf("module path '%s' contains invalid characters", path)
	}

	return nil
}

// validateProjectName validates the project name format
func validateProjectName(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("project name cannot be empty")
	}

	// Check for valid Go package name (similar to directory name rules)
	validName := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)
	if !validName.MatchString(name) {
		return fmt.Errorf("project name '%s' must start with a letter and contain only letters, numbers, underscores, and hyphens", name)
	}

	// Check for reserved names
	reserved := []string{"main", "test", "vendor", "internal"}
	for _, r := range reserved {
		if strings.EqualFold(name, r) {
			return fmt.Errorf("project name '%s' is a reserved name", name)
		}
	}

	return nil
}
