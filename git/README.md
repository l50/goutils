# goutils/v2/git

The `git` package is a part of `goutils` library.

It provides utility functions to interface with git and services
that provide git in Go.

---

## Functions

### AddFile

```go
func AddFile(filePath string) error
```

Stage the file at the given file path in its associated Git repository.
Returns an error if one occurs.

### Commit

```go
func Commit(repo *git.Repository, msg string) error
```

Create a new commit in the given repository with the provided message.
The author of the commit is retrieved from the global Git user settings.

### CloneRepo

```go
func CloneRepo(url string, clonePath string, auth transport.AuthMethod)
```

Clone the repo specified with the input url to provided clonePath.

### GetTags

```go
func GetTags(repo *git.Repository) ([]string, error)
```

Returns the tags for an input repo.

### GetGlobalUserCfg

```go
func GetGlobalUserCfg() (ConfigUserInfo, error)
```

Return the username and email from the global git user settings.

### CreateTag

```go
func CreateTag(repo *git.Repository, tag string) error
```

Create an input tag in the specified repo if it doesn't already exist.

### Push

```go
func Push(repo *git.Repository, auth transport.AuthMethod) error
```

Push the contents of the input repo to the default remote (origin).

### PushTag

```go
func PushTag(repo *git.Repository, tag string, auth transport.AuthMethod) error
```

Push a tag to remote.

### DeleteTag

```go
func DeleteTag(repo *git.Repository, tag string) error
```

Delete the local input tag from the specified repo.

### DeletePushedTag

```go
func DeletePushedTag(repo *git.Repository, tag string, auth transport.AuthMethod) error
```

Delete a tag from a given repository that has been pushed remote.
The tag is deleted from both the local repository and the remote repository.

### PullRepos

```go
func PullRepos(dirs ...string) error
```

Update all git repositories found in the given directories by pulling
changes from the upstream branch.

It looks for repositories by finding directories with a ".git" subdirectory.

If a repository is not on the default branch, it will switch to the
default branch before pulling changes.

---

## Installation

To use the `goutils/v2/git` package, you need to install it via `go get`:

```bash
go get github.com/l50/goutils/v2/git
```

---

## Usage

After installation, you can import it in your project:

```go
import "github.com/l50/goutils/v2/git"
```

---

## Tests

To run the tests for the `goutils/v2/git` package, navigate to
your `$GOPATH/src/github.com/l50/goutils/v2/git` directory
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
