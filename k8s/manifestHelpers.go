package k8s

import (
	"context"
	"fmt"
	"os"
	"strings"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

// MetadataConfig holds metadata configuration for Kubernetes resources.
//
// **Attributes:**
//
// Name: The name of the resource.
type MetadataConfig struct {
	Name string // and any other fields you expect
}

// ManifestConfig represents the configuration needed to manage Kubernetes manifests.
//
// **Attributes:**
//
// KubeConfigPath: Path to the kubeconfig file.
// ManifestPath: Path to the Kubernetes manifest file.
// Namespace: Kubernetes namespace in which the operations will be performed.
// Type: The type of manifest (raw, Helm, or Kustomize).
// Operation: The operation to perform (apply or delete).
// Metadata: Metadata related to the manifest.
// Client: The dynamic Kubernetes client interface.
// ReadFile: Function to read the manifest file from the filesystem.
type ManifestConfig struct {
	KubeConfigPath string
	ManifestPath   string
	Namespace      string
	Type           ManifestType
	Operation      ManifestOperation
	Metadata       *MetadataConfig
	Client         dynamic.Interface
	ReadFile       func(string) ([]byte, error)
}

// ManifestType defines the type of Kubernetes manifest.
//
// **Values:**
//
// ManifestRaw: Raw Kubernetes manifest.
// ManifestHelm: Helm chart.
// ManifestKustomize: Kustomize configuration.
type ManifestType int

const (
	ManifestRaw ManifestType = iota
	ManifestHelm
	ManifestKustomize
)

// ManifestOperation specifies the type of operation to perform on the manifest.
//
// **Values:**
//
// OperationApply: Apply the manifest.
// OperationDelete: Delete the manifest.
type ManifestOperation int

const (
	OperationApply ManifestOperation = iota
	OperationDelete
)

// NewManifestConfig creates a new ManifestConfig with default settings.
//
// **Returns:**
//
// *ManifestConfig: A new ManifestConfig instance with ReadFile set to os.ReadFile.
func NewManifestConfig() *ManifestConfig {
	return &ManifestConfig{
		ReadFile: os.ReadFile,
	}
}

// String returns the string representation of the ManifestType.
//
// **Returns:**
//
// string: The string representation of the ManifestType.
func (mo ManifestOperation) String() string {
	switch mo {
	case OperationApply:
		return "apply"
	case OperationDelete:
		return "delete"
	default:
		return "unknown"
	}
}

// ApplyOrDeleteManifest applies or deletes a Kubernetes manifest based on the
// ManifestConfig settings.
//
// **Parameters:**
//
// ctx: Context for the operation.
//
// **Returns:**
//
// error: Error if any issue occurs while applying or deleting the manifest.
func (mc *ManifestConfig) ApplyOrDeleteManifest(ctx context.Context) error {
	if mc.Client == nil {
		config, err := clientcmd.BuildConfigFromFlags("", mc.KubeConfigPath)
		if err != nil {
			return fmt.Errorf("error building kubeconfig: %v", err)
		}
		mc.Client, err = dynamic.NewForConfig(config)
		if err != nil {
			return fmt.Errorf("error creating dynamic client: %v", err)
		}
	}
	switch mc.Type {
	case ManifestRaw:
		return mc.HandleRawManifest(ctx, mc.Client)
	case ManifestHelm:
		return mc.handleHelmManifest()
	default:
		return fmt.Errorf("unsupported manifest type")
	}
}

// HandleRawManifest applies or deletes raw Kubernetes manifests based on the
// operation specified in ManifestConfig.
//
// **Parameters:**
//
// ctx: The context for the operation.
// dynClient: The dynamic client to perform Kubernetes operations.
//
// **Returns:**
//
// error: Error if any issue occurs while handling the raw manifest.
func (mc *ManifestConfig) HandleRawManifest(ctx context.Context, dynClient dynamic.Interface) error {
	data, err := mc.ReadFile(mc.ManifestPath)
	if err != nil {
		return fmt.Errorf("error reading manifest file: %v", err)
	}
	decoder := yaml.NewYAMLOrJSONDecoder(strings.NewReader(string(data)), 2048)
	for {
		rawObj := &unstructured.Unstructured{}
		if err := decoder.Decode(rawObj); err != nil {
			if err.Error() == "EOF" {
				break
			}
			return fmt.Errorf("error decoding YAML: %v", err)
		}

		if rawObj.Object == nil {
			continue
		}

		gvk := rawObj.GroupVersionKind()
		gvr, err := mc.groupVersionResource(gvk)
		if err != nil {
			return fmt.Errorf("error getting GroupVersionResource for %v: %v", gvk, err)
		}
		resourceClient := dynClient.Resource(gvr).Namespace(mc.Namespace)

		var operationErr error
		switch mc.Operation {
		case OperationApply:
			_, operationErr = resourceClient.Create(ctx, rawObj, metav1.CreateOptions{})
		case OperationDelete:
			operationErr = resourceClient.Delete(ctx, rawObj.GetName(), metav1.DeleteOptions{})
		}

		if operationErr != nil {
			return fmt.Errorf("failed to %s manifest: %v", strings.ToLower(mc.Operation.String()), operationErr)
		}
	}
	return nil
}

// groupVersionResource constructs a GroupVersionResource from a GroupVersionKind.
//
// **Parameters:**
//
// gvk: The GroupVersionKind to convert.
//
// **Returns:**
//
// GroupVersionResource: The constructed GroupVersionResource.
// error: Error if the kind is empty.
func (mc *ManifestConfig) groupVersionResource(gvk schema.GroupVersionKind) (schema.GroupVersionResource, error) {
	if gvk.Kind == "" {
		return schema.GroupVersionResource{}, fmt.Errorf("kind must not be empty")
	}

	resource := strings.ToLower(gvk.Kind) + "s"
	return schema.GroupVersionResource{
		Group:    gvk.Group,
		Version:  gvk.Version,
		Resource: resource,
	}, nil
}

// handleHelmManifest manages Helm chart installations or deletions based on
// the operation specified in ManifestConfig.
//
// **Parameters:**
//
// **Returns:**
//
// error: Error if any issue occurs while handling the Helm manifest.
func (mc *ManifestConfig) handleHelmManifest() error {
	settings := cli.New() // Initialize Helm settings
	actionConfig := new(action.Configuration)

	if err := actionConfig.Init(settings.RESTClientGetter(), mc.Namespace, os.Getenv("HELM_DRIVER"), nil); err != nil {
		return fmt.Errorf("failed to initialize Helm: %v", err)
	}

	switch mc.Operation {
	case OperationApply:
		return mc.installHelmChart(actionConfig)
	case OperationDelete:
		return mc.deleteHelmRelease(actionConfig)
	default:
		return fmt.Errorf("unsupported Helm operation")
	}
}

// installHelmChart installs a Helm chart using the specified action configuration.
//
// **Parameters:**
//
// actionConfig: Configuration for the Helm install action.
//
// **Returns:**
//
// error: Error if the installation fails.
func (mc *ManifestConfig) installHelmChart(actionConfig *action.Configuration) error {
	// Create an instance of the Install action
	install := action.NewInstall(actionConfig)
	install.ReleaseName = mc.Metadata.Name // The release name must be set if not automatically generated
	install.Namespace = mc.Namespace       // Namespace where the chart will be installed

	// Load the chart from the given path
	chart, err := loader.Load(mc.ManifestPath)
	if err != nil {
		return fmt.Errorf("failed to load helm chart: %v", err)
	}

	// Run the installation
	_, err = install.Run(chart, nil) // Pass nil if no custom values are needed
	if err != nil {
		return fmt.Errorf("helm install failed: %v", err)
	}

	return nil
}

// deleteHelmRelease uninstalls a Helm release using the specified action configuration.
//
// **Parameters:**
//
// actionConfig: Configuration for the Helm uninstall action.
//
// **Returns:**
//
// error: Error if the uninstallation fails.
func (mc *ManifestConfig) deleteHelmRelease(actionConfig *action.Configuration) error {
	client := action.NewUninstall(actionConfig)

	// Ensure ReleaseName is derived safely.
	if mc.Metadata == nil || mc.Metadata.Name == "" {
		return fmt.Errorf("invalid release name for deletion")
	}
	_, err := client.Run(mc.Metadata.Name)
	if err != nil {
		return fmt.Errorf("helm uninstall failed: %v", err)
	}

	return nil
}
