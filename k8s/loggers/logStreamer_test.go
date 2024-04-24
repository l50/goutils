package k8s_test

import (
	"context"
	"io"
	"os"
	"testing"

	k8s "github.com/l50/goutils/v2/k8s/loggers"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

func TestStreamLogs(t *testing.T) {
	tests := []struct {
		name          string
		namespace     string
		resourceType  string
		resourceName  string
		setupClient   func(cs *fake.Clientset)
		expectedError string
		expectOutput  bool
	}{
		{
			name:         "stream logs from pod",
			namespace:    "default",
			resourceType: "pod",
			resourceName: "test-pod",
			setupClient: func(cs *fake.Clientset) {
				if _, err := cs.CoreV1().Pods("default").Create(context.TODO(), &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{Name: "test-pod", Namespace: "default"},
				}, metav1.CreateOptions{}); err != nil {
					t.Fatalf("failed to create pod: %v", err)
				}
				cs.PrependReactor("get", "pods", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
					return true, &corev1.Pod{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-pod",
							Namespace: "default",
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{{Name: "container"}},
						},
					}, nil
				})
			},
			expectedError: "",
			expectOutput:  true,
		},
		{
			name:          "stream logs from job with no pods",
			namespace:     "default",
			resourceType:  "job",
			resourceName:  "test-job",
			setupClient:   func(cs *fake.Clientset) {},
			expectedError: "no pods found for job: test-job",
			expectOutput:  false,
		},
		{
			name:          "unsupported resource type",
			namespace:     "default",
			resourceType:  "unsupported",
			resourceName:  "test",
			setupClient:   func(cs *fake.Clientset) {},
			expectedError: "unsupported resource type: unsupported",
			expectOutput:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			clientset := fake.NewSimpleClientset()
			if tc.setupClient != nil {
				tc.setupClient(clientset)
			}

			// Redirect stdout to capture output for verification
			rescueStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := k8s.StreamLogs(clientset, tc.namespace, tc.resourceType, tc.resourceName)

			// Restore stdout
			w.Close()
			out, _ := io.ReadAll(r)
			os.Stdout = rescueStdout

			if tc.expectedError == "" {
				assert.NoError(t, err)
				if tc.expectOutput {
					assert.NotEmpty(t, out, "Expected log output but got none")
				}
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError, "Expected an error containing specific text")
			}
		})
	}
}
