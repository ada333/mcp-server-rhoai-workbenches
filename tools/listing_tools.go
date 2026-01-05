package tools

import (
	"context"
	"fmt"
	"strings"

	core "github.com/ada333/MCP-test/core"
	"github.com/ada333/MCP-test/resources"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

// variables used for mocking in tests
var GetClientSet = func() (kubernetes.Interface, error) { return core.LogIntoClusterClientSet() }

var GetDynamicClient = func() (dynamic.Interface, error) { return core.LogIntoClusterDynamic() }

func ListPods(ctx context.Context, req *mcp.CallToolRequest, input core.ListWorkbenchesInput) (*mcp.CallToolResult, core.PodsOutput, error) {
	clientset, err := GetClientSet()
	if err != nil {
		return nil, core.PodsOutput{}, err
	}

	// list pods - this should be only code in the func
	pods, err := clientset.CoreV1().Pods(input.Namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, core.PodsOutput{}, fmt.Errorf("failed to list pods: %v", err)
	}

	msg := ""
	for _, pod := range pods.Items {
		msg += fmt.Sprintf("- %s (%s)\n", pod.Name, pod.Status.Phase)
	}
	return nil, core.PodsOutput{Pods: msg}, nil
}

func ListWorkbenches(ctx context.Context, req *mcp.CallToolRequest, input core.ListWorkbenchesInput) (*mcp.CallToolResult, core.ListWorkbenchesResult, error) {

	dyn, err := GetDynamicClient()
	if err != nil {
		return nil, core.ListWorkbenchesResult{}, err
	}

	notebooks, err := dyn.Resource(core.WorkbenchesGVR).Namespace(input.Namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, core.ListWorkbenchesResult{}, fmt.Errorf("failed to list workbenches: %v", err)
	}

	msg := ""
	for _, nb := range notebooks.Items {
		name := nb.GetName()
		msg += fmt.Sprintf("- %s\n", name)
	}
	return nil, core.ListWorkbenchesResult{Workbenches: msg}, nil
}

func ListAllWorkbenches(ctx context.Context, req *mcp.CallToolRequest, input core.ListWorkbenchesInput) (*mcp.CallToolResult, core.ListWorkbenchesResult, error) {
	_, workbenches, err := ListWorkbenches(ctx, req, core.ListWorkbenchesInput{Namespace: ""})
	if err != nil {
		return nil, core.ListWorkbenchesResult{}, err
	}
	return nil, core.ListWorkbenchesResult{Workbenches: workbenches.Workbenches}, nil
}

// Lists image-display-name for every image in the cluster
func ListImages(ctx context.Context, req *mcp.CallToolRequest, input core.ListWorkbenchesInput) (*mcp.CallToolResult, core.ListImagesOutput, error) {
	images, err := resources.GetImages(ctx)
	if err != nil {
		return nil, core.ListImagesOutput{}, err
	}

	msg := ""
	for _, image := range images {
		msg += fmt.Sprintf("Image: %s\n URL: %s\n Versions: %s\n", image.Annotations["opendatahub.io/notebook-image-name"], image.URL, strings.Join(image.Versions, "\n"))
	}
	return nil, core.ListImagesOutput{Images: msg}, nil
}
