# goutils/v2/lint

The `lint` package is a collection of utility functions
designed to simplify common lint tasks.

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

### AddFencedCB(string, string)

```go
AddFencedCB(string, string) error
```

AddFencedCB addresses MD040 issues found with markdownlint by adding
the input language to fenced code blocks in the input filePath.

**Parameters:**

filePath: Path to the markdown file to modify.
language: Language to be added to the fenced code block.

**Returns:**

error: An error if the markdown file fails to be modified.

---

### ClearPCCache()

```go
ClearPCCache() error
```

ClearPCCache clears the pre-commit cache.

**Returns:**

error: An error if the cache fails to clear.

---

### InstallGoPCDeps()

```go
InstallGoPCDeps() error
```

InstallGoPCDeps installs dependencies used for pre-commit with Golang
projects.

**Returns:**

error: An error if the dependencies fail to install.

---

### InstallPCHooks()

```go
InstallPCHooks() error
```

InstallPCHooks installs pre-commit hooks locally.

**Returns:**

error: An error if the hooks fail to install.

---

### RunHookTool(string, ...string)

```go
RunHookTool(string, ...string) error
```

RunHookTool executes the specified pre-commit hook on a set of files.
It constructs a command to run 'pre-commit' with the given hook and
file arguments. If no files are provided, it defaults to "all".
The function then executes the command and handles any resulting error.

**Parameters:**

hook: A string specifying the name of the pre-commit hook to be run.
files: A variadic string slice containing file paths to be included
in the pre-commit hook execution. If no files are specified, it defaults
to running the hook on all files.

**Returns:**

error: An error if any issue occurs during the execution of the pre-commit
hook, otherwise nil if the hook runs successfully.

---

### RunPCHooks(...int)

```go
RunPCHooks(...int) error
```

RunPCHooks runs pre-commit hooks with a provided timeout.
If no timeout is provided, it defaults to 600.

**Parameters:**

timeout (optional): An integer specifying the timeout duration in seconds.

**Returns:**

error: An error if the pre-commit hook execution fails.

---

### UpdatePCHooks()

```go
UpdatePCHooks() error
```

UpdatePCHooks updates pre-commit hooks locally.

**Returns:**

error: An error if the hooks fail to update.

---

## Installation

To use the goutils/v2/lint package, you first need to install it.
Follow the steps below to install via go get.

```bash
go get github.com/l50/goutils/v2/lint
```

---

## Usage

After installation, you can import the package in your Go project
using the following import statement:

```go
import "github.com/l50/goutils/v2/lint"
```

---

## Tests

To ensure the package is working correctly, run the following
command to execute the tests for `goutils/v2/lint`:

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
