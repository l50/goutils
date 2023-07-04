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

### CreatePackageDocs

```go
CreatePackageDocs(afero.Fs, Repo, string) error
```

CreatePackageDocs generates documentation for all Go packages in the current
directory and its subdirectories. It traverses the file tree using a provided
afero.Fs and Repo to create a new README.md file in each directory containing
a Go package. It uses a specified template file for generating the README files.

It will ignore any files or directories listed in the .docgenignore file
found at the root of the repository. The .docgenignore file should contain
a list of files and directories to ignore, with each entry on a new line.

**Parameters:**

fs: An afero.Fs instance for mocking the filesystem for testing.
repo: A Repo instance representing the GitHub repository
containing the Go packages.

templatePath:  A string representing the path to the template file to be
used for generating README files.

**Returns:**

error: An error, if it encounters an issue while walking the file tree,
reading a directory, parsing Go files, or generating README.md files.

---

### FixCodeBlocks

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
