# goutils/v2/macos

The `macos` package is a collection of utility functions
designed to simplify common macos tasks.

Table of contents:

- [Functions](#functions)
- [Installation](#installation)
- [Usage](#usage)
- [Tests](#tests)
- [Contributing](#contributing)
- [License](#license)

---

## Functions

### InstallBrewDeps([]string)

```go
InstallBrewDeps([]string) error
```

InstallBrewDeps executes brew install for the input packages.
If any installation fails, it returns an error.

**Parameters:**

brewPackages: Slice of strings representing the packages to install.

**Returns:**

error: An error if any package fails to install.

---

### InstallBrewTFDeps()

```go
InstallBrewTFDeps() error
```

InstallBrewTFDeps installs dependencies for terraform projects
using homebrew. The dependencies include several shell and
terraform tools. If any installation fails, it returns an error.

**Returns:**

error: An error if any package fails to install.

---

## Installation

To use the goutils/v2/macos package, you first need to install it.
Follow the steps below to install via go get.

```bash
go get github.com/l50/goutils/v2/macos
```

---

## Usage

After installation, you can import the package in your Go project
using the following import statement:

```go
import "github.com/l50/goutils/v2/macos"
```

---

## Tests

To ensure the package is working correctly, run the following
command to execute the tests for `goutils/v2/macos`:

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
