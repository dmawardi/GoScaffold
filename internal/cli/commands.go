package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/dmawardi/goScaffold/internal/config"
	"github.com/dmawardi/goScaffold/internal/generator"
)

const (
	defaultTemplateDir = "templates/goTemplate"
	defaultOutputDir   = "."
)

func RunCreateCommand() {
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

	// Validate project name format
	if err := validateProjectName(*projectName); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Set default module path if not provided
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
		Force:       *force,
		Verbose:     *verbose,
	}

	// Validate template directory exists
	if _, err := os.Stat(cfg.TemplateDir); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Template directory '%s' does not exist\n", cfg.TemplateDir)
		os.Exit(1)
	}

	// Create output path
	outputPath := filepath.Join(cfg.OutputDir, cfg.ProjectName)

	// Check if output directory exists
	if _, err := os.Stat(outputPath); !os.IsNotExist(err) && !cfg.Force {
		fmt.Fprintf(os.Stderr, "Error: Directory '%s' already exists. Use --force to overwrite\n", outputPath)
		os.Exit(1)
	}

	if cfg.Verbose {
		fmt.Printf("Scaffolding project '%s'...\n", cfg.ProjectName)
		fmt.Printf("Template: %s\n", cfg.TemplateDir)
		fmt.Printf("Output: %s\n", outputPath)
		fmt.Printf("Module: %s\n", cfg.ModulePath)
	}

	// Create the generator and run it
	gen := generator.New(cfg)
	if err := gen.Generate(); err != nil {
		log.Fatalf("Failed to generate project: %v", err)
	}

	fmt.Printf("✅ Project '%s' successfully created at '%s'\n", cfg.ProjectName, outputPath)
	fmt.Println("\nNext steps:")
	fmt.Printf("  cd %s\n", outputPath)
	fmt.Printf("  go mod tidy\n")
	fmt.Printf("  go run cmd/%s/main.go\n", cfg.ProjectName)
}
