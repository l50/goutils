package k8s

import (
	"context"
	"fmt"
	"net/url"
	"os"

	client "github.com/l50/goutils/v2/k8s/client"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

// ExecutorCreator represents an interface that includes a method for creating
// a new SPDY executor.
//
// **Methods:**
//
// NewSPDYExecutor: Creates a new SPDY executor given a configuration, method,
// and URL.
type ExecutorCreator interface {
	// NewSPDYExecutor creates a new SPDY executor given a configuration,
	// method, and URL. It returns a remotecommand.Executor and an error.
	//
	// **Parameters:**
	//
	// config: A pointer to a rest.Config struct that includes the configuration
	// for the executor.
	// method: A string representing the HTTP method to use for the request.
	// url: A pointer to a url.URL struct that includes the URL for the request.
	//
	// **Returns:**
	//
	// remotecommand.Executor: The created SPDY executor.
	// error: An error if any issue occurs while creating the executor.
	NewSPDYExecutor(config *rest.Config, method string, url *url.URL) (remotecommand.Executor, error)
}

// DefaultExecutorCreator represents a struct that includes a method for
// creating a new SPDY executor.
type DefaultExecutorCreator struct{}

// NewSPDYExecutor creates a new SPDY executor given a configuration, method,
// and URL. It returns a remotecommand.Executor and an error.
//
// **Parameters:**
//
// config: A pointer to a rest.Config struct that includes the configuration
// for the executor.
// method: A string representing the HTTP method to use for the request.
// url: A pointer to a url.URL struct that includes the URL for the request.
//
// **Returns:**
//
// remotecommand.Executor: The created SPDY executor.
// error: An error if any issue occurs while creating the executor.
func (dec *DefaultExecutorCreator) NewSPDYExecutor(config *rest.Config, method string, url *url.URL) (remotecommand.Executor, error) {
	return remotecommand.NewSPDYExecutor(config, method, url)
}

// ExecKubernetesResources executes a command in a specified resource within a
// given namespace using the existing KubernetesClient.
//
// **Parameters:**
//
// ctx: The context to use for the request.
// kc: The KubernetesClient that includes both the standard and dynamic clients.
// namespace: The namespace of the resource where the pod is located.
// podName: The name of the pod to execute the command in.
// command: A slice of strings representing the command to execute inside the pod.
// restClient: The rest.Interface used to create the request.
// executorCreator: An ExecutorCreator interface to create the SPDYExecutor for command execution.
//
// **Returns:**
//
// string: The output from the executed command or an error message if execution fails.
// error: An error if any issue occurs during the setup or execution of the command.
func ExecKubernetesResources(ctx context.Context, kc *client.KubernetesClient, namespace, podName string, command []string, restClient rest.Interface, executorCreator ExecutorCreator) (string, error) {
	if kc == nil || kc.Clientset == nil {
		return "", fmt.Errorf("KubernetesClient or Clientset is not initialized")
	}
	req := restClient.Post().Resource("pods").Name(podName).Namespace(namespace).SubResource("exec").VersionedParams(&v1.PodExecOptions{
		Command: command,
		Stdin:   true,
		Stdout:  true,
		Stderr:  true,
		TTY:     true,
	}, metav1.ParameterCodec)
	executor, err := executorCreator.NewSPDYExecutor(kc.Config, "POST", req.URL())
	if err != nil {
		return "", fmt.Errorf("failed to initialize command executor: %v", err)
	}
	if err := executor.StreamWithContext(ctx, remotecommand.StreamOptions{Stdin: os.Stdin, Stdout: os.Stdout, Stderr: os.Stderr, Tty: true}); err != nil {
		return "", fmt.Errorf("failed to execute command: %v", err)
	}
	return "Command executed successfully", nil
}
