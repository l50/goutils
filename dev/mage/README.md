# goutils/v2/mageutils

The `mageutils` package is a collection of utility functions
designed to simplify common mageutils tasks.

Table of contents:

- [Functions](#functions)
- [Installation](#installation)
- [Usage](#usage)
- [Tests](#tests)
- [Contributing](#contributing)
- [License](#license)

---

## Functions

### Compile(string, string, string)

```go
Compile(string, string, string) error
```

Compile builds a Go application for a specified operating system and
architecture. It sets the appropriate environment variables and runs `go
build`. The compiled application is placed in the specified build path.

**Parameters:**

buildPath: The output directory for the compiled application.
goOS: The target operating system (e.g., "linux", "darwin", "windows").
goArch: The target architecture (e.g., "amd64", "arm64").

**Returns:**

error: An error if the compilation process encounters one.

---

### FindExportedFuncsWithoutTests(string)

```go
FindExportedFuncsWithoutTests(string) []string, error
```

FindExportedFuncsWithoutTests discovers all exported functions in a given
package path that lack corresponding tests.

**Parameters:**

pkgPath: A string defining the package path to search.

**Returns:**

[]string: A slice of strings containing the names of exported functions that
lack corresponding tests.

error: An error if there was a problem parsing the package or finding the tests.

---

### FindExportedFunctionsInPackage(string)

```go
FindExportedFunctionsInPackage(string) []FuncInfo, error
```

FindExportedFunctionsInPackage finds all exported functions in a given Go
package by parsing all non-test Go files in the package directory. It returns
a slice of FuncInfo structs. Each contains the file path and the name of an
exported function. If no exported functions are found in the package, an
error is returned.

**Parameters:**

pkgPath: A string representing the path to the directory containing the package
to search for exported functions.

**Returns:**

[]FuncInfo: A slice of FuncInfo structs, each containing the file path and the
name of an exported function found in the package.
error: An error if no exported functions are found.

---

### GHRelease(string)

```go
GHRelease(string) error
```

GHRelease creates a new release on GitHub using the given new version.
It requires the gh CLI tool to be available on the PATH.

**Parameters:**

newVer: A string specifying the new version, e.g., "v1.0.1"

**Returns:**

error: An error if the GHRelease function is not successful.

---

### GoReleaser()

```go
GoReleaser() error
```

GoReleaser runs the Goreleaser tool to generate all the supported binaries
specified in a .goreleaser configuration file.

**Returns:**

error: An error if the Goreleaser function is not successful.

---

### InstallGoDeps([]string)

```go
InstallGoDeps([]string) error
```

InstallGoDeps installs the specified Go dependencies by executing 'go install'
for each dependency.

**Parameters:**

deps: A slice of strings defining the dependencies to install.

**Returns:**

error: An error if the InstallGoDeps function didn't run successfully.

---

### InstallVSCodeModules()

```go
InstallVSCodeModules() error
```

InstallVSCodeModules installs the modules used by the vscode-go extension in
Visual Studio Code.

**Returns:**

error: An error if the InstallVSCodeModules function is not successful.

---

### ModUpdate(bool, bool)

```go
ModUpdate(bool, bool) error
```

ModUpdate updates go modules by running 'go get -u' or 'go get -u ./...' if
recursive is set to true. The function will run in verbose mode if 'v' is
set to true.

**Parameters:**

recursive: A boolean specifying whether to run the update recursively.
v: A boolean specifying whether to run the update in verbose mode.

**Returns:**

error: An error if the ModUpdate function is not successful.

---

### Tidy()

```go
Tidy() error
```

Tidy executes 'go mod tidy' to clear the module dependencies.

**Returns:**

error: An error if the Tidy function didn't run successfully.

---

### UpdateMageDeps(string)

```go
UpdateMageDeps(string) error
```

UpdateMageDeps modifies the dependencies in a given Magefile directory.
If no directory is provided, it falls back to the 'magefiles' directory.

**Parameters:**

magedir: A string defining the path to the magefiles directory.

**Returns:**

error: An error if the UpdateMageDeps function didn't run successfully.

---

## Installation

To use the goutils/v2/mageutils package, you first need to install it.
Follow the steps below to install via go get.

```bash
go get github.com/l50/goutils/v2/mageutils
```

---

## Usage

After installation, you can import the package in your Go project
using the following import statement:

```go
import "github.com/l50/goutils/v2/mageutils"
```

---

## Tests

To ensure the package is working correctly, run the following
command to execute the tests for `goutils/v2/mageutils`:

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
