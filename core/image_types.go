package core

type ListImagesOutput struct {
	Images string `json:"images" jsonschema_description:"the list of images"`
}

type CreateCustomImageInput struct {
	ImageLocation    string `json:"imageLocation" jsonschema_description:"the location of the image"`
	ImageName        string `json:"imageName" jsonschema_description:"the name of the image"`
	ImageDescription string `json:"imageDescription" jsonschema_description:"the description of the image"`
}

type UpdateImageInput struct {
	ImageName        string `json:"imageName" jsonschema_description:"the name of the image"`
	ImageDescription string `json:"imageDescription" jsonschema_description:"the description of the image"`
}

type DeleteImageInput struct {
	ImageName string `json:"imageName" jsonschema_description:"the name of the image"`
}
