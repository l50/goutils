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

### DefaultJobPodNameGetter.GetJobPodName(context.Context, string)

```go
GetJobPodName(context.Context, string) string, error
```

GetJobPodName retrieves the name of the first pod associated with a specific
Kubernetes job within a given namespace. It uses a label selector to find
pods that are labeled with the job's name. This method is typically used in
scenarios where jobs create a single pod or when only the first pod
is of interest.

**Parameters:**

ctx: Context for managing control flow of the request.
jobName: Name of the Kubernetes job to find pods for.
namespace: Namespace where the job and its pods are located.

**Returns:**

string: The name of the first pod found that is associated with the job
error: An error if no pods are found or if an error occurs during the pod retrieval

---

### JobsClient.ApplyKubernetesJob(string, func(string) ([]byte, error))

```go
ApplyKubernetesJob(string, func(string) ([]byte, error)) error
```

ApplyKubernetesJob applies a Kubernetes job manifest to a Kubernetes cluster
using the provided kubeconfig file. The job is applied to the specified namespace.

**Parameters:**

jobFilePath: Path to the job manifest file to apply.
namespace: Namespace where the job should be applied.

**Returns:**

error: An error if the job could not be applied.

---

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

### JobsClient.GetJobPodName(context.Context, string)

```go
GetJobPodName(context.Context, string) string, error
```

GetJobPodName retrieves the name of the first pod associated with a specific
Kubernetes job within a given namespace. It uses a label selector to find
pods that are labeled with the job's name. This method is typically used in
scenarios where jobs create a single pod or when only the first pod
is of interest.

**Parameters:**

ctx: Context for managing control flow of the request.
jobName: Name of the Kubernetes job to find pods for.
namespace: Namespace where the job and its pods are located.

**Returns:**

string: The name of the first pod found that is associated with the job
error: An error if no pods are found or if an error occurs during the pod retrieval

---

### JobsClient.JobExists(context.Context, string)

```go
JobExists(context.Context, string) bool, error
```

JobExists checks if a Kubernetes job with the specified name exists within a given namespace.

**Parameters:**

ctx: Context for managing control flow of the request.
jobName: Name of the Kubernetes job to check for existence.
namespace: Namespace where the job is located.

**Returns:**

bool: true if the job exists, false otherwise.
error: An error if the job existence check fails.

---

### JobsClient.ListKubernetesJobs(context.Context, string)

```go
ListKubernetesJobs(context.Context, string) []batchv1.Job, error
```

ListKubernetesJobs lists Kubernetes jobs from a specified namespace, or all namespaces
if no namespace is specified. This method allows for either targeted or broad job retrieval.

**Parameters:**

ctx: Context for managing control flow of the request.
namespace: Optional; specifies the namespace from which to list jobs. If empty, jobs will be listed from all namespaces.

**Returns:**

[]batchv1.Job: A slice of batchv1.Job objects containing the jobs found.
error: An error if the API call to fetch the jobs fails.

---

### JobsClient.StreamJobLogs(string)

```go
StreamJobLogs(string) error
```

StreamJobLogs monitors a Kubernetes job by waiting for it to reach
the 'Ready' state and then streams logs from the associated pod.

**Parameters:**

jobsClient: A JobsClient for managing Kubernetes jobs.
workloadName: Name of the Kubernetes job to monitor.
namespace: Namespace where the job is located.

**Returns:**

error: An error if the job monitoring fails.

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
