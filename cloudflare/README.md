# goutils/v2/cloudflare

The `cloudflare` package is a collection of utility functions
designed to simplify common cloudflare tasks.

Table of contents:

- [Functions](#functions)
- [Installation](#installation)
- [Usage](#usage)
- [Tests](#tests)
- [Contributing](#contributing)
- [License](#license)

---

## Functions

### GetDNSRecords(Cloudflare)

```go
GetDNSRecords(Cloudflare) error
```

GetDNSRecords retrieves the DNS records from Cloudflare for a
specified zone ID using the provided Cloudflare credentials.
It makes a GET request to the Cloudflare API, reads the
response, and prints the 'name' and 'content' fields of
each DNS record.

**Parameters:**

cf: A Cloudflare struct containing the necessary credentials
(email, API key) and the zone ID for which the DNS records
should be retrieved.

**Returns:**

error: An error if any issue occurs while trying to
get the DNS records.

---

## Installation

To use the goutils/v2/cloudflare package, you first need to install it.
Follow the steps below to install via go get.

```bash
go get github.com/goutils/v2/l50/cloudflare
```

---

## Usage

After installation, you can import the package in your Go project
using the following import statement:

```go
import "github.com/goutils/v2/l50/cloudflare"
```

---

## Tests

To ensure the package is working correctly, run the following
command to execute the tests for `goutils/v2/cloudflare`:

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
