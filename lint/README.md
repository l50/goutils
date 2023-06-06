# goutils/lint

The `lint` package is a part of `goutils` library. It provides utility 
functions for lint operations in Go.

---

## Functions

### InstallGoPCDeps

```go
func InstallGoPCDeps() error
```

Install dependencies used for pre-commit with Golang projects.

### InstallPCHooks

```go
func InstallPCHooks() error
```

Installs pre-commit hooks locally.

### UpdatePCHooks

```go
func UpdatePCHooks() error
```

Update pre-commit hooks locally.

### ClearPCCache

```go
func ClearPCCache() error
```

Clears the pre-commit cache.

### RunPCHooks

```go
func RunPCHooks() error
```

Run all pre-commit hooks locally.

### AddFencedCB

```go
func AddFencedCB(filePath string, language string) error
```

Addresses MD040 issues found with markdownlint by adding the 
input language to fenced code blocks in the input filePath.


---

## Installation

To use the `goutils/lint` package, you need to install it via `go get`:

```bash
go get github.com/l50/goutils/lint
```

---

## Usage

After installation, you can import it in your project:

```go
import "github.com/l50/goutils/lint"
```

---

## Tests

To run the tests for the `goutils/lint` package, navigate to
your `$GOPATH/src/github.com/l50/goutils/lint` directory
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
