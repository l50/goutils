# goutils/v2/ansible

The `ansible` package is a collection of utility functions
designed to simplify common ansible tasks.

Table of contents:

- [Functions](#functions)
- [Installation](#installation)
- [Usage](#usage)
- [Tests](#tests)
- [Contributing](#contributing)
- [License](#license)

---

## Functions

### Ping

```go
Ping(string) error
```

Ping runs the `ansible all -m ping` command against
all nodes found in the provided hosts file by using the
mage/sh package to execute the command. If the command
execution fails, an error is returned.

**Parameters:**

hostsFile: A string representing the path to the hosts
file to be used by the `ansible` command.

**Returns:**

error: An error if the `ansible` command execution fails.

---

## Installation

To use the goutils/v2/ansible package, you first need to install it.
Follow the steps below to install via go get.

```bash
go get github.com/goutils/v2/l50/ansible
```

---

## Usage

After installation, you can import the package in your Go project
using the following import statement:

```go
import "github.com/goutils/v2/l50/ansible"
```

---

## Tests

To ensure the package is working correctly, run the following
command to execute the tests for `goutils/v2/ansible`:

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
