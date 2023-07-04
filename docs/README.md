# goutils/v2/docs

The `docs` package is a collection of utility functions
designed to simplify common docs tasks.

Table of contents:

- [Functions](#functions)
- [Installation](#installation)
- [Usage](#usage)
- [Tests](#tests)
- [Contributing](#contributing)
- [License](#license)

---

## Functions

## CreatePackageDocs(afero.Fs, Repo, string, ...string)

```go
CreatePackageDocs(afero.Fs, Repo, string, ...string) error
```

CreatePackageDocs generates package documentation for a Go project using
a specified template file. It first checks if the template file exists in
the filesystem denoted by a provided afero.Fs instance. If it exists, the
function walks the project directory, excluding any specified packages,
and applies the template to each non-excluded package to generate its
documentation.

**Parameters:**

fs: An afero.Fs instance representing the filesystem.

repo: A Repo instance containing the Go project's repository details.

templatePath: A string representing the path to the template file to be
used for generating the package documentation.

excludedPackages: Zero or more strings representing the names of packages
to be excluded from documentation generation.

**Returns:**

error: An error, if it encounters an issue while checking if the template
file exists, walking the project directory, or generating the package
documentation.

---

## FixCodeBlocks(fileutils.RealFile, string)

```go
FixCodeBlocks(fileutils.RealFile, string) error
```

FixCodeBlocks processes a provided file to ensure that all code
blocks within comments are surrounded by markdown fenced code block
delimiters with the specified language.

**Parameters:**
file: An object satisfying the File interface, which is to be processed.
language: A string representing the language for the code blocks.

**Returns:**
error: An error if there's an issue reading or writing the file.

---

## Installation

To use the goutils/v2/docs package, you first need to install it.
Follow the steps below to install via go get.

```bash
go get github.com/goutils/v2/l50/docs
```

---

## Usage

After installation, you can import the package in your Go project
using the following import statement:

```go
import "github.com/goutils/v2/l50/docs"
```

---

## Tests

To ensure the package is working correctly, run the following
command to execute the tests for `goutils/v2/docs`:

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
