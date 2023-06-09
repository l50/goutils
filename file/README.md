# goutils/v2/file

The `file` package is a part of `goutils` library.

It provides utility functions for file manipulation in Go.

---

## Functions

### Append

```go
func Append(file string, text string) error
```

Appends an input text string to the end of the input file.

### CreateEmpty

```go
func CreateEmpty(name string) bool
```

Creates an empty file based on the name input. It returns true if the file was created,
otherwise it returns false.

### Create

```go
func Create(filePath string, fileContents []byte) error
```

Creates a file at the input filePath with the specified fileContents.

### CreateDirectory

```go
func CreateDirectory(path string) error
```

Creates a directory at the input path. If any part of the input path doesn't exist,
create it. Return an error if the path already exists.

### CSVToLines

```go
func CSVToLines(filename string) ([][]string, error)
```

Reads the contents of the specified CSV file and returns its contents as a two-dimensional
string slice.

### Delete

```go
func Delete(file string) error
```

Deletes the input file.

### Exists

```go
func Exists(fileLoc string) bool
```

Returns true if a file specified with fileLoc exists. If the file does not exist,
it returns false.

### ToSlice

```go
func ToSlice(fileName string) ([]string, error)
```

Reads an input file into a slice, removes blank strings, and returns it.

### Find

```go
func Find(fileName string, dirs []string) (string, error)
```

Looks for an input filename in the specified set of dirs. The filepath is returned
if the filename is found.

### ListR

```go
func ListR(path string) ([]string, error)
```

Lists the files found recursively from the input path.

### FindStr

```go
func FindStr(path string, searchStr string) (bool, error)
```

Searches for input searchStr in the input filepath.

---

## Installation

To use the `goutils/v2/file` package, you need to install it via `go get`:

```bash
go get github.com/l50/goutils/v2/file
```

---

## Usage

After installation, you can import it in your project:

```go
import "github.com/l50/goutils/v2/file"
```

---

## Tests

To run the tests for the `goutils/v2/file` package, navigate to
your `$GOPATH/src/github.com/l50/goutils/v2/file` directory and run go test:

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
