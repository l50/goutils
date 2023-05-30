package sys

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/fatih/color"
	cp "github.com/otiai10/copy"
	"github.com/shirou/gopsutil/v3/process"
)

// Signal represents a signal that can be sent to a process.
type Signal int

const (
	// SignalKill is a signal that causes the process to be killed immediately.
	SignalKill Signal = iota
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

// KillProcess sends a signal to the process with the specified PID.
//
// On Windows, it uses the taskkill command to terminate the process.
// On Unix-like systems, it sends the specified signal to the process using the syscall.Kill function.
//
// Parameters:
// - pid: The process ID of the process to kill.
// - signal: The signal to send to the process. Currently, only SignalKill is supported, which will terminate the process.
//
// Returns:
// - error: An error if the process could not be killed.
//
// Example usage:
//
// err := KillProcess(1234, SignalKill)
// if err != nil {
// fmt.Printf("Failed to kill process: %v", err)
// } else {
// fmt.Println("Process terminated successfully")
// }
//
// Note that SignalKill may not work on all platforms. For more information, see the documentation for the syscall package.
func KillProcess(pid int, signal Signal) error {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprintf("%d", pid))
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to kill process: %v", err)
		}
		return nil
	}

	var sig os.Signal
	switch signal {
	case SignalKill:
		sig = syscall.SIGKILL
	default:
		return fmt.Errorf("unsupported signal: %v", signal)
	}

	p, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process: %v", err)
	}

	if err := p.Signal(sig); err != nil {
		return fmt.Errorf("failed to send signal to process: %v", err)
	}

	return nil
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
				if err := KillProcess(int(pid), SignalKill); err != nil {
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

// RmRf removes an input path and everything in it.
// If the input path doesn't exist, an error is returned.
func RmRf(path string) error {
	if _, err := os.Stat(path); err == nil {
		if info, err := os.Stat(path); err == nil {
			if info.IsDir() {
				if err := os.RemoveAll(path); err != nil {
					return fmt.Errorf("failed to run RmRf on %s: %v", path, err)
				}
			} else {
				if err := os.Remove(path); err != nil {
					return fmt.Errorf("failed to run RmRf on %s: %v", path, err)
				}
			}
		} else {
			return fmt.Errorf("failed to os.Stat on %s: %v", path, err)
		}
	} else {
		return fmt.Errorf("failed to os.Stat on %s: %v", path, err)
	}

	return nil
}

// ExpandHomeDir expands the tilde character in a path to the user's home directory.
// The function takes a string representing a path and checks if the first character is a tilde (~).
// If it is, the function replaces the tilde with the user's home directory. The path is returned
// unchanged if it does not start with a tilde or if there's an error retrieving the user's home
// directory.
//
// Example usage:
//
//	pathWithTilde := "~/Documents/myfileutils.txt"
//	expandedPath := ExpandHomeDir(pathWithTilde)
//
// Parameters:
//
//	path: The string containing a path that may start with a tilde (~) character.
//
// Returns:
//
//	string: The expanded path with the tilde replaced by the user's home directory, or the
//	        original path if it does not start with a tilde or there's an error retrieving
//	        the user's home directory.
func ExpandHomeDir(path string) string {
	if len(path) == 0 || path[0] != '~' {
		return path
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return path
	}

	if len(path) == 1 || path[1] == '/' {
		return filepath.Join(homeDir, path[1:])
	}

	return filepath.Join(homeDir, path[1:])
}
