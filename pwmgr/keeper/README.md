# goutils/v2/keeper

The `keeper` package is a collection of utility functions
designed to simplify common keeper tasks.

Table of contents:

- [Functions](#functions)
- [Installation](#installation)
- [Usage](#usage)
- [Tests](#tests)
- [Contributing](#contributing)
- [License](#license)

---

## Functions

### Keeper.AddRecord(map[string]string)

```go
AddRecord(map[string]string) error
```

AddRecord adds a new record to the Keeper vault.

**Parameters:**

fields: A map containing the record fields.

fields.title: The title of the record.
fields.login: The username or login of the record.
fields.password: The password of the record.
fields.notes: Additional notes related to the record.

**Returns:**

error: An error if the record cannot be added.

---

### Keeper.CommanderInstalled()

```go
CommanderInstalled() bool
```

CommanderInstalled checks if the Keeper Commander tool is
installed on the current system.

**Returns:**

bool: True if the Keeper Commander tool is installed, false otherwise.

---

### Keeper.LoggedIn()

```go
LoggedIn() bool
```

LoggedIn checks if the user is logged into their Keeper vault.

**Returns:**

bool: True if the user is logged into their Keeper vault, false otherwise.

---

### Keeper.RetrieveRecord(string)

```go
RetrieveRecord(string) pwmgr.Record, error
```

RetrieveRecord retrieves a user's Keeper record using the
provided unique identifier (uid).

**Parameters:**

uid: A string representing the unique identifier of the
Keeper record to retrieve.

**Returns:**

pwmgr.Record: The retrieved Keeper record. This contains the following attributes:
- UID: The unique identifier of the record.
- Title: The title of the record.
- Username: The username associated with the record.
- Password: The password of the record.
- URL: The URL associated with the record.
- TOTP: The one-time password (if any) associated with the record.
- Note: Any additional notes associated with the record.

error: An error if the Keeper record cannot be retrieved.

---

### Keeper.SearchRecords(string)

```go
SearchRecords(string) string, error
```

SearchRecords searches the user's Keeper records for records
that match the provided search term.

**Parameters:**

searchTerm: A string representing the term to search for in the Keeper records.

**Returns:**

string: The unique identifier (UID) of the first Keeper record
that matches the search term. If multiple records match the
search term, only the UID of the first record is returned.

error: An error if the Keeper records cannot be searched or if
the search term does not match any records.

---

## Installation

To use the goutils/v2/keeper package, you first need to install it.
Follow the steps below to install via go get.

```bash
go get github.com/goutils/v2/l50/keeper
```

---

## Usage

After installation, you can import the package in your Go project
using the following import statement:

```go
import "github.com/goutils/v2/l50/keeper"
```

---

## Tests

To ensure the package is working correctly, run the following
command to execute the tests for `goutils/v2/keeper`:

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
