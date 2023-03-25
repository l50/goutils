package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/fatih/color"
	cp "github.com/otiai10/copy"
	"github.com/shirou/gopsutil/v3/process"
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

// CmdExists checks $PATH for
// for the input cmd.
// It returns true if the command is found,
// otherwise it returns false.
func CmdExists(cmd string) bool {
	if _, err := exec.LookPath(cmd); err != nil {
		return false
	}
	return true
}

// Cp is used to copy a file from `src` to `dst`.
func Cp(src string, dst string) error {
	if err := cp.Copy(src, dst); err != nil {
		fmt.Print(color.RedString("failed to copy %s to %s: %v", src, dst, err))
		return err
	}

	return nil
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

// GetOSAndArch detects the current system's OS and architecture, and returns them as strings.
// The function returns an error if the OS or architecture is not supported.
//
// Example usage:
//
//	osName, archName, err := GetOSAndArch()
//	if err != nil {
//		fmt.Printf("Error detecting OS and architecture: %v\n", err)
//	} else {
//		fmt.Printf("Detected OS: %s, architecture: %s\n", osName, archName)
//	}
//
// Returns:
//
//	string: The detected operating system name (i.e., "linux", "darwin", or "windows").
//	string: The detected architecture name (i.e., "amd64", "arm64", or "armv").
//	error: An error if the OS or architecture is not supported or cannot be detected.
func GetOSAndArch() (string, string, error) {
	osName := strings.ToLower(runtime.GOOS)
	archName := runtime.GOARCH

	switch archName {
	case "x86_64", "amd64":
		archName = "amd64"
	case "aarch64", "arm64":
		archName = "arm64"
	case "armv6", "armv7", "armv7l", "armv8":
		archName = "armv"
	default:
		return "", "", fmt.Errorf("unsupported architecture: %s", archName)
	}

	return osName, archName, nil
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
func RunCommandWithTimeout(to time.Duration, command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	// Start the cmd.
	if err := cmd.Start(); err != nil {
		return "", err
	}

	// Used to avoid race conditions on timedOut
	var mu sync.Mutex

	mu.Lock()
	timedOut := false
	mu.Unlock()

	// Create channel to grab any errors from the anonymous function below.
	errCh := make(chan error)

	// Create timer that triggers killing the created
	// cmd process once the input duration (to) has been met.
	timeout := time.AfterFunc(to, func() {
		mu.Lock()
		timedOut = true
		mu.Unlock()
		paramStr := strings.Join(args, " ")
		// Once the timer has finished, get the cmd PID(s):
		pids, _ := process.Pids()
		for _, pid := range pids {
			proc, _ := process.NewProcess(pid)
			cmd, _ := proc.Cmdline()
			if cmd == command+" "+paramStr {
				if err := syscall.Kill(int(pid), syscall.SIGKILL); err != nil {
					errCh <- err
				}
			}
		}
	})

	// Save output of the run cmd.
	var out bytes.Buffer
	if _, err := io.Copy(&out, stdout); err != nil {
		return "", err
	}
	stdout.Close()

	// Wait for cmd to exit and return any error that occurred.
	err = cmd.Wait()

	// Stop the timer created earlier.
	timeout.Stop()

	// Check if an error was sent through the channel
	select {
	case err := <-errCh:
		return "", err
	default:
	}

	// Remove newlines from output captured as a string.
	output := strings.TrimSpace(out.String())

	// If an expected timeout occurs, don't return an error.
	mu.Lock()
	if err != nil && !timedOut {
		mu.Unlock()
		return "", err
	}
	mu.Unlock()

	return output, nil
}
