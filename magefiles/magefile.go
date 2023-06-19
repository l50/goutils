//go:build mage

package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/l50/goutils/v2/dev/lint"
	mageutils "github.com/l50/goutils/v2/dev/mage"

	// mage utility functions
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

func init() {
	os.Setenv("GO111MODULE", "on")
}

// InstallDeps Installs go dependencies
func InstallDeps() error {
	fmt.Println("Installing dependencies.")

	if err := mageutils.Tidy(); err != nil {
		return fmt.Errorf("failed to install dependencies: %v", err)
	}

	if err := lint.InstallGoPCDeps(); err != nil {
		return fmt.Errorf("failed to install pre-commit dependencies: %v", err)
	}

	if err := mageutils.InstallVSCodeModules(); err != nil {
		return fmt.Errorf("failed to install vscode-go modules: %v", err)
	}

	return nil
}

// FindExportedFuncsWithoutTests finds exported functions without tests
func FindExportedFuncsWithoutTests(pkg string) ([]string, error) {
	funcs, err := mageutils.FindExportedFuncsWithoutTests(os.Args[1])

	if err != nil {
		log.Fatalf("failed to find exported functions without tests: %v", err)
	}

	for _, funcName := range funcs {
		fmt.Println(funcName)
	}

	return funcs, nil

}

// PackageDoc represents the documentation for a package.
type PackageDoc struct {
	PackageName string        // PackageName is the name of the package.
	Functions   []FunctionDoc // Functions is a slice of FunctionDoc representing the functions in the package.
	GoGetPath   string        // GoGetPath is the Go get path for the package.
}

// FunctionDoc represents the documentation for a function.
type FunctionDoc struct {
	Name        string // Name is the name of the function.
	Signature   string // Signature is the function signature, including the parameters and return types.
	Description string // Description is the description or documentation of the function.
}

// CreatePackageDocs creates package documentation
func CreatePackageDocs() error {
	return filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}

		// Skip directories without Go files
		files, err := ioutil.ReadDir(path)
		if err != nil {
			return err
		}

		hasGoFiles := false
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if strings.HasSuffix(file.Name(), ".go") &&
				!strings.HasSuffix(file.Name(), "_test.go") &&
				!strings.HasSuffix(file.Name(), "magefile.go") {
				hasGoFiles = true
				break
			}
		}
		if !hasGoFiles {
			return nil
		}

		// Parse the directory
		fset := token.NewFileSet()
		pkgs, err := parser.ParseDir(fset, path, func(info os.FileInfo) bool {
			return !strings.HasSuffix(info.Name(), "_test.go") // Ignore _test.go files
		}, parser.ParseComments)
		if err != nil {
			return err
		}

		// Extract and print the documentation
		for _, pkg := range pkgs {
			// Create the package documentation struct
			pkgDoc := &PackageDoc{
				PackageName: pkg.Name,                                              // Or whatever your package's name is
				GoGetPath:   fmt.Sprintf("github.com/l50/goutils/v2/%s", pkg.Name), // Or however your Go Get path is structured
			}

			for _, file := range pkg.Files {
				for _, decl := range file.Decls {
					if fn, isFn := decl.(*ast.FuncDecl); isFn {
						// Ignore non-exported and test functions
						if !fn.Name.IsExported() || strings.HasPrefix(fn.Name.Name, "Test") {
							continue
						}

						// Create the function documentation struct
						fnDoc := FunctionDoc{
							Name:        fn.Name.Name,
							Signature:   fmt.Sprintf("%s(%s) %s", fn.Name.Name, formatNode(fset, fn.Type.Params), formatNode(fset, fn.Type.Results)),
							Description: fn.Doc.Text(),
						}

						// Append it to the package doc
						pkgDoc.Functions = append(pkgDoc.Functions, fnDoc)
					}
				}
			}

			// Generate README.md from the template
			err = generateReadmeFromTemplate(pkgDoc, filepath.Join(path, "README.md"))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func generateReadmeFromTemplate(pkgDoc *PackageDoc, path string) error {
	// Open the template file
	tmpl, err := template.ParseFiles("magefiles/tmpl/README.md.tmpl")
	if err != nil {
		return err
	}

	// Open the output file
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	// Execute the template with the package documentation
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, pkgDoc)
	if err != nil {
		return err
	}

	// Replace &#34; with "
	readmeContent := strings.ReplaceAll(buf.String(), "&#34;", "\"")

	// Write the modified content to the README file
	_, err = out.WriteString(readmeContent)
	if err != nil {
		return err
	}

	return nil
}

func formatNode(fset *token.FileSet, node interface{}) string {
	switch n := node.(type) {
	case *ast.FieldList:
		return fieldListString(fset, n)
	default:
		var buf bytes.Buffer
		err := printer.Fprint(&buf, fset, node)
		if err != nil {
			return fmt.Sprintf("error printing syntax tree: %v", err)
		}
		return buf.String()
	}
}

func fieldListString(fset *token.FileSet, fieldList *ast.FieldList) string {
	var buf strings.Builder
	for i, field := range fieldList.List {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(formatNode(fset, field.Type))
	}
	return buf.String()
}

// RunPreCommit runs all pre-commit hooks locally
func RunPreCommit() error {
	mg.Deps(InstallDeps)

	fmt.Println("Updating pre-commit hooks.")
	if err := lint.UpdatePCHooks(); err != nil {
		return err
	}

	fmt.Println("Clearing the pre-commit cache to ensure we have a fresh start.")
	if err := lint.ClearPCCache(); err != nil {
		return err
	}

	fmt.Println("Running all pre-commit hooks locally.")
	if err := lint.RunPCHooks(); err != nil {
		return err
	}

	return nil
}

// RunTests runs all of the unit tests
func RunTests() error {
	mg.Deps(InstallDeps)

	fmt.Println("Running unit tests.")
	if err := sh.RunV(filepath.Join(".hooks", "go-unit-tests.sh"), "all"); err != nil {
		return fmt.Errorf("failed to run unit tests: %v", err)
	}

	return nil
}

// UpdateMirror updates pkg.go.goutils with the release associated with the input tag
func UpdateMirror(tag string) error {
	var err error
	fmt.Printf("Updating pkg.go.goutils with the new tag %s.", tag)

	err = sh.RunV("curl", "--silent", fmt.Sprintf(
		"https://sum.golang.org/lookup/github.com/l50/goutils/v2@%s",
		tag))
	if err != nil {
		return fmt.Errorf("failed to update proxy.golang.org: %w", err)
	}

	err = sh.RunV("curl", "--silent", fmt.Sprintf(
		"https://proxy.golang.org/github.com/l50/goutils/v2/@v/%s.info",
		tag))
	if err != nil {
		return fmt.Errorf("failed to update pkg.go.goutils: %w", err)
	}

	return nil
}
