package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	core "github.com/amaly/mcp-server-rhoai/core"
	"github.com/amaly/mcp-server-rhoai/resources"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
)

func getPVCNameFromWorkbench(wb *unstructured.Unstructured) (string, error) {
	volumes, found, err := unstructured.NestedSlice(wb.Object, "spec", "template", "spec", "volumes")
	if err != nil {
		return "", fmt.Errorf("failed to get volumes: %v", err)
	}
	if !found {
		return "", nil
	}

	for _, vol := range volumes {
		volMap, ok := vol.(map[string]interface{})
		if !ok {
			continue
		}
		if pvc, found, _ := unstructured.NestedString(volMap, "persistentVolumeClaim", "claimName"); found {
			return pvc, nil
		}
	}
	return "", nil
}

func getResourceRequestsFromWorkbench(wb *unstructured.Unstructured) (cpuRequest, memoryRequest, gpuRequest string, err error) {
	containers, found, err := unstructured.NestedSlice(wb.Object, "spec", "template", "spec", "containers")
	if err != nil {
		return "", "", "", fmt.Errorf("failed to get containers: %v", err)
	}
	if !found || len(containers) == 0 {
		return "", "", "", nil
	}

	container, ok := containers[0].(map[string]interface{})
	if !ok {
		return "", "", "", nil
	}

	requests, found, err := unstructured.NestedStringMap(container, "resources", "requests")
	if err != nil {
		return "", "", "", fmt.Errorf("failed to get resource requests: %v", err)
	}
	if !found {
		return "", "", "", nil
	}

	cpuRequest = requests["cpu"]
	memoryRequest = requests["memory"]
	gpuRequest = requests["nvidia.com/gpu"] // maybe there can be other than nvidia.com/gpu?
	return cpuRequest, memoryRequest, gpuRequest, nil
}

// lists workbenches in a given namespace
func ListWorkbenches(ctx context.Context, req *mcp.CallToolRequest, input core.ListWorkbenchesInput) (*mcp.CallToolResult, core.ListWorkbenchesResult, error) {
	dyn, err := GetDynamicClient()
	if err != nil {
		return nil, core.ListWorkbenchesResult{}, err
	}

	workbenches, err := dyn.Resource(core.WorkbenchesGVR).Namespace(input.Namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, core.ListWorkbenchesResult{}, fmt.Errorf("failed to list workbenches: %v", err)
	}

	workbenchesInfo := []core.WorkbenchInfo{}
	for _, wb := range workbenches.Items {
		name := wb.GetName()
		user := wb.GetAnnotations()["opendatahub.io/username"]
		status := wb.GetAnnotations()["kubeflow-resource-stopped"]
		imageDisplayName := wb.GetAnnotations()["opendatahub.io/image-display-name"]

		imageTag := ""
		lastImageSelection := wb.GetAnnotations()["notebooks.opendatahub.io/last-image-selection"]
		if lastImageSelection != "" {
			parts := strings.Split(lastImageSelection, ":")
			if len(parts) > 1 {
				imageTag = parts[1]
			}
		}

		hardwareProfile := wb.GetAnnotations()["opendatahub.io/hardware-profile-name"]
		namespace := wb.GetNamespace()
		pvcName, err := getPVCNameFromWorkbench(&wb)
		if err != nil {
			return nil, core.ListWorkbenchesResult{}, fmt.Errorf("failed to get PVC name for workbench %s: %v", name, err)
		}

		cpuRequest, memoryRequest, gpuRequest, err := getResourceRequestsFromWorkbench(&wb)
		if err != nil {
			return nil, core.ListWorkbenchesResult{}, fmt.Errorf("failed to get resource requests for workbench %s: %v", name, err)
		}

		if status != "" {
			status = "stopped"
		} else {
			status = "running"
		}

		uptime := ""
		if status == "running" {
			uptime, err = getUptimeFromWorkbench(wb.GetName(), wb.GetNamespace())
			if err != nil {
				return nil, core.ListWorkbenchesResult{}, fmt.Errorf("failed to get uptime for workbench %s: %v", name, err)
			}
		} else {
			uptime = "0s"
		}

		diskUsage := ""
		if pvcName != "" {
			diskUsage, err = getDiskUsageFromPVC(ctx, dyn, wb.GetNamespace(), pvcName)
			if err != nil {
				return nil, core.ListWorkbenchesResult{}, fmt.Errorf("failed to get disk usage for workbench %s: %v", name, err)
			}
		}

		workbenchInfo := core.WorkbenchInfo{
			Name:             name,
			User:             user,
			Status:           status,
			ImageDisplayName: imageDisplayName,
			ImageTag:         imageTag,
			HardwareProfile:  hardwareProfile,
			PVCName:          pvcName,
			Namespace:        namespace,
			Uptime:           uptime,
			CPUUsage:         cpuRequest,
			MemoryUsage:      memoryRequest,
			DiskUsage:        diskUsage,
			GPUUsage:         gpuRequest,
		}
		workbenchesInfo = append(workbenchesInfo, workbenchInfo)

	}
	return nil, core.ListWorkbenchesResult{Workbenches: workbenchesInfo}, nil
}

