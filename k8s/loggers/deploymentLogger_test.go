package k8s_test

import (
	"context"
	"fmt"
	"testing"

	k8sclient "github.com/l50/goutils/v2/k8s/client"
	k8s "github.com/l50/goutils/v2/k8s/loggers"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

func mockFileReader(filename string) ([]byte, error) {
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

func TestDeploymentLogger(t *testing.T) {
	tests := []struct {
		name           string
		namespace      string
		deploymentName string
		setupClient    func(*fake.Clientset)
		expectedError  string
	}{
		{
			name:           "successful deployment fetch and log",
			namespace:      "default",
			deploymentName: "test-deployment",
			setupClient: func(cs *fake.Clientset) {
				deployment := &appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{Name: "test-deployment", Namespace: "default"},
					Spec: appsv1.DeploymentSpec{
						Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": "test"}},
					},
				}
				cs.AddReactor("get", "deployments", func(action k8stesting.Action) (bool, runtime.Object, error) {
					getAction := action.(k8stesting.GetAction)
					if getAction.GetName() == deployment.Name && getAction.GetNamespace() == deployment.Namespace {
						return true, deployment, nil
					}
					return true, nil, fmt.Errorf("deployment not found")
				})
				_, err := cs.AppsV1().Deployments(deployment.Namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
				if err != nil {
					t.Fatal("Failed to create mock deployment:", err)
				}
			},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			clientset := fake.NewSimpleClientset()
			tc.setupClient(clientset)

			kc, err := k8sclient.NewKubernetesClient("apiTokenOrConfigString", mockFileReader)
			assert.NoError(t, err)
			kc.Clientset = clientset

			deploymentLogger := k8s.NewDeploymentLogger(kc, tc.namespace, tc.deploymentName)
			err = deploymentLogger.FetchAndLog(context.Background())
			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
			}
		})
	}
}
