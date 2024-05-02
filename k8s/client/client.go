package k8s

import (
	"fmt"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// FileReaderFunc defines a function signature for reading a file from a given
// path.
type FileReaderFunc func(string) ([]byte, error)

// KubernetesClient wraps a clientset to interact with Kubernetes APIs.
//
// **Attributes:**
//
// Clientset: The clientset interface provided by client-go to interact with
// Kubernetes resources.
// DynamicClient: The dynamic client interface provided by client-go to interact
// with Kubernetes resources.
// Config: The kubeconfig configuration used to create the clientset and dynamic
// client.
type KubernetesClient struct {
	Clientset     kubernetes.Interface
	DynamicClient dynamic.Interface
	Config        *rest.Config
}

// KubernetesClientInterface defines the interface for the KubernetesClient.
//
// **Methods:**
//
// NewForConfig: Creates a new clientset using the provided REST configuration.
// NewDynamicForConfig: Creates a new dynamic client using the provided REST
// configuration.
// RESTConfigFromKubeConfig: Creates a REST configuration from the provided
// kubeconfig data.
type KubernetesClientInterface interface {
	NewForConfig(config *rest.Config) (kubernetes.Interface, error)
	NewDynamicForConfig(config *rest.Config) (dynamic.Interface, error)
	RESTConfigFromKubeConfig(configData []byte) (*rest.Config, error)
}

// RealKubernetesClient implements the KubernetesClientInterface using the
// client-go library.
type RealKubernetesClient struct{}

// NewForConfig creates a new clientset using the provided REST configuration.
//
// **Parameters:**
//
// config: The REST configuration to use to create the clientset.
//
// **Returns:**
//
// *kubernetes.Clientset: A new clientset instance created using the provided
// REST configuration.
// error: An error if any issue occurs while creating the clientset.
func (r *RealKubernetesClient) NewForConfig(config *rest.Config) (*kubernetes.Clientset, error) {
	return kubernetes.NewForConfig(config)
}

// NewDynamicForConfig creates a new dynamic client using the provided REST
// configuration.
//
// **Parameters:**
//
// config: The REST configuration to use to create the dynamic client.
//
// **Returns:**
//
// dynamic.Interface: A new dynamic client instance created using the provided
// REST configuration.
// error: An error if any issue occurs while creating the dynamic client.
func (r *RealKubernetesClient) NewDynamicForConfig(config *rest.Config) (dynamic.Interface, error) {
	return dynamic.NewForConfig(config)
}

// RESTConfigFromKubeConfig creates a REST configuration from the provided
// kubeconfig data.
//
// **Parameters:**
//
// configData: The kubeconfig data to use to create the REST configuration.
//
// **Returns:**
//
// *rest.Config: A new REST configuration instance created using the provided
// kubeconfig data.
// error: An error if any issue occurs while creating the REST configuration.
func (r *RealKubernetesClient) RESTConfigFromKubeConfig(configData []byte) (*rest.Config, error) {
	return clientcmd.RESTConfigFromKubeConfig(configData)
}

// NewKubernetesClient creates a new KubernetesClient using the provided
// kubeconfig path and file reader function.
//
// **Parameters:**
//
// kubeconfig: Path to the kubeconfig file to configure access to the Kubernetes
// API.
// reader: A function to read the kubeconfig file from the specified path.
//
// **Returns:**
//
// *KubernetesClient: A new KubernetesClient instance configured with the
// specified kubeconfig.
// error: An error if any issue occurs while creating the Kubernetes client.
func NewKubernetesClient(kubeconfig string, reader FileReaderFunc, client KubernetesClientInterface) (*KubernetesClient, error) {
	configData, err := reader(kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("error reading kubeconfig: %v", err)
	}

	config, err := client.RESTConfigFromKubeConfig(configData)
	if err != nil {
		return nil, fmt.Errorf("error building kubeconfig: %v", err)
	}

	clientset, err := client.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error creating Kubernetes client: %v", err)
	}

	dynamicClient, err := client.NewDynamicForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error creating dynamic Kubernetes client: %v", err)
	}

	return &KubernetesClient{Clientset: clientset, DynamicClient: dynamicClient, Config: config}, nil
}
