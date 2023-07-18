# goutils/v2/encoderutils

The `encoderutils` package is a collection of utility functions
designed to simplify common encoderutils tasks.

Table of contents:

- [Functions](#functions)
- [Installation](#installation)
- [Usage](#usage)
- [Tests](#tests)
- [Contributing](#contributing)
- [License](#license)

---

## Functions

### Unzip(string)

```go
Unzip(string) error
```

Unzip unzips the specified zip file and extracts the contents of the
zip file to the specified destination.

**Parameters:**

src: A string representing the path to the zip file.

dest: A string representing the path to the destination directory.

**Returns:**

error: An error if any issue occurs while trying to unzip the file.

---

### Zip(string)

```go
Zip(string) error
```

Zip creates a zip file from the specified source directory and saves it to the
specified destination path.

**Parameters:**

srcDir: A string representing the path to the source directory.

destFile: A string representing the path to the destination zip file.

**Returns:**

error: An error if any issue occurs while trying to zip the file.

---

## Installation

To use the goutils/v2/encoderutils package, you first need to install it.
Follow the steps below to install via go get.

```bash
go get github.com/l50/goutils/v2/encoderutils
```

---

## Usage

After installation, you can import the package in your Go project
using the following import statement:

```go
import "github.com/l50/goutils/v2/encoderutils"
```

---

## Tests

To ensure the package is working correctly, run the following
command to execute the tests for `goutils/v2/encoderutils`:

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
