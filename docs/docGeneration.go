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
	PackageName string        // The package name.
	Functions   []FunctionDoc // A slice of FunctionDoc representing the functions.
	GoGetPath   string        // The 'go get' path for the package.
}

// Repo represents a GitHub repository.
//
// **Attributes:**
//
// Owner: The repository owner's name.
// Name:  The repository's name.
type Repo struct {
	Owner string // The repository owner's name.
	Name  string // The repository's name.
}

// FunctionDoc contains the documentation for a function within a Go package.
//
// **Attributes:**
//
// Name:        The function name.
// Signature:   The function signature, including parameters and return types.
// Description: The documentation or description of the function.
type FunctionDoc struct {
	Name        string // The function name.
	Signature   string // The function signature, including parameters and return types.
	Description string // The documentation or description of the function.
}

// FuncInfo holds information about an exported function within a Go package.
//
// **Attributes:**
//
// FilePath: The path to the source file with the function declaration.
// FuncName: The name of the exported function.
type FuncInfo struct {
	FilePath string // The path to the source file with the function declaration.
	FuncName string // The name of the exported function.
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

// CreatePackageDocs generates documentation for all Go packages in the current
// directory and its subdirectories. It traverses the file tree using a provided
// afero.Fs and Repo to create a new README.md file in each directory containing
// a Go package. It uses a specified template file for generating the README files.
//
// It will ignore any files or directories listed in the .docgenignore file
// found at the root of the repository. The .docgenignore file should contain
// a list of files and directories to ignore, with each entry on a new line.
//
// **Parameters:**
//
// fs: An afero.Fs instance for mocking the filesystem for testing.
//
// repo: A Repo instance representing the GitHub repository
// containing the Go packages.
//
// templatePath:  A string representing the path to the template file to be
// used for generating README files.
//
// **Returns:**
//
// error: An error, if it encounters an issue while walking the file tree,
// reading a directory, parsing Go files, or generating README.md files.
func CreatePackageDocs(fs afero.Fs, repo Repo, templatePath string) error {
	ignoreList, err := loadIgnoreList(fs, ".docgenignore")
	if err != nil {
		return fmt.Errorf("error loading ignore list: %w", err)
	}

	err = afero.Walk(fs, ".", func(path string, info os.FileInfo, walkErr error) error {
		// Skip hidden directories or files
		if strings.HasPrefix(filepath.Base(path), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip the README file at the root of the repository
		if path == "./README.md" {
			return nil
		}

		absPath, absErr := filepath.Abs(path)
		if absErr != nil {
			return absErr
		}

		relPath, relErr := filepath.Rel(".", absPath)
		if relErr != nil {
			return relErr
		}

		if _, ok := ignoreList[relPath]; ok {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		return handleDirectory(fs, repo, templatePath)(path, info, walkErr)
	})

	if err != nil {
		return fmt.Errorf("error walking directories: %w", err)
	}
	return nil
}

func handleDirectory(fs afero.Fs, repo Repo, templatePath string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		hasGoFiles, err := directoryContainsGoFiles(fs, path)
		if err != nil {
			return err
		}

		if !hasGoFiles {
			return nil
		}

		return processGoFiles(fs, path, repo, templatePath)
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

func processGoFiles(fs afero.Fs, path string, repo Repo, templatePath string) error {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path, nonTestFilter, parser.ParseComments)
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {
		err := generateReadmeForPackage(path, fset, pkg, repo, templatePath)
		if err != nil {
			return err
		}
	}

	return nil
}

func nonTestFilter(info os.FileInfo) bool {
	return !strings.HasSuffix(info.Name(), "_test.go")
}

func generateReadmeForPackage(path string, fset *token.FileSet, pkg *ast.Package, repo Repo, templatePath string) error {
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

	return generateReadmeFromTemplate(pkgDoc, filepath.Join(path, "README.md"), templatePath)
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

func generateReadmeFromTemplate(pkgDoc *PackageDoc, path string, templatePath string) error {
	// Open the template file
	tmpl, err := template.ParseFiles(templatePath)
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

	// Replace hard tabs with spaces
	readmeContent = strings.ReplaceAll(readmeContent, "\t", "    ")

	// Write the modified content to the README file
	if _, err := out.WriteString(readmeContent); err != nil {
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
