# goutils/v2/file

The `file` package is a collection of utility functions
designed to simplify common file tasks.

---

## Table of contents

- [Functions](#functions)
- [Installation](#installation)
- [Usage](#usage)
- [Tests](#tests)
- [Contributing](#contributing)
- [License](#license)

---

## Functions

### CSVToLines(string)

```go
CSVToLines(string) [][]string, error
```

CSVToLines reads a CSV file and returns it as a 2D string slice. Each
element in the outer slice represents a row in the CSV, each element in the
inner slice represents a value in that row. The first row of the CSV,
assumed to contain column headers, is skipped.

**Parameters:**

path: String representing the path to the CSV file.

**Returns:**

[][]string: 2D slice of strings representing the rows and values of the CSV.
error: An error if the file cannot be read or parsed.

---

### Create(string, []byte, CreateType)

```go
Create(string, []byte, CreateType) error
```

Create makes a directory, an empty file, a file with content, or a temporary file at
the specified path, depending on the createType argument.

**Parameters:**

path: Path to the directory or file. For temporary files, this serves as a pattern.
contents: Content to write to the file as a byte slice.
createType: A CreateType value representing what kind of file creation action to execute.

**Returns:**

error: An error if the directory or file can't be created, if it
already exists (except for temporary files), or if there's a problem writing to the file.

---

### Delete(string)

```go
Delete(string) error
```

Delete removes the specified file.

**Parameters:**

path: String representing the path to the file.

**Returns:**

error: An error if the file cannot be deleted.

---

### Exists(string)

```go
Exists(string) bool
```

Exists checks whether a file at the specified path exists.

**Parameters:**

fileLoc: String representing the path to the file.

**Returns:**

bool: Returns true if the file exists, otherwise false.

---

### Find(string, []string)

```go
Find(string, []string) []string, error
```

Find searches for a specified filename in a set of directories and returns
all matches found as a slice of file paths. If no matches are found, it
returns an error.

**Parameters:**

fileName: Name of the file to find.
dirs: Slice of strings representing the directories to search in.

**Returns:**

[]string: Slice of file paths if the file is found.
error: An error if the file cannot be found.

---

### HasStr(string, string)

```go
HasStr(string, string) bool, error
```

HasStr checks for the presence of a string in a specified file.

**Parameters:**

path: String representing the path to the file.
searchStr: String to look for in the file.

**Returns:**

bool: Returns true if the string is found, otherwise false.
error: An error if the file cannot be read.

---

### ListR(string)

```go
ListR(string) []string, error
```

ListR lists all files in a directory and its subdirectories.

**Parameters:**

dirPath: String representing the path to the directory.

**Returns:**

[]string: Slice of strings representing the paths of the files found.
error: An error if the files cannot be listed.

---

### RealFile.Append(string)

```go
Append(string) error
```

Append adds a string to the end of a file. If the file
doesn't exist, it's created with the default permissions.

**Parameters:**

text: String to append to the end of the file.

**Returns:**

error: An error if the file can't be opened or the string can't be
written to the file.

---

### RealFile.Open()

```go
Open() io.ReadCloser, error
```

Open is a method for the RealFile type that opens the file and
returns a io.ReadCloser and an error.

**Returns:**

io.ReadCloser: An object that allows reading from and closing the file.
error: An error if any issue occurs while trying to open the file.

---

### RealFile.Remove()

```go
Remove() error
```

Remove is a method for the RealFile type that removes the specified file or directory.
Note that it will not remove a directory unless it is empty.

**Parameters:**

name: A string representing the path to the file or directory to remove.

**Returns:**

error: An error if any issue occurs while trying to remove the file or directory.

---

### RealFile.RemoveAll()

```go
RemoveAll() error
```

RemoveAll is a method for the RealFile type that removes
a file or directory at the specified path.
If the path represents a directory, RemoveAll will remove
the directory and all its content.

**Parameters:**

path: A string representing the path to the file or directory to remove.

**Returns:**

error: An error if any issue occurs while trying to remove the file or directory.

---

### RealFile.Stat()

```go
Stat() os.FileInfo, error
```

Stat is a method for the RealFile type that retrieves the
FileInfo for the specified file or directory.

**Parameters:**

name: A string representing the path to the file or directory.

**Returns:**

os.FileInfo: FileInfo describing the named file.
error: An error if any issue occurs while trying to get the FileInfo.

---

### RealFile.Write([]byte, os.FileMode)

```go
Write([]byte, os.FileMode) error
```

Write is a method for the RealFile type that writes a slice of bytes
to the file with specified file permissions.

**Parameters:**

contents: A slice of bytes that should be written to the file.
mode: File permissions to use when creating the file.

**Returns:**

error: An error if any issue occurs while trying to write to the file.

---

### SeekAndDestroy(string, string)

```go
SeekAndDestroy(string, string) error
```

SeekAndDestroy walks through a directory and deletes all files that match the pattern

**Parameters:**

path: String representing the path to the directory.
pattern: String representing the pattern to match.

**Returns:**

error: An error if the files cannot be deleted.

---

### ToSlice(string)

```go
ToSlice(string) []string, error
```

ToSlice reads a file and returns its content as a slice of strings, each
element represents a line in the file. Blank lines are omitted.

**Parameters:**

path: String representing the path to the file.

**Returns:**

[]string: Slice of strings where each element represents a line in the file.
error: An error if the file cannot be read.

---

### WriteTempFile(string, *bytes.Buffer)

```go
WriteTempFile(string, *bytes.Buffer) string, error
```

WriteTempFile creates a temporary file in the system default temp directory,
writes the contents from the provided buffer, and returns the file path.

**Parameters:**

workloadName: A string representing the base name of the temporary file.
jobFile: A *bytes.Buffer containing the data to write to the temporary file.

**Returns:**

string: The name of the temporary file created.
error: An error if any issue occurs during file creation or writing.

---

## Installation

To use the goutils/v2/file package, you first need to install it.
Follow the steps below to install via go get.

```bash
go get github.com/l50/goutils/v2/file
```

---

## Usage

After installation, you can import the package in your Go project
using the following import statement:

```go
import "github.com/l50/goutils/v2/file"
```

---

## Tests

To ensure the package is working correctly, run the following
command to execute the tests for `goutils/v2/file`:

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
