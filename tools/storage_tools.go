package tools

import (
	"context"
	"fmt"

	"github.com/amaly/mcp-server-rhoai/core"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
)

func ListPVCs(ctx context.Context, req *mcp.CallToolRequest, input core.ListPVCsInput) (*mcp.CallToolResult, core.PVCsOutput, error) {
	dyn, err := GetDynamicClient()
	if err != nil {
		return nil, core.PVCsOutput{}, err
	}

	pvcs, err := dyn.Resource(core.PvcGVR).Namespace(input.Namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, core.PVCsOutput{}, fmt.Errorf("failed to list PVCs: %v", err)
	}

	msg := ""
	for _, pvc := range pvcs.Items {
		msg += fmt.Sprintf("- %s\n", pvc.GetName())
	}
	return nil, core.PVCsOutput{PVCs: msg}, nil
}

func CreatePVCTool(ctx context.Context, req *mcp.CallToolRequest, input core.CreatePVCInput) (*mcp.CallToolResult, core.DefaultToolOutput, error) {
	dyn, err := GetDynamicClient()
	if err != nil {
		return nil, core.DefaultToolOutput{}, err
	}

	err = createPersistentVolumeClaim(ctx, dyn, input.Namespace, input.PVCName, input.Size)
	if err != nil {
		return nil, core.DefaultToolOutput{}, fmt.Errorf("failed to create PVC: %v", err)
	}

	return nil, core.DefaultToolOutput{Message: "PVC was succesfully created!"}, nil
}

func createPersistentVolumeClaim(ctx context.Context, dyn dynamic.Interface, namespace, name, size string) error {
	pvc := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "PersistentVolumeClaim",
			"metadata": map[string]interface{}{
				"name":      name,
				"namespace": namespace,
				"labels": map[string]interface{}{
					"opendatahub.io/dashboard": "true",
				},
			},
			"spec": map[string]interface{}{
				"accessModes": []interface{}{"ReadWriteOnce"},
				"resources": map[string]interface{}{
					"requests": map[string]interface{}{
						"storage": size,
					},
				},
			},
		},
	}

	_, err := dyn.Resource(core.PvcGVR).Namespace(namespace).Create(ctx, pvc, metav1.CreateOptions{})
	if err != nil && !errors.IsAlreadyExists(err) {
		return err
	}
	return nil
}
