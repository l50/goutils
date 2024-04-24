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

### JobsClient.DeleteKubernetesJob(context.Context, string)

```go
DeleteKubernetesJob(context.Context, string) error
```

DeleteKubernetesJob deletes a specified Kubernetes job within
a given namespace. It sets the deletion propagation policy
to 'Foreground' to ensure that the delete operation waits
until the cascading delete has completed.

**Parameters:**

ctx: Context for managing control flow of the request.
jobName: Name of the Kubernetes job to delete.
namespace: Namespace where the job is located.

**Returns:**

error: An error if the job could not be deleted.

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
