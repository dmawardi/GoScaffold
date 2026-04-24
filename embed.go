package main

import (
	"embed"
	"fmt"
	"io/fs"
)

// Embed templates for moduleName and goTemplate
//
//go:embed templates/moduleName
var moduleTemplateFS embed.FS

//go:embed templates/moduleName_test.go
var moduleTestTemplateFS embed.FS

// getModuleTemplate returns the embedded module template filesystem
func getModuleTemplate() embed.FS {
	return moduleTemplateFS
}

// getModuleTestTemplate returns the embedded module e2e test template filesystem
func getModuleTestTemplate() embed.FS {
	return moduleTestTemplateFS
}

//go:embed templates/goTemplate
var goTemplateFS embed.FS

// getGoTemplate returns the embedded goTemplate filesystem
func getGoTemplate() embed.FS {
	return goTemplateFS
}

// testEmbeds prints what's embedded (for testing purposes)
func testEmbeds() {
	fmt.Println("=== Embedded Templates Test ===")

	// Test moduleName template
	fmt.Println("\nmoduleName template contents:")
	fs.WalkDir(moduleTemplateFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			fmt.Printf("  DIR:  %s\n", path)
		} else {
			fmt.Printf("  FILE: %s\n", path)
		}
		return nil
	})

	// Test goTemplate
	fmt.Println("\ngoTemplate contents:")
	fs.WalkDir(goTemplateFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			fmt.Printf("  DIR:  %s\n", path)
		} else {
			fmt.Printf("  FILE: %s\n", path)
		}
		return nil
	})
}
