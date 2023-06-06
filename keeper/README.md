# goutils/keeper

The `keeper` package is a part of `goutils` library.

It provides utility functions for interfacing with the Keeper
password manager in Go.

---

## Functions

### CommanderInstalled

```go
func CommanderInstalled() bool
```

Check if the Keeper Commander tool is installed on the current system.

### LoggedIn

```go
func LoggedIn() bool
```

Check if the user is logged in to Keeper.

### RetrieveRecord

```go
func RetrieveRecord(keeperUID string) (Record, error)
```

Retrieves a user's Keeper record using the provided unique identifier (keeperUID).

### SearchRecords

```go
func SearchRecords(searchTerm string) (string, error)
```

Searches the user's Keeper records for records that match the provided search term.

---

## Installation

To use the `goutils/keeper` package, you need to install it via `go get`:

```bash
go get github.com/l50/goutils/keeper
```

---

## Usage

After installation, you can import it in your project:

```go
import "github.com/l50/goutils/keeper"
```

---

## Tests

To run the tests for the `goutils/keeper` package, navigate to
your `$GOPATH/src/github.com/l50/goutils/keeper` directory
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
the [LICENSE](../LICENSE) file for details.
