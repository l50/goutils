# goutils/v2/web

The `web` package is a part of `goutils` library.

It provides utility functions to interact with web applications
in Go, including the means to drive a headless browser
using [chromeDP](https://github.com/chromedp/chromedp).

---

## Functions

### CancelAll

```go
func CancelAll(cancels ...func())
```

Executes all provided cancel functions. Typically used for cleaning
up or aborting operations that were started earlier and can be cancelled.

### GetRandomWait

```go
func GetRandomWait(minWait, maxWait time.Duration) (time.Duration, error)
```

Return a random duration between the specified minWait and maxWait durations.

### Wait

```go
func Wait(near float64) (time.Duration, error)
```

Generates a random period of time anchored to a given input value.
Useful for simulating more human-like interaction in the context of a web application.

---

## Installation

To use the `goutils/v2/web` package, you need to install it via `go get`:

```bash
go get github.com/l50/goutils/v2/web
```

---

## Usage

After installation, you can import it in your project:

```go
import "github.com/l50/goutils/v2/web"
```

---

## Tests

To run the tests for the `goutils/v2/web` package, navigate to
your `$GOPATH/src/github.com/l50/goutils/v2/web` directory
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
