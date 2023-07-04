# goutils/v2/git

The `git` package is a collection of utility functions
designed to simplify common git tasks.

Table of contents:

- [Functions](#functions)
- [Installation](#installation)
- [Usage](#usage)
- [Tests](#tests)
- [Contributing](#contributing)
- [License](#license)

---

## Functions

### AddFile(string)

```go
AddFile(string) error
```

AddFile adds the file located at the given file path to its
affiliated Git repository. An error is returned if an issue
happens during this process.

**Parameters:**

filePath: A string indicating the path to the file that
will be staged.

**Returns:**

error: An error if any occurs during the staging process.

---

### CloneRepo(string, string, transport.AuthMethod)

```go
CloneRepo(string, string, transport.AuthMethod) *git.Repository, error
```

CloneRepo clones a Git repository from the specified URL to
the given path, using the supplied authentication method, if
provided.

**Parameters:**

url: A string indicating the URL of the repository to clone.
clonePath: A string representing the path where the repository
will be cloned.
auth: A transport.AuthMethod interface symbolizing the
authentication method for cloning. If nil, no authentication is used.

**Returns:**

*git.Repository: A pointer to the Repository struct
representing the cloned repository.

error: An error if the repository can't be cloned or already
exists at the target path.

---

### Commit(*git.Repository, string)

```go
Commit(*git.Repository, string) error
```

Commit generates a new commit in the specified repository with
the given message. The commit's author is extracted from the
global Git user settings.

**Parameters:**

repo: A pointer to the Repository struct symbolizing the
repository where the commit should be made.
msg: A string depicting the commit message.

**Returns:**

error: An error if the commit can't be created.

---

### CreateTag(*git.Repository, string)

```go
CreateTag(*git.Repository, string) error
```

CreateTag forms a new tag in the specified repository.

**Parameters:**

repo: Pointer to the Repository struct, the repository where the tag is created.
tag: String, the name of the tag to create.

**Returns:**

error: Error if the tag can't be created, already exists, or if the global git
user settings can't be retrieved.

---

### DeletePushedTag(*git.Repository, string, transport.AuthMethod)

```go
DeletePushedTag(*git.Repository, string, transport.AuthMethod) error
```

DeletePushedTag deletes a tag from a repository that has been pushed.

**Parameters:**

repo: Repository where the tag should be deleted.
tag: The tag that should be deleted.
auth: Authentication method for the push.

**Returns:**

error: Error if the tag cannot be deleted.

---

### DeleteTag(*git.Repository, string)

```go
DeleteTag(*git.Repository, string) error
```

DeleteTag deletes the local input tag from the specified repo.

**Parameters:**

repo: Repository where the tag should be deleted.
tag: The tag that should be deleted.

**Returns:**

error: Error if the tag cannot be deleted.

---

### GetGlobalUserCfg()

```go
GetGlobalUserCfg() ConfigUserInfo, error
```

GetGlobalUserCfg fetches the username and email from the global git user
settings. It returns a ConfigUserInfo struct containing the global git
username and email. An error is returned if the global git username or email
cannot be retrieved.

**Returns:**

ConfigUserInfo: Struct containing the global git username and email.

error: Error if the global git username or email can't be retrieved.

---

### GetTags(*git.Repository)

```go
GetTags(*git.Repository) []string, error
```

GetTags returns all tags of the given repository.

**Parameters:**

repo: A pointer to the Repository struct representing
the repo from which tags are retrieved.

**Returns:**

[]string: A slice of strings, each representing a tag in the repository.
error: An error if a problem occurs while retrieving the tags.

---

### PullRepos(...string)

```go
PullRepos(...string) error
```

PullRepos updates all git repositories located in the specified directories.

**Parameters:**

dirs: Paths to directories to be searched for git repositories.

**Returns:**

error: Error if there's a problem with pulling the repositories.

---

### Push(*git.Repository, transport.AuthMethod)

```go
Push(*git.Repository, transport.AuthMethod) error
```

Push transmits the contents of the specified repository to the default
remote (origin).

**Parameters:**

repo: Pointer to the Repository struct, the repository to push.
auth: A transport.AuthMethod interface, the authentication method for the push.
If it's nil, no authentication is used.

**Returns:**

error: Error if the push fails.

---

### PushTag(*git.Repository, string, transport.AuthMethod)

```go
PushTag(*git.Repository, string, transport.AuthMethod) error
```

PushTag pushes a specific tag of the given repository to the default remote.

**Parameters:**

repo: Repository where the tag should be pushed.
tag: Name of the tag to push.
auth: Authentication method for the push. If nil, no authentication is used.

**Returns:**

error: Error if the push fails.

---

### RepoRoot()

```go
RepoRoot() string, error
```

RepoRoot finds and returns the root directory of the current Git repository.

**Returns:**

string: Absolute path to the root directory of the current Git repository.
error: Error if the Git repository root cannot be found.

---

## Installation

To use the goutils/v2/git package, you first need to install it.
Follow the steps below to install via go get.

```bash
go get github.com/goutils/v2/l50/git
```

---

## Usage

After installation, you can import the package in your Go project
using the following import statement:

```go
import "github.com/goutils/v2/l50/git"
```

---

## Tests

To ensure the package is working correctly, run the following
command to execute the tests for `goutils/v2/git`:

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
