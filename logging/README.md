# goutils/v2/logging

The `logging` package is a collection of utility functions
designed to simplify common logging tasks.

Table of contents:

- [Functions](#functions)
- [Installation](#installation)
- [Usage](#usage)
- [Tests](#tests)
- [Contributing](#contributing)
- [License](#license)

---

## Functions

### CreateLogFile

```go
CreateLogFile(afero.Fs, string, string) LogInfo, error
```

CreateLogFile creates a log file in a 'logs' subdirectory of the
specified directory. The log file's name is the provided log name
with the extension '.log'.

**Parameters:**

fs: An afero.Fs instance to mock filesystem for testing.
logDir: A string for the directory where 'logs' subdirectory and
log file should be created.
logName: A string for the name of the log file to be created.

**Returns:**

LogInfo: A LogInfo struct with information about the log file,
including its directory, file pointer, file name, and path.
error: An error, if an issue occurs while creating the directory
or the log file.

---

## Installation

To use the goutils/v2/logging package, you first need to install it.
Follow the steps below to install via go get.

```bash
go get github.com/goutils/v2/l50/logging
```

---

## Usage

After installation, you can import the package in your Go project
using the following import statement:

```go
import "github.com/goutils/v2/l50/logging"
```

---

## Tests

To ensure the package is working correctly, run the following
command to execute the tests for `goutils/v2/logging`:

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
