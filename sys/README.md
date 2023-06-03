# goutils/sys

The `sys` package is a part of `goutils` library.

These functions together provide a range of system utilities such as file handling,
directory manipulation, process management, environment variable manipulation, and
system information retrieval. They might be a part of a larger system or application
toolkit.

---

## Functions

### CheckRoot

```go
func CheckRoot() error
```

Checks if the current process is being run with root permissions. If not, it returns
an error.

### Cd

```go
func Cd(path string) error
```

Changes the current working directory to the specified path.

### CmdExists

```go
func CmdExists(cmd string) bool
```

Checks if a given command is available in the `$PATH`.

### Cp

```go
func Cp(src string, dst string) error
```

Copies a file from the source path to the destination path.

### EnvVarSet

```go
func EnvVarSet(key string) error
```

Checks if a given environment variable is set. If not, it returns an error.

### GetHomeDir

```go
func GetHomeDir() (string, error) {
```

Returns the current user's home directory.

### GetPublicSSHKey

```go
func GetSSHPubKey(keyName string, password string) (*ssh.PublicKeys, error) {
```

Retrieves the public SSH key for the given key name, decrypting the associated private
key if a password is provided.

### Gwd

```go
func Gwd() (string, error)
```

Returns the current working directory.

### GetFutureTime

```go
func GetFutureTime(years int, months int, days int) time.Time {
```

Returns the date and time of the input years, months, and days parameters from the
current time.

### GetOSAndArch

```go
func GetOSAndArch() (string, string, error) {
```

Returns the current operating system and architecture.

### IsDirEmpty

```go
func IsDirEmpty(name string) (bool, error) {
```

Checks if the input directory is empty. Returns true if it is, otherwise returns
false.

### KillProcess

```go
func KillProcess(pid int, signal Signal) error {
```

Sends a signal to the process with the specified PID.

### RunCommand

```go
func RunCommand(cmd string, args ...string) (string, error) {
```

Runs the specified command with the specified arguments and returns the output as
a string.

### RunCommandWithTimeout

```go
func RunCommandWithTimeout(to time.Duration, command string, args ...string) (string,
error) {
```

Runs a command for a specified number of seconds before timing out and returning
the output.

### RmRf

```go
func RmRf(path string) error {
```

Deletes the specified path and all of its contents.

### ExpandHomeDir

```go
func ExpandHomeDir(path string) string {
```

Expands the `~` character in the input path to the current user's home directory.

---

## Installation

To use the `goutils/sys` package, you need to install it via `go get`:

```bash
go get github.com/l50/goutils/sys
```

---

## Usage

After installation, you can import it in your project:

```go
import "github.com/l50/goutils/sys"
```

---

## Tests

To run the tests for the `goutils/sys` package, navigate to your `$GOPATH/src/github.com/l50/goutils/sys`
directory and run go test:

```bash
go test -v
```

---

## Contributing

Pull requests are welcome. For major changes, please
open an issue first to discuss what you would like to change.

---

## License

This project is licensed under the MIT License - see
the [LICENSE](../LICENSE) file for details.
