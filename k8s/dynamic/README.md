# goutils/v2/k8s

The `k8s` package is a collection of utility functions
designed to simplify common k8s tasks.

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

### GetResourceStatus(context.Context, dynamic.Interface, string, schema.GroupVersionResource)

```go
GetResourceStatus(context.Context dynamic.Interface string schema.GroupVersionResource) bool error
```

GetResourceStatus checks the status of any Kubernetes resource.

**Parameters:**

ctx: A context.Context to control the operation.
client: The dynamic.Interface client used for Kubernetes API calls.
resourceName: The name of the resource being checked.
namespace: The namespace of the resource.
gvr: The schema.GroupVersionResource that specifies the resource type.

**Returns:**

bool: true if the resource status is 'Running', false otherwise.
error: An error if the resource cannot be retrieved or the status is not found.

---

### WaitForResourceReady(context.Context, string, func(name, namespace string) (bool, error))

```go
WaitForResourceReady(context.Context string func(name namespace string) (bool error)) error
```

WaitForResourceReady waits for any Kubernetes resource to reach a ready state.

**Parameters:**

ctx: A context.Context to allow for cancellation and timeouts.
resourceName: The name of the resource to monitor.
namespace: The namespace in which the resource exists.
resourceType: The type of the resource (e.g., Pod, Service).
checkStatusFunc: A function that checks if the resource is ready.

**Returns:**

error: An error if the waiting is cancelled by context, times out, or
fails to determine readiness.

---

## Installation

To use the goutils/v2/k8s package, you first need to install it.
Follow the steps below to install via go get.

```bash
go get github.com/l50/goutils/v2/k8s
```

---

## Usage

After installation, you can import the package in your Go project
using the following import statement:

```go
import "github.com/l50/goutils/v2/k8s"
```

---

## Tests

To ensure the package is working correctly, run the following
command to execute the tests for `goutils/v2/k8s`:

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
