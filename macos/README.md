# goutils/macos

The `macos` package is a part of `goutils` library. It provides utility functions for 
macos oriented operations in Go.

---

## Functions

### InstallBrewDeps

```go
func InstallBrewDeps(brewPackages []string) error
```

Installs the input brew packages by running brew install.


### InstallBrewTFDeps

```go
func InstallBrewTFDeps() error
```

Install dependencies for terraform projects using homebrew.

---

## Installation

To use the `goutils/macos` package, you need to install it via `go get`:

```bash
go get github.com/l50/goutils/macos
```

---

## Usage

After installation, you can import it in your project:

```go
import "github.com/l50/goutils/macos"
```

---

## Tests

To run the tests for the `goutils/macos` package, navigate to
your `$GOPATH/src/github.com/l50/goutils/macos` directory
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
