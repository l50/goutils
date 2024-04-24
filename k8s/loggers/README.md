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

### DeploymentLogger.FetchAndLog(context.Context)

```go
FetchAndLog(context.Context) error
```

FetchAndLog fetches the deployment details and logs related pod events.

**Parameters:**

ctx: Context to control the request lifetime.

**Returns:**

error: An error if any occurs during fetching and logging.

---

### FetchAndLogPods(context.Context, kubernetes.Interface, string)

```go
FetchAndLogPods(context.Context, kubernetes.Interface, string) error
```

FetchAndLogPods fetches and logs pods based on the specified label selector.

**Parameters:**

ctx: Context to control the request lifetime.
clientset: Kubernetes clientset to interact with Kubernetes API.
namespace: Namespace from which to list the pods.
labelSelector: String defining the label selector for filtering pods.

**Returns:**

error: An error if any occurs during fetching and logging of pods.

---

### NewDeploymentLogger(*k8s.KubernetesClient, string)

```go
NewDeploymentLogger(*k8s.KubernetesClient, string) *DeploymentLogger
```

NewDeploymentLogger creates a new instance of DeploymentLogger.

**Parameters:**

kc: Pointer to KubernetesClient.
namespace: Namespace where the deployment is located.
deploymentName: Name of the deployment.

**Returns:**

*DeploymentLogger: A new instance of DeploymentLogger.

---

### NewServiceLogger(*k8s.KubernetesClient, string)

```go
NewServiceLogger(*k8s.KubernetesClient, string) *ServiceLogger
```

NewServiceLogger creates a new instance of ServiceLogger.

**Parameters:**

kc: Pointer to KubernetesClient.
namespace: Namespace where the service is located.
serviceName: Name of the service.

**Returns:**

*ServiceLogger: A new instance of ServiceLogger.

---

### ServiceLogger.FetchAndLog(context.Context)

```go
FetchAndLog(context.Context) error
```

FetchAndLog fetches the service details and logs related pod events.

**Parameters:**

ctx: Context to control the request lifetime.

**Returns:**

error: An error if any occurs during fetching and logging.

---

### StreamLogs(kubernetes.Interface, string)

```go
StreamLogs(kubernetes.Interface, string) error
```

StreamLogs streams logs for a specific resource within a namespace.

**Parameters:**

clientset: Kubernetes clientset to interact with Kubernetes API.
namespace: Namespace where the resource is located.
resourceType: Type of resource ('pod', 'job', or 'deployment').
resourceName: Name of the resource to stream logs from.

**Returns:**

error: An error if any occurs during the log streaming process.

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
