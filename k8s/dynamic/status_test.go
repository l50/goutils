package k8s_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	client "github.com/l50/goutils/v2/k8s/client"
	k8s "github.com/l50/goutils/v2/k8s/dynamic"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/kubernetes/scheme"
	k8stesting "k8s.io/client-go/testing"
)

func TestWaitForResourceReady(t *testing.T) {
	tests := []struct {
		name            string
		resourceName    string
		namespace       string
		resourceType    string
		setupClient     func() dynamic.Interface
		checkStatusFunc func(resourceName, namespace string) (bool, error)
		expectedError   bool
		timeout         time.Duration
	}{
		{
			name:         "successful resource ready check",
			resourceName: "test-resource",
			namespace:    "default",
			resourceType: "pod",
			setupClient: func() dynamic.Interface {
				fakeClient := fake.NewSimpleDynamicClient(scheme.Scheme)
				return fakeClient
			},
			checkStatusFunc: func(name, namespace string) (bool, error) {
				// Simulate an immediate readiness without delay.
				return true, nil
			},
			expectedError: false,
			timeout:       10 * time.Second,
		},
		{
			name:         "resource ready check times out",
			resourceName: "test-resource",
			namespace:    "default",
			resourceType: "pod",
			setupClient: func() dynamic.Interface {
				fakeClient := fake.NewSimpleDynamicClient(scheme.Scheme)
				return fakeClient
			},
			checkStatusFunc: func(name, namespace string) (bool, error) {
				// Simulate not ready without delay.
				return false, nil
			},
			expectedError: true,
			timeout:       100 * time.Millisecond, // Timeout intended for failure scenario
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tc.timeout)
			defer cancel()

			err := k8s.WaitForResourceReady(ctx, tc.resourceName, tc.namespace, tc.resourceType, func(name, namespace string) (bool, error) {
				return tc.checkStatusFunc(name, namespace)
			})

			if tc.expectedError {
				assert.Error(t, err, "Expected an error but did not get one in test case: %s", tc.name)
			} else {
				assert.NoError(t, err, "Did not expect an error but got one in test case: %s", tc.name)
			}
		})
	}
}

func TestDescribeKubernetesResource(t *testing.T) {
	tests := []struct {
		name           string
		resourceName   string
		namespace      string
		gvr            schema.GroupVersionResource
		resourceSetup  func() *unstructured.Unstructured
		expectedOutput string
		expectError    bool
	}{
		{
			name:         "successful resource description",
			resourceName: "test-pod",
			namespace:    "default",
			gvr:          schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"},
			resourceSetup: func() *unstructured.Unstructured {
				return &unstructured.Unstructured{
					Object: map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "Pod",
						"metadata": map[string]interface{}{
							"name":      "test-pod",
							"namespace": "default",
						},
						"status": map[string]interface{}{
							"phase": "Running",
						},
					},
				}
			},
			expectedOutput: "Name: test-pod\nNamespace: default\nLabels: map[]\nAnnotations: map[]\nDetails:\napiVersion: v1\nkind: Pod\nmetadata: map[name:test-pod namespace:default]\nstatus: map[phase:Running]\n",
			expectError:    false,
		},
		{
			name:         "resource not found",
			resourceName: "unknown-pod",
			namespace:    "default",
			gvr:          schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"},
			resourceSetup: func() *unstructured.Unstructured {
				return nil // correctly simulate non-existent resource
			},
			expectedOutput: "",
			expectError:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			// Setup fake dynamic client and wrap in KubernetesClient
			fakeDynamicClient := fake.NewSimpleDynamicClient(scheme.Scheme)
			kubernetesClient := &client.KubernetesClient{
				DynamicClient: fakeDynamicClient,
				// Populate other fields if needed, such as Clientset or Config
			}

			// Handling the "get" reactor for specific tests
			fakeDynamicClient.PrependReactor("get", "pods", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
				getAction, ok := action.(k8stesting.GetAction)
				if ok && getAction.GetName() == tc.resourceName && tc.resourceSetup != nil {
					resource := tc.resourceSetup()
					if resource == nil {
						return true, nil, fmt.Errorf("pod '%s' not found", getAction.GetName())
					}
					return true, resource, nil
				}
				return false, nil, fmt.Errorf("pod '%s' not found", getAction.GetName())
			})

			description, err := k8s.DescribeKubernetesResource(ctx, kubernetesClient, tc.resourceName, tc.namespace, tc.gvr)

			if tc.expectError {
				assert.Error(t, err, "Expected an error but did not get one in test case: %s", tc.name)
			} else {
				assert.NoError(t, err, "Did not expect an error but got one in test case: %s", tc.name)
				assert.Equal(t, tc.expectedOutput, description, "Unexpected description in test case: %s", tc.name)
			}
		})
	}
}
