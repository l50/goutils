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
	"github.com/l50/goutils/str"
)

// Append appends an input text string to the end of the specified file.
// If the file does not exist, it will be created.
//
// Parameters:
//
// file: A string representing the path to the file.
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

// CreateEmpty creates an empty file at the specified path.
// If a file with the same name already exists, it will be overwritten.
//
// Parameters:
//
// name: A string representing the path to the file.
//
// Returns:
//
// bool: Returns true if the file was created successfully, otherwise false.
//
// Example:
//
// filePath := "/path/to/your/file"
// success := CreateEmpty(filePath)
//
//	if !success {
//	  log.Fatalf("failed to create empty file")
//	}
func CreateEmpty(name string) bool {
	file, err := os.Create(name)
	if err != nil {
		return false
	}

	file.Close()

	return true
}

// Create creates a file at the specified path with the provided content.
// If the file does not exist, it will be created.
//
// Parameters:
//
// filePath: A string representing the path to the file.
// fileContents: A byte slice containing the content to be written to the file.
//
// Returns:
//
// error: An error if the file cannot be created or the content cannot be written to the file.
//
// Example:
//
// filePath := "/path/to/your/file"
// content := []byte("content to be written")
// err := Create(filePath, content)
//
//	if err != nil {
//	  log.Fatalf("failed to create file: %v", err)
//	}
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

// CreateDirectory creates a directory at the specified path.
// If the directory already exists, it returns an error.
//
// Parameters:
//
// path: A string representing the path to the directory.
//
// Returns:
//
// error: An error if the directory cannot be created or if it already exists.
//
// Example:
//
// dirPath := "/path/to/your/directory"
// err := CreateDirectory(dirPath)
//
//	if err != nil {
//	  log.Fatalf("failed to create directory: %v", err)
//	}
func CreateDirectory(path string) error {
	// Check if the input path is absolute
	if !filepath.IsAbs(path) {
		// If the input path is relative, attempt to convert it to an absolute path.
		absDir, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("failed to convert input "+
				"relative path to an absolute path: %v", err)
		}
		path = absDir
	}

	// Check if the directory already exists
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("%s already exists", path)
	}

	// Create the input directory if we've gotten here successfully
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf(
			"failed to create new directory at %s: %v", path, err)
	}

	return nil
}

// CSVToLines reads the contents of a CSV file and returns it as a two-dimensional string slice,
// where each element in the outer slice represents a row in the CSV file,
// and each element in the inner slice represents a value in that row.
// The first row of the CSV file, which is assumed to contain column headers, is skipped.
//
// Parameters:
//
// filename: A string representing the path to the CSV file.
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

// Delete deletes the specified file.
//
// Parameters:
//
// file: A string representing the path to the file.
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
func Delete(file string) error {
	if err := os.Remove(file); err != nil {
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
// fileName: A string representing the path to the file.
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

// Find searches for a specified filename in a set of directories.
// Returns the file path if found, or an error if the file cannot be found.
//
// Parameters:
//
// fileName: The name of the file to find.
// dirs: A slice of strings representing the directories to search in.
//
// Returns:
//
// string: The file path if the file is found.
// error: An error if the file cannot be found.
//
// Example:
//
// fileName := "file_to_find.txt"
// dirs := []string{"/path/to/first/directory", "/path/to/second/directory"}
// filePath, err := Find(fileName, dirs)
//
//	if err != nil {
//	  log.Fatalf("failed to find file: %v", err)
//	}
//
// fmt.Printf("File found at: %s\n", filePath)
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
				return "", fmt.Errorf("error - failed to locate %s: %v", fileReg, err)
			} else if m {
				return f, nil
			}
		}
	}
	return "", nil
}

// ListR lists all files in a directory and its subdirectories.
//
// Parameters:
//
// path: A string representing the path to the directory.
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
func ListR(path string) ([]string, error) {
	result, err := script.FindFiles(path).String()
	if err != nil {
		return []string{}, err
	}

	fileList := str.ToSlice(result, "\n")

	return fileList, nil
}

// FindStr searches for a string in a specified file.
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
// found, err := FindStr(filePath, searchStr)
//
//	if err != nil {
//	  log.Fatalf("failed to search file: %v", err)
//	}
//
//	if found {
//	  fmt.Printf("'%s' found in file\n", searchStr)
//	}
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
