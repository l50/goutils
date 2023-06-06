# goutils/cloudflare

The `cloudflare` package is a part of `goutils` library. It provides 
utility functions to interface with cloudflore using go.

---

## Functions

### GetDNSRecords

```go
func GetDNSRecords(cf Cloudflare) error
```

Retrieve the DNS records from Cloudflare for a specified 
zone ID using the provided Cloudflare credentials.

---

## Installation

To use the `goutils/cloudflare` package, you need to install it via `go get`:

```bash
go get github.com/l50/goutils/cloudflare
```

---

## Usage

After installation, you can import it in your project:

```go
import "github.com/l50/goutils/cloudflare"
```

---

## Tests

To run the tests for the `goutils/cloudflare` package, navigate to
your `$GOPATH/src/github.com/l50/goutils/cloudflare` directory
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
