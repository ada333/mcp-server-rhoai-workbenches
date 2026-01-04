package tools

import (
	"context"
	"fmt"

	core "github.com/ada333/MCP-test/main_logic"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func CreateCustomImage(ctx context.Context, req *mcp.CallToolRequest, input core.CreateCustomImageInput) (*mcp.CallToolResult, core.WorkbenchOutput, error) {
	dyn, err := GetDynamicClient()
	if err != nil {
		return nil, core.WorkbenchOutput{}, err
	}

	namespace := "redhat-ods-applications"

	imageStream := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       "ImageStream",
			"apiVersion": "image.openshift.io/v1",
			"metadata": map[string]interface{}{
				"name":      input.ImageName,
				"namespace": namespace,
				"annotations": map[string]interface{}{
					"opendatahub.io/notebook-image-creator":   "htpasswd-cluster-admin-user", // who should it be?
					"opendatahub.io/notebook-image-desc":      input.ImageDescription,
					"opendatahub.io/notebook-image-name":      input.ImageName,
					"opendatahub.io/notebook-image-url":       input.ImageLocation,
					"opendatahub.io/recommended-accelerators": "[]",
				},
				"labels": map[string]interface{}{
					"app.kubernetes.io/created-by":  "byon",
					"opendatahub.io/dashboard":      "true",
					"opendatahub.io/notebook-image": "true",
				},
			},
			"spec": map[string]interface{}{
				"lookupPolicy": map[string]interface{}{
					"local": true,
				},
				"tags": []interface{}{
					map[string]interface{}{
						"name": "latest",
						"annotations": map[string]interface{}{
							"opendatahub.io/notebook-python-dependencies": "[]",
							"opendatahub.io/notebook-software":            "[]",
							"openshift.io/imported-from":                  input.ImageLocation,
						},
						"from": map[string]interface{}{
							"kind": "DockerImage",
							"name": input.ImageLocation,
						},
						"importPolicy": map[string]interface{}{
							"importMode": "Legacy",
						},
						"referencePolicy": map[string]interface{}{
							"type": "Source",
						},
					},
				},
			},
		},
	}

	_, err = dyn.Resource(core.ImagesGVR).Namespace(namespace).Create(ctx, imageStream, metav1.CreateOptions{})
	if err != nil {
		return nil, core.WorkbenchOutput{}, fmt.Errorf("failed to create notebook: %v", err)
	}

	return nil, core.WorkbenchOutput{Message: "Workbench was succesfully created!"}, nil
}
