package k8s

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DeleteKubernetesJob deletes a Kubernetes Job in the specified namespace.
func (kc *KubernetesClient) DeleteKubernetesJob(jobName, namespace string) error {
	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}
	err := kc.Clientset.BatchV1().Jobs(namespace).Delete(
		context.TODO(), jobName, deleteOptions)
	if err != nil {
		return fmt.Errorf("failed to delete job '%s' in namespace '%s': %v",
			jobName, namespace, err)
	}
	return nil
}
