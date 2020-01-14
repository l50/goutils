package utils

import (
	"io"
	"os"
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

// GetHomeDir returns the path to current user's home directory
func GetHomeDir() (string, error) {
	out, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}
	return out, nil
}

// IsDirEmpty checks if an input directory (name) is empty
// Resource: https://socketloop.com/tutorials/golang-determine-if-directory-is-empty-with-os-file-readdir-function
func IsDirEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// read in ONLY one file
	_, err = f.Readdir(1)

	// and if the file is EOF... well, the dir is empty.
	if err == io.EOF {
		return true, nil
	}
	return false, err
}
