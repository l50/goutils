package mageutils_test

import (
	"fmt"
	"log"

	mageutils "github.com/l50/goutils/v2/dev/mage"
)

func ExampleCompile() {
	buildPath := "/path/to/output/directory"
	goOS := "linux"
	goArch := "amd64"

	if err := mageutils.Compile(buildPath, goOS, goArch); err != nil {
		log.Fatalf("failed to compile: %v", err)
	}

	fmt.Printf("application compiled successfully at: %s\n", buildPath)
}

func ExampleGHRelease() {
	newVer := "v1.0.1"
	if err := mageutils.GHRelease(newVer); err != nil {
		log.Fatalf("failed to create new GH release: %v", err)
	}
}

// Example GoReleaser
func ExampleGoReleaser() {
	if err := mageutils.GoReleaser(); err != nil {
		log.Fatalf("failed to run GoReleaser: %v", err)
	}
}

func ExampleInstallVSCodeModules() {
	if err := mageutils.InstallVSCodeModules(); err != nil {
		log.Fatalf("failed to install VS Code modules: %v", err)
	}
}

func ExampleModUpdate() {
	recursive := true
	verbose := true

	if err := mageutils.ModUpdate(recursive, verbose); err != nil {
		log.Fatalf("failed to update modules: %v", err)
	}
}

func ExampleTidy() {
	if err := mageutils.Tidy(); err != nil {
		log.Fatalf("failed to tidy modules: %v", err)
	}
}

func ExampleUpdateMageDeps() {
	magedir := "custom/mage/dir"

	if err := mageutils.UpdateMageDeps(magedir); err != nil {
		log.Fatalf("failed to update Mage dependencies: %v", err)
	}
}

func ExampleInstallGoDeps() {
	deps := []string{"github.com/stretchr/testify", "github.com/go-chi/chi"}

	if err := mageutils.InstallGoDeps(deps); err != nil {
		log.Fatalf("failed to install Go dependencies: %v", err)
	}
}

func ExampleFindExportedFunctionsInPackage() {
	packagePath := "/path/to/your/go/package"

	funcs, err := mageutils.FindExportedFunctionsInPackage(packagePath)
	if err != nil {
		log.Fatalf("failed to find exported functions: %v", err)
	}

	for _, f := range funcs {
		log.Printf("Exported function %s found in file %s\n", f.FuncName, f.FilePath)
	}
}

func ExampleFindExportedFuncsWithoutTests() {
	funcs, err := mageutils.FindExportedFuncsWithoutTests("github.com/myorg/mypackage")
	if err != nil {
		log.Fatalf("failed to find exported functions without tests: %v", err)
	}

	for _, funcName := range funcs {
		fmt.Println(funcName)
	}
}
