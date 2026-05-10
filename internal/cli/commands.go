package cli

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/dmawardi/goScaffold/internal/config"
	"github.com/dmawardi/goScaffold/internal/generator"
)

const (
	defaultTemplateDir = "templates/goTemplate"
	defaultOutputDir   = "."
)

// RunCreateCommand handles the "create" subcommand logic: Creating a new project from a template
func RunCreateCommand(templateFS embed.FS) {
	// Create a new FlagSet for the create subcommand
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)

	// Define flags for the create command (shown when "goScaffold create -h" is run)
	var (
		projectName = createCmd.String("name", "", "Project name (required)")
		outputDir   = createCmd.String("output", defaultOutputDir, "Output directory for the new project")
		modulePath  = createCmd.String("module", "", "Go module path (e.g., github.com/user/project)")
		templateDir = createCmd.String("template", defaultTemplateDir, "Template directory to use")
		force       = createCmd.Bool("force", false, "Force creation even if directory exists")
		verbose     = createCmd.Bool("verbose", false, "Verbose output")
		help        = createCmd.Bool("h", false, "Show help for create command")
	)

	// Override the default Usage function for the create command to provide custom help output
	createCmd.Usage = func() {
		fmt.Printf("Usage: %s create [options]\n\n", os.Args[0])
		fmt.Println("Create a new Go project from a template.")
		fmt.Println("\nOptions:")
		createCmd.PrintDefaults()
		fmt.Println("\nExamples:")
		fmt.Printf("  %s create -name myproject -module github.com/user/myproject\n", os.Args[0])
		fmt.Printf("  %s create -name myapi -output /path/to/projects -module github.com/company/myapi\n", os.Args[0])
		fmt.Printf("  %s create -name mycli -template templates/cliTemplate -module github.com/user/mycli\n", os.Args[0])
	}

	// Parse the arguments after "create"
	createCmd.Parse(os.Args[2:])

	// If help flag is set, show usage and exit
	if *help {
		createCmd.Usage()
		return
	}

	// Validate required arguments
	if *projectName == "" {
		fmt.Fprintf(os.Stderr, "Error: Project name is required\n\n")
		createCmd.Usage()
		os.Exit(1)
	}

	// Validate required arguments
	if *modulePath == "" {
		fmt.Fprintf(os.Stderr, "Error: Module path is required\n\n")
		createCmd.Usage()
		os.Exit(1)
	}

	// Validate project name format
	if err := validateProjectName(*projectName); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Set default module path if not provided (backup plan, should be caught by validation above)
	if *modulePath == "" {
		*modulePath = fmt.Sprintf("github.com/user/%s", *projectName)
		if *verbose {
			fmt.Printf("Using default module path: %s\n", *modulePath)
		}
	}

	// Validate module path format
	if err := validateModulePath(*modulePath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Create configuration
	cfg := &config.Config{
		ProjectName: *projectName,
		OutputDir:   *outputDir,
		ModulePath:  *modulePath,
		TemplateDir: *templateDir,
		TemplateFS:  templateFS, // Pass the embedded template filesystem to the config
		Force:       *force,
		Verbose:     *verbose,
	}

	// Validate template directory exists on disk only when a custom template path is provided.
	// The default template is embedded in the binary and does not exist on the filesystem.
	if *templateDir != defaultTemplateDir {
		if _, err := os.Stat(cfg.TemplateDir); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: Template directory '%s' does not exist\n", cfg.TemplateDir)
			os.Exit(1)
		}
	}

	// Create output path
	outputPath := filepath.Join(cfg.OutputDir, cfg.ProjectName)

	// Check if output directory exists
	if _, err := os.Stat(outputPath); !os.IsNotExist(err) && !cfg.Force {
		fmt.Fprintf(os.Stderr, "Error: Directory '%s' already exists. Use --force to overwrite\n", outputPath)
		os.Exit(1)
	}

	// If verbose output is enabled, print the configuration being used
	if cfg.Verbose {
		fmt.Printf("Scaffolding project '%s'...\n", cfg.ProjectName)
		fmt.Printf("Template: %s\n", cfg.TemplateDir)
		fmt.Printf("Output: %s\n", outputPath)
		fmt.Printf("Module: %s\n", cfg.ModulePath)
	}

	// Create the generator and run it
	gen := generator.NewBaseGenerator(cfg)
	if err := gen.Generate(); err != nil {
		log.Fatalf("Failed to generate project: %v", err)
	}

	fmt.Printf("✅ Project '%s' successfully created at '%s'\n", cfg.ProjectName, outputPath)
	fmt.Println("\nNext steps:")
	fmt.Printf("  cd %s\n", outputPath)
	fmt.Printf("  go mod tidy\n")
	fmt.Printf("  go run cmd/%s/main.go\n", cfg.ProjectName)
}

