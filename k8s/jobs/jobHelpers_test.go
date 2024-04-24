package k8s_test

import (
	"context"
	"testing"

	k8s "github.com/l50/goutils/v2/k8s/client"
	jobs "github.com/l50/goutils/v2/k8s/jobs"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestDeleteKubernetesJob(t *testing.T) {
	tests := []struct {
		name        string
		jobName     string
		namespace   string
		setupClient func() *jobs.JobsClient
		expectError bool
	}{
		{
			name:      "successful job deletion",
			jobName:   "test-job",
			namespace: "default",
			setupClient: func() *jobs.JobsClient {
				fakeClient := fake.NewSimpleClientset(&batchv1.Job{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-job",
						Namespace: "default",
					},
				})
				return &jobs.JobsClient{Client: &k8s.KubernetesClient{Clientset: fakeClient}}
			},
			expectError: false,
		},
		{
			name:      "failed job deletion",
			jobName:   "nonexistent-job",
			namespace: "default",
			setupClient: func() *jobs.JobsClient {
				fakeClient := fake.NewSimpleClientset() // No pre-existing job
				return &jobs.JobsClient{Client: &k8s.KubernetesClient{Clientset: fakeClient}}
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			jobsClient := tc.setupClient()
			err := jobsClient.DeleteKubernetesJob(ctx, tc.jobName, tc.namespace)
			if (err != nil) != tc.expectError {
				t.Errorf("Test %s: expected error: %v, got: %v", tc.name, tc.expectError, err)
			}
		})
	}
}
