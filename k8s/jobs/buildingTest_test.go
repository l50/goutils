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
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	fakedynamic "k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
	k8stesting "k8s.io/client-go/testing"
	"k8s.io/client-go/util/homedir"
)

// MockKubernetesClient is a mock implementation of the Kubernetes client
type MockKubernetesClient struct {
	mock.Mock
}

func (m *MockKubernetesClient) RESTConfigFromKubeConfig(kubeconfigPath string) (*rest.Config, error) {
	args := m.Called(kubeconfigPath)
	config, _ := args.Get(0).(*rest.Config)
	return config, args.Error(1)
}

func (m *MockKubernetesClient) NewDynamicForConfig(config *rest.Config) (dynamic.Interface, error) {
	args := m.Called(config)
	client, _ := args.Get(0).(dynamic.Interface)
	return client, args.Error(1)
}

// MockManifestConfig is a mock implementation of the ManifestConfig
type MockManifestConfig struct {
	mock.Mock
	ReadFile func(string) ([]byte, error)
}

func (m *MockManifestConfig) ApplyOrDeleteManifest(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestApplyKubernetesJob(t *testing.T) {
	homeDir := homedir.HomeDir()
	kubeconfig := fmt.Sprintf("%s/.kube/config", homeDir)
	os.Setenv("KUBECONFIG", kubeconfig)

	tests := []struct {
		name        string
		setupMocks  func(m *MockKubernetesClient, mockManifestConfig *MockManifestConfig)
		jobFilePath string
		namespace   string
		expectError bool
	}{
		{
			name: "successful job application",
			setupMocks: func(m *MockKubernetesClient, mockManifestConfig *MockManifestConfig) {
				config := &rest.Config{}
				dynClient := fakedynamic.NewSimpleDynamicClient(runtime.NewScheme())

				// Simulate job existence and deletion
				existingJob := &unstructured.Unstructured{}
				existingJob.SetKind("Job")
				existingJob.SetAPIVersion("batch/v1")
				existingJob.SetName("unique-job-1")
				existingJob.SetNamespace("default")
				dynClient.PrependReactor("create", "jobs", func(action k8stesting.Action) (bool, runtime.Object, error) {
					return true, nil, errors.NewAlreadyExists(schema.GroupResource{Group: "batch", Resource: "jobs"}, "unique-job-1")
				})
				dynClient.PrependReactor("get", "jobs", func(action k8stesting.Action) (bool, runtime.Object, error) {
					return true, existingJob, nil
				})
				dynClient.PrependReactor("delete", "jobs", func(action k8stesting.Action) (bool, runtime.Object, error) {
					return true, nil, nil
				})
				dynClient.PrependReactor("create", "jobs", func(action k8stesting.Action) (bool, runtime.Object, error) {
					return true, existingJob, nil
				})

				m.On("RESTConfigFromKubeConfig", kubeconfig).Return(config, nil).Once()
				m.On("NewDynamicForConfig", config).Return(dynClient, nil).Once()
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
			name: "failure in building kubeconfig",
			setupMocks: func(m *MockKubernetesClient, mockManifestConfig *MockManifestConfig) {
				m.On("RESTConfigFromKubeConfig", kubeconfig).Return(nil, fmt.Errorf("failed to build kubeconfig")).Once()
				mockManifestConfig.ReadFile = func(string) ([]byte, error) {
					return nil, nil
				}
			},
			jobFilePath: "testdata/job.yaml",
			namespace:   "default",
			expectError: true,
		},
		{
			name: "failure in creating dynamic client",
			setupMocks: func(m *MockKubernetesClient, mockManifestConfig *MockManifestConfig) {
				config := &rest.Config{}
				m.On("RESTConfigFromKubeConfig", kubeconfig).Return(config, nil).Once()
				m.On("NewDynamicForConfig", config).Return(nil, fmt.Errorf("failed to create dynamic client")).Once()
				mockManifestConfig.ReadFile = func(string) ([]byte, error) {
					return nil, nil
				}
			},
			jobFilePath: "testdata/job.yaml",
			namespace:   "default",
			expectError: true,
		},
		{
			name: "failure in applying manifest",
			setupMocks: func(m *MockKubernetesClient, mockManifestConfig *MockManifestConfig) {
				config := &rest.Config{}
				dynClient := fakedynamic.NewSimpleDynamicClient(runtime.NewScheme())
				m.On("RESTConfigFromKubeConfig", kubeconfig).Return(config, nil).Once()
				m.On("NewDynamicForConfig", config).Return(dynClient, nil).Once()
				mockManifestConfig.On("ApplyOrDeleteManifest", mock.Anything).Return(fmt.Errorf("failed to apply job")).Once()
				mockManifestConfig.ReadFile = func(string) ([]byte, error) {
					return []byte(`apiVersion: batch/v1
kind: Job
metadata:
  name: invalid-job
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
			jobFilePath: "testdata/invalid-job.yaml",
			namespace:   "default",
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := new(MockKubernetesClient)
			mockManifestConfig := new(MockManifestConfig)
			tc.setupMocks(mockClient, mockManifestConfig)

			kubeClient := &client.KubernetesClient{
				Config: &rest.Config{
					Host: kubeconfig,
				},
				Clientset: fake.NewSimpleClientset(),
			}

			jobsClient := &jobs.JobsClient{Client: kubeClient}

			require.NotNil(t, mockManifestConfig)

			err := jobsClient.ApplyKubernetesJob(tc.jobFilePath, tc.namespace, mockManifestConfig.ReadFile)
			require.Equal(t, tc.expectError, err != nil, "expected error: %v, got: %v", tc.expectError, err)

			if err != nil {
				t.Logf("error: %v", err)
			}

			// Clean up the job after test
			_ = jobsClient.DeleteKubernetesJob(context.Background(), "unique-job-1", tc.namespace)
			_ = jobsClient.DeleteKubernetesJob(context.Background(), "invalid-job", tc.namespace)

			mockClient.AssertExpectations(t)
			mockManifestConfig.AssertExpectations(t)
		})
	}
}
