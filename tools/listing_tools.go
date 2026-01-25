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

func ListNamespaces(ctx context.Context, req *mcp.CallToolRequest, input struct{}) (*mcp.CallToolResult, core.ListNamespacesOutput, error) {
	namespaces, err := GetAllNamespaces(ctx)
	if err != nil {
		return nil, core.ListNamespacesOutput{}, err
	}

	msg := ""
	for _, ns := range namespaces {
		msg += fmt.Sprintf("- %s\n", ns)
	}
	return nil, core.ListNamespacesOutput{Namespaces: msg}, nil
}

func GetAllNamespaces(ctx context.Context) ([]string, error) {
	clientset, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	namespaces, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %v", err)
	}

	var names []string
	for _, ns := range namespaces.Items {
		names = append(names, ns.Name)
	}
	return names, nil
}

func convertToString(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v
	case int:
		return fmt.Sprintf("%d", v)
	case int64:
		return fmt.Sprintf("%d", v)
	case float64:
		return fmt.Sprintf("%.0f", v)
	default:
		return ""
	}
}

func ListHardwareProfiles(ctx context.Context, req *mcp.CallToolRequest, input struct{}) (*mcp.CallToolResult, core.ListHardwareProfilesOutput, error) {
	dyn, err := GetDynamicClient()
	if err != nil {
		return nil, core.ListHardwareProfilesOutput{}, err
	}

	hardwareProfiles, err := dyn.Resource(core.HardwareProfilesGVR).Namespace(core.GetDefaultNamespace()).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, core.ListHardwareProfilesOutput{}, fmt.Errorf("failed to list hardware profiles: %v", err)
	}

	var result []core.HardwareProfile
	for _, profile := range hardwareProfiles.Items {
		identifiers, foundIdentifiers, err := unstructured.NestedSlice(profile.Object, "spec", "identifiers")
		if !foundIdentifiers || err != nil {
			continue
		}

		resources := make([]core.HardwareProfileResource, 0, len(identifiers))
		for _, identifier := range identifiers {
			identifierMap, okIdentifierMap := identifier.(map[string]interface{})
			if !okIdentifierMap {
				continue
			}

			displayName, _ := identifierMap["displayName"].(string)
			identifierStr, _ := identifierMap["identifier"].(string)
			resourceType, _ := identifierMap["resourceType"].(string)
			defaultCount := convertToString(identifierMap["defaultCount"])
			maxCount := convertToString(identifierMap["maxCount"])
			minCount := convertToString(identifierMap["minCount"])

			resources = append(resources, core.HardwareProfileResource{
				ResourceName:       displayName,
				ResourceIdentifier: identifierStr,
				ResourceType:       resourceType,
				DefaultCount:       defaultCount,
				MaxCount:           maxCount,
				MinCount:           minCount,
			})
		}

		result = append(result, core.HardwareProfile{
			HardwareProfileName: profile.GetName(),
			Resources:           resources,
		})
	}
	return nil, core.ListHardwareProfilesOutput{HardwareProfiles: result}, nil
}
