package k8s_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	client "github.com/l50/goutils/v2/k8s/client"
	dynK8s "github.com/l50/goutils/v2/k8s/dynamic"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/kubernetes/scheme"
	k8stesting "k8s.io/client-go/testing"
)

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

			description, err := dynK8s.DescribeKubernetesResource(ctx, kubernetesClient, tc.resourceName, tc.namespace, tc.gvr)

			if tc.expectError {
				assert.Error(t, err, "Expected an error but did not get one in test case: %s", tc.name)
			} else {
				assert.NoError(t, err, "Did not expect an error but got one in test case: %s", tc.name)
				assert.Equal(t, tc.expectedOutput, description, "Unexpected description in test case: %s", tc.name)
			}
		})
	}
}

func TestGetResourceStatus(t *testing.T) {
	tests := []struct {
		name          string
		resourceName  string
		namespace     string
		gvr           schema.GroupVersionResource
		setupPod      func() *unstructured.Unstructured
		expectedState bool
		expectError   bool
	}{
		{
			name:         "pod in Running state",
			resourceName: "running-pod",
			namespace:    "default",
			gvr:          schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"},
			setupPod: func() *unstructured.Unstructured {
				return &unstructured.Unstructured{
					Object: map[string]interface{}{
						"status": map[string]interface{}{
							"phase": "Running",
						},
					},
				}
			},
			expectedState: true,
			expectError:   false,
		},
		{
			name:         "pod in Failed state",
			resourceName: "failed-pod",
			namespace:    "default",
			gvr:          schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"},
			setupPod: func() *unstructured.Unstructured {
				return &unstructured.Unstructured{
					Object: map[string]interface{}{
						"status": map[string]interface{}{
							"phase": "Failed",
						},
					},
				}
			},
			expectedState: false,
			expectError:   false,
		},
		{
			name:         "pod in OOMKilled state",
			resourceName: "oomkilled-pod",
			namespace:    "default",
			gvr:          schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"},
			setupPod: func() *unstructured.Unstructured {
				return &unstructured.Unstructured{
					Object: map[string]interface{}{
						"status": map[string]interface{}{
							"phase": "OOMKilled",
						},
					},
				}
			},
			expectedState: false,
			expectError:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			// Setup fake dynamic client and wrap in KubernetesClient
			fakeDynamicClient := fake.NewSimpleDynamicClient(scheme.Scheme)
			kubernetesClient := &client.KubernetesClient{
				DynamicClient: fakeDynamicClient,
			}

			// Prepend the 'get' reactor for the dynamic client
			fakeDynamicClient.PrependReactor("get", "pods", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
				if action.(k8stesting.GetAction).GetName() == tc.resourceName {
					return true, tc.setupPod(), nil
				}
				return false, nil, fmt.Errorf("pod '%s' not found", action.(k8stesting.GetAction).GetName())
			})

			// Perform the status check
			status, err := dynK8s.GetResourceStatus(ctx, kubernetesClient, tc.resourceName, tc.namespace, tc.gvr)

			assert.Equal(t, tc.expectedState, status, "Expected state did not match for test case: %s", tc.name)
			if tc.expectError {
				assert.Error(t, err, "Expected an error but did not get one in test case: %s", tc.name)
			} else {
				assert.NoError(t, err, "Did not expect an error but got one in test case: %s", tc.name)
			}
		})
	}
}

func TestWaitForResourceState(t *testing.T) {
	tests := []struct {
		name            string
		resourceName    string
		namespace       string
		resourceType    string
		desiredState    string
		setupClient     func() dynamic.Interface
		checkStatusFunc func(dynamic.Interface, string, string) (bool, error)
		expectedError   bool
		timeout         time.Duration
	}{
		{
			name:         "successful resource state check for readiness",
			resourceName: "test-resource",
			namespace:    "default",
			resourceType: "pod",
			desiredState: "Running",
			setupClient: func() dynamic.Interface {
				return fake.NewSimpleDynamicClient(scheme.Scheme)
			},
			checkStatusFunc: func(client dynamic.Interface, name, namespace string) (bool, error) {
				return true, nil // Simulate an immediate readiness without delay.
			},
			expectedError: false,
			timeout:       2 * time.Second,
		},
		{
			name:         "resource state check times out for readiness",
			resourceName: "test-resource",
			namespace:    "default",
			resourceType: "pod",
			desiredState: "Running",
			setupClient: func() dynamic.Interface {
				return fake.NewSimpleDynamicClient(scheme.Scheme)
			},
			checkStatusFunc: func(client dynamic.Interface, name, namespace string) (bool, error) {
				return false, nil
			},
			expectedError: true,
			timeout:       20 * time.Millisecond,
		},
		{
			name:         "successful resource state check for removal",
			resourceName: "test-resource",
			namespace:    "default",
			resourceType: "pod",
			desiredState: "Deleted",
			setupClient: func() dynamic.Interface {
				fakeClient := fake.NewSimpleDynamicClient(scheme.Scheme)
				fakeClient.PrependReactor("get", "pods", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
					getAction := action.(k8stesting.GetAction)
					if getAction.GetName() == "test-resource" {
						return true, nil, errors.NewNotFound(schema.GroupResource{Group: "", Resource: "pods"}, "test-resource")
					}
					return false, nil, nil
				})
				return fakeClient
			},
			checkStatusFunc: func(client dynamic.Interface, name, namespace string) (bool, error) {
				_, err := client.Resource(schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}).Namespace(namespace).Get(context.TODO(), name, metav1.GetOptions{})
				if errors.IsNotFound(err) {
					return true, nil
				}
				return false, err
			},
			expectedError: false,
			timeout:       20 * time.Second,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tc.timeout)
			defer cancel()

			client := tc.setupClient() // Setup client for each test

			err := dynK8s.WaitForResourceState(ctx, tc.resourceName, tc.namespace, tc.resourceType, tc.desiredState, func(name, namespace string) (bool, error) {
				return tc.checkStatusFunc(client, name, namespace) // Pass the client directly
			})

			if tc.expectedError {
				assert.Error(t, err, "Expected an error but did not get one in test case: %s", tc.name)
			} else {
				assert.NoError(t, err, "Did not expect an error but got one in test case: %s", tc.name)
			}
		})
	}
}