func ListAllWorkbenches(ctx context.Context, req *mcp.CallToolRequest, input core.ListWorkbenchesInput) (*mcp.CallToolResult, core.ListWorkbenchesResult, error) {
	_, workbenches, err := ListWorkbenches(ctx, req, core.ListWorkbenchesInput{Namespace: ""})
	if err != nil {
		return nil, core.ListWorkbenchesResult{}, err
	}
	return nil, core.ListWorkbenchesResult{Workbenches: workbenches.Workbenches}, nil
}

func CreateWorkbench(ctx context.Context, req *mcp.CallToolRequest, input core.CreateWorkbenchInput) (*mcp.CallToolResult, core.DefaultToolOutput, error) {
	dyn, err := GetDynamicClient()
	if err != nil {
		return nil, core.DefaultToolOutput{}, err
	}

	repoURL, gitCommit, imageName, err := GetImageInfo(ctx, input.ImageDisplayName, input.ImageTag)
	if err != nil {
		return nil, core.DefaultToolOutput{}, fmt.Errorf("failed to lookup image info: %v", err)
	}

	if input.PVCName == "" {
		input.PVCName = input.WorkbenchName
		err = createPersistentVolumeClaim(ctx, dyn, input.Namespace, input.PVCName, "10Gi")
		if err != nil {
			return nil, core.DefaultToolOutput{}, fmt.Errorf("failed to create PVC: %v", err)
		}
	}

	notebookArgs := fmt.Sprintf(`--ServerApp.port=8888
                  --ServerApp.token=''
                  --ServerApp.password=''
                  --ServerApp.base_url=/notebook/%s/%s
                  --ServerApp.quit_button=False`, input.Namespace, input.WorkbenchName)

	imageFull := repoURL
	if input.ImageTag != "" {
		imageFull = fmt.Sprintf("%s:%s", repoURL, input.ImageTag)
	}

	var hardwareProfile core.HardwareProfile
	if input.HardwareProfile.HardwareProfileName != "" {
		hardwareProfile = input.HardwareProfile
	} else {
		hardwareProfile = resources.GetDefaultHardwareProfile()
	}

	limits := make(map[string]interface{})
	requests := make(map[string]interface{})
	for _, resource := range hardwareProfile.Resources {
		limits[resource.ResourceIdentifier] = resource.MaxCount
		requests[resource.ResourceIdentifier] = resource.DefaultCount
	}

	notebook := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "kubeflow.org/v1",
			"kind":       "Notebook",
			"metadata": map[string]interface{}{
				"name":      input.WorkbenchName,
				"namespace": input.Namespace,
				"labels": map[string]interface{}{
					"app":                        input.WorkbenchName,
					"opendatahub.io/dashboard":   "true",
					"opendatahub.io/odh-managed": "true",
				},
				"annotations": map[string]interface{}{
					"opendatahub.io/image-display-name":                                input.ImageDisplayName,
					"openshift.io/display-name":                                        input.WorkbenchName,
					"openshift.io/description":                                         "Created via MCP",
					"notebooks.opendatahub.io/inject-auth":                             "true",
					"notebooks.opendatahub.io/last-image-selection":                    fmt.Sprintf("%s:%s", imageName, input.ImageTag),
					"notebooks.opendatahub.io/last-image-version-git-commit-selection": gitCommit,
					"opendatahub.io/hardware-profile-name":                             hardwareProfile.HardwareProfileName,
					"opendatahub.io/hardware-profile-namespace":                        core.GetDefaultNamespace(),
				},
			},
			"spec": map[string]interface{}{
				"template": map[string]interface{}{
					"spec": map[string]interface{}{
						"serviceAccountName": "default",
						"enableServiceLinks": false,
						"containers": []interface{}{
							map[string]interface{}{
								"name":            input.WorkbenchName,
								"image":           imageFull,
								"imagePullPolicy": "Always",
								"workingDir":      "/opt/app-root/src",
								"ports": []interface{}{
									map[string]interface{}{
										"containerPort": 8888,
										"name":          "notebook-port",
										"protocol":      "TCP",
									},
								},
								"env": []interface{}{
									map[string]interface{}{
										"name":  "NOTEBOOK_ARGS",
										"value": notebookArgs,
									},
									map[string]interface{}{
										"name":  "JUPYTER_IMAGE",
										"value": imageFull,
									},
								},
								"resources": map[string]interface{}{
									"limits":   limits,
									"requests": requests,
								},
								"volumeMounts": []interface{}{
									map[string]interface{}{
										"mountPath": "/opt/app-root/src/",
										"name":      input.PVCName,
									},
									map[string]interface{}{
										"mountPath": "/dev/shm",
										"name":      "shm",
									},
								},
							},
						},
						"volumes": []interface{}{
							map[string]interface{}{
								"name": input.PVCName,
								"persistentVolumeClaim": map[string]interface{}{
									"claimName": input.PVCName,
								},
							},
							map[string]interface{}{
								"name": "shm",
								"emptyDir": map[string]interface{}{
									"medium":    "Memory",
									"sizeLimit": "1Gi",
								},
							},
						},
					},
				},
			},
		},
	}

	_, err = dyn.Resource(core.WorkbenchesGVR).Namespace(input.Namespace).Create(ctx, notebook, metav1.CreateOptions{})
	if err != nil {
		return nil, core.DefaultToolOutput{}, fmt.Errorf("failed to create notebook: %v", err)
	}

	return nil, core.DefaultToolOutput{Message: "Workbench was succesfully created!"}, nil
}

