package file

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bitfield/script"
	"github.com/fatih/color"
	"github.com/l50/goutils/v2/str"
)

// Append appends an input text string to
// the end of the input fileutils.
func Append(file string, text string) error {
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

// CreateEmpty creates an file based on the name input.
// It returns true if the file was created, otherwise it returns false.
func CreateEmpty(name string) bool {
	file, err := os.Create(name)
	if err != nil {
		return false
	}

	file.Close()

	return true
}

// Create creates a file at the input filePath
// with the specified fileContents.
func Create(filePath string, fileContents []byte) error {
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

// Delete deletes the input file
func Delete(file string) error {
	if err := os.Remove(file); err != nil {
		return err
	}

	return nil
}

// Exists will return true if a file specified with fileLoc
// exists. If the file does not exist, it returns false.
func Exists(fileLoc string) bool {
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

// ToSlice reads an input file into a slice,
// removes blank strings, and returns it.
func ToSlice(fileName string) ([]string, error) {
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

// Find looks for an input `filename` in the specified
// set of `dirs`. The filepath is returned if the `filename` is found.
func Find(fileName string, dirs []string) (string, error) {
	for _, d := range dirs {
		files, err := ListR(d)
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

// ListR lists the files found recursively
// from the input `path`.
func ListR(path string) ([]string, error) {
	result, err := script.FindFiles(path).String()
	if err != nil {
		return []string{}, err
	}

	fileList := str.ToSlice(result, "\n")

	return fileList, nil
}

// FindStr searches for input searchStr in
// input the input filepath.
func FindStr(path string, searchStr string) (bool, error) {
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
