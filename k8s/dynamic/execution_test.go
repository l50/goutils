package k8s_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	ktesting "k8s.io/client-go/testing"

	client "github.com/l50/goutils/v2/k8s/client"
	dynK8s "github.com/l50/goutils/v2/k8s/dynamic"
)

func TestExecKubernetesResources(t *testing.T) {
	tests := []struct {
		name          string
		namespace     string
		podName       string
		setupClient   func(*fake.Clientset)
		expectedError string
	}{
		{
			name:      "PodNotFound",
			namespace: "default",
			podName:   "my-pod",
			setupClient: func(cs *fake.Clientset) {
				cs.PrependReactor("get", "pods", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					return true, nil, errors.New("pods \"my-pod\" not found")
				})
			},
			expectedError: "pods \"my-pod\" not found",
		},
		{
			name:      "PodNotRunning",
			namespace: "default",
			podName:   "my-pod",
			setupClient: func(cs *fake.Clientset) {
				pod := &v1.Pod{
					ObjectMeta: metav1.ObjectMeta{Name: "my-pod", Namespace: "default"},
					Status:     v1.PodStatus{Phase: v1.PodPending},
				}
				_, err := cs.CoreV1().Pods("default").Create(context.Background(), pod, metav1.CreateOptions{})
				if err != nil {
					t.Fatalf("failed to create pod: %v", err)
				}
			},
			expectedError: "pod my-pod is not in running state, current state: Pending",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := fake.NewSimpleClientset()
			tc.setupClient(mockClient)

			params := dynK8s.ExecParams{
				Context:   context.Background(),
				Client:    &client.KubernetesClient{Clientset: mockClient},
				Namespace: tc.namespace,
				PodName:   tc.podName,
				Command:   []string{"ls", "-l"},
			}

			output, err := dynK8s.ExecKubernetesResources(params)

			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
				assert.NotContains(t, output, "Command executed successfully")
			} else {
				assert.NoError(t, err)
				assert.Contains(t, output, "Command executed successfully")
			}
		})
	}
}
