package tools

import (
	"context"
	"fmt"

	core "github.com/amaly/mcp-server-rhoai/core"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
)

func CreateWorkbench(ctx context.Context, req *mcp.CallToolRequest, input core.CreateWorkbenchInput) (*mcp.CallToolResult, core.DefaultToolOutput, error) {
	dyn, err := GetDynamicClient()
	if err != nil {
		return nil, core.DefaultToolOutput{}, err
	}

	repoURL, gitCommit, imageName, err := GetImageInfo(ctx, input.ImageDisplayName, input.ImageTag)
	if err != nil {
		return nil, core.DefaultToolOutput{}, fmt.Errorf("failed to lookup image info: %v", err)
	}

	err = createPersistentVolumeClaim(ctx, dyn, input.Namespace, input.WorkbenchName, "10Gi")
	if err != nil {
		return nil, core.DefaultToolOutput{}, fmt.Errorf("failed to create PVC: %v", err)
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
					"opendatahub.io/hardware-profile-name":                             "default-profile",
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
									"limits": map[string]interface{}{
										"cpu":    "2",
										"memory": "4Gi",
									},
									"requests": map[string]interface{}{
										"cpu":    "2",
										"memory": "4Gi",
									},
								},
								"volumeMounts": []interface{}{
									map[string]interface{}{
										"mountPath": "/opt/app-root/src/",
										"name":      "storage-volume",
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
								"name": "storage-volume",
								"persistentVolumeClaim": map[string]interface{}{
									"claimName": input.WorkbenchName,
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

// from display name and version, gets url, git commit and image name
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
