package k8s_test

import (
	"fmt"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	fake "k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

type fakeError []byte

func (f fakeError) Error() string {
	return string(f)
}

func TestFetchAndLogPods(t *testing.T) {
	tests := []struct {
		name           string
		namespace      string
		labelSelector  string
		mockPods       []*corev1.Pod
		logContent     string
		expectedOutput string
		expectedError  string
	}{
		{
			name:          "error fetching pod logs",
			namespace:     "default",
			labelSelector: "app=error",
			mockPods: []*corev1.Pod{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:   "pod-error",
						Labels: map[string]string{"app": "error"},
					},
				},
			},
			logContent:     "log data for pod-error",
			expectedOutput: "Attempting to list pods with label selector: 'app=error' in namespace 'default'\nFetching logs for pod: pod-error\nLogs for pod pod-error:\nlog data for pod-error\n",
			expectedError:  "",
		},
		{
			name:          "pods with empty logs",
			namespace:     "default",
			labelSelector: "app=empty",
			mockPods: []*corev1.Pod{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:   "pod-empty",
						Labels: map[string]string{"app": "empty"},
					},
				},
			},
			logContent:     "",
			expectedOutput: "Attempting to list pods with label selector: 'app=empty' in namespace 'default'\nFetching logs for pod: pod-empty\nNo logs for pod pod-empty\n",
			expectedError:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			clientset := fake.NewSimpleClientset(podsToRuntimeObjects(tc.mockPods)...)
			clientset.PrependReactor("list", "pods", func(action k8stesting.Action) (bool, runtime.Object, error) {
				// Type assertion to ListAction
				if listAction, ok := action.(k8stesting.ListAction); ok {
					var pods []*corev1.Pod
					for _, pod := range tc.mockPods {
						if listAction.GetListRestrictions().Labels.Matches(labels.Set(pod.Labels)) {
							pods = append(pods, pod)
						}
					}
					return true, &corev1.PodList{Items: podsToCoreV1Pods(pods)}, nil
				}
				return false, nil, fmt.Errorf("unexpected action type: %T", action)
			})

			clientset.PrependReactor("get", "pods/log", func(action k8stesting.Action) (bool, runtime.Object, error) {
				// Type assertion to GetAction
				if getAction, ok := action.(k8stesting.GetAction); ok {
					for _, pod := range tc.mockPods {
						if getAction.GetName() == pod.Name {
							return true, nil, fakeError([]byte(tc.logContent))
						}
					}
					return true, nil, fmt.Errorf("pod not found")
				}
				return false, nil, fmt.Errorf("unexpected action type: %T", action)
			})
		})
	}
}

func podsToCoreV1Pods(pods []*corev1.Pod) []corev1.Pod {
	result := make([]corev1.Pod, len(pods))
	for i, pod := range pods {
		result[i] = *pod
	}
	return result
}
func podsToRuntimeObjects(pods []*corev1.Pod) []runtime.Object {
	objs := make([]runtime.Object, len(pods))
	for i, pod := range pods {
		objs[i] = pod
	}
	return objs
}
