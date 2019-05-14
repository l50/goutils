package utils

import (
	"os"
)

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
