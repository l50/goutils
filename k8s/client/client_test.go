package k8s_test

import (
	"fmt"
	"testing"

	client "github.com/l50/goutils/v2/k8s/client"
	"github.com/stretchr/testify/mock"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type MockKubernetesClient struct {
	mock.Mock
}

func NewKubernetesClient(kubeconfig string, reader client.FileReaderFunc, kClient client.KubernetesClientInterface) (*client.KubernetesClient, error) {
	configData, err := reader(kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("error reading kubeconfig: %v", err)
	}

	config, err := kClient.RESTConfigFromKubeConfig(configData)
	if err != nil {
		return nil, fmt.Errorf("error building kubeconfig: %v", err)
	}

	kInterface, err := kClient.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error creating Kubernetes client: %v", err)
	}

	// Safely assert the type to *kubernetes.Clientset
	kClientset, ok := kInterface.(*kubernetes.Clientset)
	if !ok {
		return nil, fmt.Errorf("failed to assert Kubernetes interface to Clientset")
	}

	dynamicClient, err := kClient.NewDynamicForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error creating dynamic Kubernetes client: %v", err)
	}

	return &client.KubernetesClient{Clientset: kClientset, DynamicClient: dynamicClient, Config: config}, nil
}

func (m *MockKubernetesClient) NewForConfig(config *rest.Config) (kubernetes.Interface, error) {
	args := m.Called(config)
	return args.Get(0).(kubernetes.Interface), args.Error(1)
}

func (m *MockKubernetesClient) NewDynamicForConfig(config *rest.Config) (dynamic.Interface, error) {
	args := m.Called(config)
	return args.Get(0).(dynamic.Interface), args.Error(1)
}

func (m *MockKubernetesClient) RESTConfigFromKubeConfig(configData []byte) (*rest.Config, error) {
	args := m.Called(configData)
	return args.Get(0).(*rest.Config), args.Error(1)
}

func TestNewKubernetesClient(t *testing.T) {
	tests := []struct {
		name        string
		kubeconfig  string
		data        []byte
		expectError bool
	}{
		{
			name:       "valid kubeconfig",
			kubeconfig: "path/to/valid/kubeconfig",
			data: []byte(`apiVersion: v1
clusters:
- cluster:
    server: https://localhost:6443
  name: test-cluster
contexts:
- context:
    cluster: test-cluster
    user: test-user
  name: test-context
current-context: test-context
kind: Config
preferences: {}
users:
- name: test-user
  user:
    token: fake-token`),
			expectError: false,
		},
		{
			name:        "invalid kubeconfig",
			kubeconfig:  "invalid/path",
			data:        nil,
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			reader := func(path string) ([]byte, error) {
				if path == tc.kubeconfig && tc.data != nil {
					return tc.data, nil
				}
				return nil, fmt.Errorf("error reading kubeconfig")
			}

			mockClient := new(MockKubernetesClient)
			mockClient.On("RESTConfigFromKubeConfig", tc.data).Return(&rest.Config{}, nil)
			mockClient.On("NewForConfig", mock.Anything).Return(&kubernetes.Clientset{}, nil)
			mockClient.On("NewDynamicForConfig", mock.Anything).Return(dynamic.NewForConfigOrDie(&rest.Config{}), nil)
			client, err := client.NewKubernetesClient(tc.kubeconfig, reader, mockClient)
			if (err != nil) != tc.expectError {
				t.Errorf("Test '%s' failed - expected error: %v, got: %v", tc.name, tc.expectError, err)
			}
			if err == nil {
				// Perform further validation on the successful creation case
				if client == nil || client.Clientset == nil || client.DynamicClient == nil {
					t.Errorf("Test '%s' failed - clientset or dynamic client not properly initialized", tc.name)
				}
			}
		})
	}
}
