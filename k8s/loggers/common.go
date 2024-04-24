package k8s

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// FetchAndLogPods fetches and logs pods based on the specified label selector.
//
// **Parameters:**
//
// ctx: Context to control the request lifetime.
// clientset: Kubernetes clientset to interact with Kubernetes API.
// namespace: Namespace from which to list the pods.
// labelSelector: String defining the label selector for filtering pods.
//
// **Returns:**
//
// error: An error if any occurs during fetching and logging of pods.
func FetchAndLogPods(ctx context.Context, clientset kubernetes.Interface, namespace, labelSelector string) error {
	pods, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return fmt.Errorf("error listing pods: %v", err)
	}
	if len(pods.Items) == 0 {
		fmt.Println("no pods found.")
		return nil
	}
	for _, pod := range pods.Items {
		logs, err := clientset.CoreV1().Pods(namespace).GetLogs(pod.Name, &v1.PodLogOptions{}).DoRaw(ctx)
		if err != nil {
			fmt.Printf("error getting logs for pod %s: %v\n", pod.Name, err)
			continue
		}
		fmt.Printf("Logs for pod %s:\n%s\n", pod.Name, string(logs))
	}
	return nil
}
