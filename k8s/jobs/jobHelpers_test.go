package k8s_test

import (
	"context"
	"fmt"
	"testing"

	k8s "github.com/l50/goutils/v2/k8s/client"
	jobs "github.com/l50/goutils/v2/k8s/jobs"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
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

func TestListKubernetesJobs(t *testing.T) {
	tests := []struct {
		name         string
		namespace    string
		setupClient  func() *jobs.JobsClient
		expectedJobs int
		expectError  bool
	}{
		{
			name:      "list jobs from a specific namespace",
			namespace: "default",
			setupClient: func() *jobs.JobsClient {
				fakeClient := fake.NewSimpleClientset(
					&corev1.Pod{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-job-pod-123",
							Namespace: "default",
						},
					},
					&batchv1.Job{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-job-1",
							Namespace: "default",
						},
					},
					&batchv1.Job{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-job-2",
							Namespace: "default",
						},
					},
				)
				return &jobs.JobsClient{Client: &k8s.KubernetesClient{Clientset: fakeClient}}
			},
			expectedJobs: 2,
			expectError:  false,
		},
		{
			name:      "list jobs from all namespaces",
			namespace: "",
			setupClient: func() *jobs.JobsClient {
				fakeClient := fake.NewSimpleClientset(
					&batchv1.Job{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-job-1",
							Namespace: "default",
						},
					},
					&batchv1.Job{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-job-2",
							Namespace: "test",
						},
					},
				)
				return &jobs.JobsClient{Client: &k8s.KubernetesClient{Clientset: fakeClient}}
			},
			expectedJobs: 2,
			expectError:  false,
		},
		{
			name:      "error fetching jobs",
			namespace: "default",
			setupClient: func() *jobs.JobsClient {
				fakeClient := fake.NewSimpleClientset()
				fakeClient.PrependReactor("list", "jobs", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
					return true, nil, fmt.Errorf("failed to list jobs")
				})
				return &jobs.JobsClient{Client: &k8s.KubernetesClient{Clientset: fakeClient}}
			},
			expectedJobs: 0,
			expectError:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			jobsClient := tc.setupClient()
			jobs, err := jobsClient.ListKubernetesJobs(ctx, tc.namespace)
			if (err != nil) != tc.expectError {
				t.Errorf("Test '%s': expected error: %v, got: %v", tc.name, tc.expectError, err)
			}
			if len(jobs) != tc.expectedJobs {
				t.Errorf("Test '%s': expected number of jobs: %d, got: %d", tc.name, tc.expectedJobs, len(jobs))
			}
		})
	}
}

func TestJobExists(t *testing.T) {
	tests := []struct {
		name         string
		jobName      string
		namespace    string
		setupClient  func() *jobs.JobsClient
		expectExists bool
		expectError  bool
	}{
		{
			name:      "job exists",
			jobName:   "existing-job",
			namespace: "default",
			setupClient: func() *jobs.JobsClient {
				fakeClient := fake.NewSimpleClientset(&batchv1.Job{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "existing-job",
						Namespace: "default",
					},
				})
				return &jobs.JobsClient{Client: &k8s.KubernetesClient{Clientset: fakeClient}}
			},
			expectExists: true,
			expectError:  false,
		},
		{
			name:      "job does not exist",
			jobName:   "nonexistent-job",
			namespace: "default",
			setupClient: func() *jobs.JobsClient {
				fakeClient := fake.NewSimpleClientset()
				return &jobs.JobsClient{Client: &k8s.KubernetesClient{Clientset: fakeClient}}
			},
			expectExists: false,
			expectError:  false,
		},
		{
			name:      "error on retrieving job",
			jobName:   "faulty-job",
			namespace: "default",
			setupClient: func() *jobs.JobsClient {
				fakeClient := fake.NewSimpleClientset()
				fakeClient.PrependReactor("get", "jobs", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
					return true, nil, fmt.Errorf("internal error")
				})
				return &jobs.JobsClient{Client: &k8s.KubernetesClient{Clientset: fakeClient}}
			},
			expectExists: false,
			expectError:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			jobsClient := tc.setupClient()
			exists, err := jobsClient.JobExists(ctx, tc.jobName, tc.namespace)
			if (err != nil) != tc.expectError {
				t.Errorf("Test %s: expected error: %v, got: %v", tc.name, tc.expectError, err)
			}
			if exists != tc.expectExists {
				t.Errorf("Test %s: expected job existence: %v, got: %v", tc.name, tc.expectExists, exists)
			}
		})
	}
}
