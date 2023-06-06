# goutils/web/chrome

The `web/chrome` package is a part of `goutils` library.

It provides utility functions to interact with web applications 
in Go, including the means to drive a headless browser 
using [chromeDP](https://github.com/chromedp/chromedp).

---

## Functions

### Init

```go
func Init(headless bool, ignoreCertErrors bool) (web.Browser, error)
```

This function initializes a Google Chrome browser instance with the 
specified headless mode and SSL certificate error ignoring options. 
It creates contexts and associated cancel functions for browser operation.

### GetPageSource

```go
func GetPageSource(site web.Site) (string, error)
```

Retrieves the HTML source code of the currently loaded page in a site session.

### Navigate

```go
func Navigate(site web.Site, actions []InputAction, waitTime time.Duration) error
```

Navigates a site using provided actions. It enables network events 
and sets up request logging.

### ScreenShot

```go
func ScreenShot(site web.Site, imgPath string) error
```

ScreenShot takes a screenshot of the input `targetURL` and saves it to `imgPath`.

---

## Installation

To use the `goutils/chrome` package, you need to install it via `go get`:

```bash
go get github.com/l50/goutils/chrome
```

---

## Usage

After installation, you can import it in your project:

```go
import "github.com/l50/goutils/web"
```

---

## Tests

To run the tests for the `goutils/web` package, navigate to
your `$GOPATH/src/github.com/l50/goutils/web` directory
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
