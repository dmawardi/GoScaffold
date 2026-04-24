package generator

import (
	"fmt"
	"strings"
)

// insertImport adds importPath inside the existing import block.
func insertImport(content, importPath string) (string, error) {
	// Find the import block and insert the new import path before the closing parenthesis
	blockStart := strings.Index(content, "import (")
	if blockStart == -1 {
		return "", fmt.Errorf("could not find import block")
	}
	// Find the closing parenthesis of the import block to insert before it
	closingIdx := strings.Index(content[blockStart:], "\n)")
	if closingIdx == -1 {
		return "", fmt.Errorf("could not find end of import block")
	}
	// Insert the new import path before the closing parenthesis, ensuring proper formatting
	insertAt := blockStart + closingIdx
	return content[:insertAt] + "\n\t" + importPath + content[insertAt:], nil
}

// insertStateFuncsEntry appends entry as the last item in the stateFuncs slice.
func insertStateFuncsEntry(content, entry string) (string, error) {
	// Find the start of the stateFuncs slice declaration
	sliceStart := strings.Index(content, "var stateFuncs = []StateFunc{")
	if sliceStart == -1 {
		return "", fmt.Errorf("could not find stateFuncs slice")
	}
	// Find the closing brace of the slice declaration to insert before it
	closingIdx := strings.Index(content[sliceStart:], "\n}")
	if closingIdx == -1 {
		return "", fmt.Errorf("could not find end of stateFuncs slice")
	}
	// Insert the new entry before the closing brace, ensuring proper formatting
	insertAt := sliceStart + closingIdx
	return content[:insertAt] + "\n\t" + entry + content[insertAt:], nil
}

// insertMigrationEntry appends entry as the last item in the modelsToMigrate slice.
func insertMigrationEntry(content, entry string) (string, error) {
	// Find the start of the modelsToMigrate slice declaration
	sliceStart := strings.Index(content, "var modelsToMigrate = []interface{}{")
	if sliceStart == -1 {
		return "", fmt.Errorf("could not find modelsToMigrate slice")
	}
	// Find the closing brace of the slice declaration to insert before it
	closingIdx := strings.Index(content[sliceStart:], "\n}")
	if closingIdx == -1 {
		return "", fmt.Errorf("could not find end of modelsToMigrate slice")
	}
	// Insert the new entry before the closing brace, ensuring proper formatting
	insertAt := sliceStart + closingIdx
	return content[:insertAt] + "\n\t" + entry + content[insertAt:], nil
}

// insertModuleStackField adds a field to the moduleStack struct (in module.go) before its closing brace.
func insertModuleStackField(content, field string) (string, error) {
	// Find the start of the moduleStack struct declaration
	structStart := strings.Index(content, "type moduleStack struct {")
	if structStart == -1 {
		return "", fmt.Errorf("could not find moduleStack struct")
	}
	// Find the closing brace of the struct declaration to insert before it
	closingIdx := strings.Index(content[structStart:], "\n}")
	if closingIdx == -1 {
		return "", fmt.Errorf("could not find end of moduleStack struct")
	}
	// Insert the new field before the closing brace, ensuring proper formatting
	insertAt := structStart + closingIdx
	return content[:insertAt] + "\n\t" + field + content[insertAt:], nil
}

// insertModuleStackInit inserts an initialization block before `return module` in buildModuleStack (in module.go).
func insertModuleStackInit(content, block string) (string, error) {
	returnIdx := strings.Index(content, "\n\treturn module\n}")
	if returnIdx == -1 {
		return "", fmt.Errorf("could not find 'return module' in buildModuleStack")
	}
	return content[:returnIdx] + "\n\t" + block + "\n" + content[returnIdx:], nil
}

func insertRouteEntry(content, entry string) (string, error) { // Find the start of the RegisterRoutes function declaration
	funcStart := strings.Index(content, "func registerAllRoutes(router *gin.Engine, modules *moduleStack) {")
	if funcStart == -1 {
		return "", fmt.Errorf("could not find registerAllRoutes function")
	}
	// Find the closing brace of the function declaration to insert before it
	closingIdx := strings.LastIndex(content, "\n}")
	if closingIdx == -1 {
		return "", fmt.Errorf("could not find end of registerAllRoutes function")
	}
	// Insert the new entry before the closing brace, ensuring proper formatting
	insertAt := closingIdx
	return content[:insertAt] + "\n\t" + entry + content[insertAt:], nil
}
