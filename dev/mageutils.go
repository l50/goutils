package dev

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/bitfield/script"
	"github.com/magefile/mage/sh"

	fileutils "github.com/l50/goutils/file"
	"github.com/l50/goutils/sys"
)

// GHRelease creates a new release on GitHub with the given new version.
// It requires that the gh CLI tool is available on the PATH.
//
// Parameters:
//
// newVer: A string specifying the new version, e.g., "v1.0.1"
//
// Returns:
//
// error: An error if the GHRelease function was not successful.
//
// Example:
//
// newVer := "v1.0.1"
// err := GHRelease(newVer)
//
//	if err != nil {
//	  log.Fatalf("failed to create new GH release: %v", err)
//	}
func GHRelease(newVer string) error {
	cmd := "gh"
	if !sys.CmdExists("gh") {
		return fmt.Errorf("required cmd %s not found in $PATH", cmd)
	}

	cl := "CHANGELOG.md"
	// Generate CHANGELOG
	if err := sh.RunV("gh", "changelog", "new", "--next-version", newVer); err != nil {
		return fmt.Errorf("failed to create changelog for new release %s: %v", newVer, err)
	}

	// Create release using CHANGELOG
	if err := sh.RunV("gh", "release", "create", newVer, "-F", cl); err != nil {
		return fmt.Errorf("failed to create new release %s: %v", newVer, err)
	}

	// Remove created CHANGELOG file
	if err := fileutils.Delete(cl); err != nil {
		return fmt.Errorf("failed to delete generated CHANGELOG: %v", err)
	}

	return nil
}

// GoReleaser runs the Goreleaser tool to generate all the supported binaries specified in a .goreleaser configuration file.
//
// Returns:
//
// error: An error if the Goreleaser function was not successful.
//
// Example:
//
// err := GoReleaser()
//
//	if err != nil {
//	  log.Fatalf("failed to run GoReleaser: %v", err)
//	}
func GoReleaser() error {
	if fileutils.Exists(".goreleaser.yaml") || fileutils.Exists(".goreleaser.yml") {
		if sys.CmdExists("goreleaser") {
			if _, err := script.Exec("goreleaser --snapshot --rm-dist").Stdout(); err != nil {
				return fmt.Errorf("failed to run goreleaser: %v", err)
			}
		} else {
			return errors.New("goreleaser not found in $PATH")
		}
	} else {
		return errors.New("no .goreleaser file found")
	}

	return nil
}

// InstallVSCodeModules installs the modules used by the vscode-go extension in Visual Studio Code.
//
// Returns:
//
// error: An error if the InstallVSCodeModules function was not successful.
//
// Example:
//
// err := InstallVSCodeModules()
//
//	if err != nil {
//	  log.Fatalf("failed to install VS Code modules: %v", err)
//	}
func InstallVSCodeModules() error {
	fmt.Println("Installing vscode-go dependencies.")
	vscodeDeps := []string{
		"github.com/uudashr/gopkgs/v2/cmd/gopkgs",
		"github.com/ramya-rao-a/go-outline",
		"github.com/cweill/gotests/gotests",
		"github.com/fatih/gomodifytags",
		"github.com/josharian/impl",
		"github.com/haya14busa/goplay/cmd/goplay",
		"github.com/go-delve/delve/cmd/dlv",
		"honnef.co/go/tools/cmd/staticcheck",
		"golang.org/x/tools/gopls",
		"github.com/rogpeppe/godef",
	}

	if err := InstallGoDeps(vscodeDeps); err != nil {
		return fmt.Errorf("failed to install vscode-go dependencies: %v", err)
	}

	return nil
}

// ModUpdate updates go modules by running 'go get -u' or 'go get -u ./...' if recursive is set to true.
// The function will run in verbose mode if 'v' is set to true.
//
// Parameters:
//
// recursive: A boolean specifying whether to run the update recursively.
// v: A boolean specifying whether to run the update in verbose mode.
//
// Returns:
//
// error: An error if the ModUpdate function was not successful.
//
// Example:
//
// recursive := true
// verbose := true
// err := ModUpdate(recursive, verbose)
//
//	if err != nil {
//	  log.Fatalf("failed to update modules: %v", err)
//	}
func ModUpdate(recursive bool, v bool) error {
	verbose := ""
	if v {
		verbose = "-v"
	}

	if recursive {
		if err := sh.Run("go", "get", "-u", verbose, "./..."); err != nil {
			return fmt.Errorf("failed to run `go get -u %v ./...`: %v", verbose, err)
		}
	}

	if err := sh.Run("go", "get", "-u", verbose); err != nil {
		return fmt.Errorf("failed to run `go get -u %v`", err)
	}

	return nil
}

// Tidy runs 'go mod tidy' to clean up the module dependencies.
//
// Returns:
//
// error: An error if the Tidy function was not successful.
//
// Example:
//
// err := Tidy()
//
//	if err != nil {
//	  log.Fatalf("failed to tidy modules: %v", err)
//	}
func Tidy() error {
	if err := sh.Run("go", "mod", "tidy"); err != nil {
		return fmt.Errorf("failed to run `go mod tidy`: %v", err)
	}

	return nil
}

