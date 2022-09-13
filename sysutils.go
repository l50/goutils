package utils

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/bitfield/script"
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

// Cd is used to change the current working directory
// to the specified destination.
func Cd(dst string) error {
	err := os.Chdir(dst)
	if err != nil {
		fmt.Print(color.RedString("failed to change directory to %s: %v", dst, err))
		return err
	}

	return nil
}

// Cp is used to copy a file from a src to a destination
func Cp(src string, dst string) bool {
	input, err := os.ReadFile(src)
	if err != nil {
		fmt.Print(color.RedString("failed to copy %s to %s: %v", src, dst, err))
		return false
	}

	if err := os.WriteFile(dst, input, 0644); err != nil {
		fmt.Print(color.RedString("failed to copy %s to %s: %v", src, dst, err))
		return false
	}

	return true
}

// EnvVarSet checks if an input environment variable
// is set by checking the input key for
// an associated value.
// If an env var is not set, an error is returned.
func EnvVarSet(key string) error {
	_, ok := os.LookupEnv(key)
	if !ok {
		err := errors.New(key + " not set")
		return err
	}

	return nil
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

// GetFutureTime returns the date and time of the input
// years, months, and days parameters from the current time.
func GetFutureTime(years int, months int, days int) time.Time {
	t := time.Now()
	exp := t.AddDate(years, months, days)
	return exp
}

// IsDirEmpty checks if an input directory (name) is empty
func IsDirEmpty(name string) (bool, error) {
	entries, err := os.ReadDir(name)
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

// RunCommandWithTimeout runs a command for a specified number of
// seconds before timing out and returning the output.
func RunCommandWithTimeout(timeoutS string, command string) (string, error) {
	cmd := fmt.Sprintf("timeout %s %s",
		timeoutS, command)
	p := script.Exec(cmd)
	if p.Error() != nil {
		return "", err
	}

	output, err := p.String()
	if err != nil {
		// Ensure we still provide output if a timeout occurs.
		return output, err
	}

	return output, nil
}
