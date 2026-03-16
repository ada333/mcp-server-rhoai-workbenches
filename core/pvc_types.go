package core

type PVCInput struct {
	Namespace  string `json:"namespace" jsonschema_description:"the namespace of the PVC"`
	PVCName    string `json:"pvcName" jsonschema_description:"the name of the PVC to update"`
	NewPVCName string `json:"newPVCName,omitempty" jsonschema_description:"the new name to rename the PVC to (optional)"`
	Size       string `json:"size" jsonschema_description:"the size of the PVC"`
}

type DeletePVCInput struct {
	Namespace string `json:"namespace" jsonschema_description:"the namespace of the PVC"`
	PVCName   string `json:"pvcName" jsonschema_description:"the name of the PVC"`
}

type PVCsOutput struct {
	PVCs string `json:"pvcs" jsonschema_description:"the list of PVCs"`
}

type ListPVCsInput struct {
	Namespace string `json:"namespace" jsonschema_description:"the namespace of the PVC"`
}
