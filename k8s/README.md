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

### KubernetesClient.DeleteKubernetesJob(string)

```go
DeleteKubernetesJob(string) error
```

DeleteKubernetesJob deletes a Kubernetes Job in the specified namespace.

---

### ManifestConfig.ApplyOrDeleteManifest(context.Context)

```go
ApplyOrDeleteManifest(context.Context) error
```

ApplyOrDeleteManifest applies or deletes a Kubernetes manifest based on the
ManifestConfig settings.

**Parameters:**

ctx: Context for the operation.

**Returns:**

error: Error if any issue occurs while applying or deleting the manifest.

---

### ManifestConfig.HandleRawManifest(context.Context, dynamic.Interface)

```go
HandleRawManifest(context.Context, dynamic.Interface) error
```

HandleRawManifest applies or deletes raw Kubernetes manifests based on the
operation specified in ManifestConfig.

**Parameters:**

ctx: The context for the operation.
dynClient: The dynamic client to perform Kubernetes operations.

**Returns:**

error: Error if any issue occurs while handling the raw manifest.

---

### ManifestOperation.String()

```go
String() string
```

String returns the string representation of the ManifestType.

**Returns:**

string: The string representation of the ManifestType.

---

### NewKubernetesClient(string, fileReaderFunc)

```go
NewKubernetesClient(string, fileReaderFunc) *KubernetesClient, error
```

NewKubernetesClient creates a new KubernetesClient using the provided
kubeconfig path and file reader function.

**Parameters:**

kubeconfig: Path to the kubeconfig file to configure access to the Kubernetes
API.
reader: A function to read the kubeconfig file from the specified path.

**Returns:**

*KubernetesClient: A new KubernetesClient instance configured with the
specified kubeconfig.
error: An error if any issue occurs while creating the Kubernetes client.

---

### NewManifestConfig()

```go
NewManifestConfig() *ManifestConfig
```

NewManifestConfig creates a new ManifestConfig with default settings.

**Returns:**

*ManifestConfig: A new ManifestConfig instance with ReadFile set to os.ReadFile.

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
