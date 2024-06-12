package k8s

import (
	"context"
	"fmt"
	"time"

	client "github.com/l50/goutils/v2/k8s/client"
	dynK8s "github.com/l50/goutils/v2/k8s/dynamic"
	manifests "github.com/l50/goutils/v2/k8s/manifests"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
)

// DynK8sInterface defines the methods used from dynK8s
// to manage Kubernetes resources.
//
// **Methods:**
//
// WaitForResourceState: Waits for a Kubernetes resource to reach a desired state.
// GetResourceStatus: Retrieves the status of a Kubernetes resource.
type DynK8sInterface interface {
	WaitForResourceState(ctx context.Context, resourceName, namespace, resourceType, desiredState string, checkStatusFunc func(name, namespace string) (bool, error)) error
	GetResourceStatus(ctx context.Context, kc *client.KubernetesClient, resourceName, namespace string, gvr schema.GroupVersionResource) (bool, error)
}

// K8sLoggerInterface defines the methods used from k8sLogger
// to stream logs from Kubernetes pods.
//
// **Methods:**
//
// StreamLogs: Streams logs from a Kubernetes pod.
type K8sLoggerInterface interface {
	StreamLogs(clientset kubernetes.Interface, namespace, resourceType, podName string) error
}

// JobPodNameGetter defines the method to get job pod name
// by job name and namespace.
//
// **Methods:**
//
// GetJobPodName: Retrieves the name of the first pod associated with a specific
type JobPodNameGetter interface {
	GetJobPodName(ctx context.Context, jobName, namespace string) (string, error)
}

// DefaultJobPodNameGetter implements the default behavior for getting job pod
// names by using the JobsClient to fetch the pod name.
//
// **Attributes:**
//
// JC: A JobsClient for managing Kubernetes jobs.
type DefaultJobPodNameGetter struct {
	JC *JobsClient
}

// GetJobPodName retrieves the name of the first pod associated with a specific
// Kubernetes job within a given namespace. It uses a label selector to find
// pods that are labeled with the job's name. This method is typically used in
// scenarios where jobs create a single pod or when only the first pod
// is of interest.
//
// **Parameters:**
//
// ctx: Context for managing control flow of the request.
// jobName: Name of the Kubernetes job to find pods for.
// namespace: Namespace where the job and its pods are located.
//
// **Returns:**
//
// string: The name of the first pod found that is associated with the job
// error: An error if no pods are found or if an error occurs during the pod retrieval
func (d *DefaultJobPodNameGetter) GetJobPodName(ctx context.Context, jobName, namespace string) (string, error) {
	return d.JC.GetJobPodName(ctx, jobName, namespace)
}

// JobsClient represents a client for managing Kubernetes jobs
// through the Kubernetes API.
//
// **Attributes:**
//
// Client: A pointer to KubernetesClient for accessing Kubernetes API.
// DynK8s: A DynK8sInterface for managing Kubernetes resources.
// K8sLogger: A K8sLoggerInterface for streaming logs from Kubernetes pods.
// StreamLogsFn: A function for streaming logs from a Kubernetes pod.
// PodNameGetter: A JobPodNameGetter for getting job pod names.
type JobsClient struct {
	Client        *client.KubernetesClient
	DynK8s        DynK8sInterface
	K8sLogger     K8sLoggerInterface
	StreamLogsFn  func(clientset *kubernetes.Clientset, namespace, resourceType, resourceName string) error
	PodNameGetter JobPodNameGetter
}

// ApplyKubernetesJob applies a Kubernetes job manifest to a Kubernetes cluster
// using the provided kubeconfig file. The job is applied to the specified namespace.
//
// **Parameters:**
//
// jobFilePath: Path to the job manifest file to apply.
// namespace: Namespace where the job should be applied.
//
// **Returns:**
//
// error: An error if the job could not be applied.
func (jc *JobsClient) ApplyKubernetesJob(jobFilePath, namespace string, readFile func(string) ([]byte, error)) error {
	if jc.Client == nil {
		return fmt.Errorf("jobs client is not initialized")
	}
	if jc.Client.DynamicClient == nil {
		return fmt.Errorf("dynamic client is not initialized")
	}

	manifestConfig := manifests.NewManifestConfig()
	manifestConfig.KubeConfigPath = jc.Client.Config.Host
	manifestConfig.ManifestPath = jobFilePath
	manifestConfig.Namespace = namespace
	manifestConfig.Type = manifests.ManifestRaw
	manifestConfig.Operation = manifests.OperationApply
	manifestConfig.Client = jc.Client.DynamicClient
	manifestConfig.ReadFile = readFile

	if err := manifestConfig.ApplyOrDeleteManifest(context.Background()); err != nil {
		return fmt.Errorf("failed to apply job: %v", err)
	}

	return nil
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
	if jc.Client == nil {
		return fmt.Errorf("jobs client is not initialized")
	}
	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}
	err := jc.Client.Clientset.BatchV1().Jobs(namespace).Delete(ctx, jobName, deleteOptions)
	if err != nil {
		return fmt.Errorf("failed to delete job '%s' in namespace '%s': %v", jobName, namespace, err)
	}
	return nil
}

// GetJobPodName retrieves the name of the first pod associated with a specific
// Kubernetes job within a given namespace. It uses a label selector to find
// pods that are labeled with the job's name. This method is typically used in
// scenarios where jobs create a single pod or when only the first pod
// is of interest.
//
// **Parameters:**
//
// ctx: Context for managing control flow of the request.
// jobName: Name of the Kubernetes job to find pods for.
// namespace: Namespace where the job and its pods are located.
//
// **Returns:**
//
// string: The name of the first pod found that is associated with the job
// error: An error if no pods are found or if an error occurs during the pod retrieval
func (jc *JobsClient) GetJobPodName(ctx context.Context, jobName, namespace string) (string, error) {
	if jc.Client == nil {
		return "", fmt.Errorf("jobs client is not initialized")
	}
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

	return pods.Items[0].Name, nil
}

