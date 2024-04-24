package k8s

import (
	"fmt"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type fileReaderFunc func(string) ([]byte, error)

// KubernetesClient wraps a clientset to interact with Kubernetes APIs.
//
// **Attributes:**
//
// Clientset: The clientset interface provided by client-go to interact with
// Kubernetes resources.
// DynamicClient: The dynamic client interface provided by client-go to interact
// with Kubernetes resources.
type KubernetesClient struct {
	Clientset     kubernetes.Interface
	DynamicClient dynamic.Interface
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
func NewKubernetesClient(kubeconfig string, reader fileReaderFunc) (*KubernetesClient, error) {
	configData, err := reader(kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("error reading kubeconfig: %v", err)
	}

	config, err := clientcmd.RESTConfigFromKubeConfig(configData)
	if err != nil {
		return nil, fmt.Errorf("error building kubeconfig: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error creating Kubernetes client: %v", err)
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error creating dynamic Kubernetes client: %v", err)
	}

	return &KubernetesClient{Clientset: clientset, DynamicClient: dynamicClient}, nil
}
