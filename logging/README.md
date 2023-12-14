# goutils/v2/logging

The `logging` package is a collection of utility functions
designed to simplify common logging tasks.

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

### ColorLogger.Debug(...interface{})

```go
Debug(...interface{})
```

Debug for ColorLogger logs the provided arguments as a debug line
in the specified color. The arguments are handled in the manner
of fmt.Println.

---

### ColorLogger.Debugf(string, ...interface{})

```go
Debugf(string, ...interface{})
```

Debugf for ColorLogger logs the provided formatted string as a debug
line in the specified color. The format and arguments are handled
in the manner of fmt.Printf.

---

### ColorLogger.Error(...interface{})

```go
Error(...interface{})
```

Error for ColorLogger logs the provided arguments as an error line
in the specified color. The arguments are handled in the manner
of fmt.Println.

---

### ColorLogger.Errorf(string, ...interface{})

```go
Errorf(string, ...interface{})
```

Errorf for ColorLogger logs the provided formatted string as an
error line in the specified color. The format and arguments are handled
in the manner of fmt.Printf.

---

### ColorLogger.Printf(string, ...interface{})

```go
Printf(string, ...interface{})
```

Printf for ColorLogger logs the provided formatted string in
the specified color. The format and arguments are handled in the
manner of fmt.Printf.

---

### ColorLogger.Println(...interface{})

```go
Println(...interface{})
```

Println for ColorLogger logs the provided arguments as a line
in the specified color. The arguments are handled in the manner
of fmt.Println.

---

### ConfigureLogger(slog.Level, string, OutputType)

```go
ConfigureLogger(slog.Level, string, OutputType) Logger, error
```

ConfigureLogger sets up a logger based on the provided logging level,
file path, and output type. It supports both colorized and plain text
logging output, selectable via the OutputType parameter. The logger
writes log entries to both a file and standard output.

**Parameters:**

level: Logging level as a slog.Level.
path: Path to the log file.
outputType: Type of log output, either ColorOutput or PlainOutput.

**Returns:**

Logger: Configured Logger object based on provided parameters.
error: An error, if an issue occurs while setting up the logger.

---

### CreateLogFile(afero.Fs, string)

```go
CreateLogFile(afero.Fs, string) LogInfo, error
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

### L()

```go
L() Logger
```

L returns the global logger instance for use in logging operations.

**Returns:**

Logger: The global Logger instance.

---

### NewPrettyHandler(io.Writer, PrettyHandlerOptions)

```go
NewPrettyHandler(io.Writer, PrettyHandlerOptions) *PrettyHandler
```

NewPrettyHandler creates a new PrettyHandler with specified output
writer and options. It configures a PrettyHandler for colorized
logging output.

**Parameters:**

out: Output writer where log messages will be written.
opts: PrettyHandlerOptions for configuring the handler.

**Returns:**

*PrettyHandler: A new instance of PrettyHandler.

---

### PlainLogger.Debug(...interface{})

```go
Debug(...interface{})
```

Debug for PlainLogger logs the provided arguments as a debug line
using slog library.
The arguments are converted to a string using fmt.Sprint.

---

### PlainLogger.Debugf(string, ...interface{})

```go
Debugf(string, ...interface{})
```

Debugf for PlainLogger logs the provided formatted string as a debug
line using slog library.
The format and arguments are handled in the manner of fmt.Printf.

---

### PlainLogger.Error(...interface{})

```go
Error(...interface{})
```

Error for PlainLogger logs the provided arguments as an error line
using slog library.
The arguments are converted to a string using fmt.Sprint.

---

### PlainLogger.Errorf(string, ...interface{})

```go
Errorf(string, ...interface{})
```

Errorf for PlainLogger logs the provided formatted string as an error
line using slog library.
The format and arguments are handled in the manner of fmt.Printf.

---

### PlainLogger.Printf(string, ...interface{})

```go
Printf(string, ...interface{})
```

Printf for PlainLogger logs the provided formatted string using slog library.
The format and arguments are handled in the manner of fmt.Printf.

---

### PlainLogger.Println(...interface{})

```go
Println(...interface{})
```

Println for PlainLogger logs the provided arguments as a line using
slog library.
The arguments are converted to a string using fmt.Sprint.

---

### PrettyHandler.Handle(context.Context, slog.Record)

```go
Handle(context.Context, slog.Record) error
```

Handle formats and outputs a log message for PrettyHandler. It
colorizes the log level, message, and adds structured fields
to the log output.

**Parameters:**

ctx: Context for the log record.
r: The log record containing log data.

**Returns:**

error: An error if any issue occurs during log handling.

---

## Installation

To use the goutils/v2/logging package, you first need to install it.
Follow the steps below to install via go get.

```bash
go get github.com/l50/goutils/v2/logging
```

---

## Usage

After installation, you can import the package in your Go project
using the following import statement:

```go
import "github.com/l50/goutils/v2/logging"
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
