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

### CheckKubeConfig()

```go
CheckKubeConfig() error
```

CheckKubeConfig checks if the KUBECONFIG environment variable is set and
points to a valid kubeconfig file.

Returns:

error: An error if the KUBECONFIG environment variable is not set or does
not point to a valid kubeconfig file.

---

### NewKubernetesClient(string, FileReaderFunc, KubernetesClientInterface)

```go
NewKubernetesClient(string FileReaderFunc KubernetesClientInterface) *KubernetesClient error
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

### RealKubernetesClient.NewDynamicForConfig(*rest.Config)

```go
NewDynamicForConfig(*rest.Config) dynamic.Interface, error
```

NewDynamicForConfig creates a new dynamic client using the provided REST
configuration.

**Parameters:**

config: The REST configuration to use to create the dynamic client.

**Returns:**

dynamic.Interface: A new dynamic client instance created using the provided
REST configuration.
error: An error if any issue occurs while creating the dynamic client.

---

### RealKubernetesClient.NewForConfig(*rest.Config)

```go
NewForConfig(*rest.Config) kubernetes.Interface, error
```

NewForConfig creates a new clientset using the provided REST configuration.

**Parameters:**

config: The REST configuration to use to create the clientset.

**Returns:**

*kubernetes.Clientset: A new clientset instance created using the provided
REST configuration.
error: An error if any issue occurs while creating the clientset.

---

### RealKubernetesClient.RESTConfigFromKubeConfig([]byte)

```go
RESTConfigFromKubeConfig([]byte) *rest.Config, error
```

RESTConfigFromKubeConfig creates a REST configuration from the provided
kubeconfig data.

**Parameters:**

configData: The kubeconfig data to use to create the REST configuration.

**Returns:**

*rest.Config: A new REST configuration instance created using the provided
kubeconfig data.
error: An error if any issue occurs while creating the REST configuration.

---

### SetupKubeConfig(string)

```go
SetupKubeConfig(string) error
```

SetupKubeConfig sets the KUBECONFIG environment variable to the default path
if it is not already set.

**Parameters:**

defaultPath: The default path to the kubeconfig file.

**Returns:**

error: An error if the kubeconfig file is not found or cannot be accessed

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
