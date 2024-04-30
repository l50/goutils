package k8s_test

import (
	"context"
	"net/url"
	"testing"

	client "github.com/l50/goutils/v2/k8s/client"
	dynK8s "github.com/l50/goutils/v2/k8s/dynamic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/client-go/util/flowcontrol"
)

type mockSPDYExecutor struct {
	mock.Mock
}

func (m *mockSPDYExecutor) Stream(options remotecommand.StreamOptions) error {
	return m.Called(options).Error(0)
}

func (m *mockSPDYExecutor) StreamWithContext(ctx context.Context, opts remotecommand.StreamOptions) error {
	return m.Called(ctx, opts).Error(0)
}

type mockRestClient struct {
	mock.Mock
}

func (c *mockRestClient) Post() *rest.Request {
	req := rest.NewRequestWithClient(&url.URL{}, "", rest.ClientContentConfig{}, nil)
	c.Mock.Called()
	return req
}

func (m *mockRestClient) Delete() *rest.Request {
	return &rest.Request{}
}

func (m *mockRestClient) Put() *rest.Request {
	return &rest.Request{}
}

func (m *mockRestClient) Get() *rest.Request {
	return &rest.Request{}
}

func (c *mockRestClient) Patch(pt types.PatchType) *rest.Request {
	return &rest.Request{}
}

func (c *mockRestClient) Verb(verb string) *rest.Request {
	return &rest.Request{}
}

func (m *mockRestClient) APIVersion() schema.GroupVersion {
	return schema.GroupVersion{Group: "testgroup", Version: "v1"}
}

func (c *mockRestClient) GetRateLimiter() flowcontrol.RateLimiter {
	return flowcontrol.NewTokenBucketRateLimiter(1, 1)
}

type MockExecutorCreator struct {
	mock.Mock
}

func (m *MockExecutorCreator) NewSPDYExecutor(config *rest.Config, method string, url *url.URL) (remotecommand.Executor, error) {
	args := m.Called(config, method, url)
	return args.Get(0).(remotecommand.Executor), args.Error(1)
}

func TestExecKubernetesResources(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	config := &rest.Config{Host: "https://localhost:6443"}
	kc := &client.KubernetesClient{Clientset: clientset, Config: config}

	tests := []struct {
		name      string
		kc        *client.KubernetesClient
		namespace string
		podName   string
		command   []string
		expected  string
		expectErr bool
	}{
		{
			name:      "successful command execution",
			kc:        kc,
			namespace: "default",
			podName:   "test-pod",
			command:   []string{"echo", "hello"},
			expected:  "Command executed successfully",
			expectErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRestClient := new(mockRestClient)
			mockExecutor := new(mockSPDYExecutor)
			mockExecutorCreator := new(MockExecutorCreator)
			mockRestClient.On("Post").Return(&rest.Request{})
			mockExecutor.On("StreamWithContext", mock.Anything, mock.Anything).Return(nil) // Simulate success
			mockExecutorCreator.On("NewSPDYExecutor", mock.Anything, mock.Anything, mock.Anything).Return(mockExecutor, nil)

			result, err := dynK8s.ExecKubernetesResources(context.Background(), tc.kc, tc.namespace, tc.podName, tc.command, mockRestClient, mockExecutorCreator)
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}
