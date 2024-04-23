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

// StreamLogs connects to a Kubernetes cluster and streams logs from a specified pod,
// or dynamically locates and streams logs from pods associated with a job or deployment.
//
// **Parameters:**
//
// clientset: The Kubernetes client interface for connecting to the cluster.
// namespace: The namespace in which the resources are located.
// resourceType: The type of resource ('pod', 'job', or 'deployment') from which logs are to be streamed.
// resourceName: The name of the resource.
//
// **Returns:**
//
// error: An error object if an issue occurs during the log streaming process. Nil if the operation is successful.
//
// The function first determines the pod name directly if the resource type is 'pod'. For 'job' or 'deployment',
// it queries associated pods based on label selectors. Once the relevant pod is identified, it sets up a log
// streaming connection using the Kubernetes API. Logs are streamed directly to the standard output.
// Any issues during these steps, such as failure to find pods or streaming errors, result in returning an error.
func StreamLogs(clientset kubernetes.Interface, namespace, resourceType, resourceName string) error {
	podName := ""
	switch resourceType {
	case "pod":
		podName = resourceName
	case "job", "deployment":
		// Locate the associated pods for jobs or deployments by querying based on a label selector.
		pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: fmt.Sprintf("job-name=%s", resourceName), // This label selector might need to be adjusted based on actual usage.
		})
		if err != nil {
			return fmt.Errorf("failed to list pods: %v", err)
		}
		if len(pods.Items) == 0 {
			return fmt.Errorf("no pods found for %s: %s", resourceType, resourceName)
		}
		// Assume the first pod is the desired target for log streaming.
		podName = pods.Items[0].Name
	default:
		return fmt.Errorf("unsupported resource type: %s", resourceType)
	}

	// Set up log streaming options and initiate the stream.
	podLogOpts := &corev1.PodLogOptions{
		Follow: true, // Ensures the log output is streamed until the connection is closed.
	}
	req := clientset.CoreV1().Pods(namespace).GetLogs(podName, podLogOpts)
	logStream, err := req.Stream(context.TODO())
	if err != nil {
		return fmt.Errorf("error in opening stream: %v", err)
	}
	defer logStream.Close()

	// Stream logs to standard output and handle potential errors.
	_, err = io.Copy(os.Stdout, logStream)
	if err != nil && err != io.EOF {
		return fmt.Errorf("error in copying information from log to stdout: %v", err)
	}

	return nil
}