// UpdateMageDeps updates the dependencies in a specified Magefile directory.
// If no directory is provided, it defaults to the 'magefiles' directory.
//
// Parameters:
//
// magedir: A string specifying the path to the magefiles directory.
//
// Returns:
//
// error: An error if the UpdateMageDeps function was not successful.
//
// Example:
//
// magedir := "custom/mage/dir"
// err := UpdateMageDeps(magedir)
//
//	if err != nil {
//	  log.Fatalf("failed to update Mage dependencies: %v", err)
//	}
func UpdateMageDeps(magedir string) error {
	if magedir == "" {
		magedir = "magefiles"
	}

	cwd := sys.Gwd()
	recursive := false
	verbose := false

	if err := sys.Cd(magedir); err != nil {
		return fmt.Errorf("failed to cd from %s to %s: %v", cwd, magedir, err)
	}

	if err := ModUpdate(recursive, verbose); err != nil {
		return fmt.Errorf("failed to update mage dependencies in %s: %v", magedir, err)
	}

	if err := Tidy(); err != nil {
		return fmt.Errorf("failed to update mage dependencies in %s: %v", magedir, err)
	}

	if err := sys.Cd(cwd); err != nil {
		return fmt.Errorf("failed to cd from %s to %s: %v", magedir, cwd, err)
	}

	return nil
}

// InstallGoDeps installs the specified Go dependencies by running 'go install' for each dependency.
//
// Parameters:
//
// deps: A slice of strings specifying the dependencies to install.
//
// Returns:
//
// error: An error if the InstallGoDeps function was not successful.
//
// Example:
//
// deps := []string{"github.com/stretchr/testify", "github.com/go-chi/chi"}
// err := InstallGoDeps(deps)
//
//	if err != nil {
//	  log.Fatalf("failed to install Go dependencies: %v", err)
//	}
func InstallGoDeps(deps []string) error {
	var err error
	failed := false

	for _, dep := range deps {
		if err := sh.RunV("go", "install", dep+"@latest"); err != nil {
			failed = true
		}
	}

	if failed {
		return fmt.Errorf("failed to install input go dependencies: %w", err)
	}

	return nil
}

// FuncInfo contains information about an exported function, including the file path and function name.
type FuncInfo struct {
	// FilePath is the file path of the source file containing the function declaration.
	FilePath string
	// FuncName is the name of the exported function.
	FuncName string
}

// FindExportedFunctionsInPackage finds all exported functions in a given Go package by recursively parsing all non-test
// Go files in the package directory and returning a slice of FuncInfo structs, each containing the file path and the
// name of an exported function. If no exported functions are found in the package, an error is returned.
//
// Parameters:
//
// pkgPath: A string representing the path to the directory containing the package to search for exported functions.
//
// Returns:
//
// []FuncInfo: A slice of FuncInfo structs, each containing the file path and the name of an exported function found in the package.
// error: An error if no exported functions are found.
//
// Example:
//
// packagePath := "/path/to/your/go/package"
// funcs, err := FindExportedFunctionsInPackage(packagePath)
//
//	if err != nil {
//		 log.Fatalf("failed to find exported functions: %v", err)
//	}
//
//	for _, f := range funcs {
//		 log.Printf("Exported function %s found in file %s\n", f.Name, f.FilePath)
//	}
func FindExportedFunctionsInPackage(pkgPath string) ([]FuncInfo, error) {
	var funcs []FuncInfo

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, pkgPath, func(info os.FileInfo) bool {
		return !strings.HasSuffix(info.Name(), "_test.go")
	}, parser.AllErrors)

	if err != nil {
		return nil, fmt.Errorf("failed to parse directory %s: %w", pkgPath, err)
	}

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				funcDecl, ok := decl.(*ast.FuncDecl)
				if !ok || funcDecl.Recv != nil || !funcDecl.Name.IsExported() {
					continue
				}
				info := FuncInfo{
					FilePath: fset.Position(file.Package).Filename,
					FuncName: funcDecl.Name.Name,
				}
				funcs = append(funcs, info)
			}
		}
	}

	if len(funcs) == 0 {
		return nil, errors.New("no exported functions found in package")
	}

	return funcs, nil
}

// FindExportedFuncsWithoutTests finds all exported functions in a given package path that do not have corresponding tests.
//
// Parameters:
//
// pkgPath: A string specifying the package path to search.
//
// Returns:
//
// []string: A slice of strings containing the names of exported functions that do not have corresponding tests.
//
// error: An error if there was a problem parsing the package or finding the tests.
//
// Example:
//
// funcs, err := FindExportedFuncsWithoutTests("github.com/myorg/mypackage")
//
//	if err != nil {
//	  log.Fatalf("failed to find exported functions without tests: %v", err)
//	}
//
//	for _, funcName := range funcs {
//	  fmt.Println(funcName)
//	}
func FindExportedFuncsWithoutTests(pkgPath string) ([]string, error) {
	// Find all exported functions in the package
	funcs, err := FindExportedFunctionsInPackage(pkgPath)
	if err != nil {
		return nil, err
	}

	// Find all exported functions with corresponding tests
	testFuncs, err := findTestFunctions(pkgPath)
	if err != nil {
		return nil, err
	}

	testableFuncs := make(map[string]bool)
	for _, tf := range testFuncs {
		if strings.HasPrefix(tf, "Test") {
			testableFuncs[tf[4:]] = true
		}
	}

	// Find all exported functions without tests
	exportedFuncsNoTest := make([]string, 0)
	for _, f := range funcs {
		if !testableFuncs[f.FuncName] {
			exportedFuncsNoTest = append(exportedFuncsNoTest, f.FuncName)
		}
	}

	return exportedFuncsNoTest, nil
}

func findTestFunctions(pkgPath string) ([]string, error) {
	var testFuncs []string

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, pkgPath, func(info os.FileInfo) bool {
		return strings.HasSuffix(info.Name(), "_test.go")
	}, parser.AllErrors)

	if err != nil {
		return nil, fmt.Errorf("failed to parse directory %s: %w", pkgPath, err)
	}

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				funcDecl, ok := decl.(*ast.FuncDecl)
				if !ok {
					continue
				}
				testFuncs = append(testFuncs, funcDecl.Name.Name)
			}
		}
	}

	return testFuncs, nil
}