// ListKubernetesJobs lists Kubernetes jobs from a specified namespace, or all namespaces
// if no namespace is specified. This method allows for either targeted or broad job retrieval.
//
// **Parameters:**
//
// ctx: Context for managing control flow of the request.
// namespace: Optional; specifies the namespace from which to list jobs. If empty, jobs will be listed from all namespaces.
//
// **Returns:**
//
// []batchv1.Job: A slice of batchv1.Job objects containing the jobs found.
// error: An error if the API call to fetch the jobs fails.
func (jc *JobsClient) ListKubernetesJobs(ctx context.Context, namespace string) ([]batchv1.Job, error) {
	if jc.Client == nil {
		return nil, fmt.Errorf("jobs client is not initialized")
	}
	listOptions := metav1.ListOptions{}
	var jobs *batchv1.JobList
	var err error

	if namespace == "" {
		jobs, err = jc.Client.Clientset.BatchV1().Jobs("").List(ctx, listOptions)
	} else {
		jobs, err = jc.Client.Clientset.BatchV1().Jobs(namespace).List(ctx, listOptions)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get jobs: %v", err)
	}

	return jobs.Items, nil
}

// JobExists checks if a Kubernetes job with the specified name exists within a given namespace.
//
// **Parameters:**
//
// ctx: Context for managing control flow of the request.
// jobName: Name of the Kubernetes job to check for existence.
// namespace: Namespace where the job is located.
//
// **Returns:**
//
// bool: true if the job exists, false otherwise.
// error: An error if the job existence check fails.
func (jc *JobsClient) JobExists(ctx context.Context, jobName, namespace string) (bool, error) {
	if jc.Client == nil {
		return false, fmt.Errorf("jobs client is not initialized")
	}
	_, err := jc.Client.Clientset.BatchV1().Jobs(namespace).Get(ctx, jobName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil // Job does not exist
		}
		return false, fmt.Errorf("failed to get job '%s' in namespace '%s': %v", jobName, namespace, err)
	}
	return true, nil // Job exists
}

// StreamJobLogs monitors a Kubernetes job by waiting for it to reach
// the 'Ready' state and then streams logs from the associated pod.
//
// **Parameters:**
//
// jobsClient: A JobsClient for managing Kubernetes jobs.
// workloadName: Name of the Kubernetes job to monitor.
// namespace: Namespace where the job is located.
//
// **Returns:**
//
// error: An error if the job monitoring fails.
func (jc *JobsClient) StreamJobLogs(workloadName, namespace string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	fmt.Printf("Monitoring %s job in %s namespace\n", workloadName, namespace)

	// Wait for the job to reach completion
	err := jc.DynK8s.WaitForResourceState(ctx, workloadName, namespace, "job", "Complete", func(name, ns string) (bool, error) {
		jobComplete, err := jc.DynK8s.GetResourceStatus(ctx, jc.Client, name, ns, schema.GroupVersionResource{
			Group:    "batch",
			Version:  "v1",
			Resource: "jobs",
		})
		if err != nil {
			return false, fmt.Errorf("error checking status for %s job in %s namespace: %v", name, ns, err)
		}
		return jobComplete, nil
	})

	if err != nil {
		if diagErr := logJobDiagnosticInfo(jc.Client, workloadName, namespace); diagErr != nil {
			fmt.Printf("failed to log diagnostic info for %s job: %v", workloadName, diagErr)
		}
		return fmt.Errorf("error waiting for %s job to complete in %s namespace: %v", workloadName, namespace, err)
	}

	// Attempt to fetch the pod name after confirming the job has completed
	podName, err := jc.PodNameGetter.GetJobPodName(ctx, workloadName, namespace)
	if err != nil {
		return fmt.Errorf("failed to find pod associated with %s workload: %v", workloadName, err)
	}

	fmt.Printf("%s pod for %s job in %s namespace is ready and being monitored\n", podName, workloadName, namespace)

	// Stream logs from the pod, ensuring it exists
	if err := jc.K8sLogger.StreamLogs(jc.Client.Clientset, namespace, "pod", podName); err != nil {
		return fmt.Errorf("failed to stream logs for pod '%s': %v", podName, err)
	}

	return nil
}

// logJobDiagnosticInfo logs diagnostic information for a Kubernetes job and its associated pods.
//
// **Parameters:**
//
// k8sClient: A KubernetesClient for accessing the Kubernetes API.
// jobName: Name of the Kubernetes job to log diagnostic information for.
// namespace: Namespace where the job is located.
//
// **Returns:**
//
// error: An error if the diagnostic information could not be logged.
func logJobDiagnosticInfo(k8sClient *client.KubernetesClient, jobName, namespace string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	jobGVR := schema.GroupVersionResource{Group: "batch", Version: "v1", Resource: "jobs"}
	jobDescription, err := dynK8s.DescribeKubernetesResource(ctx, k8sClient, jobName, namespace, jobGVR)
	if err != nil {
		return fmt.Errorf("error describing job: %v", err)
	}

	fmt.Printf("Describe job output for '%s':\n%s", jobName, jobDescription)

	podsGVR := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}
	podsDescription, err := dynK8s.DescribeKubernetesResource(ctx, k8sClient, jobName, namespace, podsGVR)
	if err != nil {
		return fmt.Errorf("error describing pods for job: %v", err)
	}

	fmt.Printf("Describe pods output for job '%s':\n%s", jobName, podsDescription)
	return nil
}
