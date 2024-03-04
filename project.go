package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func createProject(projectName string) error {
	// Convert projectName to CamelCase for the module name
	moduleName := toCamelCase(projectName)

	// Create the project directory and the src subdirectory
	err := os.MkdirAll(filepath.Join(projectName, "src"), 0755)
	if err != nil {
		return err
	}

	// Create the rnl.rnl file with properties and dependencies
	rnlContent := fmt.Sprintf("let properties = {\n    name: \"%s\"\n    moduleName: \"%s\"\n}\n\nlet dependencies = [\n\n]\n", projectName, moduleName)
	err = os.WriteFile(filepath.Join(projectName, "rnl.rnl"), []byte(rnlContent), fs.FileMode(0644))
	if err != nil {
		return err
	}

	// Create the main.rnl file with a module declaration and a hello world main function
	mainContent := fmt.Sprintf("module %s\n\nfn main() {\n    print(\"Hello, world!\")\n}\n", moduleName)
	err = os.WriteFile(filepath.Join(projectName, "src", "main.rnl"), []byte(mainContent), fs.FileMode(0644))
	if err != nil {
		return err
	}

	return nil
}

func toCamelCase(input string) string {
	titleSpace := strings.Title(strings.Replace(input, "_", " ", -1))
	return strings.Replace(titleSpace, " ", "", -1)
}
