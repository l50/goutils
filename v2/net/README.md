# goutils/net

The `net` package is a part of `goutils` library. It provides
utility functions for network related operations in Go.

---

## Functions

### PublicIP

```go
func PublicIP(protocol uint) (string, error)
```

Obtains the public IP address of the system. This function uses
multiple external services to determine the IP address, taking as
input an integer representing the IP protocol version (4 or 6),
and returns the public IP address as a string. If the retrieval
fails, an error is returned.

---

## Installation

To use the `goutils/net` package, you need to install it via `go get`:

```bash
go get github.com/l50/goutils/v2/str
```

---

## Usage

After installation, you can import it in your project:

```go
import "github.com/l50/goutils/v2/net"
```

---

## Tests

To run the tests for the `goutils/net` package, navigate to
your `$GOPATH/src/github.com/l50/goutils/v2/net` directory
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
