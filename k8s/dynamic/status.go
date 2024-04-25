package k8s

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

// WaitForResourceReady waits for any Kubernetes resource to reach a ready state.
func WaitForResourceReady(ctx context.Context, resourceName, namespace, resourceType string, checkStatusFunc func(name, namespace string) (bool, error)) error {
	fmt.Printf("Waiting for %s '%s' in namespace '%s' to be ready...", resourceType, resourceName, namespace)
	timeout := time.After(5 * time.Minute)
	tick := time.NewTicker(10 * time.Second)
	defer tick.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled while waiting for %s '%s' to become ready", resourceType, resourceName)
		case <-timeout:
			return fmt.Errorf("timeout while waiting for %s '%s' to become ready", resourceType, resourceName)
		case <-tick.C:
			ready, err := checkStatusFunc(resourceName, namespace)
			if err != nil {
				fmt.Printf("Failed to get status for %s '%s': %v\n", resourceType, resourceName, err)
				continue // Skip to next tick
			}
			if ready {
				fmt.Printf("%s '%s' in namespace '%s' is ready.\n", resourceType, resourceName, namespace)
				return nil
			}
		}
	}
}

// GetResourceStatus checks the status of any Kubernetes resource.
func GetResourceStatus(ctx context.Context, client dynamic.Interface, resourceName, namespace string, gvr schema.GroupVersionResource) (bool, error) {
	resource, err := client.Resource(gvr).Namespace(namespace).Get(ctx, resourceName, metav1.GetOptions{})
	if err != nil {
		return false, fmt.Errorf("failed to get %s '%s' in namespace '%s': %v", gvr.Resource, resourceName, namespace, err)
	}

	status, found, err := unstructured.NestedString(resource.UnstructuredContent(), "status", "phase")
	if err != nil || !found {
		return false, fmt.Errorf("status not found for %s '%s': %v", gvr.Resource, resourceName, err)
	}

	return status == "Running", nil
}
