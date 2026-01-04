// ... existing code ...

func TestCreateCustomImage(t *testing.T) {
	orig := getDynamicClient
	defer func() { getDynamicClient = orig }()

	scheme := runtime.NewScheme()
	client := dynamicfake.NewSimpleDynamicClient(scheme)

	getDynamicClient = func() (dynamic.Interface, error) {
		return client, nil
	}

	input := CreateCustomImageInput{
		ImageName:        "test-img",
		ImageDescription: "test description",
		ImageLocation:    "ghcr.io/test/image:latest",
	}

	_, _, err := CreateCustomImage(context.Background(), nil, input)
	if err != nil {
		t.Fatalf("CreateCustomImage returned error: %v", err)
	}

	// Verify the image stream was created
	gvr := schema.GroupVersionResource{Group: "image.openshift.io", Version: "v1", Resource: "imagestreams"}
	imageStream, err := client.Resource(gvr).Namespace("redhat-ods-applications").Get(context.Background(), "test-img", metav1.GetOptions{})
	if err != nil {
		t.Fatalf("failed to get image stream: %v", err)
	}

	if imageStream == nil {
		t.Fatal("expected image stream to be created, got nil")
	}

	// Verify some fields
	annotations := imageStream.GetAnnotations()
	if annotations["opendatahub.io/notebook-image-desc"] != "test description" {
		t.Errorf("expected description 'test description', got %q", annotations["opendatahub.io/notebook-image-desc"])
	}
	if annotations["opendatahub.io/notebook-image-url"] != "ghcr.io/test/image:latest" {
		t.Errorf("expected url 'ghcr.io/test/image:latest', got %q", annotations["opendatahub.io/notebook-image-url"])
	}
}
