# goutils/dev

The `dev` package is a part of `goutils` library. It provides utility 
functions for development-oriented operations in Go.

---

## Functions

### GHRelease

```go
func GHRelease(newVer string) error
```

Create a new release on GitHub with the given new version.

### GoReleaser

```go
func GoReleaser() error
```

Run the Goreleaser tool to generate all the supported binaries specified in a .goreleaser configuration file.

### InstallVSCodeModules

```go
func InstallVSCodeModules() error
```

Installs the modules used by the vscode-go extension in Visual Studio Code.

### ModUpdate

```go
func ModUpdate(recursive bool, v bool) error
```

Updates go modules by running `go get -u` or 
`go get -u ./...` if recursive is set to true.

### Tidy

```go
func Tidy() error
```

Run `go mod tidy` to clean up the 
module dependencies.

### UpdateMageDeps

```go
func UpdateMageDeps(magedir string) error
```

Update the dependencies in a specified Magefile directory.

### InstallGoDeps

```go
func InstallGoDeps(deps []string) error
```

Install the specified Go dependencies by 
running `go install` for each dependency.

### FindExportedFunctionsInPackage

```go
func FindExportedFunctionsInPackage(pkgPath string) ([]FuncInfo, error)
```

Find all exported functions in a given Go 
package by recursively parsing all non-test
Go files in the package directory.

### FindExportedFuncsWithoutTests

```go
func FindExportedFuncsWithoutTests(pkgPath string) ([]string, error)
```

Finds all exported functions in a given package path that 
do not have corresponding tests.

---

## Installation

To use the `goutils/dev` package, you need to install it via `go get`:

```bash
go get github.com/l50/goutils/dev
```

---

## Usage

After installation, you can import it in your project:

```go
import "github.com/l50/goutils/dev"
```

---

## Tests

To run the tests for the `goutils/dev` package, navigate to
your `$GOPATH/src/github.com/l50/goutils/dev` directory
and run go test:

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
the [LICENSE](../../LICENSE) file for details.
