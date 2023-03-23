package utils

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bitfield/script"
	"github.com/fatih/color"
)

// AppendToFile appends an input text string to
// the end of the input file.
func AppendToFile(file string, text string) error {
	f, err := os.OpenFile(file,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	// Create channel to grab any errors from the anonymous function below.
	errCh := make(chan error)

	defer func(*os.File) {
		if err := f.Close(); err != nil {
			errCh <- err
		}
	}(f)

	if _, err := f.WriteString(text + "\n"); err != nil {
		return err
	}

	// Check if an error was sent through the channel
	select {
	case err := <-errCh:
		return err
	default:
	}

	return nil
}

// CreateEmptyFile creates an file based on the name input.
// It returns true if the file was created, otherwise it returns false.
func CreateEmptyFile(name string) bool {
	file, err := os.Create(name)
	if err != nil {
		return false
	}

	file.Close()

	return true
}

// CreateFile creates a file at the input filePath
// with the specified fileContents.
func CreateFile(filePath string, fileContents []byte) error {
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return fmt.Errorf("cannot create dir portion"+
			"of filepath %s: %v", filePath, err)
	}

	if err := os.WriteFile(filePath, fileContents, os.ModePerm); err != nil {
		return fmt.Errorf("cannot write to file %s: %v",
			filePath, err)
	}

	return nil
}

// CreateDirectory creates a directory at the input path.
// If any part of the input path doesn't exist, create it.
// Return an error if the path already exists.
func CreateDirectory(path string) error {
	// Check if the input path is absolute
	if !filepath.IsAbs(path) {
		// If the input path is relative, attempt to convert it to an absolute path.
		absDir, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf(color.RedString("failed to convert input "+
				"relative path to an absolute path: %v", err))
		}
		path = absDir
	}

	// Check if the directory already exists
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf(color.RedString("%s already exists", path))
	}

	// Create the input directory if we've gotten here successfully
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf(color.RedString(
			"failed to create new directory at %s: %v", path, err))
	}

	return nil
}

// CSVToLines reads the contents of the specified CSV file and returns
// its contents as a two-dimensional string slice, where each element
// in the outer slice represents a row in the CSV file, and each element
// in the inner slice represents a value in that row. The first row in
// the CSV file is skipped, as it is assumed to contain column headers.
// If the file cannot be read or parsed, an error is returned.
func CSVToLines(filename string) ([][]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return [][]string{}, err
	}

	// close the file at the end of the function call
	defer f.Close()

	r := csv.NewReader(f)
	// skip first line
	if _, err := r.Read(); err != nil {
		return [][]string{}, err
	}

	records, err := r.ReadAll()

	if err != nil {
		return [][]string{}, err
	}

	return records, nil
}

// DeleteFile deletes the input file
func DeleteFile(file string) error {
	if err := os.Remove(file); err != nil {
		return err
	}

	return nil
}

// FileExists will return true if a file specified with fileLoc
// exists. If the file does not exist, it returns false.
func FileExists(fileLoc string) bool {
	_, err := os.Stat(fileLoc)
	if err != nil {
		// `fileLoc` does not exist
		if os.IsNotExist(err) {
			return false
		}
		panic(fmt.Sprintf(
			"failed to check for the existence of %s: %v", fileLoc, err))
	}

	return true
}

// FileToSlice reads an input file into a slice,
// removes blank strings, and returns it.
func FileToSlice(fileName string) ([]string, error) {
	b, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %v", fileName, err)
	}

	lines := strings.Split(string(b), "\n")
	filteredLines := make([]string, 0, len(lines))
	for _, line := range lines {
		if len(strings.TrimSpace(line)) > 0 {
			filteredLines = append(filteredLines, line)
		}
	}

	return filteredLines, nil
}

// FindFile looks for an input `filename` in the specified
// set of `dirs`. The filepath is returned if the `filename` is found.
func FindFile(fileName string, dirs []string) (string, error) {
	for _, d := range dirs {
		files, err := ListFilesR(d)
		if err != nil {
			return "", err
		}
		for _, f := range files {
			fileReg := `/` + fileName + `$`
			m, err := regexp.MatchString(fileReg, f)
			if err != nil {
				return "", fmt.Errorf(
					color.RedString("error - failed to locate %f: %v", fileReg, err))
			} else if m {
				return f, nil
			}
		}
	}
	return "", nil
}

// ListFilesR lists the files found recursively
// from the input `path`.
func ListFilesR(path string) ([]string, error) {
	result, err := script.FindFiles(path).String()
	if err != nil {
		return []string{}, err
	}

	fileList := StringToSlice(result, "\n")

	return fileList, nil
}

// StringInFile searches for input searchStr in
// input the input filepath.
func StringInFile(path string, searchStr string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}

	// Create channel to grab any errors from the anonymous function below.
	errCh := make(chan error)

	defer func(*os.File) {
		if err := f.Close(); err != nil {
			errCh <- err
		}
	}(f)

	scanner := bufio.NewScanner(f)

	line := 1
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), searchStr) {
			return true, nil
		}
		line++
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	// Check if an error was sent through the channel
	select {
	case err := <-errCh:
		return false, err
	default:
	}

	return false, nil
}

// RmRf removes an input path and everything in it.
// If the input path doesn't exist, an error is returned.
func RmRf(path string) error {
	if _, err := os.Stat(path); err == nil {
		if info, err := os.Stat(path); err == nil {
			if info.IsDir() {
				if err := os.RemoveAll(path); err != nil {
					return fmt.Errorf("failed to run RmRf on %s: %v", path, err)
				}
			} else {
				if err := os.Remove(path); err != nil {
					return fmt.Errorf("failed to run RmRf on %s: %v", path, err)
				}
			}
		} else {
			return fmt.Errorf("failed to os.Stat on %s: %v", path, err)
		}
	} else {
		return fmt.Errorf("failed to os.Stat on %s: %v", path, err)
	}

	return nil
}

// FindExportedFuncsWithoutTests finds all exported functions in a given package path that do not have
// corresponding tests. It returns a slice of function names or an error if there is a problem parsing
// the package or finding the tests.
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
