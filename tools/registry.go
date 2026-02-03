package tools

import "github.com/modelcontextprotocol/go-sdk/mcp"

func RegisterAllTools(server *mcp.Server) {
	registerPodTools(server)
	registerNamespaceTools(server)
	registerWorkbenchTools(server)
	registerImageTools(server)
	registerHardwareProfileTools(server)
	registerStorageTools(server)
	registerResourceConsumptionTools(server)
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
		Name:        "List Workbenches",
		Description: "list the workbenches in a given project namespace",
	}, ListWorkbenches)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "List All Workbenches",
		Description: "list the workbenches across all namespaces",
	}, ListAllWorkbenches)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "Create Workbench",
		Description: "create a new workbench with given name, image and image URL in a given project namespace",
	}, CreateWorkbench)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "Delete Workbench",
		Description: "delete a workbench with given name in a given project namespace",
	}, DeleteWorkbench)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "Change Workbench Status",
		Description: "change the status of a workbench with given name in a given project namespace",
	}, ChangeWorkbenchStatus)
}

func registerImageTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "List Images",
		Description: "list all available notebook images",
	}, ListImages)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "Create Custom Image",
		Description: "create a new custom notebook image",
	}, CreateCustomImage)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "Delete Image",
		Description: "delete an image with given name",
	}, DeleteImage)
}

func registerHardwareProfileTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "List Hardware Profiles",
		Description: "list the hardware profiles in a given project namespace",
	}, ListHardwareProfiles)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "Create Hardware Profile",
		Description: "create a hardware profile with given name, description and resources",
	}, CreateHardwareProfile)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "Delete Hardware Profile",
		Description: "delete a hardware profile with given name",
	}, DeleteHardwareProfile)
}

func registerStorageTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "List PVCs",
		Description: "list the persistent volume claims in a given project namespace",
	}, ListPVCs)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "Create PVC",
		Description: "create a persistent volume claim with given name and size in a given project namespace",
	}, CreatePVCTool)
}

func registerResourceConsumptionTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "List Resource Consumption Per Workbench",
		Description: "list the resource consumption per workbench",
	}, ListResourceConsumptionPerWorkbench)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "List Resource Consumption Per Namespace",
		Description: "list the resource consumption per namespace",
	}, ListResourceConsumptionPerNamespace)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "List Resource Consumption Per User",
		Description: "list the resource consumption per user",
	}, ListResourceConsumptionPerUser)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "List Resource Consumption Per Cluster",
		Description: "list the resource consumption per cluster",
	}, ListResourceConsumptionPerCluster)
}
