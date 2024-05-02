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

### DefaultExecutorCreator.NewSPDYExecutor(*rest.Config, string, *url.URL)

```go
NewSPDYExecutor(*rest.Config, string, *url.URL) remotecommand.Executor, error
```

NewSPDYExecutor creates a new SPDY executor given a configuration, method,
and URL. It returns a remotecommand.Executor and an error.

**Parameters:**

config: A pointer to a rest.Config struct that includes the configuration
for the executor.
method: A string representing the HTTP method to use for the request.
url: A pointer to a url.URL struct that includes the URL for the request.

**Returns:**

remotecommand.Executor: The created SPDY executor.
error: An error if any issue occurs while creating the executor.

---

### DescribeKubernetesResource(context.Context, *client.KubernetesClient, string, schema.GroupVersionResource)

```go
DescribeKubernetesResource(context.Context *client.KubernetesClient string schema.GroupVersionResource) string error
```

DescribeKubernetesResource retrieves the details of a specific Kubernetes
resource using the provided dynamic client, resource name, namespace, and
GroupVersionResource (GVR).

**Parameters:**

ctx: The context to use for the request.
kc: The KubernetesClient that includes both the standard and dynamic clients.
resourceName: The name of the resource to describe.
namespace: The namespace of the resource.
gvr: The GroupVersionResource of the resource.

**Returns:**

string: A string representation of the resource, similar to `kubectl describe`.
error: An error if any issue occurs while trying to describe the resource.

---

### ExecKubernetesResources(ExecParams)

```go
ExecKubernetesResources(ExecParams) string, error
```

ExecKubernetesResources executes a command in a specified resource within a
given namespace using the existing KubernetesClient.

**Parameters:**

ctx: The context to use for the request.
kc: The KubernetesClient that includes both the standard and dynamic clients.
namespace: The namespace of the resource where the pod is located.
podName: The name of the pod to execute the command in.
command: A slice of strings representing the command to execute inside the pod.

**Returns:**

string: The output from the executed command or an error message if execution fails.
error: An error if any issue occurs during the setup or execution of the command.

---

### GetResourceStatus(context.Context, *client.KubernetesClient, string, schema.GroupVersionResource)

```go
GetResourceStatus(context.Context *client.KubernetesClient string schema.GroupVersionResource) bool error
```

GetResourceStatus checks the status of any Kubernetes resource.

**Parameters:**

ctx: A context.Context to control the operation.
kc: The KubernetesClient that includes both the standard and dynamic clients.
resourceName: The name of the resource being checked.
namespace: The namespace of the resource.
gvr: The schema.GroupVersionResource that specifies the resource type.

**Returns:**

bool: true if the resource status is 'Running', false otherwise.
error: An error if the resource cannot be retrieved or the status is not found.

---

### WaitForResourceState(context.Context, string, func(name, namespace string) (bool, error))

```go
WaitForResourceState(context.Context string func(name namespace string) (bool error)) error
```

WaitForResourceState waits for a Kubernetes resource to reach a specified state.

**Parameters:**

ctx: A context.Context to allow for cancellation and timeouts.
resourceName: The name of the resource to monitor.
namespace: The namespace in which the resource exists.
resourceType: The type of the resource (e.g., Pod, Service).
desiredState: A string representing the desired state (e.g., "Running", "Deleted").
checkStatusFunc: A function that checks if the resource is in the desired state.

**Returns:**

error: An error if the waiting is cancelled by context, times out, or
fails to determine the state.

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
