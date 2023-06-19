package file

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitfield/script"
	"github.com/l50/goutils/v2/str"
)

func openFile(file string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(file, flag, perm)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// Append appends an input text string to the end of the specified file.
// If the file does not exist, it will be created.
//
// Parameters:
//
// appendFilePath: A string representing the path to the file.
// text: A string that will be appended to the end of the file.
//
// Returns:
//
// error: An error if the file cannot be opened or the string cannot be written to the file.
//
// Example:
//
// filePath := "/path/to/your/file"
// text := "text to be appended"
// err := Append(filePath, text)
//
//	if err != nil {
//	  log.Fatalf("failed to append text to file: %v", err)
//	}
func Append(appendFilePath string, text string) error {
	f, err := openFile(appendFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(text + "\n"); err != nil {
		return err
	}
	return nil
}

// CreateType represents the type of file creation action to execute.
type CreateType int

const (
	// CreateDirectory represents a directory creation action.
	CreateDirectory CreateType = iota
	// CreateEmptyFile represents an empty file creation action.
	CreateEmptyFile
	// CreateFile represents a file creation action.
	CreateFile
)

// Create creates a directory, an empty file, or a file with content at the specified path,
// depending on the createType argument.
//
// Parameters:
//
// path: A string representing the path to the directory or file.
// contents: A byte slice representing the content to write to the file.
// createType: A CreateType value representing whether to create a directory,
// an empty file, or a file with content.
//
// Returns:
//
// error: An error if the directory or file cannot be created, if it already exists, or
// if there is a problem writing to the file.
//
// Example:
//
// filePath := "/path/to/your/file"
// err := Create(filePath, []byte("file contents"), CreateFile)
//
//	if err != nil {
//	  log.Fatalf("failed to create file: %v", err)
//	}
func Create(path string, contents []byte, createType CreateType) error {
	if Exists(path) {
		return fmt.Errorf("file or directory at path %s already exists", path)
	}
	switch createType {
	case CreateDirectory:
		return createDirectory(path)
	case CreateEmptyFile:
		return createEmptyFile(path)
	case CreateFile:
		return createFile(path, contents)
	default:
		return fmt.Errorf("invalid createType %v", createType)
	}
}

func createDirectory(path string) error {
	if !filepath.IsAbs(path) {
		absDir, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("failed to convert input relative path to an absolute path: %v", err)
		}
		path = absDir
	}

	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("%s already exists", path)
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create new directory at %s: %v", path, err)
	}
	return nil
}

func createEmptyFile(name string) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	return file.Close()
}

func createFile(filePath string, fileContents []byte) error {
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return fmt.Errorf("cannot create dir portion of filepath %s: %v", filePath, err)
	}
	if err := os.WriteFile(filePath, fileContents, os.ModePerm); err != nil {
		return fmt.Errorf("cannot write to file %s: %v", filePath, err)
	}
	return nil
}

