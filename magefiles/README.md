# goutils/v2/magefiles

`magefiles` provides utilities that would normally be managed
and executed with a `Makefile`. Instead of being written in the make language,
magefiles are crafted in Go and leverage the [Mage](https://magefile.org/) library.

---

## Table of contents

- [Functions](#functions)
- [Contributing](#contributing)
- [License](#license)

---

## Functions

### GeneratePackageDocs()

```go
GeneratePackageDocs() error
```

GeneratePackageDocs creates documentation for the various packages
in the project.

Example usage:

```go
mage generatepackagedocs
```

**Returns:**

error: An error if any issue occurs during documentation generation.

---

### InstallDeps()

```go
InstallDeps() error
```

InstallDeps installs the Go dependencies necessary for developing
on the project.

Example usage:

```go
mage installdeps
```

**Returns:**

error: An error if any issue occurs while trying to
install the dependencies.

---

### RunPreCommit()

```go
RunPreCommit() error
```

RunPreCommit updates, clears, and executes all pre-commit hooks
locally. The function follows a three-step process:

First, it updates the pre-commit hooks.
Next, it clears the pre-commit cache to ensure a clean environment.
Lastly, it executes all pre-commit hooks locally.

Example usage:

```go
mage runprecommit
```

**Returns:**

error: An error if any issue occurs at any of the three stages
of the process.

---

### RunTests()

```go
RunTests() error
```

RunTests executes all unit tests.

Example usage:

```go
mage runtests
```

**Returns:**

error: An error if any issue occurs while running the tests.

---

### UpdateMirror(string)

```go
UpdateMirror(string) error
```

UpdateMirror updates pkg.go.dev with the release associated with the
input tag

Example usage:

```go
mage updatemirror v2.0.1
```

**Parameters:**

tag: the tag to update pkg.go.dev with

**Returns:**

error: An error if any issue occurs while updating pkg.go.dev

---

### UseFixCodeBlocks(string, string)

```go
UseFixCodeBlocks(string, string) error
```

UseFixCodeBlocks fixes code blocks for the input filepath
using the input language.

Example usage:

```go
mage fixcodeblocks docs/docGeneration.go go
```

**Parameters:**

filepath: the path to the file or directory to fix

language: the language of the code blocks to fix

**Returns:**

error: an error if one occurred

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
