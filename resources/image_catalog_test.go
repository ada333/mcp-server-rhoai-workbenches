package resources

import (
	"context"
	"testing"

	core "github.com/amaly/mcp-server-rhoai/core"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	dynamicfake "k8s.io/client-go/dynamic/fake"
)

func NewUnstructuredImageForTest(name, displayName, repoURL string, versions []string) *unstructured.Unstructured {
	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(core.ImagesGVR.GroupVersion().WithKind("ImageStream"))
	u.SetName(name)
	u.SetNamespace("redhat-ods-applications")
	u.SetLabels(map[string]string{
		"opendatahub.io/notebook-image": "true",
	})
	u.SetAnnotations(map[string]string{
		"opendatahub.io/notebook-image-name": displayName,
	})

	if repoURL != "" {
		unstructured.SetNestedField(u.Object, repoURL, "status", "dockerImageRepository")
	}

	tags := make([]interface{}, len(versions))
	for i, v := range versions {
		tags[i] = map[string]interface{}{"name": v}
	}
	unstructured.SetNestedSlice(u.Object, tags, "spec", "tags")

	return u
}

func TestGetImages(t *testing.T) {
	orig := GetDynamicClient
	defer func() { GetDynamicClient = orig }()

	scheme := runtime.NewScheme()

	image1 := NewUnstructuredImageForTest("img1", "PyTorch", "quay.io/modh/pytorch", []string{"v1", "v2"})
	image2 := NewUnstructuredImageForTest("img2", "TensorFlow", "quay.io/modh/tensorflow", []string{"latest"})

	client := dynamicfake.NewSimpleDynamicClient(scheme, image1, image2)

	GetDynamicClient = func() (dynamic.Interface, error) {
		return client, nil
	}

	images, err := GetImages(context.Background())
	if err != nil {
		t.Fatalf("GetImages returned error: %v", err)
	}

	foundPyTorch := false
	foundTensorFlow := false

	for _, img := range images {
		if img.Annotations["opendatahub.io/notebook-image-name"] == "PyTorch" {
			foundPyTorch = true
			if img.URL != "quay.io/modh/pytorch" {
				t.Errorf("expected PyTorch URL quay.io/modh/pytorch, got: %q", img.URL)
			}
			// Check versions
			hasV1 := false
			hasV2 := false
			for _, v := range img.Versions {
				if v.Name == "v1" {
					hasV1 = true
				}
				if v.Name == "v2" {
					hasV2 = true
				}
			}
			if !hasV1 || !hasV2 {
				t.Errorf("expected versions v1 and v2 for PyTorch, got: %v", img.Versions)
			}
		}
		if img.Annotations["opendatahub.io/notebook-image-name"] == "TensorFlow" {
			foundTensorFlow = true
			if img.URL != "quay.io/modh/tensorflow" {
				t.Errorf("expected TensorFlow URL quay.io/modh/tensorflow, got: %q", img.URL)
			}
		}
	}

	if !foundPyTorch {
		t.Errorf("expected PyTorch image in results")
	}
	if !foundTensorFlow {
		t.Errorf("expected TensorFlow image in results")
	}
}
