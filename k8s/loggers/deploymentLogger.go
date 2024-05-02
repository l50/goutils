package k8s

import (
	"context"
	"fmt"

	k8s "github.com/l50/goutils/v2/k8s/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DeploymentLogger represents a logger specifically designed for logging
// Kubernetes deployments.
//
// **Attributes:**
//
// kc: Pointer to KubernetesClient used for API requests.
// namespace: Namespace where the deployment is located.
// deploymentName: Name of the deployment to log.
type DeploymentLogger struct {
	kc             *k8s.KubernetesClient
	namespace      string
	deploymentName string
}

// NewDeploymentLogger creates a new instance of DeploymentLogger.
//
// **Parameters:**
//
// kc: Pointer to KubernetesClient.
// namespace: Namespace where the deployment is located.
// deploymentName: Name of the deployment.
//
// **Returns:**
//
// *DeploymentLogger: A new instance of DeploymentLogger.
func NewDeploymentLogger(kc *k8s.KubernetesClient, namespace, deploymentName string) *DeploymentLogger {
	return &DeploymentLogger{kc, namespace, deploymentName}
}

// FetchAndLog fetches the deployment details and logs related pod events.
//
// **Parameters:**
//
// ctx: Context to control the request lifetime.
//
// **Returns:**
//
// error: An error if any occurs during fetching and logging.
func (d *DeploymentLogger) FetchAndLog(ctx context.Context) error {
	deployment, err := d.kc.Clientset.AppsV1().Deployments(d.namespace).Get(
		ctx, d.deploymentName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("error fetching deployment: %v", err)
	}

	var selectorString string
	if deployment.Spec.Selector != nil && len(deployment.Spec.Selector.MatchLabels) > 0 {
		// Create a LabelSelector from the map[string]string
		labelSelector := &metav1.LabelSelector{MatchLabels: deployment.Spec.Selector.MatchLabels}

		// Format the LabelSelector to string form used in Kubernetes API calls
		selectorString = metav1.FormatLabelSelector(labelSelector)
	}

	return FetchAndLogPods(ctx, d.kc.Clientset, d.namespace, selectorString)
}