func DeleteWorkbench(ctx context.Context, req *mcp.CallToolRequest, input core.DeleteWorkbenchInput) (*mcp.CallToolResult, core.DefaultToolOutput, error) {
	dyn, err := GetDynamicClient()
	if err != nil {
		return nil, core.DefaultToolOutput{}, err
	}

	err = dyn.Resource(core.WorkbenchesGVR).Namespace(input.Namespace).Delete(ctx, input.WorkbenchName, metav1.DeleteOptions{})
	if err != nil {
		return nil, core.DefaultToolOutput{}, fmt.Errorf("failed to delete workbench %s: %v", input.WorkbenchName, err)
	}

	return nil, core.DefaultToolOutput{Message: fmt.Sprintf("Workbench %s was successfully deleted", input.WorkbenchName)}, nil
}

func ChangeWorkbenchStatus(ctx context.Context, req *mcp.CallToolRequest, input core.ChangeWorkbenchStatusInput) (*mcp.CallToolResult, core.DefaultToolOutput, error) {
	dyn, err := GetDynamicClient()
	if err != nil {
		return nil, core.DefaultToolOutput{}, err
	}

	stopped, err := IsWorkbenchStopped(ctx, dyn, input.Namespace, input.WorkbenchName)
	if err != nil {
		return nil, core.DefaultToolOutput{}, err
	}
	if (input.Status == core.Stopped && stopped) || (input.Status == core.Running && !stopped) {
		return nil, core.DefaultToolOutput{Message: fmt.Sprintf("Workbench %s is already %s", input.WorkbenchName, input.Status)}, nil
	}

	patch := map[string]interface{}{}
	annotations := map[string]interface{}{}
	if input.Status == core.Stopped {
		annotations["kubeflow-resource-stopped"] = time.Now().UTC().Format(time.RFC3339)
	} else {
		annotations["kubeflow-resource-stopped"] = nil
	}
	patch["metadata"] = map[string]interface{}{
		"annotations": annotations,
	}

	patchBytes, err := json.Marshal(patch)
	if err != nil {
		return nil, core.DefaultToolOutput{}, fmt.Errorf("failed to marshal patch: %v", err)
	}

	_, err = dyn.Resource(core.WorkbenchesGVR).Namespace(input.Namespace).Patch(
		ctx,
		input.WorkbenchName,
		k8stypes.MergePatchType,
		patchBytes,
		metav1.PatchOptions{},
	)
	if err != nil {
		return nil, core.DefaultToolOutput{}, fmt.Errorf("failed to %s workbench %s: %v", input.Status, input.WorkbenchName, err)
	}

	return nil, core.DefaultToolOutput{Message: fmt.Sprintf("Workbench %s is %s", input.WorkbenchName, input.Status)}, nil
}

