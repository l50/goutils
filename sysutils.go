package utils

import (
	"os/exec"
)

// RunCommand runs a specified system command
func RunCommand(cmd string, args ...string) (string, error) {

	out, err := exec.Command(cmd, args...).CombinedOutput()

	if err != nil {
		return "", err
	}

	return string(out), nil
}
