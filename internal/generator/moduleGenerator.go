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

	// Edit config.go in cmd/projectName to add the new module's setState function to the list of module initialisation functions
	if err := g.updateSetState(); err != nil {
		return fmt.Errorf("failed to update setState: %w", err)
	}

	// Add Migration to model.go for migration
	if err := g.updateMigration(); err != nil {
		return fmt.Errorf("failed to update migration: %w", err)
	}

	// Add new module's stack to the module stack in module.go, then init¸ repo, service, & controller to buildModuleStack
	if err := g.updateModuleStack(); err != nil {
		return fmt.Errorf("failed to update module stack: %w", err)
	}

	// Register routes in routes.go
	if err := g.updateRoutes(); err != nil {
		return fmt.Errorf("failed to update routes: %w", err)
	}

	// Make e2e tests in cmd/projectName folder of current project
	if err := g.updateE2ETests(); err != nil {
		return fmt.Errorf("failed to update e2e tests: %w", err)
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

// updateSetState edits config.go in the ./cmd/<ProjectName> folder to add the new module's setState function to the list of module initialisation functions, and adds the necessary import for the module.
func (g *ModuleGenerator) updateSetState() error {
	configPath, err := g.findCmdFile("config.go", "stateFuncs")
	if err != nil {
		return err
	}

	contentBytes, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", configPath, err)
	}
	content := string(contentBytes)

	// Derive the package-level module name (lowercase first letter)
	moduleName := strings.ToLower(g.config.ModuleName[:1]) + g.config.ModuleName[1:]
	importPath := fmt.Sprintf(`"%s/internal/%s"`, g.config.ProjectPath, moduleName)

	// Add import if not already present
	if !strings.Contains(content, importPath) {
		content, err = insertImport(content, importPath)
		if err != nil {
			return fmt.Errorf("failed to insert import: %w", err)
		}
	}

	// Add SetState entry to stateFuncs if not already present
	setStateEntry := fmt.Sprintf("%s.SetState,", moduleName)
	if !strings.Contains(content, setStateEntry) {
		content, err = insertStateFuncsEntry(content, setStateEntry)
		if err != nil {
			return fmt.Errorf("failed to insert SetState entry: %w", err)
		}
	}

	return os.WriteFile(configPath, []byte(content), 0644)
}

// updateMigration edits model.go in the ./cmd/<ProjectName> folder to add the new module's model to the list of models to migrate, and adds the necessary import for the module.
func (g *ModuleGenerator) updateMigration() error {
	mainPath, err := g.findCmdFile("model.go", "modelsToMigrate")
	if err != nil {
		return err
	}

	contentBytes, err := os.ReadFile(mainPath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", mainPath, err)
	}
	content := string(contentBytes)

	// Derive names: lowercase module name for package ref, title-case for struct name
	moduleName := strings.ToLower(g.config.ModuleName[:1]) + g.config.ModuleName[1:]
	modelName := strings.ToUpper(g.config.ModuleName[:1]) + g.config.ModuleName[1:]
	importPath := fmt.Sprintf(`"%s/internal/%s"`, g.config.ProjectPath, moduleName)

	// Add import if not already present
	if !strings.Contains(content, importPath) {
		content, err = insertImport(content, importPath)
		if err != nil {
			return fmt.Errorf("failed to insert import: %w", err)
		}
	}

	// Add migration entry to modelsToMigrate if not already present
	migrationEntry := fmt.Sprintf("&%s.%s{},", moduleName, modelName)
	if !strings.Contains(content, migrationEntry) {
		content, err = insertMigrationEntry(content, migrationEntry)
		if err != nil {
			return fmt.Errorf("failed to insert migration entry: %w", err)
		}
	}

	return os.WriteFile(mainPath, []byte(content), 0644)
}

