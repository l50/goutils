package k8s_test

import (
	"context"
	"fmt"
	"testing"

	k8s "github.com/l50/goutils/v2/k8s/client"
	jobs "github.com/l50/goutils/v2/k8s/jobs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
)

type MockDynK8s struct {
	mock.Mock
}

type MockK8sLogger struct {
	mock.Mock
}

type MockKubernetesClient struct {
	mock.Mock
	kubernetes.Interface
}

type MockJobPodNameGetter struct {
	mock.Mock
}

func (m *MockDynK8s) WaitForResourceState(ctx context.Context, workloadName, namespace, resourceType, state string, checkFunc func(string, string) (bool, error)) error {
	args := m.Called(ctx, workloadName, namespace, resourceType, state, checkFunc)
	if checkFunc != nil {
		fmt.Printf("Calling checkFunc for workloadName: %s, namespace: %s\n", workloadName, namespace)
		_, err := checkFunc(workloadName, namespace)
		if err != nil {
			return err
		}
	}
	return args.Error(0)
}

func (m *MockDynK8s) GetResourceStatus(ctx context.Context, client *k8s.KubernetesClient, name, ns string, gvr schema.GroupVersionResource) (bool, error) {
	args := m.Called(ctx, client, name, ns, gvr)
	fmt.Printf("Mock GetResourceStatus called for name: %s, namespace: %s\n", name, ns)
	return args.Bool(0), args.Error(1)
}

func (m *MockK8sLogger) StreamLogs(clientset kubernetes.Interface, namespace, resourceType, podName string) error {
	args := m.Called(clientset, namespace, resourceType, podName)
	fmt.Printf("Mock StreamLogs called for podName: %s, namespace: %s\n", podName, namespace)
	return args.Error(0)
}

func (m *MockJobPodNameGetter) GetJobPodName(ctx context.Context, workloadName, namespace string) (string, error) {
	args := m.Called(ctx, workloadName, namespace)
	fmt.Printf("Mock GetJobPodName called for workloadName: %s, namespace: %s\n", workloadName, namespace)
	return args.String(0), args.Error(1)
}

func TestStreamJobLogs(t *testing.T) {
	mockDynK8s := new(MockDynK8s)
	mockK8sLogger := new(MockK8sLogger)
	mockClientset := new(MockKubernetesClient)
	mockJobPodNameGetter := new(MockJobPodNameGetter)

	tests := []struct {
		name                   string
		workloadName           string
		namespace              string
		waitForResourceError   error
		getResourceStatus      bool
		getResourceStatusError error
		streamLogsError        error
		getJobPodName          string
		getJobPodNameError     error
		expectedError          bool
	}{
		{
			name:                   "successful log streaming",
			workloadName:           "test-job",
			namespace:              "test-namespace",
			waitForResourceError:   nil,
			getResourceStatus:      true,
			getResourceStatusError: nil,
			streamLogsError:        nil,
			getJobPodName:          "test-pod",
			getJobPodNameError:     nil,
			expectedError:          false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock expectations
			mockDynK8s.On("WaitForResourceState", mock.Anything, tc.workloadName, tc.namespace, "job", "Complete", mock.Anything).Return(tc.waitForResourceError)
			mockDynK8s.On("GetResourceStatus", mock.Anything, mock.Anything, tc.workloadName, tc.namespace, mock.Anything).Return(tc.getResourceStatus, tc.getResourceStatusError)
			mockK8sLogger.On("StreamLogs", mock.Anything, tc.namespace, "pod", tc.getJobPodName).Return(tc.streamLogsError)
			mockJobPodNameGetter.On("GetJobPodName", mock.Anything, tc.workloadName, tc.namespace).Return(tc.getJobPodName, tc.getJobPodNameError)

			jc := &jobs.JobsClient{
				Client:        &k8s.KubernetesClient{Clientset: mockClientset},
				DynK8s:        mockDynK8s,
				K8sLogger:     mockK8sLogger,
				PodNameGetter: mockJobPodNameGetter,
			}

			err := jc.StreamJobLogs(tc.workloadName, tc.namespace)
			if tc.expectedError {
				assert.Error(t, err, "Expected error while streaming job logs")
				assert.EqualError(t, err, tc.streamLogsError.Error(), "Error message does not match")
			} else {
				assert.NoError(t, err, "Expected no error while streaming job logs")
			}

			mockDynK8s.AssertExpectations(t)
			mockK8sLogger.AssertExpectations(t)
			mockJobPodNameGetter.AssertExpectations(t)
		})
	}
}
