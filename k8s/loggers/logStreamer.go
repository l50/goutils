package k8s

import (
	"context"
	"fmt"
	"io"
	"os"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// StreamLogs streams logs for a specific resource within a namespace.
//
// **Parameters:**
//
// clientset: Kubernetes clientset to interact with Kubernetes API.
// namespace: Namespace where the resource is located.
// resourceType: Type of resource ('pod', 'job', or 'deployment').
// resourceName: Name of the resource to stream logs from.
//
// **Returns:**
//
// error: An error if any occurs during the log streaming process.
func StreamLogs(clientset kubernetes.Interface, namespace, resourceType, resourceName string) error {
	podName := ""
	switch resourceType {
	case "pod":
		podName = resourceName
	case "job", "deployment":

		pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: fmt.Sprintf("job-name=%s", resourceName),
		})
		if err != nil {
			return fmt.Errorf("failed to list pods: %v", err)
		}
		if len(pods.Items) == 0 {
			return fmt.Errorf("no pods found for %s: %s", resourceType, resourceName)
		}

		podName = pods.Items[0].Name
	default:
		return fmt.Errorf("unsupported resource type: %s", resourceType)
	}

	podLogOpts := &corev1.PodLogOptions{
		Follow: true,
	}
	req := clientset.CoreV1().Pods(namespace).GetLogs(podName, podLogOpts)
	logStream, err := req.Stream(context.TODO())
	if err != nil {
		return fmt.Errorf("error in opening stream: %v", err)
	}
	defer logStream.Close()

	_, err = io.Copy(os.Stdout, logStream)
	if err != nil && err != io.EOF {
		return fmt.Errorf("error in copying information from log to stdout: %v", err)
	}
	return nil
}
