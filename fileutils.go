package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fatih/color"
)

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
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf(color.RedString("failed to read %s: %v", err))
	}

	return strings.Split(string(b), "\n"), nil
}
