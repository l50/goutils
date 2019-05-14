package utils

import (
	"os"
)

func CreateEmptyFile(name string) bool {
	file, err := os.Create(name)
	if err != nil {
		return false
	}
	file.Close()
	return true
}

func FileExists(fileLoc string) bool {
	if _, err := os.Stat(fileLoc); !os.IsNotExist(err) {
		return true
	} else {
		return false
	}
}

func GetHomeDir() (string, error) {
	out, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}
	return out, nil
}
