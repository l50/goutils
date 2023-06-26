# goutils/v2/str

The `str` package is a collection of utility functions
designed to simplify common str tasks.

Table of contents:

- [Functions](#functions)
- [Installation](#installation)
- [Usage](#usage)
- [Tests](#tests)
- [Contributing](#contributing)
- [License](#license)

---

## Functions

### GenRandom

```go
GenRandom(int) string, error
```

GenRandom generates a random string of a specified length.

**Parameters:**

length: Length of the random string to be generated.

**Returns:**

string: Generated random string.
error: An error if random string generation fails.

---

### InSlice

```go
InSlice(string, []string) bool
```

InSlice determines if a specified string exists in a given slice.

**Parameters:**

strToFind: String to search for in the slice.
inputSlice: Slice of strings to be searched.

**Returns:**

bool: true if string is found in the slice, false otherwise.

---

### IsNumeric

```go
IsNumeric(string) bool
```

IsNumeric checks if a string is entirely composed of numeric characters.

**Parameters:**

s: String to check for numeric characters.

**Returns:**

bool: true if the string is numeric, false otherwise.

---

### SlicesEqual

```go
SlicesEqual([]string) bool
```

SlicesEqual compares two slices of strings for equality.

**Parameters:**

a: First string slice for comparison.
b: Second string slice for comparison.

**Returns:**

bool: true if slices are equal, false otherwise.

---

### ToInt64

```go
ToInt64(string) int64, error
```

ToInt64 converts a string to int64.

**Parameters:**

value: String to be converted to int64.

**Returns:**

int64: int64 equivalent of the string.
error: An error if the conversion fails.

---

### ToSlice

```go
ToSlice(string, string) []string
```

ToSlice converts a string to a slice of strings using a delimiter.

**Parameters:**

delimStr: String to be split into a slice.
delim: Delimiter to be used for splitting the string.

**Returns:**

[]string: Slice of strings from the split input string.

---

## Installation

To use the goutils/v2/str package, you first need to install it.
Follow the steps below to install via go get.

```bash
go get github.com/goutils/v2/l50/str
```

---

## Usage

After installation, you can import the package in your Go project
using the following import statement:

```go
import "github.com/goutils/v2/l50/str"
```

---

## Tests

To ensure the package is working correctly, run the following
command to execute the tests for `goutils/v2/str`:

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
