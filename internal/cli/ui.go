package cli

import (
	"fmt"
	"os"
)

// PrintUsage prints the usage information for the CLI tool
func PrintUsage() {
	fmt.Printf("Usage: %s <command> [options]\n\n", os.Args[0])
	fmt.Println("A CLI tool to scaffold Go projects from templates.")
	fmt.Println("\nAvailable Commands:")
	fmt.Println("  create    Create a new project from template")
	fmt.Println("  help      Show help information")
	fmt.Println("  version   Show version information")
	fmt.Println("\nFor command-specific help:")
	fmt.Printf("  %s create -h\n", os.Args[0])
	fmt.Println("\nExamples:")
	fmt.Printf("  %s create -name myproject -module github.com/user/myproject\n", os.Args[0])
	fmt.Printf("  %s create -name myapi -output /path/to/projects -module github.com/company/myapi\n", os.Args[0])
}

// printVersion prints the version information for the CLI tool
func PrintVersion() {
	fmt.Println("goScaffold v1.0.0")
	fmt.Println("A CLI tool to scaffold Go projects from templates")
}
