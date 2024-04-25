package k8s_test

import (
	"context"
	"testing"

	k8s "github.com/l50/goutils/v2/k8s/client"
	jobs "github.com/l50/goutils/v2/k8s/jobs"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestGetJobPodName(t *testing.T) {
	tests := []struct {
		name        string
		jobName     string
		namespace   string
		setupClient func() *jobs.JobsClient
		expectedPod string
		expectError bool
	}{
		{
			name:      "successful retrieval of pod name",
			jobName:   "test-job",
			namespace: "default",
			setupClient: func() *jobs.JobsClient {
				fakeClient := fake.NewSimpleClientset(&corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-job-pod-123",
						Namespace: "default",
						Labels:    map[string]string{"job-name": "test-job"},
					},
				})
				return &jobs.JobsClient{Client: &k8s.KubernetesClient{Clientset: fakeClient}}
			},
			expectedPod: "test-job-pod-123",
			expectError: false,
		},
		{
			name:      "failed retrieval due to no pods found",
			jobName:   "test-job",
			namespace: "default",
			setupClient: func() *jobs.JobsClient {
				fakeClient := fake.NewSimpleClientset() // No pods created
				return &jobs.JobsClient{Client: &k8s.KubernetesClient{Clientset: fakeClient}}
			},
			expectedPod: "",
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			jobsClient := tc.setupClient()
			podName, err := jobsClient.GetJobPodName(ctx, tc.jobName, tc.namespace)
			if (err != nil) != tc.expectError {
				t.Errorf("Test %s: expected error: %v, got: %v", tc.name, tc.expectError, err)
			}
			if podName != tc.expectedPod {
				t.Errorf("Test %s: expected pod name: %s, got: %s", tc.name, tc.expectedPod, podName)
			}
		})
	}
}
