# goutils/ansible

The `ansible` package is a part of `goutils` library. It provides
utility functions to interact with ansible in Go.

---

## Functions

### Ping

```go
func Ping(hostsFile string) error {
```

Runs the `ansible all -m ping` command against all
nodes found in the provided hosts file. If a host file is not provided,
localhost is used by default.

---

## Installation

To use the `goutils/ansible` package, you need to install it via `go get`:

```bash
go get github.com/l50/goutils/v2/ansible
```

---

## Usage

After installation, you can import it in your project:

```go
import "github.com/l50/goutils/v2/ansible"
```

---

## Tests

To run the tests for the `goutils/ansible` package, navigate to
your `$GOPATH/src/github.com/l50/goutils/v2/ansible` directory
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
