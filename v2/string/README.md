# goutils/string

The string package is a part of goutils library, which provides
utility functions for string manipulation in Go.

---

## Functions

### RandomString

```go
func RandomString(length int) (string, error)
```

Generates a random string of the specified length. This function
takes an integer input representing the length and returns a
string of hexadecimal characters. If the generation fails, an error is returned.

### StringInSlice

```go
func StringInSlice(strToFind string, inputSlice []string) bool
```

Checks if a specific string exists in a given slice. It returns true
if the string is found and false otherwise.

### StringToInt64

```go
func StringToInt64(value string) (int64, error)
```

Converts a string to an int64 value. If the conversion fails, an error is returned.

### StringToSlice

```go
func StringToSlice(delimStr string, delim string) []string
```

Splits a given string into a slice based on the provided delimiter.

### StringSlicesEqual

```go
func StringSlicesEqual(a, b []string) bool
```

Compares two string slices for equality. It returns true if the slices
have the same length and contain the same strings in the same order.
It returns false otherwise.

---

## Installation

To use the `goutils/string` package, you need to install it via go get:

```bash
go get github.com/l50/goutils/v2/string
```

---

## Usage

After installation, you can import it in your project:

```go
import "github.com/l50/goutils/v2/string"
```

---

## Tests

To run the tests for the `goutils/string` package, navigate to
your `$GOPATH/src/github.com/l50/goutils/v2/string` directory
and run go test:

```bash
go test -v
```

---

## Contributing

Pull requests are welcome. For major changes, please open an
issue first to discuss what you would like to change.

---

## License

This project is licensed under the MIT License - see the [LICENSE](../../LICENSE)
file for details.
