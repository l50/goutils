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

### ManifestConfig.CreateConfigMapFromScript(context.Context, string, string)

```go
CreateConfigMapFromScript(context.Context, string, string) error
```

CreateConfigMapFromScript creates a ConfigMap from a script
file and applies it to the Kubernetes cluster.

**Parameters:**

ctx: The context for the operation.
scriptPath: The path to the script file.
configMapName: The name of the ConfigMap to create.

**Returns:**

error: Error if any issue occurs while creating the ConfigMap.

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
