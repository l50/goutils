package k8s_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"k8s.io/apimachinery/pkg/runtime"
	dynFake "k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"

	client "github.com/l50/goutils/v2/k8s/client"
	dynK8s "github.com/l50/goutils/v2/k8s/dynamic"
)

// Mocking the SPDY Executor for remotecommand.Stream.
type mockExecutor struct {
	mock.Mock
}

func (m *mockExecutor) StreamWithContext(ctx context.Context, opts remotecommand.StreamOptions) error {
	return m.Called(ctx, opts).Error(0)
}

func TestExecKubernetesResources(t *testing.T) {
	config := &rest.Config{Host: "https://localhost:6443"}
	mockClientset := fake.NewSimpleClientset()
	mockDynamicClient := dynFake.NewSimpleDynamicClient(runtime.NewScheme())

	kc := &client.KubernetesClient{
		Clientset:     mockClientset,
		DynamicClient: mockDynamicClient,
		Config:        config,
	}

	tests := []struct {
		name      string
		params    dynK8s.ExecParams
		expectErr bool
		expected  string
	}{
		{
			name: "successful command execution",
			params: dynK8s.ExecParams{
				Context:   context.TODO(),
				Client:    kc,
				Namespace: "default",
				PodName:   "example-pod",
				Command:   []string{"echo", "Hello"},
				Stdout:    new(bytes.Buffer),
			},
			expectErr: false,
			expected:  "Hello",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockExec := new(mockExecutor)
			mockExec.On("StreamWithContext", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("remotecommand.StreamOptions")).Return(nil)
			_, err := dynK8s.ExecKubernetesResources(tc.params)

			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Contains(t, tc.params.Stdout.(*bytes.Buffer).String(), tc.expected, "Output should contain the expected command output")
			}
		})
	}
}
