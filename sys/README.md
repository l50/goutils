# goutils/v2/sys

The `sys` package is a collection of utility functions
designed to simplify common sys tasks.

---

## Table of contents

- [Functions](#functions)
- [Installation](#installation)
- [Usage](#usage)
- [Tests](#tests)
- [Contributing](#contributing)
- [License](#license)

---

## Functions

### Cd(string)

```go
Cd(string) error
```

Cd changes the current working directory to the specified path.

**Parameters:**

dst: A string specifying the path to the directory to switch to.

**Returns:**

error: An error if the current directory cannot be changed.

---

### CheckRoot()

```go
CheckRoot() error
```

CheckRoot checks if the current process is being run with root permissions.

**Returns:**

error: An error if the process is not being run as root.

---

### CmdExists(string)

```go
CmdExists(string) bool
```

CmdExists checks if a given command is available in the $PATH.

**Parameters:**

cmd: A string specifying the name of the command to look for.

**Returns:**

bool: True if the command exists in the $PATH, otherwise False.

---

### Cp(string, string)

```go
Cp(string, string) error
```

Cp copies a file from the source path to the destination path.

**Parameters:**

src: A string specifying the path of the file to be copied.
dst: A string specifying the path to where the file should be copied.

**Returns:**

error: An error if the file cannot be copied.

---

### DefaultRuntimeInfoProvider.GetArch()

```go
GetArch() string
```

GetArch returns the current architecture.

**Returns:**

string: The current architecture.

---

### DefaultRuntimeInfoProvider.GetOS()

```go
GetOS() string
```

GetOS returns the current operating system.

**Returns:**

string: The current operating system.

---

### EnvVarSet(string)

```go
EnvVarSet(string) error
```

EnvVarSet checks whether a given environment variable is set.

**Parameters:**

key: String specifying the name of the environment variable.

**Returns:**

error: Error if the environment variable is not set.

---

### ExpandHomeDir(string)

```go
ExpandHomeDir(string) string
```

ExpandHomeDir expands the tilde (~) in a path to the home
directory of the current user.

**Parameters:**

path: String representing the path to be expanded.

**Returns:**

string: The expanded path.

---

### GetFutureTime(int, int, int)

```go
GetFutureTime(int, int, int) time.Time
```

GetFutureTime calculates the date and time after the input years, months, and
days from the current time.

**Parameters:**

years: The number of years to add.
months: The number of months to add.
days: The number of days to add.

**Returns:**

time.Time: The future date and time calculated from the current time.

---

### GetHomeDir()

```go
GetHomeDir() string, error
```

GetHomeDir fetches the home directory of the current user.

**Returns:**

string: The home directory of the current user.
error: Error if there is an issue fetching the home directory.

---

### GetOSAndArch(RuntimeInfoProvider)

```go
GetOSAndArch(RuntimeInfoProvider) string, string, error
```

GetOSAndArch identifies the current system's OS and architecture, and returns
them as strings. The function returns an error if the OS or architecture is
not supported.

**Returns:**

string: Detected operating system name (e.g., "linux", "darwin", "windows").
string: Detected architecture name (e.g., "amd64", "arm64", "armv").
error: An error if the OS or architecture is not supported or cannot be detected.

---

### GetSSHPubKey(string, string)

```go
GetSSHPubKey(string, string) *ssh.PublicKeys, error
```

GetSSHPubKey retrieves the public SSH key for the given key name,
decrypting the associated private key if a password is provided.

**Parameters:**

keyName: String representing the name of the key to retrieve.
password: String for the password used to decrypt the private key.

**Returns:**

*ssh.PublicKeys: Pointer to the PublicKeys object for the retrieved key.
error: Error if one occurs during key retrieval or decryption.

---

### GetTempPath()

```go
GetTempPath() string
```

GetTempPath determines the path to the temporary directory based on the
operating system. This function is useful for retrieving a standard location
for temporary files and directories.

**Returns:**

string: The path to the temporary directory. On Windows, it returns 'C:\\Temp'.
On Unix/Linux systems, it returns '/tmp'.

---

### Gwd()

```go
Gwd() string
```

Gwd gets the current working directory (cwd). In case of failure, it logs
the error and returns an empty string.

**Returns:**

string: The current working directory or an empty string if an error occurs.

---

### IsDirEmpty(string)

```go
IsDirEmpty(string) bool, error
```

IsDirEmpty checks whether the input directory (name) is empty.

**Parameters:**

name: The path to the directory to check.

**Returns:**

bool: A flag indicating whether the directory is empty.
error: An error if there's a problem reading the directory.

---

### KillProcess(int, Signal)

```go
KillProcess(int, Signal) error
```

KillProcess sends a signal to the process with the specified PID. On Windows,
it uses the taskkill command to terminate the process. On Unix-like systems,
it sends the specified signal to the process using the syscall.Kill function.

Note that SignalKill may not work on all platforms. For more information,
see the documentation for the syscall package.

**Parameters:**

pid: The process ID to kill.
signal: The signal to send to the process. Currently, only SignalKill is
supported, which terminates the process.

**Returns:**

error: An error if the process couldn't be killed.

---

### RmRf(fileutils.File)

```go
RmRf(fileutils.File) error
```

RmRf deletes an input path and everything in it.
If the input path doesn't exist, an error is returned.

**Parameters:**

path: A string representing the path to remove.

**Returns:**

error: An error if there was any problem removing the path.

---

### RunCommand(string, ...string)

```go
RunCommand(string, ...string) string, error
```

RunCommand executes a specified system command.

**Parameters:**

cmd: A string representing the command to run.
args: A variadic parameter representing any command line arguments to the command.

**Returns:**

string: The output from the command.
error: An error if there was any problem running the command.

---

### RunCommandWithTimeout(int, string, ...string)

```go
RunCommandWithTimeout(int, string, ...string) string, error
```

RunCommandWithTimeout executes a command for a specified number of
seconds before timing out. The command will be run in its own
process group to allow for killing child processes if necessary.

**Parameters:**

to: An int representing the number of seconds to allow
the command to run before timing out.
command: A string representing the command to run.
args: A variadic parameter representing any command line arguments to the command.

**Returns:**

string: The output from the command if it completes successfully
before the timeout. If the command does not complete before the
timeout or an error occurs, an empty string is returned.
error: An error if there was any problem running the command or if the
command does not complete before the timeout.

---

## Installation

To use the goutils/v2/sys package, you first need to install it.
Follow the steps below to install via go get.

```bash
go get github.com/l50/goutils/v2/sys
```

---

## Usage

After installation, you can import the package in your Go project
using the following import statement:

```go
import "github.com/l50/goutils/v2/sys"
```

---

## Tests

To ensure the package is working correctly, run the following
command to execute the tests for `goutils/v2/sys`:

```bash
go test -v
```

---

## Contributing

Pull requests are welcome. For major changes,
please open an issue first to discuss what
you would like to change.

---

## License

This project is licensed under the MIT
License - see the [LICENSE](../LICENSE)
file for details.
