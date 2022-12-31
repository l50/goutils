package utils

import (
	"bufio"
	"fmt"
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
	defer f.Close()
	if _, err := f.WriteString(text + "\n"); err != nil {
		return err
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
func CreateFile(fileContents []byte, filePath string) error {
	err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot create dir portion"+
			"of filepath %s: %v", filePath, err)
	}

	err = os.WriteFile(filePath, fileContents, os.ModePerm)
	if err != nil {
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

// DeleteFile deletes the input file
func DeleteFile(file string) error {
	err := os.Remove(file)
	if err != nil {
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

// FileToSlice reads an input file into a slice
// and returns it.
func FileToSlice(fileName string) ([]string, error) {
	b, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %v", fileName, err)
	}

	return strings.Split(string(b), "\n"), nil
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
	defer f.Close()

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

	return false, err
}

// RmRf removes an input path and everything in it.
// If the input path doesn't exist, an error is returned.
func RmRf(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("failed to run RmRf on %s: %v", path, err)
		}
		return fmt.Errorf("failed to os.Stat on %s: %v", path, err)
	}

	return nil
}
