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
	"github.com/go-git/go-git/plumbing/transport/ssh"
	cp "github.com/otiai10/copy"
	"github.com/shirou/gopsutil/v3/process"
)

// Signal represents a signal that can be sent to a process.
type Signal int

const (
	// SignalKill is a signal that causes the process to be killed immediately.
	SignalKill Signal = iota
)

// CheckRoot checks if the current process is being run with root permissions.
//
// Returns:
//
// error: An error if the process is not being run as root.
//
// Example:
//
// err := CheckRoot()
//
//	if err != nil {
//	    fmt.Println("The process must be run as root.")
//	}
func CheckRoot() error {
	uid := os.Geteuid()
	if uid != 0 {
		return fmt.Errorf(color.RedString("this script must be run as root - current euid: %v", uid))
	}

	return nil
}

// Cd changes the current working directory to the specified path.
//
// Parameters:
//
// dst: A string specifying the path to the directory to switch to.
//
// Returns:
//
// error: An error if the current directory cannot be changed.
//
// Example:
//
// err := Cd("/path/to/dir")
//
//	if err != nil {
//	    fmt.Println("Failed to change directory.")
//	}
func Cd(dst string) error {
	err := os.Chdir(dst)
	if err != nil {
		fmt.Print(color.RedString("failed to change directory to %s: %v", dst, err))
		return err
	}

	return nil
}

// CmdExists checks if a given command is available in the $PATH.
//
// Parameters:
//
// cmd: A string specifying the name of the command to look for.
//
// Returns:
//
// bool: True if the command exists in the $PATH, otherwise False.
//
// Example:
//
//	if !CmdExists("ls") {
//	    fmt.Println("The 'ls' command is not available.")
//	}
func CmdExists(cmd string) bool {
	if _, err := exec.LookPath(cmd); err != nil {
		return false
	}
	return true
}

// Cp copies a file from the source path to the destination path.
//
// Parameters:
//
// src: A string specifying the path of the file to be copied.
// dst: A string specifying the path to where the file should be copied.
//
// Returns:
//
// error: An error if the file cannot be copied.
//
// Example:
//
// err := Cp("/path/to/src", "/path/to/dst")
//
//	if err != nil {
//	    fmt.Println("Failed to copy the file.")
//	}
func Cp(src string, dst string) error {
	if err := cp.Copy(src, dst); err != nil {
		fmt.Print(color.RedString("failed to copy %s to %s: %v", src, dst, err))
		return err
	}

	return nil
}

// EnvVarSet checks if a given environment variable is set.
//
// Parameters:
//
// key: A string specifying the name of the environment variable to check.
//
// Returns:
//
// error: An error if the environment variable is not set.
//
// Example:
//
// err := EnvVarSet("HOME")
//
//	if err != nil {
//	    fmt.Println("The HOME environment variable is not set.")
//	}
func EnvVarSet(key string) error {
	_, ok := os.LookupEnv(key)
	if !ok {
		err := errors.New(key + " not set")
		return err
	}

	return nil
}

// GetHomeDir returns the current user's home directory.
//
// Returns:
//
// string: The home directory of the current user.
// error: An error if any issue occurs while trying to get the home directory.
//
// Example:
//
// homeDir, err := GetHomeDir()
//
//	if err != nil {
//	  log.Fatalf("failed to get user's home directory: %v", err)
//	}
//
// fmt.Println("Home Directory:", homeDir)
func GetHomeDir() (string, error) {
	out, err := os.UserHomeDir()

	if err != nil {
		return "", fmt.Errorf(color.RedString("failed to get user's home directory: %v", err))
	}

	return out, nil
}

