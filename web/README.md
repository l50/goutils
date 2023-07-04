# goutils/v2/web

The `web` package is a collection of utility functions
designed to simplify common web tasks.

Table of contents:

- [Functions](#functions)
- [Installation](#installation)
- [Usage](#usage)
- [Tests](#tests)
- [Contributing](#contributing)
- [License](#license)

---

## Functions

### CancelAll(...func())

```go
CancelAll(...func())
```

CancelAll executes all provided cancel functions. It is typically used for
cleaning up or aborting operations that were started earlier and can be cancelled.

Note: The caller is responsible for handling any errors that may occur
during the execution of the cancel functions.

**Parameters:**

cancels: A slice of cancel functions, each of type func(). These are typically
functions returned by context.WithCancel, or similar functions that provide a
way to cancel an operation.

---

### GetRandomWait(int)

```go
GetRandomWait(int) time.Duration, error
```

GetRandomWait returns a random duration in seconds between the specified minWait
and maxWait durations. The function takes the minimum and maximum wait times as
arguments, creates a new random number generator with a seed based on the current
Unix timestamp, and calculates the random wait time within the given range.

**Parameters:**

minWait: The minimum duration to wait.
maxWait: The maximum duration to wait.

**Returns:**

time.Duration: A random duration between minWait and maxWait.
error: An error if the generation of the random wait time fails.

---

### IsLogMeOutEnabled(*LoginOptions)

```go
IsLogMeOutEnabled(*LoginOptions) bool
```

IsLogMeOutEnabled checks if the option to log out the user
after login is enabled in the provided login options.

**Parameters:**

opts: A pointer to a LoginOptions instance.

**Returns:**

bool: A boolean indicating whether the user is to be logged out after login.

---

### IsTwoFacEnabled(*LoginOptions)

```go
IsTwoFacEnabled(*LoginOptions) bool
```

IsTwoFacEnabled checks if two-factor authentication is enabled in the
provided login options.

**Parameters:**

opts: A pointer to a LoginOptions instance.

**Returns:**

bool: A boolean indicating whether two-factor authentication is enabled.

---

### SetLoginOptions(...LoginOption)

```go
SetLoginOptions(...LoginOption) *LoginOptions
```

SetLoginOptions applies provided login options to a new
LoginOptions instance with default values and returns a pointer
to this instance. This function is primarily used to configure login behavior
in the LoginAccount function.

**Parameters:**

options: A variadic set of LoginOption functions. Each LoginOption is a function
that takes a pointer to a LoginOptions struct and modifies it in place.

**Returns:**

*LoginOptions: A pointer to a LoginOptions struct that has been configured
with the provided options.

---

### Wait(float64)

```go
Wait(float64) time.Duration, error
```

Wait generates a random period of time anchored to a given input value.

**Parameters:**

near: A float64 value that serves as the base value for
generating the random wait time.

**Returns:**

time.Duration: The calculated random wait time in milliseconds.
error: An error if the generation of the random wait time fails.

The function is useful for simulating more human-like interaction
in the context of a web application. It first calculates a 'zoom' value by
dividing the input 'near' by 10. Then, a random number is generated in
the range of [0, zoom), and added to 95% of 'near'. This sum is then converted
to a time duration in milliseconds and returned.

---

### WithLogout(bool)

```go
WithLogout(bool) LoginOption
```

WithLogout is a function that returns a LoginOption
function which sets the logMeOut option.

**Parameters:**

enabled: Determines if the user should be logged out after login.

**Returns:**

LoginOption: A function that modifies the logMeOut
option of a LoginOptions struct.

---

### WithTwoFac(bool)

```go
WithTwoFac(bool) LoginOption
```

WithTwoFac is a function that returns a LoginOption function
which sets the twoFacEnabled option.

**Parameters:**

enabled: Determines if two-factor authentication should
be enabled during login.

**Returns:**

LoginOption: A function that modifies the twoFacEnabled
option of a LoginOptions struct.

---

## Installation

To use the goutils/v2/web package, you first need to install it.
Follow the steps below to install via go get.

```bash
go get github.com/goutils/v2/l50/web
```

---

## Usage

After installation, you can import the package in your Go project
using the following import statement:

```go
import "github.com/goutils/v2/l50/web"
```

---

## Tests

To ensure the package is working correctly, run the following
command to execute the tests for `goutils/v2/web`:

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
