package tools

import "github.com/modelcontextprotocol/go-sdk/mcp"

func RegisterWriteTools(server *mcp.Server) {
	registerPodTools(server)
	registerNamespaceTools(server)
	registerWorkbenchTools(server)
	registerImageTools(server)
	registerHardwareProfileTools(server)
	registerStorageTools(server)
}

func RegisterReadOnlyTools(server *mcp.Server) {
	registerWorkbenchListingTools(server)
	registerImageListingTools(server)
	registerHardwareProfileListingTools(server)
	registerStorageListingTools(server)
	registerResourceConsumptionListingTools(server)
}

func registerPodTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "List Pods",
		Description: "list the pods in a namespace",
	}, ListPods)
}

func registerNamespaceTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "List Namespaces",
		Description: "list all namespaces in the cluster",
	}, ListNamespaces)
}

func registerWorkbenchTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "Create Workbench",
		Description: "create a new workbench with given name, image and image URL in a given project namespace",
	}, CreateWorkbench)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "Update Workbench",
		Description: "update a workbench with given name, image and image URL in a given project namespace",
	}, UpdateWorkbench)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "Delete Workbench",
		Description: "delete a workbench with given name in a given project namespace",
	}, DeleteWorkbench)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "Change Workbench Status",
		Description: "change the status of a workbench with given name in a given project namespace",
	}, ChangeWorkbenchStatus)
}

func registerWorkbenchListingTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "List Workbenches",
		Description: "list the workbenches in a given project namespace",
	}, ListWorkbenches)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "List All Workbenches",
		Description: "list the workbenches across all namespaces",
	}, ListAllWorkbenches)
}

func registerImageTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "Create Custom Image",
		Description: "create a new custom notebook image",
	}, CreateCustomImage)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "Update Image",
		Description: "update an image with given name, description and location",
	}, UpdateImage)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "Delete Image",
		Description: "delete an image with given name",
	}, DeleteImage)
}

func registerImageListingTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "List Images",
		Description: "list all available notebook images",
	}, ListImages)
}

func registerHardwareProfileTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "Create Hardware Profile",
		Description: "create a hardware profile with given name and resources",
	}, CreateHardwareProfile)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "Update Hardware Profile",
		Description: "update a hardware profile with given name and resources",
	}, UpdateHardwareProfile)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "Delete Hardware Profile",
		Description: "delete a hardware profile with given name",
	}, DeleteHardwareProfile)
}

func registerHardwareProfileListingTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "List Hardware Profiles",
		Description: "list the hardware profiles in a given project namespace",
	}, ListHardwareProfiles)
}

func registerStorageTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "Create PVC",
		Description: "create a persistent volume claim with given name and size in a given project namespace",
	}, CreatePVC)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "Update PVC",
		Description: "update a persistent volume claim with given name and size in a given project namespace",
	}, UpdatePVC)
}

func registerStorageListingTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "List PVCs",
		Description: "list the persistent volume claims in a given project namespace",
	}, ListPVCs)
}

func registerResourceConsumptionListingTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "List Resource Consumption Per Workbench",
		Description: "list the resource consumption of given workbench in a given namespace",
	}, ListResourceConsumptionPerWorkbench)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "List Resource Consumption Per Namespace",
		Description: "list the resource consumption of all workbenches in a given namespace",
	}, ListResourceConsumptionPerNamespace)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "List Resource Consumption Per User",
		Description: "list the resource consumption of all workbenches of a given user",
	}, ListResourceConsumptionPerUser)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "List Resource Consumption Per Cluster",
		Description: "list the resource consumption of all workbenches in the cluster",
	}, ListResourceConsumptionPerCluster)
}
