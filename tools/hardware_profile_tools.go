package tools

import (
	"context"
	"fmt"
	"time"

	core "github.com/amaly/mcp-server-rhoai/core"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

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

func CreateHardwareProfile(ctx context.Context, req *mcp.CallToolRequest, input core.HardwareProfile) (*mcp.CallToolResult, core.DefaultToolOutput, error) {
	dyn, err := GetDynamicClient()
	if err != nil {
		return nil, core.DefaultToolOutput{}, err
	}

	namespace := core.GetDefaultNamespace()

	identifiers := make([]interface{}, len(input.Resources))
	for i, res := range input.Resources {
		identifiers[i] = map[string]interface{}{
			"displayName":  res.ResourceName,
			"identifier":   res.ResourceIdentifier,
			"resourceType": res.ResourceType,
			"defaultCount": res.DefaultCount,
			"maxCount":     res.MaxCount,
			"minCount":     res.MinCount,
		}
	}

	annotations := map[string]interface{}{
		"opendatahub.io/dashboard-feature-visibility": "[]",
		"opendatahub.io/disabled":                     "false",
		"opendatahub.io/display-name":                 input.HardwareProfileName,
		"opendatahub.io/modified-date":                time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
	}

	hardwareProfile := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "infrastructure.opendatahub.io/v1",
			"kind":       "HardwareProfile",
			"metadata": map[string]interface{}{
				"name":        input.HardwareProfileName,
				"namespace":   namespace,
				"annotations": annotations,
			},
			"spec": map[string]interface{}{
				"identifiers": identifiers,
			},
		},
	}

	_, err = dyn.Resource(core.HardwareProfilesGVR).Namespace(namespace).Create(ctx, hardwareProfile, metav1.CreateOptions{})
	if err != nil {
		return nil, core.DefaultToolOutput{}, fmt.Errorf("failed to create hardware profile: %v", err)
	}

	return nil, core.DefaultToolOutput{
		Message: "Hardware Profile was successfully created!",
	}, nil
}

func DeleteHardwareProfile(ctx context.Context, req *mcp.CallToolRequest, input core.DeleteHardwareProfileInput) (*mcp.CallToolResult, core.DefaultToolOutput, error) {
	dyn, err := GetDynamicClient()
	if err != nil {
		return nil, core.DefaultToolOutput{}, err
	}

	err = dyn.Resource(core.HardwareProfilesGVR).Namespace(core.GetDefaultNamespace()).Delete(ctx, input.HardwareProfileName, metav1.DeleteOptions{})
	if err != nil {
		return nil, core.DefaultToolOutput{}, fmt.Errorf("failed to delete hardware profile %s: %v", input.HardwareProfileName, err)
	}

	return nil, core.DefaultToolOutput{Message: fmt.Sprintf("Hardware Profile %s was successfully deleted", input.HardwareProfileName)}, nil
}
