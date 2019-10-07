package utils

import (
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
