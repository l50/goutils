# goutils/logging

The `logging` package is a part of `goutils` library. It provides utility 
functions for logging operations in Go.

---

## Functions

### CreateLogFile

```go
func CreateLogFile(logDir string, logName string) (LogInfo, error)
```

Create a log file in a directory named logs, which is a 
subdirectory of the given directory.

---

## Installation

To use the `goutils/logging` package, you need to install it via `go get`:

```bash
go get github.com/l50/goutils/logging
```

---

## Usage

After installation, you can import it in your project:

```go
import "github.com/l50/goutils/logging"
```

---

## Tests

To run the tests for the `goutils/logging` package, navigate to
your `$GOPATH/src/github.com/l50/goutils/logging` directory
and run go test:

```bash
go test -v
```

---

## Contributing

Pull requests are welcome. For major changes, please
open an issue first to discuss what you would like to change.

---

## License

This project is licensed under the MIT License - see
the [LICENSE](../../LICENSE) file for details.
