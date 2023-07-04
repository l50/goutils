package file

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// File is an interface representing a system file.
//
// **Methods:**
//
// Open: Opens the file, returns a io.ReadCloser and an error.
// Write: Writes contents to the file, returns an error.
// RemoveAll: Removes a file or directory at the specified path, returns an error.
// Stat: Retrieves the FileInfo for the specified file or directory, returns an os.FileInfo and an error.
// Remove: Removes the specified file or directory, returns an error.
type File interface {
	Open() (io.ReadCloser, error)
	Write(contents []byte, perm os.FileMode) error
	RemoveAll() error
	Stat() (os.FileInfo, error)
	Remove() error
}

// RealFile is a concrete implementation of the File interface.
// It's used to operate with actual system files.
type RealFile string

// Open is a method for the RealFile type that opens the file and
// returns a io.ReadCloser and an error.
//
// **Returns:**
//
// io.ReadCloser: An object that allows reading from and closing the file.
// error: An error if any issue occurs while trying to open the file.
func (rf RealFile) Open() (io.ReadCloser, error) {
	return os.Open(string(rf))
}

// RemoveAll is a method for the RealFile type that removes
// a file or directory at the specified path.
// If the path represents a directory, RemoveAll will remove
// the directory and all its content.
//
// **Parameters:**
//
// path: A string representing the path to the file or directory to remove.
//
// **Returns:**
//
// error: An error if any issue occurs while trying to remove the file or directory.
func (rf RealFile) RemoveAll() error {
	return os.RemoveAll(string(rf))
}

// Stat is a method for the RealFile type that retrieves the
// FileInfo for the specified file or directory.
//
// **Parameters:**
//
// name: A string representing the path to the file or directory.
//
// **Returns:**
//
// os.FileInfo: FileInfo describing the named file.
// error: An error if any issue occurs while trying to get the FileInfo.
func (rf RealFile) Stat() (os.FileInfo, error) {
	return os.Stat(string(rf))
}

// Remove is a method for the RealFile type that removes the specified file or directory.
// Note that it will not remove a directory unless it is empty.
//
// **Parameters:**
//
// name: A string representing the path to the file or directory to remove.
//
// **Returns:**
//
// error: An error if any issue occurs while trying to remove the file or directory.
func (rf RealFile) Remove() error {
	return os.Remove(string(rf))
}

// Write is a method for the RealFile type that writes a slice of bytes
// to the file with specified file permissions.
//
// **Parameters:**
//
// contents: A slice of bytes that should be written to the file.
// mode: File permissions to use when creating the file.
//
// **Returns:**
//
// error: An error if any issue occurs while trying to write to the file.
func (rf RealFile) Write(contents []byte, perm os.FileMode) error {
	return os.WriteFile(string(rf), contents, perm)
}

// Append adds a string to the end of a file. If the file
// doesn't exist, it's created with the default permissions.
//
// **Parameters:**
//
// text: String to append to the end of the file.
//
// **Returns:**
//
// error: An error if the file can't be opened or the string can't be
// written to the file.
func (rf RealFile) Append(text string) error {
	rc, err := rf.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	// Read the existing contents of the file
	fileContents, err := io.ReadAll(rc)
	if err != nil {
		return err
	}

	// Append the new text
	fileContents = append(fileContents, []byte(text)...)

	// Get the existing file permissions
	info, err := rf.Stat()
	if err != nil {
		return err
	}
	perm := info.Mode()

	// Write back to the file
	return rf.Write(fileContents, perm)
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

// Create makes a directory, an empty file, or a file with content at
// the specified path, depending on the createType argument.
//
// **Parameters:**
//
// path: Path to the directory or file.
// contents: Content to write to the file as a byte slice.
// createType: A CreateType value representing whether to create a
// directory, an empty file, or a file with content.
//
// **Returns:**
//
// error: An error if the directory or file can't be created, if it
// already exists, or if there's a problem writing to the file.
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

// HasStr checks for the presence of a string in a specified file.
//
// **Parameters:**
//
// path: String representing the path to the file.
// searchStr: String to look for in the file.
//
// **Returns:**
//
// bool: Returns true if the string is found, otherwise false.
// error: An error if the file cannot be read.
func HasStr(path string, searchStr string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), searchStr) {
			return true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
}

// CSVToLines reads a CSV file and returns it as a 2D string slice. Each
// element in the outer slice represents a row in the CSV, each element in the
// inner slice represents a value in that row. The first row of the CSV,
// assumed to contain column headers, is skipped.
//
// **Parameters:**
//
// path: String representing the path to the CSV file.
//
// **Returns:**
//
// [][]string: 2D slice of strings representing the rows and values of the CSV.
// error: An error if the file cannot be read or parsed.
func CSVToLines(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return [][]string{}, err
	}
	defer file.Close()

	r := csv.NewReader(file)

	// Skip header line
	_, err = r.Read()
	if err != nil {
		return [][]string{}, err
	}

	records, err := r.ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return records, nil
}

// Delete removes the specified file.
//
// **Parameters:**
//
// path: String representing the path to the file.
//
// **Returns:**
//
// error: An error if the file cannot be deleted.
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
// **Parameters:**
//
// fileLoc: String representing the path to the file.
//
// **Returns:**
//
// bool: Returns true if the file exists, otherwise false.
func Exists(fileLoc string) bool {
	if _, err := os.Stat(fileLoc); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// ToSlice reads a file and returns its content as a slice of strings, each
// element represents a line in the file. Blank lines are omitted.
//
// **Parameters:**
//
// path: String representing the path to the file.
//
// **Returns:**
//
// []string: Slice of strings where each element represents a line in the file.
// error: An error if the file cannot be read.
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
// **Parameters:**
//
// fileName: Name of the file to find.
// dirs: Slice of strings representing the directories to search in.
//
// **Returns:**
//
// []string: Slice of file paths if the file is found.
// error: An error if the file cannot be found.
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
// **Parameters:**
//
// dirPath: String representing the path to the directory.
//
// **Returns:**
//
// []string: Slice of strings representing the paths of the files found.
// error: An error if the files cannot be listed.
func ListR(dirPath string) ([]string, error) {
	fis, err := os.ReadDir(dirPath)
	if err != nil {
		return []string{}, err
	}
	fileList := make([]string, len(fis))
	for i, fi := range fis {
		fileList[i] = fi.Name()
	}
	return fileList, nil
}

func createDirectory(path string) error {
	if !filepath.IsAbs(path) {
		absDir, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("failed to convert input relative path to an absolute path: %v", err)
		}
		path = absDir
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

// SeekAndDestroy walks through a directory and deletes all files that match the pattern
//
// **Parameters:**
//
// path: String representing the path to the directory.
// pattern: String representing the pattern to match.
//
// **Returns:**
//
// error: An error if the files cannot be deleted.
func SeekAndDestroy(path string, pattern string) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		matched, err := filepath.Match(pattern, info.Name())
		if err != nil {
			return err
		}

		if matched {
			if err := os.RemoveAll(path); err != nil {
				return fmt.Errorf("failed to delete file or directory: %v", err)
			}
		}

		return nil
	})
}
