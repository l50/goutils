package k8s

import (
	"context"
	"fmt"

	k8s "github.com/l50/goutils/v2/k8s/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ServiceLogger represents a logger specifically designed for logging
// Kubernetes services.
//
// **Attributes:**
//
// kc: Pointer to KubernetesClient used for API requests.
// namespace: Namespace where the service is located.
// serviceName: Name of the service to log.
type ServiceLogger struct {
	kc          *k8s.KubernetesClient
	namespace   string
	serviceName string
}

// NewServiceLogger creates a new instance of ServiceLogger.
//
// **Parameters:**
//
// kc: Pointer to KubernetesClient.
// namespace: Namespace where the service is located.
// serviceName: Name of the service.
//
// **Returns:**
//
// *ServiceLogger: A new instance of ServiceLogger.
func NewServiceLogger(kc *k8s.KubernetesClient, namespace, serviceName string) *ServiceLogger {
	return &ServiceLogger{kc, namespace, serviceName}
}

// FetchAndLog fetches the service details and logs related pod events.
//
// **Parameters:**
//
// ctx: Context to control the request lifetime.
//
// **Returns:**
//
// error: An error if any occurs during fetching and logging.
func (s *ServiceLogger) FetchAndLog(ctx context.Context) error {
	service, err := s.kc.Clientset.CoreV1().Services(s.namespace).Get(
		ctx, s.serviceName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("error fetching service: %v", err)
	}

	var selectorString string
	if len(service.Spec.Selector) > 0 {
		// Create a LabelSelector from the map[string]string
		labelSelector := &metav1.LabelSelector{MatchLabels: service.Spec.Selector}

		// Format the LabelSelector to string form used in Kubernetes API calls
		selectorString = metav1.FormatLabelSelector(labelSelector)
	}

	return FetchAndLogPods(ctx, s.kc.Clientset, s.namespace, selectorString)
}