func RunModuleCommand(moduleTemplateFS embed.FS, moduleTestTemplateFS embed.FS) {
	// Create a new FlagSet for the module subcommand
	moduleCmd := flag.NewFlagSet("module", flag.ExitOnError)

	// Define flags for the module command (shown when "goScaffold module -h" is run)
	var (
		moduleName = moduleCmd.String("name", "", "Module name (required)")
		force      = moduleCmd.Bool("force", false, "Force creation even if directory exists")
		verbose    = moduleCmd.Bool("verbose", false, "Verbose output")
		help       = moduleCmd.Bool("h", false, "Show help for module command")
	)

	// Override the default Usage function for the module command to provide custom help output
	moduleCmd.Usage = func() {
		fmt.Printf("Usage: %s module [options]\n\n", os.Args[0])
		fmt.Println("Create a new entity (model/repo/service/controller) inside the current project.")
		fmt.Println("\nOptions:")
		moduleCmd.PrintDefaults()
		fmt.Println("\nExamples:")
		fmt.Printf("  %s module -name mymodule\n", os.Args[0])
		fmt.Printf("  %s module -name mymodule -force\n", os.Args[0])
	}

	// Parse the arguments after "module"
	moduleCmd.Parse(os.Args[2:])

	// If help flag is set, show usage and exit
	if *help {
		moduleCmd.Usage()
		return
	}

	// Validate required arguments
	if *moduleName == "" {
		fmt.Fprintf(os.Stderr, "Error: Module name is required\n\n")
		moduleCmd.Usage()
		os.Exit(1)
	}

	// Grab the project name from the current directory from go.mod file
	projectName, projectPath, err := getCurrentProjectInfo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to get current project info: %v\n", err)
		os.Exit(1)
	}

	// Create configuration
	cfg := &config.ModuleConfig{
		ProjectName:          projectName, // Taken from the current project directory
		ProjectPath:          projectPath, // Taken from the current go.mod file
		ModuleName:           *moduleName,
		OutputDir:            "./internal",         // Modules are created inside the current project's internal directory
		ModuleTemplateFS:     moduleTemplateFS,     // Pass the embedded module template filesystem to the config
		ModuleTestTemplateFS: moduleTestTemplateFS, // Pass the embedded module test template filesystem to the config
		Force:                *force,
		Verbose:              *verbose,
	}

	// If verbose output is enabled, print the configuration being used
	if cfg.Verbose {
		fmt.Printf("Scaffolding entity for project '%s'...\n", cfg.ProjectName)
		fmt.Printf("Output: %s\n", cfg.OutputDir)
		fmt.Printf("Module: %s\n", cfg.ModuleName)
	}

	// Create the generator and run it
	gen := generator.NewModuleGenerator(cfg)
	if err := gen.Generate(); err != nil {
		log.Fatalf("Failed to generate module: %v", err)
	}

	fmt.Printf("✅ Module '%s' successfully created inside the current project\n", cfg.ModuleName)
}

// Checks the current directory for go.mod and extracts the project name and module path
func getCurrentProjectInfo() (string, string, error) {
	// Check if go.mod exists in the current directory
	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		return "", "", fmt.Errorf("go.mod not found in the current directory")
	}

	// Read the go.mod file and extract the module path from the first line
	content, err := os.ReadFile("go.mod")
	if err != nil {
		return "", "", fmt.Errorf("failed to read go.mod: %w", err)
	}

	firstLine := strings.SplitN(string(content), "\n", 2)[0]
	modulePath := strings.TrimPrefix(strings.TrimSpace(firstLine), "module ")
	if modulePath == "" {
		return "", "", fmt.Errorf("could not extract module path from go.mod")
	}

	// Derive the project name as the last path segment of the module path
	projectName := modulePath[strings.LastIndex(modulePath, "/")+1:]

	return projectName, modulePath, nil
}