// Gwd returns the current working directory (cwd). If it fails to get the cwd, it prints the error and returns an empty string.
//
// Returns:
//
// string: The current working directory or an empty string if an error occurs.
//
// Example:
//
// cwd := Gwd()
//
//	if cwd == "" {
//	  log.Fatalf("failed to get cwd")
//	}
//
// fmt.Println("Current Working Directory:", cwd)
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
//
// Parameters:
//
// years: An integer representing the number of years to add.
// months: An integer representing the number of months to add.
// days: An integer representing the number of days to add.
//
// Returns:
//
// time.Time: The future date and time calculated from the current time.
//
// Example:
//
// futureTime := GetFutureTime(1, 2, 3)
// fmt.Println("Future date and time:", futureTime)
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

// GetSSHPubKey retrieves the public SSH key for the given key name, decrypting the associated private key if a password
// is provided. It returns a pointer to the public key object, or an error if one occurs.
//
// Parameters:
//
// keyName: A string representing the name of the key to retrieve.
// password: A string representing the password used to decrypt the private key.
//
// Returns:
//
// *ssh.PublicKeys: A pointer to a PublicKeys object representing the retrieved public key.
// error: An error if one occurs during key retrieval or decryption.
//
// Example:
//
// keyName := "id_rsa"
// password := "mypassword"
// publicKey, err := GetSSHPubKey(keyName, password)
//
//	if err != nil {
//	  log.Fatalf("failed to get SSH public key: %v", err)
//	}
//
// log.Printf("Retrieved public key: %v", publicKey)
func GetSSHPubKey(keyName string, password string) (*ssh.PublicKeys, error) {
	var publicKey *ssh.PublicKeys

	sshPath := filepath.Join(os.Getenv("HOME"), ".ssh", keyName)
	publicKey, err := ssh.NewPublicKeysFromFile("git", sshPath, password)
	if err != nil {
		return nil,
			fmt.Errorf(color.RedString(
				"failed to retrieve public SSH key %s: %v",
				keyName, err))
	}

	return publicKey, nil
}

// IsDirEmpty checks if an input directory (name) is empty
//
// Parameters:
//
// name: A string representing the path to the directory to check.
//
// Returns:
//
// bool: A boolean indicating whether the directory is empty or not.
// error: An error if there was any problem reading the directory.
//
// Example:
//
// isEmpty, err := IsDirEmpty("/path/to/directory")
//
//	if err != nil {
//	  log.Fatalf("Error checking directory: %v", err)
//	}
//
// fmt.Println("Is directory empty:", isEmpty)
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
// fmt.Printf("failed to kill process: %v", err)
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
//
// Parameters:
//
// cmd: A string representing the command to run.
// args: A variadic parameter representing any command line arguments to the command.
//
// Returns:
//
// string: The output from the command.
// error: An error if there was any problem running the command.
//
// Example:
//
// output, err := RunCommand("ls", "-l")
//
//	if err != nil {
//	  log.Fatalf("Error running command: %v", err)
//	}
//
// fmt.Println("Command output:", output)
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
//
// Parameters:
//
// to: A time.Duration representing the number of seconds to allow the command to run before timing out.
// command: A string representing the command to run.
// args: A variadic parameter representing any command line arguments to the command.
//
// Returns:
//
// string: The output from the command.
// error: An error if there was any problem running the command.
//
// Example:
//
// output, err := RunCommandWithTimeout(time.Second*5, "sleep", "10")
//
//	if err != nil {
//	  log.Fatalf("Error running command: %v", err)
//	}
//
// fmt.Println("Command output:", output)
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
//
// Parameters:
//
// path: A string representing the path to remove.
//
// Returns:
//
// error: An error if there was any problem removing the path.
//
// Example:
//
// err := RmRf("/path/to/remove")
//
//	if err != nil {
//	  log.Fatalf("Error removing path: %v", err)
//	}
//
// fmt.Println("Path successfully removed!")
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

// ExpandHomeDir expands the tilde (~) in a given path to the current user's home directory.
//
// Parameters:
//
// path: A string representing the path to be expanded.
//
// Returns:
//
// string: The expanded path.
//
// Example:
//
// path := "~/Documents/project"
// expandedPath := ExpandHomeDir(path)
//
// fmt.Println("Expanded Path:", expandedPath)
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
