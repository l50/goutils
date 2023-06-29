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

### ConfigureLogger

```go
ConfigureLogger(slog.Level, string) Logger, error
```

ConfigureLogger creates a logger based on the provided level.
Depending on the level, it returns a colored or plain logger.

**Parameters:**

level: Logging level as a slog.Level.
path: Path to the log file.

**Returns:**

Logger: Logger object based on provided level.
error: An error, if an issue occurs while setting up the logger.

---

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

### Error

```go
Error(...interface{})
```

Error for PlainLogger logs the provided arguments as an error line in plain text.
The arguments are handled in the manner of fmt.Println.

---

### Error

```go
Error(...interface{})
```

Error for ColoredLogger logs the provided arguments as an error line
in the specified color. The arguments are handled in the manner
of fmt.Println.

---

### Error

```go
Error(...interface{})
```

Error for SlogPlainLogger logs the provided arguments as an error line using slog library.
The arguments are converted to a string using fmt.Sprint.

---

### Error

```go
Error(...interface{})
```

Error for SlogLogger logs the provided arguments as an error line using slog library.
The arguments are converted to a string using fmt.Sprint.

---

### Errorf

```go
Errorf(string, ...interface{})
```

Errorf for SlogLogger logs the provided formatted string as an error line using slog library.
The format and arguments are handled in the manner of fmt.Printf.

---

### Errorf

```go
Errorf(string, ...interface{})
```

Errorf for PlainLogger logs the provided formatted string as an error line in plain text.
The format and arguments are handled in the manner of fmt.Printf.

---

### Errorf

```go
Errorf(string, ...interface{})
```

Errorf for ColoredLogger logs the provided formatted string as an error line in
the specified color. The format and arguments are handled in the
manner of fmt.Printf.

---

### Errorf

```go
Errorf(string, ...interface{})
```

Errorf for SlogPlainLogger logs the provided formatted string as an error line using slog library.
The format and arguments are handled in the manner of fmt.Printf.

---

### Printf

```go
Printf(string, ...interface{})
```

Printf for PlainLogger logs the provided formatted string in plain text.
The format and arguments are handled in the manner of fmt.Printf.

---

### Printf

```go
Printf(string, ...interface{})
```

Printf for SlogLogger logs the provided formatted string using slog library.
The format and arguments are handled in the manner of fmt.Printf.

---

### Printf

```go
Printf(string, ...interface{})
```

Printf for SlogPlainLogger logs the provided formatted string using slog library.
The format and arguments are handled in the manner of fmt.Printf.

---

### Printf

```go
Printf(string, ...interface{})
```

Printf for ColoredLogger logs the provided formatted string in
the specified color. The format and arguments are handled in the
manner of fmt.Printf.

---

### Println

```go
Println(...interface{})
```

Println for SlogLogger logs the provided arguments as a line using slog library.
The arguments are converted to a string using fmt.Sprint.

---

### Println

```go
Println(...interface{})
```

Println for PlainLogger logs the provided arguments as a line in plain text.
The arguments are handled in the manner of fmt.Println.

---

### Println

```go
Println(...interface{})
```

Println for ColoredLogger logs the provided arguments as a line
in the specified color. The arguments are handled in the manner
of fmt.Println.

---

### Println

```go
Println(...interface{})
```

Println for SlogPlainLogger logs the provided arguments as a line using slog library.
The arguments are converted to a string using fmt.Sprint.

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
