package k8s_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/l50/goutils/v2/k8s"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic/fake"
	k8stesting "k8s.io/client-go/testing"
)

func TestApplyOrDeleteManifest(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*k8s.ManifestConfig, *fake.FakeDynamicClient)
		wantErr bool
	}{
		{
			name: "apply raw manifest successfully",
			setup: func(mc *k8s.ManifestConfig, fdc *fake.FakeDynamicClient) {
				mc.Type = k8s.ManifestRaw
				mc.Operation = k8s.OperationApply
				mc.ReadFile = func(string) ([]byte, error) {
					return []byte("kind: Pod\napiVersion: v1"), nil
				}
				fdc.PrependReactor("create", "*", func(action k8stesting.Action) (bool, runtime.Object, error) {
					return true, nil, nil
				})
			},
			wantErr: false,
		},
		{
			name: "error on unsupported manifest type",
			setup: func(mc *k8s.ManifestConfig, fdc *fake.FakeDynamicClient) {
				mc.Type = k8s.ManifestType(999) // Invalid type
			},
			wantErr: true,
		},
		{
			name: "apply job manifest successfully",
			setup: func(mc *k8s.ManifestConfig, fdc *fake.FakeDynamicClient) {
				mc.Type = k8s.ManifestJob
				mc.Operation = k8s.OperationApply
				mc.ReadFile = func(string) ([]byte, error) {
					return []byte("kind: Job\napiVersion: batch/v1"), nil
				}
				fdc.PrependReactor("create", "*", func(action k8stesting.Action) (bool, runtime.Object, error) {
					return true, &unstructured.Unstructured{}, nil
				})
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mc := k8s.NewManifestConfig()
			fdc := fake.NewSimpleDynamicClient(runtime.NewScheme())
			if tc.setup != nil {
				tc.setup(mc, fdc)
			}
			mc.Client = fdc
			err := mc.ApplyOrDeleteManifest(context.Background())
			if (err != nil) != tc.wantErr {
				t.Errorf("ApplyOrDeleteManifest() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestHandleRawManifest(t *testing.T) {
	tests := []struct {
		name       string
		configFunc func() *k8s.ManifestConfig
		wantErr    bool
		setup      func(fdc *fake.FakeDynamicClient)
	}{
		{
			name: "successful apply of raw manifest",
			configFunc: func() *k8s.ManifestConfig {
				mc := &k8s.ManifestConfig{
					Operation:    k8s.OperationApply,
					ManifestPath: "path/to/valid/manifest.yaml",
					ReadFile: func(path string) ([]byte, error) {
						return []byte("kind: Deployment\napiVersion: apps/v1"), nil
					},
				}
				return mc
			},
			setup: func(fdc *fake.FakeDynamicClient) {
				fdc.PrependReactor("create", "deployments", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
					return true, nil, nil
				})
			},
			wantErr: false,
		},
		{
			name: "fail to read manifest file",
			configFunc: func() *k8s.ManifestConfig {
				mc := &k8s.ManifestConfig{
					ManifestPath: "invalid/path",
					ReadFile: func(path string) ([]byte, error) {
						return nil, fmt.Errorf("file not found")
					},
				}
				return mc
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mc := tc.configFunc()
			fakeClient := fake.NewSimpleDynamicClient(runtime.NewScheme())
			if tc.setup != nil {
				tc.setup(fakeClient)
			}
			mc.Client = fakeClient // Assigning fake dynamic client

			err := mc.HandleRawManifest(context.Background(), fakeClient)
			if (err != nil) != tc.wantErr {
				t.Errorf("HandleRawManifest() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
