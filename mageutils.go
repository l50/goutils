package utils

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
)

// GHRelease creates a new release with the input newVer using the gh cli tool.
// Example newVer: v1.0.1
func GHRelease(newVer string) error {
	cmd := "gh"
	if !CmdExists("gh") {
		return errors.New(color.RedString(
			"required cmd %s not found in $PATH: %v", cmd, err))
	}

	cl := "CHANGELOG.md"
	// Generate CHANGELOG
	if err := sh.RunV("gh", "changelog", "new", "--next-version", newVer); err != nil {
		return fmt.Errorf(color.RedString(
			"failed to create changelog for new release %s: %v", newVer, err))
	}

	// Create release using CHANGELOG
	if err := sh.RunV("gh", "release", "create", newVer, "-F", cl); err != nil {
		return fmt.Errorf(color.RedString(
			"failed to create new release %s: %v", newVer, err))
	}

	// Remove created CHANGELOG file
	if err := DeleteFile(cl); err != nil {
		return fmt.Errorf(color.RedString(
			"failed to delete generated CHANGELOG: %v", err))
	}

	return nil
}

// GoReleaser Runs goreleaser to generate all of the supported binaries
// specified in `.goreleaser`.
func GoReleaser() error {
	if FileExists(".goreleaser.yaml") || FileExists(".goreleaser.yml") {
		if CmdExists("goreleaser") {
			if _, err := script.Exec("goreleaser --snapshot --rm-dist").Stdout(); err != nil {
				return fmt.Errorf(color.RedString(
					"failed to run goreleaser: %v", err))
			}
		} else {
			return errors.New(color.RedString(
				"goreleaser not found in $PATH"))
		}
	} else {
		return errors.New(color.RedString(
			"no .goreleaser file found"))
	}

	return nil
}

// InstallVSCodeModules installs the modules used by the vscode-go extension in VSCode.
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

// ModUpdate runs `go get -u` or
// `go get -u ./... if `recursive` is set to true.
// The `v` parameter provides verbose output
// if set to true.
func ModUpdate(recursive bool, v bool) error {
	verbose := ""
	if v {
		verbose = "-v"
	}

	if recursive {
		if err := sh.Run("go", "get", "-u", verbose, "./..."); err != nil {
			return fmt.Errorf(
				color.RedString("failed to run `go get -u %v ./...`: %v", verbose, err))
		}
	}

	if err := sh.Run("go", "get", "-u", verbose); err != nil {
		return fmt.Errorf(
			color.RedString("failed to run `go get -u %v`", err))
	}

	return nil
}

// Tidy runs `go mod tidy`.
func Tidy() error {
	if err := sh.Run("go", "mod", "tidy"); err != nil {
		return fmt.Errorf(
			color.RedString("failed to run `go mod tidy`: %v", err))
	}

	return nil
}

// UpdateMageDeps updates mage-specific dependencies
// using the input path to the associated go.mod.
func UpdateMageDeps(magedir string) error {
	// If no input is provided, default to magefiles.
	// As per the mage docs, the magefiles directory
	// is the default location for mage.
	if magedir == "" {
		magedir = "magefiles"
	}

	cwd := Gwd()
	recursive := false
	verbose := false

	if err := Cd(magedir); err != nil {
		return fmt.Errorf(
			color.RedString(
				"failed to cd from %s to %s: %v", cwd, magedir, err))
	}

	if err := ModUpdate(recursive, verbose); err != nil {
		return fmt.Errorf(
			color.RedString(
				"failed to update mage dependencies in %s: %v", magedir, err))
	}

	if err := Tidy(); err != nil {
		return fmt.Errorf(
			color.RedString(
				"failed to update mage dependencies in %s: %v", magedir, err))
	}

	if err := Cd(cwd); err != nil {
		return fmt.Errorf(
			color.RedString(
				"failed to cd from %s to %s: %v", magedir, cwd, err))
	}

	return nil
}

// InstallGoDeps runs go install for the input dependencies.
func InstallGoDeps(deps []string) error {
	var err error
	failed := false

	for _, dep := range deps {
		if err := sh.RunV("go", "install", dep+"@latest"); err != nil {
			failed = true
		}
	}

	if failed {
		return fmt.Errorf(
			color.RedString("failed to install input go dependencies: %w", err))
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
