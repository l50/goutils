package k8s_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	client "github.com/l50/goutils/v2/k8s/client"
	loggers "github.com/l50/goutils/v2/k8s/loggers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	dynFake "k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
)

type MockDynamicClient struct {
	mock.Mock
	dynamic.Interface
	ResourceFn func(resource schema.GroupVersionResource) dynamic.NamespaceableResourceInterface
}

func (m *MockDynamicClient) Resource(resource schema.GroupVersionResource) dynamic.NamespaceableResourceInterface {
	if m.ResourceFn != nil {
		return m.ResourceFn(resource)
	}
	return nil
}

type MockNamespaceableResourceInterface struct {
	mock.Mock
	dynamic.NamespaceableResourceInterface
}

func (m *MockNamespaceableResourceInterface) Namespace(namespace string) dynamic.ResourceInterface {
	return m.Called(namespace).Get(0).(dynamic.ResourceInterface)
}

func NewMockNamespaceableResourceInterface(ctrl *gomock.Controller) *MockNamespaceableResourceInterface {
	return &MockNamespaceableResourceInterface{}
}

func mockFileReaderSvc(filename string) ([]byte, error) {
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

type MockResourceInterface struct {
	mock.Mock
	dynamic.ResourceInterface
}

func (m *MockResourceInterface) Namespace(namespace string) dynamic.ResourceInterface {
	return m.Called(namespace).Get(0).(dynamic.ResourceInterface)
}

func NewMockResourceInterface(ctrl *gomock.Controller) *MockResourceInterface {
	return &MockResourceInterface{}
}

func NewMockDynamicClient(ctrl *gomock.Controller) *MockDynamicClient {
	return &MockDynamicClient{}
}

func NewDynamicForConfig(config *rest.Config) (dynamic.Interface, error) {
	return nil, nil
}

type testCase struct {
	name          string
	namespace     string
	serviceName   string
	setupClient   func(*fake.Clientset, testCase)
	expectedError string
}

func TestServiceLogger(t *testing.T) {
	tests := []testCase{
		{
			name:        "successful service fetch and log",
			namespace:   "default",
			serviceName: "test-service",
			setupClient: func(cs *fake.Clientset, tc testCase) {
				service := &corev1.Service{
					ObjectMeta: metav1.ObjectMeta{Name: "test-service", Namespace: "default"},
				}
				_, err := cs.CoreV1().Services("default").Create(context.Background(), service, metav1.CreateOptions{})
				if err != nil {
					t.Fatalf("failed to create fake service: %v", err)
				}
			},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			clientset := fake.NewSimpleClientset()
			tc.setupClient(clientset, tc)

			mockDynamicClient := new(MockDynamicClient)
			mockDynamicClient.On("Resource", mock.Anything).Return(nil)
			dynamicClient := dynFake.NewSimpleDynamicClient(runtime.NewScheme())

			mockClient := new(MockKubernetesClient)
			mockClient.On("NewForConfig", mock.Anything).Return(clientset, nil)
			mockClient.On("NewDynamicForConfig", mock.Anything).Return(dynamicClient, nil)
			mockClient.On("RESTConfigFromKubeConfig", mock.Anything).Return(&rest.Config{}, nil)

			kc, err := client.NewKubernetesClient("fake-kubeconfig-path", mockFileReaderSvc, mockClient)
			assert.NoError(t, err)
			kc.Clientset = clientset
			kc.DynamicClient = mockDynamicClient

			serviceLogger := loggers.NewServiceLogger(kc, tc.namespace, tc.serviceName)
			err = serviceLogger.FetchAndLog(context.Background())
			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
			}
		})
	}
}
