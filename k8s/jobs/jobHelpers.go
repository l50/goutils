package k8s

import (
	"context"
	"fmt"

	k8s "github.com/l50/goutils/v2/k8s/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// JobsClient represents a client for managing Kubernetes jobs
// through the Kubernetes API.
//
// **Attributes:**
//
// Client: A pointer to KubernetesClient for accessing Kubernetes API.
type JobsClient struct {
	Client *k8s.KubernetesClient
}

// DeleteKubernetesJob deletes a specified Kubernetes job within
// a given namespace. It sets the deletion propagation policy
// to 'Foreground' to ensure that the delete operation waits
// until the cascading delete has completed.
//
// **Parameters:**
//
// ctx: Context for managing control flow of the request.
// jobName: Name of the Kubernetes job to delete.
// namespace: Namespace where the job is located.
//
// **Returns:**
//
// error: An error if the job could not be deleted.
func (kc *JobsClient) DeleteKubernetesJob(ctx context.Context, jobName, namespace string) error {
	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}
	err := kc.Client.Clientset.BatchV1().Jobs(namespace).Delete(
		ctx, jobName, deleteOptions)
	if err != nil {
		return fmt.Errorf("failed to delete job '%s' in namespace '%s': %v",
			jobName, namespace, err)
	}
	return nil
}
