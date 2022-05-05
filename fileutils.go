package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	if _, err := os.Stat(fileLoc); !os.IsNotExist(err) {
		return true
	}

	return false
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
