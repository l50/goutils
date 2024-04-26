package k8s

import (
	"context"
	"fmt"
	"os"

	client "github.com/l50/goutils/v2/k8s/client"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/remotecommand"
)

// ExecKubernetesResources executes a command in a specified resource within a given namespace using the existing KubernetesClient.
//
// **Parameters:**
//
// ctx: The context to use for the request.
// kc: The KubernetesClient that includes both the standard and dynamic clients.
// namespace: The namespace of the resource.
// podName: The name of the pod to execute the command in.
// command: A slice of strings representing the command to execute inside the resource.
//
// **Returns:**
//
// string: The output from the executed command or an error message.
// error: An error if any issue occurs during the command execution.
func ExecKubernetesResources(ctx context.Context, kc *client.KubernetesClient, namespace, podName string, command []string) (string, error) {
	req := kc.Clientset.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Command: command,
			Stdin:   true,
			Stdout:  true,
			Stderr:  true,
			TTY:     true,
		}, metav1.ParameterCodec)

	executor, err := remotecommand.NewSPDYExecutor(kc.Config, "POST", req.URL())
	if err != nil {
		return "", fmt.Errorf("failed to initialize command executor: %v", err)
	}

	if err := executor.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}); err != nil {
		return "", fmt.Errorf("failed to execute command: %v", err)
	}

	return "Command executed successfully", nil
}
