package tools

import (
	"context"
	"fmt"
	"strings"

	core "github.com/amaly/mcp-server-rhoai/core"
	"github.com/amaly/mcp-server-rhoai/resources"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func ListImages(ctx context.Context, req *mcp.CallToolRequest, input struct{}) (*mcp.CallToolResult, core.ListImagesOutput, error) {
	images, err := resources.GetImages(ctx)
	if err != nil {
		return nil, core.ListImagesOutput{}, err
	}

	msg := ""
	for _, image := range images {
		var versions []string
		for _, v := range image.Versions {
			vInfo := v.Name
			if v.PythonDependencies != "" {
				vInfo += fmt.Sprintf(" (Python: %s)", v.PythonDependencies)
			}
			if v.Software != "" {
				vInfo += fmt.Sprintf(" (Software: %s)", v.Software)
			}
			versions = append(versions, vInfo)
		}
		msg += fmt.Sprintf("Image: %s\n URL: %s\n Versions: %s\n", image.Annotations["opendatahub.io/notebook-image-name"], image.URL, strings.Join(versions, ", "))
	}
	return nil, core.ListImagesOutput{Images: msg}, nil
}

func CreateCustomImage(ctx context.Context, req *mcp.CallToolRequest, input core.CreateCustomImageInput) (*mcp.CallToolResult, core.DefaultToolOutput, error) {
	dyn, err := GetDynamicClient()
	if err != nil {
		return nil, core.DefaultToolOutput{}, err
	}

	namespace := core.GetDefaultNamespace()

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
		return nil, core.DefaultToolOutput{}, fmt.Errorf("failed to create image: %v", err)
	}

	return nil, core.DefaultToolOutput{Message: "Image was successfully created!"}, nil
}

func DeleteImage(ctx context.Context, req *mcp.CallToolRequest, input core.DeleteImageInput) (*mcp.CallToolResult, core.DefaultToolOutput, error) {
	dyn, err := GetDynamicClient()
	if err != nil {
		return nil, core.DefaultToolOutput{}, err
	}

	if isDefault, err := ImageIsDefault(ctx, input.ImageName); err != nil {
		return nil, core.DefaultToolOutput{}, err
	} else if isDefault {
		return nil, core.DefaultToolOutput{}, fmt.Errorf("image %s is a default image and cannot be deleted", input.ImageName)
	}

	if isUsed, err := ImageIsUsed(ctx, input.ImageName); err != nil {
		return nil, core.DefaultToolOutput{}, err
	} else if isUsed {
		return nil, core.DefaultToolOutput{}, fmt.Errorf("image %s is used by a workbench and cannot be deleted", input.ImageName)
	}

	err = dyn.Resource(core.ImagesGVR).Namespace(core.GetDefaultNamespace()).Delete(ctx, input.ImageName, metav1.DeleteOptions{})
	if err != nil {
		return nil, core.DefaultToolOutput{}, fmt.Errorf("failed to delete image %s: %v", input.ImageName, err)
	}

	return nil, core.DefaultToolOutput{Message: fmt.Sprintf("Image %s was successfully deleted", input.ImageName)}, nil
}

func ImageIsUsed(ctx context.Context, imageName string) (bool, error) {
	dyn, err := GetDynamicClient()
	if err != nil {
		return false, err
	}

	workbenches, err := dyn.Resource(core.WorkbenchesGVR).Namespace("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("failed to list workbenches: %v", err)
	}

	for _, workbench := range workbenches.Items {
		if workbench.GetAnnotations()["opendatahub.io/image-display-name"] == imageName {
			return true, nil
		}
	}

	return false, nil
}

func ImageIsDefault(ctx context.Context, imageName string) (bool, error) {
	dyn, err := GetDynamicClient()
	if err != nil {
		return false, err
	}

	image, err := dyn.Resource(core.ImagesGVR).Namespace(core.GetDefaultNamespace()).Get(ctx, imageName, metav1.GetOptions{})
	if err != nil {
		return false, err
	}

	return image.GetAnnotations()["internal.config.kubernetes.io/previousNamespaces"] == "default", nil
}