// findCmdFile locates a file by name under cmd/, using ProjectName as the preferred subdir.
// If the preferred path does not exist, it scans all cmd/ subdirs and returns the first
// whose contents contain marker.
func (g *ModuleGenerator) findCmdFile(filename, marker string) (string, error) {
	if g.config.ProjectName != "" {
		path := filepath.Join(".", "cmd", g.config.ProjectName, filename)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	cmdDir := filepath.Join(".", "cmd")
	entries, err := os.ReadDir(cmdDir)
	if err != nil {
		return "", fmt.Errorf("failed to read cmd directory: %w", err)
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		path := filepath.Join(cmdDir, entry.Name(), filename)
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		if strings.Contains(string(data), marker) {
			return path, nil
		}
	}
	return "", fmt.Errorf("could not find %s containing %q under cmd/", filename, marker)
}

// updateModuleStack edits module.go in the ./cmd/<ProjectName> folder to add the new module's stack to the buildModuleStack function, including initializing the repository, service, and controller.
func (g *ModuleGenerator) updateModuleStack() error {
	modulePath, err := g.findCmdFile("module.go", "buildModuleStack")
	if err != nil {
		return err
	}

	contentBytes, err := os.ReadFile(modulePath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", modulePath, err)
	}
	content := string(contentBytes)

	// Derive names: lowercase for package/field reference, title-case for type names
	moduleName := strings.ToLower(g.config.ModuleName[:1]) + g.config.ModuleName[1:]
	modelName := strings.ToUpper(g.config.ModuleName[:1]) + g.config.ModuleName[1:]
	importPath := fmt.Sprintf(`"%s/internal/%s"`, g.config.ProjectPath, moduleName)

	// Add import if not already present
	if !strings.Contains(content, importPath) {
		content, err = insertImport(content, importPath)
		if err != nil {
			return fmt.Errorf("failed to insert import: %w", err)
		}
	}

	// Add field to moduleStack struct if not already present
	structField := fmt.Sprintf("%s %s.%sModule", moduleName, moduleName, modelName)
	if !strings.Contains(content, structField) {
		content, err = insertModuleStackField(content, structField)
		if err != nil {
			return fmt.Errorf("failed to insert module stack field: %w", err)
		}
	}

	// Add initialization block before return module if not already present
	moduleInit := fmt.Sprintf("module.%s = %s.%sModule{}", moduleName, moduleName, modelName)
	if !strings.Contains(content, moduleInit) {
		initBlock := fmt.Sprintf(
			"// Create %s module\n\tmodule.%s = %s.%sModule{}\n\tmodule.%s.Repository = %s.New%sRepository(app.DbClient)\n\tmodule.%s.Service = %s.New%sService(module.%s.Repository)\n\tmodule.%s.Controller = %s.New%sController(module.%s.Service)",
			moduleName,
			moduleName, moduleName, modelName,
			moduleName, moduleName, modelName,
			moduleName, moduleName, modelName, moduleName,
			moduleName, moduleName, modelName, moduleName,
		)
		content, err = insertModuleStackInit(content, initBlock)
		if err != nil {
			return fmt.Errorf("failed to insert module stack init: %w", err)
		}
	}

	return os.WriteFile(modulePath, []byte(content), 0644)
}

// updateRoutes edits routes.go in the ./cmd/<ProjectName> folder to register the new module's routes with the router.
func (g *ModuleGenerator) updateRoutes() error {
	configPath, err := g.findCmdFile("routes.go", "stateFuncs")
	if err != nil {
		return err
	}

	contentBytes, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", configPath, err)
	}
	content := string(contentBytes)

	// Derive the package-level module name (lowercase first letter)
	moduleName := strings.ToLower(g.config.ModuleName[:1]) + g.config.ModuleName[1:]
	importPath := fmt.Sprintf(`"%s/internal/%s"`, g.config.ProjectPath, moduleName)

	// Add import if not already present
	if !strings.Contains(content, importPath) {
		content, err = insertImport(content, importPath)
		if err != nil {
			return fmt.Errorf("failed to insert import: %w", err)
		}
	}

	routeEntry := fmt.Sprintf("%s.RegisterRoutes(api, modules.%s.Controller)", moduleName, moduleName)
	if !strings.Contains(content, routeEntry) {
		content, err = insertRouteEntry(content, routeEntry)
		if err != nil {
			return fmt.Errorf("failed to insert route entry: %w", err)
		}
	}

	return os.WriteFile(configPath, []byte(content), 0644)
}

func (g *ModuleGenerator) updateE2ETests() error {
	// TODO
	return nil
}
