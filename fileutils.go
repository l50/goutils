package utils

import (
	"os"
)

func GetHomeDir() (string, error) {
	out, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}
	return out, nil
}
