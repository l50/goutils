# goutils/v2/cdpu

The `cdpu` package is a collection of utility functions
designed to simplify common cdpu tasks.

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

### CheckElement(web.Site, string, chan error)

```go
CheckElement(web.Site, string, chan error) error
```

CheckElement checks if a web page element, identified by the provided XPath,
exists within a specified timeout.

**Note:** Ensure to handle the error sent to the "done" channel in a
separate goroutine or after calling this function to avoid deadlock.

**Parameters:**

site: A web.Site struct representing the target site.
elementXPath: A string representing the XPath of the target element.
done: A channel through which the function sends an error if the
element is found or another error occurs.

**Returns:**

error: An error if the element is found, the web driver is not of
type *Driver, failed to create a random wait time, or another error occurs.

---

### Driver.GetContext()

```go
GetContext() context.Context
```

GetContext retrieves the context associated with the Driver instance.

**Returns:**

context.Context: The context associated with this Driver.

---

### Driver.SetContext(context.Context)

```go
SetContext(context.Context)
```

SetContext associates a new context with the Driver instance.

**Parameters:**

ctx (context.Context): The new context to be associated with this Driver.

---

### GetPageSource(web.Site)

```go
GetPageSource(web.Site) string, error
```

GetPageSource retrieves the HTML source code of the currently loaded
page in the provided Site's session.

**Parameters:**

site (web.Site): The site whose source code is to be retrieved.

**Returns:**

string: The source code of the currently loaded page.
error: An error if any occurred during source code retrieval.

---

### Init(bool, bool)

```go
Init(bool, bool) web.Browser, error
```

Init initializes a chrome browser instance with the specified headless mode and
SSL certificate error ignoring options, then returns the browser instance for
further operations.

**Parameters:**

headless (bool): Whether or not the browser should be in headless mode.
ignoreCertErrors (bool): Whether or not SSL certificate errors should be ignored.

**Returns:**

web.Browser: An initialized Browser instance.
error: Any error encountered during initialization.

---

### Navigate(web.Site, []InputAction, time.Duration)

```go
Navigate(web.Site, []InputAction, time.Duration) error
```

Navigate performs the provided actions sequentially on the provided Site's
session. It enables network events and sets up request logging.

**Parameters:**

site (web.Site): The site on which the actions should be performed.
actions ([]InputAction): A slice of InputAction objects which define
the actions to be performed.
waitTime (time.Duration): The time to wait between actions.

**Returns:**

error: An error if any occurred during navigation.

---

### SaveCookiesToDisk(web.Site, string)

```go
SaveCookiesToDisk(web.Site, string) error
```

SaveCookiesToDisk retrieves cookies from the current session and writes them to a file.

**Parameters:**

site (web.Site): The site from which to retrieve cookies.
filePath (string): The file path where the cookies should be saved.

**Returns:**

error: An error if any occurred during cookie retrieval or file writing.

---

### ScreenShot(web.Site, string)

```go
ScreenShot(web.Site, string) error
```

ScreenShot captures a screenshot of the currently loaded page in the
provided Site's session and writes the image data to the provided file path.

**Parameters:**

site (web.Site): The site whose page a screenshot should be taken of.
imgPath (string): The path to which the screenshot should be saved.

**Returns:**

error: An error if any occurred during screenshot capturing or saving.

---

## Installation

To use the goutils/v2/cdpu package, you first need to install it.
Follow the steps below to install via go get.

```bash
go get github.com/l50/goutils/v2/cdpu
```

---

## Usage

After installation, you can import the package in your Go project
using the following import statement:

```go
import "github.com/l50/goutils/v2/cdpu"
```

---

## Tests

To ensure the package is working correctly, run the following
command to execute the tests for `goutils/v2/cdpu`:

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
