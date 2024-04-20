package k8s

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type fileReaderFunc func(string) ([]byte, error)

// KubernetesClient wraps a clientset to interact with Kubernetes APIs.
type KubernetesClient struct {
	Clientset *kubernetes.Clientset
}

// NewKubernetesClient creates a new KubernetesClient using the provided kubeconfig path.
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

	return &KubernetesClient{Clientset: clientset}, nil
}
