package utils

import (
	"os/exec"
)

// RunCommand runs a specified system command
func RunCommand(cmd string, args ...string) (output string, outputErr string) {

	out, err := exec.Command(cmd, args...).CombinedOutput()

	if err != nil {
		outputErr = err.Error()
	}

	return string(out), outputErr
}
