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
func (jc *JobsClient) DeleteKubernetesJob(ctx context.Context, jobName, namespace string) error {
	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}
	err := jc.Client.Clientset.BatchV1().Jobs(namespace).Delete(
		ctx, jobName, deleteOptions)
	if err != nil {
		return fmt.Errorf("failed to delete job '%s' in namespace '%s': %v", jobName, namespace, err)
	}
	return nil
}

// GetJobPodName retrieves the name of the first pod associated with a specific Kubernetes job
// within a given namespace. It uses a label selector to find pods that are labeled with
// the job's name. This method is typically used in scenarios where jobs create a single pod or
// when only the first pod is of interest.
//
// **Parameters:**
//
// ctx: Context for managing control flow of the request.
// jobName: Name of the Kubernetes job to find pods for.
// namespace: Namespace where the job and its pods are located.
//
// **Returns:**
//
// string: The name of the first pod found that is associated with the job.
// error: An error if no pods are found or if an error occurs during the pod retrieval.
func (jc *JobsClient) GetJobPodName(ctx context.Context, jobName, namespace string) (string, error) {
	labelSelector := "job-name=" + jobName
	pods, err := jc.Client.Clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get pods for job '%s' in namespace '%s': %v", jobName, namespace, err)
	}

	if len(pods.Items) == 0 {
		return "", fmt.Errorf("no pod found for job '%s'", jobName)
	}

	// Assuming the first pod is the one we're interested in, as jobs usually have one pod if not paralleled
	return pods.Items[0].Name, nil
}
