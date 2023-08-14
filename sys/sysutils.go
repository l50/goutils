package sys

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	fileutils "github.com/l50/goutils/v2/file/fileutils"
	cp "github.com/otiai10/copy"
)

// Signal represents a signal that can be sent to a process.
//
// **Attributes:**
//
// SignalKill: A signal that causes the process to be killed immediately.
type Signal int

// Signal constants
const (
	// SignalKill represents a signal that kills a process immediately
	SignalKill Signal = iota
)

// CheckRoot checks if the current process is being run with root permissions.
//
// **Returns:**
//
// error: An error if the process is not being run as root.
func CheckRoot() error {
	uid := os.Geteuid()
	if uid != 0 {
		return fmt.Errorf("this script must be run as root - current euid: %v", uid)
	}

	return nil
}

// Cd changes the current working directory to the specified path.
//
// **Parameters:**
//
// dst: A string specifying the path to the directory to switch to.
//
// **Returns:**
//
// error: An error if the current directory cannot be changed.
func Cd(path string) error {
	if err := os.Chdir(path); err != nil {
		return fmt.Errorf("failed to change directory to %s: %v", path, err)
	}

	return nil
}

// CmdExists checks if a given command is available in the $PATH.
//
// **Parameters:**
//
// cmd: A string specifying the name of the command to look for.
//
// **Returns:**
//
// bool: True if the command exists in the $PATH, otherwise False.
func CmdExists(cmd string) bool {
	if _, err := exec.LookPath(cmd); err != nil {
		return false
	}
	return true
}

// Cp copies a file from the source path to the destination path.
//
// **Parameters:**
//
// src: A string specifying the path of the file to be copied.
// dst: A string specifying the path to where the file should be copied.
//
// **Returns:**
//
// error: An error if the file cannot be copied.
func Cp(src string, dst string) error {
	if err := cp.Copy(src, dst); err != nil {
		return fmt.Errorf("failed to copy %s to %s: %v", src, dst, err)
	}

	return nil
}

// EnvVarSet checks whether a given environment variable is set.
//
// **Parameters:**
//
// key: String specifying the name of the environment variable.
//
// **Returns:**
//
// error: Error if the environment variable is not set.
func EnvVarSet(key string) error {
	_, ok := os.LookupEnv(key)
	if !ok {
		return errors.New(key + " not set")
	}

	return nil
}

// ExpandHomeDir expands the tilde (~) in a path to the home
// directory of the current user.
//
// **Parameters:**
//
// path: String representing the path to be expanded.
//
// **Returns:**
//
// string: The expanded path.
func ExpandHomeDir(path string) string {
	if len(path) == 0 || path[0] != '~' {
		return path
	}

	homeDir, err := GetHomeDir()
	if err != nil {
		return path
	}

	if len(path) == 1 || path[1] == '/' {
		return filepath.Join(homeDir, path[1:])
	}

	return filepath.Join(homeDir, path[1:])
}

// GetHomeDir fetches the home directory of the current user.
//
// **Returns:**
//
// string: The home directory of the current user.
// error: Error if there is an issue fetching the home directory.
func GetHomeDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		if _, ok := err.(user.UnknownUserError); ok {
			return "", fmt.Errorf("failed to get user's home directory: user unknown")
		}
		return "", fmt.Errorf("failed to get user's home directory: %v", err)
	}
	return homeDir, nil
}

// GetSSHPubKey retrieves the public SSH key for the given key name,
// decrypting the associated private key if a password is provided.
//
// **Parameters:**
//
// keyName: String representing the name of the key to retrieve.
// password: String for the password used to decrypt the private key.
//
// **Returns:**
//
// *ssh.PublicKeys: Pointer to the PublicKeys object for the retrieved key.
// error: Error if one occurs during key retrieval or decryption.
func GetSSHPubKey(keyName string, password string) (*ssh.PublicKeys, error) {
	var publicKey *ssh.PublicKeys

	sshPath := filepath.Join(os.Getenv("HOME"), ".ssh", keyName)
	publicKey, err := ssh.NewPublicKeysFromFile("git", sshPath, password)
	if err != nil {
		return nil,
			fmt.Errorf("failed to retrieve public SSH key %s: %v",
				keyName, err)
	}

	return publicKey, nil
}

// Gwd gets the current working directory (cwd). In case of failure, it logs
// the error and returns an empty string.
//
// **Returns:**
//
// string: The current working directory or an empty string if an error occurs.
func Gwd() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("failed to get cwd: %v", err)
		return ""
	}

	return dir
}

// GetFutureTime calculates the date and time after the input years, months, and
// days from the current time.
//
// **Parameters:**
//
// years: The number of years to add.
// months: The number of months to add.
// days: The number of days to add.
//
// **Returns:**
//
// time.Time: The future date and time calculated from the current time.
func GetFutureTime(years int, months int, days int) time.Time {
	t := time.Now()
	exp := t.AddDate(years, months, days)
	return exp
}

// RuntimeInfoProvider is an interface for providing information about the
// current runtime environment.
type RuntimeInfoProvider interface {
	GetOS() string
	GetArch() string
}

// DefaultRuntimeInfoProvider is the default implementation of the
// RuntimeInfoProvider interface.
type DefaultRuntimeInfoProvider struct{}

// GetOS returns the current operating system.
//
// **Returns:**
//
// string: The current operating system.
func (p *DefaultRuntimeInfoProvider) GetOS() string {
	return strings.ToLower(runtime.GOOS)
}