func IsWorkbenchStopped(ctx context.Context, dyn dynamic.Interface, namespace, workbenchName string) (bool, error) {
	current, err := dyn.Resource(core.WorkbenchesGVR).Namespace(namespace).Get(ctx, workbenchName, metav1.GetOptions{})
	if err != nil {
		return false, fmt.Errorf("failed to get workbench %s: %v", workbenchName, err)
	}
	currentAnnotations := current.GetAnnotations()
	currentStopped := false
	if currentAnnotations != nil {
		if _, ok := currentAnnotations["kubeflow-resource-stopped"]; ok {
			currentStopped = true
		}
	}
	return currentStopped, nil
}

// GetImageInfo retrieves image information from display name and version
// Returns: repoURL, gitCommit, imageName, error
func GetImageInfo(ctx context.Context, displayName, version string) (string, string, string, error) {
	dyn, err := GetDynamicClient()
	if err != nil {
		return "", "", "", err
	}

	images, err := dyn.Resource(core.ImagesGVR).Namespace(core.GetDefaultNamespace()).List(ctx, metav1.ListOptions{
		LabelSelector: "opendatahub.io/notebook-image=true",
	})
	if err != nil {
		return "", "", "", fmt.Errorf("failed to list images: %v", err)
	}

	for _, image := range images.Items {
		annotations := image.GetAnnotations()
		if annotations["opendatahub.io/notebook-image-name"] == displayName {
			repoURL, found, err := unstructured.NestedString(image.Object, "status", "dockerImageRepository")
			if !found || err != nil {
				repoURL = "URL not available"
			}
			imageName := image.GetName()

			tagsRaw, _, _ := unstructured.NestedSlice(image.Object, "spec", "tags")
			for _, t := range tagsRaw {
				tagMap, ok := t.(map[string]interface{})
				if !ok {
					continue
				}
				tagName, _ := tagMap["name"].(string)
				if tagName == version {
					tagAnnotations, _, _ := unstructured.NestedStringMap(tagMap, "annotations")
					return repoURL, tagAnnotations["opendatahub.io/notebook-build-commit"], imageName, nil
				}
			}
		}
	}
	return "", "", "", fmt.Errorf("image not found: %s:%s", displayName, version)
}

func getUptimeFromWorkbench(workbenchName string, nameSpace string) (string, error) {
	ctx := context.Background()
	clientset, err := GetClientSet()
	if err != nil {
		return "", fmt.Errorf("failed to get client set: %v", err)
	}

	labelSelector := fmt.Sprintf("notebook-name=%s", workbenchName)
	pods, err := clientset.CoreV1().Pods(nameSpace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return "", fmt.Errorf("failed to list pods for workbench %s in namespace %s: %v", workbenchName, nameSpace, err)
	}

	if len(pods.Items) == 0 {
		return "", nil
	}

	pod := pods.Items[0]
	if pod.Status.StartTime == nil {
		return "", nil
	}
	return pod.Status.StartTime.String(), nil
}
