package docs

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/spf13/afero"
)

// PackageDoc holds the documentation for a Go package.
//
// **Attributes:**
//
// PackageName: The package name.
// Functions:   A slice of FunctionDoc instances representing the functions.
// GoGetPath:   The 'go get' path for the package.
type PackageDoc struct {
	PackageName string
	Functions   []FunctionDoc
	GoGetPath   string
}

// Repo represents a GitHub repository.
//
// **Attributes:**
//
// Owner: The repository owner's name.
// Name:  The repository's name.
type Repo struct {
	Owner string
	Name  string
}

// FunctionDoc contains the documentation for a function within a Go package.
//
// **Attributes:**
//
// Name:        The function name.
// Signature:   The function signature, including parameters and return types.
// Description: The documentation or description of the function.
// Params:      The function parameters.
type FunctionDoc struct {
	Name        string
	Signature   string
	Description string
	Params      string
}

// FuncInfo holds information about an exported function within a Go package.
//
// **Attributes:**
//
// FilePath: The path to the source file with the function declaration.
// FuncName: The name of the exported function.
type FuncInfo struct {
	FilePath string
	FuncName string
}

// CreatePackageDocs generates package documentation for a Go project using
// a specified template file. It first checks if the template file exists in
// the filesystem denoted by a provided afero.Fs instance. If it exists, the
// function walks the project directory, excluding any specified packages,
// and applies the template to each non-excluded package to generate its
// documentation.
//
// **Parameters:**
//
// fs: An afero.Fs instance representing the filesystem.
//
// repo: A Repo instance containing the Go project's repository details.
//
// templatePath: A string representing the path to the template file to be
// used for generating the package documentation.
//
// excludedPackages: Zero or more strings representing the names of packages
// to be excluded from documentation generation.
//
// **Returns:**
//
// error: An error, if it encounters an issue while checking if the template
// file exists, walking the project directory, or generating the package
// documentation.
func CreatePackageDocs(fs afero.Fs, repo Repo, templatePath string, excludedPackages ...string) error {
	excludedPackagesMap := make(map[string]struct{})
	for _, pkg := range excludedPackages {
		excludedPackagesMap[pkg] = struct{}{}
	}

	exists, err := afero.Exists(fs, templatePath)
	if err != nil {
		return fmt.Errorf("error checking if template file exists: %w", err)
	}
	if !exists {
		return fmt.Errorf("template file does not exist")
	}

	err = afero.Walk(fs, ".", handleDirectory(fs, repo, templatePath, excludedPackagesMap))
	if err != nil {
		return fmt.Errorf("error walking directories: %w", err)
	}

	return nil
}

// generateReadmeFromTemplate generates a README.md file for a Go package using
// a specified template file. It first checks if the template file exists in
// the filesystem denoted by a provided afero.Fs instance. If it exists, the
// function reads its contents, parses it as a template, and applies it to the
// provided PackageDoc to generate the README.md content.
//
// **Parameters:**
//
// fs: An afero.Fs instance for mocking the filesystem for testing.
//
// pkgDoc: A pointer to a PackageDoc instance containing the Go package
// documentation that will be used to generate the README.md file.
//
// path: A string representing the path where the README.md file should be
// created.
//
// templatePath: A string representing the path to the template file to be
// used for generating the README.md file.
//
// **Returns:**
//
// error: An error, if it encounters an issue while checking if the template
// file exists, reading the template file, parsing the template, creating the
// README.md file, or writing to the README.md file.
func generateReadmeFromTemplate(fs afero.Fs, pkgDoc *PackageDoc, path string, templatePath string) error {
	// Check if the template file exists
	exists, err := afero.Exists(fs, templatePath)
	if err != nil {
		return fmt.Errorf("error checking if template file exists: %w", err)
	}
	if !exists {
		return fmt.Errorf("template file does not exist")
	}

	// Open the template file
	templateFile, err := fs.Open(templatePath)
	if err != nil {
		return fmt.Errorf("error opening template file: %w", err)
	}
	defer templateFile.Close()

	// Read the contents of the file into a string
	templateBytes, err := afero.ReadAll(templateFile)
	if err != nil {
		return fmt.Errorf("error reading template file: %w", err)
	}

	// Parse the template file
	tmpl, err := template.New("").Parse(string(templateBytes))
	if err != nil {
		return err
	}

	// Open the output file
	out, err := fs.Create(path)
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

	// Replace hard tabs with spaces
	readmeContent = strings.ReplaceAll(readmeContent, "\t", "    ")

	// Write the modified content to the README file
	if _, err := out.WriteString(readmeContent); err != nil {
		return err
	}

	return nil
}

