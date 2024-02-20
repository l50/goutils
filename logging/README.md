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

### ColorLogger.Warn(...interface{})

```go
Warn(...interface{})
```

Warn for ColorLogger logs the provided arguments as a warning line
in the specified color. The arguments are handled in the manner of fmt.Println.

---

### ColorLogger.Warnf(string, ...interface{})

```go
Warnf(string, ...interface{})
```

Warnf for ColorLogger logs the provided formatted string as a warning
line in the specified color. The format and arguments are handled in the
manner of fmt.Printf.

---

### InitLogging(*LogConfig)

```go
InitLogging(*LogConfig) Logger, error
```

InitLogging is a convenience function that combines
the CreateLogFile and ConfigureLogger functions into one call.
It is useful for quickly setting up logging to disk.

**Parameters:**

fs: An afero.Fs instance for filesystem operations, allows mocking in tests.
logPath: The path to the log file.
level: The logging level.
outputType: The output type of the logger (PlainOutput or ColorOutput).
logToDisk: A boolean indicating whether to log to disk or not.

**Returns:**

Logger: A configured Logger object.
error: An error if any issue occurs during initialization.

---

### L()

```go
L() Logger
```

L returns the global logger instance for use in logging operations.

**Returns:**

Logger: The global Logger instance.

---

### LogAndReturnError(Logger, string)

```go
LogAndReturnError(Logger, string) error
```

LogAndReturnError logs the provided error message using the given logger and returns the error.

This utility function is helpful for scenarios where an error needs to be both logged and returned.
It simplifies the code by combining these two actions into one call.

**Parameters:**

logger: The Logger instance used for logging the error.
errMsg: The error message to log and return.

**Returns:**

error: The error created from the errMsg, after it has been logged.

---

### LogConfig.ConfigureLogger()

```go
ConfigureLogger() Logger, error
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

### LogConfig.CreateLogFile()

```go
CreateLogFile() error
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

LogConfig: A LogConfig struct with information about the log file,
including its directory, file pointer, file name, and path.
error: An error, if an issue occurs while creating the directory
or the log file.

---

### NewColorLogger(LogConfig, color.Attribute, *slog.Logger)

```go
NewColorLogger(LogConfig, color.Attribute, *slog.Logger) *ColorLogger, error
```

NewColorLogger creates a new ColorLogger instance with the specified
LogConfig, color attribute, and slog.Logger.

**Parameters:**

cfg: LogConfig object containing information about the log file.
colorAttr: A color attribute for output styling.
logger: The slog Logger instance used for logging operations.

**Returns:**

*ColorLogger: A new instance of ColorLogger.
error: An error if any issue occurs during initialization.

---

### NewPlainLogger(LogConfig, *slog.Logger)

```go
NewPlainLogger(LogConfig, *slog.Logger) *PlainLogger, error
```

NewPlainLogger creates a new PlainLogger instance with the specified
LogConfig and slog.Logger.

**Parameters:**

cfg: LogConfig object containing information about the log file.
logger: The slog Logger instance used for logging operations.

**Returns:**

*PlainLogger: A new instance of PlainLogger.
error: An error if any issue occurs during initialization.

---

### NewPrettyHandler(io.Writer, PrettyHandlerOptions)

```go
NewPrettyHandler(io.Writer, PrettyHandlerOptions) *PrettyHandler
```

NewPrettyHandler creates a new PrettyHandler with specified output
writer and options. It configures the PrettyHandler for handling
log messages with optional colorization and structured formatting.

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
PlainLogger.go

---

### PlainLogger.Warn(...interface{})

```go
Warn(...interface{})
```

Warn for PlainLogger logs the provided arguments as a warning line
using slog library.
The arguments are converted to a string using fmt.Sprint.

---

### PlainLogger.Warnf(string, ...interface{})

```go
Warnf(string, ...interface{})
```

Warnf for PlainLogger logs the provided formatted string as a warning
line using slog library.
The format and arguments are handled in the manner of fmt.Printf.

---

### PrettyHandler.Handle(context.Context, slog.Record)

```go
Handle(context.Context, slog.Record) error
```

Handle processes and outputs a log record using the PrettyHandler.
It supports both colorized and non-colorized log messages and can
output in JSON format if not writing to a terminal.

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
