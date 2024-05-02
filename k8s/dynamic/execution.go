package k8s

import (
	"context"
	"fmt"
	"io"
	"net/url"

	client "github.com/l50/goutils/v2/k8s/client"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/kubectl/pkg/scheme"
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

// ExecParams contains all the parameters needed to execute a command in a Kubernetes pod.
//
// **Attributes:**
//
// Context: The context to use for the request.
// Client: The KubernetesClient that includes both the standard and dynamic clients.
// Namespace: The namespace of the resource where the pod is located.
// PodName: The name of the pod to execute the command in.
// Command: A slice of strings representing the command to execute inside the pod.
// Stdin: An io.Reader to use as the standard input for the command.
// Stdout: An io.Writer to use as the standard output for the command.
// Stderr: An io.Writer to use as the standard error for the command.
type ExecParams struct {
	Context   context.Context
	Client    *client.KubernetesClient
	Namespace string
	PodName   string
	Command   []string
	Stdin     io.Reader
	Stdout    io.Writer
	Stderr    io.Writer
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
//
// **Returns:**
//
// string: The output from the executed command or an error message if execution fails.
// error: An error if any issue occurs during the setup or execution of the command.
func ExecKubernetesResources(params ExecParams) (string, error) {
	if params.Client == nil || params.Client.Clientset == nil {
		return "", fmt.Errorf("KubernetesClient or Clientset is not initialized")
	}

	// Fetch the pod to ensure it exists and is in a Running state
	pod, err := params.Client.Clientset.CoreV1().Pods(params.Namespace).Get(params.Context, params.PodName, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("error fetching pod %s in namespace %s: %v", params.PodName, params.Namespace, err)
	}

	if pod.Status.Phase != v1.PodRunning {
		return "", fmt.Errorf("pod %s is not in running state, current state: %s", params.PodName, pod.Status.Phase)
	}

	if params.Client.Clientset == nil || params.Client.Clientset.CoreV1() == nil {
		return "", fmt.Errorf("kubernetes clientset is not initialized")
	}

	req := params.Client.Clientset.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Name(params.PodName).
		Namespace(params.Namespace).
		SubResource("exec")

	options := &v1.PodExecOptions{
		Command: params.Command,
		Stdin:   params.Stdin != nil,
		Stdout:  params.Stdout != nil,
		Stderr:  params.Stderr != nil,
		TTY:     true,
	}

	req.VersionedParams(
		options,
		scheme.ParameterCodec,
	)

	// Check that the URL to which we are making the request is valid
	if req.URL() == nil {
		return "", fmt.Errorf("failed to form a valid request URL")
	}

	executor, err := remotecommand.NewSPDYExecutor(params.Client.Config, "POST", req.URL())
	if err != nil {
		return "", fmt.Errorf("failed to initialize command executor: %v", err)
	}

	streamOptions := remotecommand.StreamOptions{
		Stdin:  params.Stdin,
		Stdout: params.Stdout,
		Stderr: params.Stderr,
	}

	if err := executor.StreamWithContext(params.Context, streamOptions); err != nil {
		return "", fmt.Errorf("failed to execute command: %v", err)
	}

	return "Command executed successfully", nil
}
