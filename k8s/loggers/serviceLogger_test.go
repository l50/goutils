package k8s_test

import (
	"context"
	"fmt"
	"testing"

	k8sclient "github.com/l50/goutils/v2/k8s/client"
	k8s "github.com/l50/goutils/v2/k8s/loggers"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

func mockFileReaderSvc(filename string) ([]byte, error) {
	return []byte(`{
		"apiVersion": "v1",
		"kind": "Config",
		"clusters": [{
			"cluster": {
				"server": "https://fake-server"
			},
			"name": "fake-cluster"
		}],
		"contexts": [{
			"context": {
				"cluster": "fake-cluster",
				"user": "fake-user"
			},
			"name": "fake"
		}],
		"current-context": "fake",
		"users": [{
			"name": "fake-user",
			"user": {}
		}]
	}`), nil
}

func TestServiceLogger(t *testing.T) {
	tests := []struct {
		name          string
		namespace     string
		serviceName   string
		setupClient   func(*fake.Clientset)
		expectedError string
	}{
		{
			name:        "successful service fetch and log",
			namespace:   "default",
			serviceName: "test-service",
			setupClient: func(cs *fake.Clientset) {
				service := &corev1.Service{
					ObjectMeta: metav1.ObjectMeta{Name: "test-service", Namespace: "default"},
					Spec: corev1.ServiceSpec{
						Selector: map[string]string{"app": "test"},
					},
				}
				cs.AddReactor("get", "services", func(action k8stesting.Action) (bool, runtime.Object, error) {
					getAction := action.(k8stesting.GetAction)
					if getAction.GetName() == service.Name && getAction.GetNamespace() == service.Namespace {
						return true, service, nil
					}
					return true, nil, fmt.Errorf("service not found")
				})
				_, err := cs.CoreV1().Services(service.Namespace).Create(context.TODO(), service, metav1.CreateOptions{})
				if err != nil {
					t.Fatal("Failed to create mock service:", err)
				}
			},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			clientset := fake.NewSimpleClientset()
			tc.setupClient(clientset)

			kc, err := k8sclient.NewKubernetesClient("apiTokenOrConfigString", mockFileReaderSvc)
			assert.NoError(t, err)
			kc.Clientset = clientset

			serviceLogger := k8s.NewServiceLogger(kc, tc.namespace, tc.serviceName)
			err = serviceLogger.FetchAndLog(context.Background())
			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
			}
		})
	}
}
