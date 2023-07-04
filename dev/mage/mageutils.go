package mageutils

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/bitfield/script"
	"github.com/fatih/color"
	"github.com/magefile/mage/sh"

	fileutils "github.com/l50/goutils/v2/file/fileutils"
	"github.com/l50/goutils/v2/sys"
)

// FuncInfo represents information about an exported function within a Go package.
//
// **Attributes:**
//
// FilePath: A string representing the path to the source file containing the function declaration.
// FuncName: A string representing the name of the exported function.
type FuncInfo struct {
	FilePath string // FilePath is the file path of the source file containing the function declaration.
	FuncName string // FuncName is the name of the exported function.
}

// Compile builds a Go application for a specified operating system and
// architecture. It sets the appropriate environment variables and runs `go
// build`. The compiled application is placed in the specified build path.
//
// **Parameters:**
//
// buildPath: The output directory for the compiled application.
// goOS: The target operating system (e.g., "linux", "darwin", "windows").
// goArch: The target architecture (e.g., "amd64", "arm64").
//
// **Returns:**
//
// error: An error if the compilation process encounters one.
func Compile(buildPath string, goOS string, goArch string) error {
	os.Setenv("GOOS", goOS)
	os.Setenv("GOARCH", goArch)
	err := sh.RunV(
		"go",
		"build",
		"-o",
		buildPath)
	return err
}

// GHRelease creates a new release on GitHub using the given new version.
// It requires the gh CLI tool to be available on the PATH.
//
// **Parameters:**
//
// newVer: A string specifying the new version, e.g., "v1.0.1"
//
// **Returns:**
//
// error: An error if the GHRelease function is not successful.
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

// GoReleaser runs the Goreleaser tool to generate all the supported binaries
// specified in a .goreleaser configuration file.
//
// **Returns:**
//
// error: An error if the Goreleaser function is not successful.
func GoReleaser() error {
	if fileutils.Exists(".goreleaser.yaml") || fileutils.Exists(".goreleaser.yml") {
		if sys.CmdExists("goreleaser") {
			if _, err := script.Exec("goreleaser --snapshot --clean").Stdout(); err != nil {
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

// InstallVSCodeModules installs the modules used by the vscode-go extension in
// Visual Studio Code.
//
// **Returns:**
//
// error: An error if the InstallVSCodeModules function is not successful.
func InstallVSCodeModules() error {
	fmt.Println(color.YellowString("Installing vscode-go dependencies."))
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
		return fmt.Errorf(
			color.RedString("failed to install vscode-go dependencies: %v", err))
	}

	return nil
}

// ModUpdate updates go modules by running 'go get -u' or 'go get -u ./...' if
// recursive is set to true. The function will run in verbose mode if 'v' is
// set to true.
//
// **Parameters:**
//
// recursive: A boolean specifying whether to run the update recursively.
// v: A boolean specifying whether to run the update in verbose mode.
//
// **Returns:**
//
// error: An error if the ModUpdate function is not successful.
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

// Tidy executes 'go mod tidy' to clear the module dependencies.
//
// **Returns:**
//
// error: An error if the Tidy function didn't run successfully.
func Tidy() error {
	if err := sh.Run("go", "mod", "tidy"); err != nil {
		return fmt.Errorf("failed to run `go mod tidy`: %v", err)
	}

	return nil
}

// UpdateMageDeps modifies the dependencies in a given Magefile directory.
// If no directory is provided, it falls back to the 'magefiles' directory.
//
// **Parameters:**
//
// magedir: A string defining the path to the magefiles directory.
//
// **Returns:**
//
// error: An error if the UpdateMageDeps function didn't run successfully.
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

// InstallGoDeps installs the specified Go dependencies by executing 'go install'
// for each dependency.
//
// **Parameters:**
//
// deps: A slice of strings defining the dependencies to install.
//
// **Returns:**
//
// error: An error if the InstallGoDeps function didn't run successfully.
func InstallGoDeps(deps []string) error {
	for _, dep := range deps {
		if _, err := sys.RunCommand("go", "install", dep+"@latest"); err != nil {
			return fmt.Errorf("failed to install input go dependencies: %v", err)
		}
	}

	return nil
}

// FindExportedFunctionsInPackage finds all exported functions in a given Go
// package by parsing all non-test Go files in the package directory. It returns
// a slice of FuncInfo structs. Each contains the file path and the name of an
// exported function. If no exported functions are found in the package, an
// error is returned.
//
// **Parameters:**
//
// pkgPath: A string representing the path to the directory containing the package
// to search for exported functions.
//
// **Returns:**
//
// []FuncInfo: A slice of FuncInfo structs, each containing the file path and the
// name of an exported function found in the package.
// error: An error if no exported functions are found.
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

// FindExportedFuncsWithoutTests discovers all exported functions in a given
// package path that lack corresponding tests.
//
// **Parameters:**
//
// pkgPath: A string defining the package path to search.
//
// **Returns:**
//
// []string: A slice of strings containing the names of exported functions that
// lack corresponding tests.
//
// error: An error if there was a problem parsing the package or finding the tests.
func FindExportedFuncsWithoutTests(pkgPath string) ([]string, error) {
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
