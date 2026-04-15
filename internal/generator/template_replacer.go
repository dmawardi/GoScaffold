package generator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dmawardi/goScaffold/internal/config"
)

// TemplateReplacer handles advanced template string replacement
type TemplateReplacer struct {
	config       *config.Config
	replacements map[string]string
	regexCache   map[string]*regexp.Regexp
}

// NewTemplateReplacer creates a new TemplateReplacer instance
func NewTemplateReplacer(cfg *config.Config) *TemplateReplacer {
	return &TemplateReplacer{
		config:       cfg,
		replacements: cfg.GetReplacements(),
		regexCache:   make(map[string]*regexp.Regexp),
	}
}

// ReplaceInContent performs all template replacements in the given content
func (tr *TemplateReplacer) ReplaceInContent(content string) string {
	// Apply basic string replacements first
	for old, new := range tr.replacements {
		content = strings.ReplaceAll(content, old, new)
	}

	// Apply regex-based replacements for more complex patterns
	content = tr.replaceImportPaths(content)
	content = tr.replacePackageNames(content)
	content = tr.replaceVariableNames(content)

	return content
}

// replaceImportPaths handles Go import path replacement with regex
func (tr *TemplateReplacer) replaceImportPaths(content string) string {
	// Pattern to match import statements
	importPattern := `(import\s+(?:\([\s\S]*?\)|"[^"]*"|` + "`[^`]*`" + `))`

	regex, err := tr.getOrCreateRegex(importPattern)
	if err != nil {
		// Fallback to simple string replacement if regex fails
		return strings.ReplaceAll(content, "github.com/dmawardi/goTemplate", tr.config.ModulePath)
	}

	return regex.ReplaceAllStringFunc(content, func(match string) string {
		return strings.ReplaceAll(match, "github.com/dmawardi/goTemplate", tr.config.ModulePath)
	})
}

// replacePackageNames handles package name replacements
func (tr *TemplateReplacer) replacePackageNames(content string) string {
	// Replace package declarations that might reference the template name
	packagePattern := `package\s+goTemplate`

	regex, err := tr.getOrCreateRegex(packagePattern)
	if err != nil {
		return content
	}

	return regex.ReplaceAllString(content, fmt.Sprintf("package %s", tr.config.ProjectName))
}

// replaceVariableNames handles variable and function name replacements
func (tr *TemplateReplacer) replaceVariableNames(content string) string {
	// Replace camelCase variable names
	camelPattern := `\bgoTemplate([A-Z][a-zA-Z0-9]*)`

	regex, err := tr.getOrCreateRegex(camelPattern)
	if err != nil {
		return content
	}

	return regex.ReplaceAllStringFunc(content, func(match string) string {
		suffix := strings.TrimPrefix(match, "goTemplate")
		return tr.config.ProjectName + suffix
	})
}

// ReplaceInPath handles path component replacement
func (tr *TemplateReplacer) ReplaceInPath(path string) string {
	pathReplacements := tr.config.GetPathReplacements()

	for old, new := range pathReplacements {
		path = strings.ReplaceAll(path, old, new)
	}

	return path
}

// ReplaceGoModContent specifically handles go.mod file content
func (tr *TemplateReplacer) ReplaceGoModContent(content string) string {
	// Replace module declaration
	modulePattern := `module\s+[^\s]+`

	regex, err := tr.getOrCreateRegex(modulePattern)
	if err != nil {
		// Fallback to simple replacement
		return strings.ReplaceAll(content, "github.com/dmawardi/goTemplate", tr.config.ModulePath)
	}

	content = regex.ReplaceAllString(content, fmt.Sprintf("module %s", tr.config.ModulePath))

	// Apply other replacements
	return tr.ReplaceInContent(content)
}

// ReplaceDockerContent specifically handles Docker file content
func (tr *TemplateReplacer) ReplaceDockerContent(content string) string {
	// Common Docker replacements
	dockerReplacements := map[string]string{
		"goTemplate":      tr.config.ProjectName,
		"GOTEMPLATE":      strings.ToUpper(tr.config.ProjectName),
		"/app/goTemplate": fmt.Sprintf("/app/%s", tr.config.ProjectName),
	}

	for old, new := range dockerReplacements {
		content = strings.ReplaceAll(content, old, new)
	}

	return tr.ReplaceInContent(content)
}

// ReplaceMarkdownContent handles README and other markdown files
func (tr *TemplateReplacer) ReplaceMarkdownContent(content string) string {
	// Replace titles and headings
	content = strings.ReplaceAll(content, "# goTemplate", fmt.Sprintf("# %s", tr.config.ProjectName))
	content = strings.ReplaceAll(content, "## goTemplate", fmt.Sprintf("## %s", tr.config.ProjectName))

	// Replace code blocks that might contain the old module path
	codeBlockPattern := "```[\\s\\S]*?```"

	regex, err := tr.getOrCreateRegex(codeBlockPattern)
	if err != nil {
		return tr.ReplaceInContent(content)
	}

	content = regex.ReplaceAllStringFunc(content, func(match string) string {
		return strings.ReplaceAll(match, "github.com/dmawardi/goTemplate", tr.config.ModulePath)
	})

	return tr.ReplaceInContent(content)
}

// getOrCreateRegex gets a compiled regex from cache or creates and caches it
func (tr *TemplateReplacer) getOrCreateRegex(pattern string) (*regexp.Regexp, error) {
	if regex, exists := tr.regexCache[pattern]; exists {
		return regex, nil
	}

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	tr.regexCache[pattern] = regex
	return regex, nil
}

// ValidateReplacements checks if all required replacements are valid
func (tr *TemplateReplacer) ValidateReplacements() error {
	if tr.config.ProjectName == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	if tr.config.ModulePath == "" {
		return fmt.Errorf("module path cannot be empty")
	}

	// Validate that project name is a valid Go package name
	validName := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)
	if !validName.MatchString(tr.config.ProjectName) {
		return fmt.Errorf("project name '%s' is not a valid Go package name", tr.config.ProjectName)
	}

	return nil
}

// GetReplacementSummary returns a summary of all replacements that will be made
func (tr *TemplateReplacer) GetReplacementSummary() map[string]string {
	summary := make(map[string]string)

	for old, new := range tr.replacements {
		summary[old] = new
	}

	return summary
}

type Replacer struct {
	projectName  string
	username     string
	replacements map[string]string
}

func newReplacer(projectName, username string) *Replacer {

	var replacements = map[string]string{
		"goTemplate":                     projectName,
		"github.com/dmawardi/goTemplate": fmt.Sprintf("github.com/%s/%s", username, projectName),
	}

	return &Replacer{
		projectName:  projectName,
		username:     username,
		replacements: replacements,
	}
}
