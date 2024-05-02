package k8s

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	client "github.com/l50/goutils/v2/k8s/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// WaitForResourceState waits for a Kubernetes resource to reach a specified state.
//
// **Parameters:**
//
// ctx: A context.Context to allow for cancellation and timeouts.
// resourceName: The name of the resource to monitor.
// namespace: The namespace in which the resource exists.
// resourceType: The type of the resource (e.g., Pod, Service).
// desiredState: A string representing the desired state (e.g., "Running", "Deleted").
// checkStatusFunc: A function that checks if the resource is in the desired state.
//
// **Returns:**
//
// error: An error if the waiting is cancelled by context, times out, or
// fails to determine the state.
func WaitForResourceState(ctx context.Context, resourceName, namespace, resourceType, desiredState string, checkStatusFunc func(name, namespace string) (bool, error)) error {
	// Set a timeout for reaching the desired state
	timeout := time.After(5 * time.Minute)

	// Check status every second
	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()

	for {
		select {
		case <-ctx.Done():
			// If the context is cancelled, log the appropriate message
			return fmt.Errorf("context cancelled while waiting for %s (%s) in %s namespace to reach %s state", resourceName, resourceType, namespace, desiredState)
		case <-timeout:
			// Log timeout error with correct parameters
			return fmt.Errorf("timeout while waiting for %s (%s) in %s namespace to reach %s state", resourceName, resourceType, namespace, desiredState)
		case <-tick.C:
			// Check if the resource is in the desired state
			inDesiredState, err := checkStatusFunc(resourceName, namespace)
			if err != nil {
				// Log failure in checking status
				fmt.Printf("failed to get status for %s (%s) in %s namespace: %v\n", resourceName, resourceType, namespace, err)
				continue // Continue checking at next tick
			}
			if inDesiredState {
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
// kc: The KubernetesClient that includes both the standard and dynamic clients.
// resourceName: The name of the resource being checked.
// namespace: The namespace of the resource.
// gvr: The schema.GroupVersionResource that specifies the resource type.
//
// **Returns:**
//
// bool: true if the resource status is 'Running', false otherwise.
// error: An error if the resource cannot be retrieved or the status is not found.
func GetResourceStatus(ctx context.Context, kc *client.KubernetesClient, resourceName, namespace string, gvr schema.GroupVersionResource) (bool, error) {
	resource, err := kc.DynamicClient.Resource(gvr).Namespace(namespace).Get(ctx, resourceName, metav1.GetOptions{})
	if err != nil {
		return false, fmt.Errorf("failed to get %s (%s) in %s namespace: %v", resourceName, gvr.Resource, namespace, err)
	}

	status, found, err := unstructured.NestedString(resource.UnstructuredContent(), "status", "phase")
	if err != nil || !found {
		return false, fmt.Errorf("status not found for %s (%s) in %s namespace: %v", resourceName, gvr.Resource, namespace, err)
	}

	return status == "Running", nil
}

// DescribeKubernetesResource retrieves the details of a specific Kubernetes
// resource using the provided dynamic client, resource name, namespace, and
// GroupVersionResource (GVR).
//
// **Parameters:**
//
// ctx: The context to use for the request.
// kc: The KubernetesClient that includes both the standard and dynamic clients.
// resourceName: The name of the resource to describe.
// namespace: The namespace of the resource.
// gvr: The GroupVersionResource of the resource.
//
// **Returns:**
//
// string: A string representation of the resource, similar to `kubectl describe`.
// error: An error if any issue occurs while trying to describe the resource.
func DescribeKubernetesResource(ctx context.Context, kc *client.KubernetesClient, resourceName, namespace string, gvr schema.GroupVersionResource) (string, error) {
	resource, err := kc.DynamicClient.Resource(gvr).Namespace(namespace).Get(ctx, resourceName, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to get %s '%s' in namespace '%s': %v", gvr.Resource, resourceName, namespace, err)
	}

	// Make sure the resource is not nil before accessing UnstructuredContent
	if resource == nil {
		return "", fmt.Errorf("no %s '%s' found in namespace '%s'", gvr.Resource, resourceName, namespace)
	}

	return formatResourceDescription(resource), nil
}

// formatResourceDescription creates a detailed string representation of a
// Kubernetes resource similar to `kubectl describe`.
func formatResourceDescription(resource *unstructured.Unstructured) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Name: %s\n", resource.GetName()))
	sb.WriteString(fmt.Sprintf("Namespace: %s\n", resource.GetNamespace()))
	sb.WriteString(fmt.Sprintf("Labels: %v\n", resource.GetLabels()))
	sb.WriteString(fmt.Sprintf("Annotations: %v\n", resource.GetAnnotations()))
	sb.WriteString("Details:\n")

	// Sort the keys to ensure consistent order in tests and descriptions
	keys := make([]string, 0, len(resource.UnstructuredContent()))
	for key := range resource.UnstructuredContent() {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Include sorted details like status, spec, etc.
	for _, key := range keys {
		value := resource.UnstructuredContent()[key]
		sb.WriteString(fmt.Sprintf("%s: %v\n", key, value))
	}

	return sb.String()
}
