package k8s_test

import (
	"fmt"
	"testing"

	"github.com/l50/goutils/v2/k8s"
)

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
				if path == tc.kubeconfig {
					return tc.data, nil
				}
				return nil, fmt.Errorf("file not found")
			}

			_, err := k8s.NewKubernetesClient(tc.kubeconfig, reader)
			if (err != nil) != tc.expectError {
				t.Errorf("Test '%s' failed - expected error: %v, got: %v", tc.name, tc.expectError, err)
			}
		})
	}
}