package k8s_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	client "github.com/l50/goutils/v2/k8s/client"
	jobs "github.com/l50/goutils/v2/k8s/jobs"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	fakedynamic "k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/homedir"
)

// MockManifestConfig is a mock implementation of the ManifestConfig
type MockManifestConfig struct {
	mock.Mock
	ReadFile func(string) ([]byte, error)
}

func (m *MockManifestConfig) ApplyOrDeleteManifest(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockKubeConfigBuilder is a mock implementation of the KubeConfigBuilder interface.
type MockKubeConfigBuilder struct {
	mock.Mock
}

func (m *MockKubeConfigBuilder) BuildConfigFromFlags(masterUrl, kubeconfigPath string) (*rest.Config, error) {
	args := m.Called(masterUrl, kubeconfigPath)
	config, _ := args.Get(0).(*rest.Config)
	return config, args.Error(1)
}

func (m *MockKubeConfigBuilder) NewDynamicForConfig(config *rest.Config) (dynamic.Interface, error) {
	args := m.Called(config)
	client, _ := args.Get(0).(dynamic.Interface)
	return client, args.Error(1)
}

func (m *MockKubeConfigBuilder) NewForConfig(config *rest.Config) (kubernetes.Interface, error) {
	args := m.Called(config)
	client, _ := args.Get(0).(kubernetes.Interface)
	return client, args.Error(1)
}

func (m *MockKubeConfigBuilder) RESTConfigFromKubeConfig(configData []byte) (*rest.Config, error) {
	args := m.Called(configData)
	config, _ := args.Get(0).(*rest.Config)
	return config, args.Error(1)
}

// MockDynamicClientCreator is a mock implementation of the DynamicClientCreator interface.
type MockDynamicClientCreator struct {
	mock.Mock
}

func (m *MockDynamicClientCreator) NewDynamicForConfig(config *rest.Config) (dynamic.Interface, error) {
	args := m.Called(config)
	client, _ := args.Get(0).(dynamic.Interface)
	return client, args.Error(1)
}

func (m *MockDynamicClientCreator) NewForConfig(config *rest.Config) (kubernetes.Interface, error) {
	args := m.Called(config)
	client, _ := args.Get(0).(kubernetes.Interface)
	return client, args.Error(1)
}

func (m *MockDynamicClientCreator) RESTConfigFromKubeConfig(configData []byte) (*rest.Config, error) {
	args := m.Called(configData)
	config, _ := args.Get(0).(*rest.Config)
	return config, args.Error(1)
}

func TestApplyKubernetesJob(t *testing.T) {
	homeDir := homedir.HomeDir()
	kubeconfig := fmt.Sprintf("%s/.kube/config", homeDir)
	os.Setenv("KUBECONFIG", kubeconfig)

	tests := []struct {
		name        string
		setupMocks  func(mockManifestConfig *MockManifestConfig)
		jobFilePath string
		namespace   string
		expectError bool
	}{
		{
			name: "successful job application",
			setupMocks: func(mockManifestConfig *MockManifestConfig) {
				mockManifestConfig.On("ApplyOrDeleteManifest", mock.Anything).Return(nil).Once()
				mockManifestConfig.ReadFile = func(string) ([]byte, error) {
					return []byte(`apiVersion: batch/v1
kind: Job
metadata:
  name: unique-job-1
spec:
  template:
    metadata:
      labels:
        app: my-job
    spec:
      containers:
      - name: my-container
        image: my-image
      restartPolicy: Never`), nil
				}
			},
			jobFilePath: "testdata/job.yaml",
			namespace:   "default",
			expectError: false,
		},
		{
			name: "failure in applying manifest",
			setupMocks: func(mockManifestConfig *MockManifestConfig) {
				mockManifestConfig.On("ApplyOrDeleteManifest", mock.Anything).Return(fmt.Errorf("failed to apply manifest")).Once()
				mockManifestConfig.ReadFile = func(string) ([]byte, error) {
					return nil, fmt.Errorf("failed to read file")
				}
			},
			jobFilePath: "testdata/job.yaml",
			namespace:   "default",
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockManifestConfig := new(MockManifestConfig)
			tc.setupMocks(mockManifestConfig)

			kubeClient := &client.KubernetesClient{
				Config: &rest.Config{
					Host: kubeconfig,
				},
				Clientset:     fake.NewSimpleClientset(),
				DynamicClient: fakedynamic.NewSimpleDynamicClient(runtime.NewScheme()),
			}

			jobsClient := &jobs.JobsClient{
				Client: kubeClient,
			}

			require.NotNil(t, mockManifestConfig)

			// Call ReadFile before applying the job to simulate reading the job file
			_, _ = mockManifestConfig.ReadFile(tc.jobFilePath)

			err := jobsClient.ApplyKubernetesJob(tc.jobFilePath, tc.namespace, mockManifestConfig.ReadFile)
			require.Equal(t, tc.expectError, err != nil, "expected error: %v, got: %v", tc.expectError, err)

			if err != nil {
				t.Logf("error: %v", err)
			}
		})
	}
}