// GetArch returns the current architecture.
//
// **Returns:**
//
// string: The current architecture.
func (p *DefaultRuntimeInfoProvider) GetArch() string {
	return runtime.GOARCH
}

// GetOSAndArch identifies the current system's OS and architecture, and returns
// them as strings. The function returns an error if the OS or architecture is
// not supported.
//
// **Returns:**
//
// string: Detected operating system name (e.g., "linux", "darwin", "windows").
// string: Detected architecture name (e.g., "amd64", "arm64", "armv").
// error: An error if the OS or architecture is not supported or cannot be detected.
func GetOSAndArch(provider RuntimeInfoProvider) (string, string, error) {
	osName := provider.GetOS()
	archName := provider.GetArch()

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

// IsDirEmpty checks whether the input directory (name) is empty.
//
// **Parameters:**
//
// name: The path to the directory to check.
//
// **Returns:**
//
// bool: A flag indicating whether the directory is empty.
// error: An error if there's a problem reading the directory.
func IsDirEmpty(name string) (bool, error) {
	info, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false, fmt.Errorf("%s does not exist", name)
	}
	if err != nil {
		return false, fmt.Errorf("failed to determine if %s is empty: %v", name, err)
	}

	if !info.IsDir() {
		return false, fmt.Errorf("%s is not a directory", name)
	}

	entries, err := os.ReadDir(name)
	if err != nil {
		return false, fmt.Errorf("failed to read directory entries for %s: %v", name, err)
	}

	return len(entries) == 0, nil
}

// KillProcess sends a signal to the process with the specified PID. On Windows,
// it uses the taskkill command to terminate the process. On Unix-like systems,
// it sends the specified signal to the process using the syscall.Kill function.
//
// Note that SignalKill may not work on all platforms. For more information,
// see the documentation for the syscall package.
//
// **Parameters:**
//
// pid: The process ID to kill.
// signal: The signal to send to the process. Currently, only SignalKill is
// supported, which terminates the process.
//
// **Returns:**
//
// error: An error if the process couldn't be killed.
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

// RunCommand executes a specified system command.
//
// **Parameters:**
//
// cmd: A string representing the command to run.
// args: A variadic parameter representing any command line arguments to the command.
//
// **Returns:**
//
// string: The output from the command.
// error: An error if there was any problem running the command.
func RunCommand(cmd string, args ...string) (string, error) {
	execCmd := exec.Command(cmd, args...)
	execCmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true} // create new process group

	var stdoutBuf, stderrBuf bytes.Buffer
	multiStdout := io.MultiWriter(os.Stdout, &stdoutBuf) // write to both os.Stdout and stdoutBuf
	multiStderr := io.MultiWriter(os.Stderr, &stderrBuf) // write to both os.Stderr and stderrBuf

	// Attach to standard output and standard error
	execCmd.Stdout = multiStdout
	execCmd.Stderr = multiStderr

	if err := execCmd.Run(); err != nil {
		return "", fmt.Errorf("failed to run %s with args %v: stdout: %s, stderr: %s, err: %v",
			cmd, args, stdoutBuf.String(), stderrBuf.String(), err)
	}

	return stdoutBuf.String(), nil
}

// RunCommandWithTimeout executes a command for a specified number of
// seconds before timing out. The command will be run in its own
// process group to allow for killing child processes if necessary.
//
// **Parameters:**
//
// to: An int representing the number of seconds to allow
// the command to run before timing out.
// command: A string representing the command to run.
// args: A variadic parameter representing any command line arguments to the command.
//
// **Returns:**
//
// string: The output from the command if it completes successfully
// before the timeout. If the command does not complete before the
// timeout or an error occurs, an empty string is returned.
// error: An error if there was any problem running the command or if the
// command does not complete before the timeout.
func RunCommandWithTimeout(to int, cmd string, args ...string) (string, error) {
	timeout := time.Duration(to) * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	done := make(chan string, 1)
	errCh := make(chan error, 1)

	go func() {
		output, err := RunCommand(cmd, args...)
		if err != nil {
			errCh <- err
		} else {
			done <- output
		}
	}()

	select {
	case <-ctx.Done():
		// If the context is done, check the reason
		if ctx.Err() == context.DeadlineExceeded {
			// The command timed out, now force kill the process group
			return "", fmt.Errorf("command timed out")
		}
		// if ctx.Err() is not DeadlineExceeded, that means the context was cancelled
		// for some other reason, which should not happen under normal circumstances.
		return "", fmt.Errorf("unexpected context cancellation")
	case output := <-done:
		// The command completed before the timeout
		return output, nil
	case err := <-errCh:
		// There was an error running the command
		return "", err
	}
}

// RmRf deletes an input path and everything in it.
// If the input path doesn't exist, an error is returned.
//
// **Parameters:**
//
// path: A string representing the path to remove.
//
// **Returns:**
//
// error: An error if there was any problem removing the path.
func RmRf(file fileutils.File) error {
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to os.Stat: %v", err)
	}

	if info.IsDir() {
		if err := file.RemoveAll(); err != nil {
			return fmt.Errorf("failed to run RemoveAll: %v", err)
		}
	} else {
		if err := file.Remove(); err != nil {
			return fmt.Errorf("failed to run Remove: %v", err)
		}
	}

	return nil
}
