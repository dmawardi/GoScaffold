package main

import (
	"fmt"
	"os"

	"github.com/dmawardi/goScaffold/internal/cli"
)

func main() {
	if len(os.Args) < 2 {
		cli.PrintUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "create":
		cli.RunCreateCommand(getGoTemplate()) // Pass the embedded goTemplate filesystem to the create command
	case "help", "-h", "--help":
		cli.PrintUsage()
	case "version", "-v", "--version":
		cli.PrintVersion()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		cli.PrintUsage()
		os.Exit(1)
	}
}
