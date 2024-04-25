package k8s

import (
	"context"
	"fmt"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

// WaitForResourceReady waits for any Kubernetes resource to reach a ready state.
//
// **Parameters:**
//
// ctx: A context.Context to allow for cancellation and timeouts.
// resourceName: The name of the resource to monitor.
// namespace: The namespace in which the resource exists.
// resourceType: The type of the resource (e.g., Pod, Service).
// checkStatusFunc: A function that checks if the resource is ready.
//
// **Returns:**
//
// error: An error if the waiting is cancelled by context, times out, or
// fails to determine readiness.
func WaitForResourceReady(ctx context.Context, resourceName, namespace, resourceType string, checkStatusFunc func(name, namespace string) (bool, error)) error {
	fmt.Printf("Waiting for %s (%s) in %s namespace to be ready...\n", resourceName, resourceType, namespace)

	// Set a timeout for resource readiness
	timeout := time.After(5 * time.Minute)

	// Check status every 10 seconds
	tick := time.NewTicker(10 * time.Second)
	defer tick.Stop()

	for {
		select {
		case <-ctx.Done():
			// If the context is cancelled, log the appropriate message
			return fmt.Errorf("context cancelled while waiting for %s (%s) in %s namespace", resourceName, resourceType, namespace)
		case <-timeout:
			// Log timeout error with correct parameters
			return fmt.Errorf("timeout while waiting for %s (%s) in %s namespace", resourceName, resourceType, namespace)
		case <-tick.C:
			// Check if the resource is ready
			ready, err := checkStatusFunc(resourceName, namespace)
			if err != nil {
				// Log failure in checking status
				fmt.Printf("failed to get status for %s (%s) in %s namespace: %v\n", resourceName, resourceType, namespace, err)
				continue // Continue checking at next tick
			}
			if ready {
				// Log that the resource is ready
				fmt.Printf("%s (%s) in %s namespace is ready.\n", resourceName, resourceType, namespace)
				return nil
			}
		}
	}
}

// GetResourceStatus checks the status of any Kubernetes resource.
//
// **Parameters:**
//
// ctx: A context.Context to control the operation.
// client: The dynamic.Interface client used for Kubernetes API calls.
// resourceName: The name of the resource being checked.
// namespace: The namespace of the resource.
// gvr: The schema.GroupVersionResource that specifies the resource type.
//
// **Returns:**
//
// bool: true if the resource status is 'Running', false otherwise.
// error: An error if the resource cannot be retrieved or the status is not found.
func GetResourceStatus(ctx context.Context, client dynamic.Interface, resourceName, namespace string, gvr schema.GroupVersionResource) (bool, error) {
	resource, err := client.Resource(gvr).Namespace(namespace).Get(ctx, resourceName, metav1.GetOptions{})
	if err != nil {
		return false, fmt.Errorf("failed to get %s (%s) in %s namespace: %v", resourceName, gvr.Resource, namespace, err)
	}

	status, found, err := unstructured.NestedString(resource.UnstructuredContent(), "status", "phase")
	if err != nil || !found {
		return false, fmt.Errorf("status not found for %s (%s) in %s namespace: %v", resourceName, gvr.Resource, namespace, err)
	}

	return status == "Running", nil
}

func DescribeKubernetesResource(ctx context.Context, client dynamic.Interface, resourceName, namespace string, gvr schema.GroupVersionResource) (string, error) {
	resource, err := client.Resource(gvr).Namespace(namespace).Get(ctx, resourceName, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to get %s '%s' in namespace '%s': %v", gvr.Resource, resourceName, namespace, err)
	}
	if resource == nil {
		return "", fmt.Errorf("no %s '%s' found in namespace '%s'", gvr.Resource, resourceName, namespace)
	}

	// Make sure the resource is not nil before accessing UnstructuredContent
	if resource.Object == nil {
		return "", fmt.Errorf("the resource data is nil")
	}

	description := formatResourceDescription(resource)
	return description, nil
}

// formatResourceDescription creates a detailed string representation of a Kubernetes resource similar to `kubectl describe`.
func formatResourceDescription(resource *unstructured.Unstructured) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Name: %s\n", resource.GetName()))
	sb.WriteString(fmt.Sprintf("Namespace: %s\n", resource.GetNamespace()))
	sb.WriteString(fmt.Sprintf("Labels: %v\n", resource.GetLabels()))
	sb.WriteString(fmt.Sprintf("Annotations: %v\n", resource.GetAnnotations()))
	sb.WriteString("Details:\n")

	// Include additional details like status, spec, etc.
	for key, val := range resource.UnstructuredContent() {
		sb.WriteString(fmt.Sprintf("%s: %v\n", key, val))
	}

	return sb.String()
}
