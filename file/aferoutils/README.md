# goutils/v2/aferoutils

The `aferoutils` package is a collection of utility functions
designed to simplify common aferoutils tasks.

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

### Tree(afero.Fs, string, io.Writer)

```go
Tree(afero.Fs, string, io.Writer) error
```

Tree displays the directory tree structure starting from the
specified directory path in a format similar to the `tree` command.

**Parameters:**

fs: The afero.Fs representing the file system to use.
dirPath: The path of the directory to display the tree structure for.
prefix: The prefix string to use for each line of the tree structure.
indent: The indent string to use for each level of the tree structure.
out: The io.Writer to write the tree structure output to.

**Returns:**

error: An error if any issue occurs while trying to display the tree structure.

---

## Installation

To use the goutils/v2/aferoutils package, you first need to install it.
Follow the steps below to install via go get.

```bash
go get github.com/l50/goutils/v2/aferoutils
```

---

## Usage

After installation, you can import the package in your Go project
using the following import statement:

```go
import "github.com/l50/goutils/v2/aferoutils"
```

---

## Tests

To ensure the package is working correctly, run the following
command to execute the tests for `goutils/v2/aferoutils`:

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
