package k8s_test

import (
	"context"
	"testing"
	"time"

	k8s "github.com/l50/goutils/v2/k8s/dynamic"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/kubernetes/scheme"
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
