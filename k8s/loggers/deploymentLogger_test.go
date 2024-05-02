package k8s_test

import (
	"context"
	"testing"

	client "github.com/l50/goutils/v2/k8s/client"
	loggers "github.com/l50/goutils/v2/k8s/loggers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	dynFake "k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
)

func mockFileReader(filename string) ([]byte, error) {
	return []byte(`{
		"apiVersion": "v1",
		"kind": "Config",
		"clusters": [{
			"cluster": {
				"server": "https://fake-server"
			},
			"name": "fake-cluster"
		}],
		"contexts": [{
			"context": {
				"cluster": "fake-cluster",
				"user": "fake-user"
			},
			"name": "fake"
		}],
		"current-context": "fake",
		"users": [{
			"name": "fake-user",
			"user": {}
		}]
	}`), nil
}

type MockKubernetesClient struct {
	mock.Mock
	client.KubernetesClientInterface
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

func TestDeploymentLogger(t *testing.T) {
	tests := []struct {
		name           string
		namespace      string
		deploymentName string
		setupClient    func(*fake.Clientset)
		expectedError  string
	}{
		{
			name:           "successful deployment fetch and log",
			namespace:      "default",
			deploymentName: "test-deployment",
			setupClient: func(cs *fake.Clientset) {
				deployment := &appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{Name: "test-deployment", Namespace: "default"},
				}
				_, err := cs.AppsV1().Deployments("default").Create(context.Background(), deployment, metav1.CreateOptions{})
				if err != nil {
					t.Fatalf("failed to create fake deployment: %v", err)
				}
			},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			clientset := fake.NewSimpleClientset()
			dynamicClient := dynFake.NewSimpleDynamicClient(runtime.NewScheme())
			tc.setupClient(clientset)

			mockClient := new(MockKubernetesClient)
			mockClient.On("NewForConfig", mock.Anything).Return(clientset, nil)
			mockClient.On("NewDynamicForConfig", mock.Anything).Return(dynamicClient, nil)
			mockClient.On("RESTConfigFromKubeConfig", mock.Anything).Return(&rest.Config{}, nil)

			kc, err := client.NewKubernetesClient("fake-kubeconfig-path", mockFileReader, mockClient)
			assert.NoError(t, err)

			deploymentLogger := loggers.NewDeploymentLogger(kc, tc.namespace, tc.deploymentName)
			err = deploymentLogger.FetchAndLog(context.Background())
			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
			}
		})
	}
}
