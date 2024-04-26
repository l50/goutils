# goutils/v2/main

The `main` package is a collection of utility functions
designed to simplify common main tasks.

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

### FixCodeBlocks(string, string)

```go
FixCodeBlocks(string, string) error
```

FixCodeBlocks fixes code blocks for the input filepath
using the input language.

Example usage:

```go
mage fixcodeblocks go docs/docGeneration.go
```

**Parameters:**

filepath: the path to the file or directory to fix

language: the language of the code blocks to fix

**Returns:**

error: an error if one occurred

---

### GeneratePackageDocs()

```go
GeneratePackageDocs() error
```

GeneratePackageDocs creates documentation for the various packages
in the project.

Example usage:

```go
mage generatepackagedocs
```

**Returns:**

error: An error if any issue occurs during documentation generation.

---

### InstallDeps()

```go
InstallDeps() error
```

InstallDeps installs the Go dependencies necessary for developing
on the project.

Example usage:

```go
mage installdeps
```

**Returns:**

error: An error if any issue occurs while trying to
install the dependencies.

---

### RunPreCommit()

```go
RunPreCommit() error
```

RunPreCommit updates, clears, and executes all pre-commit hooks
locally. The function follows a three-step process:

First, it updates the pre-commit hooks.
Next, it clears the pre-commit cache to ensure a clean environment.
Lastly, it executes all pre-commit hooks locally.

Example usage:

```go
mage runprecommit
```

**Returns:**

error: An error if any issue occurs at any of the three stages
of the process.

---

### RunTests()

```go
RunTests() error
```

RunTests executes all unit tests.

Example usage:

```go
mage runtests
```

**Returns:**

error: An error if any issue occurs while running the tests.

---

### UpdateMirror(string)

```go
UpdateMirror(string) error
```

UpdateMirror updates pkg.go.dev with the release associated with the
input tag

Example usage:

```go
mage updatemirror v2.0.1
```

**Parameters:**

tag: the tag to update pkg.go.dev with

**Returns:**

error: An error if any issue occurs while updating pkg.go.dev

---

## Installation

To use the goutils/v2/main package, you first need to install it.
Follow the steps below to install via go get.

```bash
go get github.com/l50/goutils/v2/main
```

---

## Usage

After installation, you can import the package in your Go project
using the following import statement:

```go
import "github.com/l50/goutils/v2/main"
```

---

## Tests

To ensure the package is working correctly, run the following
command to execute the tests for `goutils/v2/main`:

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
