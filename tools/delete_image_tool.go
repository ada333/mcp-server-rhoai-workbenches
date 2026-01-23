package tools

import (
	"context"
	"fmt"

	core "github.com/amaly/mcp-server-rhoai/core"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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

	namespaces, err := GetAllNamespaces(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get namespaces: %v", err)
	}

	for _, namespace := range namespaces {
		workbenches, err := dyn.Resource(core.WorkbenchesGVR).Namespace(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			continue
		}

		for _, workbench := range workbenches.Items {
			if workbench.GetAnnotations()["opendatahub.io/image-display-name"] == imageName {
				return true, nil
			}
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