// ContainsStr searches for a string in a specified file.
//
// Parameters:
//
// path: A string representing the path to the file.
// searchStr: The string to search for in the file.
//
// Returns:
//
// bool: Returns true if the string is found, otherwise false.
// error: An error if the file cannot be read.
//
// Example:
//
// filePath := "/path/to/your/file"
// searchStr := "text to find"
// found, err := ContainsStr(filePath, searchStr)
//
//	if err != nil {
//	  log.Fatalf("failed to search file: %v", err)
//	}
//
//	if found {
//	  fmt.Printf("'%s' found in file\n", searchStr)
//	}
func ContainsStr(path string, searchStr string) (bool, error) {
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

// CSVToLines reads the contents of a CSV file and returns it as a two-dimensional string slice,
// where each element in the outer slice represents a row in the CSV file,
// and each element in the inner slice represents a value in that row.
// The first row of the CSV file, which is assumed to contain column headers, is skipped.
//
// Parameters:
//
// path: A string representing the path to the CSV file.
//
// Returns:
//
// [][]string: A two-dimensional slice of strings representing the rows and values of the CSV file.
// error: An error if the file cannot be read or parsed.
//
// Example:
//
// csvFilePath := "/path/to/your/csv/file"
// records, err := CSVToLines(csvFilePath)
//
//	if err != nil {
//	  log.Fatalf("failed to read CSV file: %v", err)
//	}
//
//	for _, row := range records {
//	  fmt.Println(row)
//	}
func CSVToLines(path string) ([][]string, error) {
	f, err := os.Open(path)
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

// Delete deletes the specified file.
//
// Parameters:
//
// path: A string representing the path to the file.
//
// Returns:
//
// error: An error if the file cannot be deleted.
//
// Example:
//
// filePath := "/path/to/your/file"
// err := Delete(filePath)
//
//	if err != nil {
//	  log.Fatalf("failed to delete file: %v", err)
//	}
func Delete(path string) error {
	if !Exists(path) {
		return fmt.Errorf("file or directory at path %s does not exist", path)
	}

	if err := os.Remove(path); err != nil {
		return err
	}

	return nil
}

// Exists checks whether a file at the specified path exists.
//
// Parameters:
//
// fileLoc: A string representing the path to the file.
//
// Returns:
//
// bool: Returns true if the file exists, otherwise false.
//
// Example:
//
// filePath := "/path/to/your/file"
// exists := Exists(filePath)
//
//	if !exists {
//	  log.Fatalf("file does not exist")
//	}
func Exists(fileLoc string) bool {
	if _, err := os.Stat(fileLoc); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

// ToSlice reads a file and returns its content as a slice of strings,
// where each element represents a line in the file. Blank lines are omitted.
//
// Parameters:
//
// path: A string representing the path to the file.
//
// Returns:
//
// []string: A slice of strings where each element represents a line in the file.
// error: An error if the file cannot be read.
//
// Example:
//
// filePath := "/path/to/your/file"
// lines, err := ToSlice(filePath)
//
//	if err != nil {
//	  log.Fatalf("failed to read file: %v", err)
//	}
//
//	for _, line := range lines {
//	  fmt.Println(line)
//	}
func ToSlice(path string) ([]string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %v", path, err)
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

// Find searches for a specified filename in a set of directories and returns
// all matches found as a slice of file paths. If no matches are found, it
// returns an error.
//
// Parameters:
//
// fileName: The name of the file to find.
// dirs: A slice of strings representing the directories to search in.
//
// Returns:
//
// []string: A slice of file paths if the file is found.
// error: An error if the file cannot be found.
//
// Example:
//
// fileName := "file_to_find.txt"
// dirs := []string{"/path/to/first/directory", "/path/to/second/directory"}
// filePaths, err := Find(fileName, dirs)
//
//	if err != nil {
//	  log.Fatalf("failed to find file: %v", err)
//	}
//
//	for _, filePath := range filePaths {
//	    fmt.Printf("File found at: %s\n", filePath)
//	}
func Find(fileName string, dirs []string) ([]string, error) {
	var files []string
	for _, dir := range dirs {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(info.Name(), fileName) {
				files = append(files, path)
			}
			return nil
		})

		// Handle potential error from filepath.Walk
		if err != nil {
			return nil, fmt.Errorf("failed to walk directory %v: %v", dir, err)
		}
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("file %v not found in directories", fileName)
	}

	return files, nil
}

// ListR lists all files in a directory and its subdirectories.
//
// Parameters:
//
// dirPath: A string representing the path to the directory.
//
// Returns:
//
// []string: A slice of strings representing the paths of the files found.
// error: An error if the files cannot be listed.
//
// Example:
//
// dirPath := "/path/to/your/directory"
// files, err := ListR(dirPath)
//
//	if err != nil {
//	  log.Fatalf("failed to list files: %v", err)
//	}
//
//	for _, file := range files {
//	  fmt.Println(file)
//	}
func ListR(dirPath string) ([]string, error) {
	result, err := script.FindFiles(dirPath).String()
	if err != nil {
		return []string{}, err
	}

	fileList := str.ToSlice(result, "\n")

	return fileList, nil
}

// Write writes a string to a file.
//
// Parameters:
//
// path: A string representing the path to the file.
// content: The string to write to the file.
//
// Returns:
//
// error: An error if the file cannot be written.
//
// Example:
//
// filePath := "/path/to/your/file"
// content := "text to write to file"
// err := Write(filePath, content)
//
//	if err != nil {
//	  log.Fatalf("failed to write to file: %v", err)
//	}
func Write(path string, content string) error {
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create file at path %s: %v", path, err)
	}
	return nil
}