func loadIgnoreList(fs afero.Fs, ignoreFilePath string) (map[string]struct{}, error) {
	ignoreList := make(map[string]struct{})

	ignoreFile, err := fs.Open(ignoreFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return ignoreList, nil // If the ignore file doesn't exist, just return the empty list.
		}
		return nil, err
	}
	defer ignoreFile.Close()

	scanner := bufio.NewScanner(ignoreFile)
	for scanner.Scan() {
		ignorePath := scanner.Text()
		ignorePath = filepath.Clean(ignorePath)
		ignoreList[ignorePath] = struct{}{}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return ignoreList, nil
}

func handleDirectory(fs afero.Fs, repo Repo, templatePath string, excludedPackagesMap map[string]struct{}) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		ignoreList, err := loadIgnoreList(fs, ".docgenignore")
		if err != nil {
			return fmt.Errorf("error loading ignore list: %w", err)
		}

		_, ignored := ignoreList[filepath.Clean(path)]
		if ignored {
			return filepath.SkipDir
		}

		hasGoFiles, err := directoryContainsGoFiles(fs, path)
		if err != nil {
			return err
		}

		if !hasGoFiles {
			return nil
		}

		return processGoFiles(fs, path, repo, templatePath, excludedPackagesMap)
	}
}

func directoryContainsGoFiles(fs afero.Fs, path string) (bool, error) {
	entries, err := afero.ReadDir(fs, path)
	if err != nil {
		return false, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), ".go") &&
			!strings.HasSuffix(entry.Name(), "_test.go") &&
			!strings.HasSuffix(entry.Name(), "magefile.go") {
			return true, nil
		}
	}

	return false, nil
}

func processGoFiles(fs afero.Fs, path string, repo Repo, templatePath string, excludedPackagesMap map[string]struct{}) error {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path, nonTestFilter, parser.ParseComments)
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {
		// check if the package name is in the excluded packages list
		if _, exists := excludedPackagesMap[pkg.Name]; exists {
			continue // skip this package
		}

		err := generateReadmeForPackage(fs, path, fset, pkg, repo, templatePath)
		if err != nil {
			return err
		}
	}

	return nil
}

func nonTestFilter(info os.FileInfo) bool {
	return !strings.HasSuffix(info.Name(), "_test.go")
}

func generateReadmeForPackage(fs afero.Fs, path string, fset *token.FileSet, pkg *ast.Package, repo Repo, templatePath string) error {
	pkgDoc := &PackageDoc{
		PackageName: pkg.Name,
		GoGetPath:   fmt.Sprintf("github.com/%s/%s/%s", repo.Name, repo.Owner, pkg.Name),
		Functions:   []FunctionDoc{},
	}

	for _, file := range pkg.Files {
		err := processFileDeclarations(fset, pkgDoc, file)
		if err != nil {
			return err
		}
	}

	// Sort function docs by function name
	sort.Slice(pkgDoc.Functions, func(i, j int) bool {
		return pkgDoc.Functions[i].Name < pkgDoc.Functions[j].Name
	})

	return generateReadmeFromTemplate(fs, pkgDoc, filepath.Join(path, "README.md"), templatePath)

}

func processFileDeclarations(fset *token.FileSet, pkgDoc *PackageDoc, file *ast.File) error {
	for _, decl := range file.Decls {
		if fn, isFn := decl.(*ast.FuncDecl); isFn {
			if !fn.Name.IsExported() || strings.HasPrefix(fn.Name.Name, "Test") {
				continue
			}

			fnDoc, err := createFunctionDoc(fset, fn)
			if err != nil {
				return err
			}

			pkgDoc.Functions = append(pkgDoc.Functions, fnDoc)
		}
	}
	return nil
}

func createFunctionDoc(fset *token.FileSet, fn *ast.FuncDecl) (FunctionDoc, error) {
	var params, results string
	if fn.Type.Params != nil {
		params = formatNode(fset, fn.Type.Params)
	}
	if fn.Type.Results != nil {
		results = formatNode(fset, fn.Type.Results)
	}

	signature := fmt.Sprintf("%s(%s) %s", fn.Name.Name, params, results)
	signature = strings.TrimRight(signature, " ") // Trim trailing space

	// Split the signature if it's too long
	const maxLineLength = 80
	if len(signature) > maxLineLength {
		signature = splitLongSignature(signature, maxLineLength)
	}

	return FunctionDoc{
		Name:        fn.Name.Name,
		Signature:   signature,
		Description: fn.Doc.Text(),
		Params:      params,
	}, nil
}

func splitLongSignature(signature string, maxLineLength int) string {
	parts := strings.Split(signature, ",")
	for i := 1; i < len(parts); i++ {
		if len(parts[i-1]) > maxLineLength {
			parts[i-1] = strings.TrimRight(parts[i-1], " ") + ","
			parts[i] = "\n" + strings.TrimLeft(parts[i], " ")
		}
	}
	return strings.Join(parts, "")
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
