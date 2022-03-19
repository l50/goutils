package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/fatih/color"
)

// CheckRoot will check to see if the process is being run as root
func CheckRoot() error {
	uid := os.Geteuid()
	if uid != 0 {
		return fmt.Errorf(color.RedString("this script must be run as root - current euid: %v", uid))
	}

	return nil
}

// Cp is used to copy a file from a src to a destination
func Cp(src string, dst string) bool {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		fmt.Print(color.RedString("failed to copy %s to %s: %v", src, dst, err))
		return false
	}

	if err := ioutil.WriteFile(dst, input, 0644); err != nil {
		fmt.Print(color.RedString("failed to copy %s to %s: %v", src, dst, err))
		return false
	}

	return true
}

// GetHomeDir returns the path to current user's home directory
func GetHomeDir() (string, error) {
	out, err := os.UserHomeDir()

	if err != nil {
		return "", fmt.Errorf(color.RedString("failed to get user's home directory: %v", err))
	}

	return out, nil
}

// Gwd will return the current working directory
func Gwd() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Print(color.RedString("failed to get cwd: %v", err))
		return ""
	}

	return dir

}

// IsDirEmpty checks if an input directory (name) is empty
func IsDirEmpty(name string) (bool, error) {
	entries, err := ioutil.ReadDir(name)
	if err != nil {
		return false, fmt.Errorf(color.RedString("failed to determine if %s is empty: %v", name, err))
	}

	return len(entries) == 0, nil

}

// RunCommand runs a specified system command
func RunCommand(cmd string, args ...string) (string, error) {
	out, err := exec.Command(cmd, args...).CombinedOutput()

	if err != nil {
		return "", fmt.Errorf(color.RedString(
			"failed to run %s: args: %s, stdout: %s, err: %v", cmd, args, out, err))
	}

	return string(out), nil

}
