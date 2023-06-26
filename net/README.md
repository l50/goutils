# goutils/v2/net

The `net` package is a collection of utility functions
designed to simplify common net tasks.

Table of contents:

- [Functions](#functions)
- [Installation](#installation)
- [Usage](#usage)
- [Tests](#tests)
- [Contributing](#contributing)
- [License](#license)

---

## Functions

### DownloadFile

```go
DownloadFile(string, string) string, error
```

DownloadFile downloads a file from the provided URL and saves it
to the specified location on the local filesystem. The function
takes the source URL and the destination path as inputs and returns the path
where the file was saved or an error.

Parameters:

url: A string representing the URL of the file to be downloaded.
dest: A string representing the destination path where the file
should be saved on the local filesystem.

Returns:

string: The path where the downloaded file was saved.
error: An error if the function fails to download the file.

---

### PublicIP

```go
PublicIP(uint) string, error
```

PublicIP uses several external services to get the public
IP address of the system running it, using github.com/GlenDC/go-external-ip.
The function takes an IP protocol version (4 or 6) as input and
returns the public IP address as a string or an error.

Parameters:

protocol: A uint representing the IP protocol version (4 or 6).

Returns:

string: The public IP address of the system in string format.
error: An error if the function fails to retrieve the public IP address.

---

## Installation

To use the goutils/v2/net package, you first need to install it.
Follow the steps below to install via go get.

```bash
go get github.com/goutils/v2/l50/net
```

---

## Usage

After installation, you can import the package in your Go project
using the following import statement:

```go
import "github.com/goutils/v2/l50/net"
```

---

## Tests

To ensure the package is working correctly, run the following
command to execute the tests for `goutils/v2/net`:

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
